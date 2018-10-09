package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type TestCase []int64

func sumOfSquares(tc TestCase, accumulator uint64) uint64 {
	if len(tc) == 0 {
		return accumulator
	}

	var square uint64 = 0
	if tc[0] > 0 {
		square = uint64(tc[0]) * uint64(tc[0])
	}

	return sumOfSquares(tc[1:], accumulator+square)
}

func extractArray(numbersStr []string, testCase TestCase) TestCase {
	if len(numbersStr) == 0 {
		return testCase
	}

	number, err := strconv.ParseInt(numbersStr[0], 10, 64)
	if err == nil {
		testCase = append(testCase, number)
	}

	return extractArray(numbersStr[1:], testCase)
}

func extractTestCases(lines []string, testCases []TestCase) []TestCase {
	if len(lines) == 0 {
		return testCases
	}

	arrLen, _ := strconv.Atoi(lines[0])
	numberStrings := strings.Split(strings.TrimSpace(lines[1]), " ")

	if arrLen == len(numberStrings) {
		tc := extractArray(numberStrings, make(TestCase, 0, arrLen))
		testCases = append(testCases, tc)
	}

	return extractTestCases(lines[2:], testCases)
}

func runTestCases(testCases []TestCase) {
	if len(testCases) == 0 {
		return
	}

	fmt.Println(sumOfSquares(testCases[0], 0))
	runTestCases(testCases[1:])
}

func panicHandler() {
	if err := recover(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	defer panicHandler()

	reader := bufio.NewReader(os.Stdin)
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}

	input := string(bytes)
	input = strings.TrimSpace(input)
	lines := strings.Split(input, "\n")
	expectedNumTestCases, _ := strconv.Atoi(lines[0])
	lines = lines[1:]

	numTestCases := len(lines) / 2

	if expectedNumTestCases != numTestCases {
		panic(fmt.Sprintf("bad input, expected: %d, got: %d", numTestCases, expectedNumTestCases))
	}

	testCases := extractTestCases(lines, make([]TestCase, 0, numTestCases))
	runTestCases(testCases)
}
