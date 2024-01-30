package main

import "fmt"

type Vertex struct {
	Lat, Long float64
}

type Student struct {
	Name string
	Age  int
}

var m = map[string]Vertex{
	"Bell Labs": Vertex{40.68433, -74.39967},
	"Google":    Vertex{37.42202, -122.08408},
}

var studentsMap = map[int64]Student{
	220002: {"Bexeit", 19},
	220003: {"Bexeit2", 20},
	220004: {"Bexeit3", 21},
	220005: {"Bexeit4", 22},
}

func main() {
	fmt.Println(m)
	fmt.Println(len(m))

	fmt.Println(studentsMap[220003])

	k, ok := studentsMap[220006]
	fmt.Println(k)
	fmt.Println(ok)
}
