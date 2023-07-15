package cube

import (
	"fmt"
	"log"

	"github.com/yarruslan/search-parker-square/internal/square"
)

type Cube [3]square.Matrix

type Generator struct {
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

func (g *Generator) CombineSquaresToCubes(s []square.Matrix) []Cube {
	var result []Cube

	if graph := g.buildGraphOfSquares(s).filter(); graph.canContainCube() {
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

// Graph is based on squares. If 2 squares are intersecting they are connected
func (g *Generator) buildGraphOfSquares(in []square.Matrix) *Graph {
	var graph Graph
	for _, sq := range in {
		exist := false
		var newNode GraphNode
		newNode.square = sq
		newNode.connections = make(map[*GraphNode]struct{})
		for _, node := range graph {
			if node.square.Same(&newNode.square) { //TODO n^2 is too long. Use index.
				exist = true
				break
			}
		}
		if !exist {
			for _, node := range graph {
				if node.square.Intersect(&sq) {
					node.connect(&newNode)
				}
			}
			graph = append(graph, &newNode)
		}
	}
	return &graph
}

func (g *Graph) canContainCube() bool { //TODO too many responsibilities in method

	connectedNodes := 0
	for _, node := range *g { //TODO this is duplicated in filter
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
		2. if 6 of its connected squares have another common intersection square - there is a cube
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
	//TODO should be exactly 1 cube
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

	for plane0, _ := range verticalPlanes {
		if plane0.square.Contains(base[0]) && plane0.square.Intersect(&secondLayer) {
			firstPlane = plane0.square
			log.Println("Getting closer", firstPlane)
			fmt.Println("Getting closer", firstPlane)
			for plane1, _ := range verticalPlanes {
				if plane1.square.Contains(base[1]) && plane1.square.Intersect(&secondLayer) { // && intersection uses different than first plane.
					secondPlane = plane1.square
					log.Println("Getting hot", firstPlane, secondPlane)
					fmt.Println("Getting hot", firstPlane, secondPlane)
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
