package git

import (
	"fmt"
	"os"
	"log"
	"path/filepath"
	"io/fs"
	// "bufio"
)

/*
	*This function is called in git.go to help record all fileNames in an index file. So we'll make the index file upon ```init``` 

*/

func logErr(err error) {
	if err != nil {
		fmt.Println("[error]\n")
		log.Fatal(err)
	}
}

func WriteIndex() {
	/*
	* os.Stat returns fileIndfo
	* type fileInfo stuct {
		Name() string,
		Size() int,
		IsDir() bool,
		Mode() FileMode,
}
*/
	index, err := os.OpenFile("stalker.txt", os.O_RDWR, 0664)
	logErr(err)
	defer index.Close()


	files, err := os.ReadDir(".")
	logErr(err)
		
	for _, file := range files {
		f, err := os.Stat(file.Name())
		logErr(err)

		data := fmt.Sprintf("name -> %s | permission -> %v | isDir -> %t\n", f.Name(), f.Mode(), f.IsDir())
		_, err = fmt.Fprintf(index, data)
		logErr(err)
	}

	fmt.Println("done writing to stalker")
	
}

func writeIndex(data string) {

	index, err := os.OpenFile("stalker.txt", os.O_RDWR, 0664)
	logErr(err)
	defer index.Close()

		_, err = fmt.Fprintf(index, data)
		logErr(err)

}

func RecursiveCheck() {
	index, err := os.OpenFile("stalker.txt", os.O_RDWR | os.O_TRUNC, 0664)
	logErr(err)
	defer index.Close()
	
	workingDir := "."
	
	err = filepath.WalkDir(workingDir, func(path string, d os.DirEntry, err error) error {
		logErr(err)
		
		// skip actual git folders and git_folder
		if d.IsDir() && (d.Name() == ".git" || d.Name() == "git_folder") {
			return fs.SkipDir
		}
		
		// fmt.Println(path)
		if !d.IsDir() {
			i, _ := d.Info() // info implements the io interface
			fmt.Printf("		%s  |		%+v \n", path, i.Mode())
			mode := fmt.Sprintf("%s || %s \n", path, i.Mode())
			_, err = fmt.Fprintf(index, mode)
			logErr(err)
		}
		return nil
	})
	
	logErr(err)
	fmt.Println("\n [done]\n")
}
