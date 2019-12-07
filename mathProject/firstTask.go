package main

import (
	"fmt"
	"log"
	"sort"
)
func printValue(i int, n int){
	fmt.Printf("Цифра %d встречается %d раз", i,n)
	fmt.Println()
}
func printResult(i int){
	fmt.Print(i)
}
func main() {
	var value int
	var residue int
	counter := make([]int, 10)
	copyOfCounter := make([]int, 10)
	log.Println("Введите любое целое число")
	fmt.Scan(&value)
	temp := value
	for ;temp>0;{
	
		residue = temp % 10
		temp = temp / 10
		counter[residue]++
	}
	fmt.Print(counter)
	Result:=make(map[int]int)
	copy(copyOfCounter, counter)
	sort.Ints(counter)
	fmt.Print(counter)
	j:=0
	for i:=9;i>=0;i--{
		for j=0;j<=9 && copyOfCounter[j] != counter[i]; j++{
		}
		Result[9-i]=j
	}

	for i:=0;i<=9;i++{
		printValue(Result[i], copyOfCounter[Result[i]])
	}
	fmt.Print("Result: ")
	for i:=0;i<=9;i++{
		for ;copyOfCounter[Result[i]] >0; copyOfCounter[Result[i]]--{
			printResult(Result[i])
		}
}

}
