package square

import (
	"fmt"
	"sort"
	"sync"

	triplet "github.com/yarruslan/search-parker-square/internal/triplet"
)

type Matrix [3]triplet.Triplet

type Iterator interface {
	Iterate() chan []triplet.Triplet
}

type Generator struct {
	threads int
	tg      Iterator
}

func (g *Generator) Init(tg Iterator, threads int) *Generator {
	g.threads = threads
	g.tg = tg
	return g
}

func (m Matrix) String() string {
	return m[0].String() + m[1].String() + m[2].String() + "(" + fmt.Sprint(m[0][0]+m[0][1]+m[0][2]) + ")"
}

// go through slice of equal squares and find match
func CombineTripletsToMatrixes(set []triplet.Triplet, searchType int) []Matrix { //TODO fix has 2 responsibilites, find square and check square

	set = triplet.FilterSubset(set, searchType)
	lenght := len(set)

	pairIndex := triplet.BuildPairIndex(set)

	var result []Matrix
	//choose 3 non overlaping triplets, and check if they make a square.
	for i := 0; i < lenght-2; i++ {
		for j := i + 1; j < lenght-1; j++ {
			if set[i].HasOverlap(set[j]) {
				continue
			}
			if !couldBeInSameSquare(set[i], set[j], pairIndex) {
				continue
			}
			for k := j + 1; k < lenght; k++ {
				if set[i].HasOverlap(set[k]) || set[j].HasOverlap(set[k]) { //i-j checked before
					continue
				}
				if !(couldBeInSameSquare(set[k], set[j], pairIndex) && couldBeInSameSquare(set[i], set[k], pairIndex)) {
					continue
				}
				candidate := mustHave1column(Matrix{set[i], set[j], set[k]})
				if (candidate != Matrix{}) {
					candidate := mustHaveAllColumns(candidate)
					if (candidate != Matrix{}) {
						result = append(result, candidate)
					}
				}
			}
		}
	}
	return result
}

func mustHave1column(x Matrix) Matrix {
	//Sum := x[0][0] + x[0][1] + x[0][2] //x[0][0] is redundant, drop it
	partialSum := x[0][1] + x[0][2]
	switch partialSum {
	case x[1][0] + x[2][0]:
	case x[1][0] + x[2][1]:
		x[2][0], x[2][1] = x[2][1], x[2][0]
	case x[1][0] + x[2][2]:
		x[2][0], x[2][2] = x[2][2], x[2][0]
	case x[1][1] + x[2][0]:
		x[1][0], x[1][1] = x[1][1], x[1][0]
	case x[1][1] + x[2][1]:
		x[1][0], x[1][1] = x[1][1], x[1][0]
		x[2][0], x[2][1] = x[2][1], x[2][0]
	case x[1][1] + x[2][2]:
		x[1][0], x[1][1] = x[1][1], x[1][0]
		x[2][0], x[2][2] = x[2][2], x[2][0]
	case x[1][2] + x[2][0]:
		x[1][0], x[1][2] = x[1][2], x[1][0]
	case x[1][2] + x[2][1]:
		x[1][0], x[1][2] = x[1][2], x[1][0]
		x[2][0], x[2][1] = x[2][1], x[2][0]
	case x[1][2] + x[2][2]:
		x[1][0], x[1][2] = x[1][2], x[1][0]
		x[2][0], x[2][2] = x[2][2], x[2][0]
	default:
		return Matrix{}
	}
	return x
}

// true if has 2 vertical match
func mustHaveAllColumns(x Matrix) Matrix {
	//sum := x[0][0] + x[0][1] + x[0][2] //x[0][1] is redundant
	partialSum := x[0][0] + x[0][2]
	switch partialSum {
	case x[1][1] + x[2][1]:
	case x[1][1] + x[2][2]:
		x[2][1], x[2][2] = x[2][2], x[2][1]
	case x[1][2] + x[2][1]:
		x[1][1], x[1][2] = x[1][2], x[1][1]
	case x[1][2] + x[2][2]:
		x[2][1], x[2][2] = x[2][2], x[2][1]
		x[1][1], x[1][2] = x[1][2], x[1][1]
	default:
		return Matrix{}
	}
	return x
}

