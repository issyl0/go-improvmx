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

	Records struct {
		Valid    bool   `json:"valid"`
		Provider string `json:"provider"`
		Advanced bool   `json:"advanced"`
		Mx       struct {
			Valid    bool     `json:"valid"`
			Expected []string `json:"expected"`
			Values   []string `json:"values"`
		} `json:"mx"`
		Spf struct {
			Valid    bool   `json:"valid"`
			Expected string `json:"expected"`
			Values   string `json:"values"`
		} `json:"spf"`
	}

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
		Domain            string `json:"domain"`
		Whitelabel        string `json:"whitelabel"`
		NotificationEmail string `json:"notification_email"`
		Aliases           []struct {
			Id int64 `json:"id"`
		}
	} `json:"domain"`

	Credential struct {
		Created  int64  `json:"created"`
		Usage    int    `json:"usage"`
		Username string `json:"username"`
	} `json:"credential"`

	Credentials []struct {
		Created  int64  `json:"created"`
		Usage    int    `json:"usage"`
		Username string `json:"username"`
	} `json:"credentials"`
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
func (client *Client) AccountDetails() Response {
	resp, _ := client.setHeaders().Get(fmt.Sprintf("%s/account", client.BaseURL))

	parsed := Response{}
	json.Unmarshal(resp.Body(), &parsed)

	return parsed
}

func (client *Client) GetDomain(domain string) Response {
	resp, _ := client.setHeaders().Get(fmt.Sprintf("%s/domains/%s", client.BaseURL, domain))

	parsed := Response{}
	json.Unmarshal(resp.Body(), &parsed)

	return parsed
}

// https://improvmx.com/api/#domains-add
func (client *Client) CreateDomain(domain, notificationEmail, whitelabel string) Response {
	domainInput, _ := json.Marshal(map[string]string{
		"domain":             domain,
		"notification_email": notificationEmail,
		"whitelabel":         whitelabel,
	})

	resp, _ := client.setHeaders().SetBody(domainInput).Post(fmt.Sprintf("%s/domains/", client.BaseURL))

	parsed := Response{}
	json.Unmarshal(resp.Body(), &parsed)

	return parsed
}

func (client *Client) UpdateDomain(domain, notificationEmail, whitelabel string) Response {
	domainInput, _ := json.Marshal(map[string]string{
		"domain":             domain,
		"notification_email": notificationEmail,
		"whitelabel":         whitelabel,
	})

	resp, _ := client.setHeaders().SetBody(domainInput).Put(fmt.Sprintf("%s/domains/%s", client.BaseURL, domain))

	parsed := Response{}
	json.Unmarshal(resp.Body(), &parsed)

	return parsed
}

// https://improvmx.com/api/#domain-delete
func (client *Client) DeleteDomain(domain string) Response {
	resp, _ := client.setHeaders().Delete(fmt.Sprintf("%s/domains/%s", client.BaseURL, domain))

	parsed := Response{}
	json.Unmarshal(resp.Body(), &parsed)

	return parsed
}

func (client *Client) GetDomainCheck(domain string) Response {
	resp, _ := client.setHeaders().Get(fmt.Sprintf("%s/domains/%s/check", client.BaseURL, domain))

	parsed := Response{}
	json.Unmarshal(resp.Body(), &parsed)

	return parsed
}

func (client *Client) CreateEmailForward(domain, alias, forward string) Response {
	emailForwardInput, _ := json.Marshal(map[string]string{"alias": alias, "forward": forward})

	resp, _ := client.setHeaders().SetBody(emailForwardInput).Post(fmt.Sprintf("%s/domains/%s/aliases", client.BaseURL, domain))

	parsed := Response{}
	json.Unmarshal(resp.Body(), &parsed)

	return parsed
}

func (client *Client) GetEmailForward(domain, alias string) Response {
	resp, _ := client.setHeaders().Get(fmt.Sprintf("%s/domains/%s/aliases/%s", client.BaseURL, domain, alias))

	parsed := Response{}
	json.Unmarshal(resp.Body(), &parsed)

	return parsed
}

func (client *Client) UpdateEmailForward(domain, alias, forward string) Response {
	emailForwardInput, _ := json.Marshal(map[string]string{"forward": forward})
	resp, _ := client.setHeaders().SetBody(emailForwardInput).Put(fmt.Sprintf("%s/domains/%s/aliases/%s", client.BaseURL, domain, alias))

	parsed := Response{}
	json.Unmarshal(resp.Body(), &parsed)

	return parsed
}

// https://improvmx.com/api/#alias-delete
func (client *Client) DeleteEmailForward(domain, alias string) Response {
	resp, _ := client.setHeaders().Delete(fmt.Sprintf("%s/domains/%s/aliases/%s", client.BaseURL, domain, alias))

	parsed := Response{}
	json.Unmarshal(resp.Body(), &parsed)

	return parsed
}

func (client *Client) GetSMTPCredential(domain string) Response {
	resp, _ := client.setHeaders().Get(fmt.Sprintf("%s/domains/%s/credentials", client.BaseURL, domain))

	parsed := Response{}
	json.Unmarshal(resp.Body(), &parsed)

	return parsed
}

func (client *Client) CreateSMTPCredential(domain, username, password string) Response {
	credentialsInput, _ := json.Marshal(map[string]string{"username": username, "password": password})
	resp, _ := client.setHeaders().SetBody(credentialsInput).Post(fmt.Sprintf("%s/domains/%s/credentials", client.BaseURL, domain))

	parsed := Response{}
	json.Unmarshal(resp.Body(), &parsed)

	return parsed
}

func (client *Client) UpdateSMTPCredential(domain, username, password string) Response {
	credentialsInput, _ := json.Marshal(map[string]string{"password": password})
	resp, _ := client.setHeaders().SetBody(credentialsInput).Put(fmt.Sprintf("%s/domains/%s/credentials/%s", client.BaseURL, domain, username))

	parsed := Response{}
	json.Unmarshal(resp.Body(), &parsed)

	return parsed
}

func (client *Client) DeleteSMTPCredential(domain, username string) Response {
	resp, _ := client.setHeaders().Delete(fmt.Sprintf("%s/domains/%s/credentials/%s", client.BaseURL, domain, username))

	parsed := Response{}
	json.Unmarshal(resp.Body(), &parsed)

	return parsed
}
