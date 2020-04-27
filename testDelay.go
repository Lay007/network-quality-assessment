package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"math"
	"net"
	"time"

	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/raw"
	"golang.org/x/net/bpf"
)

func TestDelay(id int, net_interface_name string) { //Нагрузочное тестирование задержки
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
	db.Exec("UPDATE test_latency SET rez_64=?,rez_64_max=?,rez_64_min=?,rez_128=?,rez_128_max=?,rez_128_min=?,rez_256=?,rez_256_max=?,rez_256_min=?,rez_512=?,rez_512_max=?,rez_512_min=?,rez_1024=?,rez_1024_max=?,rez_1024_min=?,rez_1280=?,rez_1280_max=?,rez_1280_min=?, rez_1518=?,rez_1518_max=?,rez_1518_min=?,datetime=?, status=? WHERE id=?", test.rez_64, test.rez_64_max, test.rez_64_min, test.rez_128, test.rez_128_max, test.rez_128_min, test.rez_256, test.rez_256_max, test.rez_256_min, test.rez_512, test.rez_512_max, test.rez_512_min, test.rez_1024, test.rez_1024_max, test.rez_1024_min, test.rez_1280, test.rez_1280_max, test.rez_1280_min, test.rez_1518, test.rez_1518_max, test.rez_1518_min, time.Now(), test.status, id)

	db.Close()
}

