package main

import (
	"fmt"
	"log"
	"math"
	"sort"
	"time"
)

type triplet struct {
	s1 int
	s2 int
	s3 int
	//sum int
}
type matrix struct {
	a triplet
	b triplet
	c triplet
}
type sumSquares int

const threads = 5

// const max int = 5000 //greatest number to put to square of squares
const startSearch sumSquares = 0
const endSearch sumSquares = 150000
const progressStep sumSquares = 100000

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

func generate(groups map[sumSquares][]triplet, index []sumSquares, windowLow, windowHigh sumSquares) (map[sumSquares][]triplet, []sumSquares, sumSquares) {

	count := 0
	start := int(math.Floor(math.Sqrt(float64(windowLow) / 3)))
	stop := int(math.Ceil(math.Sqrt(float64(windowHigh))))
	for i := start; i < stop; i++ {
		//stop uncenessary cycles early
		for j := 1; j < i; j++ {
			if sumSquares(i*i+j*j) > windowHigh {
				break
			}
			if !(sumSquares(i*i+2*j*j) < windowLow) {
				for k := 0; k < j; k++ {
					sum := sumSquares(i*i + j*j + k*k)
					if sum > windowHigh {
						break
					}
					if sum >= windowLow {
						groups[sum] = append(groups[sum], triplet{i * i, j * j, k * k})
						count++
					}
				}
			}
		}
	}

	exist := make(map[sumSquares]bool)
	for _, v := range index {
		exist[v] = true
	}
	for k, _ := range groups {
		//avoid duplicates
		if !exist[k] {
			index = append(index, k)
		}
	}
	sort.Slice(index, func(i, j int) bool {
		return index[i] < index[j]
	})
	var maxValueInWindow sumSquares
	for _, v := range index {
		if v > windowHigh {
			break
		}
		maxValueInWindow = v
	}

	return groups, index, maxValueInWindow
}

