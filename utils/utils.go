package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

    displayRows := make([]string, len(rows))
    for i := range rows {
        displayRows = append(displayRows, strings.Join(rows[i], " | "))
    }
	return strings.TrimSpace(strings.Join(displayRows, "\n")), nil
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

func ValuesToMap(params [][2]string) map[string]string {
	result := make(map[string]string)
	for _, pair := range params {
        key := pair[0]
        val := pair[1]
        result[key] = val
	}
	return result
}

func MapToSlice(m map[string]string) [][]string {
	result := [][]string{}
	for key, value := range m {
		result = append(result, []string{key, value})
	}
	return result
}

// NormalizeLinefeeds - Removes all Windows and Mac style line feeds
func NormalizeLinefeeds(str string) string {
	str = strings.Replace(str, "\r\n", "\n", -1)
	str = strings.Replace(str, "\r", "", -1)
	return str
}

// Max returns the maximum of two integers
func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// Parse function that takes a URL string and returns a slice of key-value pairs
func Parse(rawURL string) ([][]string, error) {
	// Split the URL into the base URL and the query string
	parts := strings.SplitN(rawURL, "?", 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("no query parameters found")
	}

	// Split the query string into key-value pairs
	query := parts[1]
	pairs := strings.Split(query, "&")

	// Create a slice to store the key-value pairs
	var result [][]string

	// Iterate over each pair and split into key and value
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid query parameter: %s", pair)
		}

		// Decode the key and value
		key := strings.ReplaceAll(kv[0], "%20", " ")
		value := strings.ReplaceAll(kv[1], "%20", " ")

		// Append the key-value tuple to the result slice
		result = append(result, []string{key, value})
	}

	return result, nil
}

// ToJSON function converts a map[string]interface{} to a formatted JSON string
func ToJSON(data map[string]interface{}) (string, error) {
	// Marshal the map into a JSON byte slice with indentation
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}

	// Convert the JSON byte slice to a string and return it
	return string(jsonData), nil
}

// function to format JSON data
func FormatJSON(data []byte) string {
	var out bytes.Buffer
	err := json.Indent(&out, data, "", " ")

	if err != nil {
		fmt.Println(err)
	}

	d := out.Bytes()
	return string(d)
}
