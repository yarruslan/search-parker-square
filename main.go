package main

import (
	"fmt"
	"log"
	"time"

	triplet "github.com/yarruslan/search-parker-square/internal"
)

type matrix struct {
	a triplet.Triplet
	b triplet.Triplet
	c triplet.Triplet
}

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

// go through slice of equal squares and find match
func lookupSubset(set []triplet.Triplet) []matrix {
	lenght := len(set)
	if lenght < 7 { //broader search, including 1 diagonal
		return []matrix{}
	}

	//heuristic shortcut - 9 numbers repeat at least 2 times, and 4 of them at least 3 times
	keysStat := make(map[int]int)
	for _, v := range set {
		keysStat[v.S1]++
		keysStat[v.S2]++
		keysStat[v.S3]++
	}
	var stat2, stat3 int
	for _, v := range keysStat {
		if v >= 2 {
			stat2++
		}
		if v >= 3 {
			stat3++
		}
	}
	if stat2 < 9 || stat3 < 4 {
		return []matrix{}
	}

	var result []matrix
	for i := 0; i < lenght-2; i++ {
		if keysStat[set[i].S1] < 2 || keysStat[set[i].S2] < 2 || keysStat[set[i].S3] < 2 {
			continue //heuristic - skip triplet that does not have a cross
		}
		for j := i + 1; j < lenght-1; j++ {
			if keysStat[set[i].S1] < 2 || keysStat[set[i].S2] < 2 || keysStat[set[i].S3] < 2 {
				continue
			}
			if !noOverlap2(set[i], set[j]) {
				continue
			}
			for k := j + 1; k < lenght; k++ {
				if keysStat[set[i].S1] < 2 || keysStat[set[i].S2] < 2 || keysStat[set[i].S3] < 2 {
					continue
				}
				if !noOverlap3(set[i], set[j], set[k]) {
					continue
				}
				candidate := checkCandidate1(matrix{set[i], set[j], set[k]})
				if (candidate != matrix{}) {
					candidate := checkCandidate2(candidate)
					if (candidate != matrix{}) {
						result = append(result, candidate)
					}
				}

			}

		}
	}

	return result
}

func noOverlap2(a, b triplet.Triplet) bool {
	if a.S1 == b.S1 || a.S1 == b.S2 || a.S1 == b.S3 ||
		a.S2 == b.S1 || a.S2 == b.S2 || a.S2 == b.S3 ||
		a.S3 == b.S1 || a.S3 == b.S2 || a.S3 == b.S3 {
		return false
	}
	return true
}

func noOverlap3(a, b, c triplet.Triplet) bool {
	//dumb comparison is faster
	/*	values := []int{a.S1, a.S2, a.S3,
			b.S1, b.S2, b.S3,
			c.S1, c.S2, c.S3}
		duplicate := make(map[int]bool)
		for _, v := range values {
			_, exist := duplicate[v]
			if exist {
				return false
			}
			duplicate[v] = true
		}
		return true
	*/
	if /*a.S1 == b.S1 || a.S1 == b.S2 || a.S1 == b.S3 ||
	a.S2 == b.S1 || a.S2 == b.S2 || a.S2 == b.S3 ||
	a.S3 == b.S1 || a.S3 == b.S2 || a.S3 == b.S3 ||*/ //a & b were compared earlier
	a.S1 == c.S1 || a.S1 == c.S2 || a.S1 == c.S3 ||
		a.S2 == c.S1 || a.S2 == c.S2 || a.S2 == c.S3 ||
		a.S3 == c.S1 || a.S3 == c.S2 || a.S3 == c.S3 ||
		b.S1 == c.S1 || b.S1 == c.S2 || b.S1 == c.S3 ||
		b.S2 == c.S1 || b.S2 == c.S2 || b.S2 == c.S3 ||
		b.S3 == c.S1 || b.S3 == c.S2 || b.S3 == c.S3 {
		return false
	}
	return true
}

