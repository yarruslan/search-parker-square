package cube

import (
	"fmt"
	"log"

	"github.com/yarruslan/search-parker-square/internal/square"
	"github.com/yarruslan/search-parker-square/internal/triplet"
)

type Cube [3]square.Matrix

type Generator struct {
	sq *square.Generator
}

type Generator2 struct {
	sq *square.Generator
}

type GraphNode struct {
	square      square.Matrix
	connections map[*GraphNode]struct{}
}
type Graph []*GraphNode

func (c *Cube) String() string {
	return c[0].String() + c[1].String() + c[2].String()
}

func (g *Generator) Init(sq *square.Generator) *Generator {
	g.sq = sq
	return g
}

func (g *Generator2) Init(sq *square.Generator) *Generator2 {
	g.sq = sq
	return g
}

func (g *Generator) CombineSquaresToCubes(s []square.Matrix) []Cube {
	var result []Cube

	if graph := buildGraphOfSquares(s).filter(); graph.canContainCube() {
		result = append(result, graph.getCubes()...)
	}

	return result
}

func (g *Generator) GenerateCubes(searchType int, result chan []fmt.Stringer) {
	sqChan := make(chan []fmt.Stringer)

	go g.sq.GenerateSquares(searchType, sqChan)

	for squaresStringer := range sqChan {
		var squares []square.Matrix
		var res []fmt.Stringer
		for _, sqs := range squaresStringer {
			squares = append(squares, sqs.(square.Matrix))
		}
		cubes := g.CombineSquaresToCubes(squares)
		for _, c := range cubes {
			res = append(res, &c)
		}
		result <- res
	}
}

func (g *Generator2) GenerateCubes(searchType int, result chan []fmt.Stringer) {
	sqChan := make(chan []fmt.Stringer)

	go g.sq.GenerateSquares(searchType, sqChan)

	for squaresStringer := range sqChan {
		var squares []square.Matrix
		var res []fmt.Stringer
		for _, sqs := range squaresStringer {
			squares = append(squares, sqs.(square.Matrix))
		}
		cubes := g.CombineSquaresToCubes(squares)
		for _, c := range cubes {
			res = append(res, &c)
		}
		result <- res
	}

}

func (g *Generator2) CombineSquaresToCubes(in []square.Matrix) []Cube {
	//build index of pairs of triplets.
	//loop at filtered squares, and try to combine, using index of pairs
	var result []Cube

	index := make(map[struct{ t1, t2 string }]struct{})
	for _, sq := range in {
		addSqToIndex(sq, index)
	}

	graph := buildGraphOfSquares(in).filter()

	for _, node0 := range *graph {
		plane0 := node0.square
		//TODO pick only from 2nd level connections
		secondConnectionNodes := make(map[*GraphNode]struct{})
		for firstConnection, _ := range node0.connections {
			for secondConnection, _ := range firstConnection.connections {
				secondConnectionNodes[secondConnection] = struct{}{}
			}
		}

		for node1, _ := range secondConnectionNodes {
			plane1 := node1.square
			if plane0.Touch(plane1) {
				continue
			}
			if !(fitToIndex(plane0, plane1, index) || fitToIndex(plane0, plane1.Transpose(), index)) {
				continue
			}
			if fitToIndex(plane0, plane1.Transpose(), index) {
				plane1 = plane1.Transpose()
			}
			//align rows
			plane1 = translateRowsToMatch(plane1, plane0, index)
			if (plane1 == square.Matrix{}) {
				continue
			}
			fmt.Println("Aligned by rows:", plane0, plane1)
			//align columns
			plane1 = translateRowsToMatch(plane1.Transpose(), plane0.Transpose(), index)
			plane1 = plane1.Transpose()
			if (plane1 == square.Matrix{}) {
				continue
			}
			fmt.Println("2 planes matched:", plane0, plane1)
			for node2, _ := range secondConnectionNodes {
				plane2 := node2.square
				if plane0.Touch(plane2) || plane1.Touch(plane2) {
					continue
				}
				if !(fitToIndex(plane0, plane2, index) || fitToIndex(plane0, plane2.Transpose(), index)) || !(fitToIndex(plane1, plane2, index) || fitToIndex(plane1, plane2.Transpose(), index)) {
					continue
				}
				if fitToIndex(plane0, plane2.Transpose(), index) {
					plane2.Transpose()
				}
				if !fitToIndex(plane1, plane2, index) {
					continue
				}
				//align rows
				plane2 = translateRowsToMatch(plane2, plane0, index)
				if (plane2 == square.Matrix{}) {
					continue
				}
				//align columns
				plane2 = translateRowsToMatch(plane2.Transpose(), plane0.Transpose(), index)
				plane2 = plane2.Transpose()
				if (plane2 == square.Matrix{}) {
					continue
				}

				fmt.Println("Good candidate:", plane0, plane1, plane2)
				//TODO mandatory validate
				result = append(result, Cube{plane0, plane1, plane2})
			}
		}
	}
	//g.filter

	return result
}

