package main

import (
	"github.com/fluent/fluent-logger-golang/fluent"
	"os"
	"time"
	"bufio"
	"strings"
	"flag"
	"log"
	"gopkg.in/gomail.v2"
)

func main() {
	var tag = flag.String("tag", "syslog", "fleuntd tag for logging")
	var fluent_socket = flag.String("socket", "/tmp/td-agent.sock", "fleuntd socket for logging")
	flag.Parse()
	// fmt.Println(*tag)

	logger, err := fluent.New(fluent.Config{FluentSocketPath: *fluent_socket, FluentNetwork: "unix"})
	if err != nil {
		log.Fatal(err)
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

	if err = logger.Post(*tag, data); err != nil {
		log.Fatal(err)
	}

	d := gomail.NewPlainDialer("127.0.0.1", 25, "", "")
	s, err := d.Dial()
	if err != nil {
		log.Fatal(err)
	}

	m := gomail.NewMessage()

	m.SetHeader("From", "from@mail")
	m.SetHeader("To", "to@mail")
	m.SetHeader("Subject", "go-fluentd-stdin")
	m.SetBody("text/html", strings.Join(message, "\n"))

	if err := gomail.Send(s, m); err != nil {
		log.Printf("Could not send email to %q: %v", "to@mail", err)
	}
	m.Reset()


}