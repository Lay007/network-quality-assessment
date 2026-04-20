package main

import (
	"bytes"
	"fmt"
	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/raw"
	"golang.org/x/net/bpf"
	"math"
	"net"
	"runtime"
	"time"
)

func TestDelay(id int, net_interface_name string) {
	if verboseLogs {
		fmt.Println("Latency test started")
	}
	db, err := openDB()
	if err != nil {
		db.Close()
		if verboseLogs {
			fmt.Println(" -!! Error !!-")
		}
		if verboseLogs {
			fmt.Println(err)
		}
		if verboseLogs {
			fmt.Println(" ----=====----")
		}
		return
	}
	db.Exec("UPDATE test_latency SET status=?, datatime=? WHERE id=?", time.Now().Format("2006-01-02 15:04:05"), 2, id)

	ifi, err := net.InterfaceByName(net_interface_name)
	if err != nil {
		db.Close()
		//	log.Fatalf("failed to find interface %q: %v", net_interface_name, err)
		db.Exec("UPDATE test_latency SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
		return
	}

	row, err := db.Query("select id, miss_init_test, test_type, module_first, module_second, thr_begin, count_packs, count_tests, status from test_latency where id=?", id)
	if err != nil {
		db.Close()
		row.Close()
		if verboseLogs {
			fmt.Println(" -!! Error !!-")
		}
		if verboseLogs {
			fmt.Println(err)
		}
		if verboseLogs {
			fmt.Println(" ----=====----")
		}
		return
	}
	if verboseLogs {
		fmt.Println(row)
	}
	defer row.Close()
	row.Next()
	test := new(testDelay)
	err = row.Scan(&test.id, &test.miss_init_test, &test.test_type, &test.module_first, &test.module_second, &test.thr_begin, &test.count_packs, &test.count_tests, &test.status)
	if err != nil {
		db.Close()
		if verboseLogs {
			fmt.Println(" -!! Error !!-")
		}
		if verboseLogs {
			fmt.Println(err)
		}
		if verboseLogs {
			fmt.Println(" ----=====----")
		}
		return
	}

	var ipsrcstr string
	var ipdst_1sfpsla_str string
	var ipdst_2sfpsla_str string
	test.mac_dst = make([]byte, 6)
	test.mac_dst2 = make([]byte, 6)

	row, err = db.Query("SELECT server_IP FROM global_config")
	if err != nil {
		db.Exec("UPDATE test_latency SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
		db.Close()
		row.Close()
		if verboseLogs {
			fmt.Println(" -!! Error !!-")
		}
		if verboseLogs {
			fmt.Println(err)
		}
		if verboseLogs {
			fmt.Println(" ----=====----")
		}
		return
	}
	for row.Next() {
		err = row.Scan(&ipsrcstr)
		if err != nil {

			db.Exec("UPDATE test_latency SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
			db.Close()
			row.Close()
			if verboseLogs {
				fmt.Println(" -!! Error !!-")
			}
			if verboseLogs {
				fmt.Println(err)
			}
			if verboseLogs {
				fmt.Println(" ----=====----")
			}
			return
		}
	}
	var id_sfp1, id_sfp2 int
	row, err = db.Query("SELECT module_first, module_second FROM test_latency WHERE id=?", id)
	if err != nil {
		db.Exec("UPDATE test_latency SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
		db.Close()
		row.Close()
		if verboseLogs {
			fmt.Println(" -!! Error !!-")
		}
		if verboseLogs {
			fmt.Println(err)
		}
		if verboseLogs {
			fmt.Println(" ----=====----")
		}
		return
	}
	for row.Next() {
		err = row.Scan(&id_sfp1, &id_sfp2)
		if err != nil {

			db.Exec("UPDATE test_latency SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
			db.Close()
			row.Close()
			if verboseLogs {
				fmt.Println(" -!! Error !!-")
			}
			if verboseLogs {
				fmt.Println(err)
			}
			if verboseLogs {
				fmt.Println(" ----=====----")
			}
			return
		}
		row_ip, err := db.Query("SELECT address_ip FROM modules_sfp_sla WHERE id=?", id_sfp1)
		if err != nil {
			db.Exec("UPDATE test_latency SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
			db.Close()
			row.Close()
			row_ip.Close()
			if verboseLogs {
				fmt.Println(" -!! Error !!-")
			}
			if verboseLogs {
				fmt.Println(err)
			}
			if verboseLogs {
				fmt.Println(" ----=====----")
			}
			return
		}
		defer row_ip.Close()
		for row_ip.Next() {
			err = row_ip.Scan(&ipdst_1sfpsla_str)
			if err != nil {

				db.Exec("UPDATE test_latency SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
				db.Close()
				row.Close()
				if verboseLogs {
					fmt.Println(" -!! Error !!-")
				}
				if verboseLogs {
					fmt.Println(err)
				}
				if verboseLogs {
					fmt.Println(" ----=====----")
				}
				return
			}
		}

		row_mac, err := db.Query("SELECT mac FROM modules_sfp_sla WHERE id=?", id_sfp1)
		if err != nil {
			db.Exec("UPDATE test_latency SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
			db.Close()
			row.Close()
			row_mac.Close()
			if verboseLogs {
				fmt.Println(" -!! Error !!-")
			}
			if verboseLogs {
				fmt.Println(err)
			}
			if verboseLogs {
				fmt.Println(" ----=====----")
			}
			return
		}
		defer row_mac.Close()
		var test_mac int64
		for row_mac.Next() {
			//err = row_mac.Scan(&mac_dst_str)
			err = row_mac.Scan(&test_mac)
			if err != nil {

				db.Exec("UPDATE test_latency SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
				db.Close()
				row.Close()
				if verboseLogs {
					fmt.Println(" -!! Error !!-")
				}
				if verboseLogs {
					fmt.Println(err)
				}
				if verboseLogs {
					fmt.Println(" ----=====----")
				}
				return
			}
		}
		test.mac_dst[5] = byte(test_mac & 0xFF)
		test.mac_dst[4] = byte((test_mac >> 8) & 0xFF)
		test.mac_dst[3] = byte((test_mac >> 16) & 0xFF)
		test.mac_dst[2] = byte((test_mac >> 24) & 0xFF)
		test.mac_dst[1] = byte((test_mac >> 32) & 0xFF)
		test.mac_dst[0] = byte((test_mac >> 40) & 0xFF)

		row_mac, err = db.Query("SELECT mac FROM modules_sfp_sla WHERE id=?", id_sfp2)
		if err != nil {
			db.Exec("UPDATE test_latency SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
			db.Close()
			row.Close()
			row_mac.Close()
			if verboseLogs {
				fmt.Println(" -!! Error !!-")
			}
			if verboseLogs {
				fmt.Println(err)
			}
			if verboseLogs {
				fmt.Println(" ----=====----")
			}
			return
		}
		defer row_mac.Close()
		for row_mac.Next() {
			//err = row_mac.Scan(&mac_dst_str)
			err = row_mac.Scan(&test_mac)
			if err != nil {

				db.Exec("UPDATE test_latency SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
				db.Close()
				row.Close()
				if verboseLogs {
					fmt.Println(" -!! Error !!-")
				}
				if verboseLogs {
					fmt.Println(err)
				}
				if verboseLogs {
					fmt.Println(" ----=====----")
				}
				return
			}
		}
		test.mac_dst2[5] = byte(test_mac & 0xFF)
		test.mac_dst2[4] = byte((test_mac >> 8) & 0xFF)
		test.mac_dst2[3] = byte((test_mac >> 16) & 0xFF)
		test.mac_dst2[2] = byte((test_mac >> 24) & 0xFF)
		test.mac_dst2[1] = byte((test_mac >> 32) & 0xFF)
		test.mac_dst2[0] = byte((test_mac >> 40) & 0xFF)

		row_ip, err = db.Query("SELECT address_ip FROM modules_sfp_sla WHERE id=?", id_sfp2)
		if err != nil {
			db.Exec("UPDATE test_latency SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
			db.Close()
			row.Close()
			if verboseLogs {
				fmt.Println(" -!! Error !!-")
			}
			if verboseLogs {
				fmt.Println(err)
			}
			if verboseLogs {
				fmt.Println(" ----=====----")
			}
			return
		}
		for row_ip.Next() {
			err = row_ip.Scan(&ipdst_2sfpsla_str)
			if err != nil {

				db.Exec("UPDATE test_latency SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
				db.Close()
				row.Close()
				if verboseLogs {
					fmt.Println(" -!! Error !!-")
				}
				if verboseLogs {
					fmt.Println(err)
				}
				if verboseLogs {
					fmt.Println(" ----=====----")
				}
				return
			}
		}
	}
	db.Exec("UPDATE test_latency SET status=?, datetime_start=? WHERE id=?", 2, time.Now().Format("2006-01-02 15:04:05"), id)
	test.status = 2
	if testPing(ipdst_1sfpsla_str) > 0 || testPing(ipdst_2sfpsla_str) > 0 {
		db.Exec("UPDATE test_latency SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
		db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 2, id, "Ping test failed")
		db.Close()

		return
	}

	if check_SNMP(ipdst_1sfpsla_str) > 0 || check_SNMP(ipdst_2sfpsla_str) > 0 {
		db.Exec("UPDATE test_latency SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
		db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 2, id, "SNMP test failed")
		db.Close()

		return
	}

	if verboseLogs {
		fmt.Println(ipsrcstr)
	}
	if verboseLogs {
		fmt.Println(ipdst_1sfpsla_str)
	}
	if verboseLogs {
		fmt.Println(ipdst_2sfpsla_str)
	}

	//counter := test.count
	//	numberTX = 0

	test.ipsrc = net.ParseIP(ipsrcstr)
	test.ipdst1 = net.ParseIP(ipdst_1sfpsla_str)
	test.ipdst2 = net.ParseIP(ipdst_2sfpsla_str)

	//period_gen := time.Duration(10 * time.Second)

	/*
		test_counter:=10^9
		start_test_ticker := time.Now()
		t := time.NewTicker(time.Duration(1))
		//t := time.NewTicker(1 * time.Second)
			for range t.C {
				test_counter--
				if test_counter<0{
					t.Stop()
					break
				}
			}
			min_ticker:=time.Since(start_test_ticker)
		if verboseLogs {
			fmt.Println(" -- Test ticker= ", time.Since(start_test_ticker))
		}
	*/

	//	numberTX++
	/*
		addr := &raw.Addr{
			HardwareAddr: ethernet.Broadcast,
		}
		var number uint32
	*/

	addr := &raw.Addr{
		HardwareAddr: ethernet.Broadcast,
	}
	rez := 1
	connectTestSFP, err := raw.ListenPacket(ifi, etherType, nil)

	yest, _ := db.Query("SELECT EXISTS(SELECT id FROM test_latency WHERE id=?)", id)
	yest.Next()
	t_y := 0
	yest.Scan(&t_y)
	if t_y != 1 {
		return
	}

	if test.miss_init_test == 0 {
		rez := findSFP(connectTestSFP, addr, ipsrcstr, ipdst_1sfpsla_str, ipdst_2sfpsla_str, test.mac_src, test.mac_dst, test.mac_dst2, test.id_test_type, test.test_type, 1024, int64(1024*8*1000/test.thr_begin))

		if verboseLogs {
			fmt.Printf("\n Rez find : %X \n", rez)
		}
		if (rez & 0xFFF) == 0x999 {
			db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 2, id, " Maximum throughput - 1 Gbit/s")
			if test.thr_begin > 1000 {
				test.thr_begin = 1000
				db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 2, id, " Changing maximum throughput to 1 Gbit/s")
			}
		}
		if (rez & 0xFFF) == 0x100 {
			db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 2, id, " Maximum throughput - 100 Mbit/s")
			if test.thr_begin > 100 {
				test.thr_begin = 100
				db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 2, id, " Changing maximum throughput to 100 Mbit/s")
			}
		}
		if (rez & 0xFFF) == 0x10 {
			db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 2, id, " Maximum throughput - 10 Mbit/s")
			if test.thr_begin > 10 {
				test.thr_begin = 10
				db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 2, id, " Changing maximum throughput to 10 Mbit/s")
			}
		}
		if (rez & 0xF000) == 0x0 {
			db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 2, id, " Serial module topology. Module order is correct")
		}
		if (rez & 0xF000) == 0x1000 {
			db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 2, id, " Serial module topology. Module order is incorrect; swapping during the test")
		}
		if (rez & 0xF000) == 0x2000 {
			db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 2, id, " Star topology. Load is balanced")
		}
		if (rez & 0xF000) == 0x3000 {
			db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 2, id, " Star topology. Load is unbalanced; module order is correct")
		}
		if (rez & 0xF000) == 0x4000 {
			db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 2, id, " Star topology. Load is unbalanced; module order is incorrect; swapping during the test")
		}
	}
	//db.Close()
	connectTestSFP.Close()
	if rez == 0 {
		if verboseLogs {
			fmt.Println("Error test SFP connect")
		}
		return
	}
	if ((rez & 0xF000) == 0x1000) || ((rez & 0xF000) == 0x4000) {

		tmp := ipdst_1sfpsla_str
		ipdst_1sfpsla_str = ipdst_2sfpsla_str
		ipdst_2sfpsla_str = tmp

		test.ipsrc = net.ParseIP(ipsrcstr)
		test.ipdst1 = net.ParseIP(ipdst_1sfpsla_str)
		test.ipdst2 = net.ParseIP(ipdst_2sfpsla_str)

		//db, err = openDB()
		row_mac, err := db.Query("SELECT mac FROM modules_sfp_sla WHERE address_ip=?", ipdst_1sfpsla_str)
		if err != nil {
			db.Close()
			row_mac.Close()
			if verboseLogs {
				fmt.Println(" -!! Error !!-")
			}
			if verboseLogs {
				fmt.Println(err)
			}
			if verboseLogs {
				fmt.Println(" ----=====----")
			}
			return
		}
		defer row_mac.Close()
		var test_mac int64
		for row_mac.Next() {
			err = row_mac.Scan(&test_mac)
			if err != nil {
				db.Close()
				row_mac.Close()
				if verboseLogs {
					fmt.Println(" -!! Error !!-")
				}
				if verboseLogs {
					fmt.Println(err)
				}
				if verboseLogs {
					fmt.Println(" ----=====----")
				}
				return
			}
		}
		test.mac_dst[5] = byte(test_mac & 0xFF)
		test.mac_dst[4] = byte((test_mac >> 8) & 0xFF)
		test.mac_dst[3] = byte((test_mac >> 16) & 0xFF)
		test.mac_dst[2] = byte((test_mac >> 24) & 0xFF)
		test.mac_dst[1] = byte((test_mac >> 32) & 0xFF)
		test.mac_dst[0] = byte((test_mac >> 40) & 0xFF)

		//	db.Close()

	}
	if verboseLogs {
		fmt.Println(" Rez find : ", rez)
	}

	test.id_test_type = 0x6000 + (uint16(id) & 0x1FFF)
	testTypeTemp := 0x2000 + (uint16(id) & 0x1FFF)

	test.net_interface_name = net_interface_name
	test.mac_src = ifi.HardwareAddr

	counter := make(chan uint64, 7)
	counterRes := make(chan uint64, 7)
	quit := make(chan int64, 7)

	size_p := 64

	b := packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size_p, 0, test.id_test_type, test.test_type)
	bTemp := packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size_p, 0, testTypeTemp, test.test_type)
	go genSocket(ifi.Index, b, bTemp, test.count_packs, test.thr_begin, counter, counterRes)
	test.getMonDelay(quit, size_p)
	time.Sleep(time.Second * 2)
	<-quit

	yest, _ = db.Query("SELECT EXISTS(SELECT id FROM test_latency WHERE id=?)", id)
	yest.Next()
	t_y = 0
	yest.Scan(&t_y)
	if t_y != 1 {
		return
	}
	db.Exec("UPDATE test_latency SET rez_64=?,rez_64_max=?,rez_64_min=?,rez_128=?,rez_128_max=?,rez_128_min=?,rez_256=?,rez_256_max=?,rez_256_min=?,rez_512=?,rez_512_max=?,rez_512_min=?,rez_1024=?,rez_1024_max=?,rez_1024_min=?,rez_1280=?,rez_1280_max=?,rez_1280_min=?, rez_1518=?,rez_1518_max=?,rez_1518_min=?, status=? WHERE id=?", test.rez_64, test.rez_64_max, test.rez_64_min, test.rez_128, test.rez_128_max, test.rez_128_min, test.rez_256, test.rez_256_max, test.rez_256_min, test.rez_512, test.rez_512_max, test.rez_512_min, test.rez_1024, test.rez_1024_max, test.rez_1024_min, test.rez_1280, test.rez_1280_max, test.rez_1280_min, test.rez_1518, test.rez_1518_max, test.rez_1518_min, test.status, id)

	size_p = 128

	b = packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size_p, 0, test.id_test_type, test.test_type)
	bTemp = packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size_p, 0, testTypeTemp, test.test_type)
	go genSocket(ifi.Index, b, bTemp, test.count_packs, test.thr_begin, counter, counterRes)
	test.getMonDelay(quit, size_p)
	time.Sleep(time.Second * 2)
	<-quit

	yest, _ = db.Query("SELECT EXISTS(SELECT id FROM test_latency WHERE id=?)", id)
	yest.Next()
	t_y = 0
	yest.Scan(&t_y)
	if t_y != 1 {
		return
	}
	db.Exec("UPDATE test_latency SET rez_64=?,rez_64_max=?,rez_64_min=?,rez_128=?,rez_128_max=?,rez_128_min=?,rez_256=?,rez_256_max=?,rez_256_min=?,rez_512=?,rez_512_max=?,rez_512_min=?,rez_1024=?,rez_1024_max=?,rez_1024_min=?,rez_1280=?,rez_1280_max=?,rez_1280_min=?, rez_1518=?,rez_1518_max=?,rez_1518_min=?, status=? WHERE id=?", test.rez_64, test.rez_64_max, test.rez_64_min, test.rez_128, test.rez_128_max, test.rez_128_min, test.rez_256, test.rez_256_max, test.rez_256_min, test.rez_512, test.rez_512_max, test.rez_512_min, test.rez_1024, test.rez_1024_max, test.rez_1024_min, test.rez_1280, test.rez_1280_max, test.rez_1280_min, test.rez_1518, test.rez_1518_max, test.rez_1518_min, test.status, id)

	size_p = 256

	b = packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size_p, 0, test.id_test_type, test.test_type)
	bTemp = packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size_p, 0, testTypeTemp, test.test_type)
	go genSocket(ifi.Index, b, bTemp, test.count_packs, test.thr_begin, counter, counterRes)
	test.getMonDelay(quit, size_p)
	time.Sleep(time.Second * 2)
	<-quit

	yest, _ = db.Query("SELECT EXISTS(SELECT id FROM test_latency WHERE id=?)", id)
	yest.Next()
	t_y = 0
	yest.Scan(&t_y)
	if t_y != 1 {
		return
	}
	db.Exec("UPDATE test_latency SET rez_64=?,rez_64_max=?,rez_64_min=?,rez_128=?,rez_128_max=?,rez_128_min=?,rez_256=?,rez_256_max=?,rez_256_min=?,rez_512=?,rez_512_max=?,rez_512_min=?,rez_1024=?,rez_1024_max=?,rez_1024_min=?,rez_1280=?,rez_1280_max=?,rez_1280_min=?, rez_1518=?,rez_1518_max=?,rez_1518_min=?, status=? WHERE id=?", test.rez_64, test.rez_64_max, test.rez_64_min, test.rez_128, test.rez_128_max, test.rez_128_min, test.rez_256, test.rez_256_max, test.rez_256_min, test.rez_512, test.rez_512_max, test.rez_512_min, test.rez_1024, test.rez_1024_max, test.rez_1024_min, test.rez_1280, test.rez_1280_max, test.rez_1280_min, test.rez_1518, test.rez_1518_max, test.rez_1518_min, test.status, id)

	size_p = 512

	b = packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size_p, 0, test.id_test_type, test.test_type)
	bTemp = packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size_p, 0, testTypeTemp, test.test_type)
	go genSocket(ifi.Index, b, bTemp, test.count_packs, test.thr_begin, counter, counterRes)
	test.getMonDelay(quit, size_p)
	time.Sleep(time.Second * 2)
	<-quit

	yest, _ = db.Query("SELECT EXISTS(SELECT id FROM test_latency WHERE id=?)", id)
	yest.Next()
	t_y = 0
	yest.Scan(&t_y)
	if t_y != 1 {
		return
	}
	db.Exec("UPDATE test_latency SET rez_64=?,rez_64_max=?,rez_64_min=?,rez_128=?,rez_128_max=?,rez_128_min=?,rez_256=?,rez_256_max=?,rez_256_min=?,rez_512=?,rez_512_max=?,rez_512_min=?,rez_1024=?,rez_1024_max=?,rez_1024_min=?,rez_1280=?,rez_1280_max=?,rez_1280_min=?, rez_1518=?,rez_1518_max=?,rez_1518_min=?, status=? WHERE id=?", test.rez_64, test.rez_64_max, test.rez_64_min, test.rez_128, test.rez_128_max, test.rez_128_min, test.rez_256, test.rez_256_max, test.rez_256_min, test.rez_512, test.rez_512_max, test.rez_512_min, test.rez_1024, test.rez_1024_max, test.rez_1024_min, test.rez_1280, test.rez_1280_max, test.rez_1280_min, test.rez_1518, test.rez_1518_max, test.rez_1518_min, test.status, id)

	size_p = 1024

	b = packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size_p, 0, test.id_test_type, test.test_type)
	bTemp = packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size_p, 0, testTypeTemp, test.test_type)
	go genSocket(ifi.Index, b, bTemp, test.count_packs, test.thr_begin, counter, counterRes)
	test.getMonDelay(quit, size_p)
	time.Sleep(time.Second * 2)
	<-quit

	yest, _ = db.Query("SELECT EXISTS(SELECT id FROM test_latency WHERE id=?)", id)
	yest.Next()
	t_y = 0
	yest.Scan(&t_y)
	if t_y != 1 {
		return
	}
	db.Exec("UPDATE test_latency SET rez_64=?,rez_64_max=?,rez_64_min=?,rez_128=?,rez_128_max=?,rez_128_min=?,rez_256=?,rez_256_max=?,rez_256_min=?,rez_512=?,rez_512_max=?,rez_512_min=?,rez_1024=?,rez_1024_max=?,rez_1024_min=?,rez_1280=?,rez_1280_max=?,rez_1280_min=?, rez_1518=?,rez_1518_max=?,rez_1518_min=?, status=? WHERE id=?", test.rez_64, test.rez_64_max, test.rez_64_min, test.rez_128, test.rez_128_max, test.rez_128_min, test.rez_256, test.rez_256_max, test.rez_256_min, test.rez_512, test.rez_512_max, test.rez_512_min, test.rez_1024, test.rez_1024_max, test.rez_1024_min, test.rez_1280, test.rez_1280_max, test.rez_1280_min, test.rez_1518, test.rez_1518_max, test.rez_1518_min, test.status, id)

	size_p = 1280

	b = packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size_p, 0, test.id_test_type, test.test_type)
	bTemp = packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size_p, 0, testTypeTemp, test.test_type)
	go genSocket(ifi.Index, b, bTemp, test.count_packs, test.thr_begin, counter, counterRes)
	test.getMonDelay(quit, size_p)
	time.Sleep(time.Second * 2)
	<-quit

	yest, _ = db.Query("SELECT EXISTS(SELECT id FROM test_latency WHERE id=?)", id)
	yest.Next()
	t_y = 0
	yest.Scan(&t_y)
	if t_y != 1 {
		return
	}

	db.Exec("UPDATE test_latency SET rez_64=?,rez_64_max=?,rez_64_min=?,rez_128=?,rez_128_max=?,rez_128_min=?,rez_256=?,rez_256_max=?,rez_256_min=?,rez_512=?,rez_512_max=?,rez_512_min=?,rez_1024=?,rez_1024_max=?,rez_1024_min=?,rez_1280=?,rez_1280_max=?,rez_1280_min=?, rez_1518=?,rez_1518_max=?,rez_1518_min=?, status=? WHERE id=?", test.rez_64, test.rez_64_max, test.rez_64_min, test.rez_128, test.rez_128_max, test.rez_128_min, test.rez_256, test.rez_256_max, test.rez_256_min, test.rez_512, test.rez_512_max, test.rez_512_min, test.rez_1024, test.rez_1024_max, test.rez_1024_min, test.rez_1280, test.rez_1280_max, test.rez_1280_min, test.rez_1518, test.rez_1518_max, test.rez_1518_min, test.status, id)

	size_p = 1518

	b = packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, 1500, 0, test.id_test_type, test.test_type)
	bTemp = packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, 1500, 0, testTypeTemp, test.test_type)
	go genSocket(ifi.Index, b, bTemp, test.count_packs, test.thr_begin, counter, counterRes)
	test.getMonDelay(quit, 1500)
	time.Sleep(time.Second * 2)
	<-quit

	test.status = 3
	/*
		db, err = openDB()
		if err != nil {
			db.Close()

			if verboseLogs {
				fmt.Println(" -!! Error !!-")
			}
			if verboseLogs {
				fmt.Println(err)
			}
			if verboseLogs {
				fmt.Println(" ----=====----")
			}
			return
		}
	*/
	//db.Exec("UPDATE test_throughput SET rez_64=?,rez_128=?,rez_256=?,rez_512=?,rez_1024=?,rez_1280=?, rez_1518=?,rez_4096=?,rez_9000=?,status=? WHERE id=?", test.rez_64, test.rez_128, test.rez_256, test.rez_512, test.rez_1024, test.rez_1280, test.rez_1518, test.rez_4096, test.rez_9000, test.status, id)
	db.Exec("UPDATE test_latency SET rez_64=?,rez_64_max=?,rez_64_min=?,rez_128=?,rez_128_max=?,rez_128_min=?,rez_256=?,rez_256_max=?,rez_256_min=?,rez_512=?,rez_512_max=?,rez_512_min=?,rez_1024=?,rez_1024_max=?,rez_1024_min=?,rez_1280=?,rez_1280_max=?,rez_1280_min=?, rez_1518=?,rez_1518_max=?,rez_1518_min=?,datetime_end=?, status=? WHERE id=?", test.rez_64, test.rez_64_max, test.rez_64_min, test.rez_128, test.rez_128_max, test.rez_128_min, test.rez_256, test.rez_256_max, test.rez_256_min, test.rez_512, test.rez_512_max, test.rez_512_min, test.rez_1024, test.rez_1024_max, test.rez_1024_min, test.rez_1280, test.rez_1280_max, test.rez_1280_min, test.rez_1518, test.rez_1518_max, test.rez_1518_min, time.Now().Format("2006-01-02 15:04:05"), test.status, id)

	db.Close()
}

