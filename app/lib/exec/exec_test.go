package exec_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/exec"
)

func TestWriteOutFns(t *testing.T) {
	t.Parallel()

	t.Run("no functions", func(t *testing.T) {
		t.Parallel()
		err := exec.WriteOutFns("key", []byte("data"))
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("single function", func(t *testing.T) {
		t.Parallel()
		var captured struct {
			key  string
			data []byte
		}
		fn := func(key string, b []byte) error {
			captured.key = key
			captured.data = b
			return nil
		}
		err := exec.WriteOutFns("testkey", []byte("testdata"), fn)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if captured.key != "testkey" {
			t.Errorf("expected key 'testkey', got %q", captured.key)
		}
		if string(captured.data) != "testdata" {
			t.Errorf("expected data 'testdata', got %q", string(captured.data))
		}
	})

	t.Run("multiple functions", func(t *testing.T) {
		t.Parallel()
		var calls []string
		fn1 := func(key string, _ []byte) error {
			calls = append(calls, "fn1:"+key)
			return nil
		}
		fn2 := func(key string, _ []byte) error {
			calls = append(calls, "fn2:"+key)
			return nil
		}
		err := exec.WriteOutFns("multi", []byte("data"), fn1, fn2)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(calls) != 2 {
			t.Errorf("expected 2 calls, got %d", len(calls))
		}
		if calls[0] != "fn1:multi" || calls[1] != "fn2:multi" {
			t.Errorf("unexpected call order: %v", calls)
		}
	})

	t.Run("error stops execution", func(t *testing.T) {
		t.Parallel()
		var calls int
		errFn := func(_ string, _ []byte) error {
			calls++
			return bytes.ErrTooLarge
		}
		fn2 := func(_ string, _ []byte) error {
			calls++
			return nil
		}
		err := exec.WriteOutFns("key", []byte("data"), errFn, fn2)
		if !errors.Is(err, bytes.ErrTooLarge) {
			t.Errorf("expected ErrTooLarge, got %v", err)
		}
		if calls != 1 {
			t.Errorf("expected 1 call (should stop on error), got %d", calls)
		}
	})
}

func TestWriteOutFnsString(t *testing.T) {
	t.Parallel()
	var captured struct {
		key  string
		data []byte
	}
	fn := func(key string, b []byte) error {
		captured.key = key
		captured.data = b
		return nil
	}
	exec.WriteOutFnsString("hello world", fn)
	if captured.key != "default" {
		t.Errorf("expected key 'default', got %q", captured.key)
	}
	if string(captured.data) != "hello world" {
		t.Errorf("expected data 'hello world', got %q", string(captured.data))
	}
}

func TestNewExec(t *testing.T) {
	t.Parallel()

	t.Run("simple key", func(t *testing.T) {
		t.Parallel()
		e := exec.NewExec("mykey", 1, "echo hello", "/tmp", false)
		if e.Key != "mykey" {
			t.Errorf("expected key 'mykey', got %q", e.Key)
		}
		if e.Idx != 1 {
			t.Errorf("expected idx 1, got %d", e.Idx)
		}
		if e.Cmd != "echo hello" {
			t.Errorf("expected cmd 'echo hello', got %q", e.Cmd)
		}
		if e.Path != "/tmp" {
			t.Errorf("expected path '/tmp', got %q", e.Path)
		}
		if e.Debug {
			t.Error("expected debug false")
		}
		if e.Buffer == nil {
			t.Error("expected buffer to be initialized")
		}
	})

	t.Run("key with slash prefix", func(t *testing.T) {
		t.Parallel()
		e := exec.NewExec("prefix/actualkey", 2, "ls", ".", true)
		if e.Key != "actualkey" {
			t.Errorf("expected key 'actualkey' (prefix stripped), got %q", e.Key)
		}
	})

	t.Run("with env vars", func(t *testing.T) {
		t.Parallel()
		e := exec.NewExec("key", 1, "cmd", ".", false, "FOO=bar", "BAZ=qux")
		if len(e.Env) != 2 {
			t.Errorf("expected 2 env vars, got %d", len(e.Env))
		}
		if e.Env[0] != "FOO=bar" || e.Env[1] != "BAZ=qux" {
			t.Errorf("unexpected env vars: %v", e.Env)
		}
	})
}

func TestExecWebPath(t *testing.T) {
	t.Parallel()

	t.Run("simple key", func(t *testing.T) {
		t.Parallel()
		e := exec.NewExec("build", 3, "make", ".", false)
		expected := "/admin/exec/build/3"
		if e.WebPath() != expected {
			t.Errorf("expected %q, got %q", expected, e.WebPath())
		}
	})

	t.Run("key with special characters", func(t *testing.T) {
		t.Parallel()
		e := exec.NewExec("my build", 1, "make", ".", false)
		expected := "/admin/exec/my+build/1"
		if e.WebPath() != expected {
			t.Errorf("expected %q, got %q", expected, e.WebPath())
		}
	})
}

func TestExecString(t *testing.T) {
	t.Parallel()
	e := exec.NewExec("test", 5, "cmd", ".", false)
	expected := "test:5"
	if e.String() != expected {
		t.Errorf("expected %q, got %q", expected, e.String())
	}
}

func TestExecsGet(t *testing.T) {
	t.Parallel()

	e1 := exec.NewExec("build", 1, "make", ".", false)
	e2 := exec.NewExec("build", 2, "make test", ".", false)
	e3 := exec.NewExec("deploy", 1, "deploy.sh", ".", false)
	execs := exec.Execs{e1, e2, e3}

	t.Run("found", func(t *testing.T) {
		t.Parallel()
		found := execs.Get("build", 2)
		if found != e2 {
			t.Errorf("expected e2, got %v", found)
		}
	})

	t.Run("not found - wrong key", func(t *testing.T) {
		t.Parallel()
		found := execs.Get("unknown", 1)
		if found != nil {
			t.Errorf("expected nil, got %v", found)
		}
	})

	t.Run("not found - wrong idx", func(t *testing.T) {
		t.Parallel()
		found := execs.Get("build", 99)
		if found != nil {
			t.Errorf("expected nil, got %v", found)
		}
	})
}

func TestExecsGetByKey(t *testing.T) {
	t.Parallel()

	e1 := exec.NewExec("build", 1, "make", ".", false)
	e2 := exec.NewExec("build", 2, "make test", ".", false)
	e3 := exec.NewExec("deploy", 1, "deploy.sh", ".", false)
	execs := exec.Execs{e1, e2, e3}

	t.Run("multiple matches", func(t *testing.T) {
		t.Parallel()
		builds := execs.GetByKey("build")
		if len(builds) != 2 {
			t.Errorf("expected 2 builds, got %d", len(builds))
		}
	})

	t.Run("single match", func(t *testing.T) {
		t.Parallel()
		deploys := execs.GetByKey("deploy")
		if len(deploys) != 1 {
			t.Errorf("expected 1 deploy, got %d", len(deploys))
		}
	})

	t.Run("no matches", func(t *testing.T) {
		t.Parallel()
		unknown := execs.GetByKey("unknown")
		if len(unknown) != 0 {
			t.Errorf("expected 0 matches, got %d", len(unknown))
		}
	})
}

func TestExecsRunning(t *testing.T) {
	t.Parallel()

	now := time.Now()
	e1 := exec.NewExec("build", 1, "make", ".", false)
	e2 := exec.NewExec("build", 2, "make test", ".", false)
	e2.Completed = &now
	e3 := exec.NewExec("deploy", 1, "deploy.sh", ".", false)
	execs := exec.Execs{e1, e2, e3}

	running := execs.Running()
	if running != 2 {
		t.Errorf("expected 2 running (not completed), got %d", running)
	}
}

func TestExecsSort(t *testing.T) {
	t.Parallel()

	e1 := exec.NewExec("zebra", 1, "cmd", ".", false)
	e2 := exec.NewExec("alpha", 2, "cmd", ".", false)
	e3 := exec.NewExec("alpha", 1, "cmd", ".", false)
	e4 := exec.NewExec("Beta", 1, "cmd", ".", false)
	execs := exec.Execs{e1, e2, e3, e4}

	execs.Sort()

	// Should be sorted case-insensitively by key, then by idx
	expected := []string{"alpha:1", "alpha:2", "Beta:1", "zebra:1"}
	for i, e := range execs {
		if e.String() != expected[i] {
			t.Errorf("position %d: expected %q, got %q", i, expected[i], e.String())
		}
	}
}
