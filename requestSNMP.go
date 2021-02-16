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

	"github.com/go-ping/ping"
	"github.com/tatsushid/go-fastping"
)

func testPing(ip string) int {
	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", ip)
	if err != nil {
		fmt.Println(err)
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
	}
	p.OnIdle = func() {
		fmt.Println("finish")
	}
	err = p.Run()
	if err != nil {
		fmt.Println(err)
	}
	return 0
}

func testPing_old(ip string) int {
	fmt.Println(" Start Ping")
	pinger, err := ping.NewPinger(ip)
	fmt.Println(" Ping: ", pinger)
	if err != nil {
		fmt.Println(" Error new ping: ", err)
		return 1
	}
	pinger.Count = 3
	pinger.Timeout = time.Second * 5
	er := pinger.Run()
	if er != nil {
		fmt.Println(" Error ping run: ", err)
		return 1
	}
	stats := pinger.Statistics()
	fmt.Println(" Ping st: ", stats)
	if stats.PacketsRecv == 0 {
		fmt.Println(" Ping error: ")
		return 1
	}
	return 0
}

var mux_SNMP sync.Mutex
var db_SNMP *sql.DB

func (testWay *testWaySFP) getResult(thr int64, revers bool) (float32, float32, float32, float32) {
	//var SFP1_com, SFP1_laz, SFP2_com, SFP2_laz int64
	period_nano_gen := int64(0)
	count_gen := int64(0)
	if thr > 0 {
		period_nano_gen = int64(int64(testWay.packet_size) * 8 * 1000 / thr)
		count_gen = 10000000000 / period_nano_gen
	}
	b := packetForm(testWay.ipsrc, testWay.ipdst1, testWay.ipdst2, testWay.mac_src, testWay.mac_dst, int(testWay.packet_size), 1, 0, testWay.test_type)
	if revers {
		b = packetForm(testWay.ipsrc, testWay.ipdst2, testWay.ipdst1, testWay.mac_src, testWay.mac_dst2, int(testWay.packet_size), 1, 0, testWay.test_type)
	}
	if thr > 0 {

		go func() {
			start := time.Now()
			cc := count_gen
			ticker := time.NewTicker(time.Duration(period_nano_gen) * time.Nanosecond)
			for range ticker.C {
				cc--
				if cc%10000 == 0 {
					if time.Since(start) >= (10 * time.Second) {
						return
					}
				}
				if cc > 0 {
					(*testWay).conn.WriteTo(b, testWay.addr)
				} else {
					return
				}
			}
		}()

	}
	var SFP1_com_load, SFP1_laz_load, SFP2_com_load, SFP2_laz_load int64
	var SFP1_com_load_count, SFP1_laz_load_count, SFP2_com_load_count, SFP2_laz_load_count int64
	//timer_SNMP = time.NewTicker(500 * time.Millisecond)
	oids := []string{".1.3.6.1.4.1.2010.1.13.0", ".1.3.6.1.4.1.2010.1.14.0"}

	var counter int
	start := time.Now()

	time.Sleep(2000 * time.Millisecond)

	g.Default.Timeout = time.Millisecond * 100
	g.Default.Retries = 1

	for {

		if time.Since(start) > (time.Duration(testWay.period_gen_ms) * time.Millisecond) {
			break
		}
		mux_SNMP.Lock()
		g.Default.Target = testWay.ipdst1.String()
		err := g.Default.Connect()
		if err != nil {
			//	fmt.Printf("\nConnect to SFP1 error: %v", err)
			mux_SNMP.Unlock()
		} else {

			result, err2 := g.Default.Get(oids) // Get() accepts up to g.MAX_OIDS
			if err2 != nil {
				//		fmt.Printf("\nGet() err: %v", err2)
				mux_SNMP.Unlock()
			} else {
				g.Default.Conn.Close()
				mux_SNMP.Unlock()
				for _, variable := range result.Variables {
					//	fmt.Printf("%d: oid: %s ", i, variable.Name)
					if variable.Name == ".1.3.6.1.4.1.2010.1.13.0" {
						//	fmt.Printf("\nSFP1 number: %v   Mb/s", float32(g.ToBigInt(variable.Value).Int64()*8)/1000000.0)
						if g.ToBigInt(variable.Value).Int64() > int64(testWay.SFP1_com_min*1000000)/8 {
							SFP1_com_load = SFP1_com_load + g.ToBigInt(variable.Value).Int64()
							SFP1_com_load_count++
						}
					}
					if variable.Name == ".1.3.6.1.4.1.2010.1.14.0" {
						//	fmt.Printf("\nSFP1 number: %v   Mb/s", float32(g.ToBigInt(variable.Value).Int64()*8)/1000000.0)
						if g.ToBigInt(variable.Value).Int64() > int64(testWay.SFP1_laz_min*1000000)/8 {
							SFP1_laz_load = SFP1_laz_load + g.ToBigInt(variable.Value).Int64()
							SFP1_laz_load_count++
						}
					}
				}
			}
		}
		mux_SNMP.Lock()
		g.Default.Target = testWay.ipdst2.String()
		err = g.Default.Connect()
		if err != nil {
			//	fmt.Printf("\nConnect to SFP2 error: %v", err)
			mux_SNMP.Unlock()
		} else {
			result, err2 := g.Default.Get(oids) // Get() accepts up to g.MAX_OIDS
			if err2 != nil {
				//		fmt.Printf("\nGet() err: %v", err2)
				mux_SNMP.Unlock()
			} else {
				g.Default.Conn.Close()
				mux_SNMP.Unlock()
				for _, variable := range result.Variables {
					//	fmt.Printf("%d: oid: %s ", i, variable.Name)
					if variable.Name == ".1.3.6.1.4.1.2010.1.13.0" {
						//	fmt.Printf("\nSFP2 number: %v   Mb/s", float32(g.ToBigInt(variable.Value).Int64()*8)/1000000.0)
						if g.ToBigInt(variable.Value).Int64() > int64(testWay.SFP2_com_min*1000000)/8 {
							SFP2_com_load = SFP2_com_load + g.ToBigInt(variable.Value).Int64()
							SFP2_com_load_count++
						}
					}
					if variable.Name == ".1.3.6.1.4.1.2010.1.14.0" {
						//	fmt.Printf("\nSFP2 number: %v   Mb/s\n", float32(g.ToBigInt(variable.Value).Int64()*8)/1000000.0)
						if g.ToBigInt(variable.Value).Int64() > int64(testWay.SFP2_laz_min*1000000)/8 {
							SFP2_laz_load = SFP2_laz_load + g.ToBigInt(variable.Value).Int64()
							SFP2_laz_load_count++
						}
					}
				}
			}
		}
		counter++
		time.Sleep(500 * time.Millisecond)
		//	fmt.Println()
	}
	time.Sleep(time.Duration(testWay.pause_ms) * time.Millisecond)
	if SFP1_com_load_count > 0 {
		SFP1_com_load = SFP1_com_load / SFP1_com_load_count
	}
	if SFP1_laz_load_count > 0 {
		SFP1_laz_load = SFP1_laz_load / SFP1_laz_load_count
	}
	if SFP2_com_load_count > 0 {
		SFP2_com_load = SFP2_com_load / SFP2_com_load_count
	}
	if SFP2_laz_load_count > 0 {
		SFP2_laz_load = SFP2_laz_load / SFP2_laz_load_count
	}
	return float32(SFP1_com_load*8) / 1000000.0, float32(SFP1_laz_load*8) / 1000000.0, float32(SFP2_com_load*8) / 1000000.0, float32(SFP2_laz_load*8) / 1000000.0
}

