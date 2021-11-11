package words

import (
	"fmt"
	"regexp"
	"strings"
)

type WhitSpeechPart int

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

var whitnames = map[WhitSpeechPart]string{
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

func (p WhitSpeechPart) String() string {
	return whitnames[p]
}

func genWhitSpeechParts() []string {
	arr := []string{}
	for _, v := range whitnames {
		arr = append(arr, v)
	}
	return arr
}

func genWhitSpeechInv() map[string]WhitSpeechPart {
	inv := make(map[string]WhitSpeechPart)
	for k, v := range whitnames {
		inv[v] = k
	}
	return inv
}

var WhitSpeechParts = genWhitSpeechParts()
var whitnamesinv = genWhitSpeechInv()

type Whitword struct {
	Inflection string
	english    []string
	part       WhitSpeechPart
}

func (w Whitword) String() string {
	if len(w.english) > 0 {
		return fmt.Sprintf("{%s, Inflection: '%s', english: %v}", w.part, w.Inflection, w.english)
	} else {
		return fmt.Sprintf("{%s, Inflection: '%s'}", w.part, w.Inflection)
	}
}

type stateFn func(*scannand) stateFn

type scannand struct {
	original     string
	input        []string // whitaker definition lines
	pos          int
	whitpart     *WhitSpeechPart // whitaker word type of last line (if any)
	previnfl     string
	words        []Whitword // parsed words
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
	pattern := fmt.Sprintf(`(^[A-Za-z,\.\(\)\- ]+)(?:\s+)(%s)(?:\A|\z|\s)`, strings.Join(WhitSpeechParts, "|"))
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
				// when there's a transition l is an Inflection line
				sc.words = append(sc.words, Whitword{
					Inflection: sc.previnfl,
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
	sc.words = append(sc.words, Whitword{
		Inflection: sc.previnfl,
		english:    []string{l},
		part:       *sc.whitpart,
	})
	sc.whitpart = nil
	sc.pos++
	return formwalk
}

func Parseword(word, raw string) ([]Whitword, []string, error) {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	scan := &scannand{
		original:     word,
		input:        lines,
		pos:          0,
		words:        []Whitword{},
		unknownwords: []string{},
	}
	for st := formwalk(scan); st != nil; st = st(scan) {
	}
	return scan.words, scan.unknownwords, nil
}
