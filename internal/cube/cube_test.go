package cube

import (
	"fmt"
	"testing"

	"github.com/yarruslan/search-parker-square/internal/square"
	"github.com/yarruslan/search-parker-square/internal/triplet"
)

func TestCube_String(t *testing.T) {
	tests := []struct {
		name string
		c    *Cube
		want string
	}{
		{
			"basic",
			&Cube{},
			"[0 0 0][0 0 0][0 0 0](0)[0 0 0][0 0 0][0 0 0](0)[0 0 0][0 0 0][0 0 0](0)",
		},
		{
			"sequential",
			&Cube{{{0, 1, 4}, {9, 16, 25}, {36, 49, 64}}, {{81, 100, 121}, {144, 169, 196}, {225, 256, 17 * 17}}, {{18 * 18, 19 * 19, 20 * 20}, {21 * 21, 22 * 22, 23 * 23}, {24 * 24, 25 * 25, 26 * 26}}},
			"[0 1 2][3 4 5][6 7 8](5)[9 10 11][12 13 14][15 16 17](302)[18 19 20][21 22 23][24 25 26](1085)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("Cube.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

//TODO improve test coverage

func TestGenerator_GenerateCubes(t *testing.T) {
	type args struct {
		searchType int
		result     chan []fmt.Stringer
	}
	tests := []struct {
		name string
		g    *Generator
		args args
	}{
		/*
			{
				name: "test sum 1863225",
				g:    new(Generator).Init(new(square.Generator).Init(new(triplet.Generator).Init(1863225, 1863225, 1, 1), 1)),
				args: args{
					searchType: triplet.SearchCube,
					result:     make(chan []fmt.Stringer), //nothing
				},
			},
		*/
		/*
			{
				name: "test sum 4060225",
				g:    new(Generator).Init(new(square.Generator).Init(new(triplet.Generator).Init(4060225, 4060225, 1, 1), 1)),
				args: args{
					searchType: triplet.SearchCube,
					result:     make(chan []fmt.Stringer), // 2 connected
				},
			},

			{
				name: "test sum 1729225",
				g:    new(Generator).Init(new(square.Generator).Init(new(triplet.Generator).Init(1729225, 1729225, 1, 1), 1)),
				args: args{
					searchType: triplet.SearchCube,
					result:     make(chan []fmt.Stringer), //nothing
				},
			},
		*/
		/*
			{
				name: "test sum 30525625",
				g:    new(Generator).Init(new(square.Generator).Init(new(triplet.Generator).Init(30525625, 30525625, 1, 1), 1)),
				args: args{
					searchType: triplet.SearchCube,
					result:     make(chan []fmt.Stringer), //nothing
				},
			},
		*/
		/*
			{
				name: "test sum 46580625",
				g:    new(Generator).Init(new(square.Generator).Init(new(triplet.SquareGenerator).Init(46580625, 46580625, 10), 1)),
				args: args{
					searchType: triplet.SearchCube,
					result:     make(chan []fmt.Stringer), //nothing
				},
			},
		*/
		/*
			{
				name: "test sum 59830225",
				g:    new(Generator).Init(new(square.Generator).Init(new(triplet.SquareGenerator).Init(59830225, 59830225, 10), 1)),
				args: args{
					searchType: triplet.SearchCube,
					result:     make(chan []fmt.Stringer), //nothing
				},
			},
			{
				name: "test sum 79655625",
				g:    new(Generator).Init(new(square.Generator).Init(new(triplet.SquareGenerator).Init(79655625, 79655625, 10), 1)),
				args: args{
					searchType: triplet.SearchCube,
					result:     make(chan []fmt.Stringer), //nothing
				},
			},
		*/
		/*
			{
				name: "test sum 88830625",
				g:    new(Generator).Init(new(square.Generator).Init(new(triplet.SquareGenerator).Init(88830625, 88830625, 10), 1)),
				args: args{
					searchType: triplet.SearchCube,
					result:     make(chan []fmt.Stringer), //nothing
				},
			},
		*/

		{
			name: "test sum 65155115025",
			g:    new(Generator).Init(new(square.Generator).Init(new(triplet.SquareGenerator).Init(65155115025, 65155115025, 10), 1)),
			args: args{
				searchType: triplet.SearchCube,
				result:     make(chan []fmt.Stringer), //nothing
			},
		},

		/*
			{
				name: "test sum 117831025",
				g:    new(Generator).Init(new(square.Generator).Init(new(triplet.SquareGenerator).Init(117831025, 117831025, 10), 1)),
				args: args{
					searchType: triplet.SearchCube,
					result:     make(chan []fmt.Stringer), //nothing
				},
			},
		*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go func() {
				for resStr := range tt.args.result {
					fmt.Println(resStr) //TODO make test for string got == want
				}
			}()

			tt.g.GenerateCubes(tt.args.searchType, tt.args.result)

		})
	}
}

func TestGenerator2_GenerateCubes(t *testing.T) {
	type args struct {
		searchType int
		result     chan []fmt.Stringer
	}
	tests := []struct {
		name string
		g    *Generator2
		args args
	}{

		{
			name: "test sum 65155115025",
			g:    new(Generator2).Init(new(square.Generator).Init(new(triplet.SquareGenerator).Init(65155115025, 65155115025, 10), 1)),
			args: args{
				searchType: triplet.SearchCube,
				result:     make(chan []fmt.Stringer), //nothing
			},
		},
		/*
			{
				name: "test sum 117831025",
				g:    new(Generator2).Init(new(square.Generator).Init(new(triplet.SquareGenerator).Init(117831025, 117831025, 10), 1)),
				args: args{
					searchType: triplet.SearchCube,
					result:     make(chan []fmt.Stringer), //nothing
				},
			},
		*/
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		go func() {
			for resStr := range tt.args.result {
				fmt.Println(resStr) //TODO make test for string got == want
			}
		}()

		t.Run(tt.name, func(t *testing.T) {
			tt.g.GenerateCubes(tt.args.searchType, tt.args.result)
		})
	}
}
