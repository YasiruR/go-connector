package owner

import (
	defaultErr "errors"
	"fmt"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/boot"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/models/dcat"
	"github.com/YasiruR/connector/domain/models/odrl"
	"github.com/YasiruR/connector/domain/pkg"
	"github.com/YasiruR/connector/domain/stores"
)

type Service struct {
	assignerId string
	catalog    stores.ProviderCatalog
	ofrStore   stores.OfferStore
	urn        pkg.URNService
	log        pkg.Log
}

func New(cfg boot.Config, stores domain.Stores, plugins domain.Plugins) *Service {
	return &Service{
		assignerId: cfg.DataSpace.AssignerId, // can we assign participant ID from config to assigner?
		ofrStore:   stores.OfferStore,
		catalog:    stores.ProviderCatalog,
		urn:        plugins.URNService,
		log:        plugins.Log,
	}
}

func (s *Service) CreatePolicy(target string, permissions, prohibitions []odrl.Rule) (ofrId string, err error) {
	ofrId, err = s.urn.NewURN()
	if err != nil {
		return ``, errors.PkgError(pkg.TypeURN, `NewURN`, err, `offer id`)
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

	s.ofrStore.AddOffer(ofrId, ofr)
	s.log.Trace("created and stored a new offer", ofr)
	return ofrId, nil
}

// CreateDataset currently supports only one data distribution per a dataset
func (s *Service) CreateDataset(title, format string, descriptions, keywords, endpoints, offerIds []string) (dsId string, err error) {
	// construct policies
	var ofrs []odrl.Offer
	for _, ofrId := range offerIds {
		ofr, err := s.ofrStore.Offer(ofrId)
		if err != nil {
			if defaultErr.Is(err, stores.TypeInvalidKey) {
				return ``, errors.Client(errors.InvalidKey(stores.TypeOffer, `offer id`, err))
			}
			return ``, errors.StoreFailed(stores.TypeOffer, `Offer`, err)
		}

		ofr.Target = `` // since associated dataset id represents the target implicitly
		ofrs = append(ofrs, ofr)
	}

	// construct data distribution
	var svcList []dcat.AccessService
	for _, e := range endpoints {
		accessServiceId, err := s.urn.NewURN()
		if err != nil {
			return ``, errors.PkgError(pkg.TypeURN, `NewURN`, err, `access service id`)
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
		return ``, errors.PkgError(pkg.TypeURN, `NewURN`, err, `dataset id`)
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

	if len(ofrs) == 0 {
		s.log.Trace(fmt.Sprintf(`no policy offers associated with the dataset 
			(dataset Id: %s, offer Ids: %s)`, dsId, offerIds))
	}

	ds := dcat.Dataset{
		ID:               dsId,
		Type:             dcat.TypeDataset,
		DctTitle:         title,
		DctDescription:   descs,
		DcatKeyword:      kws,
		OdrlHasPolicy:    ofrs,
		DcatDistribution: []dcat.Distribution{dist}, // support more distributions
	}

	s.catalog.AddDataset(dsId, ds)
	s.log.Trace("created and stored a new dataset", ds)
	return dsId, nil
}
