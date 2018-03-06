package main

import (
	"fmt"
	"os"
	"time"

	docopt "github.com/docopt/docopt-go"
	"github.com/kovetskiy/lorg"
	hierr "github.com/reconquest/hierr-go"
)

var (
	version        = "[manual build]"
	logger         *lorg.Log
	err            error
	OutputFileMode os.FileMode = 0644
	Help                       = `Authorize trello.com.
Open https://trello.com/app-key and take Key.
Next click generate a Token and take him.`
)

func main() {

	start := time.Now()

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
  -v --verbose                Print verbose messages.
  --version                   Show version.
  -h --help                   Show this screen.
`
	args, err := docopt.Parse(usage, nil, true, version, false)
	if err != nil {
		panic(hierr.Errorf(err, "can't parse docopt"))
	}

	mustSetupLogger(args["--verbose"].(bool))

	path := os.ExpandEnv(args["--config"].(string))

	config, err := loadConfig(path)
	if err != nil {
		hierr.Errorf(err, "problem with configuration")
	}

	switch {
	case args["-S"].(bool):
		fmt.Println(Help)
		err = handleSetup(path, config)
	default:
		err = handleAdd(config, args)
	}

	if err != nil {
		logger.Warning(err.Error())
	}

	elapsed := time.Since(start)
	logger.Debugf("execution time %s", elapsed)
}