func (test *testDelay) getMonDelay(quit chan int64, size int) {

	var delay, delayMax, delayMin int64
	var floatDelay, floatDelayMax, floatDelayMin float32

	detectPackDelay := make(chan int64, 10)
	number := 0

	timeStart := time.Now()

	ifi, err := net.InterfaceByName(test.net_interface_name)
	if err != nil {
		fmt.Println("failed to find interface %q: %v", test.net_interface_name, err)
		quit <- 1
		return
	}

	var netConf *raw.Config = new(raw.Config)

	(*netConf).Filter, _ = bpf.Assemble([]bpf.Instruction{
		// Проверка идентификатора пакета (34 бит) (xFA-от 1 ко 2, xFB – от 2 к 1, xFC – от 1 к Серверу)
		bpf.LoadAbsolute{Off: 34, Size: 1},
		bpf.JumpIf{Cond: bpf.JumpNotEqual, Val: 0xFC, SkipTrue: 5},
		// Проверка идентификатора теста
		bpf.LoadAbsolute{Off: 64, Size: 2},
		bpf.JumpIf{Cond: bpf.JumpNotEqual, Val: uint32(test.id_test_type), SkipTrue: 3},
		// Выбор одного из 1000
		bpf.LoadExtension{Num: bpf.ExtRand},
		//	bpf.JumpIf{Cond: bpf.JumpLessThan, Val: 0xFF, SkipFalse: 1},
		bpf.JumpIf{Cond: bpf.JumpGreaterThan, Val: 0x03FFFFFF, SkipTrue: 1},
		// Verdict is "send up to 4k of the packet to userspace."
		bpf.RetConstant{Val: 4096},
		// Verdict is "ignore packet."
		bpf.RetConstant{Val: 0},
	})

	c, err := raw.ListenPacket(ifi, etherType, netConf)

	SolveDelayTicker := time.NewTicker(1 * time.Millisecond)
	b := packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, size, 1, test.id_test_type, test.test_type)

	for range SolveDelayTicker.C {

		go test.receiveMessagesDelay(detectPackDelay, c, ifi.MTU, test.id_test_type, len(b))

		select {
		//fmt.Println("Wait")
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
					if delayMin > detect {
						delayMin = detect
					}
				}

			}
		default:
			if time.Since(timeStart) > time.Duration(test.count_tests)*time.Second {
				fmt.Println(" --> Number = ", number)
				fmt.Println(" --> Size = ", size)
			
				floatDelay = (float32(delay) / float32(number)) * 1000000.0 / float32(math.Pow(2, 32))
				floatDelayMax = float32(delayMax) * 1000000.0 / float32(math.Pow(2, 32))
				floatDelayMin = float32(delayMin) * 1000000.0 / float32(math.Pow(2, 32))
/*
				fmt.Println(" --> Delay = ", floatDelay)
				fmt.Println(" --> DelayMax = ", floatDelayMax)
				fmt.Println(" --> DelayMin = ", floatDelayMin)

				floatDelay = (float32(delay) / float32(number))
				floatDelayMax = float32(delayMax)
				floatDelayMin = float32(delayMin)
*/
				fmt.Println(" --> Delay = ", floatDelay)
				fmt.Println(" --> DelayMax = ", floatDelayMax)
				fmt.Println(" --> DelayMin = ", floatDelayMin)
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
	//fmt.Println("-> Begin Catch - ", start)
	c.SetReadDeadline(start.Add(time.Microsecond * 3000))
	//ExitLoop:
	for {
		select {
		case key := <-quit:
			catchDetect <- key
			//	fmt.Println("Chanel go")
			return
			//break ExitLoop
		default:

			n, _, err := c.ReadFrom(b)
			cc++
			if err != nil {
				fmt.Printf("failed to receive message: %v", err)
				if err.Error() == "resource temporarily unavailable" {

				}
				//log.Fatalf("failed to receive message: %v", err)
				c.SetReadDeadline(start.Add(time.Hour * 24))
				quit <- 0
				continue
			}

			if time.Since(start) > (time.Millisecond * 1000) {
				quit <- 0
				continue
			}

			if (n) != packetSize {
				continue
			}

			//t_time := int64(float64(time.Now().UnixNano() - delta_nano )*float64(math.Pow(2, 32)/1000000000)) - 0x55817F00000000
			delta_nano := int64((2208988800) * 1000000000)
			t_time := int64(float64(time.Now().UnixNano()-delta_nano) * float64(math.Pow(2, 32)/1000000000))
			t_time = t_time & int64(0xFFFFFFFFFFFFFF)

			//n, addr, err := c.ReadFrom(b)
			// Unpack Ethernet II frame into Go representation.
			if err := (&f).UnmarshalBinary(b[:n]); err != nil {
				fmt.Printf("failed to unmarshal ethernet frame: %v", err)
				continue
			}

			var ips [4]byte
			copy(ips[:], (test.ipdst1).To4())

			if (len(f.Payload) >= 52) && (f.Payload[20] == 0xFC) && (bytes.Equal(f.Payload[12:16], ips[:]) == true) && (bytes.Equal(f.Payload[50:52], t_ips[:]) == true) {

				//fmt.Printf("\n\n--=Packet DETECT!!!=--\n")
				//fmt.Println(time.Now())
				//fmt.Printf("\n\n--=Test %x - \n -== %x\n",f.Payload[12:15],net.ParseIP(ipdst_1sfpsla_str))
				//fmt.Printf("size: %v raw:  %x \n", len(f.Payload), f.Payload)
				//fmt.Printf("\n\rEthernet source: [%s]\n", addr.String())

				//fmt.Printf("size     %x \n", b[2:4])

				//fmt.Printf("ip sourse %v.%v.%v.%v \n", f.Payload[12], f.Payload[13], f.Payload[14], f.Payload[15])
				//fmt.Printf("ip dst    %v.%v.%v.%v \n", f.Payload[16], f.Payload[17], f.Payload[18], f.Payload[19])

				//fmt.Printf("ip SFP2   %v.%v.%v.%v \n", f.Payload[21], f.Payload[22], f.Payload[23], f.Payload[24])

				//fmt.Printf("time marker_SFP1_1 :   %x \n", f.Payload[25:32])
				//fmt.Printf("time marker_SFP2   :   %x \n", f.Payload[32:39])
				//fmt.Printf("time marker_SFP1_2 :   %x \n", f.Payload[39:46])
				//fmt.Printf("Number marker      :   %x \n", f.Payload[46:50])
				//fmt.Println(" --== End Packet ==--")
				//*/
				var markerSFP11, markerSFP12, markerSFP2 int64
				var ind uint

				for ind = 0; ind < 7; ind++ {
					markerSFP11 = markerSFP11 + int64(f.Payload[31-ind])<<(8*ind)
					markerSFP2 = markerSFP2 + int64(f.Payload[38-ind])<<(8*ind)
					markerSFP12 = markerSFP12 + int64(f.Payload[45-ind])<<(8*ind)
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
