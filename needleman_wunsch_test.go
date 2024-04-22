package main

import (
	"fmt"
	"testing"
)

var defaultScore = NewScore()

func TestAlignExactMatch(t *testing.T) {
	seqA := "GATTACATT"
	alignment := Align(seqA, seqA, defaultScore)
	doAlignTest(t, alignment, 9, "|||||||||")
}

func TestAlignMismatch(t *testing.T) {
	seqA := "GATTACA"
	seqB := "GATTAAA"
	alignment := Align(seqA, seqB, defaultScore)
	doAlignTest(t, alignment, 5, "|||||.|")
}

func TestAlignGap(t *testing.T) {
	seqA := "GATTACA"
	seqB := "GAACA"
	alignment := Align(seqA, seqB, defaultScore)
	doAlignTest(t, alignment, 1, "||--|||")
}

func TestAlignMismatchAndGap(t *testing.T) {
	seqA := "GATTACA"
	seqB := "GAAAA"
	alignment := Align(seqA, seqB, defaultScore)
	doAlignTest(t, alignment, -1, "||--|.|")
}

func TestAffineGapScore(t *testing.T) {
	seqA := "GATTTTATAAGAACCCACATTGCGTCGATTCATAAGATGTCTCGACACAGCTAATAGTTGCCCACACAACAGGG"
	seqB := "GATTCTATAAGAATGCACATTGCGTCGATTCATAAGATGTCTCGACACAGCTAATAGTTGACACAA"
	comp := "||||.||||||||..|||||||||||||||||||||||||||||||||||||||||||||---||||||-----"

	alignment := Align(seqA, seqB, Score{1, -1, -5, -2})
	doAlignTest(t, alignment, 38, comp)
}

func doAlignTest(t *testing.T, alignment Alignment, score int, comparison string) {

	t.Logf("%s : %s\n", t.Name(), alignment.Comparison)
	if len(alignment.Comparison) != len(comparison) {
		t.Error(fmt.Printf("Wanted size %d, got %d\n", len(comparison), len(alignment.Comparison)))
	}

	if alignment.Comparison != comparison {
		t.Error(fmt.Printf("Wanted comparison string %s, got %s\n", comparison, alignment.Comparison))
	}

	if alignment.Score != score {
		t.Error(fmt.Printf("Wanted score %d, got %d\n", score, alignment.Score))
	}
}
