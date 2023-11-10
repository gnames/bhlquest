package bhlquest

import "fmt"

// GetVersion provides version information of the app.
func GetVersion() string {
	version := fmt.Sprintf("Version: %s\nBuild:   %s", Version, Build)
	return version
}
