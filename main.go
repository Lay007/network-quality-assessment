package main

import (
    "time"	
	."./go-zabbix" 
	)

const (
	defaultHost  = `localhost`
	//defaultHost = `remote.fibertrade.ru`
    defaultPort  = 10051
)

func main() {
    var metrics []*Metric
    metrics = append(metrics, NewMetric("SFP-SLA_4401", "delay", "1.22", time.Now().Unix()))
    metrics = append(metrics, NewMetric("SFP-SLA_4401", "status", "OK"))

    // Create instance of Packet class
    packet := NewPacket(metrics)

    // Send packet to zabbix
    z := NewSender(defaultHost, defaultPort)
    z.Send(packet)
}

/*import (
	"encoding/json"
	"log"
	"net/http"
	
)

func main() {
	http.HandleFunc("/", mainPage)
	port := ":313"
	println("Server Listen on port", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	user := User{"Alex", "Net"}
	js, _ := json.Marshal(user)
	w.Write(js)
}
*/