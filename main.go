package main

import (
	"math/rand"
	"time"
	"fmt"
	"net"	
	."./go-zabbix" 
	"log"
	"github.com/google/gopacket/pcap"
	
	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/raw"
	)

const (
	defaultHost  = `localhost`
	//defaultHost = `remote.fibertrade.ru`
	defaultPort  = 10051
	etherType = 0x0800
)

func main() {
  go zabbixHello("SFP-SLA_4401")
  // Find all devices
  devices, err := pcap.FindAllDevs()
  if err != nil {
	  log.Fatal(err)
  }
  var net_name string
  // Print device information
  fmt.Println("Devices found:")
  for _, device := range devices {
	  fmt.Println("\nName: ", device.Name)
	  fmt.Println("Description: ", device.Description)
	  fmt.Println("Devices addresses: ", device.Description)
	  for _, address := range device.Addresses {
		  fmt.Println("- IP address: ", address.IP)
		  if address.IP.Equal(net.ParseIP("10.0.10.115")) {
			  net_name=device.Name
		  }
		  fmt.Println("- Subnet mask: ", address.Netmask)
	  }
  }
   // Open a raw socket on the specified interface, and configure it to accept

	ifi, err := net.InterfaceByName(net_name)
	if err != nil {
		log.Fatalf("failed to find interface %q: %v", net_name, err)
	}
	fmt.Println("Net_NAME: %q", net_name)
	fmt.Println("interface: %q", ifi.Name)
	c, err := raw.ListenPacket(ifi, etherType, nil)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Default message to system's hostname if empty.
	msg := []byte{0x45,			// 14
						 0x00,			// 15 QoS
						 0x00,0x32,		// 16 Size_Frame
						 0x00,0x00,		// 18 Id
						 0x00,0x00,		// 20 Frag
						 0xFF,			// 22 TTL
						 0x5e,			// 23 Protocol
						 0x88, 0x3e,	// 24 CRC
						 10,0,10,115,	// 26 IP Source (Server)
						 10,1,10,144,	// 30 IP Dist (SFP-SLA_1)
						 0xFA,			// 34 Id
						 10,1,10,140,	// 35 IP SFP-SLA_2
						 0xAA,0xAA,0xAA,
						 0xAA,0xAA,0xAA,0xAA, // 39 Time_stamp_11
						 0xBB,0xBB,0xBB,
						 0xBB,0xBB,0xBB,0xBB, // 46 Time_stamp_11
						 0xCC,0xCC,0xCC,
						 0xCC,0xCC,0xCC,0xCC, // 53 Time_stamp_11
						 0x01,0x23,0x45,0x67 } // 60 Number_count
					


	// Send messages in one goroutine, receive messages in another.
	go sendMessages(c, ifi.HardwareAddr, msg)
	go receiveMessages(c, ifi.MTU)

	// Block forever.
	select {}
}

// sendMessages continuously sends a message over a connection at regular intervals,
// sourced from specified hardware address.
func sendMessages(c net.PacketConn, source net.HardwareAddr, msg []byte) {
	// Message is broadcast to all machines in same network segment.
	
	f := &ethernet.Frame{
		Destination: ethernet.Broadcast,
		Source:      source,
		EtherType:   0x0800,
		Payload:     []byte(msg),
	}

	b, err := f.MarshalBinary()
	if err != nil {
		log.Fatalf("failed to marshal ethernet frame: %v", err)
	}

	// Required by Linux, even though the Ethernet frame has a destination.
	// Unused by BSD.
	addr := &raw.Addr{
		HardwareAddr: ethernet.Broadcast,
	}

	// Send message forever.
	t := time.NewTicker(1 * time.Second)
	for range t.C {
		if _, err := c.WriteTo(b, addr); err != nil {
			log.Fatalf("failed to send message: %v", err)
		}
	}
}

// receiveMessages continuously receives messages over a connection. The messages
// may be up to the interface's MTU in size.
func receiveMessages(c net.PacketConn, mtu int) {
	var f ethernet.Frame
	b := make([]byte, mtu)

	// Keep receiving messages forever.
	for {
		n, addr, err := c.ReadFrom(b)
		if err != nil {
			log.Fatalf("failed to receive message: %v", err)
		}

		// Unpack Ethernet II frame into Go representation.
		if err := (&f).UnmarshalBinary(b[:n]); err != nil {
			log.Fatalf("failed to unmarshal ethernet frame: %v", err)
		}

		// Display source of message and message itself.
		if f.Payload[20]==0xFA	{
		    fmt.Printf("[%s] %s", addr.String(), string(f.Payload))
		}
	}
}

  


func zabbixHello(host string){
	var delay int
	for  {
		delay = rand.Intn(1500)
		//delay:=i*100
	var metrics []*Metric
    metrics = append(metrics, NewMetric(host, "delay",fmt.Sprint(delay),time.Now().Unix()))
   
    // Create instance of Packet class
    packet := NewPacket(metrics)
	//fmt.Println(packet);
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