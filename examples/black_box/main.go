package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/zergon321/cirno"
)

// allPassed detects if all the tests
// have passed.
func allPassed(results []bool) bool {
	for _, result := range results {
		if !result {
			return false
		}
	}

	return true
}

// allPassed detects if all the tests
// have failed.
func allFailed(results []bool) bool {
	for _, result := range results {
		if result {
			return false
		}
	}

	return true
}

// strToShape creates a new shape from string parameters.
func strToShape(str []string) (cirno.Shape, error) {
	switch str[0] {
	case "circle":
		x, err := strconv.Atoi(str[1])

		if err != nil {
			return nil, err
		}

		y, err := strconv.Atoi(str[2])

		if err != nil {
			return nil, err
		}

		radius, err := strconv.Atoi(str[3])

		if err != nil {
			return nil, err
		}

		circle, err := cirno.NewCircle(cirno.
			NewVector(float64(x), float64(y)), float64(radius))

		if err != nil {
			return nil, err
		}

		return circle, nil

	case "line":
		x1, err := strconv.Atoi(str[1])

		if err != nil {
			return nil, err
		}

		y1, err := strconv.Atoi(str[2])

		if err != nil {
			return nil, err
		}

		x2, err := strconv.Atoi(str[3])

		if err != nil {
			return nil, err
		}

		y2, err := strconv.Atoi(str[4])

		if err != nil {
			return nil, err
		}

		line, err := cirno.NewLine(cirno.NewVector(float64(x1), float64(y1)),
			cirno.NewVector(float64(x2), float64(y2)))

		if err != nil {
			return nil, err
		}

		return line, nil

	case "rectangle":
		x, err := strconv.Atoi(str[1])

		if err != nil {
			return nil, err
		}

		y, err := strconv.Atoi(str[2])

		if err != nil {
			return nil, err
		}

		width, err := strconv.Atoi(str[3])

		if err != nil {
			return nil, err
		}

		height, err := strconv.Atoi(str[4])

		if err != nil {
			return nil, err
		}

		angle, err := strconv.Atoi(str[5])

		if err != nil {
			return nil, err
		}

		rectangle, err := cirno.NewRectangle(cirno.
			NewVector(float64(x), float64(y)), float64(width), float64(height), float64(angle))

		if err != nil {
			return nil, err
		}

		return rectangle, nil

	default:
		return nil, fmt.Errorf("unknown shape: %s", str[0])
	}
}

func main() {
	file, err := os.Open("shapes")
	handleError(err)
	defer file.Close()

	tests := []string{}
	results := []bool{}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		// Extract the string from the file.
		str := scanner.Text()
		parts := strings.Split(str, ", ")

		// Get parameters from the string.
		shapeOneStr := strings.Split(parts[0], " ")
		shapeTwoStr := strings.Split(parts[1], " ")
		expectedStr := parts[2]

		// Parse the string.
		expected, err := strconv.ParseBool(expectedStr)
		handleError(err)
		oneShape, err := strToShape(shapeOneStr)
		handleError(err)
		otherShape, err := strToShape(shapeTwoStr)
		handleError(err)

		// Test the obtained shapes for collision
		// and compare results.
		test := parts[0] + ", " + parts[1]
		actual, err := cirno.ResolveCollision(oneShape, otherShape, false)
		handleError(err)
		result := expected == actual

		tests = append(tests, test)
		results = append(results, result)
	}

	// Output test summary.
	if allPassed(results) {
		fmt.Println("All tests passed")
	} else if allFailed(results) {
		fmt.Println("All tests failed")
	} else {
		length := len(tests)

		for i := 0; i < length; i++ {
			fmt.Print(tests[i] + " - ")

			if results[i] {
				fmt.Println("PASSED")
			} else {
				fmt.Println("FAILED")
			}
		}
	}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
