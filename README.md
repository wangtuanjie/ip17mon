[17mon](http://www.ipip.net/) IP location data for Golang
===

[![Circle CI](https://circleci.com/gh/wangtuanjie/ip17mon.svg?style=svg)](https://circleci.com/gh/wangtuanjie/ip17mon)

## 特性

* dat/datx 只支持 ipv4 
* ipdb 支持 ipv4/ipv6

## 安装

	go get github.com/wangtuanjie/ip17mon@latest


## 使用
	import （
		"fmt"
		"github.com/wangtuanjie/ip17mon"
	）

	func init() {
		ip17mon.Init("your data file")
	}

	func main() {
		loc, err := ip17mon.Find("116.228.111.18")
		if err != nil {
			fmt.Println("err:", err)
			return
		}
		fmt.Println(loc)
	}

更多请参考[example](https://github.com/wangtuanjie/ip17mon/tree/master/cmd/qip)



## 许可证

基于 [MIT](https://github.com/wangtuanjie/ip17mon/blob/master/LICENSE) 协议发布

