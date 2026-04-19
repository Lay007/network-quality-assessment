package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/raw"
	"golang.org/x/net/bpf"
	"math"
	"net"
	"runtime"

	"sync/atomic"
	//	"runtime"
	"time"
)

func TestY1564(id int, net_interface_name string) {
	if verboseLogs {
		fmt.Println("Y.1564 test started")
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
	db.Exec("UPDATE test_y1564 SET status=?, datatime=? WHERE id=?", time.Now().Format("2006-01-02 15:04:05"), 2, id)
	ifi, err := net.InterfaceByName(net_interface_name)
	if err != nil {

		//	log.Fatalf("failed to find interface %q: %v", net_interface_name, err)
		db.Exec("UPDATE test_y1564 SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
		db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 5, id, "Network interface access failed: "+net_interface_name)
		db.Close()
		return
	}

	row, err := db.Query("select id, test_type, module_first, module_second, block_size, ToS, VLAN_priority, CIR, EIR, TP, period, step_count, max_FTD, max_FVD, max_FLR, status from test_y1564 where id=?", id)
	if err != nil {
		db.Exec("UPDATE test_y1564 SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
		db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 5, id, "Database query failed")
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
	test := new(testY1564)
	err = row.Scan(&test.id, &test.test_type, &test.module_first, &test.module_second, &test.block_size, &test.ToS, &test.VLAN_priority, &test.CIR, &test.EIR, &test.TP, &test.period, &test.step_count, &test.max_FTD, &test.max_FVD, &test.max_FLR, &test.status)
	if err != nil {
		db.Exec("UPDATE test_y1564 SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
		db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 5, id, "Database result parsing failed")
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
		db.Exec("UPDATE test_y1564 SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
		db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 5, id, "Database query failed")
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
			db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 5, id, "Database result parsing failed")
			db.Exec("UPDATE test_y1564 SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
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
	row, err = db.Query("SELECT module_first, module_second FROM test_y1564 WHERE id=?", id)
	if err != nil {
		db.Exec("UPDATE test_y1564 SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
		db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 5, id, "Database query failed")
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
			db.Exec("UPDATE test_y1564 SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
			db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 5, id, "Database result parsing failed")
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
			db.Exec("UPDATE test_y1564 SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
			db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 5, id, "Database query failed")
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
				db.Exec("UPDATE test_y1564 SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
				db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 5, id, "Database result parsing failed")
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
			db.Exec("UPDATE test_y1564 SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
			db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 5, id, "Database query failed")
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
				db.Exec("UPDATE test_y1564 SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
				db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 5, id, "Database result parsing failed")
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
			db.Exec("UPDATE test_y1564 SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
			db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 5, id, "Database result parsing failed")
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
		test_mac = 0
		for row_mac.Next() {
			//err = row_mac.Scan(&mac_dst_str)
			err = row_mac.Scan(&test_mac)
			if err != nil {

				db.Exec("UPDATE test_y1564 SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
				db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 5, id, "Database result parsing failed")
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
			db.Exec("UPDATE test_y1564 SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
			db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 5, id, "Database query failed")
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

				db.Exec("UPDATE test_y1564 SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
				db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 5, id, "Database result parsing failed")
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
	test.status = 2
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
	db.Exec("UPDATE test_y1564 SET status=?, datetime_start=? WHERE id=?", 2, time.Now().Format("2006-01-02 15:04:05"), id)

	if testPing(ipdst_1sfpsla_str) > 0 || testPing(ipdst_2sfpsla_str) > 0 {
		db.Exec("UPDATE test_y1564 SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
		db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 5, id, "Ping test failed")
		db.Close()

		return
	}

	if check_SNMP(ipdst_1sfpsla_str) > 0 || check_SNMP(ipdst_2sfpsla_str) > 0 {
		db.Exec("UPDATE test_y1564 SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
		db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 5, id, "SNMP test failed")
		db.Close()

		return
	}

	addr := &raw.Addr{
		HardwareAddr: ethernet.Broadcast,
	}
	connectTestSFP, err := raw.ListenPacket(ifi, etherType, nil)
	rez := findSFP(connectTestSFP, addr, ipsrcstr, ipdst_1sfpsla_str, ipdst_2sfpsla_str, test.mac_src, test.mac_dst, test.mac_dst2, test.id_test_type, test.test_type, 1024, int64(1024*8*1000/test.CIR))

	if verboseLogs {
		fmt.Printf("\n Rez find : %X \n", rez)
	}
	if (rez & 0xFFF) == 0x999 {
		db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 5, id, " Maximum throughput - 1 Gbit/s")

	}
	if (rez & 0xFFF) == 0x100 {
		db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 5, id, " Maximum throughput - 100 Mbit/s")

	}
	if (rez & 0xFFF) == 0x10 {
		db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 5, id, " Maximum throughput - 10 Mbit/s")

	}
	if (rez & 0xF000) == 0x0 {
		db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 5, id, " Serial module topology. Module order is correct")
	}
	if (rez & 0xF000) == 0x1000 {
		db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 5, id, " Serial module topology. Module order is incorrect; swapping during the test")
	}
	if (rez & 0xF000) == 0x2000 {
		db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 5, id, " Star topology. Load is balanced")
	}
	if (rez & 0xF000) == 0x3000 {
		db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 5, id, " Star topology. Load is unbalanced; module order is correct")
	}
	if (rez & 0xF000) == 0x4000 {
		db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 5, id, " Star topology. Load is unbalanced; module order is incorrect; swapping during the test")
	}
	connectTestSFP.Close()
	if rez == 0 {
		if verboseLogs {
			fmt.Println("Error test SFP connect")
		}
		db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 5, id, "Express path test failed")

		return
	}
	if ((rez & 0xF000) == 0x1000) || ((rez & 0xF000) == 0x4000) {

		tmp := ipdst_1sfpsla_str
		ipdst_1sfpsla_str = ipdst_2sfpsla_str
		ipdst_2sfpsla_str = tmp

		test.ipsrc = net.ParseIP(ipsrcstr)
		test.ipdst1 = net.ParseIP(ipdst_1sfpsla_str)
		test.ipdst2 = net.ParseIP(ipdst_2sfpsla_str)

		//	db, err = openDB()
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

	yest, _ := db.Query("SELECT EXISTS(SELECT id FROM test_y1564 WHERE id=?)", id)
	yest.Next()
	t_y := 0
	yest.Scan(&t_y)
	if t_y != 1 {
		return
	}

	//db.Close()

	test.id_test_type = 0xE000 + (uint16(id) & 0x1FFF)
	test.id_test_type_temp = 0x2000 + (uint16(id) & 0x1FFF)
	test.net_interface_name = net_interface_name
	test.mac_src = ifi.HardwareAddr

	counter := make(chan uint64, 1)
	counterRes := make(chan uint64, 1)
	quit := make(chan int64, 2)

	//size_p := 64

	var ToS_tag uint8
	ToS_tag = ((uint8(0x7 & test.VLAN_priority)) << 5) + uint8(test.ToS<<1)
	test.ToS = ToS_tag
	thr_s := make([]int, test.step_count)

	switch test.step_count {
	case 1:
		thr_s[0] = test.CIR
		break
	case 2:
		thr_s[0] = test.CIR / 2
		thr_s[1] = test.CIR
		break
	case 3:
		thr_s[0] = test.CIR / 2
		thr_s[1] = int(float32(test.CIR) * 0.75)
		thr_s[2] = test.CIR
		break
	case 4:
		thr_s[0] = test.CIR / 2
		thr_s[1] = int(float32(test.CIR) * 0.75)
		thr_s[2] = int(float32(test.CIR) * 0.9)
		thr_s[3] = test.CIR
		break
	}
	//	b := packetFormY1546(ToS_tag, test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, test.block_size, 0, test.id_test_type, test.test_type)

	go test.genFramesY1564(thr_s[0], counter, counterRes)
	delay, jitter := test.getMetricsY1564(quit)
	time.Sleep(time.Second * 2)
	<-quit
	PacketsRx := <-counter
	PacketsRxRes := <-counterRes
	test.rez_IR_s1 = float32(PacketsRx+PacketsRxRes) * float32(test.block_size) * 8.0 / (float32(test.period) * 1000000.0)
	test.rez_FTD_s1 = delay
	test.rez_FVD_s1 = jitter
	test.rez_FLR_s1 = float32(PacketsRxRes-uint64(test.numberRx)) / float32(PacketsRxRes)

	db.Exec("UPDATE test_y1564 SET rez_IR_s1=?, rez_FTD_s1=?,rez_FVD_s1=?,rez_FLR_s1=?, status=? WHERE id=?",
		test.rez_IR_s1, test.rez_FTD_s1, test.rez_FVD_s1, test.rez_FLR_s1, test.status, id)

	PacketsRx = 0
	test.numberRx = 0

	if test.step_count > 1 {

		yest, _ = db.Query("SELECT EXISTS(SELECT id FROM test_y1564 WHERE id=?)", id)
		yest.Next()
		t_y = 0
		yest.Scan(&t_y)
		if t_y != 1 {
			return
		}

		go test.genFramesY1564(thr_s[1], counter, counterRes)
		delay, jitter := test.getMetricsY1564(quit)
		time.Sleep(time.Second * 2)
		<-quit
		PacketsRx := <-counter
		PacketsRxRes := <-counterRes
		test.rez_IR_s2 = float32(PacketsRx+PacketsRxRes) * float32(test.block_size) * 8.0 / (float32(test.period) * 1000000.0)
		test.rez_FTD_s2 = delay
		test.rez_FVD_s2 = jitter
		test.rez_FLR_s2 = float32(PacketsRxRes-uint64(test.numberRx)) / float32(PacketsRxRes)

		db.Exec("UPDATE test_y1564 SET rez_IR_s2=?, rez_FTD_s2=?,rez_FVD_s2=?,rez_FLR_s2=?, status=? WHERE id=?",
			test.rez_IR_s2, test.rez_FTD_s2, test.rez_FVD_s2, test.rez_FLR_s2, test.status, id)

		PacketsRx = 0
		test.numberRx = 0

	}

	if test.step_count > 2 {
		yest, _ = db.Query("SELECT EXISTS(SELECT id FROM test_y1564 WHERE id=?)", id)
		yest.Next()
		t_y = 0
		yest.Scan(&t_y)
		if t_y != 1 {
			return
		}

		go test.genFramesY1564(thr_s[2], counter, counterRes)
		delay, jitter := test.getMetricsY1564(quit)
		time.Sleep(time.Second * 2)
		<-quit
		PacketsRx := <-counter
		PacketsRxRes := <-counterRes
		test.rez_IR_s3 = float32(PacketsRx+PacketsRxRes) * float32(test.block_size) * 8.0 / (float32(test.period) * 1000000.0)
		test.rez_FTD_s3 = delay
		test.rez_FVD_s3 = jitter
		test.rez_FLR_s3 = float32(PacketsRxRes-uint64(test.numberRx)) / float32(PacketsRxRes)

		db.Exec("UPDATE test_y1564 SET rez_IR_s3=?, rez_FTD_s3=?,rez_FVD_s3=?,rez_FLR_s3=?, status=? WHERE id=?",
			test.rez_IR_s3, test.rez_FTD_s3, test.rez_FVD_s3, test.rez_FLR_s3, test.status, id)

		PacketsRx = 0
		test.numberRx = 0

	}

	if test.step_count > 3 {
		yest, _ = db.Query("SELECT EXISTS(SELECT id FROM test_y1564 WHERE id=?)", id)
		yest.Next()
		t_y = 0
		yest.Scan(&t_y)
		if t_y != 1 {
			return
		}
		go test.genFramesY1564(thr_s[3], counter, counterRes)
		delay, jitter := test.getMetricsY1564(quit)
		time.Sleep(time.Second * 2)
		<-quit
		PacketsRx := <-counter
		PacketsRxRes := <-counterRes
		test.rez_IR_s4 = float32(PacketsRx+PacketsRxRes) * float32(test.block_size) * 8.0 / (float32(test.period) * 1000000.0)
		test.rez_FTD_s4 = delay
		test.rez_FVD_s4 = jitter
		test.rez_FLR_s4 = float32(PacketsRxRes-uint64(test.numberRx)) / float32(PacketsRxRes)

		db.Exec("UPDATE test_y1564 SET rez_IR_s4=?, rez_FTD_s4=?,rez_FVD_s4=?,rez_FLR_s4=?, status=? WHERE id=?",
			test.rez_IR_s4, test.rez_FTD_s4, test.rez_FVD_s4, test.rez_FLR_s4, test.status, id)

		PacketsRx = 0
		test.numberRx = 0

	}

	runtime.Gosched()
	yest, _ = db.Query("SELECT EXISTS(SELECT id FROM test_y1564 WHERE id=?)", id)
	yest.Next()
	t_y = 0
	yest.Scan(&t_y)
	if t_y != 1 {
		return
	}
	go test.genFramesY1564(test.CIR+test.EIR, counter, counterRes)
	delay, jitter = test.getMetricsY1564(quit)
	time.Sleep(time.Second * 2)
	<-quit
	PacketsRx = <-counter
	PacketsRxRes = <-counterRes
	test.rez_IR_eir = float32(PacketsRx+PacketsRxRes) * float32(test.block_size) * 8.0 / (float32(test.period) * 1000000.0)
	test.rez_FTD_eir = delay
	test.rez_FVD_eir = jitter
	test.rez_FLR_eir = float32(PacketsRxRes-uint64(test.numberRx)) / float32(PacketsRxRes)

	db.Exec("UPDATE test_y1564 SET rez_IR_eir=?, rez_FTD_eir=?,rez_FVD_eir=?,rez_FLR_eir=?, status=? WHERE id=?",
		test.rez_IR_eir, test.rez_FTD_eir, test.rez_FVD_eir, test.rez_FLR_eir, test.status, id)

	PacketsRx = 0
	test.numberRx = 0

	runtime.Gosched()
	yest, _ = db.Query("SELECT EXISTS(SELECT id FROM test_y1564 WHERE id=?)", id)
	yest.Next()
	t_y = 0
	yest.Scan(&t_y)
	if t_y != 1 {
		return
	}
	go test.genFramesY1564(test.CIR+test.TP, counter, counterRes)
	delay, jitter = test.getMetricsY1564(quit)
	time.Sleep(time.Second * 2)
	<-quit
	PacketsRx = <-counter
	PacketsRxRes = <-counterRes
	test.rez_IR_tp = float32(PacketsRx+PacketsRxRes) * float32(test.block_size) * 8.0 / (float32(test.period) * 1000000.0)
	test.rez_FTD_tp = delay
	test.rez_FVD_tp = jitter
	test.rez_FLR_tp = float32(PacketsRxRes-uint64(test.numberRx)) / float32(PacketsRxRes)
	test.status = 3

	db.Exec("UPDATE test_y1564 SET rez_IR_tp=?, rez_FTD_tp=?,rez_FVD_tp=?,rez_FLR_tp=?, datetime_end=?, status=? WHERE id=?",
		test.rez_IR_tp, test.rez_FTD_tp, test.rez_FVD_tp, test.rez_FLR_tp, time.Now().Format("2006-01-02 15:04:05"), test.status, id)

	PacketsRx = 0
	test.numberRx = 0

	//	db, err = openDB()
	//	if err != nil {
	//		db.Close()

	//		return
	//	}
	//db.Exec("UPDATE test_throughput SET rez_64=?,rez_128=?,rez_256=?,rez_512=?,rez_1024=?,rez_1280=?, rez_1518=?,rez_4096=?,rez_9000=?,status=? WHERE id=?", test.rez_64, test.rez_128, test.rez_256, test.rez_512, test.rez_1024, test.rez_1280, test.rez_1518, test.rez_4096, test.rez_9000, test.status, id)
	//	db.Exec("UPDATE test_y1564 SET rez_IR_s1=?, rez_FTD_s1=?,rez_FVD_s1=?,rez_FLR_s1=?,rez_IR_s2=?, rez_FTD_s2=?,rez_FVD_s2=?,rez_FLR_s2=?, rez_IR_s3=?, rez_FTD_s3=?,rez_FVD_s3=?,rez_FLR_s3=?,rez_IR_s4=?, rez_FTD_s4=?,rez_FVD_s4=?,rez_FLR_s4=?, rez_IR_eir=?, rez_FTD_eir=?,rez_FVD_eir=?,rez_FLR_eir=?, rez_IR_tp=?, rez_FTD_tp=?,rez_FVD_tp=?,rez_FLR_tp=?, datetime_end=?, status=? WHERE id=?", test.rez_IR_s1, test.rez_FTD_s1, test.rez_FVD_s1, test.rez_FLR_s1, test.rez_IR_s2, test.rez_FTD_s2, test.rez_FVD_s2, test.rez_FLR_s2, test.rez_IR_s3, test.rez_FTD_s3, test.rez_FVD_s3, test.rez_FLR_s3, test.rez_IR_s4, test.rez_FTD_s4, test.rez_FVD_s4, test.rez_FLR_s4, test.rez_IR_eir, test.rez_FTD_eir, test.rez_FVD_eir, test.rez_FLR_eir, test.rez_IR_tp, test.rez_FTD_tp, test.rez_FVD_tp, test.rez_FLR_tp, time.Now().Format("2006-01-02 15:04:05"), test.status, id)

	db.Close()
}

