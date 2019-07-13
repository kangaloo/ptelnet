package portscheck

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

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

// collect summary
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

func NewHost(info string) (*Host, error) {
	i := strings.Split(info, ":")
	ports := parsePorts(i[1])
	host := &Host{
		wg:    &sync.WaitGroup{},
		sumWg: &sync.WaitGroup{},
		addr:  i[0],
		ports: ports,
		summary: &summary{
			addr:           i[0],
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
// report summary
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

func (h *Host) portsCheck() error {
	for _, port := range h.ports {
		_, err := strconv.Atoi(port)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *Host) showSummary() {
	h.summary.showSummary()

	f, ok := log.StandardLogger().Out.(*os.File)
	if !ok {
		return
	}

	file := f.Name()
	fmt.Printf("see details in file %s\n", file)
}
