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
	checkout2 string
	switchto string
	logs string
}



func read_cli() CliCmds{
	var cliArgs CliCmds
	
	flag.StringVar(&cliArgs.init, "init", "", " init --> initialise new repo object < author :string >")
	flag.StringVar(&cliArgs.commit, "commit", "", " commit --> create a new commit < commitMsg :string >")
	flag.StringVar(&cliArgs.checkout, "checkout", "", " checkout --> create a new commit < branch :string >")
	flag.StringVar(&cliArgs.checkout2, "checkout2", "", " checkout2 --> create a new commit < branch :string >")
	flag.StringVar(&cliArgs.switchto, "switchto", "", " switchto --> create a new commit < branch :string >")
	flag.StringVar(&cliArgs.logs, "logs", "", " logs --> displays all logs and commit history")

	flag.Parse()
	return cliArgs
}

func handle_err(err error) {
	if err != nil {
		return
	}
}

func exe_cli(cliArgs CliCmds) {
	var repo 	git.Commit
	// var repo2 git.Commit
	// var err error

	switch {
	case cliArgs.init != "" :
		_, err := git.Init(cliArgs.init)
		handle_err(err)
		// repo = tempRepo
		// fmt.Println(repo)

	case cliArgs.commit != "" :
		git.CommitMsg(cliArgs.commit)

	case cliArgs.checkout != "" :
		git.Checkout(cliArgs.checkout)

	// still optional, if i should make it a  method
	// so i can just grab the latest clone the latst commit
	case cliArgs.checkout2 != "" :
		repo.Checkout2(cliArgs.checkout2)

	case cliArgs.switchto != "" :
		git.SwitchBranch(cliArgs.switchto)

	case cliArgs.logs == "" :
		git.Logs()

	default:
		fmt.Println("[error] -> invalid arguments passed in")
	}
}

func main() {
	// repo, _ := git.Init("pearledIvory", "pearledivory@studios.com")
	// repo.NCommit("feat: integrating cli features")

	cliArgs := read_cli()
	exe_cli(cliArgs)
}
