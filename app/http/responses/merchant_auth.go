package responses

type MerchantTokenSyncResponse struct {
	MerchantBearerToken  string `json:"merchant_bearer_token"`
	ExternalKey          string `json:"external_key"`
	MerchantClientID     string `json:"merchant_client_id"`
	MerchantClientSecret string `json:"merchant_client_secret"`
}
