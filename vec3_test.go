package raytrace

import "testing"

func TestVec3(t *testing.T) {
	tests := []struct {
		name      string
		want, got Vec3
	}{
		{
			name: "neg",
			want: Vec3{-1, -2, -3},
			got:  Vec3{1, 2, 3}.Neg(),
		},
		{
			name: "add",
			want: Vec3{3, 5, 7},
			got:  Vec3{1, 2, 3}.Add(Vec3{2, 3, 4}),
		},
		{
			name: "sub",
			want: Vec3{1, 2, 3},
			got:  Vec3{3, 3, 3}.Sub(Vec3{2, 1, 0}),
		},
		{
			name: "mul",
			want: Vec3{2, 4, 6},
			got:  Vec3{1, 2, 3}.Mul(Vec3{2, 2, 2}),
		},
		{
			name: "div",
			want: Vec3{6, 6, 6},
			got:  Vec3{48, 36, 12}.Div(Vec3{8, 6, 2}),
		},
		{
			name: "scalar_mul",
			want: Vec3{2, 2, 2},
			got:  Vec3{1, 1, 1}.ScalarMul(2),
		},
		{
			name: "scalar_div",
			want: Vec3{2, 2, 2},
			got:  Vec3{4, 4, 4}.ScalarDiv(2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want != tt.got {
				t.Fatalf("want %s, got %s", tt.want, tt.got)
			}
		})
	}
}
