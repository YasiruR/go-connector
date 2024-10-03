package daps

type MetadataResponse struct {
	Keys []struct {
		Kty string `json:"kty"`
		N   string `json:"n"`
		E   string `json:"e"`
		Kid string `json:"kid"`
		Use string `json:"use"`
	} `json:"keys"`
}

type Response struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}
