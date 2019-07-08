package logger

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

// init logger
func init() {
	// logfile workdir/ptelnet.report.log.timestamp
	ts := strconv.Itoa(int(time.Now().Unix()))

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	logFile := fmt.Sprintf("%s%s%s%s", wd, string(os.PathSeparator), "ptelnet.report.", ts)

	f, err := os.Create(logFile)
	if err != nil {
		log.SetOutput(os.Stderr)
	}

	log.SetOutput(f)
}
