package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/yarruslan/search-parker-square/internal/square"
	triplet "github.com/yarruslan/search-parker-square/internal/triplet"
)

func Test_lookupSubset(t *testing.T) { //[97, 82, 74](21609) [94, 113, 2](21609) [58, 46, 127](21609) 1
	type args struct {
		set []triplet.Triplet
	}

	set1 := new(triplet.Generator).Init(21609, 21609, 0, 11).GetSet()
	setx16 := new(triplet.Generator).Init(21609*16, 21609*16, 0, 11).GetSet()
	setx225 := new(triplet.Generator).Init(21609*225, 21609*225, 0, 11).GetSet()
	tests := []struct {
		name string
		args args
		want square.Matrix
	}{
		{"base",
			args{
				set1[21609],
			},
			square.Matrix{
				triplet.Triplet{97 * 97, 82 * 82, 74 * 74}, triplet.Triplet{94 * 94, 113 * 113, 2 * 2}, triplet.Triplet{58 * 58, 46 * 46, 127 * 127},
			},
		},
		{"base x4",
			args{
				setx16[21609*16],
			},
			square.Matrix{
				triplet.Triplet{97 * 97 * 16, 82 * 82 * 16, 74 * 74 * 16}, triplet.Triplet{94 * 94 * 16, 113 * 113 * 16, 2 * 2 * 16}, triplet.Triplet{58 * 58 * 16, 46 * 46 * 16, 127 * 127 * 16},
			},
		},
		{"base x12",
			args{
				setx225[21609*225],
			},
			square.Matrix{
				triplet.Triplet{97 * 97 * 225, 82 * 82 * 225, 74 * 74 * 225}, triplet.Triplet{94 * 94 * 225, 113 * 113 * 225, 2 * 2 * 225}, triplet.Triplet{58 * 58 * 225, 46 * 46 * 225, 127 * 127 * 225},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got square.Matrix
			for _, sq := range square.CombineTripletsToMatrixes(tt.args.set, triplet.SearchSemiMagic) {
				if sq.CountDiagonals() > 0 { //TODO refactor, no single resp atm
					got = sq
					break
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lookupSubset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findSquaresWithDiagonals(t *testing.T) {
	type args struct {
		start triplet.Square
		end   triplet.Square
		d     int
		res   chan []fmt.Stringer
	}
	tests := []struct {
		name string
		args args
		want [][]square.Matrix
	}{
		{"baser",
			args{
				21609,
				21609,
				1,
				make(chan []fmt.Stringer),
			},
			[][]square.Matrix{
				{
					{
						triplet.Triplet{97 * 97, 82 * 82, 74 * 74},
						triplet.Triplet{94 * 94, 113 * 113, 2 * 2},
						triplet.Triplet{58 * 58, 46 * 46, 127 * 127},
					},
					{
						triplet.Triplet{97 * 97, 94 * 94, 58 * 58},
						triplet.Triplet{82 * 82, 113 * 113, 46 * 46},
						triplet.Triplet{74 * 74, 2 * 2, 127 * 127},
					},
				},
			},
		},

		{"base",
			args{
				20000,
				88000,
				1,
				make(chan []fmt.Stringer),
			},
			[][]square.Matrix{
				{
					{
						triplet.Triplet{97 * 97, 82 * 82, 74 * 74},
						triplet.Triplet{94 * 94, 113 * 113, 2 * 2},
						triplet.Triplet{58 * 58, 46 * 46, 127 * 127},
					},
					{
						triplet.Triplet{97 * 97, 94 * 94, 58 * 58},
						triplet.Triplet{82 * 82, 113 * 113, 46 * 46},
						triplet.Triplet{74 * 74, 2 * 2, 127 * 127},
					},
				},
				{
					{
						triplet.Triplet{194 * 194, 164 * 164, 148 * 148},
						triplet.Triplet{188 * 188, 226 * 226, 4 * 4},
						triplet.Triplet{116 * 116, 92 * 92, 254 * 254},
					},

					{
						triplet.Triplet{194 * 194, 188 * 188, 116 * 116},
						triplet.Triplet{164 * 164, 226 * 226, 92 * 92},
						triplet.Triplet{148 * 148, 4 * 4, 254 * 254},
					},
				},
			},
		},
	}
	progressStep := triplet.Square(100000)
	threads := 11
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got [][]fmt.Stringer
			done := make(chan struct{})
			go func() {
				for pack := range tt.args.res {
					got = append(got, pack)
				}
				done <- struct{}{}
			}()
			g := new(square.Generator).Init(new(triplet.Generator).Init(tt.args.start, tt.args.end, progressStep, threads), threads)
			g.GenerateSquares(tt.args.d, tt.args.res) //async race possible? more than 1 responce?
			//for got := range tt.args.res {
			<-done
			if len(tt.want) != len(got) {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
			for i, res := range got {
				if len(tt.want[i]) != len(got[i]) {
					t.Errorf("got = %v, want %v", got, tt.want)
				}
				for j, sq := range res {
					if len(tt.want) < i+1 || len(tt.want[i]) < j+1 || !reflect.DeepEqual(sq, tt.want[i][j]) {
						t.Errorf("got = %v, want %v", got, tt.want)
					}
				}
			}
			//}
		})
	}
}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "dummy",
			args: []string{"", "-end", "100000"},
			want: `Square  [97 82 74][94 113 2][58 46 127](21609)  has 1 diagonals
Square  [97 94 58][82 113 46][74 2 127](21609)  has 1 diagonals
Square  [194 164 148][188 226 4][116 92 254](86436)  has 1 diagonals
Square  [194 188 116][164 226 92][148 4 254](86436)  has 1 diagonals
`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args

			old := os.Stdout
			r, w, err := os.Pipe()
			if err != nil {
				log.Fatal("File error")
			}
			os.Stdout = w

			main()

			outC := make(chan string)
			// copy the output in a separate goroutine so printing can't block indefinitely
			go func() {
				var buf bytes.Buffer
				io.Copy(&buf, r)
				outC <- buf.String()
			}()
			// back to normal state
			w.Close()
			os.Stdout = old // restoring the real stdout
			got := <-outC

			//os.Stdout = stdout
			if got != tt.want {
				t.Errorf("got \n%v, want \n%v", got, tt.want)
			}

		})
	}
}
