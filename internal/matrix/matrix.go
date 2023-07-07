package matrix

import (
	"fmt"

	triplet "github.com/yarruslan/search-parker-square/internal/triplet"
)

type Matrix [3]triplet.Triplet

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
