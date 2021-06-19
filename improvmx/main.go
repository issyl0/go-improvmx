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

type Client struct {
	BaseURL     string
	AccessToken string
	Http        *resty.Request
}

func NewClient(accessToken string) *Client {
	return &Client{
		BaseURL:     "https://api.improvmx.com/api/v3",
		AccessToken: accessToken,
		Http:        resty.New().R(),
	}
}

// https://improvmx.com/api/#authentication handler
func (client *Client) setHeaders() *resty.Request {
	if client.AccessToken == "" {
		fmt.Println("ERROR: An ImprovMX API access token is required. Create one at https://app.improvmx.com/api.")
		return nil
	}

	return client.Http.SetAuthScheme("Basic").SetAuthToken(fmt.Sprintf("api:%s", client.AccessToken)).SetHeader("Content-Type", "application/json")
}

// A (manual) test function to ensure a token is connected to an account.
// TODO: More data fields.
func (client *Client) AccountDetails() {
	resp, err := client.setHeaders().Get(fmt.Sprintf("%s/account", client.BaseURL))
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
func (client *Client) ListDomains() bool {
	resp, err := client.setHeaders().Get(fmt.Sprintf("%s/domains?limit=100", client.BaseURL))
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
func (client *Client) CreateDomain(domain string) bool {
	domainInput, err := json.Marshal(map[string]string{"domain": domain})
	if err != nil {
		fmt.Printf("Couldn't convert string to JSON: %v", err)
		return false
	}

	resp, err := client.setHeaders().SetBody(domainInput).Post(fmt.Sprintf("%s/domains/", client.BaseURL))
	// TODO: `err` here isn't actually the API response error, they're still in `resp`.
	if err != nil {
		fmt.Printf("%v", err)
		return false
	}

	fmt.Println(string(resp.Body()))
	return true
}

// https://improvmx.com/api/#domain-delete
func (client *Client) DeleteDomain(domain string) bool {
	resp, err := client.setHeaders().Delete(fmt.Sprintf("%s/domains/%s", client.BaseURL, domain))
	// TODO: `err` here isn't actually the API response error, they're still in `resp`.
	if err != nil {
		fmt.Printf("Couldn't delete domain, got error %v", err)
		return false
	}

	fmt.Println(string(resp.Body()))
	return true
}

// https://improvmx.com/api/#alias-add
func (client *Client) CreateEmailForward(domain, alias, forward string) bool {
	emailForwardInput, err := json.Marshal(map[string]string{"alias": alias, "forward": forward})
	if err != nil {
		fmt.Printf("Couldn't convert input to JSON, %v", err)
		return false
	}

	resp, err := client.setHeaders().SetBody(emailForwardInput).Post(fmt.Sprintf("%s/domains/%s/aliases", client.BaseURL, domain))
	// TODO: `err` here isn't actually the API response error, they're still in `resp`.
	if err != nil {
		fmt.Printf("Couldn't create email forward, got error %v", err)
	}

	fmt.Println(string(resp.Body()))
	return true
}

// https://improvmx.com/api/#alias-delete
func (client *Client) DeleteEmailForward(domain, alias string) bool {
	resp, err := client.setHeaders().Delete(fmt.Sprintf("%s/domains/%s/aliases/%s", client.BaseURL, domain, alias))
	// TODO: `err` here isn't actually the API response error, they're still in `resp`.
	if err != nil {
		fmt.Printf("Couldn't delete email forward, got error %v", err)
	}

	fmt.Println(string(resp.Body()))
	return true
}
