package triplet

import (
	"reflect"
	"testing"
)

func TestGenerate(t *testing.T) {
	type args struct {
		groups     IndexedTriplets
		index      []Square
		windowLow  Square
		windowHigh Square
	}
	tests := []struct {
		name      string
		args      args
		wantMap   IndexedTriplets
		wantSlice []Square
		wantEnd   Square
	}{
		{
			name: "TestEmpty",
			args: args{
				groups:     IndexedTriplets{},
				index:      []Square{},
				windowLow:  0,
				windowHigh: 1,
			},
			wantMap:   IndexedTriplets{},
			wantSlice: []Square{},
			wantEnd:   0,
		},
		{
			name: "Test Min",
			args: args{
				groups:     IndexedTriplets{},
				index:      []Square{},
				windowLow:  5,
				windowHigh: 5,
			},
			wantMap: IndexedTriplets{
				5: {
					{4, 1, 0},
				},
			},
			wantSlice: []Square{
				Square(5),
			},
			wantEnd: 5,
		},
		{
			name: "Test Min+1",
			args: args{
				groups:     IndexedTriplets{},
				index:      []Square{},
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
			wantSlice: []Square{
				10, 13,
			},
			wantEnd: 13,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := new(Generator).Init(tt.args.windowLow, tt.args.windowHigh, tt.args.windowHigh-tt.args.windowLow, 11)
			//gotMap, gotSlice := Generate(tt.args.groups, tt.args.index, tt.args.windowLow, tt.args.windowHigh)
			gotMap, gotSlice := g.set, g.index
			if !reflect.DeepEqual(gotMap, tt.wantMap) {
				t.Errorf("Generate() got = %v, want %v", gotMap, tt.wantMap)
			}
			if !reflect.DeepEqual(gotSlice, tt.wantSlice) {
				t.Errorf("Generate() got1 = %v, want %v", gotSlice, tt.wantSlice)
			}
		})
	}
}

func TestGenerateCheckCount(t *testing.T) {
	type args struct {
		groups     IndexedTriplets
		index      []Square
		windowLow  Square
		windowHigh Square
	}
	tests := []struct {
		name     string
		args     args
		countMap map[Square]int
	}{
		{
			name: "Test minimal",
			args: args{
				groups:     IndexedTriplets{},
				index:      []Square{},
				windowLow:  0,
				windowHigh: 5,
			},
			countMap: map[Square]int{
				5: 1,
			},
		},
		{
			name: "Test more",
			args: args{
				groups:     IndexedTriplets{},
				index:      []Square{},
				windowLow:  0,
				windowHigh: 13,
			},
			countMap: map[Square]int{
				5:  1,
				10: 1,
				13: 1,
			},
		},

		{
			name: "Test at 1st square",
			args: args{
				groups:     IndexedTriplets{},
				index:      []Square{},
				windowLow:  21609,
				windowHigh: 21609,
			},
			countMap: map[Square]int{
				21609: 40,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := new(Generator).Init(tt.args.windowLow, tt.args.windowHigh, tt.args.windowHigh-tt.args.windowLow, 11)
			//gotMap, _ := Generate(tt.args.groups, tt.args.index, tt.args.windowLow, tt.args.windowHigh)
			gotMap := g.set
			for k, v := range gotMap {
				if len(v) != tt.countMap[k] {
					t.Errorf("Generate() got = %v: %v, want %v: %v", k, len(v), k, tt.countMap[k])
				}
			}
			for k, v := range tt.countMap {
				if len(gotMap[k]) != v {
					t.Errorf("Generate() got = %v: %v, want %v: %v", k, len(gotMap[k]), k, v)
				}
			}
		})
	}
}

func TestTriplet_String(t *testing.T) {
	tests := []struct {
		name string
		tr   *Triplet
		want string
	}{
		{
			name: "Empty",
			tr:   &Triplet{},
			want: "[0 0 0]",
		},
		{
			name: "Not empty",
			tr:   &Triplet{1, 4, 9},
			want: "[1 2 3]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.String(); got != tt.want {
				t.Errorf("Triplet.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
