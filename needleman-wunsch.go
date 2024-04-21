package main

import (
	"4d63.com/strrev" // function for string reversing
	"fmt"
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

type StepState struct {
	Score      int
	Comparison int // one of Diag, Up, Left
}

type Alignment struct {
	Sequence1  string
	Sequence2  string
	Comparison string
}

func NewScore() Score {
	return Score{1, -1, -2, -2}
}

func main() {
	fmt.Println("Executing needleman-wunsch")

	longSeqA := "GATTTTATAAGAACCCACATTGCGTCGATTCATAAGATGTCTCGACACAGCTAATAGTTGCCCACACAACA"
	longSeqB := "GATTCTATAAGAATGCACATTGCGTCGATTCATAAGATGTCTCGACACAGCTAATAGTTGACACAACA"
	//longSeqA := "GGATAGTTGCCCACACAACA"
	//longSeqB := "TAGTTGACACAACA"
	seqA := "GATTACATT"
	seqB := "GAGCATT"
	score := NewScore()
	align(seqA, seqB, score)

	align(seqA, seqA, score)

	align(longSeqA, longSeqB, score)
}

func align(seqA string, seqB string, score Score) int {

	m := len(seqA)
	n := len(seqB)
	printer := NewPrinter()

	matrix := initializeMatrix(m, n, score)
	printer.printMatrix(matrix)

	iterateMatrix(matrix, seqA, seqB, score)
	printer.printMatrix(matrix)

	alignment := walkPath(matrix, m, n, seqA, seqB, score)

	fmt.Print("\n", seqA, " ", seqB)
	printer.printAlignment(alignment)

	printer.printScore(matrix)
	return matrix[m-1][n-1].Score
}

/*
Create an m x n matrix and initialise the first row and first column with gap values.
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

func iterateMatrix(matrix [][]StepState, seqA string, seqB string, score Score) {
	for i := 1; i < len(matrix); i++ {
		for j := 1; j < len(matrix[i]); j++ {
			matrix[i][j] = evaluateCell(matrix, i, j, seqA[i-1], seqB[j-1], score)
		}
	}
}

func evaluateCell(matrix [][]StepState, i int, j int, a uint8, b uint8, score Score) StepState {

	two_by_two_matrix := getCell(matrix, i, j)
	up := two_by_two_matrix[0][1] + score.openGap
	left := two_by_two_matrix[1][0] + score.openGap
	var diagScore int
	if a == b {
		diagScore = score.match
	} else {
		diagScore = score.mismatch
	}

	diag := two_by_two_matrix[0][0] + diagScore
	// TODO implement affine gap score
	return calcMax(up, left, diag)
}

func getCell(matrix [][]StepState, i, j int) [][]int {
	return [][]int{
		{matrix[i-1][j-1].Score, matrix[i-1][j].Score},
		{matrix[i][j-1].Score, matrix[i][j].Score},
	}
}

func walkPath(matrix [][]StepState, m int, n int, a string, b string, score Score) Alignment {
	result := make([]string, 3)
	for i := 0; i < 3; i++ {
		result[i] = ""
	}

	j := n
	i := m
	for i > 0 || j > 0 {
		fmt.Print(j, " ", i, " -- ")
		var state StepState
		state = matrix[i][j]

		//fmt.Printf("Step state %d - %d\n", state.Score, state.Comparison)
		if state.Comparison == Up {
			result[0] += "-"
			result[1] += "-"
			result[2] += string(b[i-1])
			i--
		} else if state.Comparison == Left {
			result[2] += "-"
			result[1] += "-"
			result[0] += string(a[j-1])
			j--
		} else if state.Comparison == Diag {
			result[0] += string(a[i-1])
			result[2] += string(b[j-1])
			i--
			j--

			// match or mismatch
			if a[i] == b[j] {
				result[1] += "|"
			} else {
				result[1] += "."
			}

		}
	}

	vals := Alignment{strrev.Reverse(result[0]), strrev.Reverse(result[2]), strrev.Reverse(result[1])}
	return vals
}

func calcMax(up, left, diag int) StepState {
	max := up // Assume a is the maximum initially
	comparison := Up
	// Compare b and c with the current max
	if left > max {
		max = left
		comparison = Left
	}
	if diag >= max {
		max = diag
		comparison = Diag
	}
	return StepState{max, comparison}
}
