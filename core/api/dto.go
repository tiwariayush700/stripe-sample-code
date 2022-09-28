package api

// PaymentRequest should never accept amount from client side.
//One Should just accept item_id and calculate the amount based on the item_id
type PaymentRequest struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

// PaymentResponse model which can accommodate fields from multiple gateways
type PaymentResponse struct {
	ID           string                 `json:"id"`
	PublicKey    string                 `json:"publicKey"`
	ClientSecret string                 `json:"clientSecret"`
	Extras       map[string]interface{} `json:"extras"`
}
