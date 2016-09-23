package main

import (
	"time"
)

type SortList []*DirectoryInfo

func (list SortList) Len() int {
	return len(list)
}

func (list SortList) Swap(i int, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list SortList) Less(i int, j int) bool {
	// ランキングロジック
	return float64(list[i].Size)*time.Since(list[i].ModTime).Hours() > float64(list[j].Size)*time.Since(list[j].ModTime).Hours()
}
