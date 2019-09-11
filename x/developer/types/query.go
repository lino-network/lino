package types

type QueryResultIDABalance struct {
	Amount   string `json:"appida_amount"`
	Unauthed bool   `json:"unauthed"`
}
