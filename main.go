package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"syscall"
	"time"

	"github.com/anmitsu/go-shlex"
	"github.com/lwf/chainlib"
)

var nonalphanumeric = regexp.MustCompile("[^a-zA-Z0-9\\-]")

type durationFlag time.Duration

func (f *durationFlag) String() string {
	return ""
}

func (f *durationFlag) Set(s string) error {
	d, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	du := durationFlag(d)
	*f = du
	return nil
}

var windowSize durationFlag = durationFlag(time.Minute * 10)
var maxRestarts int
var metadataDir string
var metadataKey string
var svCmd string

type Metadata struct {
	Restarts []time.Time
}

func assert(err error) {
	if err != nil {
		fatal(err.Error())
	}
}

func window(w []time.Time, lower time.Time, size time.Duration) []time.Time {
	nw := make([]time.Time, 0)
	for _, i := range w {
		if lower.Sub(i) < size {
			nw = append(nw, i)
		}
	}
	return nw
}

func limit(metadata *Metadata, size time.Duration, max int) bool {
	restarts := window(metadata.Restarts, time.Now(), size)
	if len(restarts) >= max {
		return true
	}
	metadata.Restarts = append(restarts, time.Now())
	return false
}

func main() {
	flag.Var(&windowSize, "window-size", "Window size in Go `duration` format (default \"10m\")")
	flag.IntVar(&maxRestarts, "max-restarts", 5, "Max `restarts` within window-size duration")
	flag.StringVar(&metadataDir, "metadata-dir", "/run/runlimit", "Metadata `dir`, where metadata files are stored")
	flag.StringVar(&metadataKey, "metadata-key", "", "Metadata key, which will form part of the metadata file name")
	flag.StringVar(&svCmd, "sv-cmd", "", "Command to use to stop a service")
	flag.Parse()

	if windowSize == 0 || maxRestarts == 0 {
		fatal("-max-restarts and/or -window-size cannot be 0")
	}

	cmdline := flag.Args()
	if len(cmdline) < 1 {
		fatal("No command supplied")
	}

	cwd, err := os.Getwd()
	assert(err)
	if metadataKey == "" {
		metadataKey = nonalphanumeric.ReplaceAllString(cwd, "_")
	}

	metafile := filepath.Join(metadataDir, fmt.Sprintf("%s.meta", metadataKey))
	f, err := os.OpenFile(metafile, os.O_RDWR|os.O_CREATE, os.FileMode(0644))
	assert(err)
	defer f.Close()
	assert(syscall.Flock(int(f.Fd()), syscall.LOCK_NB|syscall.LOCK_EX))
	metadata := &Metadata{}
	if err := json.NewDecoder(f).Decode(metadata); err != nil && err != io.EOF {
		warning("metadata corrupted, ignoring...")
	}
	if limit(metadata, time.Duration(windowSize), maxRestarts) {
		signals := make(chan os.Signal)
		signal.Notify(signals, syscall.SIGTERM)
		if svCmd != "" {
			parts, err := shlex.Split(svCmd, true)
			assert(err)
			go func() {
				if out, err := exec.Command(parts[0], parts[1:]...).Output(); err != nil {
					warning("command exited abnormally with output %s", string(out))
				}
			}()
			select {
			case <-signals:
				break
			case <-time.After(5 * time.Second):
				warning("timed out while waiting for TERM from %s", parts[0])
			}
		}
		fatal("max restart intensity reached")
	}
	assert(syscall.Ftruncate(int(f.Fd()), 0))
	_, err = syscall.Seek(int(f.Fd()), 0, 0)
	assert(err)
	if err := json.NewEncoder(f).Encode(metadata); err != nil {
		warning("could not write metadata: %s", err.Error())
	}
	assert(chainlib.Exec(cmdline, nil))
}
