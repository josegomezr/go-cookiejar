go-cookiejar
===

A save-able Cookiejar implementation.

This is nothing but a carbon copy of go's Cookiejar implementation with the
ability to save the jar's content to disk.

I'm too lazy to redesign it up to spec right now, so copy-paste driven
development would do.

It's still in alpha but somewhat usable.

Usage
---

```go
package main

import (
	"os"
	"fmt"
	"net/http"
	"github.com/josegomezr/go-cookiejar"
)

func main() {
	// 1:1 copy of net/http/cookiejar
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(fmt.Sprintf("Error creating cookiejar: %s", err))
	}

	// get some cookies
	client := &http.Client{Jar: jar,}
	resp, err := client.Get("https://httpbingo.org/cookies/set?cookie=value")
	if err != nil {
		panic(fmt.Sprintf("Error on cookie request: %s", err))
	}

	fmt.Println(resp.Request.Method, resp.Request.URL.String(), resp.Status)
	fmt.Println()
	// show'em in stdout
	jar.WriteAsCurl(os.Stdout)
}
```
