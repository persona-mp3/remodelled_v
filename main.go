package main

import (
	"remodelled/git"
)

func main() {
	vujade, _ := git.Init("vujade")

	vujade.NCommit("stone in focus", "vu_32892412478")
	vujade.NCommit("firelord zuko", "fz_832412200421")
	// vujade.NCommit("avatar angg", "ag_7912365632281")
	
	vujade.CheckoutC("aphex_miles")	
	vujade.CheckoutC("soccah")	

	vujade.SwitchTo("aphex_miles")
	
	// git.CommitHistory()


}
