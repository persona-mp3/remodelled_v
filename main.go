package main

import (
	"fmt"

	"remodelled/git"
)

func main() {
	vujade, _ := git.Init("vujade", "vujade@tokyo.com")

	// vujade.NCommit("stone in focus", "vu_32892412478")
	// vujade.NCommit("firelord zuko", "fz_832412200421")
	// vujade.NCommit("feat: ocean gate", "og_642803925")

	// vujade.NCommit("feat: writing commits", "wc_485893453")
	// vujade.NCommit("todo: folder handling", "fh_434775346894")
	// vujade.NCommit("todo: commit history from file ", "jd_3586376324")
	// vujade.NCommit("feat: updating hash id with")

	// vujade.NCommit("test: testing DAG 1")
	// vujade.NCommit("test: testing DAG, hashId of DAG1")
	// vujade.NCommit("test: testing DAG, hashId of DAG2")
	
	// vujade.CheckoutC("aphex_miles")	
	// vujade.CheckoutC("soccah")	
	//
	// vujade.SwitchTo("aphex_miles")

	fmt.Println("\n", vujade)
	git.CommitHistory()


}