// Graph is based on squares. If 2 squares are intersecting they are connected
func buildGraphOfSquares(in []square.Matrix) *Graph {
	var graph Graph
	for _, sq := range in {
		index := make(map[string]struct{})
		var newNode GraphNode
		newNode.square = sq
		newNode.connections = make(map[*GraphNode]struct{})
		if _, ok := index[newNode.square.KeyAsString()]; !ok {
			for _, node := range graph {
				if node.square.Intersect(&sq) {
					node.connect(&newNode)
				}
			}
			graph = append(graph, &newNode)
			index[newNode.square.KeyAsString()] = struct{}{}
		}
	}
	return &graph
}

func (g *Graph) canContainCube() bool { //TODO method is duplicated in filter?

	connectedNodes := 0
	for _, node := range *g {
		if len(node.connections) >= 6 { //each square shold connect to 6 other to form a cube
			connectedNodes++
		}
	}
	if connectedNodes >= 9 {
		fmt.Println("Cube candidate", connectedNodes, (*g)[0].square.String())
		return true
	}
	return false
}

func (g *Graph) getCubes() []Cube {
	/*
		graph contains at least 9 tightly connected squares
		1. Pick a square as base
		2. if 6 of its connected squares have another common intersection square - there is could be a cube
		2.1 define 3rd square
		3. Repeat at 1.
	*/
	for _, base := range *g {
		//base := node.square
		connectionsStats := make(map[*GraphNode]int)
		for verticalPlane, _ := range base.connections {
			for horisontalPlane, _ := range verticalPlane.connections {
				connectionsStats[horisontalPlane]++
			}
		}
		for secondPlane, v := range connectionsStats {
			if secondPlane != base && v >= 6 {
				result := buildCubeBy2Planes(base.square, secondPlane.square, base.connections)
				if (result != Cube{}) {
					fmt.Println("!!!!!!!!! Found cube:", result, "!!!!!!!!!") //TODO move that to return
					log.Println("!!!!!!!!! Found cube:", result, "!!!!!!!!!")
				}
			}
		}
	}

	return []Cube{}
}

func (a *GraphNode) connect(b *GraphNode) {
	a.connections[b] = struct{}{}
	b.connections[a] = struct{}{}
}

func (a *GraphNode) disconnect(b *GraphNode) {
	delete(a.connections, b)
	delete(b.connections, a)
}

// filter out weakly connected nodes
func (g *Graph) filter() *Graph {
	var result Graph
	for _, node := range *g {
		if len(node.connections) >= 6 {
			result = append(result, node)
		} else {
			for conn, _ := range node.connections {
				node.disconnect(conn)
			}
		}
	}
	if len(*g) == len(result) {
		return &result
	}
	return (&result).filter()

}

func buildCubeBy2Planes(base, secondLayer square.Matrix, verticalPlanes map[*GraphNode]struct{}) Cube {
	/*
		take 1st row of base
		find a plane containing it and intersecting second layer
		take 2nd row of base
		find a plane containing it, and intersecting second layer, and it should be parallel to previous vertical plane on second layer (maybe last is not necessary)
		choose 3rd plane by 3rd row and remaining row in second layer
		validate the 3 layers make proper cube.
	*/
	var firstPlane, secondPlane, thirdPlane square.Matrix
	var result []Cube
	var notifiedOnce bool

	for plane0, _ := range verticalPlanes {
		if plane0.square.Contains(base[0]) && plane0.square.Intersect(&secondLayer) {
			firstPlane = plane0.square
			//log.Println("Getting closer", firstPlane)
			//fmt.Println("Getting closer", firstPlane)
			for plane1, _ := range verticalPlanes {
				if plane1.square.Contains(base[1]) && plane1.square.Intersect(&secondLayer) { // && intersection uses different than first plane.
					secondPlane = plane1.square
					if !notifiedOnce {
						log.Println("Getting hot", firstPlane, secondPlane)
						fmt.Println("Getting hot", firstPlane, secondPlane)
						notifiedOnce = true
					}
					for plane2, _ := range verticalPlanes {
						if plane2.square.Contains(base[2]) && plane2.square.Intersect(&secondLayer) {
							thirdPlane = plane2.square
							result = append(result, Cube{firstPlane, secondPlane, thirdPlane})
						}
					}
				}
			}
		}
	}
	if len(result) > 1 {
		log.Println("Unexpected")
		fmt.Println("Unexpected")
	}
	if len(result) > 0 {
		log.Println("Result", result)
		fmt.Println("Result", result)
	}
	if len(result) > 0 {
		return result[0]
	}
	return Cube{}
}

