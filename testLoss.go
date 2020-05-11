package main

import (
	"database/sql"
	"fmt"
	"net"
	"time"

	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/raw"
)

func TestLoss(id int, net_interface_name string) { //Нагрузочное тестирование задержки
	fmt.Println("Тест задержки начался")
	db, err := sql.Open("mysql", db_user+":"+db_user_pass+"@/"+db_database)
	if err != nil {
		db.Close()
		fmt.Println(" -!! Error !!-")
		fmt.Println(err)
		fmt.Println(" ----=====----")
		return
	}
	db.Exec("UPDATE test_latency SET status=?, datatime=? WHERE id=?", time.Now(), 2, id) // Тест выполняется
	ifi, err := net.InterfaceByName(net_interface_name)
	if err != nil {
		db.Close()
		//	log.Fatalf("failed to find interface %q: %v", net_interface_name, err)
		db.Exec("UPDATE test_latency SET status=? WHERE id=?", 4, id) // Ошибка выполнения
		return
	}

	row, err := db.Query("select id, test_type, module_first, module_second, thr_begin, count_packs, count_tests, status from test_latency where id=?", id)
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
	test := new(testDelay)
	err = row.Scan(&test.id, &test.test_type, &test.module_first, &test.module_second, &test.thr_begin, &test.count_packs, &test.count_tests, &test.status)
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
	test.mac_dst = make([]byte, 6)

	row, err = db.Query("SELECT server_IP FROM global_config")
	if err != nil {
		db.Exec("UPDATE test_latency SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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

			db.Exec("UPDATE test_latency SET status=? WHERE id=?", 4, id) // Ошибка выполнения
			db.Close()
			row.Close()
			fmt.Println(" -!! Error !!-")
			fmt.Println(err)
			fmt.Println(" ----=====----")
			return
		}
	}
	var id_sfp1, id_sfp2 int
	row, err = db.Query("SELECT module_first, module_second FROM test_latency WHERE id=?", id)
	if err != nil {
		db.Exec("UPDATE test_latency SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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

			db.Exec("UPDATE test_latency SET status=? WHERE id=?", 4, id) // Ошибка выполнения
			db.Close()
			row.Close()
			fmt.Println(" -!! Error !!-")
			fmt.Println(err)
			fmt.Println(" ----=====----")
			return
		}
		row_ip, err := db.Query("SELECT address_ip FROM modules_sfp_sla WHERE id=?", id_sfp1)
		if err != nil {
			db.Exec("UPDATE test_latency SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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

				db.Exec("UPDATE test_latency SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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
			db.Exec("UPDATE test_latency SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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

				db.Exec("UPDATE test_latency SET status=? WHERE id=?", 4, id) // Ошибка выполнения
				db.Close()
				row.Close()
				fmt.Println(" -!! Error !!-")
				fmt.Println(err)
				fmt.Println(" ----=====----")
				return
			}
		}
		test.mac_dst[5] = byte(test_mac & 0xFF)
		test.mac_dst[4] = byte((test_mac >> 8) & 0xFF)
		test.mac_dst[3] = byte((test_mac >> 16) & 0xFF)
		test.mac_dst[2] = byte((test_mac >> 24) & 0xFF)
		test.mac_dst[1] = byte((test_mac >> 32) & 0xFF)
		test.mac_dst[0] = byte((test_mac >> 40) & 0xFF)

		row_ip, err = db.Query("SELECT address_ip FROM modules_sfp_sla WHERE id=?", id_sfp2)
		if err != nil {
			db.Exec("UPDATE test_latency SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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

				db.Exec("UPDATE test_latency SET status=? WHERE id=?", 4, id) // Ошибка выполнения
				db.Close()
				row.Close()
				fmt.Println(" -!! Error !!-")
				fmt.Println(err)
				fmt.Println(" ----=====----")
				return
			}
		}
	}
	db.Exec("UPDATE test_latency SET status=?, datetime_start=? WHERE id=?", 2, time.Now(), id) // Тест выполняется
	db.Close()

	fmt.Println(ipsrcstr)
	fmt.Println(ipdst_1sfpsla_str)
	fmt.Println(ipdst_2sfpsla_str)

	//counter := test.count
	//var numberTX uint32
	//	numberTX = 0

	test.ipsrc = net.ParseIP(ipsrcstr)
	test.ipdst1 = net.ParseIP(ipdst_1sfpsla_str)
	test.ipdst2 = net.ParseIP(ipdst_2sfpsla_str)

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
	/*
		addr := &raw.Addr{
			HardwareAddr: ethernet.Broadcast,
		}
		var number uint32
	*/

	addr := &raw.Addr{
		HardwareAddr: ethernet.Broadcast,
	}
	connectTestSFP, err := raw.ListenPacket(ifi, etherType, nil)

	rez := findSFP(connectTestSFP, addr, ipsrcstr, ipdst_1sfpsla_str, ipdst_2sfpsla_str, test.mac_src, test.mac_dst, test.id_test_type, test.test_type)
	if rez == 0 {
		fmt.Println("Error test SFP connect")
		return
	}
	if rez == 2 {

		tmp := ipdst_1sfpsla_str
		ipdst_1sfpsla_str = ipdst_2sfpsla_str
		ipdst_2sfpsla_str = tmp

		test.ipsrc = net.ParseIP(ipsrcstr)
		test.ipdst1 = net.ParseIP(ipdst_1sfpsla_str)
		test.ipdst2 = net.ParseIP(ipdst_2sfpsla_str)

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
		test.mac_dst[5] = byte(test_mac & 0xFF)
		test.mac_dst[4] = byte((test_mac >> 8) & 0xFF)
		test.mac_dst[3] = byte((test_mac >> 16) & 0xFF)
		test.mac_dst[2] = byte((test_mac >> 24) & 0xFF)
		test.mac_dst[1] = byte((test_mac >> 32) & 0xFF)
		test.mac_dst[0] = byte((test_mac >> 40) & 0xFF)

		db.Close()

	}
	fmt.Println(" Rez find : ", rez)

	test.id_test_type = 0x2000 + (uint16(id) & 0x1FFF)
	test.net_interface_name = net_interface_name
	test.mac_src = ifi.HardwareAddr

	counter := make(chan int64, 7)
	quit := make(chan int64, 7)

	size_p := 64

	b := packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size_p, 0, test.id_test_type, test.test_type)
	go genSocket(ifi.Index, b, test.count_packs, test.thr_begin, counter)
	test.getMonDelay(quit, size_p)
	time.Sleep(time.Second * 2)
	<-quit

	size_p = 128

	b = packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size_p, 0, test.id_test_type, test.test_type)
	go genSocket(ifi.Index, b, test.count_packs, test.thr_begin, counter)
	test.getMonDelay(quit, size_p)
	time.Sleep(time.Second * 2)
	<-quit

	size_p = 256

	b = packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size_p, 0, test.id_test_type, test.test_type)
	go genSocket(ifi.Index, b, test.count_packs, test.thr_begin, counter)
	test.getMonDelay(quit, size_p)
	time.Sleep(time.Second * 2)
	<-quit

	size_p = 512

	b = packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size_p, 0, test.id_test_type, test.test_type)
	go genSocket(ifi.Index, b, test.count_packs, test.thr_begin, counter)
	test.getMonDelay(quit, size_p)
	time.Sleep(time.Second * 2)
	<-quit

	size_p = 1024

	b = packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size_p, 0, test.id_test_type, test.test_type)
	go genSocket(ifi.Index, b, test.count_packs, test.thr_begin, counter)
	test.getMonDelay(quit, size_p)
	time.Sleep(time.Second * 2)
	<-quit

	size_p = 1280

	b = packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size_p, 0, test.id_test_type, test.test_type)
	go genSocket(ifi.Index, b, test.count_packs, test.thr_begin, counter)
	test.getMonDelay(quit, size_p)
	time.Sleep(time.Second * 2)
	<-quit

	size_p = 1518

	b = packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, 1500, 0, test.id_test_type, test.test_type)
	go genSocket(ifi.Index, b, test.count_packs, test.thr_begin, counter)
	test.getMonDelay(quit, 1500)
	time.Sleep(time.Second * 2)
	<-quit

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
	db.Exec("UPDATE test_latency SET rez_64=?,rez_64_max=?,rez_64_min=?,rez_128=?,rez_128_max=?,rez_128_min=?,rez_256=?,rez_256_max=?,rez_256_min=?,rez_512=?,rez_512_max=?,rez_512_min=?,rez_1024=?,rez_1024_max=?,rez_1024_min=?,rez_1280=?,rez_1280_max=?,rez_1280_min=?, rez_1518=?,rez_1518_max=?,rez_1518_min=?,datetime_end=?, status=? WHERE id=?", test.rez_64, test.rez_64_max, test.rez_64_min, test.rez_128, test.rez_128_max, test.rez_128_min, test.rez_256, test.rez_256_max, test.rez_256_min, test.rez_512, test.rez_512_max, test.rez_512_min, test.rez_1024, test.rez_1024_max, test.rez_1024_min, test.rez_1280, test.rez_1280_max, test.rez_1280_min, test.rez_1518, test.rez_1518_max, test.rez_1518_min, time.Now(), test.status, id)

	db.Close()
}
