package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/dmitryikh/leaves"
)

func main() {
	// Read actual file which we just updated
	data, err := os.ReadFile("pkg/engine/model/borehole_model.json")
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	fmt.Println("Attempting to load model from file...")
	reader := bufio.NewReader(bytes.NewReader(data))

	// Try loading with transformation = true
	_, err = leaves.XGEnsembleFromReader(reader, true)
	if err != nil {
		fmt.Printf("ERROR (loadTransformation=true): %v\n", err)

		// Try loading with transformation = false
		reader2 := bufio.NewReader(bytes.NewReader(data))
		_, err2 := leaves.XGEnsembleFromReader(reader2, false)
		if err2 != nil {
			fmt.Printf("ERROR (loadTransformation=false): %v\n", err2)
		} else {
			fmt.Println("SUCCESS with loadTransformation=false")
		}
	} else {
		fmt.Println("SUCCESS with loadTransformation=true")
	}
}
