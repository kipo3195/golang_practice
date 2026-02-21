package unitTest

import "testing"

func TestAdd(t *testing.T) {

	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		// Table Driven Test
		{"양수 더하기", 1, 2, 3},
		{"음수 더하기", -1, -1, -2},
		{"0 더하기", 5, 0, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Add(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}

}
