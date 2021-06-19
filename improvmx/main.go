package improvmx

import (
	"encoding/json"
	"fmt"

	resty "github.com/go-resty/resty/v2"
)

// https://improvmx.com/api/#domains-list data structure
type Domains struct {
	Success bool `json:"success"`
	Total   int  `json:"total"`
	Domains []struct {
		Name string `json:"domain"`
	} `json:"domains"`
}

// https://improvmx.com/api/#account data structure
type Account struct {
	Success bool `json:"success"`
	Account struct {
		Premium bool `json:"premium"`
	} `json:"account"`
}

// https://improvmx.com/api/#authentication handler
func new(accessToken string) (*resty.Request, string) {
	if accessToken == "" {
		fmt.Println("ERROR: An ImprovMX API access token is required. Create one at https://app.improvmx.com/api.")
		return nil, ""
	}

	client := resty.New().R()
	client = client.SetAuthScheme("Basic").SetAuthToken(fmt.Sprintf("api:%s", accessToken)).SetHeader("Content-Type", "application/json")

	return client, ""
}

// A (manual) test function to ensure a token is connected to an account.
// TODO: More data fields.
func AccountDetails(accessToken string) {
	client, _ := new(accessToken)
	if client == nil {
		fmt.Println("Something went wrong.")
		return
	}

	resp, err := client.Get("https://api.improvmx.com/v3/account")
	if err != nil {
		fmt.Printf("That API key doesn't work: %v", err)
		return
	}

	parsedAccount := Account{}
	json.Unmarshal(resp.Body(), &parsedAccount)

	if parsedAccount.Success {
		fmt.Println(parsedAccount.Account.Premium)
	}
}

// https://improvmx.com/api/#domains-list
func ListDomains(accessToken string) bool {
	client, _ := new(accessToken)

	resp, err := client.Get("https://api.improvmx.com/v3/domains?limit=100")
	if err != nil {
		fmt.Printf("ERROR: Couldn't get domains. %v", err)
		return false
	}

	parsedDomains := Domains{}
	json.Unmarshal(resp.Body(), &parsedDomains)

	for _, domain := range parsedDomains.Domains {
		fmt.Println(domain.Name)
	}
	return true
}

// https://improvmx.com/api/#domains-add
func CreateDomain(accessToken, domain string) bool {
	client, _ := new(accessToken)

	domainInput, err := json.Marshal(map[string]string{"domain": domain})
	if err != nil {
		fmt.Printf("Couldn't convert string to JSON: %v", err)
		return false
	}

	resp, err := client.SetBody(domainInput).Post("https://api.improvmx.com/v3/domains/")
	// TODO: `err` here isn't actually the API response error, they're still in `resp`.
	if err != nil {
		fmt.Printf("%v", err)
		return false
	}

	fmt.Println(string(resp.Body()))
	return true
}

// https://improvmx.com/api/#domain-delete
func DeleteDomain(accessToken, domain string) bool {
	client, _ := new(accessToken)

	resp, err := client.Delete(fmt.Sprintf("https://api.improvmx.com/v3/domains/%s", domain))
	// TODO: `err` here isn't actually the API response error, they're still in `resp`.
	if err != nil {
		fmt.Printf("Couldn't delete domain, got error %v", err)
		return false
	}

	fmt.Println(string(resp.Body()))
	return true
}

// https://improvmx.com/api/#alias-add
func CreateEmailForward(accessToken, domain, alias, forward string) bool {
	client, _ := new(accessToken)

	emailForwardInput, err := json.Marshal(map[string]string{"alias": alias, "forward": forward})
	if err != nil {
		fmt.Printf("Couldn't convert input to JSON, %v", err)
		return false
	}

	resp, err := client.SetBody(emailForwardInput).Post(fmt.Sprintf("https://api.improvmx.com/v3/domains/%s/aliases", domain))
	// TODO: `err` here isn't actually the API response error, they're still in `resp`.
	if err != nil {
		fmt.Printf("Couldn't create email forward, got error %v", err)
	}

	fmt.Println(string(resp.Body()))
	return true
}

// https://improvmx.com/api/#alias-delete
func DeleteEmailForward(accessToken, domain, alias string) bool {
	client, _ := new(accessToken)

	resp, err := client.Delete(fmt.Sprintf("https://api.improvmx.com/v3/domains/%s/aliases/%s", domain, alias))
	// TODO: `err` here isn't actually the API response error, they're still in `resp`.
	if err != nil {
		fmt.Printf("Couldn't delete email forward, got error %v", err)
	}

	fmt.Println(string(resp.Body()))
	return true
}
