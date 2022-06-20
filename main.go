package main

import (
	"fmt"
)

// TODO
// rewrite generate results method - better way to determine how many goals a team scores, and concedes?
// logging
// go routines
// tests
// any other refactoring?
// documentation

func main() {
	filename := getFilenameFromUser()
	if filename == "" {
		filename = generateResultsFile()
	}
	results := getResultsFromFile(filename)

	table := Table{}

	// Build table with results
	for _, result := range results {
		table.updateTableWithResult(result)
	}

	fmt.Println(table.formatTable())
}

func getFilenameFromUser() string {
	fmt.Println("Enter results file name or leave empty to auto-generate results:")
	var filename string
	fmt.Scanln(&filename)
	return filename
}
