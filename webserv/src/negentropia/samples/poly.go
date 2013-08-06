package main

import "fmt"

type Talker interface {
	Talk(words string)
}

type Cat struct {
	name string
}

type Dog struct {
	name string
}

func (c *Cat) Talk(words string) {
	fmt.Printf("Cat " + c.name + " here: " + words + "\n")
}

func (d *Dog) Talk(words string) {
	fmt.Printf("Dog " + d.name + " here: " + words + "\n")
}

func main() {
	var t1, t2 Talker

	t1 = &Cat{"Kit"}
	t2 = &Dog{"Doug"}

	t1.Talk("meow")
	t2.Talk("woof")
}
