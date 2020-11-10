package main

import (
	"bytes"
	"database/sql"
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
	//"runtime/debug"
	"time"
)

func TestY1564(id int, net_interface_name string) { //Нагрузочное тестирование задержки
	fmt.Println("Тест Y-1564 начался")
	db, err := sql.Open("mysql", db_user+":"+db_user_pass+"@/"+db_database)
	if err != nil {
		db.Close()
		fmt.Println(" -!! Error !!-")
		fmt.Println(err)
		fmt.Println(" ----=====----")
		return
	}
	db.Exec("UPDATE test_y1564 SET status=?, datatime=? WHERE id=?", time.Now().Format("2006-01-02 15:04:05"), 2, id) // Тест выполняется
	ifi, err := net.InterfaceByName(net_interface_name)
	if err != nil {
		db.Close()
		//	log.Fatalf("failed to find interface %q: %v", net_interface_name, err)
		db.Exec("UPDATE test_y1564 SET status=? WHERE id=?", 4, id) // Ошибка выполнения
		return
	}

	row, err := db.Query("select id, test_type, module_first, module_second, block_size, ToS, VLAN_priority, CIR, EIR, TP, period, step_count, max_FTD, max_FVD, max_FLR, status from test_y1564 where id=?", id)
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
	test := new(testY1564)
	err = row.Scan(&test.id, &test.test_type, &test.module_first, &test.module_second, &test.block_size, &test.ToS, &test.VLAN_priority, &test.CIR, &test.EIR, &test.TP, &test.period, &test.step_count, &test.max_FTD, &test.max_FVD, &test.max_FLR, &test.status)
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
		db.Exec("UPDATE test_y1564 SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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

			db.Exec("UPDATE test_y1564 SET status=? WHERE id=?", 4, id) // Ошибка выполнения
			db.Close()
			row.Close()
			fmt.Println(" -!! Error !!-")
			fmt.Println(err)
			fmt.Println(" ----=====----")
			return
		}
	}
	var id_sfp1, id_sfp2 int
	row, err = db.Query("SELECT module_first, module_second FROM test_y1564 WHERE id=?", id)
	if err != nil {
		db.Exec("UPDATE test_y1564 SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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

			db.Exec("UPDATE test_y1564 SET status=? WHERE id=?", 4, id) // Ошибка выполнения
			db.Close()
			row.Close()
			fmt.Println(" -!! Error !!-")
			fmt.Println(err)
			fmt.Println(" ----=====----")
			return
		}
		row_ip, err := db.Query("SELECT address_ip FROM modules_sfp_sla WHERE id=?", id_sfp1)
		if err != nil {
			db.Exec("UPDATE test_y1564 SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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

				db.Exec("UPDATE test_y1564 SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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
			db.Exec("UPDATE test_y1564 SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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

				db.Exec("UPDATE test_y1564 SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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
			db.Exec("UPDATE test_y1564 SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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

				db.Exec("UPDATE test_y1564 SET status=? WHERE id=?", 4, id) // Ошибка выполнения
				db.Close()
				row.Close()
				fmt.Println(" -!! Error !!-")
				fmt.Println(err)
				fmt.Println(" ----=====----")
				return
			}
		}
	}
	db.Exec("UPDATE test_y1564 SET status=?, datetime_start=? WHERE id=?", 2, time.Now().Format("2006-01-02 15:04:05"), id) // Тест выполняется
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

	PacketsRx = 0
	test.numberRx = 0

	if test.step_count > 1 {

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

		PacketsRx = 0
		test.numberRx = 0

	}

	if test.step_count > 2 {

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

		PacketsRx = 0
		test.numberRx = 0

	}

	if test.step_count > 3 {

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

		PacketsRx = 0
		test.numberRx = 0

	}

	runtime.Gosched()

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

	PacketsRx = 0
	test.numberRx = 0

	runtime.Gosched()

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

	PacketsRx = 0
	test.numberRx = 0

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
	db.Exec("UPDATE test_y1564 SET rez_IR_s1=?, rez_FTD_s1=?,rez_FVD_s1=?,rez_FLR_s1=?,rez_IR_s2=?, rez_FTD_s2=?,rez_FVD_s2=?,rez_FLR_s2=?, rez_IR_s3=?, rez_FTD_s3=?,rez_FVD_s3=?,rez_FLR_s3=?,rez_IR_s4=?, rez_FTD_s4=?,rez_FVD_s4=?,rez_FLR_s4=?, rez_IR_eir=?, rez_FTD_eir=?,rez_FVD_eir=?,rez_FLR_eir=?, rez_IR_tp=?, rez_FTD_tp=?,rez_FVD_tp=?,rez_FLR_tp=?, datetime_end=?, status=? WHERE id=?", test.rez_IR_s1, test.rez_FTD_s1, test.rez_FVD_s1, test.rez_FLR_s1, test.rez_IR_s2, test.rez_FTD_s2, test.rez_FVD_s2, test.rez_FLR_s2, test.rez_IR_s3, test.rez_FTD_s3, test.rez_FVD_s3, test.rez_FLR_s3, test.rez_IR_s4, test.rez_FTD_s4, test.rez_FVD_s4, test.rez_FLR_s4, test.rez_IR_eir, test.rez_FTD_eir, test.rez_FVD_eir, test.rez_FLR_eir, test.rez_IR_tp, test.rez_FTD_tp, test.rez_FVD_tp, test.rez_FLR_tp, time.Now().Format("2006-01-02 15:04:05"), test.status, id)

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
		//	t_time := int64(float64(time.Now().UnixNano())*float64(math.Pow(2, 32)/1000000000)) + 0xAABA4000000000
		//t_time := int64(float64(time.Now().UnixNano())*float64(math.Pow(2, 32)/1000000000)) - 0x55817F00000000
		//	t_time := int64(float64(time.Now().UnixNano())*float64(math.Pow(2, 32)/1000000000))
		//t_time = t_time << (4*8)
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

		fmt.Println(" -!! Error !!-")
		fmt.Println(err)
		fmt.Println(" ----=====----")
		return []byte{}
	}
	//	fmt.Print(" ==> Packet Form - ")
	//	fmt.Println(time.Now())
	return b
}

/*
func (test *testY1564) genFramesY1564(thr int, counter chan int64) int64 {
	time.Sleep(1*time.Millisecond)
	number := uint32(0)
	period_nano := int64(test.block_size * 8 * 1000 / thr)

	fmt.Printf("\n   period_nano = %d \n   thr = %d\n", period_nano, thr)

	var counter_rez int64
	ifi, err := net.InterfaceByName(test.net_interface_name)
	c, err := raw.ListenPacket(ifi, etherType, nil)
	if err != nil {
		fmt.Println("failed to listen: %v", err)
		fmt.Println(" -!! Error !!-")
		fmt.Println(err)
		fmt.Println(" ----=====----")
		return 0
	}
	//c.SetReadDeadline(time.Now().Add(time.Millisecond * 3000))
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
			//		fmt.Println("failed to listen: %v", err)
			//		time.Sleep(time.Millisecond * 1)
			//	}
				counter_rez = counter_rez + 1
			}
		}
	}()
	time.Sleep(time.Duration(test.period) * time.Second)
	ticker.Stop()
	done <- true
	time.Sleep(1 * time.Second)
	fmt.Println("Packed send - ", counter_rez)
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
	var floatDelay, floatDelayMax, floatDelayMin float32

	detectPackDelay := make(chan int64, 10)
	number := 0

	timeStart := time.Now()

	ifi, err := net.InterfaceByName(test.net_interface_name)
	if err != nil {
		fmt.Println("failed to find interface %q: %v", test.net_interface_name, err)
		quit <- 1
		return 0, 0
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
	quit_receive := make(chan int64, 2)
	go (*test).receivePacketsY(ifi.MTU, quit_receive, test.id_test_type)

	SolveDelayTicker := time.NewTicker(1 * time.Millisecond)
	//	b := packetForm(test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, test.block_size, 1, test.id_test_type, test.test_type)
	go test.receiveMessagesDelay(detectPackDelay, c, ifi.MTU, test.id_test_type, test.block_size)

	for range SolveDelayTicker.C {
		//	for	 {
		//	go test.receiveMessagesDelay(detectPackDelay, c, ifi.MTU, test.id_test_type, test.block_size)

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
			if time.Since(timeStart) > time.Duration(time.Duration(test.period)*time.Second) {
				fmt.Println(" --> Number = ", number)
				fmt.Println(" --> Size = ", test.block_size)

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
				//time.Sleep(100 * time.Millisecond)
				quit <- 0
				return floatDelay, float32(math.Max(float64(floatDelayMax-floatDelay), float64(floatDelay-floatDelayMin)))
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
	//fmt.Println("-> Begin Catch - ", start)
	c.SetReadDeadline(start.Add(time.Second * time.Duration(test.period)))
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
			at_time :=time.Now()
			cc++
			if err != nil {
				//fmt.Printf("failed to receive message: %v", err)
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

			//t_time := int64(float64(time.Now().UnixNano() - delta_nano )*float64(math.Pow(2, 32)/1000000000)) - 0x55817F00000000
		    delta_nano := int64((2208988800) * 1000000000)
		
			t_time := int64(float64(at_time.UnixNano()-delta_nano) * float64(math.Pow(2, 32)/1000000000))
			t_time = t_time & int64(0xFFFFFFFFFFFFFF)
			//fmt.Println("  Now real - ",time.Now().UnixNano());
			
			//n, addr, err := c.ReadFrom(b)
			// Unpack Ethernet II frame into Go representation.
			if err := (&f).UnmarshalBinary(b[:n]); err != nil {
				fmt.Printf("failed to unmarshal ethernet frame: %v", err)
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
	c.SetReadDeadline(start.Add(time.Second * time.Duration(1+test.period)))
	//debug.SetGCPercent(-1)
	fmt.Println("Start receive: ", time.Now())
	for {
		select {
		case <-quit:
			//quit <- k
			//	runtime.GC()
			fmt.Println("End receive: ", time.Now())
			fmt.Println("Packets receive = ", test.numberRx)
			return
		default:
		}
		_, _, err := c.ReadFrom(b)
		if err != nil {
			fmt.Println("failed to receive message: ", err)
			if err.Error() == "i/o timeout" {
				//	(*test).number++

				c.SetReadDeadline(start.Add(time.Hour * 24))
				quit <- 1
				runtime.Gosched()
			}
			continue
		}
		//*/
		// Unpack Ethernet II frame into Go representation.
		//	if err := (&f).UnmarshalBinary(b[:n]); err != nil {
		//		fmt.Println("failed to unmarshal ethernet frame: %v", err)
		//	}
		//fmt.Printf("\n\n--=Test %x - \n", f.Payload[12:16])

		//fmt.Printf("\n\n--=T_so %x - \n", ips)

		// Display source of message and message itself.
		//	if (len(f.Payload) >= 52) && (f.Payload[20] == 0xFC) && (bytes.Equal(f.Payload[12:16], ips[:]) == true) && (bytes.Equal(f.Payload[50:52], t_ips[:]) == true) {
		//count++
		//	fmt.Printf("-->>Detect")
		//	counter <-count
		//(*test).numberCounter = uint32(count)
		atomic.AddUint64(&test.numberRx, uint64(1))
		//(*test).numberRx++
		//	}
	}
}
