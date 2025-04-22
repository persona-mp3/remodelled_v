package main

import (
	"fmt"
	"errors"
	"time"
	// "log"
)

type Commit struct {
	Parent *HashId
	Id HashId
	Author string
	// Snapshot *Tree 
	CommitMsg string
	CommitedAt time.Time
}

type Branch struct {
	Name string
	LatestCommit *HashId
}

type Header struct {
	BranchName string
	ActiveBranch *Branch
}

type HashId struct {
	Id string
}

type Commits []Commit
type Branches []Branch

var commits Commits
var branches Branches
var headerPtr Header

func Init(author string) (*Commit, error) {
	if len(author) < 3 {
		fmt.Println("[error] -> Author name too short")
		return nil, errors.New("Author Name too short")
	}

	gitRepo := Commit {
		Author: author,
	}

	return &gitRepo, nil
}

func (commit *Commit) NCommit(msg, id string)  {
	if len(msg) < 2 || len(id) < 8 {
		panic("[error] -> commit msg or id too short\n")
	}

	hashId := HashId {
		Id: id,
	}

	commitedAt := time.Now()
	commit.CommitMsg = msg
	commit.CommitedAt = commitedAt
	commit.Id = hashId
	commit.Parent = &hashId
	// fmt.Println("[success]\n")

	if len(branches) > 0 {
		// skip making a new master branch and pointer logic
		commit.Parent = &commits[len(commits)-1].Id
		commits = append(commits, *commit)
		return 
	}

	master := Branch {
		Name: "master",
		LatestCommit: &hashId,
	}

	headerPtr.BranchName = master.Name
	headerPtr.ActiveBranch = &master

	branches = append(branches, master)
	commits = append(commits, *commit)

	return 
}

func (commit *Commit) CheckoutC(name string) {
	newBranch := Branch {
		Name: name,
	}

	newBranch.LatestCommit = &commit.Id
	branches = append(branches, newBranch)
	fmt.Println("\n[switching] -> switching to ",name)

	HeaderPtr(name)

}

func (commit *Commit) SwitchTo(name string) {
	fmt.Println("\nSWITCHING BRANCH TO --> ", name)

	for _, branch := range branches {
		if branch.Name == name {
			headerPtr.ActiveBranch = &branch
			headerPtr.BranchName = branch.Name
			fmt.Println("[header] ->", name)
			return 
		}
	}

	panic("[error] -> branch not found\n") 
}


func HeaderPtr(name string) {
	// var newBranch Branch
	for _, branch := range branches {
		if branch.Name == name {
			headerPtr.ActiveBranch = &branch
			headerPtr.BranchName = branch.Name
			fmt.Println("[header] ->", name)
			return 
		}
	}

	fmt.Println("[error] -> error switching header to new branch")
	return
}


func main() {
	vujade, _ := Init("vujade")

	vujade.NCommit("stone in focus", "vu_32892412478")
	vujade.NCommit("firelord zuko", "fz_832412200421")
	vujade.NCommit("avatar angg", "ag_7912365632281")
	
	vujade.CheckoutC("aphex_miles")	
	vujade.CheckoutC("soccah")	

	vujade.SwitchTo("master")

	// fmt.Println("\n[branches]")
	// fmt.Println(branches)

	// fmt.Println("\n[[commits]]")
	// fmt.Printf("%+v",commits)
	
	// for _, c := range commits {
	// 	fmt.Println("[id] ->", c.Id)
	// 	fmt.Println("[parent] ->", *c.Parent, "\n")
	// }

}
