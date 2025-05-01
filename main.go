package main

import (
	"fmt"
	"flag"
	"remodelled/git"
)

type CliCmds struct {
	init string
	commit string
	checkout string
	switchto string
	logs bool
	allb bool
}



func read_cli() CliCmds{
	var cliArgs CliCmds
	
	flag.StringVar(&cliArgs.init, "init", "",  " init --> initialise new repo object < author :string >")
	flag.StringVar(&cliArgs.commit, "commit", "", " commit --> create a new commit < commitMsg :string >")
	flag.StringVar(&cliArgs.checkout, "checkout", "", " checkout --> create a new commit < branch :string >")
	flag.StringVar(&cliArgs.switchto, "switchto", "", " switchto --> create a new commit < branch :string >")
	flag.BoolVar(&cliArgs.logs, "logs", false,  " logs --> displays all logs and commit history")
	flag.BoolVar(&cliArgs.allb, "allb", false,  " all-branches --> shows all current branches")

	flag.Parse()
	return cliArgs
}

func handle_err(err error) {
	if err != nil {
		return
	}
}

func exe_cli(cliArgs CliCmds) {
	switch {
	case cliArgs.init != "" :
		_, err := git.Init(cliArgs.init)
		handle_err(err)

	case cliArgs.commit != "" :
		git.CommitMsg(cliArgs.commit)

	case cliArgs.checkout != "" :
		git.Checkout(cliArgs.checkout)


	case cliArgs.allb :
		git.AllBranches()

	case cliArgs.switchto != "" :
		git.SwitchTo(cliArgs.switchto)

	case cliArgs.logs :
		git.Logs()


	default:
		fmt.Println("[error] -> invalid arguments passed in")
	}
}

func main() {
	cliArgs := read_cli()
	exe_cli(cliArgs)
}
