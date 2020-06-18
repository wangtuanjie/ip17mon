package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/wangtuanjie/ip17mon"
)

func stdin() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		ip := scanner.Text()
		if loc, err := ip17mon.Find(ip); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", ip, err)
		} else {
			fmt.Println(ip, loc.Country, loc.Region, loc.City, loc.Isp)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Failed:", err)
		os.Exit(1)
	}

	return
}

func main() {

	f := flag.String("f", "ipip_12_7.datx", "ip data file support dat/datx/ipdb format")
	flag.Parse()

	ip17mon.Init(*f)

	if args := flag.Args(); len(args) > 0 {
		if loc, err := ip17mon.Find(args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", args[0], err)
			os.Exit(1)
		} else {
			fmt.Println(loc.Country, loc.Region, loc.City, loc.Isp)
		}
	} else {
		stdin()
	}
}
