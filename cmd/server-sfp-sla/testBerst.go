package main

import (
	"fmt"
	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/raw"
	"net"
	"time"

	"golang.org/x/net/bpf"
)

func TestBerst(id int, net_interface_name string) {
	if verboseLogs {
		fmt.Println("Burst test started")
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
	db.Exec("UPDATE test_bert SET status=?, datatime=? WHERE id=?", time.Now().Format("2006-01-02 15:04:05"), 2, id)
	ifi, err := net.InterfaceByName(net_interface_name)
	if err != nil {
		db.Close()
		//	log.Fatalf("failed to find interface %q: %v", net_interface_name, err)
		db.Exec("UPDATE test_bert SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
		return
	}

	row, err := db.Query("select id,miss_init_test, test_type, module_first, module_second, thr_begin, count_prob_packs, count_probs, status from test_bert where id=?", id)
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
	test := new(testBert)
	err = row.Scan(&test.id, &test.miss_init_test, &test.test_type, &test.module_first, &test.module_second, &test.thr_begin, &test.count_prob_packs, &test.count_probs, &test.status)
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
	mac_dst := make([]byte, 6)
	mac_dst2 := make([]byte, 6)

	row, err = db.Query("SELECT server_IP FROM global_config")
	if err != nil {
		db.Exec("UPDATE test_bert SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
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

			db.Exec("UPDATE test_bert SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
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
	row, err = db.Query("SELECT module_first, module_second FROM test_bert WHERE id=?", id)
	if err != nil {
		db.Exec("UPDATE test_bert SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
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

			db.Exec("UPDATE test_bert SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
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
			db.Exec("UPDATE test_bert SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
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

				db.Exec("UPDATE test_bert SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
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
			db.Exec("UPDATE test_bert SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
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

				db.Exec("UPDATE test_bert SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
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
		mac_dst[5] = byte(test_mac & 0xFF)
		mac_dst[4] = byte((test_mac >> 8) & 0xFF)
		mac_dst[3] = byte((test_mac >> 16) & 0xFF)
		mac_dst[2] = byte((test_mac >> 24) & 0xFF)
		mac_dst[1] = byte((test_mac >> 32) & 0xFF)
		mac_dst[0] = byte((test_mac >> 40) & 0xFF)

		row_mac, err = db.Query("SELECT mac FROM modules_sfp_sla WHERE id=?", id_sfp2)
		if err != nil {
			db.Exec("UPDATE test_bert SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
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

				db.Exec("UPDATE test_bert SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
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
		mac_dst2[5] = byte(test_mac & 0xFF)
		mac_dst2[4] = byte((test_mac >> 8) & 0xFF)
		mac_dst2[3] = byte((test_mac >> 16) & 0xFF)
		mac_dst2[2] = byte((test_mac >> 24) & 0xFF)
		mac_dst2[1] = byte((test_mac >> 32) & 0xFF)
		mac_dst2[0] = byte((test_mac >> 40) & 0xFF)

		row_ip, err = db.Query("SELECT address_ip FROM modules_sfp_sla WHERE id=?", id_sfp2)
		if err != nil {
			db.Exec("UPDATE test_bert SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
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

				db.Exec("UPDATE test_bert SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
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

	db.Exec("UPDATE test_bert SET status=?, datetime_start=? WHERE id=?", 2, time.Now().Format("2006-01-02 15:04:05"), id)

	if testPing(ipdst_1sfpsla_str) > 0 || testPing(ipdst_2sfpsla_str) > 0 {
		db.Exec("UPDATE test_bert SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
		db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 4, id, "Ping test failed")
		db.Close()
		return
	}

	if check_SNMP(ipdst_1sfpsla_str) > 0 || check_SNMP(ipdst_2sfpsla_str) > 0 {
		db.Exec("UPDATE test_bert SET status=?, datetime_end=? WHERE id=?", 4, time.Now().Format("2006-01-02 15:04:05"), id)
		db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 4, id, "SNMP test failed")
		db.Close()

		return
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

	ipsrc := net.ParseIP(ipsrcstr)
	ipdst1 := net.ParseIP(ipdst_1sfpsla_str)
	ipdst2 := net.ParseIP(ipdst_2sfpsla_str)

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
	addr := &raw.Addr{
		HardwareAddr: ethernet.Broadcast,
	}
	var number uint32
	var test_type uint16
	test_type = 0x2000 + (uint16(id) & 0x1FFF)

	var test_c testThr
	test_c.numberCounter = uint64(test.count_prob_packs)
	test_c.testID = id

	var netConf *raw.Config = new(raw.Config)

	(*netConf).Filter, _ = bpf.Assemble([]bpf.Instruction{
		bpf.LoadAbsolute{Off: 34, Size: 1},
		bpf.JumpIf{Cond: bpf.JumpNotEqual, Val: 0xFC, SkipTrue: 3},
		bpf.LoadAbsolute{Off: 64, Size: 2},
		bpf.JumpIf{Cond: bpf.JumpNotEqual, Val: uint32(test_type), SkipTrue: 1},
		// Verdict is "send up to 4k of the packet to userspace."
		bpf.RetConstant{Val: 4096},
		// Verdict is "ignore packet."
		bpf.RetConstant{Val: 0},
	})
	rez := 1
	connectTestSFP, err := raw.ListenPacket(ifi, etherType, nil)
	if test.miss_init_test == 0 {

		rez := findSFP(connectTestSFP, addr, ipsrcstr, ipdst_1sfpsla_str, ipdst_2sfpsla_str, ifi.HardwareAddr, mac_dst, mac_dst2, test_type, test.test_type, 1024, int64(1024*8*1000/(test.thr_begin)))

		if verboseLogs {
			fmt.Printf("\n Rez find : %X \n", rez)
		}
		if (rez & 0xFFF) == 0x999 {
			db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 4, id, " Maximum throughput - 1 Gbit/s")

		}
		if (rez & 0xFFF) == 0x100 {
			db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 4, id, " Maximum throughput - 100 Mbit/s")
			if test.thr_begin > 100 {
				test.thr_begin = 100
				db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 4, id, " Changing maximum throughput to 100 Mbit/s")
			}
		}
		if (rez & 0xFFF) == 0x10 {
			db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 4, id, " Maximum throughput - 10 Mbit/s")
			if test.thr_begin > 10 {
				test.thr_begin = 10
				db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 4, id, " Changing maximum throughput to 10 Mbit/s")
			}
		}
		if (rez & 0xF000) == 0x0 {
			db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 4, id, " Serial module topology. Module order is correct")
		}
		if (rez & 0xF000) == 0x1000 {
			db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 4, id, " Serial module topology. Module order is incorrect; swapping during the test")
		}
		if (rez & 0xF000) == 0x2000 {
			db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 4, id, " Star topology. Load is balanced")
		}
		if (rez & 0xF000) == 0x3000 {
			db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 4, id, " Star topology. Load is unbalanced; module order is correct")
		}
		if (rez & 0xF000) == 0x4000 {
			db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 4, id, " Star topology. Load is unbalanced; module order is incorrect; swapping during the test")
		}
	}
	//db.Close()

	yest, _ := db.Query("SELECT EXISTS(SELECT id FROM test_bert WHERE id=?)", id)
	yest.Next()
	t_y := 0
	yest.Scan(&t_y)
	if t_y != 1 {
		return
	}

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

		ipsrc = net.ParseIP(ipsrcstr)
		ipdst1 = net.ParseIP(ipdst_1sfpsla_str)
		ipdst2 = net.ParseIP(ipdst_2sfpsla_str)

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
		mac_dst[5] = byte(test_mac & 0xFF)
		mac_dst[4] = byte((test_mac >> 8) & 0xFF)
		mac_dst[3] = byte((test_mac >> 16) & 0xFF)
		mac_dst[2] = byte((test_mac >> 24) & 0xFF)
		mac_dst[1] = byte((test_mac >> 32) & 0xFF)
		mac_dst[0] = byte((test_mac >> 40) & 0xFF)

		//db.Close()

	}

	//db.Close()
	if verboseLogs {
		fmt.Println(" Rez find : ", rez)
	}

	period_test := test.count_prob_packs
	c, _ := raw.ListenPacket(ifi, etherType, nil)
	count_probs_one_test := 1000
	packet_count_step_packets := int64(period_test)

	size := 64
	if verboseLogs {
		fmt.Println("->> test.thr_begin = ", test.thr_begin)
	}
	period_nano := int64(size*8*1000000000) / (int64(test.thr_begin * 1000 * 1000))
	packet_count_start := int64(period_test)

	if verboseLogs {
		fmt.Println("->> period_nano  = ", period_nano)
	}
	if verboseLogs {
		fmt.Println("->> packet_count = ", packet_count_start)
	}

	b := packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, size, number, test_type, test.test_type)
	packet_count_OK := int64(0)

	for k := 0; k < test.count_probs; k++ {
		packet_count := packet_count_start
		for j := 0; j < count_probs_one_test; j++ {
			test_c.numberCounter = 0
			cnt := packet_count
			quit := make(chan int)
			go func() {
				time.Sleep(100 * time.Millisecond)
				for {
					cnt--
					c.WriteTo(b, addr)
					if cnt < 0 {
						time.Sleep(100 * time.Millisecond)
						quit <- 1
						break
					}
				}
			}()
			(test_c).receivePackets(c, ifi.MTU, ipdst_1sfpsla_str, quit, test_type)
			rez_count := test_c.numberCounter
			if verboseLogs {
				fmt.Println("rez_count= ", rez_count)
			}
			if test_c.numberCounter >= uint64(packet_count) {
				continue
				packet_count += int64(packet_count_step_packets)
			} else {
				break
			}
		}
		packet_count_OK += int64(test_c.numberCounter)
	}
	test.rez_64 = float32(packet_count_OK) / float32(test.count_probs)
	if verboseLogs {
		fmt.Println("->> rez_64 = ", test.rez_64)
	}

	yest, _ = db.Query("SELECT EXISTS(SELECT id FROM test_bert WHERE id=?)", id)
	yest.Next()
	t_y = 0
	yest.Scan(&t_y)
	if t_y != 1 {
		return
	}

	db.Exec("UPDATE test_bert SET rez_64=?,rez_128=?,rez_256=?,rez_512=?,rez_1024=?,rez_1280=?, rez_1518=? WHERE id=?", test.rez_64, test.rez_128, test.rez_256, test.rez_512, test.rez_1024, test.rez_1280, test.rez_1518, id)

	size = 128
	if verboseLogs {
		fmt.Println("->> test.thr_begin = ", test.thr_begin)
	}
	period_nano = int64(size*8*1000000000) / (int64(test.thr_begin * 1000 * 1000))
	packet_count_start = int64(period_test)

	if verboseLogs {
		fmt.Println("->> period_nano  = ", period_nano)
	}
	if verboseLogs {
		fmt.Println("->> packet_count = ", packet_count_start)
	}

	b = packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, size, number, test_type, test.test_type)
	packet_count_OK = int64(0)

	for k := 0; k < test.count_probs; k++ {
		packet_count := packet_count_start
		for j := 0; j < count_probs_one_test; j++ {
			test_c.numberCounter = 0
			cnt := packet_count
			quit := make(chan int)
			go func() {
				time.Sleep(100 * time.Millisecond)
				for {
					cnt--
					c.WriteTo(b, addr)
					if cnt < 0 {
						time.Sleep(100 * time.Millisecond)
						quit <- 1
						break
					}
				}
			}()
			(test_c).receivePackets(c, ifi.MTU, ipdst_1sfpsla_str, quit, test_type)
			rez_count := test_c.numberCounter
			if verboseLogs {
				fmt.Println("rez_count= ", rez_count)
			}
			if test_c.numberCounter >= uint64(packet_count) {
				continue
				packet_count += int64(packet_count_step_packets)
			} else {
				break
			}
		}
		packet_count_OK += int64(test_c.numberCounter)
	}
	test.rez_128 = float32(packet_count_OK) / float32(test.count_probs)
	if verboseLogs {
		fmt.Println("->> rez_128 = ", test.rez_128)
	}

	yest, _ = db.Query("SELECT EXISTS(SELECT id FROM test_bert WHERE id=?)", id)
	yest.Next()
	t_y = 0
	yest.Scan(&t_y)
	if t_y != 1 {
		return
	}
	db.Exec("UPDATE test_bert SET rez_64=?,rez_128=?,rez_256=?,rez_512=?,rez_1024=?,rez_1280=?, rez_1518=? WHERE id=?", test.rez_64, test.rez_128, test.rez_256, test.rez_512, test.rez_1024, test.rez_1280, test.rez_1518, id)

	size = 256
	if verboseLogs {
		fmt.Println("->> test.thr_begin = ", test.thr_begin)
	}
	period_nano = int64(size*8*1000000000) / (int64(test.thr_begin * 1000 * 1000))
	packet_count_start = int64(period_test)

	if verboseLogs {
		fmt.Println("->> period_nano  = ", period_nano)
	}
	if verboseLogs {
		fmt.Println("->> packet_count = ", packet_count_start)
	}

	b = packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, size, number, test_type, test.test_type)
	packet_count_OK = int64(0)

	for k := 0; k < test.count_probs; k++ {
		packet_count := packet_count_start
		for j := 0; j < count_probs_one_test; j++ {
			test_c.numberCounter = 0
			cnt := packet_count
			quit := make(chan int)
			go func() {
				time.Sleep(100 * time.Millisecond)
				for {
					cnt--
					c.WriteTo(b, addr)
					if cnt < 0 {
						time.Sleep(100 * time.Millisecond)
						quit <- 1
						break
					}
				}
			}()
			(test_c).receivePackets(c, ifi.MTU, ipdst_1sfpsla_str, quit, test_type)
			rez_count := test_c.numberCounter
			if verboseLogs {
				fmt.Println("rez_count= ", rez_count)
			}
			if test_c.numberCounter >= uint64(packet_count) {
				continue
				packet_count += int64(packet_count_step_packets)
			} else {
				break
			}
		}
		packet_count_OK += int64(test_c.numberCounter)
	}
	test.rez_256 = float32(packet_count_OK) / float32(test.count_probs)
	if verboseLogs {
		fmt.Println("->> rez_256 = ", test.rez_256)
	}

	yest, _ = db.Query("SELECT EXISTS(SELECT id FROM test_bert WHERE id=?)", id)
	yest.Next()
	t_y = 0
	yest.Scan(&t_y)
	if t_y != 1 {
		return
	}
	db.Exec("UPDATE test_bert SET rez_64=?,rez_128=?,rez_256=?,rez_512=?,rez_1024=?,rez_1280=?, rez_1518=? WHERE id=?", test.rez_64, test.rez_128, test.rez_256, test.rez_512, test.rez_1024, test.rez_1280, test.rez_1518, id)

	size = 512
	if verboseLogs {
		fmt.Println("->> test.thr_begin = ", test.thr_begin)
	}
	period_nano = int64(size*8*1000000000) / (int64(test.thr_begin * 1000 * 1000))
	packet_count_start = int64(period_test)

	if verboseLogs {
		fmt.Println("->> period_nano  = ", period_nano)
	}
	if verboseLogs {
		fmt.Println("->> packet_count = ", packet_count_start)
	}

	b = packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, size, number, test_type, test.test_type)
	packet_count_OK = int64(0)

	for k := 0; k < test.count_probs; k++ {
		packet_count := packet_count_start
		for j := 0; j < count_probs_one_test; j++ {
			test_c.numberCounter = 0
			cnt := packet_count
			quit := make(chan int)
			go func() {
				time.Sleep(100 * time.Millisecond)
				for {
					cnt--
					c.WriteTo(b, addr)
					if cnt < 0 {
						time.Sleep(100 * time.Millisecond)
						quit <- 1
						break
					}
				}
			}()
			(test_c).receivePackets(c, ifi.MTU, ipdst_1sfpsla_str, quit, test_type)
			rez_count := test_c.numberCounter
			if verboseLogs {
				fmt.Println("rez_count= ", rez_count)
			}
			if test_c.numberCounter >= uint64(packet_count) {
				continue
				packet_count += int64(packet_count_step_packets)
			} else {
				break
			}
		}
		packet_count_OK += int64(test_c.numberCounter)
	}
	test.rez_512 = float32(packet_count_OK) / float32(test.count_probs)
	if verboseLogs {
		fmt.Println("->> rez_512 = ", test.rez_512)
	}

	yest, _ = db.Query("SELECT EXISTS(SELECT id FROM test_bert WHERE id=?)", id)
	yest.Next()
	t_y = 0
	yest.Scan(&t_y)
	if t_y != 1 {
		return
	}
	db.Exec("UPDATE test_bert SET rez_64=?,rez_128=?,rez_256=?,rez_512=?,rez_1024=?,rez_1280=?, rez_1518=? WHERE id=?", test.rez_64, test.rez_128, test.rez_256, test.rez_512, test.rez_1024, test.rez_1280, test.rez_1518, id)

	size = 1024
	if verboseLogs {
		fmt.Println("->> test.thr_begin = ", test.thr_begin)
	}
	period_nano = int64(size*8*1000000000) / (int64(test.thr_begin * 1000 * 1000))
	packet_count_start = int64(period_test)

	if verboseLogs {
		fmt.Println("->> period_nano  = ", period_nano)
	}
	if verboseLogs {
		fmt.Println("->> packet_count = ", packet_count_start)
	}

	b = packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, size, number, test_type, test.test_type)
	packet_count_OK = int64(0)

	for k := 0; k < test.count_probs; k++ {
		packet_count := packet_count_start
		for j := 0; j < count_probs_one_test; j++ {
			test_c.numberCounter = 0
			cnt := packet_count
			quit := make(chan int)
			go func() {
				time.Sleep(100 * time.Millisecond)
				for {
					cnt--
					c.WriteTo(b, addr)
					if cnt < 0 {
						time.Sleep(100 * time.Millisecond)
						quit <- 1
						break
					}
				}
			}()
			(test_c).receivePackets(c, ifi.MTU, ipdst_1sfpsla_str, quit, test_type)
			rez_count := test_c.numberCounter
			if verboseLogs {
				fmt.Println("rez_count= ", rez_count)
			}
			if test_c.numberCounter >= uint64(packet_count) {
				continue
				packet_count += int64(packet_count_step_packets)
			} else {
				break
			}
		}
		packet_count_OK += int64(test_c.numberCounter)
	}
	test.rez_1024 = float32(packet_count_OK) / float32(test.count_probs)
	if verboseLogs {
		fmt.Println("->> rez_1024 = ", test.rez_1024)
	}

	yest, _ = db.Query("SELECT EXISTS(SELECT id FROM test_bert WHERE id=?)", id)
	yest.Next()
	t_y = 0
	yest.Scan(&t_y)
	if t_y != 1 {
		return
	}
	db.Exec("UPDATE test_bert SET rez_64=?,rez_128=?,rez_256=?,rez_512=?,rez_1024=?,rez_1280=?, rez_1518=? WHERE id=?", test.rez_64, test.rez_128, test.rez_256, test.rez_512, test.rez_1024, test.rez_1280, test.rez_1518, id)

	size = 1280
	if verboseLogs {
		fmt.Println("->> test.thr_begin = ", test.thr_begin)
	}
	period_nano = int64(size*8*1000000000) / (int64(test.thr_begin * 1000 * 1000))
	packet_count_start = int64(period_test)

	if verboseLogs {
		fmt.Println("->> period_nano  = ", period_nano)
	}
	if verboseLogs {
		fmt.Println("->> packet_count = ", packet_count_start)
	}

	b = packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, size, number, test_type, test.test_type)
	packet_count_OK = int64(0)

	for k := 0; k < test.count_probs; k++ {
		packet_count := packet_count_start
		for j := 0; j < count_probs_one_test; j++ {
			test_c.numberCounter = 0
			cnt := packet_count
			quit := make(chan int)
			go func() {
				time.Sleep(100 * time.Millisecond)
				for {
					cnt--
					c.WriteTo(b, addr)
					if cnt < 0 {
						time.Sleep(100 * time.Millisecond)
						quit <- 1
						break
					}
				}
			}()
			(test_c).receivePackets(c, ifi.MTU, ipdst_1sfpsla_str, quit, test_type)
			rez_count := test_c.numberCounter
			if verboseLogs {
				fmt.Println("rez_count= ", rez_count)
			}
			if test_c.numberCounter >= uint64(packet_count) {
				continue
				packet_count += int64(packet_count_step_packets)
			} else {
				break
			}
		}
		packet_count_OK += int64(test_c.numberCounter)
	}
	test.rez_1280 = float32(packet_count_OK) / float32(test.count_probs)
	if verboseLogs {
		fmt.Println("->> rez_1280 = ", test.rez_1280)
	}

	yest, _ = db.Query("SELECT EXISTS(SELECT id FROM test_bert WHERE id=?)", id)
	yest.Next()
	t_y = 0
	yest.Scan(&t_y)
	if t_y != 1 {
		return
	}
	db.Exec("UPDATE test_bert SET rez_64=?,rez_128=?,rez_256=?,rez_512=?,rez_1024=?,rez_1280=?, rez_1518=? WHERE id=?", test.rez_64, test.rez_128, test.rez_256, test.rez_512, test.rez_1024, test.rez_1280, test.rez_1518, id)

	size = 1500
	if verboseLogs {
		fmt.Println("->> test.thr_begin = ", test.thr_begin)
	}
	period_nano = int64(size*8*1000000000) / (int64(test.thr_begin * 1000 * 1000))
	packet_count_start = int64(period_test)

	if verboseLogs {
		fmt.Println("->> period_nano  = ", period_nano)
	}
	if verboseLogs {
		fmt.Println("->> packet_count = ", packet_count_start)
	}

	b = packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, size, number, test_type, test.test_type)
	packet_count_OK = int64(0)

	for k := 0; k < test.count_probs; k++ {
		packet_count := packet_count_start
		for j := 0; j < count_probs_one_test; j++ {
			test_c.numberCounter = 0
			cnt := packet_count
			quit := make(chan int)
			go func() {
				time.Sleep(100 * time.Millisecond)
				for {
					cnt--
					c.WriteTo(b, addr)
					if cnt < 0 {
						time.Sleep(100 * time.Millisecond)
						quit <- 1
						break
					}
				}
			}()
			(test_c).receivePackets(c, ifi.MTU, ipdst_1sfpsla_str, quit, test_type)
			rez_count := test_c.numberCounter
			if verboseLogs {
				fmt.Println("rez_count= ", rez_count)
			}
			if test_c.numberCounter >= uint64(packet_count) {
				continue
				packet_count += int64(packet_count_step_packets)
			} else {
				break
			}
		}
		packet_count_OK += int64(test_c.numberCounter)
	}
	test.rez_1518 = float32(packet_count_OK) / float32(test.count_probs)
	if verboseLogs {
		fmt.Println("->> rez_1518 = ", test.rez_1518)
	}
	/*
		size = 128
		period_nano = int64(size * 8 * 1000000000 / (test.thr_begin * 1000 * 1000))
		packet_count := (int64(period_test * 1000000000)) / period_nano

		b = packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, 128, number, test_type, test.test_type)
		count_rez = 0
		per = 0
		//genSocket(ifi.Index, b, period_test, test.thr_begin)

		//count_rez, per = test_c.testMax(b, c, addr, ifi.MTU, ipdst_1sfpsla_str, counter, test_type)

		for k := 0; k < test.count_probs; k++ {
			test_c.numberCounter = 0
			count_rez_uni, _, per_uni = test_c.testThrGen(net_interface_name, b, b, addr, ifi.MTU, ipdst_1sfpsla_str, packet_count, period_nano, test.thr_begin, test_type)
			count_rez = count_rez + count_rez_uni
			per = per + per_uni
		}

		test.rez_128 = (float32)(128.0 * 8.0 * (float32)(count_rez) * 1000000 / (float32)(per))
		if verboseLogs {
			fmt.Println("->> rez_128 = ", count_rez, " period = ", per, " !!!  rez=", test.rez_128)
		}

		size = 256
		period_nano = int64(size * 8 * 1000000000 / (test.thr_begin * 1000 * 1000))
		packet_count = (int64(period_test * 1000000000)) / period_nano
		count_rez = 0
		per = 0
		b = packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, 256, number, test_type, test.test_type)
		for k := 0; k < test.count_probs; k++ {
			test_c.numberCounter = 0
			count_rez_uni, _, per_uni = test_c.testThrGen(net_interface_name, b, b, addr, ifi.MTU, ipdst_1sfpsla_str, packet_count, period_nano, test.thr_begin, test_type)
			count_rez = count_rez + count_rez_uni
			per = per + per_uni
		}

		test.rez_256 = (float32)(256.0 * 8.0 * (float32)(count_rez) * 1000000 / (float32)(per))
		if verboseLogs {
			fmt.Println("->> rez_256 = ", count_rez, " period = ", per, " !!!  rez=", test.rez_256)
		}

		size = 512
		period_nano = int64(size * 8 * 1000000000 / (test.thr_begin * 1000 * 1000))
		packet_count = (int64(period_test * 1000000000)) / period_nano
		count_rez = 0
		per = 0
		b = packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, 512, number, test_type, test.test_type)
		for k := 0; k < test.count_probs; k++ {
			test_c.numberCounter = 0
			count_rez_uni, _, per_uni = test_c.testThrGen(net_interface_name, b, b, addr, ifi.MTU, ipdst_1sfpsla_str, packet_count, period_nano, test.thr_begin, test_type)
			count_rez = count_rez + count_rez_uni
			per = per + per_uni
		}

		test.rez_512 = (float32)(512.0 * 8.0 * (float32)(count_rez) * 1000000 / (float32)(per))
		if verboseLogs {
			fmt.Println("->> rez_512 = ", count_rez, " period = ", per, " !!!  rez=", test.rez_512)
		}

		size = 1024
		period_nano = int64(size * 8 * 1000000000 / (test.thr_begin * 1000 * 1000))
		packet_count = (int64(period_test * 1000000000)) / period_nano
		count_rez = 0
		per = 0
		b = packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, size, number, test_type, test.test_type)
		for k := 0; k < test.count_probs; k++ {
			test_c.numberCounter = 0
			count_rez_uni, _, per_uni = test_c.testThrGen(net_interface_name, b, b, addr, ifi.MTU, ipdst_1sfpsla_str, packet_count, period_nano, test.thr_begin, test_type)
			count_rez = count_rez + count_rez_uni
			per = per + per_uni
		}

		test.rez_1024 = (float32)(float32(size) * 8.0 * (float32)(count_rez) * 1000000 / (float32)(per))
		if verboseLogs {
			fmt.Println("->> rez_1024 = ", count_rez, " period = ", per, " !!!  rez=", test.rez_1024)
		}

		size = 1280
		period_nano = int64(size * 8 * 1000000000 / (test.thr_begin * 1000 * 1000))
		packet_count = (int64(period_test * 1000000000)) / period_nano
		count_rez = 0
		per = 0
		b = packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, size, number, test_type, test.test_type)
		for k := 0; k < test.count_probs; k++ {
			test_c.numberCounter = 0
			count_rez_uni, _, per_uni = test_c.testThrGen(net_interface_name, b, b, addr, ifi.MTU, ipdst_1sfpsla_str, packet_count, period_nano, test.thr_begin, test_type)
			count_rez = count_rez + count_rez_uni
			per = per + per_uni
		}

		test.rez_1280 = (float32)(float32(size) * 8.0 * (float32)(count_rez) * 1000000 / (float32)(per))
		if verboseLogs {
			fmt.Println("->> rez_1280 = ", count_rez, " period = ", per, " !!!  rez=", test.rez_1280)
		}

		size = 1500
		period_nano = int64(size * 8 * 1000000000 / (test.thr_begin * 1000 * 1000))
		packet_count = (int64(period_test * 1000000000)) / period_nano
		count_rez = 0
		per = 0
		b = packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, size, number, test_type, test.test_type)
		for k := 0; k < test.count_probs; k++ {
			test_c.numberCounter = 0
			count_rez_uni, _, per_uni = test_c.testThrGen(net_interface_name, b, b, addr, ifi.MTU, ipdst_1sfpsla_str, packet_count, period_nano, test.thr_begin, test_type)
			count_rez = count_rez + count_rez_uni
			per = per + per_uni
		}

		test.rez_1518 = (float32)(float32(size) * 8.0 * (float32)(count_rez) * 1000000 / (float32)(per))
		if verboseLogs {
			fmt.Println("->> rez_1518 = ", count_rez, " period = ", per, " !!!  rez=", test.rez_1518)
		}
		/*
			b = packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, 256, number, test_type)
			count_rez, per = test_c.testMax(b, c, addr, ifi.MTU, ipdst_1sfpsla_str, counter, test_type)
			test.rez_256 = (float32)(256.0 * 8.0 * (float32)(count_rez) * 1000 / (float32)(per))

			b = packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, 512, number, test_type)
			count_rez, per = test_c.testMax(b, c, addr, ifi.MTU, ipdst_1sfpsla_str, counter, test_type)
			test.rez_512 = (float32)(512.0 * 8.0 * (float32)(count_rez) * 1000 / (float32)(per))

			b = packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, 1024, number, test_type)
			count_rez, per = test_c.testMax(b, c, addr, ifi.MTU, ipdst_1sfpsla_str, counter, test_type)
			test.rez_1024 = (float32)(1024.0 * 8.0 * (float32)(count_rez) * 1000 / (float32)(per))

			b = packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, 1280, number, test_type)
			count_rez, per = test_c.testMax(b, c, addr, ifi.MTU, ipdst_1sfpsla_str, counter, test_type)
			test.rez_1280 = (float32)(1280.0 * 8.0 * (float32)(count_rez) * 1000 / (float32)(per))

			b = packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, 1500, number, test_type)
			count_rez, per = test_c.testMax(b, c, addr, ifi.MTU, ipdst_1sfpsla_str, counter, test_type)
			test.rez_1518 = (float32)(1518.0 * 8.0 * (float32)(count_rez) * 1000 / (float32)(per))
	*/
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
	db.Exec("UPDATE test_bert SET rez_64=?,rez_128=?,rez_256=?,rez_512=?,rez_1024=?,rez_1280=?, rez_1518=?,datetime_end=?, status=? WHERE id=?", test.rez_64, test.rez_128, test.rez_256, test.rez_512, test.rez_1024, test.rez_1280, test.rez_1518, time.Now().Format("2006-01-02 15:04:05"), test.status, id)

	db.Close()
}
