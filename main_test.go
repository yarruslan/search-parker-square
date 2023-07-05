package main

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_lookupSubset(t *testing.T) { //[97, 82, 74](21609) [94, 113, 2](21609) [58, 46, 127](21609) 1
	type args struct {
		set []triplet
	}
	set1 := make(map[sumSquares][]triplet)
	set1, _, _ = generate(set1, []sumSquares{}, 21609, 21609)
	setx16 := make(map[sumSquares][]triplet)
	setx16, _, _ = generate(setx16, []sumSquares{}, 21609*16, 21609*16)
	setx225 := make(map[sumSquares][]triplet)
	setx225, _, _ = generate(setx225, []sumSquares{}, 21609*225, 21609*225)
	tests := []struct {
		name string
		args args
		want matrix
	}{
		{"base",
			args{
				set1[21609],
			},
			matrix{
				triplet{97 * 97, 82 * 82, 74 * 74}, triplet{94 * 94, 113 * 113, 2 * 2}, triplet{58 * 58, 46 * 46, 127 * 127},
			},
		},
		{"base x4",
			args{
				setx16[21609*16],
			},
			matrix{
				triplet{97 * 97 * 16, 82 * 82 * 16, 74 * 74 * 16}, triplet{94 * 94 * 16, 113 * 113 * 16, 2 * 2 * 16}, triplet{58 * 58 * 16, 46 * 46 * 16, 127 * 127 * 16},
			},
		},
		{"base x12",
			args{
				setx225[21609*225],
			},
			matrix{
				triplet{97 * 97 * 225, 82 * 82 * 225, 74 * 74 * 225}, triplet{94 * 94 * 225, 113 * 113 * 225, 2 * 2 * 225}, triplet{58 * 58 * 225, 46 * 46 * 225, 127 * 127 * 225},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got matrix
			for _, sq := range lookupSubset(tt.args.set) {
				if countDiagonals(sq) > 0 { //TODO refactor, no single resp atm
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
		start sumSquares
		end   sumSquares
		d     int
		res   chan []fmt.Stringer
	}
	resChan := make(chan []fmt.Stringer)
	tests := []struct {
		name string
		args args
		want [][]matrix
	}{
		{"base",
			args{
				1,
				100000,
				1,
				resChan,
			},
			[][]matrix{
				{
					{
						triplet{97 * 97, 82 * 82, 74 * 74},
						triplet{94 * 94, 113 * 113, 2 * 2},
						triplet{58 * 58, 46 * 46, 127 * 127},
					},
					{
						triplet{97 * 97, 94 * 94, 58 * 58},
						triplet{82 * 82, 113 * 113, 46 * 46},
						triplet{74 * 74, 2 * 2, 127 * 127},
					},
				},
				{
					{
						triplet{194 * 194, 164 * 164, 148 * 148},
						triplet{188 * 188, 226 * 226, 4 * 4},
						triplet{116 * 116, 92 * 92, 254 * 254},
					},

					{
						triplet{194 * 194, 188 * 188, 116 * 116},
						triplet{164 * 164, 226 * 226, 92 * 92},
						triplet{148 * 148, 4 * 4, 254 * 254},
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
			for i, res := range got {
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
