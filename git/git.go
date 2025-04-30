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

func handle_err(err error) {
	if err != nil {
		fmt.Println("[error] ->", err)
		panic(err)
	}
}

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

	file, err := os.OpenFile("git_folder/HEAD.txt", os.O_CREATE | os.O_RDWR, 0660)
	handle_err(err)	
	defer file.Close()

	// point to master branch by default
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


	path, active_branch, updated_commit := get_latest_commit(commit)
	update_branch(hashId, path)
	update_branch_log(updated_commit, active_branch)
	return 
}

/*
	path is returned by get_latest_commit()
	path = git_folder/refs/heads/{active branch}
	crawl the path to the branch file now and write the commit id
*/
func update_branch(hashId HashId, path string) {

	branch, err := os.OpenFile(path, os.O_RDWR | os.O_TRUNC | os.O_CREATE, 0660)
	handle_err(err)
	defer branch.Close()

	// write the latest commit to the branch file
	n, err := fmt.Fprintf(branch, hashId.Id)
	handle_err(err)

	fmt.Printf("[total bytes written] -> %d", n)
	fmt.Println("[success] -> branch updated")

}

/*
	commit, active_branch is returned by get_latest_commit()
	active_branch = master (default)
	commit is "updated commit" that has its parent
*/
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


/*
	i.  Get latest commit from the active branch
			if:
	ii. 	latest_commit -eq 0000000 or "" then; parent commit.Parent === 000000000 
			else:
	iii.	commit.Parent = latest_commit
			fi 

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
	
	return path, active_branch, commit

	
	/* we can now return these values to 
			update_branch(hashId, path)
			update_branch_log(commit, active_branch)
	*/ 
}

// this function is always called whenever a ```git --checkout``` or 
// ```git --switchto``` is ran.
// it gets the name of the branch, and just overwrites what branch its pointing to 
// also using the same structure ref: refs/heads/name

func update_header(name string){
	path := "git_folder/HEAD.txt"
	// new_branch := path + name
	ref := "ref:refs/heads/" + name
	
	f, err := os.OpenFile(path, os.O_CREATE | os.O_RDWR | os.O_TRUNC , 0660)
	handle_err(err)
	defer f.Close()

	n, err := fmt.Fprintf(f, ref)
	handle_err(err)
	
	fmt.Println("\n[header] -> header pointer has been updated")
	fmt.Printf("[status] --> %d written", n)
}

func  Checkout(name string) {
	if len(name) < 2 {
		fmt.Println("[error] -> branch name too short")
		return
	}

	path := "git_folder/refs/heads/" + name

	/*reading latest commit from active branch*/
	head, err := os.OpenFile("git_folder/HEAD.txt", os.O_RDONLY, 0660)
	handle_err(err)
	defer head.Close()
	
	var curr_branch string
	scanner := bufio.NewScanner(head)
	for scanner.Scan() {
		ref := scanner.Text()
		_, curr_branch, _ = strings.Cut(ref, ":")
	}
	
	to_curr := "git_folder/" + curr_branch
	f, err := os.OpenFile(to_curr, os.O_RDONLY, 0660)
	handle_err(err)
	defer f.Close()
	
	branch_scanner := bufio.NewScanner(f)	
	var latest_commit string

	for branch_scanner.Scan() {
		latest_commit = branch_scanner.Text()
	}
	
	/*creating new branch*/
	new_branch, err := os.OpenFile(path, os.O_CREATE | os.O_RDWR | os.O_TRUNC, 0660)
	handle_err(err)
	defer new_branch.Close()
	
	n, err := fmt.Fprintf(new_branch, latest_commit)
	handle_err(err)

	fmt.Printf("[success] -> new branch made\n  written %d bytes", n)
	update_header(name)
}

/* 
	use ReadDir to get all files in the current directory 
	func ReadDir(name string) ([]DirEntry, err)
	ii. Update header 	
	iii. Print out the latest commit && the log books
*/

func SwitchTo(name string) {
	is_found, branch := find_branch(name)

	if !is_found {
		return
	}

	fmt.Println("[success] -> branch found ", branch)
	update_header(branch)
	return
}

func find_branch(name string) (bool, string) {
	git_path := "git_folder/refs/heads/"
	fs, err := os.ReadDir(git_path)
	handle_err(err)

	is_found := false
	var branch string
	
	for _, f := range fs {
		if name == f.Name() {
			branch = f.Name()	
			is_found = true
			return is_found, branch
		}
	}

	fmt.Println("\n			[error] -> the branch was not found")
	fmt.Println("\n			try making a new branch with --checkout")
	return is_found, branch
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
