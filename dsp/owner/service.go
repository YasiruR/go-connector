package owner

import (
	"fmt"
	"github.com/YasiruR/connector/core/pkg"
	"github.com/YasiruR/connector/core/protocols/dcat"
	"github.com/YasiruR/connector/core/protocols/odrl"
	"github.com/YasiruR/connector/core/stores"
	"github.com/YasiruR/connector/pkg/urn"
)

type Owner struct {
	host         string
	policyStore  stores.Policy
	datasetStore stores.Dataset
	urn          pkg.URN
	log          pkg.Log
}

func New(ps stores.Policy, ds stores.Dataset, log pkg.Log) *Owner {
	return &Owner{
		host:         `http://localhost:`,
		policyStore:  ps,
		datasetStore: ds,
		urn:          urn.NewGenerator(), // pass this
		log:          log,
	}
}

func (o *Owner) CreatePolicy(t odrl.Target, permissions, prohibitions []odrl.Rule) (policyId string, err error) {
	policyId, err = o.urn.New()
	if err != nil {
		return ``, fmt.Errorf("generate new URN failed - %w", err)
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
	o.log.Info("new policy was created and stored", ofr)
	return policyId, nil
}

// CreateDataset currently supports only one data distribution per a dataset
func (o *Owner) CreateDataset(title string, descriptions, keywords, endpoints, policyIds []string) (datasetId string, err error) {
	// construct policies (handle policies than offers later)
	var policies []odrl.Offer
	for _, pId := range policyIds {
		ofr, err := o.policyStore.GetOffer(pId)
		if err != nil {
			return ``, fmt.Errorf("get offer failed - %w", err)
		}

		policies = append(policies, ofr)
	}

	// construct data distribution
	var svcList []dcat.AccessService
	for _, e := range endpoints {
		accessServiceId, err := o.urn.New()
		if err != nil {
			return ``, fmt.Errorf("generate new URN for access service failed - %w", err)
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
		return ``, fmt.Errorf("generate new URN for dataset failed - %w", err)
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
