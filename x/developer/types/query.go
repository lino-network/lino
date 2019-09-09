package types

type QueryResultIDABalance struct {
	Amount   string `json:"amount"`
	Unauthed bool   `json:"unauthed"`
}
