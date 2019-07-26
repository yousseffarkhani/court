package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"courtdb"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	courtStore, err := courtdb.NewCourtStore()
	CheckError(err)
	defer courtStore.Close()

	// server, err := server.NewBasketServer(courtStore)
	// CheckError(err)
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		courtss := courtStore.GetAllCourts()
		json.NewEncoder(w).Encode(courtss)
	})
	fmt.Printf("Starting the server on port: %s\n", port)
	// log.Fatal(http.ListenAndServe(":"+port, server))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func CheckError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
