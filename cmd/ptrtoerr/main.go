package main

import (
	"github.com/makiuchi-d/ptrtoerr"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(ptrtoerr.Analyzer) }