func findSFP(c net.PacketConn, addr net.Addr, ip_server string, ip_1sfpsla_str string, ip_2sfpsla_str string, mac_src []byte, mac_dst []byte, mac_dst2 []byte, test_type uint16, test_Way int, packet_size int, period_nano_gen_s int64) int {

	//  Последовательное соединение: [сервер]--(модуль SFP_SLA 1)--(модуль SFP_SLA 2)
	//  Соединенеие звездой: (модуль SFP_SLA 1)--[сервер]--(модуль SFP_SLA 2)
	//  Результат выполнения: XYYY
	//  X:
	//     0 - Последоватльное соединение. Расположение правильное
	//     1 - Последовательное соенинение. Расположение не правильное. Меняем местами
	//     2 - Соединенеие звездой. Нагрузка одинаковая.
	//     3 - Соединенеие звездой. Нагрузка неравномерная. Расположение правильное
	//     4 - Соединенеие звездой. Нагрузка неравномерная. Расположение не правильное. Меняем местами
	//  Y:
	//	   010 - 10 Мб/с
	//     100 - 100 Мб/с
	//     999 - 1 Гб/с

	step1 := int64(200) //200
	step2 := int64(70)  //70
	step3 := int64(7)   //7

	period_gen_ms := int64(10000) // период генерации 10 сек
	pause_ms := int64(2000)       // период генерации 2 сек

	ipsrc := net.ParseIP(ip_server)
	ipdst1 := net.ParseIP(ip_1sfpsla_str)
	ipdst2 := net.ParseIP(ip_2sfpsla_str)

	var testWay testWaySFP

	testWay.test_type = test_Way
	testWay.ipsrc = ipsrc
	testWay.ipdst1 = ipdst1
	testWay.ipdst2 = ipdst2
	testWay.packet_size = int(1280)
	testWay.test_type = test_Way
	testWay.addr = addr

	testWay.period_gen_ms = period_gen_ms
	testWay.pause_ms = pause_ms

	testWay.mac_src = mac_src
	testWay.mac_dst = mac_dst
	testWay.mac_dst2 = mac_dst2

	testWay.conn = c

	revers := false

	fmt.Println(" ==> TEST SFP way ==")
	SFP1_com_init, SFP1_laz_init, SFP2_com_init, SFP2_laz_init := testWay.getResult(0, revers)

	fmt.Printf("\n  SFP1_com - %v   SFP1_laz - %v ", SFP1_com_init, SFP1_laz_init)
	fmt.Printf("\n  SFP2_com - %v   SFP2_laz - %v \n", SFP2_com_init, SFP2_laz_init)

	testWay.SFP1_com_min = SFP1_com_init * 12 / 10
	testWay.SFP1_laz_min = SFP1_laz_init * 12 / 10

	testWay.SFP2_com_min = SFP2_com_init * 12 / 10
	testWay.SFP2_laz_min = SFP2_laz_init * 12 / 10

	SFP1_com, SFP1_laz, SFP2_com, SFP2_laz := testWay.getResult(step1, revers)

	fmt.Printf("\n\n  SFP1_com - %v   SFP1_laz - %v ", SFP1_com, SFP1_laz)
	fmt.Printf("\n  SFP2_com - %v   SFP2_laz - %v \n", SFP2_com, SFP2_laz)

	SFP1_com -= SFP1_com_init
	SFP1_laz -= SFP1_laz_init
	SFP2_com -= SFP2_com_init
	SFP2_laz -= SFP2_laz_init

	revers = true
	SFP1_com_rev, SFP1_laz_rev, SFP2_com_rev, SFP2_laz_rev := testWay.getResult(step1, revers)

	fmt.Printf("\n\n  SFP1_com - %v   SFP1_laz - %v ", SFP1_com_rev, SFP1_laz_rev)
	fmt.Printf("\n  SFP2_com - %v   SFP2_laz - %v \n", SFP2_com_rev, SFP2_laz_rev)

	SFP1_com_rev -= SFP1_com_init
	SFP1_laz_rev -= SFP1_laz_init
	SFP2_com_rev -= SFP2_com_init
	SFP2_laz_rev -= SFP2_laz_init

	if (((SFP1_com > 0.8*float32(step1)) || (SFP1_laz > 0.8*float32(step1))) || // throuth OK
		((SFP2_com > 0.8*float32(step1)) || (SFP2_laz > 0.8*float32(step1)))) && (((SFP1_com_rev > 0.8*float32(step1)) || (SFP1_laz_rev > 0.8*float32(step1))) ||
		((SFP2_com_rev > 0.8*float32(step1)) || (SFP2_laz_rev > 0.8*float32(step1)))) {

		if (((SFP1_com > 0.8*float32(step1) && SFP1_com < 1.2*float32(step1)) && (SFP1_laz > 0.8*float32(step1) && SFP1_laz < 1.2*float32(step1))) ||
			((SFP2_com > 0.8*float32(step1) && SFP2_com < 1.2*float32(step1)) && (SFP2_laz > 0.8*float32(step1) && SFP2_laz < 1.2*float32(step1)))) &&
			((((SFP1_laz_rev > 0) && (SFP1_com_rev/(SFP1_laz_rev) > 1.7)) || ((SFP1_com_rev > 0) && (SFP1_laz_rev/(SFP1_com_rev) > 1.7))) ||
				(((SFP2_laz_rev > 0) && (SFP2_com_rev/(SFP2_laz_rev) > 1.7)) || ((SFP2_com_rev > 0) && (SFP2_laz_rev/(SFP2_com_rev) > 1.7)))) {
			return 999
		} else {
			if (((SFP1_com_rev > 0.8*float32(step1) && SFP1_com_rev < 1.2*float32(step1)) && (SFP1_laz_rev > 0.8*float32(step1) && SFP1_laz_rev < 1.2*float32(step1))) ||
				((SFP2_com_rev > 0.8*float32(step1) && SFP2_com_rev < 1.2*float32(step1)) && (SFP2_laz_rev > 0.8*float32(step1) && SFP2_laz_rev < 1.2*float32(step1)))) &&
				((((SFP1_laz > 0) && (SFP1_com/(SFP1_laz) > 1.7)) || ((SFP1_com > 0) && (SFP1_laz/(SFP1_com) > 1.7))) ||
					(((SFP2_laz > 0) && (SFP2_com/(SFP2_laz) > 1.7)) || ((SFP2_com > 0) && (SFP2_laz/(SFP2_com) > 1.7)))) {
				return 1999
			}
		}

		if ((SFP1_com > 1.6*float32(step1) && SFP1_laz < 0.3*float32(step1)) || (SFP1_laz > 1.6*float32(step1) && SFP1_com < 0.3*float32(step1))) &&
			(SFP2_com < 1.1*float32(step1) && SFP2_laz < 1.1*float32(step1)) &&
			(((SFP2_com_rev > 1.6*float32(step1) && SFP2_laz_rev < 0.3*float32(step1)) || (SFP2_laz_rev > 1.6*float32(step1) && SFP2_com_rev < 0.3*float32(step1))) &&
				(SFP1_com_rev < 1.1*float32(step1) && SFP1_laz_rev < 1.1*float32(step1))) {
			return 2999
		}

	}

	revers = false
	SFP1_com_s2, SFP1_laz_s2, SFP2_com_s2, SFP2_laz_s2 := testWay.getResult(step2, revers)

	fmt.Printf("\n\n  SFP1_com - %v   SFP1_laz - %v ", SFP1_com_s2, SFP1_laz_s2)
	fmt.Printf("\n  SFP2_com - %v   SFP2_laz - %v \n", SFP2_com_s2, SFP2_laz_s2)

	revers = true
	SFP1_com_s2_rev, SFP1_laz_s2_rev, SFP2_com_s2_rev, SFP2_laz_s2_rev := testWay.getResult(step2, revers)

	fmt.Printf("\n\n  SFP1_com - %v   SFP1_laz - %v ", SFP1_com_s2_rev, SFP1_laz_s2_rev)
	fmt.Printf("\n  SFP2_com - %v   SFP2_laz - %v \n", SFP2_com_s2_rev, SFP2_laz_s2_rev)

	revers = false
	SFP1_com_s3, SFP1_laz_s3, SFP2_com_s3, SFP2_laz_s3 := testWay.getResult(step3, revers)

	fmt.Printf("\n\n  SFP1_com - %v   SFP1_laz - %v ", SFP1_com_s3, SFP1_laz_s3)
	fmt.Printf("\n  SFP2_com - %v   SFP2_laz - %v \n", SFP2_com_s3, SFP2_laz_s3)

	revers = true
	SFP1_com_s3_rev, SFP1_laz_s3_rev, SFP2_com_s3_rev, SFP2_laz_s3_rev := testWay.getResult(step3, revers)

	fmt.Printf("\n\n  SFP1_com - %v   SFP1_laz - %v ", SFP1_com_s3_rev, SFP1_laz_s3_rev)
	fmt.Printf("\n  SFP2_com - %v   SFP2_laz - %v \n", SFP2_com_s3_rev, SFP2_laz_s3_rev)

	/*

			period_nano_gen_s = int64(packet_size_generate * 8 * 1000 / step1)
			//period_nano_gen := (int64)(5*packet_size_generate*period_nano_gen_s) / (int64)(packet_size) // 20% от исходной
			period_nano_gen := (int64)(packet_size_generate*period_nano_gen_s) / (int64)(packet_size) // 100% от исходной
			count_gen := 10000000000 / period_nano_gen

			fmt.Println(" ==> TEST SFP way ==")

			var SFP1_com_init, SFP1_laz_init, SFP2_com_init, SFP2_laz_init int64

			oids := []string{".1.3.6.1.4.1.2010.1.13.0", ".1.3.6.1.4.1.2010.1.14.0"}

			counter := 0

			timer_SNMP := time.NewTicker(500 * time.Millisecond)
			for range timer_SNMP.C {
				counter++
				if counter > 20 {
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

			fmt.Printf("\n  SFP1_com - %v   SFP1_laz - %v ", float32(SFP1_com_init*8)/10000000.0, float32(SFP1_laz_init*8)/10000000.0)
			fmt.Printf("\n  SFP2_com - %v   SFP2_laz - %v ", float32(SFP2_com_init*8)/10000000.0, float32(SFP2_laz_init*8)/10000000.0)

			fmt.Printf("\n->200 Mb/s<-\n\n")
			fmt.Printf("\n->Generate normal<-\n\n")

			ipsrc := net.ParseIP(ip_server)
			ipdst1 := net.ParseIP(ip_1sfpsla_str)
			ipdst2 := net.ParseIP(ip_2sfpsla_str)



			SFP1_com_load = SFP1_com_load - SFP1_com_init
			SFP1_laz_load = SFP1_laz_load - SFP1_laz_init
			SFP2_com_load = SFP2_com_load - SFP2_com_init
			SFP2_laz_load = SFP2_laz_load - SFP2_laz_init

			fmt.Printf("\n  SFP1_com - %v   SFP1_laz - %v ", float32(SFP1_com_load*8)/10000000.0, float32(SFP1_laz_load*8)/10000000.0)
			fmt.Printf("\n  SFP2_com - %v   SFP2_laz - %v \n", float32(SFP2_com_load*8)/10000000.0, float32(SFP2_laz_load*8)/10000000.0)

			time.Sleep(2 * time.Second)

			fmt.Printf("\n->Generate revers<-\n\n")

			ipsrc = net.ParseIP(ip_server)
			ipdst2 = net.ParseIP(ip_1sfpsla_str)
			ipdst1 = net.ParseIP(ip_2sfpsla_str)

			b = packetForm(ipsrc, ipdst1, ipdst2, mac_src, mac_dst2, int(packet_size_generate), 1, test_type, testWay)

			go func() {
				cc := count_gen
				ticker := time.NewTicker(time.Duration(period_nano_gen) * time.Nanosecond)
				for range ticker.C {
					cc--
					if cc > 0 {
						c.WriteTo(b, addr)
					} else {
						return
					}
				}
			}()

			var SFP1_com_load_rev, SFP1_laz_load_rev, SFP2_com_load_rev, SFP2_laz_load_rev int64

			counter = 0

			timer_SNMP = time.NewTicker(500 * time.Millisecond)
			for range timer_SNMP.C {
				counter++
				if counter > 20 {
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
						SFP1_com_load_rev = SFP1_com_load_rev + g.ToBigInt(variable.Value).Int64()
					}
					if variable.Name == ".1.3.6.1.4.1.2010.1.14.0" {
						//	fmt.Printf("SFP1 number: %v   Mb/s\n", float32(g.ToBigInt(variable.Value).Int64()*8)/1000000.0)
						SFP1_laz_load_rev = SFP1_laz_load_rev + g.ToBigInt(variable.Value).Int64()
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
						SFP2_com_load_rev = SFP2_com_load_rev + g.ToBigInt(variable.Value).Int64()
					}
					if variable.Name == ".1.3.6.1.4.1.2010.1.14.0" {
						//	fmt.Printf("SFP2 number: %v   Mb/s\n", float32(g.ToBigInt(variable.Value).Int64()*8)/1000000.0)
						SFP2_laz_load_rev = SFP2_laz_load_rev + g.ToBigInt(variable.Value).Int64()
					}
				}

			}

			SFP1_com_load_rev = SFP1_com_load_rev - SFP1_com_init
			SFP1_laz_load_rev = SFP1_laz_load_rev - SFP1_laz_init
			SFP2_com_load_rev = SFP2_com_load_rev - SFP2_com_init
			SFP2_laz_load_rev = SFP2_laz_load_rev - SFP2_laz_init

			fmt.Printf("\n  SFP1_com - %v SFP1_laz - %v ", float32(SFP1_com_load_rev*8)/10000000.0, float32(SFP1_laz_load_rev*8)/10000000.0)
			fmt.Printf("\n  SFP2_com - %v SFP2_laz - %v \n", float32(SFP2_com_load_rev*8)/10000000.0, float32(SFP2_laz_load_rev*8)/10000000.0)

			time.Sleep(2 * time.Second)

			fmt.Printf("\n->70 Mb/s<-\n\n")
			fmt.Printf("\n->Generate normal<-\n\n")

			ipsrc = net.ParseIP(ip_server)
			ipdst1 = net.ParseIP(ip_1sfpsla_str)
			ipdst2 = net.ParseIP(ip_2sfpsla_str)

			period_nano_gen_s = int64(packet_size_generate * 8 * 1000 / step2)
			period_nano_gen = (int64)(packet_size_generate*period_nano_gen_s) / (int64)(packet_size) // 100% от исходной
			count_gen = 10000000000 / period_nano_gen

			b = packetForm(ipsrc, ipdst1, ipdst2, mac_src, mac_dst, int(packet_size_generate), 1, test_type, testWay)

			go func() {
				cc := count_gen
				ticker := time.NewTicker(time.Duration(period_nano_gen) * time.Nanosecond)
				for range ticker.C {
					cc--
					if cc > 0 {
						c.WriteTo(b, addr)
					} else {
						return
					}
				}
			}()

			SFP1_com_load = 0
			SFP1_laz_load = 0
			SFP2_com_load = 0
			SFP2_laz_load = 0

			counter = 0

			timer_SNMP = time.NewTicker(500 * time.Millisecond)
			for range timer_SNMP.C {
				counter++
				if counter > 20 {
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

			fmt.Printf("\n  SFP1_com - %v   SFP1_laz - %v ", float32(SFP1_com_load*8)/10000000.0, float32(SFP1_laz_load*8)/10000000.0)
			fmt.Printf("\n  SFP2_com - %v   SFP2_laz - %v \n", float32(SFP2_com_load*8)/10000000.0, float32(SFP2_laz_load*8)/10000000.0)

			time.Sleep(2 * time.Second)

			fmt.Printf("\n->Generate revers<-\n\n")

			ipsrc = net.ParseIP(ip_server)
			ipdst2 = net.ParseIP(ip_1sfpsla_str)
			ipdst1 = net.ParseIP(ip_2sfpsla_str)

			b = packetForm(ipsrc, ipdst1, ipdst2, mac_src, mac_dst2, int(packet_size_generate), 1, test_type, testWay)

			go func() {
				cc := count_gen
				ticker := time.NewTicker(time.Duration(period_nano_gen) * time.Nanosecond)
				for range ticker.C {
					cc--
					if cc > 0 {
						c.WriteTo(b, addr)
					} else {
						return
					}
				}
			}()

			SFP1_com_load_rev = 0
			SFP1_laz_load_rev = 0
			SFP2_com_load_rev = 0
			SFP2_laz_load_rev = 0

			counter = 0

			timer_SNMP = time.NewTicker(500 * time.Millisecond)
			for range timer_SNMP.C {
				counter++
				if counter > 20 {
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
						SFP1_com_load_rev = SFP1_com_load_rev + g.ToBigInt(variable.Value).Int64()
					}
					if variable.Name == ".1.3.6.1.4.1.2010.1.14.0" {
						//	fmt.Printf("SFP1 number: %v   Mb/s\n", float32(g.ToBigInt(variable.Value).Int64()*8)/1000000.0)
						SFP1_laz_load_rev = SFP1_laz_load_rev + g.ToBigInt(variable.Value).Int64()
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
						SFP2_com_load_rev = SFP2_com_load_rev + g.ToBigInt(variable.Value).Int64()
					}
					if variable.Name == ".1.3.6.1.4.1.2010.1.14.0" {
						//	fmt.Printf("SFP2 number: %v   Mb/s\n", float32(g.ToBigInt(variable.Value).Int64()*8)/1000000.0)
						SFP2_laz_load_rev = SFP2_laz_load_rev + g.ToBigInt(variable.Value).Int64()
					}
				}

			}

			SFP1_com_load_rev = SFP1_com_load_rev - SFP1_com_init
			SFP1_laz_load_rev = SFP1_laz_load_rev - SFP1_laz_init
			SFP2_com_load_rev = SFP2_com_load_rev - SFP2_com_init
			SFP2_laz_load_rev = SFP2_laz_load_rev - SFP2_laz_init

			fmt.Printf("\n  SFP1_com - %v SFP1_laz - %v ", float32(SFP1_com_load_rev*8)/10000000.0, float32(SFP1_laz_load_rev*8)/10000000.0)
			fmt.Printf("\n  SFP2_com - %v SFP2_laz - %v \n", float32(SFP2_com_load_rev*8)/10000000.0, float32(SFP2_laz_load_rev*8)/10000000.0)

			time.Sleep(2 * time.Second)

			fmt.Printf("\n->7 Mb/s<-\n\n")
			fmt.Printf("\n->Generate normal<-\n\n")

			ipsrc = net.ParseIP(ip_server)
			ipdst1 = net.ParseIP(ip_1sfpsla_str)
			ipdst2 = net.ParseIP(ip_2sfpsla_str)

			period_nano_gen_s = int64(packet_size_generate * 8 * 1000 / step3)
			period_nano_gen = (int64)(packet_size_generate*period_nano_gen_s) / (int64)(packet_size) // 100% от исходной
			count_gen = 10000000000 / period_nano_gen

			b = packetForm(ipsrc, ipdst1, ipdst2, mac_src, mac_dst, int(packet_size_generate), 1, test_type, testWay)

			go func() {
				cc := count_gen
				ticker := time.NewTicker(time.Duration(period_nano_gen) * time.Nanosecond)
				for range ticker.C {
					cc--
					if cc > 0 {
						c.WriteTo(b, addr)
					} else {
						return
					}
				}
			}()

			SFP1_com_load = 0
			SFP1_laz_load = 0
			SFP2_com_load = 0
			SFP2_laz_load = 0

			counter = 0

			timer_SNMP = time.NewTicker(500 * time.Millisecond)



			for range timer_SNMP.C {
				counter++
				if counter > 20 {
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

			fmt.Printf("\n  SFP1_com - %v   SFP1_laz - %v ", float32(SFP1_com_load*8)/10000000.0, float32(SFP1_laz_load*8)/10000000.0)
			fmt.Printf("\n  SFP2_com - %v   SFP2_laz - %v \n", float32(SFP2_com_load*8)/10000000.0, float32(SFP2_laz_load*8)/10000000.0)
			time.Sleep(2 * time.Second)

			fmt.Printf("\n->Generate revers<-\n\n")

			ipsrc = net.ParseIP(ip_server)
			ipdst2 = net.ParseIP(ip_1sfpsla_str)
			ipdst1 = net.ParseIP(ip_2sfpsla_str)

			b = packetForm(ipsrc, ipdst1, ipdst2, mac_src, mac_dst2, int(packet_size_generate), 1, test_type, testWay)

			go func() {
				cc := count_gen
				ticker := time.NewTicker(time.Duration(period_nano_gen) * time.Nanosecond)
				for range ticker.C {
					cc--
					if cc > 0 {
						c.WriteTo(b, addr)
					} else {
						return
					}
				}
			}()

			SFP1_com_load_rev = 0
			SFP1_laz_load_rev = 0
			SFP2_com_load_rev = 0
			SFP2_laz_load_rev = 0

			counter = 0

			timer_SNMP = time.NewTicker(500 * time.Millisecond)
			for range timer_SNMP.C {
				counter++
				if counter > 20 {
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
						SFP1_com_load_rev = SFP1_com_load_rev + g.ToBigInt(variable.Value).Int64()
					}
					if variable.Name == ".1.3.6.1.4.1.2010.1.14.0" {
						//	fmt.Printf("SFP1 number: %v   Mb/s\n", float32(g.ToBigInt(variable.Value).Int64()*8)/1000000.0)
						SFP1_laz_load_rev = SFP1_laz_load_rev + g.ToBigInt(variable.Value).Int64()
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
						SFP2_com_load_rev = SFP2_com_load_rev + g.ToBigInt(variable.Value).Int64()
					}
					if variable.Name == ".1.3.6.1.4.1.2010.1.14.0" {
						//	fmt.Printf("SFP2 number: %v   Mb/s\n", float32(g.ToBigInt(variable.Value).Int64()*8)/1000000.0)
						SFP2_laz_load_rev = SFP2_laz_load_rev + g.ToBigInt(variable.Value).Int64()
					}
				}

			}

			SFP1_com_load_rev = SFP1_com_load_rev - SFP1_com_init
			SFP1_laz_load_rev = SFP1_laz_load_rev - SFP1_laz_init
			SFP2_com_load_rev = SFP2_com_load_rev - SFP2_com_init
			SFP2_laz_load_rev = SFP2_laz_load_rev - SFP2_laz_init

		fmt.Printf("\n  SFP1_com - %v SFP1_laz - %v ", float32(SFP1_com_load_rev*8)/10000000.0, float32(SFP1_laz_load_rev*8)/10000000.0)
		fmt.Printf("\n  SFP2_com - %v SFP2_laz - %v \n", float32(SFP2_com_load_rev*8)/10000000.0, float32(SFP2_laz_load_rev*8)/10000000.0)

		min_load := 7 * (1024 * 1000000000 / period_nano_gen)
		fmt.Printf("\n  Min_load - %v \n", min_load)

		if (SFP1_com_load > min_load) && (SFP1_laz_load > min_load) && ((float32(SFP1_laz_load)/float32(SFP1_com_load) > 1.7) || (float32(SFP1_com_load)/float32(SFP1_laz_load) > 1.7)) {
			fmt.Printf("\n  Modules state change \n")
			return 2
		} else {
			fmt.Printf("\n  SFP2_laz / SFP2_com %f \n", float32(SFP2_laz_load)/float32(SFP2_com_load))
			fmt.Printf("\n  SFP2_com / SFP2_laz %f \n", float32(SFP2_com_load)/float32(SFP2_laz_load))
			if (SFP2_com_load > min_load) && (SFP2_laz_load > min_load) && ((float32(SFP2_laz_load)/float32(SFP2_com_load) > 1.7) || (float32(SFP2_com_load)/float32(SFP2_laz_load) > 1.7)) {
				fmt.Printf("\n  Modules state change \n")
				return 2
			} else {
				fmt.Printf("\n  Modules state OK \n")
				return 1
			}
		}
	*/
	return 1
}

