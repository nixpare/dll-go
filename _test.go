//dll Print(s string) (n int) = ./print.dll
package main

func main() {
	fmt.Println(Print("Hello, World!"))
	fmt.Println(PrintDLL("Hello, World! (from dll)"))
}

func Print(s string) int {
	n, _ := fmt.Println(s)
	return n
}
