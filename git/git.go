package git

import (
	"fmt"
	"errors"
	"time"
	"os"
	"encoding/json"
	"github.com/aquasecurity/table"
)

type Commit struct {
	Parent *HashId  				`json:"parent"`
	Id HashId 							`json:"id"`
	Author string 					`json:"author"` 
	// Snapshot *Tree 
	CommitMsg string				`json:"commitMsg"`
	CommitedAt time.Time		`json:"commitedAt"`
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

	f, _ := os.Create("commits.json")
	fmt.Print(f)

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
	fmt.Println("new commit added")

	if len(branches) > 0 {
		// skip making a new master branch and pointer logic
		commit.Parent = &commits[len(commits)-1].Id
		commits = append(commits, *commit)
		writeCommit(commit)
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
	// write commit to file

	return 
}

func writeCommit(commit *Commit)  {
	// file, err := os.Create("commits.json")
	// if err != nil {
	// 	panic(err)
	// }

	jsonData, err := json.Marshal(commit)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonData))
}

func (commit *Commit) CheckoutC(name string) {
	newBranch := Branch {
		Name: name,
	}

	newBranch.LatestCommit = &commit.Id
	branches = append(branches, newBranch)
	fmt.Println("\n[branching] -> branching to ",name)

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

func CommitHistory() {
	table := table.New(os.Stdout)
	table.SetHeaders("#id", "parent", "commitMsg","commitedAt")

	fmt.Println("\n[RENDERING COMMIT HISTORY]")

	for _, commit := range commits {
			
		parent := commit.Parent.Id
		commitedAt := commit.CommitedAt.Format("Jan 2, 2006, 3:04 PM")
		id := commit.Id.Id
		msg := commit.CommitMsg

		table.AddRow(id, parent, msg, commitedAt)
	}

	table.Render()
}
