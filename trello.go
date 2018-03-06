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
	s      = spinner.New(spinner.CharSets[35], 100*time.Millisecond)
	id int = 0
)

func selectId(max int, msg string) int {
	for {
		fmt.Fprintf(os.Stderr, "\n:: Enter %s id? [%d]: ", msg, id)
		fmt.Scanln(&id)
		if id <= max {
			break
		}
	}
	return id
}

func selectVar(value, msg string) string {
	for {
		fmt.Fprintf(os.Stderr, "\n:: Enter %s string? [%s]: ", msg, value)
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

	logger.Debugf("successfully get boards %#v", boards)

	printBoardLine()
	fmt.Printf("\n|%-3s|%-20s|%-40s|", "ID", "Board", "Team")
	printBoardLine()

	for i, board := range boards {
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
		fmt.Printf("\n|%-3d|%-20s|%-40s|", i, board.Name, team)
	}
	printBoardLine()

	id := selectId(len(boards), "board")

	if id > len(boards) {
		return "", fmt.Errorf("Incorrect Board Id.")
	}

	return boards[id].ID, nil
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

	logger.Debugf("successfully get lists %#v", lists)

	printListLine()
	fmt.Printf("\n|%-3s|%-61s|", "ID", "List")
	printListLine()

	for i, list := range lists {
		fmt.Printf("\n|%-3d|%-61s|", i, list.Name)
	}
	printListLine()

	id := selectId(len(lists), "list")

	if id > len(lists) {
		return "", fmt.Errorf("Incorrect List Id.")
	}

	return lists[id].ID, nil
}

func handleAdd(
	config *Config,
	args map[string]interface{},
) error {

	if len(config.Trello.List) == 0 {
		return fmt.Errorf("Firs run -S and setup your config")
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

	logger.Debugf(
		"successfuly add card '%s' with description '%s'",
		name,
		description,
	)

	return nil
}
