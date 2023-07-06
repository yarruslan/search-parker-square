package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/yarruslan/search-parker-square/internal/matrix"
	triplet "github.com/yarruslan/search-parker-square/internal/triplet"
)

func Test_lookupSubset(t *testing.T) { //[97, 82, 74](21609) [94, 113, 2](21609) [58, 46, 127](21609) 1
	type args struct {
		set []triplet.Triplet
	}
	set1 := make(triplet.IndexedTriplets)
	set1, _, _ = triplet.Generate(set1, []triplet.SumSquares{}, 21609, 21609)
	setx16 := make(triplet.IndexedTriplets)
	setx16, _, _ = triplet.Generate(setx16, []triplet.SumSquares{}, 21609*16, 21609*16)
	setx225 := make(triplet.IndexedTriplets)
	setx225, _, _ = triplet.Generate(setx225, []triplet.SumSquares{}, 21609*225, 21609*225)
	tests := []struct {
		name string
		args args
		want matrix.Matrix
	}{
		{"base",
			args{
				set1[21609],
			},
			matrix.Matrix{
				triplet.Triplet{97 * 97, 82 * 82, 74 * 74}, triplet.Triplet{94 * 94, 113 * 113, 2 * 2}, triplet.Triplet{58 * 58, 46 * 46, 127 * 127},
			},
		},
		{"base x4",
			args{
				setx16[21609*16],
			},
			matrix.Matrix{
				triplet.Triplet{97 * 97 * 16, 82 * 82 * 16, 74 * 74 * 16}, triplet.Triplet{94 * 94 * 16, 113 * 113 * 16, 2 * 2 * 16}, triplet.Triplet{58 * 58 * 16, 46 * 46 * 16, 127 * 127 * 16},
			},
		},
		{"base x12",
			args{
				setx225[21609*225],
			},
			matrix.Matrix{
				triplet.Triplet{97 * 97 * 225, 82 * 82 * 225, 74 * 74 * 225}, triplet.Triplet{94 * 94 * 225, 113 * 113 * 225, 2 * 2 * 225}, triplet.Triplet{58 * 58 * 225, 46 * 46 * 225, 127 * 127 * 225},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got matrix.Matrix
			for _, sq := range matrix.LookupSubset(tt.args.set) {
				if matrix.CountDiagonals(sq) > 0 { //TODO refactor, no single resp atm
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
		start triplet.SumSquares
		end   triplet.SumSquares
		d     int
		res   chan []fmt.Stringer
	}
	resChan := make(chan []fmt.Stringer)
	tests := []struct {
		name string
		args args
		want [][]matrix.Matrix
	}{
		{"base",
			args{
				20000,
				88000,
				1,
				resChan,
			},
			[][]matrix.Matrix{
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
			findSquaresWithDiagonals(tt.args.start, tt.args.end, tt.args.d, tt.args.res) //async race possible? more than 1 responce?
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
