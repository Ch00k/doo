package main

import "fmt"

func main() {
	db := SetupDB()
	r := SetupRouter(db)

	httpHost := GetEnvVar("DOO_HTTP_HOST", "localhost")
	httpPort := GetEnvVar("DOO_HTTP_PORT", "8080")
	r.Run(fmt.Sprintf("%s:%s", httpHost, httpPort))
}
