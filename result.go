package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

const results_filename_prefix string = "results_data/"

type Result struct {
	homeTeam, awayTeam   string
	homeGoals, awayGoals int
}

func (r Result) IsHomeWin() bool {
	return r.homeGoals > r.awayGoals
}

func (r Result) IsAwayWin() bool {
	return r.homeGoals < r.awayGoals
}

func getResultsFromFile(filename string) []Result {
	log.Printf("Retrieving results from file. Filename: %v", filename)

	// Open file and read CSV records
	file, err := os.Open(results_filename_prefix + filename)
	if err != nil {
		log.Fatalf("Error while opening results file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	csvRecords, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error while reading results file: %v", err)
	}

	log.Printf("Converting CSV records to Results. Number of results: %v", len(csvRecords))

	// Convert CSV records to Results
	var results []Result
	for i, csvRecord := range csvRecords {
		homeGoals, err := strconv.Atoi(csvRecord[1])
		if err != nil {
			log.Fatalf("Home goals value couldn't be converted to int. Value: '%v', record index: %v", csvRecord[1], i)
		}

		awayGoals, err := strconv.Atoi(csvRecord[2])
		if err != nil {
			log.Fatalf("Away goals value couldn't be converted to int. Value: '%v', record index: %v", csvRecord[1])
		}

		results = append(results,
			Result{
				homeTeam:  csvRecord[0],
				awayTeam:  csvRecord[3],
				homeGoals: homeGoals,
				awayGoals: awayGoals,
			})
	}

	return results
}
