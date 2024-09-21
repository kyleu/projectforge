package util

import (
	"database/sql"
	"encoding/json"
	"time"
)

var EmptyStruct = struct{}{}

type NilBool struct {
	sql.NullBool
}

func (n NilBool) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return ToJSONBytes(n.Bool, false), nil
	}
	return ToJSONBytes(nil, false), nil
}

func (n *NilBool) UnmarshalJSON(data []byte) error {
	var x *bool
	if err := FromJSON(data, &x); err != nil {
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
		return ToJSONBytes(n.Float64, false), nil
	}
	return ToJSONBytes(nil, false), nil
}

func (n *NilFloat64) UnmarshalJSON(data []byte) error {
	var x *float64
	if err := FromJSON(data, &x); err != nil {
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
		return ToJSONBytes(n.Int32, false), nil
	}
	return ToJSONBytes(nil, false), nil
}

func (n *NilInt32) UnmarshalJSON(data []byte) error {
	var x *int32
	if err := FromJSON(data, &x); err != nil {
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
		return ToJSONBytes(n.Int64, false), nil
	}
	return ToJSONBytes(nil, false), nil
}

func (n *NilInt64) UnmarshalJSON(data []byte) error {
	var x *int64
	if err := FromJSON(data, &x); err != nil {
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

type NilJSON struct {
	sql.Null[json.RawMessage]
}

func (n NilJSON) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.V, nil
	}
	return ToJSONBytes(nil, false), nil
}

func (n *NilJSON) UnmarshalJSON(data []byte) error {
	if data != nil {
		n.Valid = true
		n.V = data
	} else {
		n.Valid = false
		n.V = nil
	}
	return nil
}

type NilString struct {
	sql.NullString
}

func (n NilString) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return ToJSONBytes(n.String, false), nil
	}
	return ToJSONBytes(nil, false), nil
}

func (n *NilString) UnmarshalJSON(data []byte) error {
	var x *string
	if err := FromJSON(data, &x); err != nil {
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
		return ToJSONBytes(n.Time, false), nil
	}
	return ToJSONBytes(nil, false), nil
}

func (n *NilTime) UnmarshalJSON(data []byte) error {
	var x *time.Time
	if err := FromJSON(data, &x); err != nil {
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
