package types

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimeParam_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
		validate  func(t *testing.T, tp TimeParam)
	}{
		{
			name:      "RFC3339 format",
			input:     `"2023-12-25T10:30:00Z"`,
			wantError: false,
			validate: func(t *testing.T, tp TimeParam) {
				expected := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)
				if !tp.Time.Equal(expected) {
					t.Errorf("Expected %v, got %v", expected, tp.Time)
				}
			},
		},
		{
			name:      "RFC3339 format without quotes",
			input:     `2023-12-25T10:30:00Z`,
			wantError: false,
			validate: func(t *testing.T, tp TimeParam) {
				expected := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)
				if !tp.Time.Equal(expected) {
					t.Errorf("Expected %v, got %v", expected, tp.Time)
				}
			},
		},
		{
			name:      "Week duration - 1 week",
			input:     `1w`,
			wantError: false,
			validate: func(t *testing.T, tp TimeParam) {
				now := time.Now()
				expected := now.Add(-7 * 24 * time.Hour)
				// Allow for small time differences due to execution time
				if tp.Time.Sub(expected).Abs() > time.Second {
					t.Errorf("Expected approximately %v, got %v", expected, tp.Time)
				}
			},
		},
		{
			name:      "Week duration - 2 weeks",
			input:     `2w`,
			wantError: false,
			validate: func(t *testing.T, tp TimeParam) {
				now := time.Now()
				expected := now.Add(-14 * 24 * time.Hour)
				if tp.Time.Sub(expected).Abs() > time.Second {
					t.Errorf("Expected approximately %v, got %v", expected, tp.Time)
				}
			},
		},
		{
			name:      "Day duration - 1 day",
			input:     `1d`,
			wantError: false,
			validate: func(t *testing.T, tp TimeParam) {
				now := time.Now()
				expected := now.Add(-24 * time.Hour)
				if tp.Time.Sub(expected).Abs() > time.Second {
					t.Errorf("Expected approximately %v, got %v", expected, tp.Time)
				}
			},
		},
		{
			name:      "Day duration - 7 days",
			input:     `7d`,
			wantError: false,
			validate: func(t *testing.T, tp TimeParam) {
				now := time.Now()
				expected := now.Add(-7 * 24 * time.Hour)
				if tp.Time.Sub(expected).Abs() > time.Second {
					t.Errorf("Expected approximately %v, got %v", expected, tp.Time)
				}
			},
		},
		{
			name:      "Standard duration - hours",
			input:     `2h`,
			wantError: false,
			validate: func(t *testing.T, tp TimeParam) {
				now := time.Now()
				expected := now.Add(-2 * time.Hour)
				if tp.Time.Sub(expected).Abs() > time.Second {
					t.Errorf("Expected approximately %v, got %v", expected, tp.Time)
				}
			},
		},
		{
			name:      "Standard duration - minutes",
			input:     `30m`,
			wantError: false,
			validate: func(t *testing.T, tp TimeParam) {
				now := time.Now()
				expected := now.Add(-30 * time.Minute)
				if tp.Time.Sub(expected).Abs() > time.Second {
					t.Errorf("Expected approximately %v, got %v", expected, tp.Time)
				}
			},
		},
		{
			name:      "Standard duration - complex",
			input:     `2h30m15s`,
			wantError: false,
			validate: func(t *testing.T, tp TimeParam) {
				now := time.Now()
				expected := now.Add(-(2*time.Hour + 30*time.Minute + 15*time.Second))
				if tp.Time.Sub(expected).Abs() > time.Second {
					t.Errorf("Expected approximately %v, got %v", expected, tp.Time)
				}
			},
		},
		{
			name:      "Zero duration",
			input:     `0s`,
			wantError: false,
			validate: func(t *testing.T, tp TimeParam) {
				now := time.Now()
				if tp.Time.Sub(now).Abs() > time.Second {
					t.Errorf("Expected approximately %v, got %v", now, tp.Time)
				}
			},
		},
		{
			name:      "Invalid week duration - non-numeric",
			input:     `abcw`,
			wantError: false,
			validate: func(t *testing.T, tp TimeParam) {
				assert.Equal(t, tp.Time, time.Time{})
			},
		},
		{
			name:      "Invalid day duration - non-numeric",
			input:     `abcd`,
			wantError: false,
			validate: func(t *testing.T, tp TimeParam) {
				assert.Equal(t, tp.Time, time.Time{})
			},
		},
		{
			name:      "Invalid duration format",
			input:     `invalid`,
			wantError: false,
			validate: func(t *testing.T, tp TimeParam) {
				assert.Equal(t, tp.Time, time.Time{})
			},
		},
		{
			name:      "Invalid RFC3339 format",
			input:     `2023-13-45T25:70:90Z`,
			wantError: false,
			validate: func(t *testing.T, tp TimeParam) {
				assert.Equal(t, tp.Time, time.Time{})
			},
		},
		{
			name:      "Empty string",
			input:     ``,
			wantError: false,
			validate: func(t *testing.T, tp TimeParam) {
				assert.Equal(t, tp.Time, time.Time{})
			},
		},
		{
			name:      "Negative week duration",
			input:     `-1w`,
			wantError: false,
			validate: func(t *testing.T, tp TimeParam) {
				now := time.Now().Add(7 * 24 * time.Hour)
				if now.Sub(tp.Time).Abs() > time.Second {
					t.Errorf("Expected approximately %v, got %v", now, tp.Time)
				}
			},
		},
		{
			name:      "Negative day duration",
			input:     `-1d`,
			wantError: false,
			validate: func(t *testing.T, tp TimeParam) {
				now := time.Now().Add(24 * time.Hour)
				if now.Sub(tp.Time).Abs() > time.Second {
					t.Errorf("Expected approximately %v, got %v", now, tp.Time)
				}
			},
		},
		{
			name:      "Week duration - zero",
			input:     `0w`,
			wantError: false,
			validate: func(t *testing.T, tp TimeParam) {
				now := time.Now()
				if tp.Time.Sub(now).Abs() > time.Second {
					t.Errorf("Expected approximately %v, got %v", now, tp.Time)
				}
			},
		},
		{
			name:      "Day duration - zero",
			input:     `0d`,
			wantError: false,
			validate: func(t *testing.T, tp TimeParam) {
				now := time.Now()
				if tp.Time.Sub(now).Abs() > time.Second {
					t.Errorf("Expected approximately %v, got %v", now, tp.Time)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tp TimeParam
			err := tp.UnmarshalJSON([]byte(tt.input))

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tt.validate != nil {
				tt.validate(t, tp)
			}
		})
	}
}

