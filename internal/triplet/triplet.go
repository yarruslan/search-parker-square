package triplet

import (
	"fmt"
	"log"
	"math"
	"sort"
	//"github.com/yarruslan/search-parker-square/internal/square"
)

type IndexedTriplets map[Square][]Triplet
type Square int
type Triplet [3]Square
type Generator struct {
	set          IndexedTriplets
	index        []Square //sorted list of generated sums
	id           int      // index[id] is the last returned sum
	bufferWindow Square
	showProgress Square //TODO split to 2 classes and interface
	lastProgress Square //TODO split to 2 classes and interface
	goal         Square
	iterator     chan []Triplet
	MapLock      chan struct{}
	start        Square //TODO split to 2 classes and interface
}

const SearchPureMagic int = 2
const SearchSemiMagic int = 1
const SearchNoMagic int = 0
const SearchCube int = 100
const SearchCubeInSquares int = 101

func (t *Triplet) getRoot() (ret [3]int) {
	ret[0] = int(math.Sqrt(float64(t[0])))
	ret[1] = int(math.Sqrt(float64(t[1])))
	ret[2] = int(math.Sqrt(float64(t[2])))
	return
}

func (t *Triplet) String() string {
	return fmt.Sprint(t.getRoot())
}

func (g *Generator) Init(start, goal, window Square, readerThreads int) *Generator {
	g.set = make(IndexedTriplets)
	g.bufferWindow = window
	g.goal = goal
	g.index = []Square{}
	g.MapLock = make(chan struct{}, readerThreads)
	g.generate(start)
	g.updateIndex()
	g.start = start
	g.showProgress = window
	return g
}

/*
	func (g *Generator) minIndexed() Square {
		if len(g.index) > 0 {
			return g.index[0]
		}
		return 0
	}
*/
func (g *Generator) maxIndexed() Square {
	if len(g.index) > 0 {
		return g.index[len(g.index)-1]
	}
	return 0
}

func (g *Generator) Iterate() chan []Triplet {
	if g.iterator == nil {
		g.iterator = make(chan []Triplet)
		go func() {
			defer close(g.iterator)
			for ; g.index[g.id] <= g.goal; g.id = g.next() {
				g.iterator <- g.set[g.index[g.id]]
			}
		}()
	}
	return g.iterator
}

func (g *Generator) next() int {
	if len(g.index) == g.id+1 {
		log.Println("Processed sums up to: ", g.index[g.id])
		g.exclusiveLock()
		g.deleteProcessed()
		g.generate(g.index[g.id] + 1)
		g.updateIndex()
		g.releaseExclusiveLock()
		log.Println("Generated next portion up to:", g.maxIndexed())
	}
	return g.id + 1
}

func (g *Generator) IterateSquares() chan []Triplet { //TODO split to 2 classes and interface
	g.bufferWindow = 0
	if g.iterator == nil {
		g.iterator = make(chan []Triplet)
		go func() {
			defer close(g.iterator)
			for g.id = int(math.Sqrt(float64((g.start)))); Square(g.id*g.id) <= g.goal; g.nextSquare() {
				g.iterator <- g.set[Square(g.id*g.id)]
			}
		}()
	}
	return g.iterator
}

func (g *Generator) nextSquare() {

	if Square(g.id*g.id) > g.lastProgress+g.showProgress {
		log.Println("Processed sums up to: ", Square(g.id*g.id))
		g.lastProgress = Square(g.id * g.id)
	}
	g.exclusiveLock()
	//g.deleteProcessed()
	g.id++
	g.generate(Square(g.id * g.id)) //TODO almost single threaded. how to make it parallel?
	g.releaseExclusiveLock()
	//log.Println("Generated next portion up to:", g.maxIndexed())

}

func (g *Generator) updateIndex() {

	exist := make(map[Square]bool)
	for _, v := range g.index {
		exist[v] = true
	}
	for k := range g.set {
		//avoid duplicates
		if !exist[k] {
			g.index = append(g.index, k)
		}
	}
	sort.Slice(g.index, func(i, j int) bool {
		return g.index[i] < g.index[j]
	})

}

func (g *Generator) deleteProcessed() {
	if len(g.index) > g.id+1 {
		log.Fatal("Partial clean-up not implemented")
	}
	g.set = make(IndexedTriplets)
}

func (g *Generator) generate(windowLow Square) {
	windowHigh := windowLow + g.bufferWindow
	count := 0
	start := int(math.Floor(math.Sqrt(float64(windowLow) / 3)))
	stop := int(math.Ceil(math.Sqrt(float64(windowHigh))))

	for i := start; i < stop; i++ {
		//stop uncenessary cycles early
		for j := 1; j < i; j++ {
			if Square(i*i+j*j) > windowHigh {
				break
			}
			if !(Square(i*i+2*j*j) < windowLow) {
				for k := 0; k < j; k++ {
					sum := Square(i*i + j*j + k*k)
					if sum > windowHigh {
						break
					}
					if sum >= windowLow {
						g.set[sum] = append(g.set[sum], Triplet{Square(i * i), Square(j * j), Square(k * k)})
						count++
					}
				}
			}
		}
	}

}

func (a *Triplet) HasOverlap(b Triplet) bool {
	if a[0] == b[0] || a[0] == b[1] || a[0] == b[2] ||
		a[1] == b[0] || a[1] == b[1] || a[1] == b[2] ||
		a[2] == b[0] || a[2] == b[1] || a[2] == b[2] {
		return true
	}
	return false
}

func FilterSubset(in []Triplet, searchType int) []Triplet { //TODO split to different filters. Make it a method, or ref to func
	var minTriplets int
	switch searchType {
	case SearchPureMagic:
		minTriplets = 8
	case SearchSemiMagic:
		minTriplets = 7
	case SearchNoMagic:
		minTriplets = 6
	case SearchCube, SearchCubeInSquares:
		minTriplets = 27
	}
	if len(in) < minTriplets {
		return []Triplet{}
	}
	if searchType == SearchCubeInSquares && !is_square(in[0][0]+in[0][1]+in[0][2]) { //Empirical: all solutions of intersecting no-magic squares are when summ of 3 items is square itself. don't know why
		return []Triplet{}
	}

	keysStat := make(map[Square]int)
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
	if (searchType == SearchCube || searchType == SearchCubeInSquares) && stat3 < 27 {
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
		case SearchCube, SearchCubeInSquares:
			dimensions = 3
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

func (g *Generator) exclusiveLock() {
	for i := 0; i < cap(g.MapLock); i++ {
		g.MapLock <- struct{}{}
	}

}

func (g *Generator) releaseExclusiveLock() {
	for i := 0; i < cap(g.MapLock); i++ {
		<-g.MapLock
	}
}

func (g *Generator) GetSet() IndexedTriplets {
	return g.set
}

// Returns true if 2 triplets share same numbers
func (a *Triplet) Same(b *Triplet) bool {
	//012 021 102 120 210 201
	if (a[0] == b[0] && a[1] == b[1] && a[2] == b[2]) ||
		(a[0] == b[0] && a[1] == b[2] && a[2] == b[1]) ||
		(a[0] == b[1] && a[1] == b[0] && a[2] == b[2]) ||
		(a[0] == b[1] && a[1] == b[2] && a[2] == b[0]) ||
		(a[0] == b[2] && a[1] == b[1] && a[2] == b[0]) ||
		(a[0] == b[2] && a[1] == b[0] && a[2] == b[1]) {
		return true
	}
	return false
}

func is_square(n Square) bool {
	root := int(math.Sqrt(float64(n)))
	return Square(root*root) == n
}
