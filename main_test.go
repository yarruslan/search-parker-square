package main

import (
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
	setx144 := make(map[sumSquares][]triplet)
	setx144, _, _ = generate(setx144, []sumSquares{}, 21609*144, 21609*144)
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
				setx144[21609*144],
			},
			matrix{
				triplet{97 * 97 * 144, 82 * 82 * 144, 74 * 74 * 144}, triplet{94 * 94 * 144, 113 * 113 * 144, 2 * 2 * 144}, triplet{58 * 58 * 144, 46 * 46 * 144, 127 * 127 * 144},
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
