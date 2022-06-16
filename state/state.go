package state

type AppState struct {
	Marketplace    string
	LoggedIn       bool
	Email          string
	AccountID      string
	SubscriptionID string
	ContractID     string
}

var Status = AppState{
	Marketplace:    "canonical-ua",
	LoggedIn:       false,
	Email:          "",
	AccountID:      "",
	SubscriptionID: "",
	ContractID:     "",
}

const CRED = "\033[91m"
const CGREEN = "\033[32m"
const CYELLOW = "\033[33m"
const CBLUE = "\033[34m"
const CGREY = "\033[90m"
const CEND = "\033[0m"
