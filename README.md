# go-improvmx

A Golang API client for [ImprovMX](https://improvmx.com), my email forwarding service of choice.

## Usage

```golang
package main

import (
	"os"

	"github.com/issyl0/go-improvmx/improvmx"
)

func main() {
	accessToken := os.Getenv("IMPROVMX_API_TOKEN")
	improvmx.ListDomains(accessToken)
	improvmx.CreateDomain(accessToken, "example.com")
	improvmx.CreateEmailForward(accessToken, "example.com", "hi", "me@realdomain.com")
	improvmx.DeleteEmailForward(accessToken, "example.com", "hello")
	improvmx.DeleteDomain(accessToken, "example.com")
}
```

## TODO

- [ ] Avoid making users pass `accessToken` on every operation.
- [ ] Tests.
- [ ] Better error handling.
- [ ] Actual human output to the user, not just the JSON.
- [ ] The future side project now this is done... a Terraform provider for ImprovMX email forwards.
