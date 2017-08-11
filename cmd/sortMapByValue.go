package main

import (
	"fmt"
	"sort"
)

func main() {
	mapStrInt := map[string]int{
		"a": 23,
		"b": 43,
		"c": 4,
		"d": 98,
		"e": 20,
		"f": 54,
	}

	fmt.Println(mapStrInt, sortMapByValue(mapStrInt, "desc"))
}

func sortMapByValue(mapsi map[string]int, direct string) PairList {
	//总结：优先使用方式二初始化slice
	/*
    //方式一
	var pl PairList
	for k, v := range mapsi {
		pl = append(pl, Pair{k, v})
	}
    */

	//方式二
	pl := make(PairList, len(mapsi))
	l := 0
	for k, v := range mapsi {
		pl[l] = Pair{k, v}
		l++
	}
	if direct == "asc" {
		sort.Sort(pl)
	} else {
		sort.Sort(sort.Reverse(pl))
	}
	return pl
}

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (pl PairList) Len() int { return len(pl) }
func (pl PairList) Less(i, j int) bool { return pl[i].Value < pl[j].Value }
func (pl PairList) Swap(i, j int) { pl[i], pl[j] = pl[j], pl[i] }

