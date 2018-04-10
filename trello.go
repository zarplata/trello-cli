package main

import (
	"fmt"
	"os"
	"time"

	trello "github.com/adlio/trello"
	"github.com/briandowns/spinner"
	hierr "github.com/reconquest/hierr-go"
)

var (
	s = spinner.New(spinner.CharSets[35], 100*time.Millisecond)
)

func selectID(max int, message string) int {
	var ID int
	for {
		fmt.Fprintf(os.Stderr, "\n:: Enter %s ID? [%d]: ", message, ID)
		fmt.Scanln(&ID)
		if ID <= max {
			break
		}
	}
	return ID
}

func selectVar(value, message string) string {
	for {
		fmt.Fprintf(os.Stderr, "\n:: Enter %s string? [%s]: ", message, value)
		fmt.Scanln(&value)
		if len(value) > 0 {
			break
		}
	}
	return value
}

func printBoardLine() {
	fmt.Printf(
		"\n|%-3s+%-20s+%-40s|",
		"---",
		"--------------------",
		"----------------------------------------",
	)
}
func printListLine() {
	fmt.Printf(
		"\n|%-3s+%-61s|",
		"---",
		"-------------------------------------------------------------",
	)
}

func handleSetup(path string, config *Config) error {

	config.Trello.AppKey = selectVar(config.Trello.AppKey, "AppKey")
	config.Trello.Token = selectVar(config.Trello.Token, "Token")

	client := trello.NewClient(config.Trello.AppKey, config.Trello.Token)

	s.Start()

	member, err := client.GetToken(config.Trello.Token, trello.Defaults())
	if err != nil {
		return hierr.Errorf(err, "can't get author")
	}
	config.Trello.Member = member.IDMember

	s.Stop()

	config.Trello.Board, err = handleBoard(client, config)
	if err != nil {
		return err
	}

	config.Trello.List, err = handleList(client, config)
	if err != nil {
		return err
	}

	err = saveConfig(path, config)
	if err != nil {
		return err
	}

	return nil
}

func handleBoard(
	client *trello.Client,
	config *Config,
) (string, error) {

	s.Start()

	myTokenID, err := client.GetToken(config.Trello.Token, trello.Defaults())
	if err != nil {
		return "", hierr.Errorf(err, "can't get token ID")
	}

	member, err := client.GetMember(myTokenID.IDMember, trello.Defaults())
	if err != nil {
		return "", hierr.Errorf(err, "can't get member")
	}

	boards, err := member.GetBoards(trello.Defaults())
	if err != nil {
		return "", hierr.Errorf(err, "can't get boards")
	}

	s.Stop()

	printBoardLine()
	fmt.Printf("\n|%-3s|%-20s|%-40s|", "ID", "Board", "Team")
	printBoardLine()

	for boardID, board := range boards {
		team := ""
		if board.IdOrganization != "" {
			org, err := client.GetOrganization(
				board.IdOrganization,
				trello.Defaults(),
			)
			if err != nil {
				return "", hierr.Errorf(err, "can't get organization")
			}
			team = fmt.Sprintf("%s (%s)", org.DisplayName, org.Name)
		}
		fmt.Printf("\n|%-3d|%-20s|%-40s|", boardID, board.Name, team)
	}
	printBoardLine()

	ID := selectID(len(boards), "board")

	return boards[ID].ID, nil
}

func handleList(
	client *trello.Client,
	config *Config,
) (string, error) {

	s.Start()

	board, err := client.GetBoard(config.Trello.Board, trello.Defaults())
	if err != nil {
		return "", hierr.Errorf(err, "can't get your board")
	}

	lists, err := board.GetLists(trello.Defaults())
	if err != nil {
		return "", hierr.Errorf(err, "can't get lists")
	}

	s.Stop()

	printListLine()
	fmt.Printf("\n|%-3s|%-61s|", "ID", "List")
	printListLine()

	for listID, list := range lists {
		fmt.Printf("\n|%-3d|%-61s|", listID, list.Name)
	}
	printListLine()

	ID := selectID(len(lists), "list")

	return lists[ID].ID, nil
}

func handleAdd(
	config *Config,
	args map[string]interface{},
) error {

	if len(config.Trello.List) == 0 {
		return fmt.Errorf("first run -S and setup your config")
	}

	client := trello.NewClient(config.Trello.AppKey, config.Trello.Token)

	s.Start()

	name, _ := args["<name>"].(string)
	description, _ := args["<description>"].(string)

	card := trello.Card{
		Name:      name,
		Desc:      description,
		IDList:    config.Trello.List,
		IDMembers: []string{config.Trello.Member},
	}

	err = client.CreateCard(&card, trello.Defaults())
	if err != nil {
		return hierr.Errorf(err, "can't add card to list")
	}

	s.Stop()

	return nil
}
