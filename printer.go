package main

import "fmt"

/*
Printer - interface for printing common objects when running needleman_wunsch
*/
type Printer interface {
	printAlignment(alignment Alignment)

	printMatrix(matrix [][]stepState)

	printScore(matrix [][]stepState)

	printSequences(seqA, seqB string)
}

type simplePrinter struct {
}

func (p simplePrinter) printMatrix(matrix [][]stepState) {

	if len(matrix) > 20 {
		fmt.Println("Matrix too large, skipping printing matrix")
		return
	}
	fmt.Print("------------------------------\n")
	for _, row := range matrix {
		for _, col := range row {
			fmt.Printf("|%3d ", col.Score)
		}
		fmt.Print("|\n")
	}
	fmt.Print("------------------------------\n")

}

func (p simplePrinter) printAlignment(alignment Alignment) {
	fmt.Printf("\n%s\n%s\n%s\n", alignment.Sequence1, alignment.Comparison, alignment.Sequence2)
}

func (p simplePrinter) printScore(matrix [][]stepState) {
	fmt.Printf("Max score : %d\n", matrix[len(matrix)-1][len(matrix[0])-1])
}

func (p simplePrinter) printSequences(seqA, seqB string) {
	fmt.Printf("\nSequence 1: %s\nSequence 2: %s\n", seqA, seqB)
}

func newPrinter() Printer {
	return simplePrinter{}
}

func noopPrinter() Printer {
	return nullPrinter{}
}

type nullPrinter struct {
}

func (p nullPrinter) printAlignment(alignment Alignment) {

}

func (p nullPrinter) printMatrix(matrix [][]stepState) {

}

func (p nullPrinter) printScore(matrix [][]stepState) {

}

func (p nullPrinter) printSequences(seqA, seqB string) {

}
