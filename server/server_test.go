package server

import (
	"fmt"
	"testing"
)

func TestRun(t *testing.T) {
	run("ls -ll")
}

func TestReadSettings(t *testing.T) {
	s, err := readSettings("../settings.json")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(s.Commands["ls"][1])
}

func BenchmarkReadSettings(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := readSettings("../settings.json")
		if err != nil {
			b.Fatal(err)
		}
	}

}