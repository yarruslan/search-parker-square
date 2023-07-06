package triplet

import (
	"fmt"
	"math"
	"sort"
)

type Triplet struct {
	S1 int
	S2 int
	S3 int
	//sum int
}
type IndexedTriplets map[SumSquares][]Triplet
type SumSquares int

func (t *Triplet) getRoot() (ret [3]int) {
	//ret = make([]int, 3)
	ret[0] = int(math.Sqrt(float64(t.S1)))
	ret[1] = int(math.Sqrt(float64(t.S2)))
	ret[2] = int(math.Sqrt(float64(t.S3)))
	return
}

func (t *Triplet) String() string {
	return fmt.Sprint(t.getRoot()) // "{" + fmt.Sprint(math.Sqrt(float64(t.S1))) + ", " + fmt.Sprint(math.Sqrt(float64(t.S2))) + ", " + fmt.Sprint(math.Sqrt(float64(t.S3))) + "}"
}

func Generate(groups IndexedTriplets, index []SumSquares, windowLow, windowHigh SumSquares) (IndexedTriplets, []SumSquares, SumSquares) {

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
						groups[sum] = append(groups[sum], Triplet{i * i, j * j, k * k})
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
