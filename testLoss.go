package main

import (
	"database/sql"
	"fmt"
	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/raw"
	"net"
	"runtime"
	"sync/atomic"
	//"runtime"
	//"runtime/debug"
	"time"

	"golang.org/x/net/bpf"
)

func TestLoss(id int, net_interface_name string) { //Нагрузочное тестирование задержки
	fmt.Println("Тест потери пакетов начался")
	db, err := sql.Open("mysql", db_user+":"+db_user_pass+"@/"+db_database)
	if err != nil {
		db.Close()
		fmt.Println(" -!! Error !!-")
		fmt.Println(err)
		fmt.Println(" ----=====----")
		return
	}
	db.Exec("UPDATE test_frame_loss SET status=?, datatime=? WHERE id=?", time.Now().Format("2006-01-02 15:04:05"), 2, id) // Тест выполняется
	ifi, err := net.InterfaceByName(net_interface_name)
	if err != nil {
		db.Close()
		//	log.Fatalf("failed to find interface %q: %v", net_interface_name, err)
		db.Exec("UPDATE test_frame_loss SET status=? WHERE id=?", 4, id) // Ошибка выполнения
		return
	}

	row, err := db.Query("select id, test_type, module_first, module_second, thr_begin, step, count_frames, count_steps, status from test_frame_loss where id=?", id)
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
	test := new(testLoss)
	err = row.Scan(&test.id, &test.test_type, &test.module_first, &test.module_second, &test.thr_begin, &test.step, &test.count_frames, &test.count_steps, &test.status)
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
	test.mac_dst2 = make([]byte, 6)

	row, err = db.Query("SELECT server_IP FROM global_config")
	if err != nil {
		db.Exec("UPDATE test_frame_loss SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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

			db.Exec("UPDATE test_frame_loss SET status=? WHERE id=?", 4, id) // Ошибка выполнения
			db.Close()
			row.Close()
			fmt.Println(" -!! Error !!-")
			fmt.Println(err)
			fmt.Println(" ----=====----")
			return
		}
	}
	var id_sfp1, id_sfp2 int
	row, err = db.Query("SELECT module_first, module_second FROM test_frame_loss WHERE id=?", id)
	if err != nil {
		db.Exec("UPDATE test_frame_loss SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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

			db.Exec("UPDATE test_frame_loss SET status=? WHERE id=?", 4, id) // Ошибка выполнения
			db.Close()
			row.Close()
			fmt.Println(" -!! Error !!-")
			fmt.Println(err)
			fmt.Println(" ----=====----")
			return
		}
		row_ip, err := db.Query("SELECT address_ip FROM modules_sfp_sla WHERE id=?", id_sfp1)
		if err != nil {
			db.Exec("UPDATE test_frame_loss SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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

				db.Exec("UPDATE test_frame_loss SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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
			db.Exec("UPDATE test_frame_loss SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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

				db.Exec("UPDATE test_frame_loss SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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

		row_mac, err = db.Query("SELECT mac FROM modules_sfp_sla WHERE id=?", id_sfp2)
		if err != nil {
			db.Exec("UPDATE test_frame_loss SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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
		test_mac = 0
		for row_mac.Next() {
			//err = row_mac.Scan(&mac_dst_str)
			err = row_mac.Scan(&test_mac)
			if err != nil {

				db.Exec("UPDATE test_frame_loss SET status=? WHERE id=?", 4, id) // Ошибка выполнения
				db.Close()
				row.Close()
				fmt.Println(" -!! Error !!-")
				fmt.Println(err)
				fmt.Println(" ----=====----")
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
			db.Exec("UPDATE test_frame_loss SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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

				db.Exec("UPDATE test_frame_loss SET status=? WHERE id=?", 4, id) // Ошибка выполнения
				db.Close()
				row.Close()
				fmt.Println(" -!! Error !!-")
				fmt.Println(err)
				fmt.Println(" ----=====----")
				return
			}
		}
	}
	db.Exec("UPDATE test_frame_loss SET status=?, datetime_start=? WHERE id=?", 2, time.Now().Format("2006-01-02 15:04:05"), id) // Тест выполняется

	if testPing(ipdst_1sfpsla_str) > 0 || testPing(ipdst_2sfpsla_str) > 0 {
		db.Exec("UPDATE test_frame_loss SET status=? WHERE id=?", 4, id) // Ошибка выполнения
		db.Close()

		return
	}

	if check_SNMP(ipdst_1sfpsla_str) > 0 || check_SNMP(ipdst_2sfpsla_str) > 0 {
		db.Exec("UPDATE test_frame_loss SET status=? WHERE id=?", 4, id) // Ошибка выполнения
		db.Close()

		return
	}

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

	rez := findSFP(connectTestSFP, addr, ipsrcstr, ipdst_1sfpsla_str, ipdst_2sfpsla_str, test.mac_src, test.mac_dst, test.mac_dst2, test.id_test_type, test.test_type, 1024, int64(1024*8*1000/test.thr_begin))
	fmt.Printf("\n Rez find : %X \n", rez)
	if (rez & 0xFFF) == 0x999 {
		db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 3, id, " Максимальная пропускная способность - 1 Гбит/с")
		if test.thr_begin > 1500 {
			test.thr_begin = 1500
			db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 3, id, " Изменяем максимальную пропускную способность на 1.5 Гбит/с")
		}
	}
	if (rez & 0xFFF) == 0x100 {
		db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 3, id, " Максимальная пропускная способность - 100 Мбит/с")
		
		if test.thr_begin > 150 {
			test.thr_begin = 150
			db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 3, id, " Изменяем максимальную пропускную способность на 150 Мбит/с")
		}
	}
	if (rez & 0xFFF) == 0x10 {
		db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 3, id, " Максимальная пропускная способность - 10 Мбит/с")
	
		if test.thr_begin > 15 {
			test.thr_begin = 15
			db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 3, id, " Изменяем максимальную пропускную способность на 15 Мбит/с")
		}
	}
	if (rez & 0xF000) == 0x0 {
		db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 3, id, " Последовательное расположение модулей. Порядок модулей правильный")
	}
	if (rez & 0xF000) == 0x1000 {
		db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 3, id, " Последовательное расположение модулей. Порядок модулей неправильный. Изменяем при тестироваинии")
	}
	if (rez & 0xF000) == 0x2000 {
		db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 3, id, " Соединенеие звездой. Нагрузка одинаковая")
	}
	if (rez & 0xF000) == 0x3000 {
		db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 3, id, "  Соединенеие звездой. Нагрузка неравномерная. Расположение правильное")
	}
	if (rez & 0xF000) == 0x4000 {
		db.Exec("INSERT INTO message (date,test_type, test_id, message) VALUES(NOW(),?, ?, ?)", 3, id, "  Соединенеие звездой. Нагрузка неравномерная. Расположение неправильное. Изменяем при тестироваинии")
	}
	connectTestSFP.Close()
	if rez == 0 {
		fmt.Println("Error test SFP connect")
		return
	}
	if ((rez & 0xF000) == 0x1000) || ((rez & 0xF000) == 0x4000) {

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

	test.id_test_type = 0x8000 + (uint16(id) & 0x1FFF)
	testTypeTemp := 0x2000 + (uint16(id) & 0x1FFF)
	test.net_interface_name = net_interface_name
	test.mac_src = ifi.HardwareAddr
	for step := 0; step < test.count_steps; step++ {

		testRez := new(testLossRez)
		testRez.id_test = id
		testRez.step_number = step

		thr_step := test.thr_begin - step*(int(float32(test.thr_begin)*float32(test.step)/100.0))

		counter := make(chan uint64, 7)
		counterRes := make(chan uint64, 7)
		quit := make(chan int64, 7)

		size_p := 64
		test.numberRx = 0
		b := packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size_p, 0, test.id_test_type, test.test_type)
		bTemp := packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size_p, 0, testTypeTemp, test.test_type)
		go test.receivePacketsLoss(ifi.MTU, quit)
		fmt.Println("*")
		time.Sleep(time.Second * 1)
		go genSocket(ifi.Index, b, bTemp, test.count_frames, thr_step, counter, counterRes)
		fmt.Println("**")

		fmt.Println("***")
		PacketsTx := <-counter
		PacketsTxRes := <-counterRes
		//quit <- 1
		runtime.Gosched()
		time.Sleep(time.Second * 3)
		testRez.rez_64 = float32(PacketsTxRes-uint64(test.numberRx)) / float32(PacketsTxRes)

		fmt.Println(" Packet size - ", size_p)
		fmt.Println(" Send Packets - ", PacketsTx)
		fmt.Println(" Send PacketsRes - ", PacketsTxRes)
		fmt.Println(" Receive Packets - ", test.numberRx)

		size_p = 128
		test.numberRx = 0
		test.numberRxRes = 0
		b = packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size_p, 0, test.id_test_type, test.test_type)
		bTemp = packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size_p, 0, testTypeTemp, test.test_type)
		go test.receivePacketsLoss(ifi.MTU, quit)
		time.Sleep(time.Second * 1)
		go genSocket(ifi.Index, b, bTemp, test.count_frames, thr_step, counter, counterRes)
		PacketsTx = <-counter
		PacketsTxRes = <-counterRes
		//quit <- 1
		runtime.Gosched()
		time.Sleep(time.Second * 3)
		testRez.rez_128 = float32(PacketsTxRes-uint64(test.numberRx)) / float32(PacketsTxRes)

		fmt.Println(" Packet size - ", size_p)
		fmt.Println(" Send Packets - ", PacketsTx)
		fmt.Println(" Send PacketsRes - ", PacketsTxRes)
		fmt.Println(" Receive Packets - ", test.numberRx)

		size_p = 256
		test.numberRx = 0
		test.numberRxRes = 0
		b = packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size_p, 0, test.id_test_type, test.test_type)
		bTemp = packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size_p, 0, testTypeTemp, test.test_type)
		go test.receivePacketsLoss(ifi.MTU, quit)
		time.Sleep(time.Second * 1)
		go genSocket(ifi.Index, b, bTemp, test.count_frames, thr_step, counter, counterRes)

		PacketsTx = <-counter
		PacketsTxRes = <-counterRes
		//quit <- 1
		time.Sleep(time.Second * 3)
		testRez.rez_256 = float32(PacketsTxRes-uint64(test.numberRx)) / float32(PacketsTxRes)

		fmt.Println(" Packet size - ", size_p)
		fmt.Println(" Send Packets - ", PacketsTx)
		fmt.Println(" Send PacketsRes - ", PacketsTxRes)
		fmt.Println(" Receive Packets - ", test.numberRx)

		size_p = 512
		test.numberRx = 0
		test.numberRxRes = 0
		b = packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size_p, 0, test.id_test_type, test.test_type)
		bTemp = packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size_p, 0, testTypeTemp, test.test_type)
		go test.receivePacketsLoss(ifi.MTU, quit)
		time.Sleep(time.Second * 1)
		go genSocket(ifi.Index, b, bTemp, test.count_frames, thr_step, counter, counterRes)
		PacketsTx = <-counter
		PacketsTxRes = <-counterRes
		//quit <- 1
		time.Sleep(time.Second * 3)
		testRez.rez_512 = float32(PacketsTxRes-uint64(test.numberRx)) / float32(PacketsTxRes)

		fmt.Println(" Packet size - ", size_p)
		fmt.Println(" Send Packets - ", PacketsTx)
		fmt.Println(" Send PacketsRes - ", PacketsTxRes)
		fmt.Println(" Receive Packets - ", test.numberRx)

		size_p = 1024
		test.numberRx = 0
		test.numberRxRes = 0
		b = packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size_p, 0, test.id_test_type, test.test_type)
		bTemp = packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size_p, 0, testTypeTemp, test.test_type)
		go test.receivePacketsLoss(ifi.MTU, quit)
		time.Sleep(time.Second * 1)
		go genSocket(ifi.Index, b, bTemp, test.count_frames, thr_step, counter, counterRes)
		PacketsTx = <-counter
		PacketsTxRes = <-counterRes
		//quit <- 1
		time.Sleep(time.Second * 3)
		testRez.rez_1024 = float32(PacketsTxRes-uint64(test.numberRx)) / float32(PacketsTxRes)

		fmt.Println(" Packet size - ", size_p)
		fmt.Println(" Send Packets - ", PacketsTx)
		fmt.Println(" Send PacketsRes - ", PacketsTxRes)
		fmt.Println(" Receive Packets - ", test.numberRx)

		size_p = 1280
		test.numberRx = 0
		test.numberRxRes = 0
		b = packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size_p, 0, test.id_test_type, test.test_type)
		bTemp = packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size_p, 0, testTypeTemp, test.test_type)
		go test.receivePacketsLoss(ifi.MTU, quit)
		time.Sleep(time.Second * 1)
		go genSocket(ifi.Index, b, bTemp, test.count_frames, thr_step, counter, counterRes)
		PacketsTx = <-counter
		PacketsTxRes = <-counterRes
		//quit <- 1
		time.Sleep(time.Second * 3)
		testRez.rez_1280 = float32(PacketsTxRes-uint64(test.numberRx)) / float32(PacketsTxRes)

		fmt.Println(" Packet size - ", size_p)
		fmt.Println(" Send Packets - ", PacketsTx)
		fmt.Println(" Send PacketsRes - ", PacketsTxRes)
		fmt.Println(" Receive Packets - ", test.numberRx)

		size_p = 1518
		test.numberRx = 0
		test.numberRxRes = 0
		b = packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, 1500, 0, test.id_test_type, test.test_type)
		bTemp = packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, 1500, 0, testTypeTemp, test.test_type)
		go test.receivePacketsLoss(ifi.MTU, quit)
		time.Sleep(time.Second * 1)
		go genSocket(ifi.Index, b, bTemp, test.count_frames, thr_step, counter, counterRes)
		PacketsTx = <-counter
		PacketsTxRes = <-counterRes
		//quit <- 1
		time.Sleep(time.Second * 3)
		testRez.rez_1518 = float32(PacketsTxRes-uint64(test.numberRx)) / float32(PacketsTxRes)

		fmt.Println(" Packet size - ", size_p)
		fmt.Println(" Send Packets - ", PacketsTx)
		fmt.Println(" Send PacketsRes - ", PacketsTxRes)
		fmt.Println(" Receive Packets - ", test.numberRx)

		db, err = sql.Open("mysql", db_user+":"+db_user_pass+"@/"+db_database)
		if err != nil {
			db.Close()
			fmt.Println(" -!! Error !!-")
			fmt.Println(err)
			fmt.Println(" ----=====----")
			return
		}
		db.Exec("INSERT INTO test_frame_loss_rez (id_test, step_number, rez_64, rez_128, rez_256, rez_512, rez_1024, rez_1280, rez_1518, rez_4096, rez_9000) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", testRez.id_test, testRez.step_number, testRez.rez_64, testRez.rez_128, testRez.rez_256, testRez.rez_512, testRez.rez_1024, testRez.rez_1280, testRez.rez_1518, testRez.rez_4096, testRez.rez_9000)
		db.Close()
	}
	test.status = 3

	db, err = sql.Open("mysql", db_user+":"+db_user_pass+"@/"+db_database)
	if err != nil {
		db.Close()
		fmt.Println(" -!! Error !!-")
		fmt.Println(err)
		fmt.Println(" ----=====----")
		return
	}
	db.Exec("UPDATE test_frame_loss SET datetime_end=?, status=? WHERE id=?", time.Now().Format("2006-01-02 15:04:05"), test.status, id)
	db.Close()

}

