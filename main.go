package main

// from stdin, find most repetitive patterns enlight them
// v0.36 : change -S to specify separators
// v0.35 : add -S to choose usual separators splitting
// v0.34 : without func for assigning colors, and a color struct for init
// v0.32 : timeout and colors error
// v0.3 : multiple colors
// v0.2 : find multiple patterns in one line
// v0.1 : find only the longuest pattern within one line
// of course with a lot of bugs.. so more testing/debugging is always needed.
// ( Always check  Windows 10 : $Env:GOPATH )

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
)

var sep string

func main() {

	var (
		input                     string
		minlen, views             int
		patternsonly, notimelimit bool
		maxtimelimit              = 15 * time.Second
		patterns                  = make(map[string]int)
		output                    = make([]string, 0)
	)

	// Parsing args
	flag.IntVar(&minlen, "l", 7, "min pattern length")
	flag.IntVar(&views, "o", 3, "min occurrences")
	flag.BoolVar(&patternsonly, "P", false, "only print found patterns")
	flag.BoolVar(&notimelimit, "T", false, "no time limit")
	flag.StringVar(&sep, "S", "", "use separators \"space+\t,;/\" [default char]")
	flag.StringVar(&input, "i", "", "input [default:sdtin]")

	flag.Parse()

	if !notimelimit {
		go func() { // Don't wait forever if data is huge or complex.
			time.Sleep(maxtimelimit)
			error := color.New(color.BgRed, color.FgHiWhite).SprintFunc()
			errormsg := fmt.Sprintf("Sorry, but timeout is set to %2.fs, please retry with simplified data.", maxtimelimit.Seconds())
			fmt.Fprintf(color.Output, error(errormsg))
			fmt.Printf("\n")
			os.Exit(1)
		}()
	}
	// defining input
	var s *bufio.Scanner
	if input != "" {
		f, err := os.Open(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to open %s :%v", input, err)
			os.Exit(1)
		}
		s = bufio.NewScanner(f)
	} else {
		s = bufio.NewScanner(os.Stdin)
	}

	// Reading input , putting all in []string
	for s.Scan() {
		line := s.Text()
		output = append(output, line)
		if sep != "" { // faster
			splitusingseparators(minlen, line, patterns)
		} else { // much longer char by chars
			split(minlen, line, patterns)
		}
	}

	// trimming patterns
	buildmap(views, patterns)
	// replacing patterns map with hash ordered by length
	hash := rankByWordCount(patterns)

	// if we only want to see patterns
	if patternsonly {
		fmt.Printf(" #/len \t>Pattern by length (min len:%d #:%d)<\n", minlen, views)
		for _, v := range hash {
			fmt.Printf("%2d/%2d \t>%s<\n", v.Value, len(v.Key), v.Key)
		}
		os.Exit(0)
	}

	// Attributing colors to patterns
	idxcolors := initcolors()

	keycolors := make(map[string]*color.Color)
	for i, p := range hash {
		keycolors[p.Key] = idxcolors[i%len(idxcolors)]
	}

	// printing input with color
	for _, v := range output {
		works := v
		for len(works) > minlen {
			idx, bestPLen, bestPIdx, bestP := 0, 0, 0, ""
			// recherche de la chaine la plus longue
			for _, k := range hash {
				if idx = strings.Index(works, k.Key); (idx != -1) && ((idx < bestPIdx) || (len(k.Key) > bestPLen)) {
					bestP, bestPLen, bestPIdx = k.Key, len(k.Key), idx
				}
			}

			if bestPLen > 0 { // found a pattern
				fmt.Fprintf(color.Output, "%s%s",
					works[0:bestPIdx],
					keycolors[bestP].Sprintf(bestP))
				works = works[bestPIdx+bestPLen:]
			} else { // no more patterns
				fmt.Fprintf(color.Output, "%s", works)
				works = ""
			}

		} // end loop len(workstring)
		if len(works) != 0 {
			fmt.Fprintf(color.Output, "%s", works)
		}
		fmt.Fprintf(color.Output, "\n")
	} // end loop output
	os.Exit(0)
} // end of main()

// buildmap : trim duplicate values
func buildmap(views int, values map[string]int) {
	for k := range values { // fast loop first
		if values[k] < views {
			//fmt.Printf("deleting insufficient views k(%d)[%d] [%s] \n", values[k], len(k), k)
			delete(values, k)
		}
	}
	for k := range values { //N^2 loop
		//fmt.Printf(" %s %d \n", k, v)
		for k2 := range values {
			if values[k2] == values[k] && len(k2) != len(k) && strings.Contains(k2, k) {
				//fmt.Printf("deleting k(%d)[%d] found in k2(%d)[%d] : [%s] in [%s]\n", values[k], len(k), values[k2], len(k2), k, k2)
				delete(values, k)
				// break
			}
		}
	}

} // end of buildmap

// func separators(r rune) bool {
// 	return r == ' ' || r == '	' || r == '/' || r == ',' || r == ';'
// }

func separators(r rune) bool {
	for _, v := range sep { // sep is a string containing separators
		if r == rune(v) {
			return true
		}
	}
	return false
}

func splitusingseparators(minlen int, line string, patterns map[string]int) {
	if len(line) < minlen {
		return
	}
	parts := strings.FieldsFunc(line, separators)
	for i := range parts {
		if len(parts[i]) >= minlen {
			patterns[parts[i]]++
		}
	}
}

// split a string of minlen into a map[string]int
func split(minlen int, line string, patterns map[string]int) {
	l := len(line)
	if l < minlen {
		return
	}
	for j := 0; j <= l-minlen; j++ {
		for i := j + minlen; i <= l; i++ {
			patterns[line[j:i]]++
		}
	}
}

// initcolors : initialize an array with a few colors
func initcolors() []*color.Color {
	var h []*color.Color

	var mycolors = []struct {
		fg color.Attribute
		bg color.Attribute
	}{
		{color.FgWhite, color.BgGreen},
		{color.FgBlack, color.BgYellow},
		{color.FgWhite, color.BgBlue},
		{color.FgBlack, color.BgCyan},
		{color.FgWhite, color.BgRed},
		{color.FgRed, color.BgWhite},
		{color.FgBlue, color.BgWhite},
		{color.FgYellow, color.BgBlue},
	}

	for _, v := range mycolors {
		h = append(h, color.New(v.fg, v.bg))
	}
	return h
}

// dealing with sorted output
func rankByWordCount(wordFrequencies map[string]int) PairList {
	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	// sort.Sort(pl)

	return pl
}

// Pair ...
type Pair struct {
	Key   string
	Value int
}

// PairList ...
type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return len(p[i].Key) < len(p[j].Key) }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// end of main.go
