package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")

	// Init data source, example Redis repo instance

	// Init shortener service instance ( in needs the data source passed to it as dependency)

	// Init our API handlers

	// Start out API/http server

	// Handle kill signals gracefully

}
