package main

import (
	"fmt"
	"sort"
	"strconv"
	"time"
)

// сюда писать код
func SingleHash(in, out chan interface{}) {
	for val := range in {
		strval := strconv.Itoa(val.(int))
		res := DataSignerCrc32(strval) + "~" + DataSignerCrc32(DataSignerMd5(strval))
		fmt.Println("SingleHash result ", res)
		out <- res
	}
}

func MultiHash(in, out chan interface{}) {
	hashvals := []int{0, 1, 2, 3, 4, 5}
	for val := range in {
		var res string
		for _, hv := range hashvals {
			res += DataSignerCrc32(strconv.Itoa(hv) + val.(string))
		}
		fmt.Println("MultiHash result ", res)
		out <- res
	}
}

func CombineResults(in, out chan interface{}) {
	var vals []string
	for val := range in {
		vals = append(vals, val.(string))
	}
	sort.Strings(vals)
	var res string
	for idx, val := range vals {
		res += val
		if idx != len(vals)-1 {
			res += "_"
		}
	}
	fmt.Println("CombineResults ", res)
	out <- res
}

func ExecutePipeline(jobs ...job) {
	chans := make([]chan interface{}, len(jobs)+1)
	for idx := 0; idx < len(jobs)+1; idx++ {
		chans[idx] = make(chan interface{}, 1)
	}

	for idx, worker := range jobs {
		go worker(chans[idx], chans[idx+1])
	}
	time.Sleep(3 * time.Second)
	for _, channel := range chans {
		close(channel)
		time.Sleep(time.Second)
	}
}

func main() {
}
