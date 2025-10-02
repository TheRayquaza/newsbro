package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <input_file> [output_file]")
		fmt.Println("If output_file is not provided, results will be printed to stdout")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	var outputFile string
	if len(os.Args) >= 3 {
		outputFile = os.Args[2]
	}

	// Read file and remove duplicates
	lines, err := readAndDedup(inputFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Write output
	if outputFile != "" {
		err = writeToFile(outputFile, lines)
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully removed duplicates. Original lines: %d, Unique lines: %d\n", 
			countLines(inputFile), len(lines))
		fmt.Printf("Output written to: %s\n", outputFile)
	} else {
		for _, line := range lines {
			fmt.Println(line)
		}
	}
}

func readAndDedup(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	seen := make(map[string]bool)
	var uniqueLines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !seen[line] {
			seen[line] = true
			uniqueLines = append(uniqueLines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return uniqueLines, nil
}

func writeToFile(filename string, lines []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

func countLines(filename string) int {
	file, err := os.Open(filename)
	if err != nil {
		return 0
	}
	defer file.Close()

	count := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		count++
	}
	return count
}