// go through slice of equal squares and find match
func lookupSubset(set []triplet) []matrix {
	lenght := len(set)
	if lenght < 7 { //broader search, including 1 diagonal
		return []matrix{}
	}

	//heuristic shortcut - 9 numbers repeat at least 2 times, and 4 of them at least 3 times
	keysStat := make(map[int]int)
	for _, v := range set {
		keysStat[v.s1]++
		keysStat[v.s2]++
		keysStat[v.s3]++
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
		if keysStat[set[i].s1] < 2 || keysStat[set[i].s2] < 2 || keysStat[set[i].s3] < 2 {
			continue //heuristic - skip triplet that does not have a cross
		}
		for j := i + 1; j < lenght-1; j++ {
			if keysStat[set[i].s1] < 2 || keysStat[set[i].s2] < 2 || keysStat[set[i].s3] < 2 {
				continue
			}
			if !noOverlap2(set[i], set[j]) {
				continue
			}
			for k := j + 1; k < lenght; k++ {
				if keysStat[set[i].s1] < 2 || keysStat[set[i].s2] < 2 || keysStat[set[i].s3] < 2 {
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

func noOverlap2(a, b triplet) bool {
	if a.s1 == b.s1 || a.s1 == b.s2 || a.s1 == b.s3 ||
		a.s2 == b.s1 || a.s2 == b.s2 || a.s2 == b.s3 ||
		a.s3 == b.s1 || a.s3 == b.s2 || a.s3 == b.s3 {
		return false
	}
	return true
}

func noOverlap3(a, b, c triplet) bool {
	//dumb comparison is faster
	/*	values := []int{a.s1, a.s2, a.s3,
			b.s1, b.s2, b.s3,
			c.s1, c.s2, c.s3}
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
	if /*a.s1 == b.s1 || a.s1 == b.s2 || a.s1 == b.s3 ||
	a.s2 == b.s1 || a.s2 == b.s2 || a.s2 == b.s3 ||
	a.s3 == b.s1 || a.s3 == b.s2 || a.s3 == b.s3 ||*/ //a & b were compared earlier
	a.s1 == c.s1 || a.s1 == c.s2 || a.s1 == c.s3 ||
		a.s2 == c.s1 || a.s2 == c.s2 || a.s2 == c.s3 ||
		a.s3 == c.s1 || a.s3 == c.s2 || a.s3 == c.s3 ||
		b.s1 == c.s1 || b.s1 == c.s2 || b.s1 == c.s3 ||
		b.s2 == c.s1 || b.s2 == c.s2 || b.s2 == c.s3 ||
		b.s3 == c.s1 || b.s3 == c.s2 || b.s3 == c.s3 {
		return false
	}
	return true
}

// true if has 1 vertical match
func checkCandidate1(x matrix) matrix {
	sum := x.a.s1 + x.a.s2 + x.a.s3
	if x.a.s1+x.b.s1+x.c.s1 == sum {
		return x
	}
	if x.a.s1+x.b.s1+x.c.s2 == sum {
		x.c.s1, x.c.s2 = x.c.s2, x.c.s1
		return x
	}
	if x.a.s1+x.b.s1+x.c.s3 == sum {
		x.c.s1, x.c.s3 = x.c.s3, x.c.s1
		return x
	}
	if x.a.s1+x.b.s2+x.c.s1 == sum {
		x.b.s1, x.b.s2 = x.b.s2, x.b.s1
		return x
	}
	if x.a.s1+x.b.s2+x.c.s2 == sum {
		x.b.s1, x.b.s2 = x.b.s2, x.b.s1
		x.c.s1, x.c.s2 = x.c.s2, x.c.s1
		return x
	}
	if x.a.s1+x.b.s2+x.c.s3 == sum {
		x.b.s1, x.b.s2 = x.b.s2, x.b.s1
		x.c.s1, x.c.s3 = x.c.s3, x.c.s1
		return x
	}
	if x.a.s1+x.b.s3+x.c.s1 == sum {
		x.b.s1, x.b.s3 = x.b.s3, x.b.s1
		return x
	}
	if x.a.s1+x.b.s3+x.c.s2 == sum {
		x.b.s1, x.b.s3 = x.b.s3, x.b.s1
		x.c.s1, x.c.s2 = x.c.s2, x.c.s1
		return x
	}
	if x.a.s1+x.b.s3+x.c.s3 == sum {
		x.b.s1, x.b.s3 = x.b.s3, x.b.s1
		x.c.s1, x.c.s3 = x.c.s3, x.c.s1
		return x
	}
	return matrix{}
}

// true if has 2 vertical match
func checkCandidate2(x matrix) matrix {
	sum := x.a.s1 + x.a.s2 + x.a.s3
	if x.a.s2+x.b.s2+x.c.s2 == sum {
		return x
	}
	if x.a.s2+x.b.s2+x.c.s3 == sum {
		x.c.s2, x.c.s3 = x.c.s3, x.c.s2
		return x
	}
	if x.a.s2+x.b.s3+x.c.s2 == sum {
		x.b.s2, x.b.s3 = x.b.s3, x.b.s2
		return x
	}
	if x.a.s2+x.b.s3+x.c.s3 == sum {
		x.c.s2, x.c.s3 = x.c.s3, x.c.s2
		x.b.s2, x.b.s3 = x.b.s3, x.b.s2
		return x
	}
	return matrix{}
}

/*
func sumSquare(a, b, c int) sumSquares {
	return a*a + b*b + c*c
}*/

func (t *triplet) String() string {
	return "{" + fmt.Sprint(math.Sqrt(float64(t.s1))) + ", " + fmt.Sprint(math.Sqrt(float64(t.s2))) + ", " + fmt.Sprint(math.Sqrt(float64(t.s3))) + "}"
}

func (m matrix) String() string {
	return m.a.String() + m.b.String() + m.c.String() + "(" + fmt.Sprint(m.a.s1+m.a.s2+m.a.s3) + ")"
}

func countDiagonals(x matrix) int {
	sum := x.a.s1 + x.a.s2 + x.a.s3
	nrDiagonals := 0
	if x.a.s1+x.b.s2+x.c.s3 == sum {
		nrDiagonals++
	}
	if x.a.s2+x.b.s3+x.c.s1 == sum {
		nrDiagonals++
	}
	if x.a.s3+x.b.s1+x.c.s2 == sum {
		nrDiagonals++
	}
	if x.a.s3+x.b.s2+x.c.s1 == sum {
		nrDiagonals++
	}
	if x.a.s2+x.b.s1+x.c.s3 == sum {
		nrDiagonals++
	}
	if x.a.s1+x.b.s3+x.c.s2 == sum {
		nrDiagonals++
	}

	return nrDiagonals
}

func findSquaresWithDiagonals(start, end sumSquares, d int, res chan []fmt.Stringer) {
	tasklist := make(chan sumSquares)
	mapLock := make(chan struct{}, threads)
	var index []sumSquares
	groupedTriplets := make(map[sumSquares][]triplet)
	groupedTriplets, index, completeSum := generate(groupedTriplets, index, start, start+progressStep)
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
			groupedTriplets, index, completeSum = generate(groupedTriplets, index, completeSum, completeSum+progressStep)
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
