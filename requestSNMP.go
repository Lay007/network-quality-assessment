package main

import (
	"fmt"
	//	"log"
	. "./go-zabbix"
	"database/sql"
	g "github.com/soniah/gosnmp"
	"net"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var mux_SNMP sync.Mutex
var db_SNMP *sql.DB

func findSFP(c net.PacketConn, addr net.Addr, ip_server string, ip_1sfpsla_str string, ip_2sfpsla_str string, mac_src []byte, mac_dst []byte, test_type uint16, testWay int) int {

	fmt.Println(" ==> TEST SFP way ==")

	var SFP1_com_init, SFP1_laz_init, SFP2_com_init, SFP2_laz_init int64

	oids := []string{".1.3.6.1.4.1.2010.1.13.0", ".1.3.6.1.4.1.2010.1.14.0"}

	counter := 0

	timer_SNMP := time.NewTicker(1000 * time.Millisecond)
	for range timer_SNMP.C {
		counter++
		if counter > 10 {
			timer_SNMP.Stop()
			break
		}
		mux_SNMP.Lock()
		g.Default.Target = ip_1sfpsla_str
		err := g.Default.Connect()
		if err != nil {
			fmt.Printf("Connect to SFP1 error: %v", err)
			mux_SNMP.Unlock()
			return 0
		}

		result, err2 := g.Default.Get(oids) // Get() accepts up to g.MAX_OIDS
		if err2 != nil {
			fmt.Printf("Get() err: %v", err2)
			mux_SNMP.Unlock()
		}
		for _, variable := range result.Variables {
			//	fmt.Printf("%d: oid: %s ", i, variable.Name)
			if variable.Name == ".1.3.6.1.4.1.2010.1.13.0" {
				//	fmt.Printf("SFP1 number: %v   Mb/s\n", float32(g.ToBigInt(variable.Value).Int64()*8)/1000000.0)
				SFP1_com_init = SFP1_com_init + g.ToBigInt(variable.Value).Int64()
			}
			if variable.Name == ".1.3.6.1.4.1.2010.1.14.0" {
				//	fmt.Printf("SFP1 number: %v   Mb/s\n", float32(g.ToBigInt(variable.Value).Int64()*8)/1000000.0)
				SFP1_laz_init = SFP1_laz_init + g.ToBigInt(variable.Value).Int64()
			}
		}
		g.Default.Conn.Close()

		g.Default.Target = ip_2sfpsla_str
		err = g.Default.Connect()
		if err != nil {
			fmt.Printf("Connect to SFP2 error: %v", err)
			mux_SNMP.Unlock()
			return 0
		}

		result, err2 = g.Default.Get(oids) // Get() accepts up to g.MAX_OIDS
		if err2 != nil {
			fmt.Printf("Get() err: %v", err2)
			mux_SNMP.Unlock()
		}
		g.Default.Conn.Close()
		mux_SNMP.Unlock()
		for _, variable := range result.Variables {
			//	fmt.Printf("%d: oid: %s ", i, variable.Name)
			if variable.Name == ".1.3.6.1.4.1.2010.1.13.0" {
				//	fmt.Printf("SFP2 number: %v   Mb/s\n", float32(g.ToBigInt(variable.Value).Int64()*8)/1000000.0)
				SFP2_com_init = SFP2_com_init + g.ToBigInt(variable.Value).Int64()
			}
			if variable.Name == ".1.3.6.1.4.1.2010.1.14.0" {
				//	fmt.Printf("SFP2 number: %v   Mb/s\n", float32(g.ToBigInt(variable.Value).Int64()*8)/1000000.0)
				SFP2_laz_init = SFP2_laz_init + g.ToBigInt(variable.Value).Int64()
			}
		}

	}

	fmt.Printf("\n  SFP1_com - %v SFP1_laz - %v ", float32(SFP1_com_init*8)/10000000.0, float32(SFP1_laz_init*8)/10000000.0)
	fmt.Printf("\n  SFP2_com - %v SFP2_laz - %v ", float32(SFP2_com_init*8)/10000000.0, float32(SFP2_laz_init*8)/10000000.0)

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
			if cc > 0 {
				c.WriteTo(b, addr)
			} else {
				return
			}
		}
	}()

	var SFP1_com_load, SFP1_laz_load, SFP2_com_load, SFP2_laz_load int64

	counter = 0

	timer_SNMP = time.NewTicker(1000 * time.Millisecond)
	for range timer_SNMP.C {
		counter++
		if counter > 10 {
			timer_SNMP.Stop()
			break
		}
		mux_SNMP.Lock()
		g.Default.Target = ip_1sfpsla_str
		err := g.Default.Connect()
		if err != nil {
			fmt.Printf("Connect to SFP1 error: %v", err)
			mux_SNMP.Unlock()
			return 0
		}

		result, err2 := g.Default.Get(oids) // Get() accepts up to g.MAX_OIDS
		if err2 != nil {
			fmt.Printf("Get() err: %v", err2)
			mux_SNMP.Unlock()
		}
		g.Default.Conn.Close()
		mux_SNMP.Unlock()
		for _, variable := range result.Variables {
			//	fmt.Printf("%d: oid: %s ", i, variable.Name)
			if variable.Name == ".1.3.6.1.4.1.2010.1.13.0" {
				//	fmt.Printf("SFP1 number: %v   Mb/s\n", float32(g.ToBigInt(variable.Value).Int64()*8)/1000000.0)
				SFP1_com_load = SFP1_com_load + g.ToBigInt(variable.Value).Int64()
			}
			if variable.Name == ".1.3.6.1.4.1.2010.1.14.0" {
				//	fmt.Printf("SFP1 number: %v   Mb/s\n", float32(g.ToBigInt(variable.Value).Int64()*8)/1000000.0)
				SFP1_laz_load = SFP1_laz_load + g.ToBigInt(variable.Value).Int64()
			}
		}

		mux_SNMP.Lock()
		g.Default.Target = ip_2sfpsla_str
		err = g.Default.Connect()
		if err != nil {
			fmt.Printf("Connect to SFP2 error: %v", err)
			mux_SNMP.Unlock()
			return 0
		}

		result, err2 = g.Default.Get(oids) // Get() accepts up to g.MAX_OIDS
		if err2 != nil {
			fmt.Printf("Get() err: %v", err2)
			mux_SNMP.Unlock()
		}
		g.Default.Conn.Close()
		mux_SNMP.Unlock()
		for _, variable := range result.Variables {
			//	fmt.Printf("%d: oid: %s ", i, variable.Name)
			if variable.Name == ".1.3.6.1.4.1.2010.1.13.0" {
				//	fmt.Printf("SFP2 number: %v   Mb/s\n", float32(g.ToBigInt(variable.Value).Int64()*8)/1000000.0)
				SFP2_com_load = SFP2_com_load + g.ToBigInt(variable.Value).Int64()
			}
			if variable.Name == ".1.3.6.1.4.1.2010.1.14.0" {
				//	fmt.Printf("SFP2 number: %v   Mb/s\n", float32(g.ToBigInt(variable.Value).Int64()*8)/1000000.0)
				SFP2_laz_load = SFP2_laz_load + g.ToBigInt(variable.Value).Int64()
			}
		}

	}

	SFP1_com_load = SFP1_com_load - SFP1_com_init
	SFP1_laz_load = SFP1_laz_load - SFP1_laz_init
	SFP2_com_load = SFP2_com_load - SFP2_com_init
	SFP2_laz_load = SFP2_laz_load - SFP2_laz_init

	fmt.Printf("\n  SFP1_com - %v SFP1_laz - %v ", float32(SFP1_com_load*8)/10000000.0, float32(SFP1_laz_load*8)/10000000.0)
	fmt.Printf("\n  SFP2_com - %v SFP2_laz - %v \n", float32(SFP2_com_load*8)/10000000.0, float32(SFP2_laz_load*8)/10000000.0)

	fmt.Printf("\n  SFP1_com - %v SFP1_laz - %v ", SFP1_com_load, SFP1_laz_load)
	fmt.Printf("\n  SFP2_com - %v SFP2_laz - %v \n", SFP2_com_load, SFP2_laz_load)

	if (SFP1_com_load > 40000000) && (SFP1_laz_load > 40000000) && ((float32(SFP1_laz_load)/float32(SFP1_com_load) > 1.7) || (float32(SFP1_com_load)/float32(SFP1_laz_load) > 1.7)) {
		return 2
	} else {
		if (SFP2_com_load > 40000000) && (SFP2_laz_load > 40000000) && ((float32(SFP2_laz_load)/float32(SFP2_com_load) > 1.7) || (float32(SFP2_com_load)/float32(SFP2_laz_load) > 1.7)) {

			return 2
		} else {
			return 1
		}
	}
}

