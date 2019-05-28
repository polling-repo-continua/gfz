package libgrequest

import (
	"fmt"
	"sync"
)

// Processor : channel controlcenter
func Processor(s *State) {
	N := TotalRequests(s.Fuzzer.Maxes)
	payloadCh := make(chan string)
	codeCh := make(chan *int, s.Threads)
	procWG := new(sync.WaitGroup)
	procWG.Add(s.Threads)

	go func() { GetURL(s, 0, s.Commandline, payloadCh) }() // Payload is just a string with 'FUZZ'

	for i := 0; i < s.Threads; i++ {
		go func() {
			for {
				// if s.Terminate == true {break}
				code, _ := s.Request(s, s.URL, s.Cookies, <-payloadCh)
				codeCh <- code
			}
		}()
		procWG.Done()
	}
	for r := 0; r < N; r++ {
		<-codeCh
		fmt.Printf("[+] requests: %d/%d\r", s.Counter.v, N)
	}
	procWG.Wait()
}

//TODO: add --recursive
