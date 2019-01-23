package main

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"
)

func crc32Async(wg *sync.WaitGroup, data string, out chan<- string) {
	defer wg.Done()
	out <- DataSignerCrc32(data)
}

func singleHashWorker(wg *sync.WaitGroup, data1 string, data2 string, out chan<- interface{}) {
	defer wg.Done()
	resChan := make(chan string, 2)
	wg.Add(1)
	go crc32Async(wg, data1, resChan)
	time.Sleep(time.Millisecond)
	wg.Add(1)
	go crc32Async(wg, data2, resChan)
	res := <-resChan + "~" + <-resChan
	out <- res
	fmt.Println("SingleHash result ", res)
}

func multiHashWorker(wg *sync.WaitGroup, input []string, out chan<- interface{}) {
	defer wg.Done()
	resChan := make(chan string, 6)
	for _, inp := range input {
		wg.Add(1)
		go crc32Async(wg, inp, resChan)
		time.Sleep(time.Millisecond)
	}

	var res string
	for idx := 0; idx < len(input); idx++ {
		res += <-resChan
	}
	out <- res
	fmt.Println("MultiHash result ", res)
}

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for val := range in {
		strval := strconv.Itoa(val.(int))
		md5 := DataSignerMd5(strval)
		wg.Add(1)
		go singleHashWorker(wg, strval, md5, out)
	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	hashvals := []int{0, 1, 2, 3, 4, 5}
	wg := &sync.WaitGroup{}
	for val := range in {
		var inStr []string
		for _, hv := range hashvals {
			inStr = append(inStr, strconv.Itoa(hv)+val.(string))
		}
		wg.Add(1)
		go multiHashWorker(wg, inStr, out)
	}
	wg.Wait()
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
	out <- res
	fmt.Println("CombineResults ", res)
}
func jobWorker(wg *sync.WaitGroup, jobFn job, in, out chan interface{}) {
	defer wg.Done()
	defer close(out)

	jobFn(in, out)
}
func ExecutePipeline(jobs ...job) {
	chans := make([]chan interface{}, len(jobs)+1)
	for idx := 0; idx < len(jobs)+1; idx++ {
		chans[idx] = make(chan interface{}, 1)
	}

	wg := &sync.WaitGroup{}
	for idx, worker := range jobs {
		wg.Add(1)
		go jobWorker(wg, worker, chans[idx], chans[idx+1])
	}

	wg.Wait()
}

func main() {
}
