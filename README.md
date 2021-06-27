# go-improvmx

A Golang API client for [ImprovMX](https://improvmx.com), my email forwarding service of choice.

## Usage

```golang
package main

import (
	"fmt"
	"os"

	improvmx "github.com/issyl0/go-improvmx"
)

func main() {
	client := improvmx.NewClient(os.Getenv("IMPROVMX_API_TOKEN"))

	_, err := client.AccountDetails()
	if err != "" {
		fmt.Println(err)
	}

	_, err = client.CreateDomain("example.com")
	if err != "" {
		fmt.Println(err)
	}

	_, err = client.CreateEmailForward("example.com", "hi", "hi@realdomain.com")
	if err != "" {
		fmt.Println(err)
	}

	_, err = client.DeleteEmailForward("example.com", "hi")
	if err != "" {
		fmt.Println(err)
	}

	_, err = client.DeleteDomain("example.com")
	if err != "" {
		fmt.Println(err)
	}
}
```

## TODO

- [ ] Tests.
- [ ] The future side project now this is done... a Terraform provider for ImprovMX email forwards.
