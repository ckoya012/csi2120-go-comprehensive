package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
)

// MatrixInt struct
type MatrixInt struct {
	n       int
	matrix  [][]BaseNumber
	allRows []BaseRow
	allCols []BaseCol
}

type BaseNumber struct {
	val        int
	assigned   bool
	crossedOut bool
}

type BaseRow struct {
	row  []BaseNumber
	tick bool
	line bool
}

type BaseCol struct {
	col  []BaseNumber
	tick bool
	line bool
}

func findMinimum(arr []BaseNumber) int {
	min := arr[0].val
	for _, val := range arr {
		if val.val < min {
			min = val.val
		}
	}
	return min
}

// STEP 1: Row Reduction
func step1RowReduction(row []BaseNumber, wg *sync.WaitGroup) {
	defer wg.Done()

	minValue := findMinimum(row)
	for i := 0; i < len(row); i++ {
		row[i].val -= minValue
	}
}

// STEP 2: Column Reduction
func (m MatrixInt) step2ColumnReduction(min int, index int, wg *sync.WaitGroup) {
	defer wg.Done()

	for row := 0; row < len(m.matrix); row++ {
		m.matrix[row][index].val -= min
	}
}

// scan row L to R and keep the 1st 0
// cross out all other zeroes
// go up and down the column & loop over row and cross out any othe zero

// 3b if there's more than 0
// row is ticked, if A(r,C) is assigned sicne each row we find an assignment
// then untick again so theres nothing ticked at teh end
// tick is to ensure that when we have 1+ zero, then we handle in next loop

// STEP 3a
func (m *MatrixInt) step3OptimalSln() {
	for row := 0; row < len(m.matrix); row++ {
		first := true
		for col := 0; col < len(m.matrix); col++ {
			if m.matrix[row][col].val == 0 {
				if first {
					first = false
					m.matrix[row][col].assigned = true

					for rr := 0; rr < len(m.matrix); rr++ {
						// fmt.Println("--------------")
						// fmt.Println("ROW: ", row)
						// fmt.Println("RR: ", rr)
						// fmt.Println("CC: ", col)
						// // fmt.Println(m.allCols)
						// fmt.Println(m.matrix[col][rr])

						if (m.matrix[rr][col].assigned != true || m.matrix[row][col].assigned != true) && m.matrix[rr][col].val == 0 && row != rr {
							m.matrix[rr][col].crossedOut = true
						}

					}
				} else {
					m.matrix[row][col].crossedOut = true
				}
			}
			fmt.Println(col, row, m.matrix[col][row])

		}
	}

	m.createBaseCols()
	// fmt.Println("-------STEP3 ROWS--------")
	// fmt.Println(m.allRows)
	// fmt.Println("-------STEP3 COLS--------")
	// fmt.Println(m.allCols)

	m.tickRows()

}

// Step 3b Ticking Rows
func (m *MatrixInt) tickRows() {
	for row := 0; row < len(m.matrix); row++ {
		m.allRows[row].tick = true
		for col := 0; col < len(m.matrix[row]); col++ {
			if m.matrix[row][col].val == 0 {
				if m.matrix[row][col].assigned == true && m.matrix[row][col].crossedOut == false {
					m.allRows[row].tick = false
				}
			}
		}
	}

	m.createBaseCols()
	m.tickCols()
}

// Step 3b
func (m *MatrixInt) tickCols() {
	for row := 0; row < len(m.matrix); row++ {
		if m.allRows[row].tick != true {
			continue
		}

		for col := 0; col < len(m.matrix[row]); col++ {
			if m.matrix[row][col].crossedOut == true {
				m.allCols[col].tick = true
			}
		}
	}
	m.createBaseCols()

	// Step 3c
	// newTick := false

	for col := 0; col < len(m.matrix); col++ {
		if m.allCols[col].tick != true {
			continue
		}

		for row := 0; row < len(m.matrix[col]); row++ {
			if m.matrix[row][col].assigned == true {
				m.allRows[row].tick = true
				// newTick = true
			}
		}
	}

	// Step 3d
	// if newTick == true {
	// 	m.tickRows()
	// }

	m.drawLines()
}

