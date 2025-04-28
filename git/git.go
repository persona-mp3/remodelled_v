package git

import (
	"fmt"
	"errors"
	"time"
	"os"
	"reflect"
	// "encoding/json"

	"github.com/aquasecurity/table"
	"github.com/google/uuid"
)

type Commit struct {
	Parent *HashId `json:"parent"`
	Id HashId      `json:"id"`
	Author string  `json:"author"` 
	Email string   `json:"email"`
	// Snapshot *Tree 
	CommitMsg string     `json:"commitMsg"`
	CommitedAt time.Time `json:"commitedAt"`
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



func Init(email string) (*Commit, error) {
	if len(email) < 3 {
		fmt.Println("[error] -> Author name too short")
		return nil, errors.New("Author Name too short")
	}

	gitRepo := Commit {
		Author: "persona-mp3",
		Email : email,
	}
	// fmt.Printf("Repository initalised as %s with %s", author, email )
	fmt.Printf("\nRepository initalised as %s with < %s >\n", gitRepo.Author, gitRepo.Email )

	return &gitRepo, nil
}

func (commit *Commit) NCommit(msg string)  {
	if len(msg) < 2  {
		panic("[error] -> commit msg or id too short\n")
	}

	hashId := HashId {
		Id: uuid.New().String(),
	}

	commitedAt := time.Now()

	commit.CommitMsg = msg
	commit.CommitedAt = commitedAt
	commit.Id = hashId
	commit.Parent = &hashId
	fmt.Println("\nnew commit added")

	if len(branches) > 0 {
		// skip making a new master branch and pointer logic
		commit.Parent = &commits[len(commits)-1].Id
		commits = append(commits, *commit)
		// record_commit(commit)
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
	// record_commit(commit)

	return 
}

// for now, the commit is basically the same as the repo created
func CommitMsg(msg string) {
	if len(msg) < 2  {
		fmt.Println("[error] -> Commit message must be more than 2 characters")
		return
	}

	var commit Commit

	hashId := HashId {
		Id: uuid.New().String(),
	}

	commitedAt := time.Now()

	commit.Author = "persona-mp3"
	commit.Email = "peraledivory@studios.com"
	commit.CommitMsg = msg
	commit.CommitedAt = commitedAt
	commit.Id = hashId
	commit.Parent = &hashId
	fmt.Println("\nnew commit added")

	// if len(branches) > 0 {
	// 	// skip making a new master branch and pointer logic
	// 	commit.Parent = &commits[len(commits)-1].Id
	// 	commits = append(commits, *commit)
	// 	record_commit(commit)
	// 	return 
	// }

	// master := Branch {
	// 	Name: "master",
	// 	LatestCommit: &hashId,
	// }
	//

	// headerPtr.BranchName = master.Name
	// headerPtr.ActiveBranch = &master
	

	// branches = append(branches, master)
	// commits = append(commits, *commit)
	record_commit(commit)

	return 
}

func record_commit(commit Commit) {
	file, err := os.OpenFile("commits.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0660)
	checker(err)
	defer file.Close()


	type Format struct {
		Parent string
		HashId string
		Author string
		Email string
		CommitMsg string
		CommitedAt string
	}

	var formatted Format
	formatted.HashId = commit.Id.Id
	formatted.Parent = commit.Parent.Id
	formatted.Author = commit.Author
	formatted.Email = fmt.Sprintf("< %s >",commit.Email)
	formatted.CommitMsg = commit.CommitMsg
	formatted.CommitedAt = commit.CommitedAt.Format("Jan 2, 2006, 3:04 PM")
	
	types := reflect.TypeOf(formatted)
	values := reflect.ValueOf(formatted)
	
	
	n_total := 0
	for i :=0; i < types.NumField(); i++ {
		field := types.Field(i)
		value := values.Field(i).Interface()
		
		line := fmt.Sprintf("%s : %v\n", field.Name, value)
		n, err := fmt.Fprintf(file, line)

		checker(err)
		n_total += n

	}
	fmt.Printf("[recorded] -> %d bytes", n_total)
}

func checker(err error) {
	if err != nil {
		fmt.Println("[error] ->", err)
		return
	}
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

// TODO: update to get branches from file
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

// TODO: update to read from file instead
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
