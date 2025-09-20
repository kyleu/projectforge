package util_test

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"projectforge.dev/projectforge/app/util"
)

type parseTest[T any] struct {
	name       string
	input      any
	path       string
	allowEmpty bool
	want       T
	wantErr    bool
}

func TestParseBool(t *testing.T) {
	t.Parallel()
	tests := []parseTest[bool]{
		{"valid bool true", true, "test", false, true, false},
		{"valid bool false", false, "test", false, false, false},
		{"valid string true", "true", "test", false, true, false},
		{"valid string false", "false", "test", false, false, false},
		{"nil with allowEmpty", nil, "test", true, false, false},
		{"nil without allowEmpty", nil, "test", false, false, true},
		{"invalid type", 123, "test", false, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := util.ParseBool(tt.input, tt.path, tt.allowEmpty)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseInt64(t *testing.T) {
	t.Parallel()
	tests := []parseTest[int64]{
		{"valid int", 42, "test", false, 42, false},
		{"valid int32", int32(42), "test", false, 42, false},
		{"valid int64", int64(42), "test", false, 42, false},
		{"valid float64", 42.0, "test", false, 42, false},
		{"valid string", "42", "test", false, 42, false},
		{"invalid string", "not a number", "test", false, 0, true},
		{"nil with allowEmpty", nil, "test", true, 0, false},
		{"nil without allowEmpty", nil, "test", false, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := util.ParseInt64(tt.input, tt.path, tt.allowEmpty)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseTime(t *testing.T) {
	t.Parallel()
	now := util.TimeCurrent()
	tests := []parseTest[*time.Time]{
		{"valid time", now, "test", false, &now, false},
		{"valid pointer", &now, "test", false, &now, false},
		{"invalid string", "not a time", "test", false, nil, true},
		{"nil with allowEmpty", nil, "test", true, nil, false},
		{"nil without allowEmpty", nil, "test", false, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := util.ParseTime(tt.input, tt.path, tt.allowEmpty)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil && tt.want != nil {
				if !got.Equal(*tt.want) {
					t.Errorf("ParseTime() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestParseUUID(t *testing.T) {
	t.Parallel()
	validUUID := uuid.New()
	tests := []parseTest[*uuid.UUID]{
		{"valid UUID", validUUID, "test", false, &validUUID, false},
		{"valid UUID pointer", &validUUID, "test", false, &validUUID, false},
		{"valid string", validUUID.String(), "test", false, &validUUID, false},
		{"invalid string", "not-a-uuid", "test", false, nil, true},
		{"nil with allowEmpty", nil, "test", true, nil, false},
		{"nil without allowEmpty", nil, "test", false, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := util.ParseUUID(tt.input, tt.path, tt.allowEmpty)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil && tt.want != nil {
				if *got != *tt.want {
					t.Errorf("ParseUUID() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestParseMap(t *testing.T) {
	t.Parallel()
	validMap := util.ValueMap{"key": "value"}
	tests := []parseTest[util.ValueMap]{
		{"valid map", validMap, "test", false, validMap, false},
		{"empty map not allowed", util.ValueMap{}, "test", false, nil, true},
		{"empty map allowed", util.ValueMap{}, "test", true, util.ValueMap{}, false},
		{"valid JSON string", `{"key":"value"}`, "test", false, validMap, false},
		{"invalid JSON string", "not-json", "test", false, nil, true},
		{"nil with allowEmpty", nil, "test", true, nil, false},
		{"nil without allowEmpty", nil, "test", false, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := util.ParseMap(tt.input, tt.path, tt.allowEmpty)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("ParseMap() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
