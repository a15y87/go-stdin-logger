package main

import (
	"github.com/fluent/fluent-logger-golang/fluent"
	"fmt"
	"os"
	"time"
	"bufio"
	"strings"
	"flag"
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

	bio := bufio.NewReader(os.Stdin)
	for {
		line, _, err := bio.ReadLine();
		if err != nil {
			break
		}
		message = append(message, string(line))
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