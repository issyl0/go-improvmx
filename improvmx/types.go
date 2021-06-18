package improvmx

type Account struct {
	Success bool `json:"success"`
	Account struct {
		Premium bool `json:"premium"`
	} `json:"account"`
}

type Domains struct {
	Success bool `json:"success"`
	Total   int  `json:"total"`
	Domains []struct {
		Name string `json:"domain"`
	} `json:"domains"`
}
