package main

import (
	"bufio"
	"net"
	"encoding/json"
	"container/list"
	"log"
	"os"
	"sync"
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
var subsOut = list.New()
var subsIn = list.New()
func main() {
	log.SetOutput(os.Stdout)

	go emit()

  go listenIn()

	lnOut, errOut  := net.Listen("tcp", net.JoinHostPort("", genPort()))
	if errOut != nil {
		log.Fatal(errOut)
	}
	defer lnOut.Close()
	log.Printf("Listen on %s", lnOut.Addr())
	for {
		conn, err := lnOut.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Connection from %s", conn.RemoteAddr())
		addSubscriberOut(conn)
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

func addSubscriberOut(conn net.Conn) {
	mu.Lock()
	defer mu.Unlock()
	subsOut.PushBack(NewSubscriber(conn))
}

func addSubscriberIn(conn net.Conn) {
	mu.Lock()
	defer mu.Unlock()
	subsIn.PushBack(NewSubscriber(conn))
}


func listenIn(){
  lnIn, erriIn  := net.Listen("tcp", net.JoinHostPort(listenIp(), listenPort()))
	if erriIn != nil {
		log.Fatal(erriIn)
	}
	defer lnIn.Close()
	log.Printf("Listen on %s", lnIn.Addr())
	for {
		conn, err := lnIn.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Connection from %s", conn.RemoteAddr())
		addSubscriberIn(conn)
	}
}

func emit(){
	for{
		if subsIn.Len()!=0{
		e := subsIn.Front()
		s := e.Value.(*Subscriber)
		status := s.ReadMeas()

		err := json.Unmarshal([]byte(status), &meas)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal([]byte(meas.Value), &smpvs)
		if err != nil {
			log.Fatal(err)
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
				updateSubscribersStr(stringObj)
				string = string.Next()
			}

			inv := invertersToSend.Front()
			for inv != nil {
				invObj := inv.Value.(*Inverter)
				updateSubscribersInv(invObj)
				inv = inv.Next()
			}

			farm := farmsToSend.Front()
			for farm != nil {
				farmObj := farm.Value.(*Farm)
				updateSubscribersFarm(farmObj)
				farm = farm.Next()
			}
		}
	}
}
func (s *Subscriber) ReadMeas() string {
	status, err := s.r.ReadString('\n')
	if err!=nil {
		log.Fatal(err)
	}
	return status
}


type Subscriber struct {
	conn net.Conn
	w    *bufio.Writer
	r    *bufio.Reader
	enc  *json.Encoder
}

func NewSubscriber(conn net.Conn) *Subscriber {
	w := bufio.NewWriter(conn)
	r := bufio.NewReader(conn)
	enc := json.NewEncoder(w)
	return &Subscriber{conn: conn, w: w, r: r, enc: enc}
}

func (s *Subscriber) Close() {
	s.conn.Close()
}
