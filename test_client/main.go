package main

import (
	"fmt"
)

func main() {
	fmt.Println(renderMarkdown("# 标题1"))
	fmt.Println(makeQuery("select * from test.test;"))
}
