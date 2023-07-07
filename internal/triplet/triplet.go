package triplet

import (
	"fmt"
	"math"
	"sort"
)

type IndexedTriplets map[SumSquares][]Triplet
type SumSquares int
type Triplet [3]SumSquares
type Generator struct {
	set          IndexedTriplets
	index        []SumSquares
	id           int // index[id] is the last returned sum
	bufferWindow SumSquares
	goal         SumSquares
}

const SearchPureMagic int = 2
const SearchSemiMagic int = 1
const SearchNoMagic int = 0

func (t *Triplet) getRoot() (ret [3]int) {
	ret[0] = int(math.Sqrt(float64(t[0])))
	ret[1] = int(math.Sqrt(float64(t[1])))
	ret[2] = int(math.Sqrt(float64(t[2])))
	return
}

func (t *Triplet) String() string {
	return fmt.Sprint(t.getRoot())
}

func (g *Generator) Init(start, goal, window SumSquares) *Generator {
	g.set = make(IndexedTriplets)
	g.bufferWindow = window
	g.goal = goal
	g.index = []SumSquares{}
	g.generate(start)
	return g
}

func (g *Generator) minIndexed() SumSquares {
	if len(g.index) > 0 {
		return g.index[0]
	}
	return 0
}

func (g *Generator) maxIndexed() SumSquares {
	if len(g.index) > 0 {
		return g.index[len(g.index)-1]
	}
	return 0
}

func (g *Generator) Iterate() chan []Triplet {
	iterator := make(chan []Triplet)

	go func() {
		for ; g.index[g.id] < g.goal; g.id = g.next() {
			iterator <- g.set[g.index[g.id]]
		}
	}()
	return iterator
}

func (g *Generator) next() int {
	if len(g.index) == g.id {
		g.generate(g.index[g.id])
		g.updateIndex()
	}
	return g.id + 1
}

func (g *Generator) updateIndex() {
	return
}

func (g *Generator) generate(start SumSquares) {
	g.set, g.index, _ = Generate(g.set, g.index, start, start+g.bufferWindow)
}

func Generate(groups IndexedTriplets, index []SumSquares, windowLow, windowHigh SumSquares) (IndexedTriplets, []SumSquares, SumSquares) {
	//TODO not single responsibility. Refactor to object Generator with internal state
	count := 0
	start := int(math.Floor(math.Sqrt(float64(windowLow) / 3)))
	stop := int(math.Ceil(math.Sqrt(float64(windowHigh))))
	for i := start; i < stop; i++ {
		//stop uncenessary cycles early
		for j := 1; j < i; j++ {
			if SumSquares(i*i+j*j) > windowHigh {
				break
			}
			if !(SumSquares(i*i+2*j*j) < windowLow) {
				for k := 0; k < j; k++ {
					sum := SumSquares(i*i + j*j + k*k)
					if sum > windowHigh {
						break
					}
					if sum >= windowLow {
						groups[sum] = append(groups[sum], Triplet{SumSquares(i * i), SumSquares(j * j), SumSquares(k * k)})
						count++
					}
				}
			}
		}
	}

	exist := make(map[SumSquares]bool)
	for _, v := range index {
		exist[v] = true
	}
	for k := range groups {
		//avoid duplicates
		if !exist[k] {
			index = append(index, k)
		}
	}
	sort.Slice(index, func(i, j int) bool {
		return index[i] < index[j]
	})
	var maxValueInWindow SumSquares
	for _, v := range index {
		if v > windowHigh {
			break
		}
		maxValueInWindow = v
	}

	return groups, index, maxValueInWindow
}

func (a *Triplet) HasOverlap2(b Triplet) bool {
	if a[0] == b[0] || a[0] == b[1] || a[0] == b[2] ||
		a[1] == b[0] || a[1] == b[1] || a[1] == b[2] ||
		a[2] == b[0] || a[2] == b[1] || a[2] == b[2] {
		return true
	}
	return false
}

func FilterSubset(in []Triplet, searchType int) []Triplet {
	var minDiagonals int
	switch searchType {
	case SearchPureMagic:
		minDiagonals = 8
	case SearchSemiMagic:
		minDiagonals = 7
	case SearchNoMagic:
		minDiagonals = 6
	}
	if len(in) < minDiagonals {
		return []Triplet{}
	}
	keysStat := make(map[SumSquares]int)
	for _, t := range in {
		keysStat[t[0]]++
		keysStat[t[1]]++
		keysStat[t[2]]++
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

	//heuristic shortcut - 9 numbers repeat at least 2 times, and 4 of them at least 3 times
	if (searchType == SearchPureMagic || searchType == SearchSemiMagic || searchType == SearchNoMagic) && stat2 < 9 {
		return []Triplet{}
	}
	if (searchType == SearchPureMagic) && stat3 < 4 {
		return []Triplet{}
	}
	if (searchType == SearchSemiMagic) && stat3 < 2 {
		return []Triplet{}
	}

	var filtered []Triplet
	var dimensions int

	for _, t := range in {
		switch searchType {
		case SearchPureMagic, SearchSemiMagic, SearchNoMagic:
			dimensions = 2 //heuristic shortcut: each number should be in at least 2 triplets to be part of the square
		}
		if keysStat[t[0]] >= dimensions {
			filtered = append(filtered, t)
		}
	}
	if len(in) == len(filtered) {
		return in
	}
	return FilterSubset(filtered, searchType) //recursively apply filter until it doesn't change result

}
