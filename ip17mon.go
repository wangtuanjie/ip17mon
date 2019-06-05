package ip17mon

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io/ioutil"
	"net"
	"strings"
)

const Null = "N/A"

var (
	ErrInvalidIp = errors.New("invalid ip format")
	std          *Locator
)

// Init defaut locator with dataFile
func Init(dataFile string) (err error) {
	if std != nil {
		return
	}
	std, err = NewLocator(dataFile)
	return
}

// Init defaut locator with data
func InitWithData(data []byte) {
	if std != nil {
		return
	}
	std = NewLocatorWithData(data)
	return
}

// Init defaut locator with datx
func InitWithDataX(data []byte) {
	if std != nil {
		return
	}
	std = NewLocatorWithDataX(data)
	return
}

// Find locationInfo by ip string
// It will return err when ipstr is not a valid format
func Find(ipstr string) (*LocationInfo, error) {
	return std.Find(ipstr)
}

// Find locationInfo by uint32
func FindByUint(ip uint32) *LocationInfo {
	return std.FindByUint(ip)
}

func Dump() (rs []Range, locs []*LocationInfo) {
	return std.Dump()
}

//-----------------------------------------------------------------------------

// New locator with dataFile
func NewLocator(dataFile string) (loc *Locator, err error) {
	data, err := ioutil.ReadFile(dataFile)
	if err != nil {
		return
	}
	if strings.HasSuffix(dataFile, ".datx") {
		loc = NewLocatorWithDataX(data)
	} else {
		loc = NewLocatorWithData(data)
	}
	return
}

// New locator with data
func NewLocatorWithData(data []byte) (loc *Locator) {
	loc = new(Locator)
	loc.init(data)
	return
}

// New locator with datax
func NewLocatorWithDataX(data []byte) (loc *Locator) {
	loc = new(Locator)
	loc.initX(data)
	return
}

type Locator struct {
	index           [256]int
	indexData       []uint32
	textStartIndex  []int
	textLengthIndex []int
	textData        []byte
}

type LocationInfo struct {
	Country string
	Region  string
	City    string
	Isp     string
}

type Range struct {
	Start net.IP
	End   net.IP
}

// Find locationInfo by ip string
// It will return err when ipstr is not a valid format
func (loc *Locator) Find(ipstr string) (info *LocationInfo, err error) {
	ip := net.ParseIP(ipstr)
	if ip == nil || ip.To4() == nil {
		err = ErrInvalidIp
		return
	}
	info = loc.FindByUint(binary.BigEndian.Uint32([]byte(ip.To4())))
	return
}

// Find locationInfo by uint32
func (loc *Locator) FindByUint(ip uint32) (info *LocationInfo) {

	idx := loc.findTextIndex(ip, loc.index[ip>>24])
	start := loc.textStartIndex[idx]
	return newLocationInfo(loc.textData[start : start+loc.textLengthIndex[idx]])
}

func (loc *Locator) Dump() (rs []Range, locs []*LocationInfo) {

	rs = make([]Range, 0, len(loc.indexData))
	locs = make([]*LocationInfo, 0, len(loc.indexData))

	for i := 1; i < len(loc.indexData); i++ {
		s, e := loc.indexData[i-1], loc.indexData[i]
		off := loc.textStartIndex[i]
		l := newLocationInfo(loc.textData[off : off+loc.textLengthIndex[i]])
		rs = append(rs, Range{ipOf(s), ipOf(e)})
		locs = append(locs, l)
	}
	return
}

func ipOf(n uint32) net.IP {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, n)
	return net.IP(b)
}

// binary search
func (loc *Locator) findTextIndex(ip uint32, start int) int {

	end := len(loc.indexData) - 1
	for start < end {
		mid := (start + end) / 2
		if ip > loc.indexData[mid] {
			start = mid + 1
		} else {
			end = mid
		}
	}

	if loc.indexData[end] >= ip {
		return end
	} else {
		return start
	}

}

func (loc *Locator) init(data []byte) {

	offset := int(binary.BigEndian.Uint32(data[:4]))
	textOff := offset - 1024

	loc.textData = data[textOff:]
	for i := 0; i < 256; i++ {
		off := 4 + i*4
		loc.index[i] = int(binary.LittleEndian.Uint32(data[off : off+4]))
	}

	nidx := (textOff - 4 - 1024) / 8

	loc.indexData = make([]uint32, nidx)
	loc.textStartIndex = make([]int, nidx)
	loc.textLengthIndex = make([]int, nidx)

	for i := 0; i < nidx; i++ {
		off := 4 + 1024 + i*8
		loc.indexData[i] = binary.BigEndian.Uint32(data[off : off+4])
		loc.textStartIndex[i] = int(uint32(data[off+4]) | uint32(data[off+5])<<8 | uint32(data[off+6])<<16)
		loc.textLengthIndex[i] = int(data[off+7])
	}
	return
}

// datx format
func (loc *Locator) initX(data []byte) {

	offset := int(binary.BigEndian.Uint32(data[:4]))
	textOff := offset - 256*256*4
	loc.textData = data[textOff:]
	for i := 0; i < 256; i++ {
		// datx格式使用了ipv4的前两个字节做为索引字段，出于对data格式兼容考虑这里只使用首字节做为索引字段
		// 由于我们使用二分查找, 这个点上认为对性能不会有任何影响
		off := 4 + i*256*4
		loc.index[i] = int(binary.LittleEndian.Uint32(data[off : off+4]))
	}

	nidx := (textOff - 4 - 256*256*4) / 9

	loc.indexData = make([]uint32, nidx)
	loc.textStartIndex = make([]int, nidx)
	loc.textLengthIndex = make([]int, nidx)

	for i := 0; i < nidx; i++ {
		off := 4 + 256*256*4 + i*9
		loc.indexData[i] = binary.BigEndian.Uint32(data[off : off+4])
		loc.textStartIndex[i] = int(uint32(data[off+4]) | uint32(data[off+5])<<8 | uint32(data[off+6])<<16)
		loc.textLengthIndex[i] = int(uint32(data[off+8]) | uint32(data[off+7])<<8)
	}
	return
}

func newLocationInfo(str []byte) *LocationInfo {

	var info *LocationInfo

	fields := bytes.Split(str, []byte("\t"))
	if len(fields) < 4 {
		panic("unexpected ip info:" + string(str))
	}
	info = &LocationInfo{
		Country: string(fields[0]),
		Region:  string(fields[1]),
		City:    string(fields[2]),
	}
	if len(fields) >= 5 {
		info.Isp = string(fields[4])
	}

	{
		if len(info.Country) == 0 {
			info.Country = Null
		}
		if len(info.Region) == 0 {
			info.Region = Null
		}
		if len(info.City) == 0 {
			info.City = Null
		}
		if len(info.Isp) == 0 {
			info.Isp = Null
		}
	}

	return info
}
