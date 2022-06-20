package main

import (
	"fmt"
	"log"
)

// TODO
// logging
// split getResultsFromFile into separate methods to aid testing
// go routines
// tests
// any other refactoring?
// documentation

func main() {
	filename := getFilenameFromUser()
	if filename == "" {
		filename = generateResultsToFile()
	}
	results := getResultsFromFile(filename)

	table := Table{}

	// Build table with results
	log.Print("Building table with results")
	for _, result := range results {
		table.updateTableWithResult(result)
	}

	fmt.Println(table.formatTable())
}

func getFilenameFromUser() string {
	fmt.Println("Enter results file name or leave empty to auto-generate results:")
	var filename string
	fmt.Scanln(&filename)
	log.Printf("User entered filename: '%v'", filename)
	return filename
}