func TestTimeParam_UnmarshalJSON_Integration(t *testing.T) {
	// Test integration with JSON unmarshaling
	tests := []struct {
		name      string
		jsonInput string
		wantError bool
	}{
		{
			name:      "JSON with RFC3339 time",
			jsonInput: `{"start_time": "2023-12-25T10:30:00Z"}`,
			wantError: false,
		},
		{
			name:      "JSON with week duration",
			jsonInput: `{"start_time": "2w"}`,
			wantError: false,
		},
		{
			name:      "JSON with day duration",
			jsonInput: `{"start_time": "5d"}`,
			wantError: false,
		},
		{
			name:      "JSON with standard duration",
			jsonInput: `{"start_time": "2h30m"}`,
			wantError: false,
		},
		{
			name:      "JSON with invalid duration",
			jsonInput: `{"start_time": "invalid"}`,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var params struct {
				StartTime TimeParam `json:"start_time"`
			}

			err := json.Unmarshal([]byte(tt.jsonInput), &params)

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			//
			//// Verify that time was properly set
			//if params.StartTime.Time.IsZero() {
			//	t.Errorf("Time was not properly set")
			//}
		})
	}
}

// Test cases specifically for quoted JSON strings
func TestTimeParam_UnmarshalJSON_QuotedStrings(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
		validate  func(t *testing.T, tp TimeParam)
	}{
		{
			name:      "Quoted RFC3339 format",
			input:     `"2023-12-25T10:30:00Z"`,
			wantError: false,
			validate: func(t *testing.T, tp TimeParam) {
				expected := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)
				if !tp.Time.Equal(expected) {
					t.Errorf("Expected %v, got %v", expected, tp.Time)
				}
			},
		},
		{
			name:      "Quoted week duration",
			input:     `"1w"`,
			wantError: false,
			validate: func(t *testing.T, tp TimeParam) {
				now := time.Now()
				expected := now.Add(-7 * 24 * time.Hour)
				if tp.Time.Sub(expected).Abs() > time.Second {
					t.Errorf("Expected approximately %v, got %v", expected, tp.Time)
				}
			},
		},
		{
			name:      "Quoted day duration",
			input:     `"1d"`,
			wantError: false,
			validate: func(t *testing.T, tp TimeParam) {
				now := time.Now()
				expected := now.Add(-24 * time.Hour)
				if tp.Time.Sub(expected).Abs() > time.Second {
					t.Errorf("Expected approximately %v, got %v", expected, tp.Time)
				}
			},
		},
		{
			name:      "Quoted standard duration",
			input:     `"2h"`,
			wantError: false,
			validate: func(t *testing.T, tp TimeParam) {
				now := time.Now()
				expected := now.Add(-2 * time.Hour)
				if tp.Time.Sub(expected).Abs() > time.Second {
					t.Errorf("Expected approximately %v, got %v", expected, tp.Time)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tp TimeParam
			err := tp.UnmarshalJSON([]byte(tt.input))

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tt.validate != nil {
				tt.validate(t, tp)
			}
		})
	}
}

func TestNewTimeParam(t *testing.T) {
	testTime := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)
	tp := NewTimeParam(testTime)

	if !tp.Time.Equal(testTime) {
		t.Errorf("Expected %v, got %v", testTime, tp.Time)
	}
}
