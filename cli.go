package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/pborman/getopt"
)

const (
	WORDS_DIR string = "/usr/local/words"

	MAX_FILES int = 20

	FO_DELAY = time.Millisecond * 1
)

type filelock struct {
	n  int
	mu sync.Mutex
}

func (f *filelock) open() {
	f.mu.Lock()
	if f.n >= MAX_FILES {
		f.mu.Unlock()
		time.Sleep(FO_DELAY)
		f.open()
		return
	}
	f.n++
	f.mu.Unlock()
}

func (f *filelock) close() {
	f.mu.Lock()
	f.n--
	f.mu.Unlock()
}

type wordset struct {
	k []whitword
	u []string
}

func root(word, wordsdir string, ch chan wordset, wg *sync.WaitGroup, fl *filelock) {
	fl.open()
	cmd := exec.Command("./words", word)
	cmd.Dir = wordsdir
	out, err := cmd.Output()
	if err != nil {
		panic(fmt.Sprintf("%s: open files %d", err, fl.n))
	}
	fl.close()
	words, unknown, err := parseword(word, string(out))
	if err != nil {
		panic(fmt.Sprintf("Unable to parse Whitaker definition: %v\n", err))
	}
	ch <- wordset{words, unknown}
	wg.Done()
}

func readlist(path string) ([]string, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Unable to read file")
	}
	firstlist := strings.Split(string(raw), "\n")
	list := []string{}
	for _, s := range firstlist {
		if len(strings.TrimSpace(s)) > 0 {
			list = append(list, s)
		}
	}
	return list, nil
}

func readwordsdir(dir string) error {
	if fi, err := os.Stat(dir); err != nil {
		return fmt.Errorf("Unable to get path stats: %s", err)
	} else if !fi.Mode().IsDir() {
		return fmt.Errorf("Words directory given not directory")
	}
	return nil
}

func main() {
	wordsdir := getopt.StringLong("words-dir", 'w', WORDS_DIR, "path to the words directory")
	prunknown := getopt.BoolLong("print-unknown", 'u', "whether to print unknown words")
	getopt.Parse()
	args := getopt.Args()
	if len(args) == 0 {
		log.Fatalln("You must input at least one book path")
	}
	if err := readwordsdir(*wordsdir); err != nil {
		log.Fatalln(err)
	}
	for _, p := range args {
		list, err := readlist(p)
		if err != nil {
			log.Fatalf("Unable to parse file '%s': %v\n", p, err)
		}
		fl := filelock{0, sync.Mutex{}}
		wg := sync.WaitGroup{}
		ch := make(chan wordset)
		for _, raww := range list {
			wg.Add(1)
			go root(raww, *wordsdir, ch, &wg, &fl)
		}
		go func() {
			wg.Wait()
			close(ch)
		}()
		for ws := range ch {
			for _, w := range ws.k {
				fmt.Println(w.inflection)
			}
			if *prunknown {
				for _, w := range ws.u {
					fmt.Println(w)
				}
			}
		}
	}
}