func (test *testDelay) getMonDelay(quit chan int64, size int) {

	var delay, delayMax, delayMin int64
	var floatDelay, floatDelayMax, floatDelayMin float32

	detectPackDelay := make(chan int64, 1000)
	number := 0

	timeStart := time.Now()

	ifi, err := net.InterfaceByName(test.net_interface_name)
	if err != nil {
		if verboseLogs {
			fmt.Printf("failed to find interface %q: %v\n", test.net_interface_name, err)
		}
		quit <- 1
		return
	}

	var netConf *raw.Config = new(raw.Config)

	(*netConf).Filter, _ = bpf.Assemble([]bpf.Instruction{
		bpf.LoadAbsolute{Off: 34, Size: 1},
		bpf.JumpIf{Cond: bpf.JumpNotEqual, Val: 0xFC, SkipTrue: 5},
		bpf.LoadAbsolute{Off: 64, Size: 2},
		bpf.JumpIf{Cond: bpf.JumpNotEqual, Val: uint32(test.id_test_type), SkipTrue: 3},
		bpf.LoadExtension{Num: bpf.ExtRand},
		//	bpf.JumpIf{Cond: bpf.JumpLessThan, Val: 0xFF, SkipFalse: 1},
		//bpf.JumpIf{Cond: bpf.JumpGreaterThan, Val: 0x0FFFFFFF, SkipTrue: 1},
		bpf.JumpIf{Cond: bpf.JumpGreaterThan, Val: 0xFFFFFFFF, SkipTrue: 1},
		// Verdict is "send up to 4k of the packet to userspace."
		bpf.RetConstant{Val: 4096},
		// Verdict is "ignore packet."
		bpf.RetConstant{Val: 0},
	})

	c, err := raw.ListenPacket(ifi, etherType, netConf)

	SolveDelayTicker := time.NewTicker(5 * time.Millisecond)
	b := packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size, 1, test.id_test_type, test.test_type)

	for range SolveDelayTicker.C {

		go test.receiveMessagesDelay(detectPackDelay, c, ifi.MTU, test.id_test_type, len(b))

		select {
		case detect := <-detectPackDelay:

			if detect > 0 {
				number++
				delay = delay + detect
				if number == 1 {
					delayMax = detect
					delayMin = detect
				} else {
					if delayMax < detect {
						delayMax = detect
					}
					if (delayMin > detect) && (float32(detect) > (float32(delay)/float32(number))*0.70) {
						delayMin = detect
					}
				}

			}
		default:
			if time.Since(timeStart) > time.Duration(test.count_packs)*time.Second {
				if verboseLogs {
					fmt.Println(" --> Number = ", number)
				}
				if verboseLogs {
					fmt.Println(" --> Size = ", size)
				}

				floatDelay = (float32(delay) / float32(number)) * 1000000.0 / float32(math.Pow(2, 32))
				floatDelayMax = float32(delayMax) * 1000000.0 / float32(math.Pow(2, 32))
				floatDelayMin = float32(delayMin) * 1000000.0 / float32(math.Pow(2, 32))
				/*
					if verboseLogs {
						fmt.Println(" --> Delay = ", floatDelay)
					}
					if verboseLogs {
						fmt.Println(" --> DelayMax = ", floatDelayMax)
					}
					if verboseLogs {
						fmt.Println(" --> DelayMin = ", floatDelayMin)
					}

					floatDelay = (float32(delay) / float32(number))
					floatDelayMax = float32(delayMax)
					floatDelayMin = float32(delayMin)
				*/
				if verboseLogs {
					fmt.Println(" --> Delay = ", floatDelay)
				}
				if verboseLogs {
					fmt.Println(" --> DelayMax = ", floatDelayMax)
				}
				if verboseLogs {
					fmt.Println(" --> DelayMin = ", floatDelayMin)
				}
				switch size {
				case 64:
					test.rez_64 = floatDelay
					test.rez_64_max = floatDelayMax
					test.rez_64_min = floatDelayMin
				case 128:
					test.rez_128 = floatDelay
					test.rez_128_max = floatDelayMax
					test.rez_128_min = floatDelayMin
				case 256:
					test.rez_256 = floatDelay
					test.rez_256_max = floatDelayMax
					test.rez_256_min = floatDelayMin
				case 512:
					test.rez_512 = floatDelay
					test.rez_512_max = floatDelayMax
					test.rez_512_min = floatDelayMin
				case 1024:
					test.rez_1024 = floatDelay
					test.rez_1024_max = floatDelayMax
					test.rez_1024_min = floatDelayMin
				case 1280:
					test.rez_1280 = floatDelay
					test.rez_1280_max = floatDelayMax
					test.rez_1280_min = floatDelayMin
				case 1500:
					test.rez_1518 = floatDelay
					test.rez_1518_max = floatDelayMax
					test.rez_1518_min = floatDelayMin
				}

				quit <- 1
				return
			}
		}
	}

}

