package main

import (
	"4d63.com/strrev" // function for string reversing
	"fmt"
	"log"
	"time"
)

type Score struct {
	match       int
	mismatch    int
	openGap     int
	continueGap int
}

const (
	Diag = iota
	Up
	Left
)

const (
	MatchChar    = "|"
	MismatchChar = "."
	GapChar      = "-"
)

type StepState struct {
	Score      int
	Comparison int // one of Diag, Up, Left
}

type Alignment struct {
	Sequence1  string
	Sequence2  string
	Comparison string
	Score      int
	Path       []StepState
}

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
Implementation of Needleman-Wunsch

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
	printer := NewPrinter()

	matrix := initializeMatrix(m, n, score)
	printer.printMatrix(matrix)

	iterateMatrix(matrix, seqA, seqB, score)
	printer.printMatrix(matrix)

	alignment := walkPath(matrix, m, n, seqA, seqB)

	fmt.Print("\n", seqA, " ", seqB)
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
func initializeMatrix(rows int, cols int, score Score) [][]StepState {
	matrix := make([][]StepState, rows+1)

	gapVal := score.continueGap
	for i := range matrix {
		// initialise matrix of size cols
		matrix[i] = make([]StepState, cols+1)
		// initialise first column with continue gap value
		matrix[i][0] = StepState{gapVal, Up}
		gapVal += score.continueGap
	}

	// initialise first row with continue gap value
	gapVal = score.continueGap
	for i := range matrix[0] {
		matrix[0][i] = StepState{gapVal, Left}
		gapVal += score.continueGap
	}

	// special values
	matrix[0][0] = StepState{0, Diag}
	matrix[0][1] = StepState{score.openGap, Left}
	matrix[1][0] = StepState{score.openGap, Up}

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
func iterateMatrix(matrix [][]StepState, seqA string, seqB string, score Score) {
	for i := 1; i < len(matrix); i++ {
		for j := 1; j < len(matrix[i]); j++ {
			matrix[i][j] = evaluateCell(matrix, i, j, seqA[i-1], seqB[j-1], score)
		}
	}
}

func evaluateCell(matrix [][]StepState, i int, j int, a uint8, b uint8, score Score) StepState {

	// affine gap score calculation
	upStateStep := matrix[i-1][j]
	upGapScore := score.openGap
	if upStateStep.Comparison == Up { // if the up cell is up, then this is a continuing gap
		upGapScore = score.continueGap
	}

	leftStateStep := matrix[i][j-1]
	leftGapScore := score.openGap
	if leftStateStep.Comparison == Left { // if the left cell is left, then this is an opening gap
		leftGapScore = score.continueGap
	}

	two_by_two_matrix := getCell(matrix, i, j)
	up := two_by_two_matrix[0][1] + upGapScore
	left := two_by_two_matrix[1][0] + leftGapScore
	var diagScore int
	if a == b {
		diagScore = score.match
	} else {
		diagScore = score.mismatch
	}

	diag := two_by_two_matrix[0][0] + diagScore

	return calcMax(up, left, diag)
}

func getCell(matrix [][]StepState, i, j int) [][]int {
	return [][]int{
		{matrix[i-1][j-1].Score, matrix[i-1][j].Score},
		{matrix[i][j-1].Score, matrix[i][j].Score},
	}
}

func calcMax(up, left, diag int) StepState {
	maximum := up // Assume up is the maximum initially
	comparison := Up
	// Compare left and diag with the current max
	if left > maximum {
		maximum = left
		comparison = Left
	}
	if diag >= maximum { // if diagonal is equal max, prefer it
		maximum = diag
		comparison = Diag
	}
	return StepState{maximum, comparison}
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
func walkPath(matrix [][]StepState, m int, n int, a string, b string) Alignment {
	result := make([]string, 3)
	for i := 0; i < 3; i++ {
		result[i] = ""
	}
	var path []StepState

	j := n
	i := m
	for i > 0 || j > 0 {
		fmt.Print(j, " ", i, " -- ")
		state := matrix[i][j]
		path = append(path, state)

		//fmt.Printf("Step state %d - %d\n", state.Score, state.Comparison)
		if state.Comparison == Up {
			result[0] += GapChar // gap character
			result[1] += GapChar
			result[2] += string(a[i-1])
			i--
		} else if state.Comparison == Left {
			result[2] += GapChar
			result[1] += GapChar
			result[0] += string(b[j-1])
			j--
		} else if state.Comparison == Diag {
			result[0] += string(a[i-1])
			result[2] += string(b[j-1])
			i--
			j--

			// match or mismatch
			if a[i] == b[j] {
				result[1] += MatchChar
			} else {
				result[1] += MismatchChar
			}

		}
	}

	vals := Alignment{strrev.Reverse(result[0]), strrev.Reverse(result[2]), strrev.Reverse(result[1]), matrix[m][n].Score, path}
	return vals
}
