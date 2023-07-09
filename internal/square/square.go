package square

import (
	"fmt"
	"sync"

	triplet "github.com/yarruslan/search-parker-square/internal/triplet"
)

type Matrix [3]triplet.Triplet

type Generator struct {
	threads int
	tg      *triplet.Generator
}

func (g *Generator) Init(tg *triplet.Generator, threads int) *Generator {
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

	var result []Matrix
	for i := 0; i < lenght-2; i++ {
		for j := i + 1; j < lenght-1; j++ {
			if set[i].HasOverlap(set[j]) {
				continue
			}
			for k := j + 1; k < lenght; k++ {
				if set[i].HasOverlap(set[k]) || set[j].HasOverlap(set[k]) { //i-j checked before
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
			g.tg.MapLock <- struct{}{}
			squares := CombineTripletsToMatrixes(task, searchType)
			<-g.tg.MapLock
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
	if searchType == triplet.SearchCube && len(in) < 9 {
		return []Matrix{}
	}
	switch searchType {
	case triplet.SearchPureMagic:
		diagonals = 2
	case triplet.SearchSemiMagic:
		diagonals = 1
	case triplet.SearchNoMagic, triplet.SearchCube:
		diagonals = 0
	}
	for _, sq := range in {
		if sq.CountDiagonals() >= diagonals {
			out = append(out, sq)
		}
	}
	if len(in) != len(out) {
		out = filter(out, searchType) //repeat recursively
	}
	return out
}
