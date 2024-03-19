package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

type Student struct {
	Person
	School string
}

func main() {

	stu := Student{Person: Person{
		Name: "xiaoming",
		Age:  12,
	}, School: "xxxxx soll"}
	fmt.Println(stu.Name, stu.Age, stu.School)
}