func (m *Matrix) CountDiagonals() int {
	sum := m[0][0] + m[0][1] + m[0][2]
	nrDiagonals := 0
	if m[0][0]+m[1][1]+m[2][2] == sum {
		nrDiagonals++
	}
	if m[0][1]+m[1][2]+m[2][0] == sum {
		nrDiagonals++
	}
	if m[0][2]+m[1][0]+m[2][1] == sum {
		nrDiagonals++
	}
	if m[0][2]+m[1][1]+m[2][0] == sum {
		nrDiagonals++
	}
	if m[0][1]+m[1][0]+m[2][2] == sum {
		nrDiagonals++
	}
	if m[0][0]+m[1][2]+m[2][1] == sum {
		nrDiagonals++
	}

	return nrDiagonals
}

func (g *Generator) GenerateSquares(searchType int, res chan []fmt.Stringer) { //TODO fix has 2 responsibilites, find square and check square

	wg := &sync.WaitGroup{}
	worker := func(tasklist chan []triplet.Triplet) {
		defer wg.Done()
		for task := range tasklist {
			var ret []fmt.Stringer
			var squares []Matrix
			switch gen := g.tg.(type) {
			case *triplet.Generator:
				gen.MapLock <- struct{}{}
				squares = CombineTripletsToMatrixes(task, searchType)
				<-gen.MapLock
			case *triplet.SquareGenerator:
				squares = CombineTripletsToMatrixes(task, searchType)
			}
			squares = filter(squares, searchType)
			for _, sq := range squares {
				ret = append(ret, sq)
			}
			if len(ret) > 0 {
				res <- ret
			}
		}
	}

	for i := 0; i < g.threads; i++ {
		wg.Add(1)
		go worker(g.tg.Iterate())
	}

	defer func() {
		wg.Wait()
		close(res)
	}()
}

func filter(in []Matrix, searchType int) []Matrix {
	var out []Matrix

	var diagonals int
	if (searchType == triplet.SearchCube || searchType == triplet.SearchCubeInSquares) && len(in) < 9 {
		return []Matrix{}
	}
	switch searchType {
	case triplet.SearchPureMagic:
		diagonals = 2
	case triplet.SearchSemiMagic:
		diagonals = 1
	case triplet.SearchNoMagic, triplet.SearchCube, triplet.SearchCubeInSquares:
		diagonals = 0
	}
	uniqueSquares := make(map[string]struct{})
	for _, sq := range in {
		if _, exist := uniqueSquares[sq.KeyAsString()]; !exist && sq.CountDiagonals() >= diagonals {
			out = append(out, sq)
			uniqueSquares[sq.KeyAsString()] = struct{}{}
		}
	}
	var tripletUsageIndex = make(map[triplet.Triplet]int)
	if searchType == triplet.SearchCube || searchType == triplet.SearchCubeInSquares {
		var out2 []Matrix
		for _, sq := range out {
			addTripletsToIndex(sq, tripletUsageIndex)
		}
		for _, sq := range out {
			if checkSatisfiesIndex(sq, tripletUsageIndex, searchType) {
				out2 = append(out2, sq)
			}
		}
		out = out2
	}
	if len(in) != len(out) {
		out = filter(out, searchType) //repeat recursively
	}
	return out
}

// Returns true if one matrix can become another via by rearranging numbers
func (a *Matrix) Same(b *Matrix) bool {
	//lazy check - check 2 squares contain same set of numbers. Lasy is enough given that source is a magic square
	return a.KeyAsString() == b.KeyAsString()
}

// 2 Squares intersect if they have matching triplet
func (a *Matrix) Intersect(b *Matrix) bool {
	tripletsA := []triplet.Triplet{a[0], a[1], a[2], a.Column(0), a.Column(1), a.Column(2)}
	tripletsB := []triplet.Triplet{b[0], b[1], b[2], b.Column(0), b.Column(1), b.Column(2)}
	for _, t1 := range tripletsA {
		for _, t2 := range tripletsB {
			if t1.Same(&t2) {
				return true
			}
		}
	}
	return false
}

/*
	func (s *Matrix) rotate(up, right int) *Matrix {
		//TODO
		return &Matrix{}
	}
*/
func (s *Matrix) Transpose() Matrix {
	out := *s
	out[0][1], out[1][0] = out[1][0], out[0][1]
	out[0][2], out[2][0] = out[2][0], out[0][2]
	out[1][2], out[2][1] = out[2][1], out[1][2]
	return out
}

