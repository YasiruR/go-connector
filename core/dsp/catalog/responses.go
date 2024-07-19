package catalog

// Data models required for the Catalog Protocol as specified by
// https://docs.internationaldataspaces.org/ids-knowledgebase/v/dataspace-protocol/catalog/catalog.protocol

type Response struct {
	Context        string `json:"@context"`
	ID             string `json:"@id"`
	Type           string `json:"@type"`
	DctTitle       string `json:"dct:title"`
	DctDescription []struct {
		Value    string `json:"@value"`
		Language string `json:"@language"`
	} `json:"dct:description"`
	DspaceParticipantID string   `json:"dspace:participantId"`
	DcatKeyword         []string `json:"dcat:keyword"`
	DcatService         []struct {
		ID                      string `json:"@id"`
		Type                    string `json:"@type"`
		DcatEndpointDescription string `json:"dcat:endpointDescription"`
		DcatEndpointURL         string `json:"dcat:endpointURL"`
	} `json:"dcat:service"`
	DcatDataset []struct {
		ID             string `json:"@id"`
		Type           string `json:"@type"`
		DctTitle       string `json:"dct:title"`
		DctDescription []struct {
			Value    string `json:"@value"`
			Language string `json:"@language"`
		} `json:"dct:description"`
		DcatKeyword   []string `json:"dcat:keyword"`
		OdrlHasPolicy []struct {
			ID             string `json:"@id"`
			Type           string `json:"@type"`
			OdrlAssigner   string `json:"odrl:assigner"`
			OdrlPermission []struct {
				OdrlAction     string `json:"odrl:action"`
				OdrlConstraint []struct {
					OdrlLeftOperand  string `json:"odrl:leftOperand"`
					OdrlOperator     string `json:"odrl:operator"`
					OdrlRightOperand string `json:"odrl:rightOperand"`
				} `json:"odrl:constraint"`
				OdrlDuty struct {
					OdrlAction string `json:"odrl:action"`
				} `json:"odrl:duty"`
			} `json:"odrl:permission"`
		} `json:"odrl:hasPolicy"`
		DcatDistribution []struct {
			Type              string `json:"@type"`
			DctFormat         string `json:"dct:format"`
			DcatAccessService []struct {
				ID              string `json:"@id"`
				Type            string `json:"@type"`
				DcatEndpointURL string `json:"dcat:endpointURL"`
			} `json:"dcat:accessService"`
		} `json:"dcat:distribution"`
	} `json:"dcat:dataset"`
}

type DatasetResponse struct {
	Context        string `json:"@context"`
	ID             string `json:"@id"`
	Type           string `json:"@type"`
	DctTitle       string `json:"dct:title"`
	DctDescription []struct {
		Value    string `json:"@value"`
		Language string `json:"@language"`
	} `json:"dct:description"`
	DcatKeyword   []string `json:"dcat:keyword"`
	OdrlHasPolicy []struct {
		Type           string `json:"@type"`
		ID             string `json:"@id"`
		OdrlAssigner   string `json:"odrl:assigner"`
		OdrlPermission []struct {
			OdrlAction     string `json:"odrl:action"`
			OdrlConstraint []struct {
				OdrlLeftOperand  string `json:"odrl:leftOperand"`
				OdrlRightOperand string `json:"odrl:rightOperand"`
				OdrlOperator     string `json:"odrl:operator"`
			} `json:"odrl:constraint"`
		} `json:"odrl:permission"`
	} `json:"odrl:hasPolicy"`
	DcatDistribution []struct {
		Type      string `json:"@type"`
		DctFormat struct {
			ID string `json:"@id"`
		} `json:"dct:format"`
		DcatAccessService []struct {
			ID              string `json:"@id"`
			DcatEndpointURL string `json:"dcat:endpointURL"`
		} `json:"dcat:accessService"`
	} `json:"dcat:distribution"`
}

type ErrorResponse struct {
	Context      string `json:"@context"`
	Type         string `json:"@type"`
	DspaceCode   string `json:"dspace:code"`
	DspaceReason []struct {
		Value    string `json:"@value"`
		Language string `json:"@language"`
	} `json:"dspace:reason"`
}
