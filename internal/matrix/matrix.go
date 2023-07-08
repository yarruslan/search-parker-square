package matrix

import (
	"fmt"
	"sync"
	"time"

	triplet "github.com/yarruslan/search-parker-square/internal/triplet"
)

type Matrix [3]triplet.Triplet

type Generator struct {
	//set          IndexedTriplets
	//index        []SumSquares
	//id           int // index[id] is the last returned sum
	bufferWindow triplet.SumSquares
	start        triplet.SumSquares
	goal         triplet.SumSquares
	iterator     chan []Matrix
	//MapLock      chan struct{}
	threads int
}

func (g *Generator) Init(start, goal, window triplet.SumSquares, readerThreads int) *Generator {
	//g.set = make(IndexedTriplets)
	g.bufferWindow = window
	g.start = start
	g.goal = goal
	g.threads = readerThreads
	//g.index = []SumSquares{}
	//g.MapLock = make(chan struct{}, readerThreads)
	//g.generate(start)
	return g
}

func (g *Generator) GetSquares() {

}

func (m Matrix) String() string {
	return m[0].String() + m[1].String() + m[2].String() + "(" + fmt.Sprint(m[0][0]+m[0][1]+m[0][2]) + ")"
}

// go through slice of equal squares and find match
func LookupSubset(set []triplet.Triplet, searchType int) []Matrix {

	set = triplet.FilterSubset(set, searchType)
	lenght := len(set)

	var result []Matrix
	for i := 0; i < lenght-2; i++ {
		for j := i + 1; j < lenght-1; j++ {
			if set[i].HasOverlap2(set[j]) {
				continue
			}
			for k := j + 1; k < lenght; k++ {
				if set[i].HasOverlap2(set[k]) || set[j].HasOverlap2(set[k]) { //i-j checked before
					continue
				}
				candidate := checkCandidate1(Matrix{set[i], set[j], set[k]})
				if (candidate != Matrix{}) {
					candidate := checkCandidate2(candidate)
					if (candidate != Matrix{}) {
						result = append(result, candidate)
					}
				}
			}
		}
	}
	return result
}

// true if has 1 vertical match
func checkCandidate1(x Matrix) Matrix {
	sum := x[0][0] + x[0][1] + x[0][2]
	//TODO change to switch. x[0][0] is common, unnecessary
	if x[0][0]+x[1][0]+x[2][0] == sum {
		return x
	}
	if x[0][0]+x[1][0]+x[2][1] == sum {
		x[2][0], x[2][1] = x[2][1], x[2][0]
		return x
	}
	if x[0][0]+x[1][0]+x[2][2] == sum {
		x[2][0], x[2][2] = x[2][2], x[2][0]
		return x
	}
	if x[0][0]+x[1][1]+x[2][0] == sum {
		x[1][0], x[1][1] = x[1][1], x[1][0]
		return x
	}
	if x[0][0]+x[1][1]+x[2][1] == sum {
		x[1][0], x[1][1] = x[1][1], x[1][0]
		x[2][0], x[2][1] = x[2][1], x[2][0]
		return x
	}
	if x[0][0]+x[1][1]+x[2][2] == sum {
		x[1][0], x[1][1] = x[1][1], x[1][0]
		x[2][0], x[2][2] = x[2][2], x[2][0]
		return x
	}
	if x[0][0]+x[1][2]+x[2][0] == sum {
		x[1][0], x[1][2] = x[1][2], x[1][0]
		return x
	}
	if x[0][0]+x[1][2]+x[2][1] == sum {
		x[1][0], x[1][2] = x[1][2], x[1][0]
		x[2][0], x[2][1] = x[2][1], x[2][0]
		return x
	}
	if x[0][0]+x[1][2]+x[2][2] == sum {
		x[1][0], x[1][2] = x[1][2], x[1][0]
		x[2][0], x[2][2] = x[2][2], x[2][0]
		return x
	}
	return Matrix{}
}

// true if has 2 vertical match
func checkCandidate2(x Matrix) Matrix {
	sum := x[0][0] + x[0][1] + x[0][2]
	if x[0][1]+x[1][1]+x[2][1] == sum {
		return x
	}
	if x[0][1]+x[1][1]+x[2][2] == sum {
		x[2][1], x[2][2] = x[2][2], x[2][1]
		return x
	}
	if x[0][1]+x[1][2]+x[2][1] == sum {
		x[1][1], x[1][2] = x[1][2], x[1][1]
		return x
	}
	if x[0][1]+x[1][2]+x[2][2] == sum {
		x[2][1], x[2][2] = x[2][2], x[2][1]
		x[1][1], x[1][2] = x[1][2], x[1][1]
		return x
	}
	return Matrix{}
}

func CountDiagonals(x Matrix) int {
	sum := x[0][0] + x[0][1] + x[0][2]
	nrDiagonals := 0
	if x[0][0]+x[1][1]+x[2][2] == sum {
		nrDiagonals++
	}
	if x[0][1]+x[1][2]+x[2][0] == sum {
		nrDiagonals++
	}
	if x[0][2]+x[1][0]+x[2][1] == sum {
		nrDiagonals++
	}
	if x[0][2]+x[1][1]+x[2][0] == sum {
		nrDiagonals++
	}
	if x[0][1]+x[1][0]+x[2][2] == sum {
		nrDiagonals++
	}
	if x[0][0]+x[1][2]+x[2][1] == sum {
		nrDiagonals++
	}

	return nrDiagonals
}

// Main logic for searching magic squares: Generate triplets, try to combine them, count diagonals
func (g *Generator) FindSquaresWithDiagonals(searchType int, res chan []fmt.Stringer) {

	generator := new(triplet.Generator).Init(g.start, g.goal, g.bufferWindow, g.threads)
	wg := &sync.WaitGroup{}

	worker := func(tasklist chan []triplet.Triplet) {
		defer wg.Done()
		wg.Add(1)
		for task := range tasklist {
			var ret []fmt.Stringer
			generator.MapLock <- struct{}{} //mapLock <- struct{}{}
			squares := LookupSubset(task, searchType)
			for _, sq := range squares {
				diagonals := CountDiagonals(sq)
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

	for i := 0; i < g.threads; i++ {
		go worker(generator.Iterate())
	}

	defer func() {
		time.Sleep(time.Second * 1) //TODO - race here :(
		wg.Wait()
		close(res)
	}()
}
