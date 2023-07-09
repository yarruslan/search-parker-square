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
	self      *square.Matrix
	connected []*square.Matrix
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
	//TODO implement
	/*
		loop at nodes
			loop at added nodes to check it is not duplicate //duplicate = *square.Same(square2)
			append to result
			loop at added nodes and add bi-diretional connection on intersecting squares // connection = *square.Intersect(square2)
				'same' and 'intersect' do "rotate"(?) and "revolve"*2 and "mirror" both squares, to check all variations
		}
	*/
	return &Graph{}
}

func (g *Graph) containsCube() bool {
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
