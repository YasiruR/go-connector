package catalog

type Request struct {
	Context      string   `json:"@context"`
	Type         string   `json:"@type"`
	DspaceFilter []string `json:"dspace:filter"` // optional
}