func (module *module_sfp) startSNMP(conf global_config) {

	fmt.Println("Start SNMP_check - ", (*module).address_ip)
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
			fmt.Println("End SNMP_check - ", (*module).address_ip)
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
				id_test := 0
				row_mod := db_SNMP.QueryRow("SELECT id FROM modules_sfp_sla WHERE id=?", (*module).id)
				err = row_mod.Scan(&id_test)

				if err != nil || id_test == 0 {
					fmt.Println(" -!! Error check module !!-")
					return
				}

				db_SNMP.Exec("INSERT INTO modules_sfp_sla_load_rez (module_id, datatime, load_to_lazer, load_to_com) VALUES(?, NOW(), ?, ?)", (*module).id, SFP_laz, SFP_com)

				//fmt.Println("Rez add module metric: ",rez)
				//fmt.Println("Error add module metric: ",er)
				//	fmt.Println(" add metric - ", time.Now())
				//db_SNMP.Close()

			}
		}
	}
}

func check_SNMP(ip string) int {

	envTarget := ip
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

	rez := 0

	to := time.After(5 * time.Second)

	done := make(chan bool, 1)

	oids := []string{".1.3.6.1.4.1.2010.1.13.0"}

	go func() {
		timer_SNMP := time.NewTicker(900 * time.Millisecond)
		for range timer_SNMP.C {
			select {
			case <-to:
				done <- true
				timer_SNMP.Stop()
				return
			default:
				mux_SNMP.Lock()
				err := params.Connect()
				if err != nil {
					fmt.Print("Module : ", envTarget)
					fmt.Println("  Connect() err: ", err)
					params.Conn.Close()
					mux_SNMP.Unlock()
					continue
				}

				result, err2 := params.Get(oids) // Get() accepts up to g.MAX_OIDS
				if err2 != nil {
					fmt.Print("Module : ", envTarget)
					fmt.Println("Get() err: ", err2)
					params.Conn.Close()
					mux_SNMP.Unlock()
					continue
				}

				for range result.Variables {

					rez++
				}
				mux_SNMP.Unlock()

			}
		}
	}()

	<-done
	fmt.Println("\nModule : %s, rez = %d", ip, rez)
	if rez > 3 {
		return 0
	} else {
		return 1
	}

}
