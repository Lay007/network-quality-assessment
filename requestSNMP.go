package main

import (
	"fmt"
	//	"log"
	"net"
	"time"

	g "github.com/soniah/gosnmp"
)

func findSFP(c net.PacketConn, addr net.Addr,ip_server string, ip_1sfpsla_str string, ip_2sfpsla_str string, mac_src []byte, mac_dst []byte, test_type uint16, testWay int) bool {

	fmt.Println(" ==> TEST SFP way ==")

	SFP1_com := []int64{}
	SFP1_laz := []int64{}

	SFP2_com := []int64{}
	SFP2_laz := []int64{}

	oids := []string{".1.3.6.1.4.1.2010.1.13.0", ".1.3.6.1.4.1.2010.1.14.0"}

	counter := 0

	timer_SNMP := time.NewTicker(1000 * time.Millisecond)
	for range timer_SNMP.C {
		counter++
		if counter > 10 {
			timer_SNMP.Stop()
			break
		}
		g.Default.Target = ip_1sfpsla_str
		err := g.Default.Connect()
		if err != nil {
			fmt.Printf("Connect to SFP1 error: %v", err)
			return true
		}

		result, err2 := g.Default.Get(oids) // Get() accepts up to g.MAX_OIDS
		if err2 != nil {
			fmt.Printf("Get() err: %v", err2)
		}		
		for i, variable := range result.Variables {
			fmt.Printf("%d: oid: %s ", i, variable.Name)
			if variable.Name == ".1.3.6.1.4.1.2010.1.13.0" {
				fmt.Printf("SFP1 number: %v   Mb/s\n", float32(g.ToBigInt(variable.Value).Int64()*8)/1000000.0)
				SFP1_com = append(SFP1_com, g.ToBigInt(variable.Value).Int64())
			}
			if variable.Name == ".1.3.6.1.4.1.2010.1.14.0" {
				fmt.Printf("SFP1 number: %v   Mb/s\n", float32(g.ToBigInt(variable.Value).Int64()*8)/1000000.0)
				SFP1_laz = append(SFP1_laz, g.ToBigInt(variable.Value).Int64())
			}
		}
		g.Default.Conn.Close()

		g.Default.Target = ip_2sfpsla_str
		err = g.Default.Connect()
		if err != nil {
			fmt.Printf("Connect to SFP2 error: %v", err)
			return true
		}

		result, err2 = g.Default.Get(oids) // Get() accepts up to g.MAX_OIDS
		if err2 != nil {
			fmt.Printf("Get() err: %v", err2)
		}
		g.Default.Conn.Close()
		for i, variable := range result.Variables {
			fmt.Printf("%d: oid: %s ", i, variable.Name)
			if variable.Name == ".1.3.6.1.4.1.2010.1.13.0" {
				fmt.Printf("SFP2 number: %v   Mb/s\n", float32(g.ToBigInt(variable.Value).Int64()*8)/1000000.0)
				SFP2_com = append(SFP2_com, g.ToBigInt(variable.Value).Int64())
			}
			if variable.Name == ".1.3.6.1.4.1.2010.1.14.0" {
				fmt.Printf("SFP2 number: %v   Mb/s\n", float32(g.ToBigInt(variable.Value).Int64()*8)/1000000.0)
				SFP2_laz = append(SFP2_laz, g.ToBigInt(variable.Value).Int64())
			}
		}

		g.Default.Conn.Close()
	}

	fmt.Printf("->Generate<-")

	ipsrc := net.ParseIP(ip_server)
	ipdst1 := net.ParseIP(ip_1sfpsla_str)
	ipdst2 := net.ParseIP(ip_2sfpsla_str)

	b := packetForm(ipsrc, ipdst1, ipdst2, mac_src, mac_dst, 1024, 1, test_type, testWay)

	go func() {
		cc := 100000
		ticker := time.NewTicker(100 * time.Microsecond)
		for range ticker.C {
			cc--
			if cc > 10 {
				c.WriteTo(b, addr)
			} else {
				return
			}
		}
	}()
	counter = 0

	timer_SNMP = time.NewTicker(1000 * time.Millisecond)
	for range timer_SNMP.C {
		counter++
		if counter > 10 {
			timer_SNMP.Stop()
			break
		}
		g.Default.Target = ip_1sfpsla_str
		err := g.Default.Connect()
		if err != nil {
			fmt.Printf("Connect to SFP1 error: %v", err)
			return true
		}

		result, err2 := g.Default.Get(oids) // Get() accepts up to g.MAX_OIDS
		if err2 != nil {
			fmt.Printf("Get() err: %v", err2)
		}
		g.Default.Conn.Close()
		for i, variable := range result.Variables {
			fmt.Printf("%d: oid: %s ", i, variable.Name)
			if variable.Name == ".1.3.6.1.4.1.2010.1.13.0" {
				fmt.Printf("SFP1 number: %v   Mb/s\n", float32(g.ToBigInt(variable.Value).Int64()*8)/1000000.0)
				SFP1_com = append(SFP1_com, g.ToBigInt(variable.Value).Int64())
			}
			if variable.Name == ".1.3.6.1.4.1.2010.1.14.0" {
				fmt.Printf("SFP1 number: %v   Mb/s\n", float32(g.ToBigInt(variable.Value).Int64()*8)/1000000.0)
				SFP1_laz = append(SFP1_laz, g.ToBigInt(variable.Value).Int64())
			}
		}

		g.Default.Target = ip_2sfpsla_str
		err = g.Default.Connect()
		if err != nil {
			fmt.Printf("Connect to SFP2 error: %v", err)
			return true
		}

		result, err2 = g.Default.Get(oids) // Get() accepts up to g.MAX_OIDS
		if err2 != nil {
			fmt.Printf("Get() err: %v", err2)
		}
		g.Default.Conn.Close()
		for i, variable := range result.Variables {
			fmt.Printf("%d: oid: %s ", i, variable.Name)
			if variable.Name == ".1.3.6.1.4.1.2010.1.13.0" {
				fmt.Printf("SFP2 number: %v   Mb/s\n", float32(g.ToBigInt(variable.Value).Int64()*8)/1000000.0)
				SFP2_com = append(SFP2_com, g.ToBigInt(variable.Value).Int64())
			}
			if variable.Name == ".1.3.6.1.4.1.2010.1.14.0" {
				fmt.Printf("SFP2 number: %v   Mb/s\n", float32(g.ToBigInt(variable.Value).Int64()*8)/1000000.0)
				SFP2_laz = append(SFP2_laz, g.ToBigInt(variable.Value).Int64())
			}
		}

		g.Default.Conn.Close()
	}

	return true
}