func (module *module_sfp) startSNMP(conf global_config) {

	//fmt.Println("Start SNMP_check - ", (*module).address_ip)
	if conf.net_interface_name == "" {
		fmt.Println("null net_interface_name")
		return
	}
	envTarget := (*module).address_ip
	envPort := "161"

	port, _ := strconv.ParseUint(envPort, 10, 16)

	// Build our own GoSNMP struct, rather than using g.Default.
	// Do verbose logging of packets.
	params := &g.GoSNMP{
		Target:    envTarget,
		Port:      uint16(port),
		Community: "public",
		Version:   g.Version2c,
		Timeout:   time.Millisecond * 100,
		Retries:   1,
	}

	var SFP_laz, SFP_com int64
	oids := []string{".1.3.6.1.4.1.2010.1.13.0", ".1.3.6.1.4.1.2010.1.14.0"}

	timer_SNMP := time.NewTicker(1000 * time.Millisecond)
	//	timer_SNMP.C <- time.Now()

	for range timer_SNMP.C {
		select {
		case <-(*module).chan_stop:
			//fmt.Println("End SNMP_check - ", (*module).address_ip)
			timer_SNMP.Stop()
			return
		default:
			{
				mux_SNMP.Lock()
				//	g.Default.Retries=1
				//	g.Default.Target = (*module).address_ip
				//	err := g.Default.Connect()
				err := params.Connect()
				if err != nil {
					fmt.Print("Module : ", envTarget)
					fmt.Println("Connect() err: ", err)
					mux_SNMP.Unlock()
					continue
				}
				defer params.Conn.Close()
				result, err2 := params.Get(oids) // Get() accepts up to g.MAX_OIDS
				if err2 != nil {
					fmt.Print("Module : ", envTarget)
					fmt.Println("Get() err: ", err2)
					params.Conn.Close()
					mux_SNMP.Unlock()
					continue
				}
				var metrics []*Metric
				for _, variable := range result.Variables {
					//	fmt.Printf("%d: oid: %s ", i, variable.Name)

					if variable.Name == ".1.3.6.1.4.1.2010.1.13.0" {

						SFP_com = g.ToBigInt(variable.Value).Int64() * 8
						metrics = append(metrics, NewMetric((*module).zabbix_node, "band_to_lazer", fmt.Sprint(SFP_com), time.Now().Unix()))
					}
					if variable.Name == ".1.3.6.1.4.1.2010.1.14.0" {
						//	fmt.Printf("SFP1 number: %v   Mb/s\n", float32(g.ToBigInt(variable.Value).Int64()*8)/1000000.0)

						SFP_laz = g.ToBigInt(variable.Value).Int64() * 8
						metrics = append(metrics, NewMetric((*module).zabbix_node, "band_to_comm", fmt.Sprint(SFP_laz), time.Now().Unix()))
					}
				}
				//params.Conn.Close()
				//g.Default.Conn.Close()
				mux_SNMP.Unlock()
				if (*module).zabbix_node != "" {
					packet := NewPacket(metrics)
					// Send packet to zabbix
					z := NewSender(conf.zabbix_server_name, conf.zabbix_server_port)
					z.Send(packet)
				}
				/*
					db_SNMP, err := sql.Open("mysql", db_user+":"+db_user_pass+"@/"+db_database)
					if err != nil {
						db_SNMP.Close()
						mux_SNMP.Unlock()
						fmt.Println(" -!! Error !!-")
						fmt.Println(err)
						fmt.Println(" ----=====----")
						return
					}*/
				//rez,er:=
				db_SNMP.Exec("INSERT INTO modules_sfp_sla_load_rez (module_id, datatime, load_to_lazer, load_to_com) VALUES(?, NOW(), ?, ?)", (*module).id, SFP_laz, SFP_com)
				//fmt.Println("Rez add module metric: ",rez)
				//fmt.Println("Error add module metric: ",er)
				//	fmt.Println(" add metric - ", time.Now())
				//db_SNMP.Close()

			}
		}
	}
}
