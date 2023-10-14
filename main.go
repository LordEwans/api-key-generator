package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lordewans/api-key-generator/internal/routes"
	"github.com/lordewans/api-key-generator/pkg/generate"
)

const defaultPort = "9080"

func main() {
	ch := make(chan bool)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	r := gin.Default()
	go routes.UseRoute(r)
	go fmt.Println(generate.GenerateKey())

	go log.Printf("Server started at http://localhost:%s", port)
	go log.Fatal(r.Run(":" + port))

	<-ch
}
