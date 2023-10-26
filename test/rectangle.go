package main

import (
	"fmt"
	"log"
)

type Rectangle struct {
	width  float32
	height float32
}

func (r *Rectangle) Area() {
	r.width = 5
	log.Println("正方形修改中：", r)
	fmt.Println("面积=", r.width*r.height)
}
