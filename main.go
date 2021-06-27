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
	Success bool `json:"success"`
	Total   int  `json:"total"`
	Code    int  `json:"code"`

	Errors struct {
		Domain  []string `json:"domain"`
		Alias   []string `json:"alias"`
		Account []string `json:"account"`
	} `json:"errors"`

	Account struct {
		Plan struct {
			Display string `json:"display"`
		} `json:"plan"`
	} `json:"account"`

	Alias struct {
		Forward string `json:"forward"`
		Alias   string `json:"alias"`
		Id      int64  `json:"id"`
	} `json:"alias"`

	Domain struct {
		Domain     string `json:"domain"`
		Whitelabel string `json:"whitelabel"`
		Id         int64  `json:"id"`
	} `json:"domain"`
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
		return nil
	}

	return client.Http.SetAuthScheme("Basic").SetAuthToken(fmt.Sprintf("api:%s", client.AccessToken)).SetHeader("Content-Type", "application/json")
}

// https://improvmx.com/api/#account handler
func (client *Client) AccountDetails() (bool, string) {
	resp, _ := client.setHeaders().Get(fmt.Sprintf("%s/account", client.BaseURL))

	parsed := Response{}
	json.Unmarshal(resp.Body(), &parsed)

	if parsed.Success {
		return true, ""
	} else {
		return false, parsed.Errors.Account[0]
	}
}

func (client *Client) GetDomain(domain string) Response {
	resp, _ := client.setHeaders().Get(fmt.Sprintf("%s/domains/%s", client.BaseURL, domain))

	parsed := Response{}
	json.Unmarshal(resp.Body(), &parsed)

	return parsed
}

// https://improvmx.com/api/#domains-add
func (client *Client) CreateDomain(domain string) (bool, string) {
	domainInput, err := json.Marshal(map[string]string{"domain": domain})
	if err != nil {
		return false, fmt.Sprintf("Couldn't convert string to JSON: %v\n", err)
	}

	resp, _ := client.setHeaders().SetBody(domainInput).Post(fmt.Sprintf("%s/domains/", client.BaseURL))

	parsed := Response{}
	json.Unmarshal(resp.Body(), &parsed)

	if parsed.Success {
		return true, ""
	} else {
		return false, parsed.Errors.Domain[0]
	}
}

// https://improvmx.com/api/#domain-delete
func (client *Client) DeleteDomain(domain string) (bool, string) {
	resp, _ := client.setHeaders().Delete(fmt.Sprintf("%s/domains/%s", client.BaseURL, domain))

	parsed := Response{}
	json.Unmarshal(resp.Body(), &parsed)

	return parsed.Success, ""
}

// https://improvmx.com/api/#alias-add
func (client *Client) CreateEmailForward(domain, alias, forward string) (bool, string) {
	emailForwardInput, err := json.Marshal(map[string]string{"alias": alias, "forward": forward})
	if err != nil {
		return false, fmt.Sprintf("Couldn't convert input to JSON, %v", err)
	}

	resp, _ := client.setHeaders().SetBody(emailForwardInput).Post(fmt.Sprintf("%s/domains/%s/aliases", client.BaseURL, domain))

	parsed := Response{}
	json.Unmarshal(resp.Body(), &parsed)

	if parsed.Success {
		return true, ""
	} else {
		return false, parsed.Errors.Alias[0]
	}
}

func (client *Client) GetEmailForward(domain, alias string) Response {
	resp, _ := client.setHeaders().Get(fmt.Sprintf("%s/domains/%s/aliases/%s", client.BaseURL, domain, alias))

	parsed := Response{}
	json.Unmarshal(resp.Body(), &parsed)

	return parsed
}

// https://improvmx.com/api/#alias-delete
func (client *Client) DeleteEmailForward(domain, alias string) (bool, string) {
	resp, _ := client.setHeaders().Delete(fmt.Sprintf("%s/domains/%s/aliases/%s", client.BaseURL, domain, alias))

	parsed := Response{}
	json.Unmarshal(resp.Body(), &parsed)

	return parsed.Success, ""
}
