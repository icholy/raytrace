package raytrace

import "testing"

func TestVec3(t *testing.T) {
	tests := []struct {
		want, got Vec3
	}{
		{
			want: Vec3{1, 1, 1},
			got:  Unit(),
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if tt.want != tt.got {
				t.Fatalf("want %s, got %s", tt.want, tt.got)
			}
		})
	}
}
