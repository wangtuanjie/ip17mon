package ip17mon

import (
	"math/rand"
	"testing"
	"time"
)

const data = "17monipdb.dat"

func TestFind(t *testing.T) {
	if err := Init(data); err != nil {
		t.Fatal("Init failed:", err)
	}
	info, err := Find("115.231.237.124")

	if err != nil {
		t.Fatal("Find failed:", err)
	}

	if info.Country != "中国" {
		t.Fatal("country expect = 中国, but actual =", info.Country)
	}

	if info.Region != "浙江" {
		t.Fatal("region expect = 浙江, but actual =", info.Region)
	}

	if info.City != "嘉兴" {
		t.Fatal("city expect = 嘉兴, but actual =", info.City)
	}

	if info.Isp != Null {
		t.Fatal("isp expect = Null, but actual =", info.Isp)
	}
}

//-----------------------------------------------------------------------------

// Benchmark command
//	go test -bench=Find
// Benchmark result
//	BenchmarkFind 2000000       830 ns/op
// about 120w/s
func BenchmarkFind(b *testing.B) {
	b.StopTimer()
	if err := Init(data); err != nil {
		b.Fatal("Init failed:", err)
	}
	rand.Seed(time.Now().UnixNano())
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		if FindByUint(rand.Uint32()) == nil {
			b.Fatal("FindByUint found nil val")
		}
	}
}
