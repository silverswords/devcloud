package main

import (
	"os"

	"github.com/silverswords/devcloud/pkg/rest"
)

func main() {
	os.Setenv("MONGO_SERVER_URL", "mongodb://localhost:27017")
	rest.Connect()
	rest.Run()
}
