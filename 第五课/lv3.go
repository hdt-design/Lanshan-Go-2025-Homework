//全文AI，不知道怎么写，要求也未完全看懂，此后会再抽时间理解

package main

import (
	"bufio"
	"fmt"
	"goprojects/第五课/mod"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ./catch [directory] [keyword]")
		return
	}

	rootdir := os.Args[1]
	keyword := os.Args[2]

	files := []string{}
	filepath.Walk(rootdir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	p := mod.NewMod(10, len(files))

	for _, f := range files {
		filepath := f
		p.Submit(func() {
			searchfile(filepath, keyword)
		})
	}

	p.Wait()
	p.Close()
}

func searchfile(filepath, keyword string) {
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	linenum := 1
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, keyword) {
			fmt.Printf("%s:%d:%s\n", filepath, linenum, line)
		}
		linenum++
	}
}
