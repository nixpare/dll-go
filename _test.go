//go:generate dll-go -output _test-dll.go _test.go
package main

func main() {
	fmt.Println(Print("Hello, World!" ))
	fmt.Println(PrintDLL("Hello, World! (from dll)"))
}

//dll Print(s string) (n int, echo *string) = ./print.dll
func Print(msg message) (int, error) {
	n, _ := fmt.Println(msg.m)
	return n, nil
}
