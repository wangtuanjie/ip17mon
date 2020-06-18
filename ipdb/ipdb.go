package ipdb

import (
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"net"
	"strings"

	. "github.com/wangtuanjie/ip17mon/internal/proto"
)

func New(dataFile string) (Locator, error) {

	b, err := ioutil.ReadFile(dataFile)
	if err != nil {
		return nil, err
	}
	return NewWith(b)
}

func NewWith(b []byte) (Locator, error) {

	loc := new(locator)
	if err := loc.init(b); err != nil {
		return nil, err
	}
	return loc, nil
}

type locator struct {
	data []byte
	md   metadata

	v4Offset uint32

	countryField int
	regionField  int
	cityField    int
	ispField     int
}

type metadata struct {
	Build     int64          `json:"build"`
	IPVersion uint16         `json:"ip_version"`
	Languages map[string]int `json:"languages"`
	NodeCount uint32         `json:"node_count"`
	TotalSize int            `json:"total_size"`
	Fields    []string       `json:"fields"`
}

func (l *locator) Find(ipstr string) (*LocationInfo, error) {

	ip := net.ParseIP(ipstr)
	if ip == nil {
		return nil, ErrUnsupportedIP
	}

	var node uint32
	bitCount := 128

	if ipv4 := ip.To4(); ipv4 != nil {
		if l.md.IPVersion&0x01 == 0 {
			return nil, ErrUnsupportedIP
		}
		node = l.v4Offset
		bitCount = 32
		ip = ipv4
	} else if l.md.IPVersion&0x02 == 0 {
		return nil, ErrUnsupportedIP
	}

	for i := 0; i < bitCount; i++ {
		if node > l.md.NodeCount {
			return l.newLocationInfo(node), nil
		}
		node = l.nextNode(node, ((0xFF&int(ip[i>>3]))>>uint(7-(i%8)))&1 == 1)
	}

	return nil, ErrUnsupportedIP
}

func (l *locator) init(b []byte) error {

	metaSize := int(binary.BigEndian.Uint32(b[:4]))
	b = b[4:]

	if err := json.Unmarshal(b[:metaSize], &l.md); err != nil {
		return err
	}
	l.data = b[metaSize:]

	var node uint32
	for i := 0; i < 96 && node < l.md.NodeCount; i++ {
		if i >= 80 {
			node = l.nextNode(node, true)
		} else {
			node = l.nextNode(node, false)
		}
	}
	l.v4Offset = node

	l.countryField = -1
	l.regionField = -1
	l.cityField = -1
	l.ispField = -1

	for i, f := range l.md.Fields {
		switch f {
		case "country_name":
			l.countryField = i
		case "region_name":
			l.regionField = i
		case "city_name":
			l.cityField = i
		case "isp_domain":
			l.ispField = i
		}
	}

	return nil
}

func (l *locator) nextNode(node uint32, right bool) uint32 {

	off := node * 8
	if right {
		off += 4
	}
	return binary.BigEndian.Uint32(l.data[off : off+4])
}

func (l *locator) newLocationInfo(node uint32) *LocationInfo {

	off := node - l.md.NodeCount + l.md.NodeCount*8
	size := uint32(binary.BigEndian.Uint16(l.data[off : off+2]))
	text := l.data[off+2 : off+2+size]
	fields := strings.Split(string(text), "\t")

	fieldOf := func(i int) string {
		if i >= 0 && i < len(fields) && fields[i] != "" {
			return fields[i]
		}
		return Null
	}

	return &LocationInfo{
		Country: fieldOf(l.countryField),
		Region:  fieldOf(l.regionField),
		City:    fieldOf(l.cityField),
		Isp:     fieldOf(l.ispField),
	}
}
