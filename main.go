package main

import (
	"fmt"
	"os"
)

func main() {
	// notest: tested through module test
	router := CreateRouter()

	err := router.Run("0.0.0.0:8080")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running server: %s\n", err)
	}
}
