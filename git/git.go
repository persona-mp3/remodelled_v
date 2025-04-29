package git

import (
	"fmt"
	"errors"
	"time"
	"os"
	"reflect"
	"bufio"
	"strings"
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


type Format struct {
	Parent string
	HashId string
	Author string
	Email string
	CommitMsg string
	CommitedAt string
}
type Commits []Commit
type Branches []Branch

var commits Commits
var branches Branches
var headerPtr Header



// when we initalised a new git repository, we'd want to just set all the defaults
// git folders, head files, log folders and all
// although this might not be case for the real implementation of git, this will just be used instead
func Init(email string) (*Commit, error) {
	if len(email) < 3 {
		fmt.Println("[error] -> Author name too short")
		return nil, errors.New("Author Name too short")
	}

	gitRepo := Commit {
		Author: "persona-mp3",
		Email : email,
	}
	fmt.Printf("\nRepository initalised as %s with < %s >\n", gitRepo.Author, gitRepo.Email )

	init_folders()

	return &gitRepo, nil
}

func init_folders() {
	err := os.Mkdir("git_folder", 0777)
	handle_err(err)	

	mk_head()

	sub_folders := [3]string{"refs", "objects", "logs"}

	for _, sub_folder := range sub_folders {
		path := "git_folder/" + sub_folder
		fmt.Println("made 1")
		err := os.Mkdir(path, 0777)
		handle_err(err)
	}
	fmt.Println("\n[main folders initialised]")

	defer init_refs_folders()
	defer mk_log_folders()
}

func mk_head() {
	// the HEAD file in the root dir simply contains a the path to the current working branch 
	// it will show something like this:
	// ref: refs/heads/master
	// so anytime i do a new checkout this is all that happens:
	// a new branch is made in the refs/heads/ directory --> new_branch

	// -- now new_branch contains the latest commit sha-1 id, 
	// and is recorded in the logs dir as logs/heads/new_branch {commit details}

	// and now, the header now points to the new branch path, were the updates are made 
	// -- HEAD now shows this: ref: refs/heads/new_branch

	// upon a new commit on this branch, the new_branch will now update to this
	// -- refs/heads/new_branch --> 7d_q388374237
	// -- logs/refs/heads/new_branch --> "broke prod at 5am" 7d_q388374237 vibecodes < vibecoder@yahoo.com >
	// HEAD still remains the same

	file, err := os.OpenFile("git_folder/HEAD.txt", os.O_CREATE | os.O_RDWR, 0660)
	handle_err(err)	
	defer file.Close()

	// point to master branch
	n, err := fmt.Fprintf(file, "ref:refs/heads/master")
	handle_err(err)	

	fmt.Printf("[HEAD CREATED] -> %d written", n)
}

func init_refs_folders() {
	sub_folders := [3]string{"heads", "remotes", "tags"}

	for _, sub_folder := range sub_folders {
		path := "git_folder/refs/" + sub_folder 
		err := os.Mkdir(path, 0777)
		handle_err(err)	
	}

	fmt.Println("\n[refs sub folders initialised]")
}

func mk_log_folders() {
	// we just need to make the refs folder 
	err := os.Mkdir("git_folder/logs/refs", 0770)
	handle_err(err)

	f, err2 := os.OpenFile("git_folder/logs/HEAD.txt", os.O_APPEND | os.O_CREATE | os.O_RDWR,0660)
	handle_err(err2)
	defer f.Close()

	sub_folders := [3]string{"heads", "remotes", "tags"}

	for _, sub_folder := range sub_folders {
		path := "git_folder/logs/refs/" + sub_folder 
		err := os.Mkdir(path, 0777)
		handle_err(err)	
	}

	fmt.Println("\n[log sub folders initialised]")
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


	// active_branch := update_branch(hashId)
	path, active_branch, updated_commit := get_latest_commit(commit)
	update_branch(hashId, path)
	update_branch_log(updated_commit, active_branch)
	return 
}

// this function just returns the path to the current branch so it can be used to update the logs
// this will be further be referred to as "active_branch" anywhere outside this function
func update_branch(hashId HashId, path string) string{

	// crawl the path to the branch file now and write the commit id
	// branch_path := "git_folder/"	+ path
	branch, err := os.OpenFile(path, os.O_RDWR | os.O_TRUNC | os.O_CREATE, 0660)
	handle_err(err)
	defer branch.Close()

	// write the latest commit to the branch file
	n, err := fmt.Fprintf(branch, hashId.Id)
	handle_err(err)

	fmt.Printf("[total bytes written] -> %d", n)
	fmt.Println("[success] -> branch updated")

  return path 
}


func update_branch_log(commit Commit, active_branch string) {
	path := "git_folder/logs/refs/heads/" + active_branch 

	file, err := os.OpenFile(path, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0660)
	handle_err(err)	
	defer file.Close()

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


		handle_err(err)	
		n_total += n

	}

	// fmt.Printf("[update_log] -> %d bytes", n_total)
	fmt.Println("[updated_log] -> success")
}

