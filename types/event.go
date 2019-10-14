package types

// Event - event executed in app.go
type Event interface{}

// Minute -> TimeEventList
type TimeEventList struct {
	Events []Event `json:"events"`
}


// ReturnCoinEvent - return a certain amount of coin to an account
type ReturnCoinEvent struct {
	Username   AccountKey         `json:"username"`
	Amount     Coin               `json:"amount"`
	ReturnType TransferDetailType `json:"return_type"`
}


