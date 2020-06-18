package ipdb

import (
	"testing"
)

func TestIpdb(t *testing.T) {

	l, err := New("city.free.ipdb")
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

	if info.Isp != "N/A" {
		t.Fatal("isp expect = N/A, but actual =", info.Isp)
	}
}