func packetFormY1546(ToS uint8, ipsrc net.IP, ipdst1 net.IP, ipdst2 net.IP, mac_src []byte, mac_dst []byte, size int, number uint32, test_type uint16, testWay int) []byte {
	ip := iphdr{
		vhl:   0x45,
		tos:   ToS,
		id:    0x0000, // the kernel overwrites id if it is zero
		off:   0,
		ttl:   0xFF,
		proto: 0x5E,
	}
	copy(ip.src[:], ipsrc.To4())
	copy(ip.dst[:], ipdst1.To4())
	sfpdat := sfpsla{
		//	id: 0xFC,
	}
	copy(sfpdat.dst[:], ipdst2.To4())

	sfpdat.number = number
	sfpdat.test_type = test_type
	var ind uint
	//ip.iplen = uint16(20 + 26 + 4)
	if testWay == 2 {
		delta_nano := int64((2208988800) * 1e9)
		t_time := int64(float64(time.Now().UnixNano()-delta_nano) * float64(math.Pow(2, 32)/float64(1e9)))
		t_time = t_time & int64(0xFFFFFFFFFFFFFF)

		for ind = 0; ind < 7; ind++ {
			sfpdat.merkertime2[6-ind] = byte((t_time >> (8 * ind)) & 0xFF)
		}
		copy(sfpdat.dst[:], ipdst1.To4())
		sfpdat.id = 0xFB
	}
	payloadAdd := make([]byte, 0)
	if size > 66 {
		payloadAdd = make([]byte, size-66)
		for h := 0; h < len(payloadAdd); h++ {
			payloadAdd[h] = byte(h)
		}
	}
	//	ip.iplen = uint16(unsafe.Sizeof(ip) + unsafe.Sizeof(sfpdat) + unsafe.Sizeof(payloadAdd))
	if size > 66 {
		ip.iplen = uint16(size - 14)
	} else {
		ip.iplen = 52
	}
	ip.checksum()

	var bin_buf bytes.Buffer
	binary.Write(&bin_buf, binary.BigEndian, ip)
	binary.Write(&bin_buf, binary.BigEndian, sfpdat)
	binary.Write(&bin_buf, binary.BigEndian, payloadAdd)

	msg := bin_buf.Bytes()
	f := &ethernet.Frame{
		//Destination: ethernet.Broadcast,
		Destination: mac_dst,
		//Destination: []byte{0x5A, 0x11, 0x22, 0x33, 0x44, 0x01},
		//Destination: []byte{0x64, 0xD1, 0x54, 0x17, 0xF6, 0x82},
		Source:    mac_src,
		EtherType: 0x0800,
		Payload:   []byte(msg),
	}

	b, err := f.MarshalBinary()
	if err != nil {
		//log.Fatalf("failed to marshal ethernet frame: %v", err)

		if verboseLogs {
			fmt.Println(" -!! Error !!-")
		}
		if verboseLogs {
			fmt.Println(err)
		}
		if verboseLogs {
			fmt.Println(" ----=====----")
		}
		return []byte{}
	}
	return b
}

