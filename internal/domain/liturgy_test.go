package domain

import "testing"

  func TestVerseContains(t *testing.T) {
  	verse := Verse{
  		StartMs: 1000,
  		EndMs:   2000,
  	}

  	tests := []struct {
  		name     string
  		position int
  		expected bool
  	}{
  		{
  			name:     "before start",
  			position: 999,
  			expected: false,
  		},
  		{
  			name:     "at start",
  			position: 1000,
  			expected: true,
  		},
  		{
  			name:     "between start and end",
  			position: 1500,
  			expected: true,
  		},
  		{
  			name:     "at end",
  			position: 2000,
  			expected: false,
  		},
  		{
  			name:     "after end",
  			position: 2001,
  			expected: false,
  		},
  	}

  	for _, tt := range tests {
  		t.Run(tt.name, func(t *testing.T) {
  			got := verse.Contains(tt.position)

  			if got != tt.expected {
  				t.Errorf(
  					"Contains(%d) = %v; expected %v",
  					tt.position,
  					got,
  					tt.expected,
  				)
  			}
  		})
  	}
  }
