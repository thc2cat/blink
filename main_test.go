package main

import (
	"log"
	"testing"
)

//go test
func TestSplit(t *testing.T) {
	//func split(minlen int, line string, patterns map[string]int) {
	testpat := make(map[string]int)
	split(3, "1234", testpat)
	patterns := []string{"1234", "123", "234"}

	for _, key := range patterns {
		if v, ok := testpat[key]; v != 1 || !ok {
			log.Fatal("expected value 1234 found")
		}
	}
}

func TestBuildmap(t *testing.T) {
	testmap := map[string]int{
		"1234":    1, // should disapear <2
		"123":     2, // should disapear included in 12345
		"12345":   2, // should be kept =! 0123456 (2<3)
		"0123456": 3, // should kept
	}
	buildmap(2, testmap) // Keep at least 2 occurrences
	//func buildmap(views int, values map[string]int) {

	patterns := map[string]bool{
		"12345":   true,
		"0123456": true,
		"1234":    false,
	}
	for key := range patterns {
		if _, ok := testmap[key]; ok != patterns[key] {
			log.Fatal("expected value 1234 found")
		}
	}

}

// go test -bench=.
func BenchmarkBuildmap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		testpat := make(map[string]int)
		split(3, "1234567890", testpat)
	}
}

func BenchmarkSplit(b *testing.B) {
	for n := 0; n < b.N; n++ {
		testmap := map[string]int{
			"1234":    1, // should disapear <2
			"123":     2, // should disapear included in 12345
			"12345":   2, // should be kept =! 0123456 (2<3)
			"0123456": 3, // should kept
		}
		buildmap(2, testmap)
	}
}

//go test -cover