/*
	func (test *testY1564) genFramesY1564(thr int, counter chan int64) int64 {
		time.Sleep(1*time.Millisecond)
		number := uint32(0)
		period_nano := int64(test.block_size * 8 * 1000 / thr)

		if verboseLogs {
			fmt.Printf("\n   period_nano = %d \n   thr = %d\n", period_nano, thr)
		}

		var counter_rez int64
		ifi, err := net.InterfaceByName(test.net_interface_name)
		c, err := raw.ListenPacket(ifi, etherType, nil)
		if err != nil {
			if verboseLogs {
				fmt.Println("failed to listen: %v", err)
			}
			if verboseLogs {
				fmt.Println(" -!! Error !!-")
			}
			if verboseLogs {
				fmt.Println(err)
			}
			if verboseLogs {
				fmt.Println(" ----=====----")
			}
			return 0
		}
		addr := &raw.Addr{
			HardwareAddr: ethernet.Broadcast,
		}

		star_gen := time.Now()
		var rez_time int64
		ticker := time.NewTicker(time.Duration(period_nano))
		done := make(chan bool)
		b := packetFormY1546(test.ToS, test.ipsrc, test.ipdst1, test.ipdst2, ifi.HardwareAddr, test.mac_dst, test.block_size, number, test.id_test_type, test.test_type)

		go func() {
			for {
				select {
				case <-done:
					rez_time = (int64)(time.Since(star_gen))
					return
				case <-ticker.C:
					number++
				//	b := packetFormY1546(test.ToS, test.ipsrc, test.ipdst1, test.ipdst2, ifi.HardwareAddr, test.mac_dst, test.block_size, number, test.id_test_type, test.test_type)

				//	for { //	time.Sleep(time.Millisecond * 1)
					//	n, err := c.WriteTo(b, addr)
						c.WriteTo(b, addr)
				//		if n == len(b) && err == nil {
				//			break
				//		}
				//	}
					counter_rez = counter_rez + 1
				}
			}
		}()
		time.Sleep(time.Duration(test.period) * time.Second)
		ticker.Stop()
		done <- true
		time.Sleep(1 * time.Second)
		if verboseLogs {
			fmt.Println("Packed send - ", counter_rez)
		}
		counter <- counter_rez
		return rez_time
	}
*/
func (test *testY1564) genFramesY1564(thr int, counter chan uint64, counterRes chan uint64) int64 {
	time.Sleep(time.Millisecond * 10)
	ifi, _ := net.InterfaceByName(test.net_interface_name)
	b := packetFormY1546(test.ToS, test.ipsrc, test.ipdst1, test.ipdst2, ifi.HardwareAddr, test.mac_dst, test.block_size, 1, test.id_test_type, test.test_type)
	bTemp := packetFormY1546(test.ToS, test.ipsrc, test.ipdst1, test.ipdst2, ifi.HardwareAddr, test.mac_dst, test.block_size, 1, test.id_test_type_temp, test.test_type)
	//counter := make(chan int64, 1)
	rez_time := genSocket(ifi.Index, b, bTemp, test.period, thr, counter, counterRes)

	return rez_time
}
func (test *testY1564) getMetricsY1564(quit chan int64) (float32, float32) {

	var delay, delayMax, delayMin int64
	var jitter, delayPrep int64
	var floatDelay, floatDelayMax, floatDelayMin, floatJitter float32

	detectPackDelay := make(chan int64, 10)
	number := 0

	timeStart := time.Now()

	ifi, err := net.InterfaceByName(test.net_interface_name)
	if err != nil {
		if verboseLogs {
			fmt.Printf("failed to find interface %q: %v\n", test.net_interface_name, err)
		}
		quit <- 1
		return 0, 0
	}

	var netConf *raw.Config = new(raw.Config)

	(*netConf).Filter, _ = bpf.Assemble([]bpf.Instruction{
		bpf.LoadAbsolute{Off: 34, Size: 1},
		bpf.JumpIf{Cond: bpf.JumpNotEqual, Val: 0xFC, SkipTrue: 3},
		bpf.LoadAbsolute{Off: 64, Size: 2},
		bpf.JumpIf{Cond: bpf.JumpNotEqual, Val: uint32(test.id_test_type), SkipTrue: 1},
		// Verdict is "send up to 4k of the packet to userspace."
		bpf.RetConstant{Val: 4096},
		// Verdict is "ignore packet."
		bpf.RetConstant{Val: 0},
	})

	c, err := raw.ListenPacket(ifi, etherType, netConf)
	quit_receive := make(chan int64, 2)
	go (*test).receivePacketsY(ifi.MTU, quit_receive, test.id_test_type)

	SolveDelayTicker := time.NewTicker(1 * time.Millisecond)
	//	b := packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, test.block_size, 1, test.id_test_type, test.test_type)
	go test.receiveMessagesDelay(detectPackDelay, c, ifi.MTU, test.id_test_type, test.block_size)

	for range SolveDelayTicker.C {
		//	for	 {
		//	go test.receiveMessagesDelay(detectPackDelay, c, ifi.MTU, test.id_test_type, test.block_size)

		select {
		case detect := <-detectPackDelay:

			if detect > 0 {
				number++
				delay = delay + detect
				if number > 1 {
					jitter += int64(math.Abs(float64(detect - delayPrep)))
				}
				delayPrep = detect
				if number == 1 {
					delayMax = detect
					delayMin = detect
				} else {
					if delayMax < detect {
						delayMax = detect
					}
					if delayMin > detect {
						delayMin = detect
					}
				}

			}
		default:
			if time.Since(timeStart) > time.Duration(time.Duration(test.period)*time.Second) {
				if verboseLogs {
					fmt.Println("\n --> Number = ", number)
				}
				if verboseLogs {
					fmt.Println(" --> Size = ", test.block_size)
				}

				if number > 0 {
					floatDelay = (float32(delay) / float32(number)) * 1000000.0 / float32(math.Pow(2, 32))
					floatJitter = (float32(jitter) / float32(number)) * 1000000.0 / float32(math.Pow(2, 32))

				}
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

				/*
					switch test.block_size {
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
				*/
				time.Sleep(1000 * time.Millisecond)
				quit_receive <- 1
				quit <- 0
				//	return floatDelay, float32(math.Max(float64(floatDelayMax-floatDelay), float64(floatDelay-floatDelayMin)))
				return floatDelay, floatJitter
			}
		}
	}
	return 0, 0

}

