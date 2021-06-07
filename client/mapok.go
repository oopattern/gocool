package main

import "fmt"

func main() {
	oids := make(map[string]int64)
	oids["xxx"] = 10
	oids["yyy"] = 20
	oids["zzz"] = 30
	fmt.Printf("xxx=%v\n", oids["xxx"]/10)
	fmt.Printf("yyy=%v\n", oids["yyy"]/10)
	fmt.Printf("111=%v\n", oids["111"]/10)
	fmt.Printf("222=%v\n", oids["222"]/10)
}
