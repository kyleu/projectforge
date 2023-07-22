package data

import (
	"strconv"
	"strings"
)

const unknownKey = "Unknown"

type Field struct {
	Name        string `json:"name"`
	RowID       int    `json:"rowID,omitempty"`
	Type        string `json:"type"`
	Key         string `json:"key,omitempty"`
	KeyID       string `json:"key_id,omitempty"`
	KeyOffset   int    `json:"key_offset,omitempty"`
	Enum        any    `json:"enum,omitempty"`
	Unique      bool   `json:"unique,omitempty"`
	FilePath    bool   `json:"file_path,omitempty"`
	FileExt     any    `json:"file_ext,omitempty"`
	Display     any    `json:"display,omitempty"`
	DisplayType string `json:"display_type,omitempty"`
	Description string `json:"description,omitempty"`
}

func (f *Field) IsUnknown() bool {
	return strings.HasPrefix(f.Name, unknownKey)
}

func (f *Field) UnknownIdx() int {
	if !strings.HasPrefix(f.Name, unknownKey) {
		return -1
	}
	s := strings.TrimPrefix(f.Name, unknownKey)
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return -2
	}
	return int(i)
}

type Fields []*Field
