package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const resultsGenerationFilePrefix string = "results_generation/"

type teamRating struct {
	Team    string
	Attack  float64
	Defence float64
}

func generateResultsToFile() string {
	log.Printf("Generating new result set to file")
	csvResults := generateResults()

	// Create CSV file
	log.Printf("Creating CSV file to store generated results")
	filename := fmt.Sprintf("results_%v.csv", time.Now().UnixMilli())
	newFile, err := os.Create(results_filename_prefix + filename)
	if err != nil {
		log.Fatalf("Error while creating results CSV file: %v", err)
	}
	defer newFile.Close()

	// Write records to CSV file
	writer := csv.NewWriter(newFile)
	defer writer.Flush()
	for _, record := range csvResults {
		if err := writer.Write(record); err != nil {
			log.Fatalf("Error while writing record to results CSV file: %v", err)
		}
	}
	log.Printf("Finished writing results records to CSV file. Filename: %v", filename)

	return filename
}

func generateResults() [][]string {
	teamRatings := getDataFromJsonFile[[]teamRating]("teams_and_ratings.json")
	goalDistribution := getDataFromJsonFile[[]int]("goal_distribution.json")
	bonusGoals := getDataFromJsonFile[[]int]("bonus_goals.json")

	rand.Seed(time.Now().UnixMilli())

	var csv [][]string

	log.Printf("Generating results for teams")
	for i, tr := range teamRatings {
		// Put other teams into new slice
		teamsCopy := make([]teamRating, len(teamRatings))
		copy(teamsCopy, teamRatings)
		otherTeamRatings := teamsCopy[i+1:]

		// Home games against other teams
		for _, otr := range otherTeamRatings {
			homeGoals := generateGoalsScoredForTeam(tr, otr, &goalDistribution, &bonusGoals)
			awayGoals := generateGoalsScoredForTeam(otr, tr, &goalDistribution, &bonusGoals)
			csvRow := formatResultAsCsvRecord(tr.Team, otr.Team, homeGoals, awayGoals)
			csv = append(csv, csvRow)
		}

		// Away games against other teams
		for _, otr := range otherTeamRatings {
			homeGoals := generateGoalsScoredForTeam(otr, tr, &goalDistribution, &bonusGoals)
			awayGoals := generateGoalsScoredForTeam(tr, otr, &goalDistribution, &bonusGoals)
			csvRow := formatResultAsCsvRecord(otr.Team, tr.Team, homeGoals, awayGoals)
			csv = append(csv, csvRow)
		}
	}
	return csv
}

func getDataFromJsonFile[T any](filename string) T {
	log.Printf("Getting data from file '%v'", filename)

	// Open file
	file, err := os.Open(resultsGenerationFilePrefix + filename)
	if err != nil {
		log.Fatalf("Error while opening %v file: %v", filename, err)
	}
	defer file.Close()

	// Read file
	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error while reading %v file: %v", filename, err)
	}

	// Parse JSON to objects
	var response T
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		log.Fatalf("Error while parsing %v file: %v", filename, err)
	}

	return response
}

func formatResultAsCsvRecord(homeTeam string, awayTeam string, homeGoals int, awayGoals int) []string {
	// Format: home team, home goals, away goals, away team
	return []string{
		homeTeam,
		strconv.Itoa(homeGoals),
		strconv.Itoa(awayGoals),
		awayTeam,
	}
}

func generateGoalsScoredForTeam(attackingTeam teamRating, defendingTeam teamRating, goalDistribution *[]int, bonusGoals *[]int) int {
	// Pick random int from goal distribution array (G)
	// Adjust using attacking team's attack rating (AR), then round
	// Adjust using defending team's defence rating (DR), then round
	// Add a random int from bonus goals array (B)
	// i.e. (G x AR) x (1 - DR) + B
	goals := (*goalDistribution)[rand.Intn(len(*goalDistribution))]
	adjustedForAttack := math.Round(float64(goals) * attackingTeam.Attack)
	adjustedForDefence := math.Round(adjustedForAttack * (1 - defendingTeam.Defence))
	addedBonusGoals := int(adjustedForDefence) + (*bonusGoals)[rand.Intn(len(*bonusGoals))]
	return addedBonusGoals
}
