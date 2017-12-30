package main

// from stdin, find most repetitive patterns and enlight them
// v0.2 working on the first part of the string and only the longuest pattern.
// ( wich of course will not work on very long lines )
// of course with a lot of bugs.. so more testing needed.

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

var (
	minlen, views int
	patternsonly  bool
	input         string
	patterns      = make(map[string]int)
)

func main() {

	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	info := color.New(color.FgWhite, color.BgGreen).SprintFunc()
	cyanunder := color.New(color.FgCyan).Add(color.Underline).SprintFunc()
	redbold := color.New(color.FgRed).Add(color.Bold).SprintFunc()
	fmt.Fprintf(color.Output,
		"This is a %s and this is %s %s %s %s.\n",
		yellow("warning"), red("error"), info("info"),
		redbold("redbold"), cyanunder("underlined ?"))

	flag.IntVar(&minlen, "l", 7, "min pattern length")
	flag.IntVar(&views, "o", 3, "min occurences")
	flag.BoolVar(&patternsonly, "P", false, "only print found patterns")
	flag.StringVar(&input, "i", "", "input [default:sdtin]")

	flag.Parse()

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

	var output = make([]string, 0, 0)

	// Reading input
	for s.Scan() {
		line := s.Text()
		output = append(output, line)
		split(line)
	}

	buildmap(patterns)

	if patternsonly {
		hash := rankByWordCount(patterns)
		fmt.Printf(" #/len \t>Pattern by length<\n")
		for _, v := range hash {
			fmt.Printf("%2d/%2d \t>%s<\n", v.Value, len(v.Key), v.Key)
		}
		os.Exit(0)
	}

	for _, v := range output {
		works := v
		for len(works) > minlen {
			idx, bestPLen, bestPIdx, bestP := 0, 0, 0, ""
			// recherche de la chaine la plus longue
			for k := range patterns {
				//fmt.Printf("looking %s in %s\n", k, v)
				if idx = strings.Index(works, k); (idx != -1) && (len(k) > bestPLen) {
					bestP, bestPLen, bestPIdx = k, len(k), idx
					//fmt.Printf("found new best %s in %s at idx %d\n", k, v, idx)
				}
			}
			if bestPLen > 0 {
				fmt.Fprintf(color.Output, "%s%s",
					works[0:bestPIdx],
					info(bestP))
				works = works[bestPIdx+bestPLen:]
			} else {
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
func buildmap(values map[string]int) {
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
func split(line string) {
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
