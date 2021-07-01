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

	client.CreateDomain("example.com")
	client.CreateEmailForward("example.com", "hi", "person@realdomain.com")
	client.UpdateEmailForward("example.com", "hi", "person@realdomain.com,another@realdomain.com")
	client.DeleteEmailForward("example.com", "hi")
	client.DeleteDomain("example.com")
}
```

## TODO

- [ ] Tests.
