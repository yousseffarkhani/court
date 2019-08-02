package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"courtdb"
	"server"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	courtStore, err := courtdb.NewCourtStore()
	CheckError(err)
	defer courtStore.Close()

	server, err := server.NewBasketServer(courtStore)
	CheckError(err)

	fmt.Printf("Starting the server on port: %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, server))
}

func CheckError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