func (test *testDelay) receiveMessagesDelay(catchDetect chan int64, c net.PacketConn, mtu int, test_type uint16, packetSize int) {
	var f ethernet.Frame
	b := make([]byte, mtu)
	cc := 0
	var t_ips [2]byte
	t_ips[1] = byte(test_type & 0xFF)
	t_ips[0] = byte((test_type >> 8) & 0xFF)
	start := time.Now()
	quit := make(chan int64, 10)
	c.SetReadDeadline(start.Add(time.Microsecond * 3000))
	//ExitLoop:
	for {
		select {
		case key := <-quit:
			catchDetect <- key
			return
			//break ExitLoop
		default:

			n, _, err := c.ReadFrom(b)
			cc++
			if err != nil {
				if err.Error() == "resource temporarily unavailable" {

				}
				//log.Fatalf("failed to receive message: %v", err)
				c.SetReadDeadline(start.Add(time.Hour * 24))
				quit <- 0
				runtime.Gosched()
				continue
			}

			if time.Since(start) > (time.Millisecond * 1000) {
				quit <- 0
				continue
			}

			if (n) != packetSize {
				continue
			}

			delta_nano := int64((2208988800) * 1000000000)
			t_time := int64(float64(time.Now().UnixNano()-delta_nano) * float64(math.Pow(2, 32)/1000000000))
			t_time = t_time & int64(0xFFFFFFFFFFFFFF)

			//n, addr, err := c.ReadFrom(b)
			// Unpack Ethernet II frame into Go representation.
			if err := (&f).UnmarshalBinary(b[:n]); err != nil {
				if verboseLogs {
					fmt.Printf("failed to unmarshal ethernet frame: %v", err)
				}
				continue
			}

			var ips [4]byte
			copy(ips[:], (test.ipdst1).To4())

			if (len(f.Payload) >= 52) && (f.Payload[20] == 0xFC) && (bytes.Equal(f.Payload[12:16], ips[:]) == true) && (bytes.Equal(f.Payload[50:52], t_ips[:]) == true) {

				/*	fmt.Printf("\n\n--=Packet DETECT!!!=--\n")

					if verboseLogs {
						fmt.Printf("size     %x \n", b[2:4])
					}

					if verboseLogs {
						fmt.Printf("ip sourse %v.%v.%v.%v \n", f.Payload[12], f.Payload[13], f.Payload[14], f.Payload[15])
					}
					if verboseLogs {
						fmt.Printf("ip dst    %v.%v.%v.%v \n", f.Payload[16], f.Payload[17], f.Payload[18], f.Payload[19])
					}

					if verboseLogs {
						fmt.Printf("ip SFP2   %v.%v.%v.%v \n", f.Payload[21], f.Payload[22], f.Payload[23], f.Payload[24])
					}

					if verboseLogs {
						fmt.Printf("time marker_SFP1_1 :   %x \n", f.Payload[25:32])
					}
					if verboseLogs {
						fmt.Printf("time marker_SFP2   :   %x \n", f.Payload[32:39])
					}
					if verboseLogs {
						fmt.Printf("time marker_SFP1_2 :   %x \n", f.Payload[39:46])
					}
					if verboseLogs {
						fmt.Println(" --== End Packet ==--")
					}
					//*/
				var markerSFP11, markerSFP12, markerSFP2 int64
				var ind uint

				for ind = 0; ind < 7; ind++ {
					markerSFP11 = markerSFP11 + int64(f.Payload[31-ind])<<(8*ind)
					markerSFP2 = markerSFP2 + int64(f.Payload[38-ind])<<(8*ind)
					markerSFP12 = markerSFP12 + int64(f.Payload[45-ind])<<(8*ind)
				}

				//floatDelay := (float32(markerSFP12 - markerSFP11)) * 1000000.0 / float32(math.Pow(2, 32))

				/*
					if floatDelay < 25 {

						if verboseLogs {
							fmt.Printf("size     %x \n", b[2:4])
						}
						if verboseLogs {
							fmt.Printf("Packet: %x ", f.Payload)
						}
						if verboseLogs {
							fmt.Printf("ip sourse %v.%v.%v.%v \n", f.Payload[12], f.Payload[13], f.Payload[14], f.Payload[15])
						}
						if verboseLogs {
							fmt.Printf("ip dst    %v.%v.%v.%v \n", f.Payload[16], f.Payload[17], f.Payload[18], f.Payload[19])
						}

						if verboseLogs {
							fmt.Printf("ip SFP2   %v.%v.%v.%v \n", f.Payload[21], f.Payload[22], f.Payload[23], f.Payload[24])
						}

						if verboseLogs {
							fmt.Printf("time marker_SFP1_1 :   %x \n", f.Payload[25:32])
						}
						if verboseLogs {
							fmt.Printf("time marker_SFP2   :   %x \n", f.Payload[32:39])
						}
						if verboseLogs {
							fmt.Printf("time marker_SFP1_2 :   %x \n", f.Payload[39:46])
						}
						if verboseLogs {
							fmt.Println(" delay = ", floatDelay)
						}
						if verboseLogs {
							fmt.Println(" --== End Packet ==--")
						}
					}
				*/
				if (test.test_type == 2) && (markerSFP2 == 0) {
					continue
				}

				if test.test_type == 1 {
					quit <- (markerSFP12 - markerSFP11)
				} else {
					quit <- (t_time - markerSFP2)
				}

			}
			//else {
			//	//fmt.Printf("\n\n\r[%s] %v %x", addr.String(), len(f.Payload), f.Payload[:25])

		}

	}
}
