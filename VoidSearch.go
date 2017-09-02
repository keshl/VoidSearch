package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

var (
	regex = `[^>]*>([^<]*)-([^<-]*)_([0-9]+)\.([^.]*)\.xbps<\/a>.*\s([0-9]+)$`
	repos = []string{
		"https://repo.voidlinux.eu/current/",
		"https://repo.voidlinux.eu/current/multilib/",
		"https://repo.voidlinux.eu/current/nonfree/",
		"https://repo.voidlinux.eu/current/multilib/nonfree/",
		"https://repo.voidlinux.eu/current/debug/",
		"https://repo.voidlinux.eu/current/musl/",
		"https://repo.voidlinux.eu/current/musl/nonfree/",
		"https://repo.voidlinux.eu/current/musl/debug/",
	}
	repoNames = []string{
		"current",
		"multilib",
		"nonfree",
		"multilib/nonfree",
		"debug",
		"musl/current",
		"musl/nonfree",
		"musl/debug",
	}
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage:", os.Args[0], "<package>")
		os.Exit(1)
	}
	pak := os.Args[1]

	var compiled = regexp.MustCompile(regex)
	var bodyReader *bufio.Reader
	var table = tablewriter.NewWriter(os.Stdout)

	table.SetCenterSeparator("")
	table.SetRowSeparator("")
	table.SetColumnSeparator("")
	table.SetHeader([]string{"Repo", "Name", "Version", "Revision", "Arch", "Size"})

	ready := make(chan struct{})
	found := make(chan struct{})

	go func() {
		var n int = 0
		for {
			select {
			case <-ready:
				fmt.Println("\r\r")
				return
			case <-found:
				n++
			case <-time.After(time.Millisecond * 300):
				fmt.Printf("\r\r%d packages found.", n)
			}
		}
	}()

	for p := 0; p < len(repos); p++ {
		req, err := http.Get(repos[p])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		bodyReader = bufio.NewReader(req.Body)
		for {
			l, _, err := bodyReader.ReadLine()
			if err != nil {
				break
			}

			readed := string(l)

			exp := compiled.FindAllString(readed, -1)
			if len(exp) > 0 {
				s := regexp.MustCompile(regex).FindStringSubmatch(readed)
				if strings.Contains(s[1], pak) {
					found <- struct{}{}
					s[0] = repoNames[p]
					table.Append(s)
				}
			}
		}
	}
	ready <- struct{}{}
	table.Render()
}
