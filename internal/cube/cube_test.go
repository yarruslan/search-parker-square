package cube

import "testing"

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
