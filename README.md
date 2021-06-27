# go-improvmx

A Golang API client for [ImprovMX](https://improvmx.com), my email forwarding service of choice.

## Usage

```golang
package main

import (
	"os"

	improvmx "github.com/issyl0/go-improvmx"
)

func main() {
	client := improvmx.NewClient(os.Getenv("IMPROVMX_API_TOKEN"))
	client.AccountDetails()
	client.ListDomains()
	client.CreateDomain("example.com")
	client.CreateEmailForward("example.com", "hi", "hi@realdomain.com")
	client.DeleteEmailForward("example.com", "hi")
	client.DeleteDomain("example.com")
}
```

## TODO

- [ ] Tests.
- [ ] The future side project now this is done... a Terraform provider for ImprovMX email forwards.
