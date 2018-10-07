package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"

	"github.com/lovego/gospec/problems"
	"github.com/lovego/gospec/rules"
)

func main() {
	traDirs, dirs, files := processArgs()
	for _, dir := range traDirs {
		traverseDir(dir)
	}
	for _, dir := range dirs {
		checkDir(dir)
	}
	if len(files) > 0 {
		checkFiles(files)
	}

	if problems.Count() > 0 {
		problems.Render()
		os.Exit(1)
	}
}

func traverseDir(dir string) {
	f, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	list, err := f.Readdir(-1)
	if err != nil {
		panic(err)
	}
	for _, d := range list {
		if d.IsDir() && d.Name()[0] != '.' {
			traverseDir(filepath.Join(dir, d.Name()))
		}
	}
	checkDir(dir)
}

func checkDir(dir string) {
	f, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	names, err := f.Readdirnames(-1)
	if err != nil {
		panic(err)
	}
	files := make([]string, 0, len(names))
	for _, name := range names {
		if willBuild(name) {
			files = append(files, filepath.Join(dir, name))
		}
	}
	if len(files) > 0 {
		rules.Check(dir, files)
	}
}

func checkFiles(paths []string) {
	dirs := make(map[string][]string)
	for _, p := range paths {
		dir := filepath.Dir(p)
		dirs[dir] = append(dirs[dir], p)
	}
	for dir, files := range dirs {
		rules.Check(dir, files)
	}
}

func processArgs() (traDirs, dirs, files []string) {
	for _, path := range flag.Args() {
		traverse := strings.HasSuffix(path, "/...")
		if traverse {
			path = strings.TrimSuffix(path, "/...")
		}

		switch mode := fileMode(path); {
		case mode.IsDir():
			if traverse {
				traDirs = append(traDirs, filepath.Clean(path))
			} else {
				dirs = append(dirs, filepath.Clean(path))
			}
		case mode.IsRegular():
			if willBuild(path) {
				files = append(files, filepath.Clean(path))
			}
		}
	}
	return
}

func willBuild(name string) bool {
	return filepath.Ext(name) == `.go` && filepath.Base(name)[0] != '.' && name[0] != '_'
}

func fileMode(path string) os.FileMode {
	if fi, err := os.Stat(path); err == nil {
		return fi.Mode()
	} else {
		panic(err)
	}
}
