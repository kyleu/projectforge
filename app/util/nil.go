// Package util - Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"database/sql"
	"encoding/json"
	"time"
)

type NilBool struct {
	sql.NullBool
}

func (n NilBool) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Bool)
	}
	return json.Marshal(nil)
}

func (n *NilBool) UnmarshalJSON(data []byte) error {
	var x *bool
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		n.Valid = true
		n.Bool = *x
	} else {
		n.Valid = false
		n.Bool = false
	}
	return nil
}

type NilFloat64 struct {
	sql.NullFloat64
}

func (n NilFloat64) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Float64)
	}
	return json.Marshal(nil)
}

func (n *NilFloat64) UnmarshalJSON(data []byte) error {
	var x *float64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		n.Valid = true
		n.Float64 = *x
	} else {
		n.Valid = false
		n.Float64 = 0
	}
	return nil
}

type NilInt32 struct {
	sql.NullInt32
}

func (n NilInt32) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Int32)
	}
	return json.Marshal(nil)
}

func (n *NilInt32) UnmarshalJSON(data []byte) error {
	var x *int32
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		n.Valid = true
		n.Int32 = *x
	} else {
		n.Valid = false
		n.Int32 = 0
	}
	return nil
}

type NilInt64 struct {
	sql.NullInt64
}

func (n NilInt64) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Int64)
	}
	return json.Marshal(nil)
}

func (n *NilInt64) UnmarshalJSON(data []byte) error {
	var x *int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		n.Valid = true
		n.Int64 = *x
	} else {
		n.Valid = false
		n.Int64 = 0
	}
	return nil
}

type NilString struct {
	sql.NullString
}

func (n NilString) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.String)
	}
	return json.Marshal(nil)
}

func (n *NilString) UnmarshalJSON(data []byte) error {
	var x *string
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		n.Valid = true
		n.String = *x
	} else {
		n.Valid = false
		n.String = ""
	}
	return nil
}

type NilTime struct {
	sql.NullTime
}

func (n NilTime) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Time)
	}
	return json.Marshal(nil)
}

func (n *NilTime) UnmarshalJSON(data []byte) error {
	var x *time.Time
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		n.Valid = true
		n.Time = *x
	} else {
		n.Valid = false
		n.Time = time.Time{}
	}
	return nil
}
