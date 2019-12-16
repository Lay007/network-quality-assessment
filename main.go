package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"time"

	. "./go-zabbix"
	"github.com/google/gopacket/pcap"

	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/raw"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const (
	db_user      = `sfp_user`
	db_user_pass = `rootsfp`
	db_database  = `server_sfp_sla`

	debugV = true

	defaultHost = `localhost`
	//defaultHost = `remote.fibertrade.ru`
	defaultPort = 10051
	etherType   = 0x0800

	ipsrcstr          = "10.0.10.115"
	ipdst_1sfpsla_str = "10.0.10.172"
	ipdst_2sfpsla_str = "10.0.10.175"
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

type global_config struct {
	server_ip          string
	zabbix_server_name string
	vlan               int
	vlan_number        int
}
type module_sfp struct {
	id          int
	name        string
	address_mac string
	address_ip  string
	version     string
	location    string
}

var numberTx, numberCounter uint32

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
	//time.Sleep(100 * time.Second)

	db, err := sql.Open("mysql", db_user+":"+db_user_pass+"@/"+db_database)
	if err != nil {
		panic(err)
	}

	db.Exec("DELETE FROM net_interfaces_from_server_sla")
	db.Exec("ALTER TABLE net_interfaces_from_server_sla AUTO_INCREMENT = 1")

	devices, err := pcap.FindAllDevs()
	fmt.Println(err)
	if err != nil {
		log.Fatal(err)
	}

	var net_name string

	for _, device := range devices {
		fmt.Println(device.Name)
		netInterface, err := net.InterfaceByName(device.Name)
		var addressMac net.HardwareAddr
		if err == nil {
			addressMac = netInterface.HardwareAddr
		}
		net_name = device.Name
		for _, address := range device.Addresses {
			db.Exec("INSERT INTO net_interfaces_from_server_sla (name, address_IP, address_mac) VALUES(?, ?, ?)", device.Name, address.IP.String(), addressMac.String())
		}
	}

	t := time.NewTicker(30 * time.Second)
	for range t.C {
		//	row_test_real, err := db.Query("select * from global_config where status=1")

	}
	// -------=======

	row, err := db.Query("select * from global_config")
	if err != nil {
		panic(err)
	}
	defer row.Close()
	row.Next()
	conf := new(global_config)
	err = row.Scan(&conf.server_ip, &conf.zabbix_server_name, &conf.vlan, &conf.vlan_number)
	if err != nil {
		fmt.Println(err)
	}

	rows, err := db.Query("select * from modules_sfp_sla")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	modules := []module_sfp{}

	for rows.Next() {
		m := module_sfp{}
		err = rows.Scan(&m.id, &m.name, &m.address_mac, &m.address_ip, &m.version, &m.location)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(m.address_ip)
		modules = append(modules, m)
	}

	defer db.Close()

	//go zabbixHello("SFP-SLA_4401")

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

	// Send messages in one goroutine, receive messages in another.
	go sendMessages(c, ifi.HardwareAddr)
	go receiveMessages(c, ifi.MTU)

	// Block forever.
	select {}
}

// sendMessages continuously sends a message over a connection at regular intervals,
// sourced from specified hardware address.
func sendMessages(c net.PacketConn, source net.HardwareAddr) {

	t := time.NewTicker(1 * time.Second)
	for range t.C {

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
			id: 0xFC,
		}
		copy(sfpdat.dst[:], ipdst2.To4())
		numberTx++
		sfpdat.number = numberTx
		ip.iplen = uint16(20 + 26 + 4)
		ip.checksum()

		var bin_buf bytes.Buffer
		binary.Write(&bin_buf, binary.BigEndian, ip)
		binary.Write(&bin_buf, binary.BigEndian, sfpdat)

		msg := bin_buf.Bytes()
		f := &ethernet.Frame{
			//Destination: ethernet.Broadcast,
			Destination: []byte{0x5A, 0x11, 0x22, 0x33, 0x44, 0x00},
			//Destination: []byte{0x64, 0xD1, 0x54, 0x17, 0xF6, 0x82},
			Source:    source,
			EtherType: 0x0800,
			Payload:   []byte(msg),
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

		fmt.Printf("raw:  %x \n", b)
		fmt.Println(" --== Packet send ==--")
		fmt.Printf("mac dst  %x \n", b[0:6])
		fmt.Printf("mac src  %x \n", b[6:12])
		fmt.Printf("type eth %x \n", b[12:14])
		fmt.Printf("size     %v \n", b[16:18])

		fmt.Printf("ip sourse %v.%v.%v.%v \n", b[26], b[27], b[28], b[29])
		fmt.Printf("ip dst    %v.%v.%v.%v \n", b[30], b[31], b[32], b[33])
		fmt.Println(" --== End Packet ==--")

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
		fmt.Printf("\n\n--=Test %x - \n", f.Payload[12:16])
		var ips [4]byte
		copy(ips[:], (net.ParseIP(ipdst_1sfpsla_str)).To4())
		fmt.Printf("\n\n--=T_so %x - \n", ips)

		// Display source of message and message itself.
		if (f.Payload[20] == 0xFC) && (bytes.Equal(f.Payload[12:16], ips[:]) == true) {
			numberCounter++
			fmt.Printf("\n\n--=Packet DETECT!!!=--\n")
			//fmt.Printf("\n\n--=Test %x - \n -== %x\n",f.Payload[12:15],net.ParseIP(ipsrcstr))
			fmt.Printf("size: %v raw:  %x \n", len(f.Payload), f.Payload)
			fmt.Printf("\n\rEthernet source: [%s]\n", addr.String())

			fmt.Printf("size     %x \n", b[2:4])

			fmt.Printf("ip sourse %v.%v.%v.%v \n", f.Payload[12], f.Payload[13], f.Payload[14], f.Payload[15])
			fmt.Printf("ip dst    %v.%v.%v.%v \n", f.Payload[16], f.Payload[17], f.Payload[18], f.Payload[19])

			fmt.Printf("ip SFP2   %v.%v.%v.%v \n", f.Payload[21], f.Payload[22], f.Payload[23], f.Payload[24])

			fmt.Printf("time marker_SFP1_1 :   %x \n", f.Payload[25:32])
			fmt.Printf("time marker_SFP2   :   %x \n", f.Payload[32:39])
			fmt.Printf("time marker_SFP1_2 :   %x \n", f.Payload[39:46])
			fmt.Printf("Number marker      :   %x \n", f.Payload[46:50])
			fmt.Println(" --== End Packet ==--")

			var markerSFP11, markerSFP12, markerSFP2 int64
			var ind uint

			for ind = 0; ind < 7; ind++ {
				markerSFP11 = markerSFP11 + int64(f.Payload[31-ind])<<(8*ind)
				markerSFP2 = markerSFP2 + int64(f.Payload[38-ind])<<(8*ind)
				markerSFP12 = markerSFP12 + int64(f.Payload[45-ind])<<(8*ind)
			}

			var numberR uint32
			for ind = 0; ind < 4; ind++ {
				numberR += uint32(f.Payload[49-ind]) << (8 * ind)
			}

			zabbix_delay("SFP-SLA_4401", markerSFP12-markerSFP11)
			zabbix_jitter("SFP-SLA_4401", getJitter(markerSFP12-markerSFP11))
			zabbix_error("SFP-SLA_4401", float32(numberR-numberCounter)/float32(numberR))

		} else {
			fmt.Printf("\n\n\r[%s] %v %x", addr.String(), len(f.Payload), f.Payload[:25])

		}
	}
}

