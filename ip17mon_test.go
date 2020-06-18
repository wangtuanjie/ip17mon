package ip17mon

import (
	"testing"
)

func TestAll(t *testing.T) {

	l1, err := New("datx/17monipdb.dat")
	if err != nil {
		t.Fatal(err)
	}

	l2, err := New("ipdb/city.free.ipdb")
	if err != nil {
		t.Fatal(err)
	}

	info1, err := l1.Find("115.231.237.124")
	if err != nil {
		t.Fatal(err)
	}

	info2, err := l2.Find("115.231.237.124")
	if err != nil {
		t.Fatal(err)
	}

	if *info1 != *info2 {
		t.Fatal("xxxoo")
	}
}
