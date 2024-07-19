package owner

import (
	"fmt"
	"github.com/YasiruR/connector/core/models/odrl"
	"github.com/YasiruR/connector/core/pkg"
	"github.com/YasiruR/connector/core/stores"
	"github.com/YasiruR/connector/pkg/urn"
	"strconv"
)

type Owner struct {
	callbackAddr string
	policyStore  stores.Policy
	urn          pkg.URN
	log          pkg.Log
}

func New(port int, ps stores.Policy, log pkg.Log) *Owner {
	return &Owner{
		callbackAddr: `http://localhost:` + strconv.Itoa(port),
		policyStore:  ps,
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
		Assigner:     odrl.Assigner(o.callbackAddr),
		Permissions:  permissions,
		Prohibitions: prohibitions,
	}

	o.policyStore.SetOffer(policyId, ofr)
	o.log.Trace("new policy was created and stored", ofr)
	return policyId, nil
}

func (o *Owner) CreateContractDefinition() {

}
