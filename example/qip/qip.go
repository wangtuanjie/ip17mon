package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/wangtuanjie/ip17mon"
)

func init() {
	ip17mon.InitWithData(data)
}

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
	if len(os.Args) > 1 {
		if loc, err := ip17mon.Find(os.Args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[1], err)
			os.Exit(1)
		} else {
			fmt.Println(loc.Country, loc.Region, loc.City, loc.Isp)
		}
	} else {
		stdin()
	}
}
