package ui

import (
	"embed"
)

// embed files as part of go binary
// prevents having to read files from disk at runtime

//go:embed "html" "static"
var Files embed.FS