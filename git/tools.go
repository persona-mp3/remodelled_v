package git

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func logErr(err error, message string) {
	if err != nil {
		fmt.Println(message)
		log.Fatal(err)
	}
}

type FilePath struct {
	Name string
	HasParent bool
	ParentDir string
}
/*the purpose of this was because i could not get the files from funWalkDir to retrun as a slice of strings to write to from walkFn*/
var files []FilePath


func GitAdd(files string) {
	if files != "." {
		fmt.Println("user wants specific foldersn")
		return
	}

	fmt.Println("adding all files to staging area")
	// read all files in working directory and add write them to index file inside git folder
	i, err := os.OpenFile("git_folder/index.txt", os.O_CREATE | os.O_APPEND | os.O_RDWR, 0666)
	logErr(err, "Error occured in opening index file")
	i.Close()

	filepath.WalkDir(".", walkFn)
   
}


func factory(name string, hasParent bool, parent string) {
	file := FilePath {
		Name: name,
		HasParent: hasParent,
		ParentDir: parent,
	}

	files = append(files, file)
	i, err := os.OpenFile("git_folder/index.txt", os.O_APPEND | os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf( "Index file had not been made yet!\n")
		panic(err)
	}	

	defer i.Close()
	// info := fmt.Sprintf("%s\t\t\t\t%t\t\t\t%s\n", file.Name, file.HasParent, file.ParentDir)
	info := fmt.Sprintf("%-20s %-6t %-30s\n", file.Name, file.HasParent, file.ParentDir)
	n, err := fmt.Fprintf(i, info)
	logErr(err, "Error occured while writing to index file")
	fmt.Println("writing bytes ->", n)
}



func walkFn(path string, d os.DirEntry, err error ) error{
	logErr(err, "Error occured???")

	if d.Name() == "git_folder" || d.Name() == ".git" {
		return fs.SkipDir 
	}
	
	if !d.IsDir() {
		_, err := d.Info()
		logErr(err, "Error in getting file info")

		if path != d.Name() {
			factory(d.Name(), true, path)
		}else {
			factory(d.Name(), false, "root")
		}
		// fmt.Printf("f -> %s + %s\n", d.Name(), path)
	}
	return nil 
}

