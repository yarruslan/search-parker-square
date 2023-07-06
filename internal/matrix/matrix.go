package matrix

import (
	"fmt"

	"github.com/yarruslan/search-parker-square/internal/triplet"
)

type Matrix struct {
	A triplet.Triplet
	B triplet.Triplet
	C triplet.Triplet
}

func (m Matrix) String() string {
	return m.A.String() + m.B.String() + m.C.String() + "(" + fmt.Sprint(m.A.S1+m.A.S2+m.A.S3) + ")"
}

// go through slice of equal squares and find match
func LookupSubset(set []triplet.Triplet) []Matrix {
	lenght := len(set)
	if lenght < 7 { //broader search, including 1 diagonal
		return []Matrix{}
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
		return []Matrix{}
	}

	var result []Matrix
	for i := 0; i < lenght-2; i++ {
		if keysStat[set[i].S1] < 2 || keysStat[set[i].S2] < 2 || keysStat[set[i].S3] < 2 {
			continue //heuristic - skip triplet that does not have a cross
		}
		for j := i + 1; j < lenght-1; j++ {
			if keysStat[set[i].S1] < 2 || keysStat[set[i].S2] < 2 || keysStat[set[i].S3] < 2 {
				continue
			}
			if !triplet.NoOverlap2(set[i], set[j]) {
				continue
			}
			for k := j + 1; k < lenght; k++ {
				if keysStat[set[i].S1] < 2 || keysStat[set[i].S2] < 2 || keysStat[set[i].S3] < 2 {
					continue
				}
				if !triplet.NoOverlap3(set[i], set[j], set[k]) {
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
	sum := x.A.S1 + x.A.S2 + x.A.S3
	if x.A.S1+x.B.S1+x.C.S1 == sum {
		return x
	}
	if x.A.S1+x.B.S1+x.C.S2 == sum {
		x.C.S1, x.C.S2 = x.C.S2, x.C.S1
		return x
	}
	if x.A.S1+x.B.S1+x.C.S3 == sum {
		x.C.S1, x.C.S3 = x.C.S3, x.C.S1
		return x
	}
	if x.A.S1+x.B.S2+x.C.S1 == sum {
		x.B.S1, x.B.S2 = x.B.S2, x.B.S1
		return x
	}
	if x.A.S1+x.B.S2+x.C.S2 == sum {
		x.B.S1, x.B.S2 = x.B.S2, x.B.S1
		x.C.S1, x.C.S2 = x.C.S2, x.C.S1
		return x
	}
	if x.A.S1+x.B.S2+x.C.S3 == sum {
		x.B.S1, x.B.S2 = x.B.S2, x.B.S1
		x.C.S1, x.C.S3 = x.C.S3, x.C.S1
		return x
	}
	if x.A.S1+x.B.S3+x.C.S1 == sum {
		x.B.S1, x.B.S3 = x.B.S3, x.B.S1
		return x
	}
	if x.A.S1+x.B.S3+x.C.S2 == sum {
		x.B.S1, x.B.S3 = x.B.S3, x.B.S1
		x.C.S1, x.C.S2 = x.C.S2, x.C.S1
		return x
	}
	if x.A.S1+x.B.S3+x.C.S3 == sum {
		x.B.S1, x.B.S3 = x.B.S3, x.B.S1
		x.C.S1, x.C.S3 = x.C.S3, x.C.S1
		return x
	}
	return Matrix{}
}

// true if has 2 vertical match
func checkCandidate2(x Matrix) Matrix {
	sum := x.A.S1 + x.A.S2 + x.A.S3
	if x.A.S2+x.B.S2+x.C.S2 == sum {
		return x
	}
	if x.A.S2+x.B.S2+x.C.S3 == sum {
		x.C.S2, x.C.S3 = x.C.S3, x.C.S2
		return x
	}
	if x.A.S2+x.B.S3+x.C.S2 == sum {
		x.B.S2, x.B.S3 = x.B.S3, x.B.S2
		return x
	}
	if x.A.S2+x.B.S3+x.C.S3 == sum {
		x.C.S2, x.C.S3 = x.C.S3, x.C.S2
		x.B.S2, x.B.S3 = x.B.S3, x.B.S2
		return x
	}
	return Matrix{}
}

func CountDiagonals(x Matrix) int {
	sum := x.A.S1 + x.A.S2 + x.A.S3
	nrDiagonals := 0
	if x.A.S1+x.B.S2+x.C.S3 == sum {
		nrDiagonals++
	}
	if x.A.S2+x.B.S3+x.C.S1 == sum {
		nrDiagonals++
	}
	if x.A.S3+x.B.S1+x.C.S2 == sum {
		nrDiagonals++
	}
	if x.A.S3+x.B.S2+x.C.S1 == sum {
		nrDiagonals++
	}
	if x.A.S2+x.B.S1+x.C.S3 == sum {
		nrDiagonals++
	}
	if x.A.S1+x.B.S3+x.C.S2 == sum {
		nrDiagonals++
	}

	return nrDiagonals
}