// Step 3e
func (m *MatrixInt) drawLines() {

	for row := 0; row < len(m.matrix); row++ {
		if m.allRows[row].tick != true {
			m.allRows[row].line = true
		} else {
			m.allRows[row].line = false
		}
	}

	for col := 0; col < len(m.matrix); col++ {
		if m.allCols[col].tick == true {
			m.allCols[col].line = true
		} else {
			m.allCols[col].line = false
		}
	}
	m.createBaseCols()

	// fmt.Println("-----BASE COLS----")
	// fmt.Println(m.allCols, "\n")
	// fmt.Println("-----BASE ROWS----")
	// fmt.Println(m.allRows, "\n")

}

func (m *MatrixInt) checkForStep4() {
	var counter int
	for row := 0; row < len(m.matrix); row++ {
		if m.allRows[row].line == true {
			counter++
		}
		if m.allCols[row].line == true {
			counter++
		}
	}

	if counter < len(m.allCols) || counter < len(m.allRows) {
		m.step4ShiftZeroes()
	} else {
		return
	}
}

// Find smallest uncovered value = 5
// Subtract from all uncovered values
// Add to intersection of the 2 lines
// Lines removed and go back to step 3
func (m *MatrixInt) step4ShiftZeroes() {
	// var validBaseNums []BaseNumber
	var minValue int
	for rowIndex, row := range m.allRows {
		for colIndex, col := range m.allCols {
			if row.line == false && col.line == false {
				// validBaseNums = append(validBaseNums, m.matrix[rowIndex][colIndex])

				if minValue == 0 {
					minValue = m.matrix[rowIndex][colIndex].val
				} else if m.matrix[rowIndex][colIndex].val < minValue {
					minValue = m.matrix[rowIndex][colIndex].val
				}
			}
		}
	}

	for rowIndex, row := range m.allRows {
		for colIndex, col := range m.allCols {
			if row.line == false && col.line == false {
				m.matrix[rowIndex][colIndex].val -= 5
			} else if row.line == true && col.line == true {
				m.matrix[rowIndex][colIndex].val += 5
			}
		}
	}

	for r := 0; r < len(m.allRows); r++ {
		for c := 0; r < len(m.allCols); r++ {
			m.allRows[r].line = false
			m.allRows[r].tick = false
			m.allCols[c].line = false
			m.allCols[c].tick = false
		}
	}

	m.createBaseCols()
	m.step3OptimalSln()

}

func (m *MatrixInt) step5FinalAssignment(original MatrixInt) ([][]int, int) {
	var values [][]int
	for row := 0; row < len(m.matrix); row++ {
		for col := 0; col < len(m.matrix[row]); col++ {
			if m.matrix[row][col].assigned == true && m.matrix[row][col].crossedOut == false {
				index := []int{row, col}
				values = append(values, index)
			}
		}
	}

	var result int
	for i := 0; i < len(values); i++ {
		r := values[i][0]
		c := values[i][1]
		result += original.matrix[r][c].val
	}
	return values, result
}

// Prints a matrix of strings
func printMatrix(matrix [][]BaseNumber) {
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[i]); j++ {
			fmt.Print(matrix[i][j].val)
			fmt.Print(" ")
		}
		fmt.Println()
	}
}

// Create BaseRow and BaseCol arrays
func (matrixInt *MatrixInt) createBaseRows() {
	for row := 0; row < len(matrixInt.matrix); row++ {
		var newBaseRow BaseRow
		newBaseRow.row = matrixInt.matrix[row]
		matrixInt.allRows = append(matrixInt.allRows, newBaseRow)
	}
}

func (matrixInt *MatrixInt) createBaseCols() {

	for index, _ := range matrixInt.allRows {
		if matrixInt.allCols != nil {
			matrixInt.allCols[index].col = nil

			columnValues := []BaseNumber{}
			for col := 0; col < len(matrixInt.matrix); col++ {
				columnValues = append(columnValues, matrixInt.matrix[col][index])
			}
			var newBaseCol BaseCol
			newBaseCol.col = columnValues
			newBaseCol.tick = matrixInt.allCols[index].tick
			newBaseCol.line = matrixInt.allCols[index].line

			matrixInt.allCols[index] = newBaseCol
		} else {
			for row := 0; row < len(matrixInt.matrix); row++ {
				columnValues := []BaseNumber{}
				for col := 0; col < len(matrixInt.matrix); col++ {
					columnValues = append(columnValues, matrixInt.matrix[col][row])
				}
				var newBaseCol BaseCol
				newBaseCol.col = columnValues

				matrixInt.allCols = append(matrixInt.allCols, newBaseCol)

			}
		}
	}
	// matrixInt.allCols = nil

	// for row := 0; row < len(matrixInt.matrix); row++ {
	// 	columnValues := []BaseNumber{}
	// 	for col := 0; col < len(matrixInt.matrix); col++ {
	// 		columnValues = append(columnValues, matrixInt.matrix[col][row])
	// 	}
	// 	var newBaseCol BaseCol
	// 	newBaseCol.col = columnValues

	// 	matrixInt.allCols = append(matrixInt.allCols, newBaseCol)
	// 	// newBaseCol.line = matrixInt.allCols[row].line
	// 	// newBaseCol.tick = matrixInt.allCols[row].tick
	// }

}

