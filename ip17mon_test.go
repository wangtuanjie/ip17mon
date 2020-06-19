package ip17mon

import (
	"testing"
)

func TestAll(t *testing.T) {

	l1, err := New("datx/17monipdb.dat")
	if err != nil {
		t.Fatal("New failed:", err)
	}

	l2, err := New("ipdb/city.free.ipdb")
	if err != nil {
		t.Fatal("New failed:", err)
	}

	ip := "115.231.237.124"

	info1, err := l1.Find(ip)
	if err != nil {
		t.Fatal("l1.Find failed:", err)
	}

	info2, err := l2.Find(ip)
	if err != nil {
		t.Fatal("l2.Find failed:", err)
	}

	if *info1 != *info2 {
		t.Fatalf("info: %v != %v", *info1, *info2)
	}
}
