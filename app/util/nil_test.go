//go:build test_all || !func_test
// +build test_all !func_test

package util_test

import (
	"strconv"
	"testing"

	"projectforge.dev/projectforge/app/util"
)

const jsonNull = "null"

func TestIsNil(t *testing.T) {
	t.Parallel()

	t.Run("nil interface is nil", func(t *testing.T) {
		t.Parallel()
		if !util.IsNil(nil) {
			t.Error("expected true for nil")
		}
	})

	t.Run("nil pointer is nil", func(t *testing.T) {
		t.Parallel()
		var p *int
		if !util.IsNil(p) {
			t.Error("expected true for nil pointer")
		}
	})

	t.Run("non-nil pointer is not nil", func(t *testing.T) {
		t.Parallel()
		x := 42
		if util.IsNil(&x) {
			t.Error("expected false for non-nil pointer")
		}
	})

	t.Run("nil slice is nil", func(t *testing.T) {
		t.Parallel()
		var s []int
		if !util.IsNil(s) {
			t.Error("expected true for nil slice")
		}
	})

	t.Run("empty slice is not nil", func(t *testing.T) {
		t.Parallel()
		s := []int{}
		if util.IsNil(s) {
			t.Error("expected false for empty slice")
		}
	})

	t.Run("nil map is nil", func(t *testing.T) {
		t.Parallel()
		var m map[string]int
		if !util.IsNil(m) {
			t.Error("expected true for nil map")
		}
	})

	t.Run("empty map is not nil", func(t *testing.T) {
		t.Parallel()
		m := map[string]int{}
		if util.IsNil(m) {
			t.Error("expected false for empty map")
		}
	})

	t.Run("nil channel is nil", func(t *testing.T) {
		t.Parallel()
		var c chan int
		if !util.IsNil(c) {
			t.Error("expected true for nil channel")
		}
	})

	t.Run("nil func is nil", func(t *testing.T) {
		t.Parallel()
		var f func()
		if !util.IsNil(f) {
			t.Error("expected true for nil func")
		}
	})

	t.Run("non-nil value types are not nil", func(t *testing.T) {
		t.Parallel()
		if util.IsNil(0) {
			t.Error("expected false for int")
		}
		if util.IsNil("") {
			t.Error("expected false for empty string")
		}
		if util.IsNil(false) {
			t.Error("expected false for bool")
		}
	})
}

func TestNilBool_JSON(t *testing.T) {
	t.Parallel()

	t.Run("marshal valid true", func(t *testing.T) {
		t.Parallel()
		nb := util.NilBool{}
		nb.Valid = true
		nb.Bool = true
		data := util.ToJSON(nb)
		if data != "true" {
			t.Errorf("expected 'true', got %s", data)
		}
	})

	t.Run("marshal invalid", func(t *testing.T) {
		t.Parallel()
		nb := util.NilBool{}
		data := util.ToJSON(nb)
		if data != jsonNull {
			t.Errorf("expected '%s', got %s", jsonNull, data)
		}
	})

	t.Run("unmarshal true", func(t *testing.T) {
		t.Parallel()
		var nb util.NilBool
		err := util.FromJSON([]byte("true"), &nb)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !nb.Valid || !nb.Bool {
			t.Errorf("expected valid true, got %+v", nb)
		}
	})

	t.Run("unmarshal null", func(t *testing.T) {
		t.Parallel()
		var nb util.NilBool
		err := util.FromJSON([]byte("null"), &nb)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if nb.Valid {
			t.Error("expected invalid")
		}
	})
}

func TestNilFloat64_JSON(t *testing.T) {
	t.Parallel()

	t.Run("marshal valid", func(t *testing.T) {
		t.Parallel()
		nf := util.NilFloat64{}
		nf.Valid = true
		nf.Float64 = 3.14
		data := util.ToJSON(nf)
		if data != "3.14" {
			t.Errorf("expected '3.14', got %s", data)
		}
	})

	t.Run("marshal invalid", func(t *testing.T) {
		t.Parallel()
		nf := util.NilFloat64{}
		data := util.ToJSON(nf)
		if data != jsonNull {
			t.Errorf("expected '%s', got %s", jsonNull, data)
		}
	})

	t.Run("unmarshal number", func(t *testing.T) {
		t.Parallel()
		var nf util.NilFloat64
		err := util.FromJSON([]byte("2.718"), &nf)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !nf.Valid || nf.Float64 != 2.718 {
			t.Errorf("expected valid 2.718, got %+v", nf)
		}
	})

	t.Run("unmarshal null", func(t *testing.T) {
		t.Parallel()
		var nf util.NilFloat64
		err := util.FromJSON([]byte("null"), &nf)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if nf.Valid {
			t.Error("expected invalid")
		}
	})
}

