package crypto

import (
	"testing"
)

func TestHash(t *testing.T) {
	s := "test"
	want := "36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab80"
	result := Hash(s)
	if result != want {
		t.Errorf("Hash(%s) = %s, want %s", s, result, want)
	}
}
