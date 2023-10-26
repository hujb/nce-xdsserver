package main

import (
	"errors"
	"fmt"
	rs "github.com/hujb/nce-xdsserver/common/resource"
	"log"
)

func main() {
	//b := a
	//a = false
	//fmt.Println(a)
	//fmt.Println(b)
	rs1 := &rs.ResourceSnapshot{}
	rs1.InitResourceSnapshot()
	fmt.Println(rs1.GetVersion())

	rs2 := &rs.ResourceSnapshot{}
	rs2.InitResourceSnapshot()
	fmt.Println(rs2.GetVersion())

	//m := make(map[string]Shape)
	//m["a"] = &Rectangle{width: 1.0, height: 2.0}
	r := &Rectangle{width: 1.0, height: 3.0}
	log.Println("正方形修改前：", r)
	//fmt.Println("面积=", r.Area())
	r.Area()
	log.Println("正方形修改后：", r)
}

type Shape interface {
	Area() float32
}

func getRectangle() (rec *Rectangle, err error) {
	return nil, errors.New("aaa")
}