func TestNilInt32_JSON(t *testing.T) {
	t.Parallel()

	t.Run("marshal valid", func(t *testing.T) {
		t.Parallel()
		ni := util.NilInt32{}
		ni.Valid = true
		ni.Int32 = 42
		data := util.ToJSON(ni)
		if data != "42" {
			t.Errorf("expected '42', got %s", data)
		}
	})

	t.Run("unmarshal number", func(t *testing.T) {
		t.Parallel()
		var ni util.NilInt32
		err := util.FromJSON([]byte("123"), &ni)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ni.Valid || ni.Int32 != 123 {
			t.Errorf("expected valid 123, got %+v", ni)
		}
	})
}

func TestNilInt64_JSON(t *testing.T) {
	t.Parallel()

	t.Run("marshal valid", func(t *testing.T) {
		t.Parallel()
		ni := util.NilInt64{}
		ni.Valid = true
		ni.Int64 = 9223372036854775807
		data := util.ToJSON(ni)
		if data != "9223372036854775807" {
			t.Errorf("expected '9223372036854775807', got %s", data)
		}
	})

	t.Run("unmarshal number", func(t *testing.T) {
		t.Parallel()
		var ni util.NilInt64
		err := util.FromJSON([]byte("999"), &ni)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ni.Valid || ni.Int64 != 999 {
			t.Errorf("expected valid 999, got %+v", ni)
		}
	})
}

func TestNilString_JSON(t *testing.T) {
	t.Parallel()

	t.Run("marshal valid", func(t *testing.T) {
		t.Parallel()
		ns := util.NilString{}
		ns.Valid = true
		ns.String = boolTestHello
		data := util.ToJSON(ns)
		if data != strconv.Quote(boolTestHello) {
			t.Errorf("expected %q, got %s", strconv.Quote(boolTestHello), data)
		}
	})

	t.Run("marshal invalid", func(t *testing.T) {
		t.Parallel()
		ns := util.NilString{}
		data := util.ToJSON(ns)
		if data != jsonNull {
			t.Errorf("expected '%s', got %s", jsonNull, data)
		}
	})

	t.Run("unmarshal string", func(t *testing.T) {
		t.Parallel()
		var ns util.NilString
		err := util.FromJSON([]byte(`"world"`), &ns)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ns.Valid || ns.String != "world" {
			t.Errorf("expected valid 'world', got %+v", ns)
		}
	})

	t.Run("unmarshal null", func(t *testing.T) {
		t.Parallel()
		var ns util.NilString
		err := util.FromJSON([]byte("null"), &ns)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ns.Valid {
			t.Error("expected invalid")
		}
	})
}

func TestNilTime_JSON(t *testing.T) {
	t.Parallel()

	t.Run("marshal invalid", func(t *testing.T) {
		t.Parallel()
		nt := util.NilTime{}
		data := util.ToJSON(nt)
		if data != jsonNull {
			t.Errorf("expected '%s', got %s", jsonNull, data)
		}
	})

	t.Run("unmarshal null", func(t *testing.T) {
		t.Parallel()
		var nt util.NilTime
		err := util.FromJSON([]byte("null"), &nt)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if nt.Valid {
			t.Error("expected invalid")
		}
	})

	t.Run("unmarshal time string", func(t *testing.T) {
		t.Parallel()
		var nt util.NilTime
		err := util.FromJSON([]byte(`"2024-06-15T14:30:00Z"`), &nt)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !nt.Valid {
			t.Error("expected valid time")
		}
	})
}

func TestNilJSON_JSON(t *testing.T) {
	t.Parallel()

	t.Run("marshal invalid", func(t *testing.T) {
		t.Parallel()
		nj := util.NilJSON{}
		data := util.ToJSON(nj)
		if data != jsonNull {
			t.Errorf("expected '%s', got %s", jsonNull, data)
		}
	})

	t.Run("unmarshal object", func(t *testing.T) {
		t.Parallel()
		var nj util.NilJSON
		err := util.FromJSON([]byte(`{"key":"value"}`), &nj)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !nj.Valid {
			t.Error("expected valid")
		}
	})
}
