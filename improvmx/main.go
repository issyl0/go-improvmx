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
