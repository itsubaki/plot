package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range records {
		fmt.Println(r)
	}
}
