package main

// from stdin, find most repetitive patterns and enlight them
// v0.3 : multiple colors
// v0.2 : find multiple patterns in one line
// v0.1 : find only the longuest pattern within one line
// of course with a lot of bugs.. so more testing/debugging is always needed.

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func main() {

	var (
		input         string
		minlen, views int
		patternsonly  bool
		patterns      = make(map[string]int)
	)

	// Parsing args
	flag.IntVar(&minlen, "l", 7, "min pattern length")
	flag.IntVar(&views, "o", 3, "min occurences")
	flag.BoolVar(&patternsonly, "P", false, "only print found patterns")
	flag.StringVar(&input, "i", "", "input [default:sdtin]")

	flag.Parse()

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
	var output = make([]string, 0, 0)

	for s.Scan() {
		line := s.Text()
		output = append(output, line)
		split(minlen, line, patterns)
	}

	// trimming patterns
	buildmap(views, patterns)
	// replacing patterns map with hash ordered by length
	hash := rankByWordCount(patterns)

	// if we only want to see patterns
	if patternsonly {
		fmt.Printf(" #/len \t>Pattern by length<\n")
		for _, v := range hash {
			fmt.Printf("%2d/%2d \t>%s<\n", v.Value, len(v.Key), v.Key)
		}
		os.Exit(0)
	}

	// Attributing colors to patterns
	idxcolors := initcolors()
	keycolors := assignPatternColor(hash, idxcolors)

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
} // fin main()

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

} // fin buildmap

// split a string of minlen into a map[string]int
func split(minlen int, line string, patterns map[string]int) {
	l := len(line)
	if l < minlen {
		return
	}
	for j := 0; j <= l-minlen; j++ {
		for i := j + minlen; i <= l; i++ {
			patterns[line[j:i]]++
			//fmt.Printf("split : i=%d j=%d  %s\t -> %d\n", i, j, line[j:i], occurences[line[j:i]])
		}
	}
} // fin split

// initcolors : initialize an array with a few colors
func initcolors() []*color.Color {
	var h []*color.Color
	h = append(h, color.New(color.FgWhite, color.BgGreen))
	h = append(h, color.New(color.FgBlack, color.BgYellow))
	h = append(h, color.New(color.FgWhite, color.BgBlue))
	h = append(h, color.New(color.FgBlack, color.BgCyan))
	h = append(h, color.New(color.FgWhite, color.BgRed))
	h = append(h, color.New(color.FgWhite, color.BgMagenta))
	h = append(h, color.New(color.FgBlack, color.BgWhite))
	return h
}

// assignPatternColor : try to associate != colors for * patternts
func assignPatternColor(hash PairList, idxcolor []*color.Color) map[string]*color.Color {
	m := make(map[string]*color.Color)
	for i, p := range hash {
		m[p.Key] = idxcolor[i%len(idxcolor)]
	}
	return m
}

// end of main.go
