//go:build tools

package main

// This file imports packages that are used in build scripts but not otherwise referenced in the code.
import (
	_ "golang.org/x/tools/cmd/goimports"
)