func main() {

	var wg sync.WaitGroup

	// Open the file
	csvfile, err := os.Open("cost_matrix_5.csv")
	if err != nil {
		log.Fatalln("Error opening the CSV file", err)
	}
	// Parse the file
	r := csv.NewReader(csvfile)

	var matrixInt MatrixInt
	var matrixStr [][]string
	var originalMatrix MatrixInt

	// Iterate through the records and read each record
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		matrixStr = append(matrixStr, record)

		var r1 = []BaseNumber{}
		var r2 = []BaseNumber{}
		for _, i := range record {
			j, err := strconv.Atoi(i)
			if err == nil {
				var b BaseNumber
				b.val = j
				r1 = append(r1, b)
				r2 = append(r2, b)
			}
		}

		originalMatrix.matrix = append(originalMatrix.matrix, r1)
		matrixInt.matrix = append(matrixInt.matrix, r2)
		matrixInt.n = len(matrixInt.matrix) - 1
	}

	// Removing empty array at index 0
	matrixInt.matrix = matrixInt.matrix[1:]
	originalMatrix.matrix = originalMatrix.matrix[1:]
	matrixInt.createBaseRows()

	fmt.Println("************************\nMatrix generated from CSV file: ")
	printMatrix(matrixInt.matrix)

	fmt.Println("************************\nStarting Step 1 ")

	for i := 0; i < len(matrixInt.matrix); i++ {
		wg.Add(1)
		go step1RowReduction(matrixInt.matrix[i], &wg)
	}

	wg.Wait()

	fmt.Println("************************\nAfter Step 1: ")
	// printMatrix(matrixInt.matrix)

	fmt.Println("************************\nStarting Step 2")
	for row := 0; row < len(matrixInt.matrix); row++ {
		columnValues := []BaseNumber{}
		for col := 0; col < len(matrixInt.matrix); col++ {
			columnValues = append(columnValues, matrixInt.matrix[col][row])
		}
		wg.Add(1)
		minValue := findMinimum(columnValues)
		go matrixInt.step2ColumnReduction(minValue, row, &wg)
	}
	wg.Wait()

	fmt.Println("************************\nAfter Step 2: ")
	// printMatrix(matrixInt.matrix)
	matrixInt.createBaseCols()

	fmt.Println("************************\nStarting Step 3")
	matrixInt.step3OptimalSln()
	fmt.Println("************************\nAfter Step 3: ")
	printMatrix(matrixInt.matrix)

	fmt.Println("************************\nStarting Step 4")
	matrixInt.checkForStep4()
	// printMatrix(matrixInt.matrix)

	fmt.Println("************************\nStarting Step 5")
	index, result := matrixInt.step5FinalAssignment(originalMatrix)
	fmt.Println("************************\nAfter Step 5: ")
	fmt.Println("The result is: ", result)
	fmt.Println("The indexes are: ", index, "\n")
	// printMatrix(originalMatrix.matrix)

	// // Create table
	// content := ""
	// for i := 1; i < len(matrixStr); i++ {
	// 	content += matrixStr[i][0]
	// 	content += fmt.Sprint(" ", index[i-1][1], "\n")
	// }

	// // Create output file
	// filename := fmt.Sprint("tracker_go_", matrixInt.n, ".csv")
	// file, err := os.Create(filename)
	// if err != nil {
	// 	panic(err)
	// }
	// defer file.Close()

	// file.WriteString(content)
	// if err != nil {
	// 	panic(err)
	// }

	// for i := 0; i < len(matrixInt.allCols); i++ {
	// 	fmt.Println(matrixInt.allRows[i])
	// }
	// fmt.Println("----------------")
	// for i := 0; i < len(matrixInt.allCols); i++ {
	// 	fmt.Println(matrixInt.allCols[i])
	// }
	// fmt.Println("----------------")
}
