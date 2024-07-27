package main

import (
	"log"
	"time"
)

/*
Score datatype capturing a scoring schema for the needleman wunsch algorithm.

Gap and mismatch scores should be negative, and match scores should be non-negative
*/
type Score struct {
	match       int
	mismatch    int
	openGap     int
	continueGap int
}

type direction int

const (
	diag direction = iota
	up
	left
)

const (
	matchChar    = "|"
	mismatchChar = "."
	gapChar      = "-"
)

type stepState struct {
	Score      int
	Comparison direction // one of diag, up, left
}

/*
Alignment - result of an align operation between two sequences, calculates
the match score and path through the matrix and produces a comparison string.
*/
type Alignment struct {
	Sequence1  string
	Sequence2  string
	Comparison string
	Score      int
	Path       []stepState
}

/*
NewScore Create a default scoring scheme
*/
func NewScore() Score {
	return Score{1, -1, -2, -2}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func main() {
	log.Println("Executing needleman-wunsch")

	//longSeqA := "GGATAGTTGCCCACACAACA"
	//longSeqB := "TAGTTGACACAACA"
	seqA := "GATTACATT"
	seqB := "GAGCATT"
	score := NewScore()
	Align(seqA, seqB, score)

	Align(seqA, seqA, score)

	//align(longlongA, longlongB, score)
}

/*
Align Implementation of Needleman-Wunsch

Initialization:

	F(i,0)=F(i−1,0)−d
	F(0,j)=F(0,j−1)−d

	Iteration:F(i,j)=max{
		F(i−1,j)−d insert gap in S
		F(i,j−1)−d insert gap in T
		F(i−1,j−1)+s(xi,yj) match or mutation
*/
func Align(seqA string, seqB string, score Score) Alignment {
	defer timeTrack(time.Now(), "align")
	m := len(seqA)
	n := len(seqB)
	printer := noopPrinter()

	matrix := initializeMatrix(m, n, score)
	printer.printMatrix(matrix)

	iterateMatrix(matrix, seqA, seqB, score)
	printer.printMatrix(matrix)

	alignment := walkPath(matrix, m, n, seqA, seqB)

	printer.printSequences(seqA, seqB)

	printer.printAlignment(alignment)

	printer.printScore(matrix)
	return alignment
}

/*
Create an m + 1 x n + 1 matrix and initialise the first row and first column with gap values.
Based on :

Initialization:

	F(i,0)=F(i−1,0)−d
	F(0,j)=F(0,j−1)−d

i.e.

------------------------------
|  0 | -2 | -6 | -8 |-10 |-12 |
| -2 |  0 |  0 |  0 |  0 |  0 |
| -6 |  0 |  0 |  0 |  0 |  0 |
| -8 |  0 |  0 |  0 |  0 |  0 |
|-10 |  0 |  0 |  0 |  0 |  0 |
|-12 |  0 |  0 |  0 |  0 |  0 |
|-14 |  0 |  0 |  0 |  0 |  0 |
|-16 |  0 |  0 |  0 |  0 |  0 |
------------------------------
*/
func initializeMatrix(rows int, cols int, score Score) [][]stepState {
	matrix := make([][]stepState, rows+1)

	gapVal := score.continueGap
	for i := range matrix {
		// initialise matrix of size cols
		matrix[i] = make([]stepState, cols+1)
		// initialise first column with continue gap value
		matrix[i][0] = stepState{gapVal, up}
		gapVal += score.continueGap
	}

	// initialise first row with continue gap value
	gapVal = score.continueGap
	for i := range matrix[0] {
		matrix[0][i] = stepState{gapVal, left}
		gapVal += score.continueGap
	}

	// special values
	matrix[0][0] = stepState{0, diag}
	matrix[0][1] = stepState{score.openGap, left}
	matrix[1][0] = stepState{score.openGap, up}

	return matrix
}

/*
Iterate through the matrix and evaluate each cell based on the max of:

	Iteration:F(i,j)=max{
		F(i−1,j)−d insert gap in S
		F(i,j−1)−d insert gap in T
		F(i−1,j−1)+s(xi,yj) match or mutation

i.e.
------------------------------
|  0 | -2 | -6 | -8 |-10 |-12 |
| -2 |  1 | -1 | -3 | -5 | -7 |
| -6 | -1 |  2 |  0 | -2 | -4 |
| -8 | -3 |  0 |  1 | -1 | -3 |
|-10 | -5 | -2 | -1 |  0 | -2 |
|-12 | -7 | -4 | -1 |  0 |  1 |
|-14 | -9 | -6 | -3 | -2 | -1 |
|-16 |-11 | -8 | -5 | -2 | -1 |
------------------------------

Record the cell path, and the direction that was taken in the matrix
*/
func iterateMatrix(matrix [][]stepState, seqA string, seqB string, score Score) {
	for i := 1; i < len(matrix); i++ {
		for j := 1; j < len(matrix[i]); j++ {
			matrix[i][j] = evaluateCell(matrix, i, j, seqA[i-1], seqB[j-1], score)
		}
	}
}

func evaluateCell(matrix [][]stepState, i int, j int, a uint8, b uint8, score Score) stepState {

	// affine gap score calculation
	upStateStep := matrix[i-1][j]
	upGapScore := score.openGap
	if upStateStep.Comparison == up { // if the up cell is up, then this is a continuing gap
		upGapScore = score.continueGap
	}

	leftStateStep := matrix[i][j-1]
	leftGapScore := score.openGap
	if leftStateStep.Comparison == left { // if the left cell is left, then this is an opening gap
		leftGapScore = score.continueGap
	}

	twoByTwoMatrix := getCell(matrix, i, j)
	up := twoByTwoMatrix[0][1] + upGapScore
	left := twoByTwoMatrix[1][0] + leftGapScore
	var diagScore int
	if a == b {
		diagScore = score.match
	} else {
		diagScore = score.mismatch
	}

	diag := twoByTwoMatrix[0][0] + diagScore

	return calcMax(up, left, diag)
}

func getCell(matrix [][]stepState, i, j int) [][]int {
	return [][]int{
		{matrix[i-1][j-1].Score, matrix[i-1][j].Score},
		{matrix[i][j-1].Score, matrix[i][j].Score},
	}
}

func calcMax(upScore, leftScore, diagScore int) stepState {
	maximum := upScore // Assume up is the maximum initially
	comparison := up
	// Compare left and diag with the current max
	if leftScore > maximum {
		maximum = leftScore
		comparison = left
	}
	if diagScore >= maximum { // if diagonal is equal max, prefer it
		maximum = diagScore
		comparison = diag
	}
	return stepState{maximum, comparison}
}

/*
Follow the path from the bottom right back to the top left. This is a matter of following the directional
state elements we calculated before.

This builds up the sequence alignments and comparison string in reverse.
i.e. For two sequences GATTACA GAAAA

5 7 -- 4 6 -- 3 5 -- 2 4 -- 2 3 -- 2 2 -- 1 1 --

GA--ACA
||--|.|
GATTAAA
Score -1

returns an alignment object, with the 3 strings, the max score and the directional path from bottom right to top left
*/
func walkPath(matrix [][]stepState, m int, n int, a string, b string) Alignment {
	result := make([]string, 3)
	for i := 0; i < 3; i++ {
		result[i] = ""
	}
	var path []stepState

	j := n
	i := m
	for i > 0 || j > 0 {
		//fmt.Print(j, " ", i, " -- ")
		state := matrix[i][j]
		path = append(path, state)

		//fmt.Printf("Step state %d - %d\n", state.Score, state.Comparison)
		if state.Comparison == up {
			result[0] += gapChar // gap character
			result[1] += gapChar
			result[2] += string(a[i-1])
			i--
		} else if state.Comparison == left {
			result[2] += gapChar
			result[1] += gapChar
			result[0] += string(b[j-1])
			j--
		} else if state.Comparison == diag {
			result[0] += string(a[i-1])
			result[2] += string(b[j-1])
			i--
			j--

			// match or mismatch
			if a[i] == b[j] {
				result[1] += matchChar
			} else {
				result[1] += mismatchChar
			}

		}
	}

	vals := Alignment{reverseString(result[0]), reverseString(result[2]), reverseString(result[1]), matrix[m][n].Score, path}
	return vals
}

func reverseString(s string) string {
	// Convert the string to a slice of runes to handle multi-byte characters correctly
	runes := []rune(s)
	// Get the length of the rune slice
	n := len(runes)
	// Reverse the rune slice in place
	for i := 0; i < n/2; i++ {
		runes[i], runes[n-1-i] = runes[n-1-i], runes[i]
	}
	// Convert the reversed rune slice back to a string
	return string(runes)
}
