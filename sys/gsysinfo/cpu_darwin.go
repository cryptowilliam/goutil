package gsysinfo

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
)

func CpuLimit(processName string) {
	flag.Parse()
	targets := make([]string, 0, 1)
	if *fExe != "" {
		targets = append(targets, *fExe)
	}
	var err error
	oneSecond := time.Second
	mtx := sync.Mutex{}
	running := true
	procMap := make(map[int]*os.Process, 16)
	if *fPid > 0 {
		var err error
		if procMap[*fPid], err = os.FindProcess(*fPid); err != nil {
			log.Fatalf("cannot find %d: %v", *fPid, err)
		}
	} else {

		if flag.NArg() > 0 {
			targets = append(targets, flag.Args()...)
		}
		for i, exe := range targets {
			if exe[0] != '/' {
				exe, err = exec.LookPath(exe)
				if err != nil {
					log.Printf("cannot find full path for %q: %s", exe, err)
					continue
				}
				targets[i] = exe
			}
		}
		go func() {
			var (
				ok        bool
				null      struct{}
				processes []*os.Process
				oldpids   = make(map[int]struct{}, 16)
				times     int
			)
			for {
				processes = getProcesses(processes[:0], targets)
				if len(processes) == 0 {
					if *fTimeout > 0 {
						times++
						if times > *fTimeout {
							log.Println("no more processes to watch, timeout reached - exiting.")
							running = false
							return
						}
					}
				} else {
					mtx.Lock()
					for k := range procMap {
						oldpids[k] = null
					}
					for _, p := range processes {
						if _, ok = procMap[p.Pid]; !ok {
							log.Printf("new process %d", p.Pid)
						}
						procMap[p.Pid] = p
						delete(oldpids, p.Pid)
					}
					for k := range oldpids {
						log.Printf("%d exited", k)
						delete(procMap, k)
						delete(oldpids, k)
					}
					mtx.Unlock()
				}
				time.Sleep(oneSecond)
			}
		}()
	}

	stopped := false
	var (
		sig   os.Signal
		sleep time.Duration
		n     int64
	)
	tbd := make([]int, 0, 2)
	run := time.Duration(10*(*fLimit)) * time.Millisecond
	freeze := time.Duration(1000)*time.Millisecond - run
	for running {
		mtx.Lock()
		n = int64(len(procMap))
		if n == 0 {
			sleep = oneSecond
		} else {
			if stopped {
				sig, stopped, sleep = syscall.SIGCONT, false, time.Duration(int64(run)/n)
			} else {
				sig, stopped, sleep = syscall.SIGSTOP, true, freeze
			}
			tbd = tbd[:0]
			for pid, p := range procMap {
				if err = p.Signal(sig); err != nil {
					if strings.HasSuffix(err.Error(), "no such process") {
						log.Printf("%d vanished.", pid)
					} else {
						log.Printf("error signaling %d: %s", pid, err)
					}
					tbd = append(tbd, pid)
				}
			}
			if len(tbd) > 0 {
				for _, pid := range tbd {
					delete(procMap, pid)
				}
			}
		}
		mtx.Unlock()
		time.Sleep(sleep)
	}
}
