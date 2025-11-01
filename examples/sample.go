package main

import "fmt"

// Sample Go code for testing the MCP server
// This file contains intentional issues for demonstration

func main() {
	// Issue: using println instead of fmt.Println
	println("Hello, World!")

	// Issue: unused variable
	x := 42

	// Issue: error not checked
	fmt.Println("This is a sample")

	// Good code
	result := add(1, 2)
	fmt.Printf("Result: %d\n", result)
}

// add adds two numbers
func add(a, b int) int {
	return a + b
}

// Issue: unused function
func unused() {
	fmt.Println("This function is never called")
}
