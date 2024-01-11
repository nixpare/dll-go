//go:generate dll-go -output _test-dll.go _test.go
package main

func main() {
	fmt.Println(Print("Hello, World!"))
	fmt.Println(PrintDLL("Hello, World! (from dll)"))
}

//dll Print(s string) (n int) = ./print.dll
func Print(s string) int {
	n, _ := fmt.Println(s)
	return n
}
