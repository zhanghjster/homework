package main

import "fmt"

type Triangle struct {
	A int64
	B int64
	C int64
	Base int64
	Height int64
}

func (t Triangle) Area() int64 {
	return (t.Base * t.Height)/2
}

type Rectangle struct {
	Width int64
	Height int64
}

func (r Rectangle) Area() int64 {
	return r.Height*r.Width
}

func main() {

	t := Triangle{
		Base: 4, Height:3,
	}

	r := Rectangle{
		Width:4, Height: 3,
	}

	SaveShapeArea(t)
	SaveShapeArea(r)
}

type Shape interface {
	Area() int64
}

func SaveShapeArea(s Shape) {
	fmt.Println(s.Area())
}