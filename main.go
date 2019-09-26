package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/yousseffarkhani/court/database"
	"github.com/yousseffarkhani/court/server"
)

func main() {
	file, err := openCourtFile("assets/courts.json")
	if err != nil {
		log.Println(fmt.Errorf("Problem opening %s, %v", file.Name(), err))
	}

	courtStore, err := database.NewCourtStore(file)
	DisplayError(err)
	defer courtStore.Close()

	server, err := server.NewBasketServer(courtStore)
	DisplayError(err)

	port := os.Getenv("PORT")

	if os.Getenv("APP_ENV") == "production" {
		if port == "" {
			port = "443"
		}
		pwd, _ := os.Getwd()
		fmt.Printf("Starting the server on port: %s\n", port)
		pathToCertFile := os.Getenv("CERTFILE")
		pathToPrivKey := os.Getenv("PRIVKEY")
		log.Fatal(http.ListenAndServeTLS(":"+port, filepath.Join(pwd, pathToCertFile), filepath.Join(pwd, pathToPrivKey), server))
	}
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Starting the server on port: %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, server))

}

func DisplayError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// Opens the file. If file non existent creates a JSON file.
func openCourtFile(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Problem opening %s, %v", file.Name(), err)
	}

	file.Seek(0, 0)

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("Couldn't get file info from file : %s", file.Name())
	}

	if fileInfo.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, 0)
	}

	return file, nil
}
