package main

import (
	"fmt"
	"log"
)

func main() {
	db := SetupDB()
	r := SetupRouter(db)

	httpHost := GetEnvVar("DOO_HTTP_HOST", "localhost")
	httpPort := GetEnvVar("DOO_HTTP_PORT", "8080")

	err := r.Run(fmt.Sprintf("%s:%s", httpHost, httpPort))
	if err != nil {
		log.Fatal(err)
	}
}
