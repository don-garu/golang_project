package main

import (
	"Database"
	"HttpAction"
)

func main() {
	Database.Initialize()
	HttpAction.HandleRequest()

	defer Database.Close()
}