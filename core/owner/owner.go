package owner

import (
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/boot"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/models/dcat"
	"github.com/YasiruR/connector/domain/models/odrl"
	"github.com/YasiruR/connector/domain/pkg"
	"github.com/YasiruR/connector/domain/stores"
)

type Service struct {
	assignerId  string
	catalog     stores.Catalog
	policyStore stores.Policy
	urn         pkg.URNService
	log         pkg.Log
}

func New(cfg boot.Config, stores domain.Stores, plugins domain.Plugins) *Service {
	return &Service{
		assignerId:  cfg.DataSpace.AssignerId, // can we assign participant ID from config to assigner?
		policyStore: stores.Policy,
		catalog:     stores.Catalog,
		urn:         plugins.URNService,
		log:         plugins.Log,
	}
}

func (s *Service) CreatePolicy(target string, permissions, prohibitions []odrl.Rule) (ofrId string, err error) {
	ofrId, err = s.urn.NewURN()
	if err != nil {
		return ``, errors.URNFailed(`offerId`, `NewURN`, err)
	}

	// handle other policy types
	ofr := odrl.Offer{
		Id:           ofrId,
		Type:         odrl.TypeOffer,
		Target:       odrl.Target(target),
		Assigner:     odrl.Assigner(s.assignerId),
		Permissions:  permissions,
		Prohibitions: prohibitions,
	}

	s.policyStore.SetOffer(ofrId, ofr)
	s.log.Trace("created and stored a new offer", ofr)
	return ofrId, nil
}

// CreateDataset currently supports only one data distribution per a dataset
func (s *Service) CreateDataset(title, format string, descriptions, keywords, endpoints, offerIds []string) (dsId string, err error) {
	// construct policies (handle policies than offers later)
	var policies []odrl.Offer
	for _, pId := range offerIds {
		ofr, err := s.policyStore.GetOffer(pId)
		if err != nil {
			return ``, errors.StoreFailed(stores.TypePolicy, `GetOffer`, err)
		}

		policies = append(policies, ofr)
	}

	// construct data distribution
	var svcList []dcat.AccessService
	for _, e := range endpoints {
		accessServiceId, err := s.urn.NewURN()
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
		DctFormat:         format,
		DcatAccessService: svcList,
	}

	// construct and store final dataset
	dsId, err = s.urn.NewURN()
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
		ID:               dsId,
		Type:             dcat.TypeDataset,
		DctTitle:         title,
		DctDescription:   descs,
		DcatKeyword:      kws,
		OdrlHasPolicy:    policies,
		DcatDistribution: []dcat.Distribution{dist}, // support more distributions
	}

	s.catalog.AddDataset(dsId, ds)
	s.log.Trace("created and stored a new dataset", ds)
	return dsId, nil
}
