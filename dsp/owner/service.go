package owner

import (
	"github.com/YasiruR/connector/core/errors"
	"github.com/YasiruR/connector/core/pkg"
	"github.com/YasiruR/connector/core/protocols/dcat"
	"github.com/YasiruR/connector/core/protocols/odrl"
	"github.com/YasiruR/connector/core/stores"
)

type Owner struct {
	host         string
	policyStore  stores.Policy
	datasetStore stores.Dataset
	urn          pkg.URN
	log          pkg.Log
}

func New(ps stores.Policy, ds stores.Dataset, urn pkg.URN, log pkg.Log) *Owner {
	return &Owner{
		host:         `http://localhost:`,
		policyStore:  ps,
		datasetStore: ds,
		urn:          urn,
		log:          log,
	}
}

func (o *Owner) CreatePolicy(t odrl.Target, permissions, prohibitions []odrl.Rule) (policyId string, err error) {
	policyId, err = o.urn.New()
	if err != nil {
		return ``, errors.URNFailed(`policyId`, `New`, err)
	}

	// handle other policy types
	ofr := odrl.Offer{
		Id:           policyId,
		Type:         odrl.TypeOffer,
		Target:       t,
		Assigner:     odrl.Assigner(o.host),
		Permissions:  permissions,
		Prohibitions: prohibitions,
	}

	o.policyStore.SetOffer(policyId, ofr)
	o.log.Info("created and stored the new policy", ofr)
	return policyId, nil
}

// CreateDataset currently supports only one data distribution per a dataset
func (o *Owner) CreateDataset(title string, descriptions, keywords, endpoints, policyIds []string) (datasetId string, err error) {
	// construct policies (handle policies than offers later)
	var policies []odrl.Offer
	for _, pId := range policyIds {
		ofr, err := o.policyStore.GetOffer(pId)
		if err != nil {
			return ``, errors.StoreFailed(stores.TypePolicy, `GetOffer`, err)
		}

		policies = append(policies, ofr)
	}

	// construct data distribution
	var svcList []dcat.AccessService
	for _, e := range endpoints {
		accessServiceId, err := o.urn.New()
		if err != nil {
			return ``, errors.URNFailed(`accessServiceId`, `New`, err)
		}

		svcList = append(svcList, dcat.AccessService{
			ID:              accessServiceId,
			Type:            "", // add type
			DcatEndpointURL: e,
		})
	}

	dist := dcat.Distribution{
		Type:              dcat.TypeDistribution,
		DctFormat:         "", // add format (e.g. dspace:s3+push)
		DcatAccessService: svcList,
	}

	// construct and store final dataset
	datasetId, err = o.urn.New()
	if err != nil {
		return ``, errors.URNFailed(`datasetId`, `New`, err)
	}

	var descs []dcat.Description
	for _, desc := range descriptions {
		descs = append(descs, dcat.Description{
			Value:    desc,
			Language: dcat.LanguageEnglish, // support other languages
		})
	}

	var kws []dcat.Keyword
	for _, kw := range keywords {
		kws = append(kws, dcat.Keyword(kw))
	}

	ds := dcat.Dataset{
		ID:               datasetId,
		Type:             dcat.TypeDataset,
		DctTitle:         title,
		DctDescription:   descs,
		DcatKeyword:      kws,
		OdrlHasPolicy:    policies,
		DcatDistribution: []dcat.Distribution{dist}, // support more distributions
	}

	o.datasetStore.Set(datasetId, ds)
	o.log.Info("new dataset was created and stored", ds)
	return datasetId, nil
}
