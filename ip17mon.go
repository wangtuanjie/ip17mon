package ip17mon

import (
	"path/filepath"
	"strings"

	"github.com/wangtuanjie/ip17mon/datx"
	"github.com/wangtuanjie/ip17mon/ipdb"

	. "github.com/wangtuanjie/ip17mon/internal/proto"
)

var def Locator

func Init(dataFile string) {
	var err error
	def, err = New(dataFile)
	if err != nil {
		panic(err)
	}
}

func InitWithDatx(b []byte) {
	def = datx.NewWithDatx(b)
}

func InitWithIpdb(b []byte) {
	var err error
	def, err = ipdb.NewWith(b)
	println(def)
	if err != nil {
		panic(err)
	}
}

func Find(ipstr string) (*LocationInfo, error) {
	return def.Find(ipstr)
}

func New(dataFile string) (loc Locator, err error) {

	switch strings.ToLower(filepath.Ext(dataFile)) {
	case ".dat", ".datx":
		return datx.New(dataFile)
	case ".ipdb":
		return ipdb.New(dataFile)
	default:
		return nil, ErrUnsupportedFile
	}
}
