package api

type PaymentRequest struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

// PaymentResponse model which can accommodate fields from multiple gateways
type PaymentResponse struct {
	PublicKey    string                 `json:"publicKey"`
	ClientSecret string                 `json:"clientSecret"`
	Extras       map[string]interface{} `json:"extras"`
}
