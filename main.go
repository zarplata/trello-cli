package main

import (
	"fmt"
	"os"

	docopt "github.com/docopt/docopt-go"
	hierr "github.com/reconquest/hierr-go"
)

var (
	version        = "[manual build]"
	err            error
	OutputFileMode os.FileMode = 0644
	Help                       = `Authorize trello.com.
Open https://trello.com/app-key and take Key.
Next click generate a Token and take him.`
)

func main() {

	usage := `trello-cli

Usage:
  trello-cli [options] <name> [<description>]
  trello-cli [options] -S

Options:
  -c --config <path>          Use specified configuration file 
                                [default: $HOME/.config/trello.conf].

Commands:
  -S                          Setup AppKey, token, board, list.

Add card options:
  <name>                      Card name.
  <description>               Card description (optionals).

Misc options:
  --version                   Show version.
  -h --help                   Show this screen.
`
	args, err := docopt.Parse(usage, nil, true, version, false)
	if err != nil {
		panic(hierr.Errorf(err, "can't parse docopt"))
	}

	path := os.ExpandEnv(args["--config"].(string))

	config, err := loadConfig(path)
	if err != nil {
		hierr.Fatalf(err, "problem with configuration, run setup")
	}

	switch {
	case args["-S"].(bool):
		fmt.Println(Help)
		err = handleSetup(path, config)
	default:
		err = handleAdd(config, args)
	}

	if err != nil {
		hierr.Errorf(err, "an error occurred during the process")
	}
}