func (s *Matrix) Column(id int) triplet.Triplet {
	return triplet.Triplet{s[id][0], s[id][1], s[id][2]}
}

func (s *Matrix) Contains(t triplet.Triplet) bool {
	var c0, c1, c2 triplet.Triplet
	c0 = s.Column(0)
	c1 = s.Column(1)
	c2 = s.Column(2)
	if t.Same(&s[0]) || t.Same(&s[1]) || t.Same(&s[2]) ||
		t.Same(&c0) || t.Same(&c1) || t.Same(&c2) {
		return true
	}
	return false
}

func couldBeInSameSquare(a, b triplet.Triplet, index map[triplet.Pair]struct{}) bool {
	//Check each of a[0],a[1],a[2] have a pair with b.
	return numberCouldBeInSquareWithTriplet(a[0], b, index) && numberCouldBeInSquareWithTriplet(a[1], b, index) && numberCouldBeInSquareWithTriplet(a[2], b, index)
}

func numberCouldBeInSquareWithTriplet(num triplet.Square, tr triplet.Triplet, index map[triplet.Pair]struct{}) bool { //TODO don't like the naming
	//index pair is build with {bigger, smaller}
	var pair0 triplet.Pair = triplet.Pair{num, tr[0]}
	if pair0.A < pair0.B {
		pair0.A, pair0.B = pair0.B, pair0.A
	}
	var pair1 triplet.Pair = triplet.Pair{num, tr[1]}
	if pair1.A < pair1.B {
		pair1.A, pair1.B = pair1.B, pair1.A
	}
	var pair2 triplet.Pair = triplet.Pair{num, tr[2]}
	if pair2.A < pair2.B {
		pair2.A, pair2.B = pair2.B, pair2.A
	}
	//1 pair is enough
	if _, ok := index[pair0]; ok == true {
		return true
	}
	if _, ok := index[pair1]; ok == true {
		return true
	}
	if _, ok := index[pair2]; ok == true {
		return true
	}
	return false
}

func (s *Matrix) KeyAsString() string {
	numbers := []int{int(s[0][0]), int(s[0][1]), int(s[0][2]), int(s[1][0]), int(s[1][1]), int(s[1][2]), int(s[2][0]), int(s[2][1]), int(s[2][2])}
	sort.Ints(numbers)
	return fmt.Sprint(numbers)
}

func addTripletsToIndex(sq Matrix, index map[triplet.Triplet]int) {
	index[sq[0].Sorted()]++
	index[sq[1].Sorted()]++
	index[sq[2].Sorted()]++

	index[sq.Column(0).Sorted()]++
	index[sq.Column(1).Sorted()]++
	index[sq.Column(2).Sorted()]++
}

func checkSatisfiesIndex(sq Matrix, index map[triplet.Triplet]int, searchType int) bool {
	if searchType == triplet.SearchCubeInSquares || searchType == triplet.SearchCube {
		//each row and column is belong to at least 2 squares
		nrOfIndexed := 0
		if index[sq[0].Sorted()] >= 2 {
			nrOfIndexed++
		}
		if index[sq[1].Sorted()] >= 2 {
			nrOfIndexed++
		}
		if index[sq[2].Sorted()] >= 2 {
			nrOfIndexed++
		}
		if index[sq.Column(0).Sorted()] >= 2 {
			nrOfIndexed++
		}
		if index[sq.Column(1).Sorted()] >= 2 {
			nrOfIndexed++
		}
		if index[sq.Column(2).Sorted()] >= 2 {
			nrOfIndexed++
		}
		if nrOfIndexed == 6 {
			return true
		}
	}
	return false
}

// 2 square touch if they share a number
func (a *Matrix) Touch(b Matrix) bool {
	setA := []int{int(a[0][0]), int(a[0][1]), int(a[0][2]), int(a[1][0]), int(a[1][1]), int(a[1][2]), int(a[2][0]), int(a[2][1]), int(a[2][2])}
	setB := []int{int(b[0][0]), int(b[0][1]), int(b[0][2]), int(b[1][0]), int(b[1][1]), int(b[1][2]), int(b[2][0]), int(b[2][1]), int(b[2][2])}

	for numA := range setA {
		for numB := range setB {
			if numA == numB {
				return true
			}
		}
	}
	return false
}
