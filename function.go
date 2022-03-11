package importer

import (
	"context"
	"fmt"

	"github.com/xuri/excelize/v2"
)

type ColumnType int

const (
	Set = iota
	CardNumber
	Description
	TeamCity
	TeamName
	Rookie
	Auto
	Mem
	Serial
	Odds
	Point
)

type PubSubMessage struct {
	Data []byte `json:"data"`
}

type Card struct {
	Set        string
	CardNumber string
	Player     string
	TeamCity   string
	TeamName   string
	IsRookie   bool
	HasAuto    bool
	HasMem     bool
	NumberedTo string
	Odds       string
	Point      string
}

func ImportCheckList(ctx context.Context, m PubSubMessage) error {
	//checklistFile := string(m.Data)

	return nil
}

func parseChecklist(f string) (map[string][]Card, error) {
	file, err := excelize.OpenFile(f)

	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	rows, err := file.Rows("Checklist")

	if err != nil {
		return nil, err
	}

	var collectedRows [][]string

	for rows.Next() {
		columns, _ := rows.Columns()
		collectedRows = append(collectedRows, trimSliceSpace(columns))
	}

	m := mapRowsToStruct(collectedRows)

	return m, nil
}

func mapRowsToStruct(rows [][]string) map[string][]Card {
	m := make(map[string][]Card)

	for i, row := range rows {
		if i == 0 {
			continue // this is the header data
		}

		set := row[Set]
		cardList := m[set]
		card := new(Card)

		card.Set = row[Set]
		card.CardNumber = row[CardNumber]
		card.Player = row[Description]
		card.TeamCity = row[TeamCity]
		card.TeamName = row[TeamName]
		card.NumberedTo = row[Serial]
		card.Odds = row[Odds]
		card.Point = row[Point]

		if row[Rookie] == "" {
			card.IsRookie = false
		} else {
			card.IsRookie = true
		}

		if row[Auto] == "" {
			card.HasAuto = false
		} else {
			card.HasAuto = true
		}

		if row[Mem] == "" {
			card.HasMem = false
		} else {
			card.HasMem = true
		}

		cardList = append(cardList, *card)
		m[set] = cardList
	}

	return m
}

func trimSliceSpace(s []string) []string {
	for {
		if len(s) > 0 && s[len(s)-1] == "" {
			s = s[:len(s)-1]
		} else {
			break
		}
	}
	return s
}
