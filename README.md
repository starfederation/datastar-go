<p align="center"><img width="150" height="150" src="https://data-star.dev/static/images/rocket-512x512.png"></p>

# Datastar Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/starfederation/datastar-go.svg)](https://pkg.go.dev/github.com/starfederation/datastar-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/starfederation/datastar-go)](https://goreportcard.com/report/github.com/starfederation/datastar-go)

This package provides a Go SDK for working with Datastar.

## License

This package is licensed for free under the [MIT License](LICENSE).

## Requirements

This package requires Go 1.24 or later.

## Installation

```bash
go get github.com/starfederation/datastar-go
```

## Usage

```go
import (
    "net/http"
    "github.com/starfederation/datastar-go/datastar"
)

// Read signals from request
type Store struct {
    Message string `json:"message"`
    Count   int    `json:"count"`
}

func handler(w http.ResponseWriter, r *http.Request) {
    // Read signals from the request
    store := &Store{}
    if err := datastar.ReadSignals(r, store); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Create a Server-Sent Event writer
    sse := datastar.NewSSE(w, r)

    // Patch elements in the DOM
    sse.PatchElements(`<div id="output">Hello from Datastar!</div>`)

    // Remove elements from the DOM
    sse.RemoveElements("#temporary-element")

    // Patch signals (update client-side state)
    sse.PatchSignals(map[string]any{
        "message": "Updated message",
        "count":   store.Count + 1,
    })

    // Execute JavaScript in the browser
    sse.ExecuteScript(`console.log("Hello from server!")`)

    // Redirect the browser
    sse.Redirect("/new-page")
}
```

## Examples

See the [examples directory](cmd/examples) for complete working examples.

## Testing

```bash
go test ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.