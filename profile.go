package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Stat struct {
	code int
	time int
	size int
}

func handleProfile(url url.URL, profileSize int) {
	times := make([]int, profileSize)
	minSize := math.MaxInt32
	maxSize := 0
	codes := make(map[int]int)
	timeSum := 0
	results := make(chan Stat, profileSize)

	for i := 0; i < profileSize; i++ {
		go performRequestJob(url, results)
	}

	for i := 0; i < profileSize; i++ {
		stat := <-results
		times[i] = stat.time
		timeSum += stat.time
		maxSize = max(maxSize, stat.size)
		minSize = min(maxSize, stat.size)
		if val, ok := codes[stat.code]; ok {
			codes[stat.code] = val + 1
		} else {
			codes[stat.code] = 1
		}
	}
	sort.Ints(times)

	avg := (float64(timeSum)) / (float64(profileSize))
	median := times[profileSize/2]
	if profileSize%2 == 0 {
		median += times[profileSize/2-1]
		median /= 2
	}

	fmt.Println("--------------Profile--------------")
	fmt.Printf("Number of requests: %d\n", profileSize)
	fmt.Printf("Fastest request: %d miliseconds\n", times[0])
	fmt.Printf("Slowest request: %d miliseconds\n", times[len(times)-1])
	fmt.Printf("Avg request time: %f miliseconds\n", avg)
	fmt.Printf("Median request time: %d miliseconds\n", median)
	fmt.Printf("Percentage of successful requests: %d\n", profileSize)
	fmt.Println("Error codes:")
	for k, v := range codes {
		if k != 200 {
			fmt.Printf("\t%d occurence: %d\n", k, v)
		}
	}
	fmt.Printf("Size of smallest response : %d bytes\n", minSize)
	fmt.Printf("Size of largest response : %d bytes\n", maxSize)
	fmt.Println("-----------------------------------")
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a int, b int) int {
	if a < b {
		return b
	}
	return a
}

func handlePrintResponse(url url.URL) {
	resp, _ := performRequest(url)
	fmt.Printf("%s", resp)
}

func performRequestJob(url url.URL, results chan Stat) {
	start := time.Now()
	resp, code := performRequest(url)
	duration := time.Since(start)
	stat := Stat{code: code, time: int(duration.Milliseconds()), size: len(resp)}
	results <- stat
}

func performRequest(url url.URL) (r []byte, c int) {
	timeout := time.Second * 5
	dialer := net.Dialer{
		Timeout: timeout,
	}
	conn, err := tls.DialWithDialer(&dialer, "tcp", url.Hostname()+":https", nil)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	conn.Write([]byte("GET " + url.Path + " HTTP/1.0\r\nHost: " + url.Hostname() + "\r\n\r\n"))

	resp, err := ioutil.ReadAll(conn)
	response := string(resp)
	if err != nil {
		log.Fatal(err)
	}
	code, err := strconv.Atoi(response[9:12])
	if err != nil {
		log.Fatal(err)
	}
	conn.Close()

	if code > 300 && code < 400 {
		// get redirect url
		searchTerm := "Location: "
		start := strings.Index(response, searchTerm) + len(searchTerm)
		end := start + strings.Index(response[start:], "\r")
		redirectURLString := response[start:end]
		redirectURL, err := url.Parse(redirectURLString)
		if err != nil {
			log.Fatal(err)
		}
		return performRequest(*redirectURL)
	} else {
		return resp, code
	}
}