func addSqToIndex(sq square.Matrix, index map[struct{ t1, t2 string }]struct{}) {
	var tr1, tr2 string
	tr1 = fmt.Sprint(sq[0].Sorted())
	tr2 = fmt.Sprint(sq[1].Sorted())
	if tr1 > tr2 {
		tr1, tr2 = tr2, tr1
	}
	index[struct {
		t1 string
		t2 string
	}{tr1, tr2}] = struct{}{}

	tr1 = fmt.Sprint(sq[0].Sorted())
	tr2 = fmt.Sprint(sq[2].Sorted())
	if tr1 > tr2 {
		tr1, tr2 = tr2, tr1
	}
	index[struct {
		t1 string
		t2 string
	}{tr1, tr2}] = struct{}{}

	tr1 = fmt.Sprint(sq[1].Sorted())
	tr2 = fmt.Sprint(sq[2].Sorted())
	if tr1 > tr2 {
		tr1, tr2 = tr2, tr1
	}
	index[struct {
		t1 string
		t2 string
	}{tr1, tr2}] = struct{}{}

	tr1 = fmt.Sprint(sq.Column(0).Sorted())
	tr2 = fmt.Sprint(sq.Column(1).Sorted())
	if tr1 > tr2 {
		tr1, tr2 = tr2, tr1
	}
	index[struct {
		t1 string
		t2 string
	}{tr1, tr2}] = struct{}{}

	tr1 = fmt.Sprint(sq.Column(0).Sorted())
	tr2 = fmt.Sprint(sq.Column(2).Sorted())
	if tr1 > tr2 {
		tr1, tr2 = tr2, tr1
	}
	index[struct {
		t1 string
		t2 string
	}{tr1, tr2}] = struct{}{}

	tr1 = fmt.Sprint(sq.Column(1).Sorted())
	tr2 = fmt.Sprint(sq.Column(2).Sorted())
	if tr1 > tr2 {
		tr1, tr2 = tr2, tr1
	}
	index[struct {
		t1 string
		t2 string
	}{tr1, tr2}] = struct{}{}

}

func fitToIndex(sq0, sq1 square.Matrix, index map[struct{ t1, t2 string }]struct{}) bool {

	//2 squares only can fit into a cube, if each of 3 rows of 1st square can be paired with 3 rows or 3 columns in 2nd square via index
	var matches int = 0
	for i := 0; i <= 2; i++ {
		var hasMatch bool = false
		for j := 0; j <= 2; j++ {
			tr1, tr2 := fmt.Sprint(sq0[i].Sorted()), fmt.Sprint(sq1[j].Sorted())
			if tr1 > tr2 {
				tr1, tr2 = tr2, tr1
			}
			if _, ok := index[struct {
				t1 string
				t2 string
			}{tr1, tr2}]; ok {
				hasMatch = true
			}
		}
		if hasMatch {
			matches++
		}
	}
	return matches == 3
}

func translateRowsToMatch(in, match square.Matrix, index map[struct{ t1, t2 string }]struct{}) square.Matrix {
	var out square.Matrix
	for i := 0; i <= 2; i++ {
		for j := 0; j <= 2; j++ {
			tr1, tr2 := fmt.Sprint(in[i].Sorted()), fmt.Sprint(match[j].Sorted())
			if tr1 > tr2 {
				tr1, tr2 = tr2, tr1
			}
			if _, ok := index[struct {
				t1 string
				t2 string
			}{tr1, tr2}]; ok { //TODO do not take duplicates
				duplicate := false
				for k := 0; k <= 2; k++ {
					if out[k] == in[j] {
						duplicate = true
					}
				}
				if !duplicate {
					out[i] = in[j]
					break
				}
			}
		}
		if (out[i] == triplet.Triplet{}) {
			return square.Matrix{}
		}
	}

	return out
}
