package main

import "fmt"

func main() {
	anon := make(map[string]struct{})
	anon["an"] = struct{}{}

	if _, ok := anon["an"]; ok {
		fmt.Println("found")
	}
	if anon["ban"] == struct{}{} {
		fmt.Println("second ")
	}
}
