package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/starfederation/datastar-go/datastar"

	_ "embed"
)

const (
	serverAddress = "localhost:9001"
)

//go:embed hotreload.html
var indexHTML []byte

var hotReloadOnlyOnce sync.Once

func HotReloadHandler(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	hotReloadOnlyOnce.Do(func() {
		// Refresh the client page as soon as connection
		// is established. This will occur only once
		// after the server starts.
		sse.ExecuteScript("window.location.reload()")
	})

	// Freeze the event stream until the connection
	// is lost for any reason. This will force the client
	// to attempt to reconnect after the server reboots.
	<-r.Context().Done()
}

func PageWithHotReload(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write(indexHTML)
}

func main() {
	// Hot reload requires a file system watcher and
	// a refresh script. [Reflex](https://github.com/cespare/reflex)
	// is one of the best tools for running a command on
	// each file change.
	//
	// Setup/Run the example with `task hotreload`
	//
	// The refresh script is a Datastar handler
	// that emits a page refresh event only once
	// for each server start.
	//
	// When the the file watcher forces the server to restart,
	// Datastar client will lose the network connection to the
	// server and attempt to reconnect. Once the connection is
	// established, the client will receive the refresh event.
	http.HandleFunc("/hotreload", HotReloadHandler)
	http.HandleFunc("/", PageWithHotReload)
	slog.Info(fmt.Sprintf(
		"Open your browser to: http://%s/",
		serverAddress,
	))
	http.ListenAndServe(serverAddress, nil)

	// Tip: read the reflex documentation to see advanced usage
	// options like responding to specific file changes by filter.
	//
	// $ reflex --help
}
