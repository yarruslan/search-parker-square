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

func NoOverlap2(a, b Triplet) bool {
	if a.S1 == b.S1 || a.S1 == b.S2 || a.S1 == b.S3 ||
		a.S2 == b.S1 || a.S2 == b.S2 || a.S2 == b.S3 ||
		a.S3 == b.S1 || a.S3 == b.S2 || a.S3 == b.S3 {
		return false
	}
	return true
}

func NoOverlap3(a, b, c Triplet) bool {
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
