package main

import (
	"fmt"
	"sort"
	"strings"
)

type Table map[string]*TableRecord

func (table *Table) updateTableWithResult(result Result) {
	table.ensureTeamInTable(result.homeTeam)
	table.ensureTeamInTable(result.awayTeam)

	homeTeam, awayTeam, homeGoals, awayGoals := result.homeTeam, result.awayTeam, result.homeGoals, result.awayGoals

	// Increment goal stats
	(*table)[homeTeam].goalsFor += homeGoals
	(*table)[homeTeam].goalsAgainst += awayGoals
	(*table)[awayTeam].goalsFor += awayGoals
	(*table)[awayTeam].goalsAgainst += homeGoals

	// Increment W/D/L
	if result.IsHomeWin() {
		(*table)[homeTeam].won++
		(*table)[awayTeam].lost++
		return
	}
	if result.IsAwayWin() {
		(*table)[awayTeam].won++
		(*table)[homeTeam].lost++
		return
	}
	(*table)[awayTeam].drawn++
	(*table)[homeTeam].drawn++
}

func (table *Table) ensureTeamInTable(team string) {
	if _, ok := (*table)[team]; !ok {
		(*table)[team] = &TableRecord{}
	}
}

func (table *Table) formatTable() string {
	// Column headings
	var sb strings.Builder
	sb.WriteString("| Pos | Team                    | P        | W        | D        | L        | GF       | GA       | GD       | Points   |\n")
	sb.WriteString("|-----|-------------------------|----------|----------|----------|----------|----------|----------|----------|----------|\n")

	// Table records
	teams := (*table).getSortedTeamNames()
	for i, team := range teams {
		sb.WriteString(fmt.Sprintf("|%-5v|%-25v|%-10v|%-10v|%-10v|%-10v|%-10v|%-10v|%-10v|%-10v|\n",
			i+1, team, (*table)[team].matchesPlayed(), (*table)[team].won, (*table)[team].drawn, (*table)[team].lost, (*table)[team].goalsFor,
			(*table)[team].goalsAgainst, (*table)[team].goalDifference(), (*table)[team].points()))
	}

	return sb.String()
}

func (table *Table) getSortedTeamNames() []string {
	// Get a slice of all team names
	var teams []string
	for k := range *table {
		teams = append(teams, k)
	}

	// Sort by points
	sort.SliceStable(teams, func(i, j int) bool {
		return (*table)[teams[i]].points() > (*table)[teams[j]].points()
	})
	// Sort by GD where points are equal
	sort.SliceStable(teams, func(i, j int) bool {
		return (*table)[teams[i]].points() == (*table)[teams[j]].points() &&
			(*table)[teams[i]].goalDifference() > (*table)[teams[j]].goalDifference()
	})
	// Sort by GF where points and GD are equal
	sort.SliceStable(teams, func(i, j int) bool {
		return (*table)[teams[i]].points() == (*table)[teams[j]].points() &&
			(*table)[teams[i]].goalDifference() == (*table)[teams[j]].goalDifference() &&
			(*table)[teams[i]].goalsFor > (*table)[teams[j]].goalsFor
	})

	return teams
}

type TableRecord struct {
	won, drawn, lost, goalsFor, goalsAgainst int
}

func (tr TableRecord) matchesPlayed() int {
	return tr.won + tr.drawn + tr.lost
}

func (tr TableRecord) goalDifference() int {
	return tr.goalsFor - tr.goalsAgainst
}

func (tr TableRecord) points() int {
	return (tr.won * 3) + (tr.drawn)
}
