# Faye + Go

Original code forked from: https://github.com/roncohen/faye-go

Experimental

## Usage

```go

package main

import (
	"git.xuvasi.com/gocode/faye-go"
	"git.xuvasi.com/gocode/faye-go/adapters"
	"net/http"
)

func main() {
	fayeServer := faye.NewServer(faye.NewEngine())
	http.Handle("/faye", adapters.FayeHandler(fayeServer))

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
```
