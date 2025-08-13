package main

import (
	"fmt"
	"os"
)

const (
	cdn = "https://cdn.jsdelivr.net/gh/starfederation/datastar@develop/bundles/datastar.js"
)

// generates an hotreload.html file to be used to play with the hotreload example
func main() {
	os.WriteFile("hotreload.html", []byte(fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">

<head>
	<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=0" />
	<script type="module" defer src="%s"></script>
</head>

<!-- next line mounts the refresh script -->
<body data-on-load="@get('/hotreload', {retryInterval: 100})">
	<main>
		<p>
			This page will automatically reload on any change to the hotreload.html file. Update this paragraph, save changes, and
			switch back to the open browser tab to observe
			the update.
		</p>
	</main>
</body>

</html>`, cdn)), 0644)
}
