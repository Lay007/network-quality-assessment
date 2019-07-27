package main

import (
	"math/rand"
	"time"
	"fmt"	
	."./go-zabbix" 
	"log"
    "github.com/google/gopacket/pcap"
	)

const (
	defaultHost  = `localhost`
	//defaultHost = `remote.fibertrade.ru`
    defaultPort  = 10051
)

func main() {
  go zabbixHello("SFP-SLA_4401")
  // Find all devices
  devices, err := pcap.FindAllDevs()
  if err != nil {
	  log.Fatal(err)
  }

  // Print device information
  fmt.Println("Devices found:")
  for _, device := range devices {
	  fmt.Println("\nName: ", device.Name)
	  fmt.Println("Description: ", device.Description)
	  fmt.Println("Devices addresses: ", device.Description)
	  for _, address := range device.Addresses {
		  fmt.Println("- IP address: ", address.IP)
		  fmt.Println("- Subnet mask: ", address.Netmask)
	  }
  }
  
}

func zabbixHello(host string){
	for  {
		var delay = rand.Intn(1500)
		//delay:=i*100
	var metrics []*Metric
    metrics = append(metrics, NewMetric(host, "delay",fmt.Sprint(delay),time.Now().Unix()))
   
    // Create instance of Packet class
    packet := NewPacket(metrics)
	fmt.Println(packet);
    // Send packet to zabbix
    z := NewSender(defaultHost, defaultPort)
	z.Send(packet)
	time.Sleep(5 * time.Second)
	}
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