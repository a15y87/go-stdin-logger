package main

import (
	"fmt"
	"time"
	"flag"
	"os/exec"
	"log"
	"bytes"
	"os"
	"syscall"
	"github.com/samuel/go-zookeeper/zk"
	"crypto/sha256"
	"encoding/base64"
)

func main() {
//	var tag = flag.String("tag", "syslog", "fluentd tag for logging")
//	var fluent_socket = flag.String("socket", "/tmp/td-agent.sock", "fluentd socket for logging")
	var task_name = flag.String("task", "test", "task name")
	var task_timer =  flag.Int("timer", 0, "task timer")
//	var message_subject =  flag.String("subject", "subject", "message subject")
//	var message_enable = flag.Bool("debug", false, "send email")


	zk_c, _, err := zk.Connect([]string{"wz-zk-1.dol.cx:2181", "wz-zk-2.dol.cx:2181", "wz-zk-3.dol.cx:2181"}, time.Second)
	if err != nil {
		panic(err)
	}
	children, stat, _, err := zk_c.ChildrenW("/")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v %+v\n", children, stat)
	//e := <-ch
	//fmt.Printf("%+v\n", e)
	zk_c.Close()


	var app_name = flag.String("app", "./test.sh", "app")
	flag.Parse()

	cmd := exec.Command(*app_name, "2>&1")
	var cmd_result int

	randomBytes := &bytes.Buffer{}
	cmd.Stdout = randomBytes

	// Start command asynchronously
	if err = cmd.Start(); err != nil {	log.Fatal(err) }

	hasher := sha256.New()
	task_hash := base64.StdEncoding.EncodeToString( hasher.Sum([]byte(*task_name)) )

	log.Print(*task_timer, task_hash)
	time_start := time.Now()

	// Create a ticker that outputs elapsed time
//	ticker := time.NewTicker(time.Second)
//	go func(ticker *time.Ticker) {
//		now := time.Now()
//		for _ = range ticker.C {
//			log.Print (
//				fmt.Sprintf("%s:\n%s", time.Since(now), string(randomBytes.Bytes()) ),
//			)
//		}
//	}(ticker)

	if *task_timer > 0 {
		timer := time.NewTimer(time.Second * time.Duration(*task_timer))
		go func (timer *time.Timer, cmd *exec.Cmd) {
		for _ = range timer.C {
			err = cmd.Process.Signal(os.Kill)
			if err != nil {	log.Fatal(err) }
		}
		}(timer, cmd)
	}


	if err = cmd.Wait(); err != nil {
		exit_code, _ :=err.(*exec.ExitError)
		status, _ :=  exit_code.Sys().(syscall.WaitStatus)
		cmd_result = status.ExitStatus()
	}


	time_end := time.Now()
	time_delta := time_end.Sub(time_start)

	log.Print(fmt.Sprintf("%d bytes generated!, %d \ntime:\n%s\n%s\n%s", len(randomBytes.Bytes()), cmd_result, time_start, time_end, time_delta))


}
