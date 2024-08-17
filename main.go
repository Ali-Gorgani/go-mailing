package main

import goMailing "go-mailing/cmd/go-mailing"

func main() {
	err := goMailing.StartServer()
	if err != nil {
		panic(err)
	}
}
