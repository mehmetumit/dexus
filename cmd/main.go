package main

import "fmt"

var(
	Version string
	Commit string
)

func main(){
	fmt.Printf("Hello dexus\nVersion: %s | Commit: %s\n", Version, Commit)
}
