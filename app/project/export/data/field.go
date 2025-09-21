package data

import (
	"strconv"
	"strings"
)

const unknownKey = "Unknown"

type Field struct {
	Name        string `json:"name"`
	RowID       int    `json:"rowID,omitzero"`
	Type        string `json:"type"`
	Key         string `json:"key,omitzero"`
	KeyID       string `json:"key_id,omitzero"`
	KeyOffset   int    `json:"key_offset,omitzero"`
	Enum        any    `json:"enum,omitzero"`
	Unique      bool   `json:"unique,omitzero"`
	FilePath    bool   `json:"file_path,omitzero"`
	FileExt     any    `json:"file_ext,omitzero"`
	Display     any    `json:"display,omitzero"`
	DisplayType string `json:"display_type,omitzero"`
	Description string `json:"description,omitzero"`
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
