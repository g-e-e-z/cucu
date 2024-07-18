package utils

import (
	"errors"
	"strings"
)

func RenderComponent(rows [][]string) (string, error) {
	if len(rows) == 0 {
		return "", nil
	}
	if !displayArraysAligned(rows) {
		return "", errors.New("Each item must return the same number of strings to display")
	}

	// columnPadWidths := getPadWidths(rows)
	// paddedDisplayRows := getPaddedDisplayStrings(rows, columnPadWidths)
	//
	// return strings.Join(paddedDisplayRows, "\n"), nil
    // [
    //  [GET, RequestName1],
    //  [GET, RequestName2]
    // ]
    displayRows := make([]string, len(rows))
    for i := range rows {
        displayRows = append(displayRows, strings.Join(rows[i], " | "))
    }
	return strings.Join(displayRows, "\n"), nil
}

// displayArraysAligned returns true if every string array returned from our
// list of displayables has the same length
func displayArraysAligned(stringArrays [][]string) bool {
	for _, strings := range stringArrays {
		if len(strings) != len(stringArrays[0]) {
			return false
		}
	}
	return true
}

// func getPadWidths(rows [][]string) []int {
// 	if len(rows[0]) <= 1 {
// 		return []int{}
// 	}
// 	columnPadWidths := make([]int, len(rows[0])-1)
// 	for i := range columnPadWidths {
// 		for _, cells := range rows {
// 			uncoloredCell := Decolorise(cells[i])
//
// 			if runewidth.StringWidth(uncoloredCell) > columnPadWidths[i] {
// 				columnPadWidths[i] = runewidth.StringWidth(uncoloredCell)
// 			}
// 		}
// 	}
// 	return columnPadWidths
// }
//
// func getPaddedDisplayStrings(rows [][]string, columnPadWidths []int) []string {
// 	paddedDisplayRows := make([]string, len(rows))
// 	for i, cells := range rows {
// 		for j, columnPadWidth := range columnPadWidths {
// 			paddedDisplayRows[i] += WithPadding(cells[j], columnPadWidth) + " "
// 		}
// 		paddedDisplayRows[i] += cells[len(columnPadWidths)]
// 	}
// 	return paddedDisplayRows
// }
//
