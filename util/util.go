package util

import (
	"bufio"
	"io"
	"os"
	"strconv"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func ReadLines(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	var result []string
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	return result, scanner.Err()
}

func ReadLinesFromFile(name string) ([]string, error) {
	var err error
	var reader io.Reader
	if reader, err = os.Open(name); err != nil {
		return nil, err
	}
	var result []string
	if result, err = ReadLines(reader); err != nil {
		return nil, err
	}
	return result, nil
}

func ReadInts(r io.Reader) ([]int, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	var result []int
	for scanner.Scan() {
		x, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return result, err
		}
		result = append(result, x)
	}
	return result, scanner.Err()
}

func ReadIntsFromFile(name string) ([]int, error) {
	var err error
	var reader io.Reader
	if reader, err = os.Open(name); err != nil {
		return nil, err
	}
	var result []int
	if result, err = ReadInts(reader); err != nil {
		return nil, err
	}
	return result, nil
}

func MinInt(first int, rest ...int) int {
	result := first
	for _, i := range rest {
		if i < result {
			result = i
		}
	}
	return result
}

func MaxInt(first int, rest ...int) int {
	result := first
	for _, i := range rest {
		if i > result {
			result = i
		}
	}
	return result
}

func AbsInt(i int) int {
	if i < 0 {
		return -i
	} else {
		return i
	}
}
