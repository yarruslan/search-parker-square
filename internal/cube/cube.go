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
	square      *square.Matrix
	connections map[*GraphNode]struct{}
}
type Graph []GraphNode

func (c *Cube) String() string {
	return c[0].String() + c[1].String() + c[2].String()
}

func (g *Generator) Init(sq *square.Generator) *Generator {
	g.sq = sq
	return g
}

func (g *Generator) CombineSquaresToCubes(s []square.Matrix) []Cube {
	var result []Cube

	if graph := g.buildGraphOfSquares(s); graph.containsCube() {
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
		newNode.square = &sq
		newNode.connections = make(map[*GraphNode]struct{})
		for _, node := range graph {
			if node.square.Same(&sq) {
				exist = true
				break
			}
		}
		if !exist {
			graph = append(graph, newNode)

			for _, node := range graph {
				if node.square.Intersect(&sq) {
					node.connect(&newNode)
				}
			}
		}
	}
	return &graph
}

func (g *Graph) containsCube() bool {

	strongNodes := 0
	potentialNodes := 0
	for _, node := range *g {
		if len(node.connections) >= 3 {
			strongNodes++
		}
		if len(node.connections) >= 2 {
			potentialNodes++
		}
	}
	if potentialNodes > 1 {
		fmt.Println(len(*g), potentialNodes, strongNodes, (*g)[0].square.String()) //Further implementation make sense if strongly interconnected squares exist
	}
	//TODO implement
	/*
		candidate to cubes has connected graph of 9 connected squares, each of them has connection to 6 others
		so
		1.find loosly connected, remove them.
		2.rinse, repeat.
		3.thoroughly check the 9-s+ connection graph
		4.validate
	*/
	return false
}

func (g *Graph) getCubes() []Cube {
	//TODO implement
	return []Cube{}
}

func (a *GraphNode) connect(b *GraphNode) {
	a.connections[b] = struct{}{}
	b.connections[a] = struct{}{}
}

/*
func (a *GraphNode) disconnect(b *GraphNode) {
	delete(a.connections, b)
	delete(b.connections, a)
}
*/
