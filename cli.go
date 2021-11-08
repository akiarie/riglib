package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/pborman/getopt/v2"
)

type name []string

func (n name) String() string {
	return strings.Join(n, " ")
}

type definition struct {
	Name     name
	Meanings []string
	Children []definition
}

func parsedef(fname string, data []byte) (*definition, error) {
	lines := strings.Split(string(data), "\n")
	if len(lines) < 1 {
		return nil, fmt.Errorf("Definition must have at least one line")
	}

	def := &definition{
		Name:     []string{fname},
		Meanings: make([]string, 0),
		Children: make([]definition, 0),
	}
	curdef := def

	newdef := false
	start := 0
	// check for custom spelling in first word
	if l0 := strings.TrimSpace(lines[0]); l0[0] == '.' {
		if len(l0) > 1 && l0[1] == ' ' {
			l0 = fmt.Sprintf(".%s%s", def.Name, l0[1:])
		}
		def.Name = strings.Split(l0[1:], " ")
		start = 1
	}
	for i := start; i < len(lines); i++ {
		l := strings.TrimSpace(lines[i])
		if l == "" {
			if strings.TrimSpace(strings.Join(lines[i:], "")) == "" {
				break
			}
			newdef = true
			continue
		}
		if newdef {
			def.Children = append(def.Children, definition{Name: strings.Split(l, " ")})
			curdef = &def.Children[len(def.Children)-1]
			newdef = false
			continue
		}
		curdef.Meanings = append(curdef.Meanings, l)
		continue
	}
	return def, nil
}

func readdict(p string) ([]definition, error) {
	if fi, err := os.Stat(p); err != nil {
		return nil, fmt.Errorf("Unable to get path stats: %s", err)
	} else if !fi.Mode().IsDir() {
		return nil, fmt.Errorf("Dictionary path must be directory")
	}

	files, err := ioutil.ReadDir(p)
	if err != nil {
		return nil, fmt.Errorf("Error reading dictionary path: %s", err)
	}

	defs := []definition{}
	for _, fi := range files {
		if !fi.Mode().IsRegular() || strings.Contains(fi.Name(), ".swp") {
			continue
		}
		data, err := ioutil.ReadFile(filepath.Join(p, fi.Name()))
		if err != nil {
			return nil, fmt.Errorf("Unable to read file: %s", err)
		}
		def, err := parsedef(fi.Name(), data)
		if err != nil {
			return nil, fmt.Errorf("Cannot parse file %s: %s", filepath.Join(p, fi.Name()), err)
		}
		defs = append(defs, *def)
	}
	return defs, nil
}

const (
	FORMAT_PATH string = "format.tex"
	LDELIM             = "{{"
	RDELIM             = "}}"
)

func main() {
	formatfile := getopt.StringLong("format-file", 'f', FORMAT_PATH, "path to the format template file")
	ldelim := getopt.StringLong("left-delimiter", 'l', LDELIM, "the left-delimiter for template actions")
	rdelim := getopt.StringLong("right-delimiter", 'r', RDELIM, "the right-delimiter for template actions")
	getopt.Parse()
	args := getopt.Args()
	if len(args) == 0 {
		fmt.Printf("You must include the path: %v\n", args)
		os.Exit(1)
	}

	// read template
	var fns = template.FuncMap{
		"last": func(x int, a interface{}) bool {
			return x == reflect.ValueOf(a).Len()-1
		},
	}
	tmpl, err := template.New(path.Base(*formatfile)).Delims(*ldelim, *rdelim).Funcs(fns).ParseFiles(*formatfile)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	defs, err := readdict(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	for _, def := range defs {
		if err := tmpl.Execute(os.Stdout, def); err != nil {
			fmt.Println(err)
			os.Exit(4)
		}
	}
}
