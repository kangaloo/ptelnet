package portscheck

import (
	"fmt"
	"github.com/fatih/color"
	"sync"
)

type summary struct {
	sync.WaitGroup
	addr           string
	chanSuccessful chan string
	chanFailed     chan string
	successful     []string
	failed         []string
}

func (s *summary) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	s.Add(2)
	go func() {
		defer s.Done()
		for {
			port, ok := <-s.chanSuccessful
			if !ok {
				return
			}
			s.successful = append(s.successful, port)
		}
	}()

	go func() {
		defer s.Done()
		for {
			port, ok := <-s.chanFailed
			if !ok {
				return
			}
			s.failed = append(s.failed, port)
		}
	}()

	s.Wait()
}

func (s *summary) finish() {
	close(s.chanSuccessful)
	close(s.chanFailed)
}

func (s *summary) showSummary() {
	fmt.Println()
	fmt.Println(color.YellowString("Host: %s", s.addr))
	fmt.Println(color.GreenString("--------successful--------"))
	for _, port := range s.successful {
		fmt.Printf("%s %s %s\n", s.addr, port, color.GreenString("LISTEN"))
	}
	fmt.Println()

	fmt.Println(color.RedString("----------failed----------"))
	for _, port := range s.failed {
		fmt.Printf("%s %s %s\n", s.addr, port, color.RedString("FAILED"))
	}
	fmt.Println()
	fmt.Println()
}
