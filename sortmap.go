package main

import "sort"

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
