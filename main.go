package main

import (
	"email-template/api"
)

// @title Tradewinds API
// @version 1.0
// @description This is a Tradewinds Admin API.

// @contact.name API Support
// @contact.url tradewinds.com
// @contact.email support@tradewinds.com

// @host localhost:8888
// @BasePath /
// @schemes http
func main() {

	api.Run()

}