var mass_solve []int64

func getJitter(in_solve int64) float32 {
	var jitter, mean float32
	var size_s int
	var max, min int64

	size_s = 100
	mass_solve = append(mass_solve, in_solve)
	if len(mass_solve) < (size_s + 1) {
		return 0
	}
	mass_solve = mass_solve[1:(size_s + 1)]

	max = mass_solve[0]
	min = max
	mean = float32(mass_solve[0]) / float32(size_s)

	for ind := 1; ind < size_s; ind++ {
		if max < mass_solve[ind] {
			max = mass_solve[ind]
		}
		if min > mass_solve[ind] {
			min = mass_solve[ind]
		}
		mean = mean + (float32(mass_solve[ind]) / float32(size_s))
	}
	if (float32(max) - mean) > (mean - float32(min)) {
		jitter = float32(max) - mean
	} else {
		jitter = mean - float32(min)
	}

	fmt.Printf(" --== Jitter debug ==-- \n")
	fmt.Printf(" --== Slice: %x \n", mass_solve)
	fmt.Printf(" --== Max = %x \n", max)
	fmt.Printf(" --== Min = %x \n", min)
	fmt.Printf(" --== Mean = %f \n", mean)
	fmt.Printf(" --== Jitter = %f \n", jitter)
	fmt.Printf(" --== End Jitter debug ==-- \n")

	return jitter
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

func zabbix_delay(host string, delay int64) {

	//delay = delay * 8 // [mks] 125 MGz - clock, => T = 8 mks
	//delay = int64( float64(delay) * 1000000 / (math.Pow(2, 32))) // [mks] 125 MGz - clock, => T = 8 mks
	var delfloat float32
	delfloat = float32(delay) * 1000000 / float32(math.Pow(2, 32))
	var metrics []*Metric
	metrics = append(metrics, NewMetric(host, "delay", fmt.Sprint(delfloat), time.Now().Unix()))

	// Create instance of Packet class
	packet := NewPacket(metrics)
	//fmt.Println(packet);
	// Send packet to zabbix
	z := NewSender(defaultHost, defaultPort)
	z.Send(packet)

}

func zabbix_jitter(host string, jitter float32) {

	//delay = delay * 8 // [mks] 125 MGz - clock, => T = 8 mks
	if jitter != 0 {
		jitter = jitter * 1000000 / float32(math.Pow(2, 32)) // [mks] 125 MGz - clock, => T = 8 mks
	}
	var metrics []*Metric
	metrics = append(metrics, NewMetric(host, "jitter", fmt.Sprint(jitter), time.Now().Unix()))

	// Create instance of Packet class
	packet := NewPacket(metrics)
	//fmt.Println(packet);
	// Send packet to zabbix
	z := NewSender(defaultHost, defaultPort)
	z.Send(packet)

}

func zabbix_error(host string, err float32) {

	var metrics []*Metric
	metrics = append(metrics, NewMetric(host, "error_probability", fmt.Sprint(err), time.Now().Unix()))

	// Create instance of Packet class
	packet := NewPacket(metrics)
	//fmt.Println(packet);
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
