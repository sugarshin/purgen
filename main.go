package main

import "github.com/sugarshin/purgen/app"

func main() {
	if _, err := app.Run(); err != nil {
		panic(err)
	}
}
