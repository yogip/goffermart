package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(time.Duration(3) * time.Second)
	fmt.Println(time.Duration(3))
	fmt.Println("--- ", time.Now())
	time.Sleep(time.Duration(3) * time.Second)
	fmt.Println("--- ", time.Now())
}
