package portscheck

import (
	log "github.com/sirupsen/logrus"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

func parsePorts(portsString string) []string {
	var ports []string
	s := strings.Split(portsString, ",")

	for _, p := range s {
		if !strings.Contains(p, "-") {
			ports = append(ports, p)
		} else {
			ps := strings.Split(p, "-")
			start, err := strconv.Atoi(ps[0])
			if err != nil {
				log.Fatalf("parse range ports failed, %s", err)
			}

			end, err := strconv.Atoi(ps[1])
			if err != nil {
				log.Fatalf("parse range ports failed, %s", err)
			}

			for i := start; i <= end; i++ {
				ports = append(ports, strconv.Itoa(i))
			}
		}
	}

	return ports
}

func connect(addr string, gw *sync.WaitGroup, sum *summary, timeout int) {
	defer gw.Done()
	now := time.Now()
	conn, err := net.DialTimeout("tcp", addr, time.Second*time.Duration(timeout))
	if err != nil {
		log.WithField("reason", err.Error()).Warnf("connect to %s failed after %v", addr, time.Since(now))
		sum.chanFailed <- strings.Split(addr, ":")[1]
		return
	}
	_ = conn.Close()
	log.Infof("connected to %s in %v", addr, time.Since(now))
	sum.chanSuccessful <- strings.Split(addr, ":")[1]
}
