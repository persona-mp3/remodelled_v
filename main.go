package main

import (
	"fmt"
	"flag"
	"remodelled/git"
)

type CliCmds struct {
	init string
	commit string
}



func read_cli() CliCmds{
	var cliArgs CliCmds
	
	flag.StringVar(&cliArgs.init, "init", "", " init -> initialise new repo object < author :string >")
	flag.StringVar(&cliArgs.commit, "commit", "", " commit -> create a new commit < commitMsg :string >")

	flag.Parse()
	return cliArgs
}

func handle_err(err error) {
	if err != nil {
		return
	}
}

func exe_cli(cliArgs CliCmds) {
	var repo 	*git.Commit
	// var repo2 git.Commit
	// var err error

	switch {
	case cliArgs.init != "" :
		tempRepo, err := git.Init(cliArgs.init)
		handle_err(err)
		repo = tempRepo
		fmt.Println(repo)
	

	case cliArgs.commit != "" :
		// repo.NCommit(cliArgs.commit)
		git.CommitMsg(cliArgs.commit)

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
