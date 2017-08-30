package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"regexp"
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
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Package: ")
	_, _, err := reader.ReadLine()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var compiled = regexp.MustCompile(regex)
	var bodyReader *bufio.Reader
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

			fmt.Println(repoNames[p], compiled.FindAllString(string(l), -1))
		}
	}
}
