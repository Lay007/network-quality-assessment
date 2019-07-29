package main

import (
	. "./go-zabbix"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/google/gopacket/pcap"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/raw"
)

const (
	defaultHost = `localhost`
	//defaultHost = `remote.fibertrade.ru`
	defaultPort = 10051
	etherType   = 0x0800
)

type iphdr struct {
	vhl   uint8
	tos   uint8
	iplen uint16
	id    uint16
	off   uint16
	ttl   uint8
	proto uint8
	csum  uint16
	src   [4]byte
	dst   [4]byte
}

type sfpsla struct {
	id          uint8
	dst         [4]byte
	merkertime1 [7]byte
	merkertime2 [7]byte
	merkertime3 [7]byte
	number      uint32
}

func checksum(buf []byte) uint16 {
	sum := uint32(0)

	for ; len(buf) >= 2; buf = buf[2:] {
		sum += uint32(buf[0])<<8 | uint32(buf[1])
	}
	if len(buf) > 0 {
		sum += uint32(buf[0]) << 8
	}
	for sum > 0xffff {
		sum = (sum >> 16) + (sum & 0xffff)
	}
	csum := ^uint16(sum)
	/*
	 * From RFC 768:
	 * If the computed checksum is zero, it is transmitted as all ones (the
	 * equivalent in one's complement arithmetic). An all zero transmitted
	 * checksum value means that the transmitter generated no checksum (for
	 * debugging or for higher level protocols that don't care).
	 */
	if csum == 0 {
		csum = 0xffff
	}
	return csum
}

func (h *iphdr) checksum() {
	h.csum = 0
	var b bytes.Buffer
	binary.Write(&b, binary.BigEndian, h)
	h.csum = checksum(b.Bytes())
}

func main() {
	go zabbixHello("SFP-SLA_4401")
	// Find all devices
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	var net_name string
	// Print device information
	// fmt.Println("Devices found:")
	for _, device := range devices {
		//	  fmt.Println("\nName: ", device.Name)
		//	  fmt.Println("Description: ", device.Description)
		//	  fmt.Println("Devices addresses: ", device.Description)
		for _, address := range device.Addresses {
			//		  fmt.Println("- IP address: ", address.IP)
			if address.IP.Equal(net.ParseIP("10.0.10.115")) {
				net_name = device.Name
			}
			//		  fmt.Println("- Subnet mask: ", address.Netmask)
		}
	}
	// Open a raw socket on the specified interface, and configure it to accept

	ifi, err := net.InterfaceByName(net_name)
	if err != nil {
		log.Fatalf("failed to find interface %q: %v", net_name, err)
	}
	fmt.Println("Net_NAME: ", net_name)
	fmt.Println("interface: ", ifi.Name)
	c, err := raw.ListenPacket(ifi, etherType, nil)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	ipsrcstr := "10.0.10.115"
	ipdst_1sfpsla_str := "10.1.10.145"
	ipdst_2sfpsla_str := "10.1.10.140"

	ipsrc := net.ParseIP(ipsrcstr)
	ipdst1 := net.ParseIP(ipdst_1sfpsla_str)
	ipdst2 := net.ParseIP(ipdst_2sfpsla_str)

	// Default message to system's hostname if empty.
	ip := iphdr{
		vhl:   0x45,
		tos:   0,
		id:    0x0000, // the kernel overwrites id if it is zero
		off:   0,
		ttl:   0xFF,
		proto: 0x5E,
	}
	copy(ip.src[:], ipsrc.To4())
	copy(ip.dst[:], ipdst1.To4())
	sfpdat := sfpsla{
		id: 0xFA,
	}
	copy(sfpdat.dst[:], ipdst2.To4())
	ip.iplen = uint16(20 + 26)
	ip.checksum()

	//msg := make([]byte, ip.iplen)
	//ipd:=[]byte(fmt.Sprintf("%v",ip))
	//dats:=[]byte(fmt.Sprintf("%v",sfpdat))
	//msg=append(ipd,dats...)

	var bin_buf bytes.Buffer
	binary.Write(&bin_buf, binary.BigEndian, ip)
	binary.Write(&bin_buf, binary.BigEndian, sfpdat)

    msg := bin_buf.Bytes()
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
		//Destination: ethernet.Broadcast,
		Destination: []byte{0x64,0xD1,0x54,0x17,0xF6,0x82},
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
	fmt.Printf(" --== %x", b)
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
		if f.Payload[20] == 0xFA {
			fmt.Printf("\n\n--=Yes=--\n\r\r[%s] %s", addr.String(), string(f.Payload))
		} else {
			fmt.Printf("\n\n\r[%s] %v %x", addr.String(), len(f.Payload),f.Payload[:25])

		}
	}
}

func zabbixHello(host string) {
	var delay int
	for {
		delay = rand.Intn(1500)
		//delay:=i*100
		var metrics []*Metric
		metrics = append(metrics, NewMetric(host, "delay", fmt.Sprint(delay), time.Now().Unix()))

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
