package main

import (
	"math"
	"bytes"
	"encoding/binary"
	"database/sql"
	"fmt"
	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/raw"
	"net"
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
	db.Exec("UPDATE test_y1564 SET status=?, datatime=? WHERE id=?", time.Now(), 2, id) // Тест выполняется
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
	err = row.Scan(&test.id, &test.test_type, &test.module_first, &test.module_second, &test.block_size, &test.ToS, &test.VLAN_priority,&test.CIR,&test.EIR, &test.TP, &test.period, &test.step_count, &test.max_FTD, &test.max_FVD, &test.max_FLR, &test.status)
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
	db.Exec("UPDATE test_y1564 SET status=?, datetime_start=? WHERE id=?", 2, time.Now(), id) // Тест выполняется
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
	test.net_interface_name = net_interface_name
	test.mac_src = ifi.HardwareAddr

	//counter := make(chan int64, 7)
	//quit := make(chan int64, 7)

	//size_p := 64

	var ToS_tag uint8
	ToS_tag=((uint8(0x7&test.VLAN_priority))<<5)+uint8(test.ToS << 1)

	b := packetFormY1546(ToS_tag,test.ipsrc, test.ipdst1, test.ipdst2, test.mac_src, test.mac_dst, test.block_size, 0, test.id_test_type, test.test_type)
	//go genSocket(ifi.Index, b, test.count_packs, test.thr_begin, counter)
	//test.getMonDelay(quit, size_p)
	//time.Sleep(time.Second * 2)
	//<-quit

	test.rez_FLR_s2=0.02



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
	db.Exec("UPDATE test_y1564 SET rez_IR_s1=?, rez_FTD_s1=?,rez_FVD_s1=?,rez_FLR_s1=?,rez_IR_s2=?, rez_FTD_s2=?,rez_FVD_s2=?,rez_FLR_s2=?, rez_IR_s3=?, rez_FTD_s3=?,rez_FVD_s3=?,rez_FLR_s3=?,rez_IR_s4=?, rez_FTD_s4=?,rez_FVD_s4=?,rez_FLR_s4=?, rez_IR_eir=?, rez_FTD_eir=?,rez_FVD_eir=?,rez_FLR_eir=?, rez_IR_tp=?, rez_FTD_tp=?,rez_FVD_tp=?,rez_FLR_tp=?, datetime_end=?, status=? WHERE id=?", test.rez_IR_s1, test.rez_FTD_s1,test.rez_FVD_s1,test.rez_FLR_s1, test.rez_IR_s2, test.rez_FTD_s2,test.rez_FVD_s2,test.rez_FLR_s2, test.rez_IR_s3, test.rez_FTD_s3,test.rez_FVD_s3,test.rez_FLR_s3, test.rez_IR_s4, test.rez_FTD_s4,test.rez_FVD_s4,test.rez_FLR_s4, test.rez_IR_eir, test.rez_FTD_eir,test.rez_FVD_eir,test.rez_FLR_eir, test.rez_IR_tp, test.rez_FTD_tp,test.rez_FVD_tp,test.rez_FLR_tp, time.Now(), test.status, id)								   

	db.Close()
}

func packetFormY1546(ToS uint8,ipsrc net.IP, ipdst1 net.IP, ipdst2 net.IP, mac_src []byte, mac_dst []byte, size int, number uint32, test_type uint16, testWay int) []byte {
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