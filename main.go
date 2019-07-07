package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// todo 增加summary功能

func main() {
	file := "file.txt"
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	hosts, err := NewHosts(f)
	if err != nil {
		panic(err)
	}

	hosts.Check()
}

type Hosts struct {
	hosts []*Host
}

func NewHosts(f *os.File) (*Hosts, error) {
	defer f.Close()
	var hosts []*Host
	reader := bufio.NewReader(f)
	for {
		l, err := reader.ReadString('\n')
		l = strings.TrimSuffix(l, "\n")
		l = strings.TrimSpace(l)

		// 去空行
		// todo 支持去掉#注释的行
		if len(l) != 0 {
			h, err := NewHost(l)
			if err != nil {
				return nil, err
			}
			hosts = append(hosts, h)
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}

	return &Hosts{hosts: hosts}, nil
}

// todo 收集summary
func (hs *Hosts) Check() {
	for _, host := range hs.hosts {
		host.PortsConnectTest()
		host.showSummary()
	}
}

type Host struct {
	wg      *sync.WaitGroup
	sumWg   *sync.WaitGroup
	addr    string
	ports   []string
	summary *summary
}

func (h *Host) showSummary() {
	fmt.Println()
	fmt.Println(color.YellowString("Host: %s", h.addr))
	h.summary.showSummary()
}

type summary struct {
	sync.WaitGroup
	chanSuccessful chan string
	chanFailed     chan string
	successful     []string
	failed         []string
}

func (s *summary) showSummary() {
	fmt.Println(color.GreenString("--------successful--------"))
	for _, port := range s.successful {
		fmt.Printf("%s, ", port)
	}
	fmt.Println()

	fmt.Println(color.RedString("----------failed----------"))
	for _, port := range s.failed {
		fmt.Printf("%s, ", port)
	}
	fmt.Println()
	fmt.Println()
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

func NewHost(info string) (*Host, error) {
	log.Printf("%s\n", info)
	i := strings.Split(info, ":")
	ports := parsePorts(i[1])
	host := &Host{
		wg:    &sync.WaitGroup{},
		sumWg: &sync.WaitGroup{},
		addr:  i[0],
		ports: ports,
		summary: &summary{
			chanSuccessful: make(chan string, 64),
			chanFailed:     make(chan string, 64),
		},
	}

	if err := host.portsCheck(); err != nil {
		return nil, err
	}
	return host, nil
}

// PortsTest
// 上报summary
func (h *Host) PortsConnectTest() {
	// start count goroutine
	h.sumWg.Add(1)
	go h.summary.Start(h.sumWg)

	for _, port := range h.ports {
		h.wg.Add(1)
		go connect(h.addr+":"+port, h.wg, h.summary)
	}

	h.wg.Wait()

	// close summary channel and count goroutine will exit
	h.summary.finish()
	h.sumWg.Wait()

}

func parsePorts(portsString string) []string {
	var ports []string
	s := strings.Split(portsString, ",")

	for _, p := range s {
		if !strings.Contains(p, "-") {
			ports = append(ports, p)
		} else {

			// todo 检查错误
			ps := strings.Split(p, "-")
			start, _ := strconv.Atoi(ps[0])
			end, _ := strconv.Atoi(ps[1])

			for i := start; i <= end; i++ {
				ports = append(ports, strconv.Itoa(i))
			}
		}
	}

	return ports
}

func (h *Host) portsCheck() error {
	for _, port := range h.ports {
		_, err := strconv.Atoi(port)
		if err != nil {
			return err
		}
	}
	return nil
}

func connect(addr string, gw *sync.WaitGroup, sum *summary) {
	defer gw.Done()
	conn, err := net.DialTimeout("tcp", addr, time.Second*10)
	if err != nil {
		// verbose参数打开时，输出详细日志
		//log.Println(color.RedString("connect to %s failed: %v", addr, err))
		sum.chanFailed <- strings.Split(addr, ":")[1]
		return
	}
	_ = conn.Close()
	sum.chanSuccessful <- strings.Split(addr, ":")[1]
	// verbose参数打开时，输出详细日志
	//log.Println(color.GreenString("connect to %s successful", addr))

}