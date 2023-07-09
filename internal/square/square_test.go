package square

import (
	"testing"

	"github.com/yarruslan/search-parker-square/internal/triplet"
)

func TestMatrix_String(t *testing.T) {
	tests := []struct {
		name string
		m    Matrix
		want string
	}{
		{
			name: "Null matrix text",
			m:    Matrix{},
			want: "[0 0 0][0 0 0][0 0 0](0)",
		},
		{
			name: "Squares matrix text",
			m: Matrix{
				triplet.Triplet{0, 1, 4},
				triplet.Triplet{9, 16, 25},
				triplet.Triplet{36, 49, 64},
			},
			want: "[0 1 2][3 4 5][6 7 8](5)",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.want {
				t.Errorf("Matrix.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
