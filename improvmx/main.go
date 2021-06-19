package improvmx

import (
	"encoding/json"
	"fmt"

	resty "github.com/go-resty/resty/v2"
)

type Client struct {
	BaseURL     string
	AccessToken string
	Http        *resty.Request
}

type Response struct {
	Success bool   `json:"success"`
	Total   int    `json:"total"`
	Code    int    `json:"code"`
	Error   string `json:"error"`

	Account struct {
		Premium bool `json:"premium"`
	} `json:"account"`

	Domains []struct {
		Name string `json:"domain"`
	} `json:"domains"`
}

func NewClient(accessToken string) *Client {
	return &Client{
		BaseURL:     "https://api.improvmx.com/v3",
		AccessToken: accessToken,
		Http:        resty.New().R(),
	}
}

// https://improvmx.com/api/#authentication handler
func (client *Client) setHeaders() *resty.Request {
	if client.AccessToken == "" {
		fmt.Println("ERROR: An ImprovMX API access token is required. Create one at https://app.improvmx.com/api.")
	}

	return client.Http.SetAuthScheme("Basic").SetAuthToken(fmt.Sprintf("api:%s", client.AccessToken)).SetHeader("Content-Type", "application/json")
}

// A (manual) test function to ensure a token is connected to an account.
// TODO: More data fields.
func (client *Client) AccountDetails() {
	resp, _ := client.setHeaders().Get(fmt.Sprintf("%s/account", client.BaseURL))

	parsed := Response{}
	json.Unmarshal(resp.Body(), &parsed)

	if parsed.Success {
		fmt.Println(parsed.Account.Premium)
	} else {
		fmt.Printf("ERROR: Couldn't find account details. Message: %v", parsed.Error)
	}
}

// https://improvmx.com/api/#domains-list
func (client *Client) ListDomains() bool {
	resp, _ := client.setHeaders().Get(fmt.Sprintf("%s/domains?limit=100", client.BaseURL))

	parsed := Response{}
	json.Unmarshal(resp.Body(), &parsed)

	if parsed.Success {
		for _, domain := range parsed.Domains {
			fmt.Println(domain.Name)
		}
		return true
	} else {
		fmt.Printf("ERROR: Couldn't get domains. Message: %s", parsed.Error)
		return false
	}
}

// https://improvmx.com/api/#domains-add
func (client *Client) CreateDomain(domain string) bool {
	domainInput, err := json.Marshal(map[string]string{"domain": domain})
	if err != nil {
		fmt.Printf("Couldn't convert string to JSON: %v", err)
		return false
	}

	resp, _ := client.setHeaders().SetBody(domainInput).Post(fmt.Sprintf("%s/domains/", client.BaseURL))

	parsed := Response{}
	json.Unmarshal(resp.Body(), &parsed)

	if parsed.Success {
		fmt.Printf("SUCCESS: Domain %s created.", domain)
		return true
	} else {
		fmt.Printf("ERROR: Couldn't create domain. Message: %s", parsed.Error)
		return false
	}
}

// https://improvmx.com/api/#domain-delete
func (client *Client) DeleteDomain(domain string) bool {
	resp, _ := client.setHeaders().Delete(fmt.Sprintf("%s/domains/%s", client.BaseURL, domain))

	parsed := Response{}
	json.Unmarshal(resp.Body(), &parsed)

	if parsed.Success {
		fmt.Printf("SUCCESS: Domain %s and all its forwards deleted.", domain)
		return true
	} else {
		fmt.Printf("ERROR: Couldn't delete domain. Message: %s", parsed.Error)
		return false
	}
}

// https://improvmx.com/api/#alias-add
func (client *Client) CreateEmailForward(domain, alias, forward string) bool {
	emailForwardInput, err := json.Marshal(map[string]string{"alias": alias, "forward": forward})
	if err != nil {
		fmt.Printf("Couldn't convert input to JSON, %v", err)
		return false
	}

	resp, _ := client.setHeaders().SetBody(emailForwardInput).Post(fmt.Sprintf("%s/domains/%s/aliases", client.BaseURL, domain))

	parsed := Response{}
	json.Unmarshal(resp.Body(), &parsed)

	if parsed.Success {
		fmt.Printf("SUCCESS: Created email forward from %s@%s to %s.", alias, domain, forward)
		return true
	} else {
		fmt.Printf("ERROR: Couldn't create email forward. Message: %s", parsed.Error)
		return false
	}
}

// https://improvmx.com/api/#alias-delete
func (client *Client) DeleteEmailForward(domain, alias string) bool {
	resp, _ := client.setHeaders().Delete(fmt.Sprintf("%s/domains/%s/aliases/%s", client.BaseURL, domain, alias))

	parsed := Response{}
	json.Unmarshal(resp.Body(), &parsed)

	if parsed.Success {
		fmt.Printf("SUCCESS: Deleted email forward %s@%s.", alias, domain)
		return true
	} else {
		fmt.Printf("ERROR: Couldn't delete email forward. Message: %s", parsed.Error)
		return false
	}
}