// we might need to fully refactor update_branch and include the logic
/*
	i.  Get latest commit from the active branch(upon initiaition, might just set latest commit to 00000)
			if:
	ii. 	latest_commit == 0000000 then; parent commit.Parent === 000000000 
			else:
	iii.	commit.Parent = latest_commit
			fi 

	iv. Basic format logic	
*/

func get_latest_commit(commit Commit) (string, string, Commit){
	/* active branch is given to us by the HEAD.file */
	head, err:= os.OpenFile("git_folder/HEAD.txt", os.O_RDONLY, 0660)
	handle_err(err)
	defer head.Close()

	/* read the head file and get path to the active branch*/
	head_scanner := bufio.NewScanner(head)
	var active_branch string

	for head_scanner.Scan() {
		ref := head_scanner.Text()
		_, active_branch, _ = strings.Cut(ref, ":")
	}
	
	/*path = git_folder/ + refs/heads/{name_of_branch}*/
	path := "git_folder/" + active_branch


	/* opening the branch file and reading the latest commit id*/
	branch, err := os.OpenFile(path,  os.O_RDONLY , 0660)	
	handle_err(err)
	defer branch.Close()

	branch_scanner := bufio.NewScanner(branch)
	var latest_commit string
	var hashId HashId

	for branch_scanner.Scan() {
		latest_commit = branch_scanner.Text()
		fmt.Println("this is the latest commit on branch")
		fmt.Println(latest_commit)
	}

	/*logic for deciding parent*/

	if latest_commit != "00000000000000000"{
		hashId.Id = latest_commit 
		commit.Parent = &hashId
	} else {
		hashId.Id = "0000000000000000"
		commit.Parent = &hashId
	}

	_, active_branch, _ = strings.Cut(path, "heads/")
	// fmt.Println("\nTHE ACTIVE BRANCH")
	// fmt.Println(active_branch)
	
	return path, active_branch, commit

	// update_branch_log(commit, active_branch)

}

func handle_err(err error) {
	if err != nil {
		fmt.Println("[error] ->", err)
		panic(err)
		// return
	}
}


func Checkout(name string) {
	if len(name) < 2 {
		fmt.Println("[error] -> branch name too short")
		return
	}

	path := "git_folder/refs/heads/" + name
	temp_commit := "rb_453124124835"

	f, err := os.OpenFile(path, os.O_CREATE | os.O_RDWR, 0660)
	handle_err(err)
	defer f.Close()
	
	n, err := fmt.Fprintf(f, temp_commit)
	handle_err(err)

	fmt.Println("\n[branch] -> new branch made")
	fmt.Printf("[new bytes] -> %d", n)

}

// this function is always called whenever a ```git --checkout``` or 
// ```git --switchto``` is ran.
// it gets the name of the branch, and just overwrites what branch its pointing to 
// also using the same structure ref: refs/heads/name

func update_header(name string){
	path := "git_folder/HEAD.txt"
	// new_branch := path + name
	ref := "ref: refs/heads/" + name
	
	f, err := os.OpenFile(path, os.O_CREATE | os.O_RDWR | os.O_TRUNC , 0660)
	handle_err(err)
	defer f.Close()

	n, err := fmt.Fprintf(f, ref)
	handle_err(err)
	//
	fmt.Println("\n[header] -> header pointer has been updated")
	fmt.Printf("[new bytes] -> %d", n)

}

func (repo Commit) Checkout2(name string) {
	if len(name) < 2 {
		fmt.Println("[error] -> branch name too short")
		return
	}

	path := "git_folder/refs/heads/" + name
	temp_commit := uuid.New().String()
	
	
	// repo.Email = "just_edited to check if i can use it as a method"
	f, err := os.OpenFile(path, os.O_CREATE | os.O_RDWR, 0660)
	handle_err(err)
	defer f.Close()
	
	n, err := fmt.Fprintf(f, temp_commit)
	handle_err(err)

	fmt.Println("[branch] -> new branch made")
	fmt.Printf("[new bytes] -> %d", n)
	
	update_header(name)
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

func SwitchBranch(name string) {
	path := "git_folder/refs/heads/" + name
	// we'll need to walk the tree directory and rebuild the folder with contents
	// this will be developed when zlib is being introduced
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		fmt.Println("[error] -> switching to ", name)
		fmt.Println(err)
		return
	}

	fmt.Printf("[success] -> branch was found %+v", *file)
	update_header(name)
}

func Logs() {
	// this will be later configured to read from the "logs" folder in the git_folder
	// for testing, we are just using the text file
	file, err := os.Open("commits.txt")
	handle_err(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// i is the line counter	
	i := 0
	for scanner.Scan() {
		commit := scanner.Text()
		fmt.Println(commit)
		if i % 5 == 0 {
			fmt.Println("\n")
		}
		// fmt.Println(commit)
		i++
	}
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
