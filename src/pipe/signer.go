package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

func doProcedure(procedure job, wg *sync.WaitGroup, in, out chan interface{}) {
	defer close(out)
	defer wg.Done()
	procedure(in, out)
}

func ExecutePipeline(jobs ...job) {

	wg := &sync.WaitGroup{}
	in := make(chan interface{}, 100)

	for _, job := range jobs {
		wg.Add(1)
		out := make(chan interface{}, 100)
		go doProcedure(job, wg, in, out)
		in = out
	}
	wg.Wait()
}

func handleDataSignerCrc32(str string, in chan<- string) {
	in <- DataSignerCrc32(str)
}

func SingleHash(in, out chan interface{}) {

	wg := &sync.WaitGroup{}
	md5Wait := &sync.Mutex{}

	for val := range in {
		wg.Add(1)
		valueString := strconv.Itoa(val.(int))

		go func(value string, md5Wait *sync.Mutex) {
			defer wg.Done()

			md5Ch := make(chan string)
			crcCh := make(chan string, 1)
			crcMd5Ch := make(chan string, 1)

			go func(val string, md5Ch chan<- string, md5Wait *sync.Mutex) {
				md5Wait.Lock()
				md5Ch <- DataSignerMd5(val)
				md5Wait.Unlock()
			}(value, md5Ch, md5Wait)

			go handleDataSignerCrc32(value, crcCh)
			v := <-md5Ch
			go handleDataSignerCrc32(v, crcMd5Ch)

			result := <-crcCh + "~" + <-crcMd5Ch
			out <- result

		}(valueString, md5Wait)
	}
	wg.Wait()

}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}

	for val := range in {
		wg.Add(1)
		valAsStr := val.(string)
		go func() {
			var hashesChan = [6]chan string{}

			for th := 0; th <= 5; th++ {
				hashesChan[th] = make(chan string)
				inp := strconv.Itoa(th) + valAsStr
				go handleDataSignerCrc32(inp, hashesChan[th])
			}
			var result string
			for th := 0; th <= 5; th++ {
				result += <-hashesChan[th]
			}
			out <- result
			wg.Done()
		}()
	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}) {

	var hashs []string

	for data := range in {
		hashs = append(hashs, data.(string))
	}
	sort.Strings(hashs)
	out <- strings.Join(hashs, "_")
}
