package main

import (
	"flag"
	"fmt"

	"github.com/yarruslan/search-parker-square/internal/matrix"
	triplet "github.com/yarruslan/search-parker-square/internal/triplet"
)

//const memoryTarget int = 10000 //TODO target amount of triplets in memory, not window

func main() {

	startSearch, endSearch, progressStep, threads, searchType := getParametersFromFlags()

	//Start listener for results channel
	resultChan := make(chan []fmt.Stringer)
	go func() {
		for res := range resultChan {
			for _, sq := range res {
				fmt.Println("Square ", sq, " has 1 diagonals")
			}
		}
	}()

	new(matrix.Generator).Init(startSearch, endSearch, progressStep, threads).FindSquaresWithDiagonals(searchType, resultChan)

}

func getParametersFromFlags() (start, end, progress triplet.SumSquares, threads int, searchType int) {
	fStart := flag.Int("start", 1, "Sum of squares in line to start search")
	fEnd := flag.Int("end", 1000000, "Sum of squares in line to end search")
	fProgress := flag.Int("progress", 100000, "Report progress at section of this size")
	fThreads := flag.Int("threads", 11, "Number of go-routines performing calculations in parallel")
	fMode := flag.String("mode", "1diag", "Type of search \"0diag\"|\"1diag\"|\"2diag\" ") //Int("Starging value", 1, "Sum of squares in line to start search")
	flag.Parse()

	if fStart != nil {
		start = triplet.SumSquares(*fStart)
	} else {
		start = 0
	}
	if fEnd != nil {
		end = triplet.SumSquares(*fEnd)
	} else {
		end = 1000000
	}
	if fProgress != nil {
		progress = triplet.SumSquares(*fProgress)
	} else {
		progress = 100000
	}
	if fThreads != nil {
		threads = *fThreads
	} else {
		threads = 11
	}
	if fMode != nil {
		switch *fMode {
		case "0diag":
			searchType = triplet.SearchNoMagic
		case "1diag":
			searchType = triplet.SearchSemiMagic
		case "2diag":
			searchType = triplet.SearchPureMagic
		default:
			searchType = triplet.SearchSemiMagic
		}
	} else {
		searchType = triplet.SearchSemiMagic
	}

	return
}
