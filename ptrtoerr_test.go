package ptrtoerr_test

import (
	"testing"

	"github.com/makiuchi-d/ptrtoerr"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ptrtoerr.Analyzer, "a")
}