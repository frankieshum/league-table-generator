package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// TODO
// sort table by other criteria
// tidy formatting method
// tidy file processing methods
// refactoring
// tidy file/folder structure (separate files?)
// prompt for user input (specify csv or create new)

type TableRecord struct {
	played, won, drawn, lost, goalsFor, goalsAgainst, goalDifference, points int
}

func (tr TableRecord) calculateMatchesPlayed() int {
	return tr.won + tr.drawn + tr.lost
}

func (tr TableRecord) calculateGoalDifference() int {
	return tr.goalsFor - tr.goalsAgainst
}

func (tr *TableRecord) calculatePoints() {
	(*tr).points = ((*tr).won * 3) + ((*tr).drawn)
}

type Table map[string]*TableRecord

func main() {
	fileName := generateResultsFile()
	results := getResultsFromFile(fileName)
	table := Table{}

	for _, result := range results {
		ensureTeamInTable(result.homeTeam, table)
		ensureTeamInTable(result.awayTeam, table)
		incrementTableStats(result, table)
	}

	for _, record := range table {
		record.calculateMatchesPlayed()
		record.calculateGoalDifference()
		record.calculatePoints()
	}

	fmt.Println(formatTable(table))
}

func ensureTeamInTable(team string, table Table) {
	if _, ok := table[team]; !ok {
		table[team] = &TableRecord{}
	}
}

func incrementTableStats(result Result, table Table) {
	homeTeam, awayTeam, homeGoals, awayGoals := result.homeTeam, result.awayTeam, result.homeGoals, result.awayGoals

	// Increment goal stats
	table[homeTeam].goalsFor += homeGoals
	table[homeTeam].goalsAgainst += awayGoals
	table[awayTeam].goalsFor += awayGoals
	table[awayTeam].goalsAgainst += homeGoals

	// Increment W/D/L
	outcome := result.GetOutcome()
	if outcome == "HW" {
		table[homeTeam].won++
		table[awayTeam].lost++
		return
	}
	if outcome == "AW" {
		table[awayTeam].won++
		table[homeTeam].lost++
		return
	}
	table[awayTeam].drawn++
	table[homeTeam].drawn++
}

func formatTable(table Table) string {
	// Sort table by points
	var teams []string
	for k := range table {
		teams = append(teams, k)
	}
	sort.SliceStable(teams, func(i, j int) bool {
		return table[teams[i]].points > table[teams[j]].points
	})
	// TODO then sort by GD, GF etc

	// Format table as string
	// Column headings
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("|%-5v|%-25v|%-10v|%-10v|%-10v|%-10v|%-10v|%-10v|%-10v|%-10v|\n", "Pos", "Team", "P", "W", "D", "L", "GF", "GA", "GD", "Points"))

	// Separator between headings and body
	sb.WriteString(fmt.Sprintf("|%v|%v|%v|%v|%v|%v|%v|%v|%v|%v|\n", strings.Repeat("-", 5), strings.Repeat("-", 25), strings.Repeat("-", 10), strings.Repeat("-", 10),
		strings.Repeat("-", 10), strings.Repeat("-", 10), strings.Repeat("-", 10), strings.Repeat("-", 10), strings.Repeat("-", 10), strings.Repeat("-", 10)))

	// Table records
	for i, team := range teams {
		sb.WriteString(fmt.Sprintf("|%-5v|%-25v|%-10v|%-10v|%-10v|%-10v|%-10v|%-10v|%-10v|%-10v|\n",
			i+1, team, table[team].played, table[team].won, table[team].drawn, table[team].lost, table[team].goalsFor, table[team].goalsAgainst,
			table[team].goalDifference, table[team].points))
	}

	return sb.String()
}

type Result struct {
	homeTeam, awayTeam   string
	homeGoals, awayGoals int
}

func (r *Result) GetOutcome() string {
	if r.homeGoals > r.awayGoals {
		return "HW"
	}
	if r.awayGoals > r.homeGoals {
		return "AW"
	}
	return "D"
}
func getResultsFromFile(fileName string) []Result {
	// Open file and parse CSV records
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error while opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	csvRecords, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error while reading file: %v", err)
	}

	// Convert CSV records to Results
	results := []Result{}
	for _, csvRecord := range csvRecords {
		homeGoals, _ := strconv.Atoi(csvRecord[1])
		awayGoals, _ := strconv.Atoi(csvRecord[2])
		results = append(results, Result{
			homeTeam:  csvRecord[0],
			homeGoals: homeGoals,
			awayGoals: awayGoals,
			awayTeam:  csvRecord[3],
		})
	}

	return results
}

type TeamAndRating struct {
	Team   string
	Rating float64
}

func generateResultsFile() string {
	// TODO - return in string format, then write to file separately for easier testing

	typicalScores := []int{
		0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 4, 4, 4, 5, 5, 6,
	}

	// Read teams and ratings from file
	file, err := os.Open("teams_and_ratings.json")
	if err != nil {
		log.Fatalf("Error while opening teams file: %v", err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Error while reading teams file: %v", err)
	}

	teamsAndRatings := []TeamAndRating{}
	json.Unmarshal(bytes, &teamsAndRatings)

	rand.Seed(time.Now().UnixMilli())
	// var sb strings.Builder

	// Create CSV
	fileName := fmt.Sprintf("results_%v.csv", time.Now().UnixMilli())
	newFile, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Error while creating file: %v", err)
	}
	defer newFile.Close()

	for i, tr := range teamsAndRatings {
		// Copy remaining teams into new slice
		teamsCopy := make([]TeamAndRating, len(teamsAndRatings))
		copy(teamsCopy, teamsAndRatings)
		remainingTeams := teamsCopy[i+1:]

		writer := csv.NewWriter(newFile)
		defer writer.Flush()

		// Home games
		for _, otr := range remainingTeams {
			// // This team
			// sb.WriteString(tr.Team)
			// sb.WriteString(",")
			// sb.WriteString(strconv.Itoa(int(math.Round(float64(rand.Intn(7)) * tr.Rating))))
			// sb.WriteString(",")
			// // Other team
			// sb.WriteString(strconv.Itoa(int(math.Round(float64(rand.Intn(7)) * otr.Rating))))
			// sb.WriteString(",")
			// sb.WriteString(otr.Team)
			// sb.WriteString(",\n")
			record := []string{
				tr.Team, strconv.Itoa(int(math.Round(float64(typicalScores[rand.Intn(len(typicalScores))]) * tr.Rating))), strconv.Itoa(int(math.Round(float64(typicalScores[rand.Intn(len(typicalScores))]) * otr.Rating))), otr.Team,
			}
			writer.Write(record)
		}

		// Away games
		for _, otr := range remainingTeams {
			// Other team
			// sb.WriteString(otr.Team)
			// sb.WriteString(",")
			// sb.WriteString(strconv.Itoa(int(math.Round(float64(rand.Intn(7)) * otr.Rating))))
			// sb.WriteString(",")
			// // This team
			// sb.WriteString(strconv.Itoa(int(math.Round(float64(rand.Intn(7)) * tr.Rating))))
			// sb.WriteString(",")
			// sb.WriteString(tr.Team)
			// sb.WriteString(",\n")
			record := []string{
				otr.Team, strconv.Itoa(int(math.Round(float64(typicalScores[rand.Intn(len(typicalScores))]) * otr.Rating))), strconv.Itoa(int(math.Round(float64(typicalScores[rand.Intn(len(typicalScores))]) * tr.Rating))), tr.Team,
			}
			writer.Write(record)
		} // TODO error handling on write
	}
	// fmt.Println(sb.String())
	return fileName
}
