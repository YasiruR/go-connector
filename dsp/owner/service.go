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
	us           pkg.URNService
	log          pkg.Log
}

func New(ps stores.Policy, ds stores.Dataset, urn pkg.URNService, log pkg.Log) *Owner {
	return &Owner{
		host:         `http://localhost:`,
		policyStore:  ps,
		datasetStore: ds,
		us:           urn,
		log:          log,
	}
}

func (o *Owner) CreatePolicy(target string, permissions, prohibitions []odrl.Rule) (policyId string, err error) {
	policyId, err = o.us.NewURN()
	if err != nil {
		return ``, errors.URNFailed(`policyId`, `NewURN`, err)
	}

	// handle other policy types
	ofr := odrl.Offer{
		Id:           policyId,
		Type:         odrl.TypeOffer,
		Target:       odrl.Target(target),
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
		ofr, err := o.policyStore.Offer(pId)
		if err != nil {
			return ``, errors.StoreFailed(stores.TypePolicy, `Offer`, err)
		}

		policies = append(policies, ofr)
	}

	// construct data distribution
	var svcList []dcat.AccessService
	for _, e := range endpoints {
		accessServiceId, err := o.us.NewURN()
		if err != nil {
			return ``, errors.URNFailed(`accessServiceId`, `NewURN`, err)
		}

		svcList = append(svcList, dcat.AccessService{
			ID:          accessServiceId,
			Type:        dcat.TypeDataService,
			EndpointURL: e,
		})
	}

	dist := dcat.Distribution{
		Type:              dcat.TypeDistribution,
		DctFormat:         "", // add format (e.g. dspace:s3+push)
		DcatAccessService: svcList,
	}

	// construct and store final dataset
	datasetId, err = o.us.NewURN()
	if err != nil {
		return ``, errors.URNFailed(`datasetId`, `NewURN`, err)
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
