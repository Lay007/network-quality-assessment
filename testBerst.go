package main

import (
	"database/sql"
	"fmt"
	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/raw"
	"net"
	"time"

	"golang.org/x/net/bpf"
)

func TestBerst(id int, net_interface_name string) { //Нагрузочное тестирование пропускной способности
	fmt.Println("Тест максимальной пропускной способности начался")
	db, err := sql.Open("mysql", db_user+":"+db_user_pass+"@/"+db_database)
	if err != nil {
		db.Close()
		fmt.Println(" -!! Error !!-")
		fmt.Println(err)
		fmt.Println(" ----=====----")
		return
	}
	db.Exec("UPDATE test_bert SET status=?, datatime=? WHERE id=?", time.Now().Format("2006-01-02 15:04:05"), 2, id) // Тест выполняется
	ifi, err := net.InterfaceByName(net_interface_name)
	if err != nil {
		db.Close()
		//	log.Fatalf("failed to find interface %q: %v", net_interface_name, err)
		db.Exec("UPDATE test_bert SET status=? WHERE id=?", 4, id) // Ошибка выполнения
		return
	}

	row, err := db.Query("select id, test_type, module_first, module_second, thr_begin, count_prob_packs, count_probs, status from test_bert where id=?", id)
	if err != nil {
		db.Close()
		row.Close()
		fmt.Println(" -!! Error !!-")
		fmt.Println(err)
		fmt.Println(" ----=====----")
		return
	}
	fmt.Println(row)
	defer row.Close()
	row.Next()
	test := new(testBert)
	err = row.Scan(&test.id, &test.test_type, &test.module_first, &test.module_second, &test.thr_begin, &test.count_prob_packs, &test.count_probs, &test.status)
	if err != nil {
		db.Close()
		fmt.Println(" -!! Error !!-")
		fmt.Println(err)
		fmt.Println(" ----=====----")
		return
	}

	var ipsrcstr string
	var ipdst_1sfpsla_str string
	var ipdst_2sfpsla_str string
	mac_dst := make([]byte, 6)

	row, err = db.Query("SELECT server_IP FROM global_config")
	if err != nil {
		db.Exec("UPDATE test_bert SET status=? WHERE id=?", 4, id) // Ошибка выполнения
		db.Close()
		row.Close()
		fmt.Println(" -!! Error !!-")
		fmt.Println(err)
		fmt.Println(" ----=====----")
		return
	}
	for row.Next() {
		err = row.Scan(&ipsrcstr)
		if err != nil {

			db.Exec("UPDATE test_bert SET status=? WHERE id=?", 4, id) // Ошибка выполнения
			db.Close()
			row.Close()
			fmt.Println(" -!! Error !!-")
			fmt.Println(err)
			fmt.Println(" ----=====----")
			return
		}
	}
	var id_sfp1, id_sfp2 int
	row, err = db.Query("SELECT module_first, module_second FROM test_bert WHERE id=?", id)
	if err != nil {
		db.Exec("UPDATE test_bert SET status=? WHERE id=?", 4, id) // Ошибка выполнения
		db.Close()
		row.Close()
		fmt.Println(" -!! Error !!-")
		fmt.Println(err)
		fmt.Println(" ----=====----")
		return
	}
	for row.Next() {
		err = row.Scan(&id_sfp1, &id_sfp2)
		if err != nil {

			db.Exec("UPDATE test_bert SET status=? WHERE id=?", 4, id) // Ошибка выполнения
			db.Close()
			row.Close()
			fmt.Println(" -!! Error !!-")
			fmt.Println(err)
			fmt.Println(" ----=====----")
			return
		}
		row_ip, err := db.Query("SELECT address_ip FROM modules_sfp_sla WHERE id=?", id_sfp1)
		if err != nil {
			db.Exec("UPDATE test_bert SET status=? WHERE id=?", 4, id) // Ошибка выполнения
			db.Close()
			row.Close()
			row_ip.Close()
			fmt.Println(" -!! Error !!-")
			fmt.Println(err)
			fmt.Println(" ----=====----")
			return
		}
		defer row_ip.Close()
		for row_ip.Next() {
			err = row_ip.Scan(&ipdst_1sfpsla_str)
			if err != nil {

				db.Exec("UPDATE test_bert SET status=? WHERE id=?", 4, id) // Ошибка выполнения
				db.Close()
				row.Close()
				fmt.Println(" -!! Error !!-")
				fmt.Println(err)
				fmt.Println(" ----=====----")
				return
			}
		}

		row_mac, err := db.Query("SELECT mac FROM modules_sfp_sla WHERE id=?", id_sfp1)
		if err != nil {
			db.Exec("UPDATE test_bert SET status=? WHERE id=?", 4, id) // Ошибка выполнения
			db.Close()
			row.Close()
			row_mac.Close()
			fmt.Println(" -!! Error !!-")
			fmt.Println(err)
			fmt.Println(" ----=====----")
			return
		}
		defer row_mac.Close()
		//var mac_dst_str string
		var test_mac int64
		for row_mac.Next() {
			//err = row_mac.Scan(&mac_dst_str)
			err = row_mac.Scan(&test_mac)
			if err != nil {

				db.Exec("UPDATE test_bert SET status=? WHERE id=?", 4, id) // Ошибка выполнения
				db.Close()
				row.Close()
				fmt.Println(" -!! Error !!-")
				fmt.Println(err)
				fmt.Println(" ----=====----")
				return
			}
		}
		mac_dst[5] = byte(test_mac & 0xFF)
		mac_dst[4] = byte((test_mac >> 8) & 0xFF)
		mac_dst[3] = byte((test_mac >> 16) & 0xFF)
		mac_dst[2] = byte((test_mac >> 24) & 0xFF)
		mac_dst[1] = byte((test_mac >> 32) & 0xFF)
		mac_dst[0] = byte((test_mac >> 40) & 0xFF)

		row_ip, err = db.Query("SELECT address_ip FROM modules_sfp_sla WHERE id=?", id_sfp2)
		if err != nil {
			db.Exec("UPDATE test_bert SET status=? WHERE id=?", 4, id) // Ошибка выполнения
			db.Close()
			row.Close()
			fmt.Println(" -!! Error !!-")
			fmt.Println(err)
			fmt.Println(" ----=====----")
			return
		}
		for row_ip.Next() {
			err = row_ip.Scan(&ipdst_2sfpsla_str)
			if err != nil {

				db.Exec("UPDATE test_bert SET status=? WHERE id=?", 4, id) // Ошибка выполнения
				db.Close()
				row.Close()
				fmt.Println(" -!! Error !!-")
				fmt.Println(err)
				fmt.Println(" ----=====----")
				return
			}
		}
	}

	db.Exec("UPDATE test_bert SET status=?, datetime_start=? WHERE id=?", 2, time.Now().Format("2006-01-02 15:04:05"), id) // Тест выполняется
	db.Close()

	fmt.Println(ipsrcstr)
	fmt.Println(ipdst_1sfpsla_str)
	fmt.Println(ipdst_2sfpsla_str)

	//counter := test.count
	//var numberTX uint32
	//	numberTX = 0

	ipsrc := net.ParseIP(ipsrcstr)
	ipdst1 := net.ParseIP(ipdst_1sfpsla_str)
	ipdst2 := net.ParseIP(ipdst_2sfpsla_str)

	//period_min := time.Duration(time.Duration(int(period_nano)) * time.Nanosecond)
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
		fmt.Println(" -- Test ticker= ", time.Since(start_test_ticker))
	*/

	//fmt.Println("period_min= ", period_min)
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
		// Проверка идентификатора пакета (34 бит) (xFA-от 1 ко 2, xFB – от 2 к 1, xFC – от 1 к Серверу)
		bpf.LoadAbsolute{Off: 34, Size: 1},
		bpf.JumpIf{Cond: bpf.JumpNotEqual, Val: 0xFC, SkipTrue: 3},
		// Проверка идентификатора теста
		bpf.LoadAbsolute{Off: 64, Size: 2},
		bpf.JumpIf{Cond: bpf.JumpNotEqual, Val: uint32(test_type), SkipTrue: 1},
		// Verdict is "send up to 4k of the packet to userspace."
		bpf.RetConstant{Val: 4096},
		// Verdict is "ignore packet."
		bpf.RetConstant{Val: 0},
	})

	period_test := test.count_prob_packs // период теста - 10 секунд

	size := 64
	fmt.Println("->> test.thr_begin = ", test.thr_begin)
	period_nano := int64(size*8*1000000000) / (int64(test.thr_begin * 1000 * 1000))
	packet_count := (int64(period_test * 1000000000)) / period_nano

	fmt.Println("->> period_nano  = ", period_nano)
	fmt.Println("->> packet_count = ", packet_count)

	connectTestSFP, err := raw.ListenPacket(ifi, etherType, nil)

	rez := findSFP(connectTestSFP, addr, ipsrcstr, ipdst_1sfpsla_str, ipdst_2sfpsla_str, ifi.HardwareAddr, mac_dst, test_type, test.test_type, 1024, int64(1024*8*1000/(test.thr_begin)))
	if rez == 0 {
		fmt.Println("Error test SFP connect")
		return
	}
	if rez == 2 {

		tmp := ipdst_1sfpsla_str
		ipdst_1sfpsla_str = ipdst_2sfpsla_str
		ipdst_2sfpsla_str = tmp

		ipsrc = net.ParseIP(ipsrcstr)
		ipdst1 = net.ParseIP(ipdst_1sfpsla_str)
		ipdst2 = net.ParseIP(ipdst_2sfpsla_str)

		db, err = sql.Open("mysql", db_user+":"+db_user_pass+"@/"+db_database)
		row_mac, err := db.Query("SELECT mac FROM modules_sfp_sla WHERE address_ip=?", ipdst_1sfpsla_str)
		if err != nil {
			db.Close()
			row_mac.Close()
			fmt.Println(" -!! Error !!-")
			fmt.Println(err)
			fmt.Println(" ----=====----")
			return
		}
		defer row_mac.Close()
		var test_mac int64
		for row_mac.Next() {
			err = row_mac.Scan(&test_mac)
			if err != nil {
				db.Close()
				row_mac.Close()
				fmt.Println(" -!! Error !!-")
				fmt.Println(err)
				fmt.Println(" ----=====----")
				return
			}
		}
		mac_dst[5] = byte(test_mac & 0xFF)
		mac_dst[4] = byte((test_mac >> 8) & 0xFF)
		mac_dst[3] = byte((test_mac >> 16) & 0xFF)
		mac_dst[2] = byte((test_mac >> 24) & 0xFF)
		mac_dst[1] = byte((test_mac >> 32) & 0xFF)
		mac_dst[0] = byte((test_mac >> 40) & 0xFF)

		db.Close()

	}

	if testPing(ipdst_1sfpsla_str) > 0 || testPing(ipdst_2sfpsla_str) > 0 {
		db.Exec("UPDATE test_bert SET status=? WHERE id=?", 4, id) // Ошибка выполнения
		db.Close()

		return
	}

	if check_SNMP(ipdst_1sfpsla_str) > 0 || check_SNMP(ipdst_2sfpsla_str) > 0 {
		db.Exec("UPDATE test_bert SET status=? WHERE id=?", 4, id) // Ошибка выполнения
		db.Close()

		return
	}

	db.Close()
	fmt.Println(" Rez find : ", rez)

	var count_rez, count_rez_uni int
	var per, per_uni int64
	size = 64

	b := packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, size, number, test_type, test.test_type)
	count_rez = 0
	per = 0
	for k := 0; k < test.count_probs; k++ {
		test_c.numberCounter = 0
		count_rez_uni, _, per_uni = test_c.testThrGen(net_interface_name, b, b, addr, ifi.MTU, ipdst_1sfpsla_str, packet_count, period_nano, test.thr_begin, test_type)
		count_rez = count_rez + count_rez_uni
		per = per + per_uni
	}
	test.rez_64 = (float32)(64.0 * 8.0 * (float32)(count_rez) * 1000000 / (float32)(per))
	fmt.Println("->> rez_64 = ", count_rez, " period = ", per, " !!!  rez=", test.rez_64)

	size = 128
	period_nano = int64(size * 8 * 1000000000 / (test.thr_begin * 1000 * 1000))
	packet_count = (int64(period_test * 1000000000)) / period_nano

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
	fmt.Println("->> rez_128 = ", count_rez, " period = ", per, " !!!  rez=", test.rez_128)

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
	fmt.Println("->> rez_256 = ", count_rez, " period = ", per, " !!!  rez=", test.rez_256)

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
	fmt.Println("->> rez_512 = ", count_rez, " period = ", per, " !!!  rez=", test.rez_512)

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
	fmt.Println("->> rez_1024 = ", count_rez, " period = ", per, " !!!  rez=", test.rez_1024)

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
	fmt.Println("->> rez_1280 = ", count_rez, " period = ", per, " !!!  rez=", test.rez_1280)

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
	fmt.Println("->> rez_1518 = ", count_rez, " period = ", per, " !!!  rez=", test.rez_1518)
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

	db, err = sql.Open("mysql", db_user+":"+db_user_pass+"@/"+db_database)
	if err != nil {
		db.Close()

		fmt.Println(" -!! Error !!-")
		fmt.Println(err)
		fmt.Println(" ----=====----")
		return
	}
	//db.Exec("UPDATE test_throughput SET rez_64=?,rez_128=?,rez_256=?,rez_512=?,rez_1024=?,rez_1280=?, rez_1518=?,rez_4096=?,rez_9000=?,status=? WHERE id=?", test.rez_64, test.rez_128, test.rez_256, test.rez_512, test.rez_1024, test.rez_1280, test.rez_1518, test.rez_4096, test.rez_9000, test.status, id)
	db.Exec("UPDATE test_bert SET rez_64=?,rez_128=?,rez_256=?,rez_512=?,rez_1024=?,rez_1280=?, rez_1518=?,datetime_end=?, status=? WHERE id=?", test.rez_64, test.rez_128, test.rez_256, test.rez_512, test.rez_1024, test.rez_1280, test.rez_1518, time.Now().Format("2006-01-02 15:04:05"), test.status, id)

	db.Close()
}
