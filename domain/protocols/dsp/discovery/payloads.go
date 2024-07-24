package discovery

// .well-known endpoint

type Response struct {
	Ctx      string            `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Versions []ProtocolVersion `json:"protocolVersions"`
}

type ProtocolVersion struct {
	Version string `json:"version"`
	Path    string `json:"path"`
}
