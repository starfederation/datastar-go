module github.com/starfederation/datastar-go/cmd/examples/helloworld

go 1.24

replace github.com/starfederation/datastar-go => ../../../

require (
	github.com/go-chi/chi/v5 v5.2.2
	github.com/starfederation/datastar-go v1.0.1
)

require (
	github.com/CAFxX/httpcompression v0.0.9 // indirect
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
)
