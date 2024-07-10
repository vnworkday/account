package main

import (
	_ "github.com/lib/pq"
	"github.com/vnworkday/account/cmd/app"
)

func main() {
	app.Run()
}
