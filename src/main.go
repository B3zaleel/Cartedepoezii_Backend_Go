package main

import (
	"fmt"
	"os"

	"github.com/B3zaleel/Cartedepoezii_Backend_Go/src/configs"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	host := "0.0.0.0"

	if len(os.Getenv("HOST")) > 0 {
		host = os.Getenv("")
	}
	configs.AddEndpoints(server)
	server.Run(fmt.Sprintf("%s:5000", host))
}
