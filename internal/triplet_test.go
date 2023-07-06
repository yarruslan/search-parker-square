package triplet

import (
	"reflect"
	"testing"
)

func TestGenerate(t *testing.T) {
	type args struct {
		groups     IndexedTriplets
		index      []SumSquares
		windowLow  SumSquares
		windowHigh SumSquares
	}
	tests := []struct {
		name      string
		args      args
		wantMap   IndexedTriplets
		wantSlice []SumSquares
		wantEnd   SumSquares
	}{
		{
			name: "TestEmpty",
			args: args{
				groups:     IndexedTriplets{},
				index:      []SumSquares{},
				windowLow:  0,
				windowHigh: 1,
			},
			wantMap:   IndexedTriplets{},
			wantSlice: []SumSquares{},
			wantEnd:   0,
		},
		{
			name: "Test Min",
			args: args{
				groups:     IndexedTriplets{},
				index:      []SumSquares{},
				windowLow:  5,
				windowHigh: 5,
			},
			wantMap: IndexedTriplets{
				5: {
					{4, 1, 0},
				},
			},
			wantSlice: []SumSquares{
				SumSquares(5),
			},
			wantEnd: 5,
		},
		{
			name: "Test Min+1",
			args: args{
				groups:     IndexedTriplets{},
				index:      []SumSquares{},
				windowLow:  6,
				windowHigh: 13,
			},
			wantMap: IndexedTriplets{
				10: {
					{9, 1, 0},
				},
				13: {
					{9, 4, 0},
				},
			},
			wantSlice: []SumSquares{
				10, 13,
			},
			wantEnd: 13,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMap, gotSlice, gotEnd := Generate(tt.args.groups, tt.args.index, tt.args.windowLow, tt.args.windowHigh)
			if !reflect.DeepEqual(gotMap, tt.wantMap) {
				t.Errorf("Generate() got = %v, want %v", gotMap, tt.wantMap)
			}
			if !reflect.DeepEqual(gotSlice, tt.wantSlice) {
				t.Errorf("Generate() got1 = %v, want %v", gotSlice, tt.wantSlice)
			}
			if !reflect.DeepEqual(gotEnd, tt.wantEnd) {
				t.Errorf("Generate() got2 = %v, want %v", gotEnd, tt.wantEnd)
			}
		})
	}
}
