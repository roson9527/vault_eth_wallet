package modules

// AccountAddress stores the name of the account to allow reverse lookup by address
type AccountAddress struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type Address struct {
	Name string `json:"name"`
}
