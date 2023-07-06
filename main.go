package main

import (
	"fmt"
	"log"
	"time"

	"github.com/yarruslan/search-parker-square/internal/matrix"
	triplet "github.com/yarruslan/search-parker-square/internal/triplet"
)

//type sumSquares triplet.SumSquares

const threads = 11

// const max int = 5000 //greatest number to put to square of squares
const startSearch triplet.SumSquares = 0
const endSearch triplet.SumSquares = 150000
const progressStep triplet.SumSquares = 100000

//const memoryTarget int = 10000 //TODO target amount of triplets in memory

func main() {
	//TODO refactor for better test coverage
	//separate IO and init from core
	resultChan := make(chan []fmt.Stringer)
	go func() {
		for res := range resultChan {
			for _, sq := range res {
				fmt.Println("Square ", sq, " has 1 diagonals")
			}
		}
	}()
	findSquaresWithDiagonals(startSearch, endSearch, 1, resultChan)

}

func findSquaresWithDiagonals(start, end triplet.SumSquares, d int, res chan []fmt.Stringer) {
	tasklist := make(chan triplet.SumSquares)
	mapLock := make(chan struct{}, threads)
	var index []triplet.SumSquares
	groupedTriplets := make(triplet.IndexedTriplets)
	groupedTriplets, index, completeSum := triplet.Generate(groupedTriplets, index, start, start+progressStep)
	worker := func() {
		for task := range tasklist {
			var ret []fmt.Stringer
			mapLock <- struct{}{}
			squares := matrix.LookupSubset(groupedTriplets[task])
			for _, sq := range squares {
				diagonals := matrix.CountDiagonals(sq)
				if diagonals >= d {
					ret = append(ret, sq)
				}
			}

			<-mapLock
			if len(ret) > 0 {
				res <- ret
			}
		}
	}
	for i := 0; i < threads; i++ {
		go worker()
	}
	var progress = start
	for step := 0; progress < end; step++ { //TODO fix end condition. should be exact target
		sum := index[step]
		if sum > completeSum {
			panic("Missing generated values for " + fmt.Sprint(sum))
		} //Should not happen
		if sum == completeSum {
			log.Println("Processing sum: ", sum, " Timestamp: ", time.Now())
			//wait till all workers finish task
			for i := 0; i < threads; i++ {
				mapLock <- struct{}{}
			}
			for i := progress; i < sum; i++ {
				//free processed
				delete(groupedTriplets, i)
			}
			//generate more
			groupedTriplets, index, completeSum = triplet.Generate(groupedTriplets, index, completeSum, completeSum+progressStep)
			log.Println("Generated next portion up to:", completeSum)
			//release lock
			for i := 0; i < threads; i++ {
				<-mapLock
			}

		}
		progress = sum
		tasklist <- sum
	}
	//at the end wait for completion and close
	for i := 0; i < threads; i++ {
		mapLock <- struct{}{}
	}
	close(res)
}