func (test *testLoss) receivePacketsLoss(mtu int, quit chan int64) { //, counter chan<- int) {

	ifi, err := net.InterfaceByName(test.net_interface_name)
	if err != nil {
		fmt.Println("failed to find interface %q: %v", test.net_interface_name, err)
		return
	}

	var netConf *raw.Config = new(raw.Config)

	(*netConf).Filter, _ = bpf.Assemble([]bpf.Instruction{
		// Проверка идентификатора пакета (34 бит) (xFA-от 1 ко 2, xFB – от 2 к 1, xFC – от 1 к Серверу)
		bpf.LoadAbsolute{Off: 34, Size: 1},
		bpf.JumpIf{Cond: bpf.JumpNotEqual, Val: 0xFC, SkipTrue: 3},
		// Проверка идентификатора теста
		bpf.LoadAbsolute{Off: 64, Size: 2},
		bpf.JumpIf{Cond: bpf.JumpNotEqual, Val: uint32(test.id_test_type), SkipTrue: 1},
		// Verdict is "send up to 4k of the packet to userspace."
		bpf.RetConstant{Val: 4096},
		// Verdict is "ignore packet."
		bpf.RetConstant{Val: 0},
	})

	c, err := raw.ListenPacket(ifi, etherType, netConf)

	//	var f ethernet.Frame
	b := make([]byte, mtu)
	var ips [4]byte
	copy(ips[:], (test.ipdst1).To4())

	var t_ips [2]byte
	t_ips[1] = byte(test.id_test_type & 0xFF)
	t_ips[0] = byte((test.id_test_type >> 8) & 0xFF)
	start := time.Now()
	c.SetReadDeadline(start.Add(time.Second * time.Duration(2+test.count_frames)))
	//debug.SetGCPercent(-1)
	fmt.Println("Start receive: ", time.Now())
	for {
		select {
		case <-quit:
			//quit <- k
			//runtime.GC()
			fmt.Println("End receive: ", time.Now())
			fmt.Println("Packets receive = ", test.numberRx)
			return
		default:

			_, _, err := c.ReadFrom(b)
			if err != nil {
				fmt.Println("receivePacketsLoss: failed to receive message: ", err)
				fmt.Println("receivePacketsLoss: ", time.Now())
				if err.Error() != "resource temporarily unavailable" {
					c.SetReadDeadline(start.Add(time.Hour * 24))
					//runtime.GC()
					quit <- 1
					runtime.Gosched()
					continue
				}
				continue
			}

			//(*test).numberRx++
			atomic.AddUint64(&test.numberRx, uint64(1))
		}
	}
}
