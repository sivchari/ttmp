package main

import (
	"github.com/sivchari/ttmp"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(ttmp.Analyzer) }