// true if has 1 vertical match
func checkCandidate1(x matrix) matrix {
	sum := x.a.S1 + x.a.S2 + x.a.S3
	if x.a.S1+x.b.S1+x.c.S1 == sum {
		return x
	}
	if x.a.S1+x.b.S1+x.c.S2 == sum {
		x.c.S1, x.c.S2 = x.c.S2, x.c.S1
		return x
	}
	if x.a.S1+x.b.S1+x.c.S3 == sum {
		x.c.S1, x.c.S3 = x.c.S3, x.c.S1
		return x
	}
	if x.a.S1+x.b.S2+x.c.S1 == sum {
		x.b.S1, x.b.S2 = x.b.S2, x.b.S1
		return x
	}
	if x.a.S1+x.b.S2+x.c.S2 == sum {
		x.b.S1, x.b.S2 = x.b.S2, x.b.S1
		x.c.S1, x.c.S2 = x.c.S2, x.c.S1
		return x
	}
	if x.a.S1+x.b.S2+x.c.S3 == sum {
		x.b.S1, x.b.S2 = x.b.S2, x.b.S1
		x.c.S1, x.c.S3 = x.c.S3, x.c.S1
		return x
	}
	if x.a.S1+x.b.S3+x.c.S1 == sum {
		x.b.S1, x.b.S3 = x.b.S3, x.b.S1
		return x
	}
	if x.a.S1+x.b.S3+x.c.S2 == sum {
		x.b.S1, x.b.S3 = x.b.S3, x.b.S1
		x.c.S1, x.c.S2 = x.c.S2, x.c.S1
		return x
	}
	if x.a.S1+x.b.S3+x.c.S3 == sum {
		x.b.S1, x.b.S3 = x.b.S3, x.b.S1
		x.c.S1, x.c.S3 = x.c.S3, x.c.S1
		return x
	}
	return matrix{}
}

// true if has 2 vertical match
func checkCandidate2(x matrix) matrix {
	sum := x.a.S1 + x.a.S2 + x.a.S3
	if x.a.S2+x.b.S2+x.c.S2 == sum {
		return x
	}
	if x.a.S2+x.b.S2+x.c.S3 == sum {
		x.c.S2, x.c.S3 = x.c.S3, x.c.S2
		return x
	}
	if x.a.S2+x.b.S3+x.c.S2 == sum {
		x.b.S2, x.b.S3 = x.b.S3, x.b.S2
		return x
	}
	if x.a.S2+x.b.S3+x.c.S3 == sum {
		x.c.S2, x.c.S3 = x.c.S3, x.c.S2
		x.b.S2, x.b.S3 = x.b.S3, x.b.S2
		return x
	}
	return matrix{}
}

/*
func sumSquare(a, b, c int) sumSquares {
	return a*a + b*b + c*c
}*/

func (m matrix) String() string {
	return m.a.String() + m.b.String() + m.c.String() + "(" + fmt.Sprint(m.a.S1+m.a.S2+m.a.S3) + ")"
}

func countDiagonals(x matrix) int {
	sum := x.a.S1 + x.a.S2 + x.a.S3
	nrDiagonals := 0
	if x.a.S1+x.b.S2+x.c.S3 == sum {
		nrDiagonals++
	}
	if x.a.S2+x.b.S3+x.c.S1 == sum {
		nrDiagonals++
	}
	if x.a.S3+x.b.S1+x.c.S2 == sum {
		nrDiagonals++
	}
	if x.a.S3+x.b.S2+x.c.S1 == sum {
		nrDiagonals++
	}
	if x.a.S2+x.b.S1+x.c.S3 == sum {
		nrDiagonals++
	}
	if x.a.S1+x.b.S3+x.c.S2 == sum {
		nrDiagonals++
	}

	return nrDiagonals
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
			squares := lookupSubset(groupedTriplets[task])
			for _, sq := range squares {
				diagonals := countDiagonals(sq)
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
