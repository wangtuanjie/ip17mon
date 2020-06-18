package datx

import (
	"testing"

	. "github.com/wangtuanjie/ip17mon/internal/proto"
)

const data = "17monipdb.dat"

func TestFind(t *testing.T) {
	l, err := New(data)
	if err != nil {
		t.Fatal("New failed:", err)
	}
	info, err := l.Find("115.231.237.124")

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
