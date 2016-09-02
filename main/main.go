package main

import (
	"bufio"
	"net"
	"encoding/json"
	"container/list"
	"log"
	"os"
	"sync"
	"fmt"
)
var (
	meas Sample
	smpvs SMPVS

	smpvToSend = list.New()
	stringsToSend = list.New()
	invertersToSend = list.New()
	farmsToSend = list.New()
)

var mu sync.Mutex
var subs = list.New()
func main() {


	log.SetOutput(os.Stdout)

	go emit()

	ln, err := net.Listen("tcp", net.JoinHostPort("", genPort()))
	if err != nil {

	}
	defer ln.Close()
	log.Printf("Listen on %s", ln.Addr())
	for {
		conn, err := ln.Accept()
		if err != nil {

		}
		log.Printf("Connection from %s", conn.RemoteAddr())
		addSubscriber(conn)
	}
}

func genPort() string {
	port := os.Getenv("GEN_PORT")
	if port != "" {
		return port
	}
	// TODO: remove backwards compatibility
	port = os.Getenv("GENERATOR_PORT")
	if port != "" {
		return port
	}
	return "3009"
}

func listenPort() string {
	port := os.Getenv("LISTEN_PORT")
	if port != "" {
		return port
	}
	return "3000"
}

func listenIp() string {
	ip := os.Getenv("LISTEN_IP")
	if ip != "" {
		return ip
	}
	return "127.0.0.1"
}

func addSubscriber(conn net.Conn) {
	mu.Lock()
	defer mu.Unlock()
	subs.PushBack(NewSubscriber(conn))
}


func emit(){
	conn, err := net.Dial("tcp", net.JoinHostPort(listenIp(), listenPort()))
	if err != nil {
		// handle error
	}
	for{
		status, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {

		}
		err = json.Unmarshal([]byte(status), &meas)
		if err != nil {

		}
		err = json.Unmarshal([]byte(meas.Value), &smpvs)
		if err != nil {

		}
		smpvToSend = list.New()
		stringsToSend = list.New()
		invertersToSend = list.New()
		farmsToSend = list.New()

		for x := 0; x < len(smpvs.Data); x++ {
				buildStructure(smpvs.Data[x])
			}

			updateString()
			updateInv()
			updateFarms()

			pv := smpvToSend.Front()
			for pv != nil {
				smpv := pv.Value.(*SMPVToSend)
				updateSubscribersPv(smpv)
				pv = pv.Next()
			}

			string := stringsToSend.Front()
			for string != nil {
				stringObj := string.Value.(*Sstring)
				fmt.Println(stringObj.SMPV_Pch.V)
				fmt.Println(stringObj.SMPV_Uch.V)
				fmt.Println(stringObj.SMPV_Ich.V)
				updateSubscribersStr(stringObj)
				string = string.Next()
			}

			inv := invertersToSend.Front()
			for inv != nil {
				invObj := inv.Value.(*Inverter)
				fmt.Println(invObj.SMPV_Pch.V)
				fmt.Println(invObj.SMPV_Uch.V)
				fmt.Println(invObj.SMPV_Ich.V)
				updateSubscribersInv(invObj)
				inv = inv.Next()
			}

			farm := farmsToSend.Front()
			for farm != nil {
				farmObj := farm.Value.(*Farm)
				fmt.Println(farmObj.SMPV_Pch.V)
				fmt.Println(farmObj.SMPV_Uch.V)
				fmt.Println(farmObj.SMPV_Ich.V)
				updateSubscribersFarm(farmObj)
				farm = farm.Next()
			}

	}
}

type Subscriber struct {
	conn net.Conn
	w    *bufio.Writer
	enc  *json.Encoder
}

func NewSubscriber(conn net.Conn) *Subscriber {
	w := bufio.NewWriter(conn)
	enc := json.NewEncoder(w)
	return &Subscriber{conn: conn, w: w, enc: enc}
}

func (s *Subscriber) Close() {
	s.conn.Close()
}
