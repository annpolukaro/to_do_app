package main

import (
	"fmt"
	"to_do_app/http_server"
)

func main() {

	fmt.Println("HTTP run!")
	err := http_server.StartHTTPServer()
	if err != nil {
		fmt.Println("In prosses work server, has been mistackes")
	}
}
