package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/yarruslan/search-parker-square/internal/matrix"
	triplet "github.com/yarruslan/search-parker-square/internal/triplet"
)

const threads = 1

//var mapLock = make(chan struct{}, threads) //TODO - what is a good way to get rid from global lock?

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

	generator := new(triplet.Generator).Init(start, end, progressStep, threads)
	wg := &sync.WaitGroup{}

	worker := func(tasklist chan []triplet.Triplet) {
		defer wg.Done()
		wg.Add(1)
		for task := range tasklist {
			var ret []fmt.Stringer
			generator.MapLock <- struct{}{} //mapLock <- struct{}{}
			squares := matrix.LookupSubset(task, searchType)
			for _, sq := range squares {
				diagonals := matrix.CountDiagonals(sq)
				if diagonals >= searchType {
					ret = append(ret, sq)
				}
			}
			<-generator.MapLock
			if len(ret) > 0 {
				res <- ret
			}
		}
	}

	for i := 0; i < threads; i++ {
		go worker(generator.Iterate())
	}

	defer func() {
		time.Sleep(time.Second * 1) //TODO - race here :(
		wg.Wait()
		close(res)
	}()
}
