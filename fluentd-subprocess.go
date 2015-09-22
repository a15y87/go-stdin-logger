package main

import (
	"github.com/fluent/fluent-logger-golang/fluent"
	"fmt"

	"time"
	"bufio"
	"strings"
	"flag"
	"os/exec"
)

func main() {
	var tag = flag.String("tag", "syslog", "fleuntd tag for logging")
	flag.Parse()
	// fmt.Println(*tag)

	logger, err := fluent.New(fluent.Config{FluentSocketPath: "/tmp/td-agent.sock", FluentNetwork: "unix"})
	if err != nil {
		fmt.Println(err)
	}

	var message []string

	cmd := exec.Command("./test.sh", "2>&1")
	output, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}
	cmd.Start()
	bio := bufio.NewReader(output)
//	for {
//		line, _, err := bio.ReadLine();
//		if err != nil {
//			break
//		}
//
//	}
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(3 * time.Second):

		output.Close()
		if err := cmd.Process.Kill(); err != nil {
		}
		<-done // allow goroutine to exit
	case err := <-done:
		if err !=nil {
			line, _, err := bio.ReadLine()
			if err != nil {
				break
			}
			message = append(message, string(line))
		}
	}


	// fmt.Println(time.Now(), "\n", strings.Join(message, "\n") )

	defer logger.Close()

	var data = map[string]string {
		"message": strings.Join(message, "\n" ),
		"timestamp":  time.Now().String() }

	error := logger.Post(*tag, data)
	if error != nil {
		panic(error)
	}
}