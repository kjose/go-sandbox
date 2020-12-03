package main

import "fmt"

var x int

type car struct {
	brand        string
	hiddenNumber int
}

type person struct {
	name string
	age  int
	car  car
}

func (p person) speak() {
	fmt.Println(p.name, `says "Hello I am`, p.age, "and I have a car", p.car.brand, p.car.hiddenNumber, `"`)
}

func main() {
	y := 10
	fmt.Println(x)
	fmt.Println(y)

	s := []string{"haha", "hihi"}
	fmt.Println(s)

	m := map[string]int{
		"haha": 10,
		"hihi": 20,
	}
	fmt.Println(m)

	c := car{
		"Opel",
		123,
	}

	p := person{}
	p.name = "Kévin"
	p.age = 20
	p.car = c
	p.speak()
}
