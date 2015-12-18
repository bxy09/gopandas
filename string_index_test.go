package gopandas_test

import (
	"github.com/bxy09/gopandas"
	"testing"
)

func TestStringIndex(t *testing.T) {
	array := []string{"a", "b", "c"}
	si := gopandas.NewStringIndex(array, false)
	if si.Length() != len(array) {
		t.Fatal("length error")
	}
	for i, item := range array {
		if si.String(i) != item {
			t.Fatal("error on String ", i, ":", t)
		}
		if si.Index(item) != i {
			t.Fatal("error on Index ", i, ":", t)
		}
	}
	if si.Index("") != -1 {
		t.Fatal("error on not exist string")
	}
	if si.String(-1) != "" {
		t.Fatal("error on not valid inde")
	}
	if si.String(len(array)) != "" {
		t.Fatal("error on not valid inde")
	}
}
