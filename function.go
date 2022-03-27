package function

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
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
	Data       []byte            `json:"data"`
	Attributes map[string]string `json:"attributes"`
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

type SubsetCards struct {
	Cards []Card
}

func ImportChecklist(ctx context.Context, m PubSubMessage) error {
	league := m.Attributes["league"]
	checklistUrl := m.Attributes["checklistUrl"]
	setName := m.Attributes["set"]

	fileExtension := strings.Split(checklistUrl, ".")[1]

	target := fmt.Sprintf("/tmp/%s.%s", setName, fileExtension)
	err := downloadFile(target, checklistUrl)

	if err != nil {
		return err
	}

	checklistData, err := parseChecklist(target)

	if err != nil {
		return err
	}

	err = writeToFirestore(league, setName, checklistData)

	if err != nil {
		return err
	}

	return nil
}

func downloadFile(target string, url string) error {
	resp, err := http.Get(url)

	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	out, err := os.Create(target)

	if err != nil {
		log.Fatalln(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	return err
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
		collectedRows = append(collectedRows, columns)
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

func writeToFirestore(league string, set string, data map[string][]Card) error {
	ctx := context.TODO()
	firestoreClient, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))

	if err != nil {
		return err
	}
	defer firestoreClient.Close()

	failed := false
	leagueCollection := firestoreClient.Collection(league)
	for subset, cards := range data {
		checklist := leagueCollection.Doc(set).Collection(subset).Doc("checklist")
		_, err := checklist.Create(ctx, SubsetCards{Cards: cards})

		if err != nil {
			fmt.Println(err)
			failed = true
		}
	}

	if failed {
		return errors.New("encountered a problem when writing to firestore")
	}

	return nil
}
