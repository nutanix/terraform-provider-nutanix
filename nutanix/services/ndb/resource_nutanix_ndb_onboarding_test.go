package ndb

import "testing"

func TestExpandStringList(t *testing.T) {
	in := []interface{}{"a", "", "b"}
	out := expandStringList(in)
	if len(out) != 2 {
		t.Fatalf("expected 2 values, got %d", len(out))
	}
	if out[0] != "a" || out[1] != "b" {
		t.Fatalf("unexpected output: %#v", out)
	}
}
