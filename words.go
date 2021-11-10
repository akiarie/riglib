package main

import (
	"fmt"
	"regexp"
	"strings"
)

type whitSpeechPart int

const (
	X      = iota // all, none, or unknown
	N             // Noun
	PRON          // PRONoun
	PACK          // PACKON -- artificial for code
	ADJ           // ADJective
	NUM           // NUMeral
	ADV           // ADVerb
	V             // Verb
	VPAR          // Verb PARticiple
	SUPINE        // SUPINE
	PREP          // PREPosition
	CONJ          // CONJunction
	INTERJ        // INTERJection
	TACKON        // TACKON --  artificial for code
	PREFIX        // PREFIX --  here artificial for code
	SUFFIX        //  SUFFIX --  here artificial for code
)

var whitnames = map[whitSpeechPart]string{
	X:      "X",
	N:      "N",
	PRON:   "PRON",
	PACK:   "PACK",
	ADJ:    "ADJ",
	NUM:    "NUM",
	ADV:    "ADV",
	V:      "V",
	VPAR:   "VPAR",
	SUPINE: "SUPINE",
	PREP:   "PREP",
	CONJ:   "CONJ",
	INTERJ: "INTERJ",
	TACKON: "TACKON",
	PREFIX: "PREFIX",
	SUFFIX: "SUFFIX",
}

func (p whitSpeechPart) String() string {
	return whitnames[p]
}

func genWhitSpeechParts() []string {
	arr := []string{}
	for _, v := range whitnames {
		arr = append(arr, v)
	}
	return arr
}

func genWhitSpeechInv() map[string]whitSpeechPart {
	inv := make(map[string]whitSpeechPart)
	for k, v := range whitnames {
		inv[v] = k
	}
	return inv
}

var whitspeechparts = genWhitSpeechParts()
var whitnamesinv = genWhitSpeechInv()

type whitword struct {
	inflection string
	english    []string
	part       whitSpeechPart
}

func (w whitword) String() string {
	if len(w.english) > 0 {
		return fmt.Sprintf("{%s, inflection: '%s', english: %v}", w.part, w.inflection, w.english)
	} else {
		return fmt.Sprintf("{%s, inflection: '%s'}", w.part, w.inflection)
	}
}

type stateFn func(*scannand) stateFn

type scannand struct {
	original     string
	input        []string // whitaker definition lines
	pos          int
	whitpart     *whitSpeechPart // whitaker word type of last line (if any)
	previnfl     string
	words        []whitword // parsed words
	unknownwords []string
}

func formwalk(sc *scannand) stateFn {
	if !(sc.pos < len(sc.input)) {
		return nil
	}
	l := strings.TrimSpace(sc.input[sc.pos])
	switch l {
	case "", "*":
		sc.pos++
		return formwalk
	}
	if matches := regexp.MustCompile(`([A-Za-z]+)\s+==+`).FindStringSubmatch(l); len(matches) == 2 {
		sc.unknownwords = append(sc.unknownwords, matches[1])
		sc.pos++
		return formwalk
	}
	pattern := fmt.Sprintf(`(^[A-Za-z,\.\(\)\- ]+)(?:\s+)(%s)(?:\A|\z|\s)`, strings.Join(whitspeechparts, "|"))
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(l)
	if len(matches) == 3 {
		if wp, ok := whitnamesinv[matches[2]]; ok {
			prevprevinfl := sc.previnfl
			sc.previnfl = strings.TrimSpace(strings.ReplaceAll(matches[1], ".", ""))
			if sc.whitpart == nil {
				sc.whitpart = &wp
				sc.pos++
				return formwalk
			}
			if *sc.whitpart != wp {
				// tackon fix
				if strings.Contains(sc.previnfl, "TACKON") {
					sc.previnfl = prevprevinfl
					wp = *sc.whitpart
				}
				// when there's a transition l is an inflection line
				sc.words = append(sc.words, whitword{
					inflection: sc.previnfl,
					part:       wp,
				})
				sc.whitpart = nil
				sc.pos++
				return formwalk
			}
			// then we know *sc.whitpart == wp
			sc.pos++
			return formwalk
		} else {
			panic(fmt.Sprintf("Invalid speech part: %s", matches[1]))
		}
	}

	for i := 0; i < len(sc.words); i++ {
		if sc.words[i].english == nil {
			sc.words[i].english = []string{l}
		} else {
			sc.words[i].english = append(sc.words[i].english, l)
		}
	}
	if sc.whitpart == nil {
		if len(sc.words) > 0 {
			sc.pos++
			return formwalk
		}
		ignored := []string{
			"syncop", "word mod", "an internal", "an initial",
			"two words", "may be", "slur", "bad roman numeral",
			"it is very", "a terminal",
		}
		for _, s := range ignored {
			if strings.Contains(strings.ToLower(l), s) {
				sc.pos++
				return formwalk
			}
		}
		panic(fmt.Sprintf("Halted on input '%s', line: %s\n", sc.original, l))
	}

	// we should be on English line, so set all previous unset, if necessary
	sc.words = append(sc.words, whitword{
		inflection: sc.previnfl,
		english:    []string{l},
		part:       *sc.whitpart,
	})
	sc.whitpart = nil
	sc.pos++
	return formwalk
}

func parseword(word, raw string) ([]whitword, []string, error) {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	scan := &scannand{
		original:     word,
		input:        lines,
		pos:          0,
		words:        []whitword{},
		unknownwords: []string{},
	}
	for st := formwalk(scan); st != nil; st = st(scan) {
	}
	return scan.words, scan.unknownwords, nil
}
