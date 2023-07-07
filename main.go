package main

import (
	"fmt"
	"log"

	"github.com/yarruslan/search-parker-square/internal/matrix"
	triplet "github.com/yarruslan/search-parker-square/internal/triplet"
)

const threads = 11

// const max int = 5000 //greatest number to put to square of squares
const startSearch triplet.SumSquares = 0
const endSearch triplet.SumSquares = 150000
const progressStep triplet.SumSquares = 100000

//const memoryTarget int = 10000 //TODO target amount of triplets in memory

func main() {

	//Start listener for results channel
	resultChan := make(chan []fmt.Stringer)
	go func() {
		for res := range resultChan {
			for _, sq := range res {
				fmt.Println("Square ", sq, " has 1 diagonals")
			}
		}
	}()
	//Start calculations
	findSquaresWithDiagonals(startSearch, endSearch, triplet.SearchSemiMagic, resultChan)

}

// Main logic for searching magic squares: Generate triplets, try to combine them, count diagonals
func findSquaresWithDiagonals(start, end triplet.SumSquares, searchType int, res chan []fmt.Stringer) {
	tasklist := make(chan triplet.SumSquares)
	mapLock := make(chan struct{}, threads)
	var index []triplet.SumSquares
	groupedTriplets := make(triplet.IndexedTriplets)
	groupedTriplets, index, completeSum := triplet.Generate(groupedTriplets, index, start, start+progressStep)
	worker := func() {
		for task := range tasklist {
			var ret []fmt.Stringer
			mapLock <- struct{}{}
			squares := matrix.LookupSubset(groupedTriplets[task], searchType)
			for _, sq := range squares {
				diagonals := matrix.CountDiagonals(sq)
				if diagonals >= searchType {
					ret = append(ret, sq)
				}
			}

			<-mapLock
			if len(ret) > 0 {
				res <- ret
			}
		}
	}
	workerCloser := func() {
		//at the end wait for completion and close. TODO convert to defer func
		for i := 0; i < threads; i++ {
			mapLock <- struct{}{}
		}
		close(res)
	}
	for i := 0; i < threads; i++ {
		go worker()
	}
	defer workerCloser()
	var progress = start
	for step := 0; progress < end; step++ { //TODO fix end condition. should be exact target
		sum := index[step]
		if sum > completeSum { //TODO should avoid logic here. Iterator giving tasks should do the job
			panic("Missing generated values for " + fmt.Sprint(sum))
		} //Should not happen
		if sum == completeSum {
			log.Println("Processed sums up to: ", sum)
			//steal locks from workers to prevend starting new jobs. Map is prone to read-write race.
			for i := 0; i < threads; i++ {
				mapLock <- struct{}{}
			}
			for i := progress; i < sum; i++ {
				//free processed triplets
				delete(groupedTriplets, i)
			}
			//generate more
			groupedTriplets, index, completeSum = triplet.Generate(groupedTriplets, index, completeSum, completeSum+progressStep)
			log.Println("Generated next portion up to:", completeSum)
			//release locks
			for i := 0; i < threads; i++ {
				<-mapLock
			}

		}
		progress = sum
		tasklist <- sum
	}

}
