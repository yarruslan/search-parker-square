package main

import (
	"reflect"
	"testing"
)

func Test_lookupSubset(t *testing.T) { //[97, 82, 74](21609) [94, 113, 2](21609) [58, 46, 127](21609) 1
	type args struct {
		set []triplet
	}
	set200 := make(map[int][]triplet)
	set200, _ = generate(set200, []int{}, 2, 200)
	set510 := make(map[int][]triplet)
	set510, _ = generate(set510, []int{}, 97*4, 510)
	tests := []struct {
		name string
		args args
		want matrix
	}{
		{"base",
			args{
				set200[21609],
			},
			matrix{
				triplet{97 * 97, 82 * 82, 74 * 74}, triplet{94 * 94, 113 * 113, 2 * 2}, triplet{58 * 58, 46 * 46, 127 * 127},
			},
		},
		{"base x4",
			args{
				set510[21609*16],
			},
			matrix{
				triplet{97 * 97 * 16, 82 * 82 * 16, 74 * 74 * 16}, triplet{94 * 94 * 16, 113 * 113 * 16, 2 * 2 * 16}, triplet{58 * 58 * 16, 46 * 46 * 16, 127 * 127 * 16},
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lookupSubset(tt.args.set); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lookupSubset() = %v, want %v", got, tt.want)
			}
		})
	}
}
