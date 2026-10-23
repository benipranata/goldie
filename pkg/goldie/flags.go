package goldie

import "flag"

var (
	updateFlag = flag.Bool("update", false, "write actual output to golden files")
	cleanFlag  = flag.Bool("clean", false, "with -update, remove existing golden files first")
)
