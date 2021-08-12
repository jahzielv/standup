package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"
)

type excludeFlag map[string]bool
type includeFlag []string

func (i *excludeFlag) String() string {
	return "my string representation"
}

func (i *excludeFlag) Set(value string) error {
	names := strings.Split(value, ",")
	for _, n := range names {
		(*i)[n] = true
	}
	return nil
}

func (i *includeFlag) String() string {
	return "flag"
}

func (i *includeFlag) Set(value string) error {
	*i = strings.Split(value, ",")
	return nil
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func main() {
	var excludeNames excludeFlag = make(map[string]bool)
	var includeNames includeFlag = make([]string, 0)
	flag.Var(&excludeNames, "e", "exclude names from list")
	flag.Var(&includeNames, "i", "include names")
	flag.Parse()
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	lines, err := readLines(path.Join(home, ".config/standup/.standuprc"))
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	if len(includeNames) > 0 {
		lines = append(lines, includeNames...)
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	r1.Shuffle(len(lines), func(i, j int) { lines[i], lines[j] = lines[j], lines[i] })
	offset := 1
	for i, l := range lines {
		if excludeNames[l] {
			offset--
			continue
		}

		fmt.Println(i+offset, l)
	}

}
