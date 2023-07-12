package cube

import (
	"fmt"

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
			if node.square.Same(&newNode.square) {
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
	unfilteredNodes := len(*g)
	for _, node := range *g {
		if len(node.connections) >= 6 { //each square shold connect to 6 other to form a cube
			connectedNodes++
		}
	}
	if connectedNodes >= 1 {
		gf := g.filter()
		g = gf //TODO, bug. it does not change input
		filteredNodes := len(*gf)
		//if filteredNodes >= 9 {
		fmt.Println(unfilteredNodes, connectedNodes, filteredNodes, (*g)[0].square.String()) //Further implementation make sense if strongly interconnected squares exist
		//}
	}
	//TODO implement
	/*
		candidate to cubes has connected graph of 9 connected squares, each of them has connection to 6 others
		so
		1.find loosly connected, remove them.
		2.rinse, repeat.
		//3.thoroughly check the 9-s+ connection graph //not here
		//4.validate //not here
	*/
	if connectedNodes >= 9 { //TODO check based on filtered result
		return true
	}
	return false
}

func (g *Graph) getCubes() []Cube {
	//TODO implement
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
		for k, v := range connectionsStats {
			if k != base && v >= 6 {
				//TODO thats a bit more compliacted
				/*
					sum := base.square[0][0] + base.square[0][1] + base.square[0][2]
					top := square.Matrix{
						{sum - base.square[0][0] - k.square[0][0], sum - base.square[0][1] - k.square[0][1], sum - base.square[0][2] - k.square[0][2]},
						{sum - base.square[1][0] - k.square[1][0], sum - base.square[1][1] - k.square[1][1], sum - base.square[1][2] - k.square[1][2]},
						{sum - base.square[2][0] - k.square[2][0], sum - base.square[2][1] - k.square[2][1], sum - base.square[2][2] - k.square[2][2]},
					}
				*/
				fmt.Println("Found cube:", base.square, k.square, "[? ? ?]")
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
