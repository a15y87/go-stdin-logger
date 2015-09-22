package main

import (
	"fmt"
	"time"
	"flag"
//	"os/exec"
	"log"
//	"bytes"
//	"os"
//	"syscall"
	"github.com/samuel/go-zookeeper/zk"
	"crypto/sha256"
	"encoding/base64"
	"strings"
	"github.com/satori/go.uuid"
)

func zkConnect() *zk.Conn {
	zks := []string{"wz-zk-1.dol.cx:2181", "wz-zk-2.dol.cx:2181", "wz-zk-3.dol.cx:2181"}
	conn, _, err := zk.Connect(zks, time.Second)
	if err != nil {	log.Fatal(err)	}
	return conn
}


func zkSimpleParty(task_name string) error {
	acl := zk.WorldACL(zk.PermAll)
	var err error = nil

	my_uuid := uuid.NewV4().String()
	hasher := sha256.New()
	task_hash := base64.URLEncoding.EncodeToString( hasher.Sum([]byte(task_name)) )

	zk_path := strings.Join ([]string{"/cron_", task_hash}, "")
	zk_path_uuid := strings.Join ([]string{zk_path, "/", my_uuid}, "")
	log.Print(zk_path, "\n", zk_path_uuid)

	zk_c := zkConnect()
	defer zk_c.Close()


	exist, stat, err := zk_c.Exists(zk_path)

	if err != nil {	log.Fatal(err); return err }
	if exist != true {
		msg, err := zk_c.Create(zk_path, []byte{}, 0, acl)
		if err != nil {	log.Fatal(err); return err}
		log.Print(msg)
	} else {
		if stat.NumChildren == 0 {
			if err := zk_c.Delete(zk_path, -1) ; err != nil  {	log.Fatal(err) }
			msg, err := zk_c.Create(zk_path, []byte{}, 0, acl)
			if err != nil {	log.Fatal(err); return err }
			log.Print(msg)
		}
	}

	msg, err := zk_c.Create(zk_path_uuid, []byte{}, int32(zk.FlagEphemeral + zk.FlagSequence), acl)
	if err != nil {	log.Fatal(err); return err }
	log.Print(msg)

	children, stat, err := zk_c.Children(zk_path)
	if err != nil {	log.Fatal(err); return err	}
	log.Print(fmt.Sprintf("%+v %+v\n", children, stat))

	return err
}


func main() {

	var task_name = flag.String("task", "test", "task name")
//	var task_timer =  flag.Int("timer", 0, "task timer")
//	var app_name = flag.String("app", "./test.sh", "app")
	flag.Parse()

	for i:=0; i<30000000; i++ { go zkSimpleParty(*task_name) }




}

