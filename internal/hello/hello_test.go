package hello

import "testing"

func TestHelloWorld(t *testing.T) {
	got := HelloWorld()
	want := "hello world"

	if got != want {
		t.Fatalf("HelloWorld() = %q, want %q", got, want)
	}
}