func (test *testY1564) receiveMessagesDelay(catchDetect chan int64, c net.PacketConn, mtu int, test_type uint16, packetSize int) {
	var f ethernet.Frame
	b := make([]byte, mtu)
	cc := 0
	var t_ips [2]byte
	t_ips[1] = byte(test_type & 0xFF)
	t_ips[0] = byte((test_type >> 8) & 0xFF)
	start := time.Now()
	quit := make(chan int64, 10)
	c.SetReadDeadline(start.Add(time.Second * time.Duration(test.period)))
	if packetSize == 64 {
		packetSize += 2
	}
	//ExitLoop:
	for {
		select {
		case key := <-quit:
			catchDetect <- key
			return
			//break ExitLoop
		default:

			n, _, err := c.ReadFrom(b)
			at_time := time.Now()
			cc++
			if err != nil {
				if err.Error() == "i/o timeout" {
					c.SetReadDeadline(start.Add(time.Hour * 24))
					quit <- 0
					runtime.Gosched()
				}
				//log.Fatalf("failed to receive message: %v", err)

				//quit <- 1
				//runtime.Gosched()
				continue
				//quit <- 0

			}

			if (n) != packetSize {
				continue
			}

			delta_nano := int64((2208988800) * 1000000000)

			t_time := int64(float64(at_time.UnixNano()-delta_nano) * float64(math.Pow(2, 32)/1000000000))
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
				//	(*test).numberRx++
				//if ((*test).numberRx % 100) != 0 {
				//	return
				//}
				var markerSFP11, markerSFP12, markerSFP2 int64
				var ind uint

				for ind = 0; ind < 7; ind++ {
					markerSFP11 = markerSFP11 + int64(f.Payload[31-ind])<<(8*ind)
					markerSFP2 = markerSFP2 + int64(f.Payload[38-ind])<<(8*ind)
					markerSFP12 = markerSFP12 + int64(f.Payload[45-ind])<<(8*ind)
				}
				if markerSFP2 == 0 {
					continue
				}
				if test.test_type == 1 {
					catchDetect <- (markerSFP12 - markerSFP11)
				} else {
					catchDetect <- (t_time - markerSFP2)
				}
				runtime.Gosched()
			}
		}

	}
}
func (test *testY1564) receivePacketsY(mtu int, quit chan int64, t_type uint16) { //, counter chan<- int) {

	ifi, err := net.InterfaceByName(test.net_interface_name)
	if err != nil {
		if verboseLogs {
			fmt.Printf("failed to find interface %q: %v\n", test.net_interface_name, err)
		}
		return
	}

	var netConf *raw.Config = new(raw.Config)

	(*netConf).Filter, _ = bpf.Assemble([]bpf.Instruction{
		bpf.LoadAbsolute{Off: 34, Size: 1},
		bpf.JumpIf{Cond: bpf.JumpNotEqual, Val: 0xFC, SkipTrue: 3},
		bpf.LoadAbsolute{Off: 64, Size: 2},
		bpf.JumpIf{Cond: bpf.JumpNotEqual, Val: uint32(test.id_test_type), SkipTrue: 1},
		// Verdict is "send up to 4k of the packet to userspace."
		bpf.RetConstant{Val: 4096},
		// Verdict is "ignore packet."
		bpf.RetConstant{Val: 0},
	})

	c, err := raw.ListenPacket(ifi, etherType, netConf)

	b := make([]byte, mtu)
	var ips [4]byte
	copy(ips[:], (test.ipdst1).To4())

	var t_ips [2]byte
	t_ips[1] = byte(test.id_test_type & 0xFF)
	t_ips[0] = byte((test.id_test_type >> 8) & 0xFF)
	start := time.Now()
	c.SetReadDeadline(start.Add(time.Second * time.Duration(1+test.period)))
	if verboseLogs {
		fmt.Println("Start receive: ", time.Now())
	}
	for {
		select {
		case <-quit:
			if verboseLogs {
				fmt.Println("End receive: ", time.Now())
			}
			if verboseLogs {
				fmt.Println("Packets receive = ", test.numberRx)
			}
			return
		default:
		}
		_, _, err := c.ReadFrom(b)
		if err != nil {
			if verboseLogs {
				fmt.Println("failed to receive message: ", err)
			}
			if err.Error() == "i/o timeout" {

				c.SetReadDeadline(start.Add(time.Hour * 24))
				quit <- 1
				runtime.Gosched()
			}
			continue
		}
		atomic.AddUint64(&test.numberRx, uint64(1))
	}
}
