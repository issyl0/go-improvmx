package improvmx

import (
	"encoding/json"
	"fmt"

	resty "github.com/go-resty/resty/v2"
)

func authorize(accessToken string) (*resty.Request, string) {
	if accessToken == "" {
		fmt.Println("ERROR: An ImprovMX API access token is required. Create one at https://app.improvmx.com/api.")
		return nil, ""
	}

	client := resty.New().R()
	client = client.SetAuthScheme("Basic").SetAuthToken(fmt.Sprintf("api:%s", accessToken)).SetHeader("Content-Type", "application/json")

	return client, ""
}

func Authenticate(accessToken string) {
	client, _ := authorize(accessToken)
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

func ListDomains(accessToken string) bool {
	client, _ := authorize(accessToken)

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

func CreateDomain(accessToken, domain string) bool {
	client, _ := authorize(accessToken)

	domainInput, err := json.Marshal(map[string]string{"domain": domain})
	if err != nil {
		fmt.Println("Couldn't convert string to JSON: %v", err)
		return false
	}

	resp, err := client.SetBody(domainInput).Post("https://api.improvmx.com/v3/domains/")
	if err != nil {
		fmt.Printf("%v", err)
		return false
	}

	fmt.Println(string(resp.Body()))
	return true
}

func DeleteDomain(accessToken, domain string) bool {
	client, _ := authorize(accessToken)

	resp, err := client.Delete(fmt.Sprintf("https://api.improvmx.com/v3/domains/%s", domain))
	if err != nil {
		fmt.Printf("Couldn't delete domain, got error %v", err)
		return false
	}

	fmt.Println(string(resp.Body()))
	return true
}

func CreateEmailForward(accessToken, domain, alias, forward string) bool {
	client, _ := authorize(accessToken)

	emailForwardInput, err := json.Marshal(map[string]string{"alias": alias, "forward": forward})
	if err != nil {
		fmt.Printf("Couldn't convert input to JSON, %v", err)
		return false
	}

	resp, err := client.SetBody(emailForwardInput).Post(fmt.Sprintf("https://api.improvmx.com/v3/domains/%s/aliases", domain))
	if err != nil {
		fmt.Printf("Couldn't create email forward, got error %v", err)
	}

	fmt.Println(string(resp.Body()))
	return true
}

func DeleteEmailForward(accessToken, domain, alias string) bool {
	client, _ := authorize(accessToken)

	resp, err := client.Delete(fmt.Sprintf("https://api.improvmx.com/v3/domains/%s/aliases/%s", domain, alias))
	if err != nil {
		fmt.Printf("Couldn't delete email forward, got error %v", err)
	}

	fmt.Println(string(resp.Body()))
	return true
}
