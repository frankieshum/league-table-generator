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

type teamRating struct {
	Team    string
	Attack  float64
	Defence float64
}

func generateResultsFile() string {
	// TODO - return in string format, then write to file separately for easier testing

	// Read teams and ratings from file
	file, err := os.Open("results_generation/teams_and_ratings.json")
	if err != nil {
		log.Fatalf("Error while opening teams file: %v", err)
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error while reading teams file: %v", err)
	}
	var teamRatings []teamRating
	err = json.Unmarshal(bytes, &teamRatings)
	if err != nil {
		log.Fatalf("Error while parsing teams and ratings file: %v", err)
	}
	fmt.Printf("teamsAndRatings2: %v\n", teamRatings)

	// Read goal distribution
	file, err = os.Open("results_generation/goal_distribution.json")
	if err != nil {
		log.Fatalf("Error while opening goal_distribution file: %v", err)
	}
	defer file.Close()
	bytes, err = io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error while reading goal_distribution file: %v", err)
	}
	var goalDistribution []int
	err = json.Unmarshal(bytes, &goalDistribution)
	if err != nil {
		log.Fatalf("Error while parsing goal_distribution file: %v", err)
	}

	// Read bonus goals
	file, err = os.Open("results_generation/bonus_goals.json")
	if err != nil {
		log.Fatalf("Error while opening bonus_goals file: %v", err)
	}
	defer file.Close()
	bytes, err = io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error while reading bonus_goals file: %v", err)
	}
	var bonusGoals []int
	err = json.Unmarshal(bytes, &bonusGoals)
	if err != nil {
		log.Fatalf("Error while parsing bonus_goals file: %v", err)
	}

	rand.Seed(time.Now().UnixMilli())

	// Create CSV
	filename := fmt.Sprintf("results_%v.csv", time.Now().UnixMilli())
	newFile, err := os.Create(results_filename_prefix + filename)
	if err != nil {
		log.Fatalf("Error while creating file: %v", err)
	}
	defer newFile.Close()

	for i, tr := range teamRatings {
		// Copy remaining teams into new slice
		teamsCopy := make([]teamRating, len(teamRatings))
		copy(teamsCopy, teamRatings)
		remainingTeams := teamsCopy[i+1:]

		writer := csv.NewWriter(newFile)
		defer writer.Flush()

		// Home games
		for _, otr := range remainingTeams {
			record := []string{
				tr.Team,
				strconv.Itoa(generateGoal(tr, otr, &goalDistribution, &bonusGoals)),
				strconv.Itoa(generateGoal(otr, tr, &goalDistribution, &bonusGoals)),
				otr.Team,
			}
			writer.Write(record)
		}

		// Away games
		for _, otr := range remainingTeams {
			record := []string{
				otr.Team,
				strconv.Itoa(generateGoal(otr, tr, &goalDistribution, &bonusGoals)),
				strconv.Itoa(generateGoal(tr, otr, &goalDistribution, &bonusGoals)),
				tr.Team,
			}
			writer.Write(record)
		} // TODO error handling on write
	}
	return filename

	// typicalScores := []int{
	// 	0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 4, 4, 4, 5, 5, 6,
	// }

	// // Read teams and ratings from file
	// file, err := os.Open("results_generation/teams_and_ratings.json")
	// if err != nil {
	// 	log.Fatalf("Error while opening teams file: %v", err)
	// }
	// defer file.Close()

	// bytes, err := io.ReadAll(file)
	// if err != nil {
	// 	log.Fatalf("Error while reading teams file: %v", err)
	// }

	// var teamsAndRatings []teamAndRating
	// err = json.Unmarshal(bytes, &teamsAndRatings)
	// if err != nil {
	// 	log.Fatalf("Error while parsing teams and ratings file: %v", err)
	// }

	// rand.Seed(time.Now().UnixMilli())
	// // var sb strings.Builder

	// // Create CSV
	// filename := fmt.Sprintf("results_%v.csv", time.Now().UnixMilli())
	// newFile, err := os.Create(results_filename_prefix + filename)
	// if err != nil {
	// 	log.Fatalf("Error while creating file: %v", err)
	// }
	// defer newFile.Close()

	// for i, tr := range teamsAndRatings {
	// 	// Copy remaining teams into new slice
	// 	teamsCopy := make([]teamAndRating, len(teamsAndRatings))
	// 	copy(teamsCopy, teamsAndRatings)
	// 	remainingTeams := teamsCopy[i+1:]

	// 	writer := csv.NewWriter(newFile)
	// 	defer writer.Flush()

	// 	// Home games
	// 	for _, otr := range remainingTeams {
	// 		// // This team
	// 		// sb.WriteString(tr.Team)
	// 		// sb.WriteString(",")
	// 		// sb.WriteString(strconv.Itoa(int(math.Round(float64(rand.Intn(7)) * tr.Rating))))
	// 		// sb.WriteString(",")
	// 		// // Other team
	// 		// sb.WriteString(strconv.Itoa(int(math.Round(float64(rand.Intn(7)) * otr.Rating))))
	// 		// sb.WriteString(",")
	// 		// sb.WriteString(otr.Team)
	// 		// sb.WriteString(",\n")
	// 		record := []string{
	// 			tr.Team, strconv.Itoa(int(math.Round(float64(typicalScores[rand.Intn(len(typicalScores))]) * tr.Rating))), strconv.Itoa(int(math.Round(float64(typicalScores[rand.Intn(len(typicalScores))]) * otr.Rating))), otr.Team,
	// 		}
	// 		writer.Write(record)
	// 	}

	// 	// Away games
	// 	for _, otr := range remainingTeams {
	// 		// Other team
	// 		// sb.WriteString(otr.Team)
	// 		// sb.WriteString(",")
	// 		// sb.WriteString(strconv.Itoa(int(math.Round(float64(rand.Intn(7)) * otr.Rating))))
	// 		// sb.WriteString(",")
	// 		// // This team
	// 		// sb.WriteString(strconv.Itoa(int(math.Round(float64(rand.Intn(7)) * tr.Rating))))
	// 		// sb.WriteString(",")
	// 		// sb.WriteString(tr.Team)
	// 		// sb.WriteString(",\n")
	// 		record := []string{
	// 			otr.Team, strconv.Itoa(int(math.Round(float64(typicalScores[rand.Intn(len(typicalScores))]) * otr.Rating))), strconv.Itoa(int(math.Round(float64(typicalScores[rand.Intn(len(typicalScores))]) * tr.Rating))), tr.Team,
	// 		}
	// 		writer.Write(record)
	// 	} // TODO error handling on write
	// }
	// // fmt.Println(sb.String())
	// return filename
}

func generateGoal(t teamRating, ot teamRating, goalDistribution *[]int, bonusGoals *[]int) int {

	randomGoal := (*goalDistribution)[rand.Intn(len(*goalDistribution))]
	adjustedForAttack := math.Round(float64(randomGoal) * t.Attack)
	adjustedForDefence := math.Round(adjustedForAttack * (1 - ot.Defence))
	adjustedWithBonus := int(adjustedForDefence) + (*bonusGoals)[rand.Intn(len(*bonusGoals))]
	return adjustedWithBonus
}
