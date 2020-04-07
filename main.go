package main

import (
	. "./go-zabbix"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	//"runtime"
	"time"
	//"unsafe"

	"github.com/google/gopacket/pcap"

	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/raw"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const (
	db_user      = `sfp_user`
	db_user_pass = `rootsfp`
	db_database  = `server_sfp_sla`

	debugV = true

	//	defaultHost = `localhost`
	//defaultHost = `remote.fibertrade.ru`
	//	defaultPort = 10051
	etherType = 0x0800

//	ipsrcstr          = "10.0.10.115"
//	ipdst_1sfpsla_str = "10.0.10.172"
//	ipdst_2sfpsla_str = "10.0.10.175"
)

func checksum(buf []byte) uint16 {
	sum := uint32(0)

	for ; len(buf) >= 2; buf = buf[2:] {
		sum += uint32(buf[0])<<8 | uint32(buf[1])
	}
	if len(buf) > 0 {
		sum += uint32(buf[0]) << 8
	}
	for sum > 0xffff {
		sum = (sum >> 16) + (sum & 0xffff)
	}
	csum := ^uint16(sum)
	/*
	 * From RFC 768:
	 * If the computed checksum is zero, it is transmitted as all ones (the
	 * equivalent in one's complement arithmetic). An all zero transmitted
	 * checksum value means that the transmitter generated no checksum (for
	 * debugging or for higher level protocols that don't care).
	 */
	if csum == 0 {
		csum = 0xffff
	}
	return csum
}

func (h *iphdr) checksum() {
	(*h).csum = 0
	var b bytes.Buffer
	binary.Write(&b, binary.BigEndian, *h)
	(*h).csum = checksum(b.Bytes())
}

func main() {

	//runtime.GOMAXPROCS(1024)

	//time.Sleep(100 * time.Second)

	// подключение к БД и обновление списка сетевых интерфейсов
	db, err := sql.Open("mysql", db_user+":"+db_user_pass+"@/"+db_database)
	if err != nil {
		panic(err)
	}
	db.Exec("DELETE FROM net_interfaces_from_server_sla")
	db.Exec("ALTER TABLE net_interfaces_from_server_sla AUTO_INCREMENT = 1")

	devices, err := pcap.FindAllDevs() // считывание перечня сетевых интерфейсов
	fmt.Println(err)
	if err != nil {
		log.Fatal(err)
	}

	// запись перечня сетевых интерфейсов в БД
	for _, device := range devices {
		fmt.Println(device.Name)
		netInterface, err := net.InterfaceByName(device.Name)
		var addressMac net.HardwareAddr
		if err == nil {
			addressMac = netInterface.HardwareAddr
		}
		//net_name = device.Name
		for _, address := range device.Addresses {
			db.Exec("INSERT INTO net_interfaces_from_server_sla (name, address_IP, address_mac) VALUES(?, ?, ?)", device.Name, address.IP.String(), addressMac.String())
		}
	}

	// Оценка выполнения тестов
	db.Exec("UPDATE test_throughput SET status=4 WHERE status=2")
	db.Exec("UPDATE test_sla_real SET status=1 WHERE status=2")

	db.Close()

	t := time.NewTicker(10 * time.Second) //проверка один раз в 10 секунд
	for range t.C {

		db, err = sql.Open("mysql", db_user+":"+db_user_pass+"@/"+db_database)
		if err != nil {
			db.Close()
			fmt.Println(" -!! Error !!-")
			fmt.Println(err)
			fmt.Println(" ----=====----")
			continue
		}
		// считывание из БД глобальных параметров
		row_gc, err := db.Query("select * from global_config")
		defer row_gc.Close()
		if err != nil {
			db.Close()
			fmt.Println(" -!! Error !!-")
			fmt.Println(err)
			fmt.Println(" ----=====----")
			continue
		}
		row_gc.Next()
		conf := new(global_config)
		err = row_gc.Scan(&conf.server_ip, &conf.net_interface_name, &conf.zabbix_server_name, &conf.zabbix_server_port, &conf.vlan, &conf.vlan_number, &conf.QinQ, &conf.QinQ_number)
		if err != nil {
			db.Close()

			fmt.Println(" -!! Error !!-")
			fmt.Println(err)
			fmt.Println(" ----=====----")
			continue
		}
		row_gc.Close()

		if conf.net_interface_name == "" {
			db.Close()
			fmt.Println("Net interface is not selected")
			continue
		}

		fmt.Println(conf.server_ip)

		row_modules, err := db.Query("select * from modules_sfp_sla")
		defer row_modules.Close()
		if err != nil {
			db.Close()

			fmt.Println(" -!! Error !!-")
			fmt.Println(err)
			fmt.Println(" ----=====----")
			continue
		}

		modules := []module_sfp{}

		for row_modules.Next() {
			m := module_sfp{}
			err = row_modules.Scan(&m.id, &m.addres_mac, &m.name, &m.address_ip, &m.version, &m.location)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(m.address_ip)
			modules = append(modules, m)
		}
		row_modules.Close()
		//defer db.Close()

		//go zabbixHello("SFP-SLA_4401")

		if conf.net_interface_name == "0" {
			fmt.Println("  Net interface not changed")
			db.Close()
			continue
		}

		ifi, err := net.InterfaceByName(conf.net_interface_name)
		if err != nil {
			//log.Fatalf("failed to find interface %q: %v", conf.net_interface_name, err)
			db.Close()
			fmt.Println(" -!! Error !!-")
			fmt.Println("failed to find interface %q: %v", conf.net_interface_name, err)
			fmt.Println(" ----=====----")
			continue
		}

		fmt.Println("Net_NAME: ", conf.net_interface_name)
		fmt.Println("interface: ", ifi.Name)
		fmt.Println("Time now: ", time.Now())

		// проверка тестов пропускной способности
		row_test_thr, err := db.Query("SELECT id FROM test_throughput WHERE status=1")
		defer row_test_thr.Close()
		if err != nil {
			db.Close()
			row_test_thr.Close()
			fmt.Println(" -!! Error !!-")
			fmt.Println(err)
			fmt.Println(" ----=====----")
			continue
		}

		// проверка тестов SLA в реальном времени
		row_test_real, err := db.Query("SELECT id FROM test_sla_real WHERE status=1")
		defer row_test_real.Close()
		if err != nil {
			db.Close()
			row_test_real.Close()
			fmt.Println(" -!! Error !!-")
			fmt.Println(err)
			fmt.Println(" ----=====----")
			continue
		}

		db.Close()

		for row_test_real.Next() {
			var id int
			err = row_test_real.Scan(&id)
			if err != nil {
				db.Close()
				row_test_real.Close()
				fmt.Println(" -!! Error !!-")
				fmt.Println(err)
				fmt.Println(" ----=====----")
				continue
			}
			fmt.Println("-== Test id  = ", id)
			go TestReal(id, conf.net_interface_name, conf.zabbix_server_name, conf.zabbix_server_port)
		}
		row_test_real.Close()

		for row_test_thr.Next() {
			var id int
			err = row_test_thr.Scan(&id)
			if err != nil {
				db.Close()
				row_test_thr.Close()
				fmt.Println(" -!! Error !!-")
				fmt.Println(err)
				fmt.Println(" ----=====----")
				continue
			}
			TestThroughput(id, conf.net_interface_name)
		}
		row_test_thr.Close()
		//	row_test_real, err := db.Query("select * from global_config where status=1")

		//go Test_SLA_real_go()

	}

	// -------=======
	/*
		c, err := raw.ListenPacket(ifi, etherType, nil)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		// Send messages in one goroutine, receive messages in another.
		go sendMessages(c, ifi.HardwareAddr)
		go receiveMessages(c, ifi.MTU)

		// Block forever.
		select {}
	*/
}

var count_recive int

func TestThroughput(id int, net_interface_name string) { //Нагрузочное тестирование пропускной способности
	fmt.Println("Тест пропускной способности начался")
	db, err := sql.Open("mysql", db_user+":"+db_user_pass+"@/"+db_database)
	if err != nil {
		db.Close()
		fmt.Println(" -!! Error !!-")
		fmt.Println(err)
		fmt.Println(" ----=====----")
		return
	}
	db.Exec("UPDATE test_throughput SET status=?, datatime=? WHERE id=?", time.Now(), 2, id) // Тест выполняется
	ifi, err := net.InterfaceByName(net_interface_name)
	if err != nil {
		db.Close()
		log.Fatalf("failed to find interface %q: %v", net_interface_name, err)
		db.Exec("UPDATE test_throughput SET status=? WHERE id=?", 4, id) // Ошибка выполнения
		return
	}

	row, err := db.Query("select id, test_type, module_first, module_second, thr_begin, count, ch_type, max_loss, status from test_throughput where id=?", id)
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
	test := new(testThroughput)
	err = row.Scan(&test.id, &test.test_type, &test.module_first, &test.module_second, &test.thr_begin, &test.count, &test.ch_type, &test.max_loss, &test.status)
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
		db.Exec("UPDATE test_throughput SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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

			db.Exec("UPDATE test_throughput SET status=? WHERE id=?", 4, id) // Ошибка выполнения
			db.Close()
			row.Close()
			fmt.Println(" -!! Error !!-")
			fmt.Println(err)
			fmt.Println(" ----=====----")
			return
		}
	}
	var id_sfp1, id_sfp2 int
	row, err = db.Query("SELECT module_first, module_second FROM test_throughput WHERE id=?", id)
	if err != nil {
		db.Exec("UPDATE test_throughput SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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

			db.Exec("UPDATE test_throughput SET status=? WHERE id=?", 4, id) // Ошибка выполнения
			db.Close()
			row.Close()
			fmt.Println(" -!! Error !!-")
			fmt.Println(err)
			fmt.Println(" ----=====----")
			return
		}
		row_ip, err := db.Query("SELECT address_ip FROM modules_sfp_sla WHERE id=?", id_sfp1)
		if err != nil {
			db.Exec("UPDATE test_throughput SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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

				db.Exec("UPDATE test_throughput SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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
			db.Exec("UPDATE test_sla_real SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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

				db.Exec("UPDATE test_sla_real SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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
			db.Exec("UPDATE test_throughput SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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

				db.Exec("UPDATE test_throughput SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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

	c, err := raw.ListenPacket(ifi, etherType, nil)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		db.Close()

		fmt.Println(" -!! Error !!-")
		fmt.Println(err)
		fmt.Println(" ----=====----")
		return
	}

	ipsrc := net.ParseIP(ipsrcstr)
	ipdst1 := net.ParseIP(ipdst_1sfpsla_str)
	ipdst2 := net.ParseIP(ipdst_2sfpsla_str)

	//period_min := time.Duration(time.Duration(int(period_nano)) * time.Nanosecond)
	//period_gen := time.Duration(10 * time.Second)

	//t := time.NewTicker(time.Duration(int(period_nano)) * time.Nanosecond)
	//t := time.NewTicker(1 * time.Second)
	//	for range t.C {
	//		counter--
	//fmt.Println("counter= ", counter)
	//fmt.Println("period_min= ", period_min)
	//	numberTX++
	addr := &raw.Addr{
		HardwareAddr: ethernet.Broadcast,
	}
	var number uint32
	var test_type uint16
	test_type = 0xA000 + (uint16(id) & 0x1FFF)

	var test_c testThr
	test_c.numberCounter = uint32(test.count)

	period_test := test.count // период теста - 10 секунд
	size := 64
	fmt.Println("->> test.thr_begin = ", test.thr_begin)
	period_nano := int64(size*8*1000000000) / (int64(test.thr_begin * 1000 * 1000))
	packet_count := (int64(period_test * 1000000000)) / period_nano

	fmt.Println("->> period_nano  = ", period_nano)
	fmt.Println("->> packet_count = ", packet_count)

	b := packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, 64, number, test_type, test.test_type)
	//count_rez, per := test_c.testMax(b, c, addr, ifi.MTU, ipdst_1sfpsla_str, counter, test_type)
	count_rez, per := test_c.testThrGen(b, c, addr, ifi.MTU, ipdst_1sfpsla_str, packet_count, period_nano, test_type)
	test.rez_64 = (float32)(64.0 * 8.0 * (float32)(count_rez) * 1000000 / (float32)(per))
	fmt.Println("->> rez_64 = ", count_rez, " period = ", per, " !!!  rez=", test.rez_64)

	size = 128
	period_nano = int64(size * 8 * 1000000000 / (test.thr_begin * 1000 * 1000))
	packet_count = (int64(period_test * 1000000000)) / period_nano
	test_c.numberCounter = 0
	b = packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, 128, number, test_type, test.test_type)

	//count_rez, per = test_c.testMax(b, c, addr, ifi.MTU, ipdst_1sfpsla_str, counter, test_type)
	count_rez, per = test_c.testThrGen(b, c, addr, ifi.MTU, ipdst_1sfpsla_str, packet_count, period_nano, test_type)
	test.rez_128 = (float32)(128.0 * 8.0 * (float32)(count_rez) * 1000000 / (float32)(per))
	fmt.Println("->> rez_128 = ", count_rez, " period = ", per, " !!!  rez=", test.rez_128)

	size = 256
	period_nano = int64(size * 8 * 1000000000 / (test.thr_begin * 1000 * 1000))
	packet_count = (int64(period_test * 1000000000)) / period_nano
	test_c.numberCounter = 0
	b = packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, 256, number, test_type, test.test_type)

	//count_rez, per = test_c.testMax(b, c, addr, ifi.MTU, ipdst_1sfpsla_str, counter, test_type)
	count_rez, per = test_c.testThrGen(b, c, addr, ifi.MTU, ipdst_1sfpsla_str, packet_count, period_nano, test_type)
	test.rez_256 = (float32)(256.0 * 8.0 * (float32)(count_rez) * 1000000 / (float32)(per))
	fmt.Println("->> rez_256 = ", count_rez, " period = ", per, " !!!  rez=", test.rez_256)

	size = 512
	period_nano = int64(size * 8 * 1000000000 / (test.thr_begin * 1000 * 1000))
	packet_count = (int64(period_test * 1000000000)) / period_nano
	test_c.numberCounter = 0
	b = packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, 512, number, test_type, test.test_type)

	//count_rez, per = test_c.testMax(b, c, addr, ifi.MTU, ipdst_1sfpsla_str, counter, test_type)
	count_rez, per = test_c.testThrGen(b, c, addr, ifi.MTU, ipdst_1sfpsla_str, packet_count, period_nano, test_type)
	test.rez_512 = (float32)(512.0 * 8.0 * (float32)(count_rez) * 1000000 / (float32)(per))
	fmt.Println("->> rez_512 = ", count_rez, " period = ", per, " !!!  rez=", test.rez_512)

	size = 1024
	period_nano = int64(size * 8 * 1000000000 / (test.thr_begin * 1000 * 1000))
	packet_count = (int64(period_test * 1000000000)) / period_nano
	test_c.numberCounter = 0
	b = packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, size, number, test_type, test.test_type)
	count_rez, per = test_c.testThrGen(b, c, addr, ifi.MTU, ipdst_1sfpsla_str, packet_count, period_nano, test_type)
	test.rez_1024 = (float32)(float32(size) * 8.0 * (float32)(count_rez) * 1000000 / (float32)(per))
	fmt.Println("->> rez_1024 = ", count_rez, " period = ", per, " !!!  rez=", test.rez_1024)

	size = 1280
	period_nano = int64(size * 8 * 1000000000 / (test.thr_begin * 1000 * 1000))
	packet_count = (int64(period_test * 1000000000)) / period_nano
	test_c.numberCounter = 0
	b = packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, size, number, test_type, test.test_type)
	count_rez, per = test_c.testThrGen(b, c, addr, ifi.MTU, ipdst_1sfpsla_str, packet_count, period_nano, test_type)
	test.rez_1280 = (float32)(float32(size) * 8.0 * (float32)(count_rez) * 1000000 / (float32)(per))
	fmt.Println("->> rez_1280 = ", count_rez, " period = ", per, " !!!  rez=", test.rez_1280)

	size = 1500
	period_nano = int64(size * 8 * 1000000000 / (test.thr_begin * 1000 * 1000))
	packet_count = (int64(period_test * 1000000000)) / period_nano
	test_c.numberCounter = 0
	b = packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, size, number, test_type, test.test_type)
	count_rez, per = test_c.testThrGen(b, c, addr, ifi.MTU, ipdst_1sfpsla_str, packet_count, period_nano, test_type)
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
	db.Exec("UPDATE test_throughput SET rez_64=?,rez_128=?,rez_256=?,rez_512=?,rez_1024=?,rez_1280=?, rez_1518=?, status=? WHERE id=?", test.rez_64, test.rez_128, test.rez_256, test.rez_512, test.rez_1024, test.rez_1280, test.rez_1518, test.status, id)
	
	db.Close()
}

func TestReal(id int, net_interface_name string, host_zabbix string, port_zabbix int) {
	db, err := sql.Open("mysql", db_user+":"+db_user_pass+"@/"+db_database)
	if err != nil {
		db.Close()
		fmt.Println(" -!! Error !!-")
		fmt.Println(err)
		fmt.Println(" ----=====----")
		return
	}
	db.Exec("UPDATE test_sla_real SET status=? WHERE id=?", 2, id) // Тест выполняется
	ifi, err := net.InterfaceByName(net_interface_name)
	if err != nil {
		db.Close()
		log.Fatalf("failed to find interface %q: %v", net_interface_name, err)
		db.Exec("UPDATE test_sla_real SET status=? WHERE id=?", 4, id) // Ошибка выполнения
		return
	}

	db.Exec("UPDATE test_sla_real SET data_start=? WHERE id=?", time.Now(), id) // Добавление времени начала
	row, err := db.Query("SELECT id, test_type, module_first, module_second, block_size, clock, count, node_zabbix, test_delay,test_delay_jitter, test_loss, test_delay_1,test_delay1_jitter FROM test_sla_real WHERE id=?", id)
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
	test := new(testReal)
	err = row.Scan(&test.id, &test.test_type, &test.module_first, &test.module_second, &test.block_size, &test.clock, &test.count, &test.node_zabbix, &test.test_delay, &test.test_delay_jitter, &test.test_loss, &test.test_delay_1, &test.test_delay_1_jitter)
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

	row, err = db.Query("SELECT server_IP FROM global_config")
	if err != nil {
		db.Exec("UPDATE test_sla_real SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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

			db.Exec("UPDATE test_sla_real SET status=? WHERE id=?", 4, id) // Ошибка выполнения
			db.Close()
			row.Close()
			fmt.Println(" -!! Error !!-")
			fmt.Println(err)
			fmt.Println(" ----=====----")
			return
		}
	}
	var id_sfp1, id_sfp2 int
	var testMax testRealMax

	mac_dst := make([]byte, 6)
	row, err = db.Query("SELECT module_first, module_second, delay_max, jitter_max, delay1_max, jitter1_max, loss_max FROM test_sla_real WHERE id=?", id)
	if err != nil {
		db.Exec("UPDATE test_sla_real SET status=? WHERE id=?", 4, id) // Ошибка выполнения
		db.Close()
		row.Close()
		fmt.Println(" -!! Error !!-")
		fmt.Println(err)
		fmt.Println(" ----=====----")
		return
	}
	for row.Next() {
		err = row.Scan(&id_sfp1, &id_sfp2, &testMax.delayMax, &testMax.jitterMax, &testMax.delayOneMax, &testMax.jitterOneMax, &testMax.lossMax)
		if err != nil {

			db.Exec("UPDATE test_sla_real SET status=? WHERE id=?", 4, id) // Ошибка выполнения
			db.Close()
			row.Close()
			fmt.Println(" -!! Error !!-")
			fmt.Println(err)
			fmt.Println(" ----=====----")
			return
		}
		row_ip, err := db.Query("SELECT address_ip FROM modules_sfp_sla WHERE id=?", id_sfp1)
		if err != nil {
			db.Exec("UPDATE test_sla_real SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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

				db.Exec("UPDATE test_sla_real SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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
			db.Exec("UPDATE test_sla_real SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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

				db.Exec("UPDATE test_sla_real SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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
			db.Exec("UPDATE test_sla_real SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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

				db.Exec("UPDATE test_sla_real SET status=? WHERE id=?", 4, id) // Ошибка выполнения
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
	var test_this testSLA

	var test_type uint16
	test_type = 0x2000 + (uint16(id) & 0x1FFF)

	c, err := raw.ListenPacket(ifi, etherType, nil)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		db.Close()

		fmt.Println(" -!! Error !!-")
		fmt.Println(err)
		fmt.Println(" ----=====----")
		return
	}

	addr := &raw.Addr{
		HardwareAddr: ethernet.Broadcast,
	}
	ipsrc := net.ParseIP(ipsrcstr)
	ipdst1 := net.ParseIP(ipdst_1sfpsla_str)
	ipdst2 := net.ParseIP(ipdst_2sfpsla_str)

	var number uint32

	number = 1
	//b := packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, test.block_size, number)
	t := time.NewTicker(time.Duration(test.clock) * time.Millisecond)
	defer t.Stop()
	var circ bool
	if test.count == 0 {
		circ = true
	} else {
		circ = false
	}
	check_count := 10
	counter := test.count
	for range t.C {
		fmt.Println(time.Now())
		if !circ {
			counter--
			if counter <= 0 {
				break
			}
		}

		number++
		b := packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, test.block_size, number, test_type, test.test_type)

		go func() {
			c.WriteTo(b, addr)
		}()

		go test_this.receiveMessages(id, c, ipdst_1sfpsla_str, test.node_zabbix, host_zabbix, port_zabbix, ifi.MTU, *test, test_type, testMax, len(b))
		//time.Sleep(time.Duration(test.clock/2) * time.Millisecond)
		check_count--
		if check_count < 0 {
			db, err = sql.Open("mysql", db_user+":"+db_user_pass+"@/"+db_database)
			if err == nil {
				rez_f, er := db.Query("SELECT id FROM test_sla_real WHERE id=?", id)
				if er != nil {
					rez_f.Close()
					break
				}

				if !rez_f.Next() {
					rez_f.Close()
					break
				}
				rez_f.Close()
			}
			db.Close()
			check_count = 10
		}
	}

}

func packetForm(ipsrc net.IP, ipdst1 net.IP, ipdst2 net.IP, mac_src []byte, mac_dst []byte, size int, number uint32, test_type uint16, testWay int) []byte {
	ip := iphdr{
		vhl:   0x45,
		tos:   0,
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
		t_time := int64(float64(time.Now().UnixNano())*float64(math.Pow(2, 32)/1000000000)) - 0x55817800000000
		//	t_time := int64(float64(time.Now().UnixNano())*float64(math.Pow(2, 32)/1000000000))
		//t_time = t_time << (4*8)
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
		log.Fatalf("failed to marshal ethernet frame: %v", err)

		fmt.Println(" -!! Error !!-")
		fmt.Println(err)
		fmt.Println(" ----=====----")
		return []byte{}
	}
	return b
}

//var min_period = time.Duration(15) * time.Microsecond

func (test *testThr) testThrGen(b []byte, c *raw.Conn, addr *raw.Addr, mtu int, ipdst_1sfpsla_str string, cnt int64, period_nano int64, t_type uint16) (int, int64) {
	time_to_gen := ((cnt * period_nano) * 150) / 100
	time_gen := time.Duration(cnt * period_nano)
	fmt.Println("		-- period_to_generate [ms] = ", (cnt*period_nano)/1000000)
	fmt.Println("		-- time_to_gen [ms]        = ", time_to_gen/1000000)
	fmt.Println("		-- cnt start= ", cnt)

	test_count := 10000
	gen_test_pps_start := time.Now()
	for {
		test_count--
		c.WriteTo(b, addr)
		if test_count <= 0 {
			break
		}
	}
	pps_rez := 10000 * 1000000000 / int(time.Since(gen_test_pps_start))
	fmt.Println("		 -*- max pps = ", pps_rez)

	test_count = 1000
	gen_test_min_period_start := time.Now()
	ticker := time.NewTicker(time.Duration(1))
	for range ticker.C {
		test_count--
		c.WriteTo(b, addr)
		if test_count <= 0 {
			break
		}
	}
	min_per_rez := int(time.Since(gen_test_min_period_start)) / (1000 * 1000)
	fmt.Println("		 -*- min period [mks] = ", min_per_rez)

	gen_start := time.Now()
	var rez_time int64
	test.numberCounter = 0
	//for i := 0; i < 32; i++ {

	/*

		go func() {
			//ticker := time.NewTicker(time.Duration(period_nano))
			//timer := time.NewTimer(time.Microsecond * 10)
			//period := time.Duration(period_nano)
			//fmt.Println("Start")
			//for range ticker.C {
			for {
				//	select {
				//	case <-ticker.C:
				//time.Sleep(period)
				cnt--
				c.WriteTo(b, addr)
				if cnt <= 0 {
					break
				}
				if cnt%10000 == 0 {
					if time.Since(gen_start) >= time_gen {
						break
					}
					//		}
				}
			}
			rez_time = (int64)(time.Since(gen_start))
			//ticker.Stop()
			fmt.Println("		 -*- rez_time = ", rez_time)
		}()
		//}
	*/
	quit := make(chan int)
	blen := len(b)
	go func() {
	ExitLoop:
		for {
			select {
			case <-quit:
				break ExitLoop
			default:
				n, err := c.WriteTo(b, addr)
				if err != nil {
					fmt.Printf("%v", err)
					continue
				}
				if n < blen {
					fmt.Printf("Partial write: %d", n)
					continue
				}
				cnt--
				if cnt <= 0 {
					break ExitLoop
				}
				if cnt%10000 == 0 {
					if time.Since(gen_start) >= time_gen {
						break ExitLoop
					}					
				}
			}
		}
	}()
	go (*test).receivePackets(c, mtu, ipdst_1sfpsla_str, quit, t_type)
	time.Sleep(time.Duration(time_to_gen))

	rez_count := test.numberCounter
	fmt.Println("		 --->> rez_count= ", rez_count)
	quit <- 1
	return int(rez_count), rez_time
}

func (test *testThr) testMax(b []byte, c *raw.Conn, addr *raw.Addr, mtu int, ipdst_1sfpsla_str string, cnt int, t_type uint16) (int, int64) {
	var min_period = time.Duration(15) * time.Microsecond
	time_to_Sleep := time.Duration(cnt) * min_period * 2
	fmt.Println("time_to_Sleep= ", time_to_Sleep)
	fmt.Println("cnt start= ", cnt)
	gen_start := time.Now()
	var rez_time int64
	test.numberCounter = 0
	//for i := 0; i < 32; i++ {
	go func() {
		//timer := time.NewTimer(time.Nanosecond * 10)
		//timer := time.NewTimer(time.Microsecond * 10)
		//fmt.Println("Start")
		//for range timer.C {
		for {
			cnt--
			c.WriteTo(b, addr)
			if cnt < 0 {
				break
			}
		}
		rez_time = (int64)(time.Since(gen_start))
		fmt.Println("rez_time = ", rez_time)
	}()
	//}

	quit := make(chan int)
	go (*test).receivePackets(c, mtu, ipdst_1sfpsla_str, quit, t_type)
	time.Sleep(time_to_Sleep)

	rez_count := test.numberCounter
	fmt.Println("rez_count= ", rez_count)
	quit <- 1
	return int(rez_count), rez_time
}

func (test *testThr) sendPackets(c net.PacketConn, source net.HardwareAddr, dist net.HardwareAddr, ipsrc net.IP, ipdst1 net.IP, ipdst2 net.IP, numberTX uint32, size uint16) {

	//	ipsrc := net.ParseIP(ipsrcstr)
	//	ipdst1 := net.ParseIP(ipdst_1sfpsla_str)
	//	ipdst2 := net.ParseIP(ipdst_2sfpsla_str)

	// Default message to system's hostname if empty.
	ip := iphdr{
		vhl:   0x45,
		tos:   0,
		id:    0x0000, // the kernel overwrites id if it is zero
		off:   0,
		ttl:   0xFF,
		proto: 0x5E,
	}
	copy(ip.src[:], ipsrc.To4())
	copy(ip.dst[:], ipdst1.To4())
	sfpdat := sfpsla{
		id: 0xFC,
	}
	copy(sfpdat.dst[:], ipdst2.To4())
	(*test).numberTx++
	sfpdat.number = test.numberTx
	//ip.iplen = uint16(20 + 26 + 4)
	ip.iplen = uint16(size)
	ip.checksum()
	payloadAdd := make([]byte, size-64)
	var bin_buf bytes.Buffer
	binary.Write(&bin_buf, binary.BigEndian, ip)
	binary.Write(&bin_buf, binary.BigEndian, sfpdat)
	binary.Write(&bin_buf, binary.BigEndian, payloadAdd)

	msg := bin_buf.Bytes()
	f := &ethernet.Frame{
		//Destination: ethernet.Broadcast,
		//Destination: []byte{0x5A, 0x11, 0x22, 0x33, 0x44, 0x00},
		Destination: dist,
		//Destination: []byte{0x64, 0xD1, 0x54, 0x17, 0xF6, 0x82},
		Source:    source,
		EtherType: 0x0800,
		Payload:   []byte(msg),
	}

	b, err := f.MarshalBinary()
	if err != nil {
		log.Fatalf("failed to marshal ethernet frame: %v", err)
	}

	// Required by Linux, even though the Ethernet frame has a destination.
	// Unused by BSD.
	addr := &raw.Addr{
		HardwareAddr: ethernet.Broadcast,
	}
	/*
		fmt.Printf("raw:  %x \n", b)
		fmt.Println(" --== Packet send ==--")
		fmt.Printf("mac dst  %x \n", b[0:6])
		fmt.Printf("mac src  %x \n", b[6:12])
		fmt.Printf("type eth %x \n", b[12:14])
		fmt.Printf("size     %v \n", b[16:18])

		fmt.Printf("ip sourse %v.%v.%v.%v \n", b[26], b[27], b[28], b[29])
		fmt.Printf("ip dst    %v.%v.%v.%v \n", b[30], b[31], b[32], b[33])
		fmt.Println(" --== End Packet ==--")
	*/

	//  t := time.NewTicker(

	for i := 0; i < 1000; i++ {
		if _, err := c.WriteTo(b, addr); err != nil {
			log.Fatalf("failed to send message: %v", err)
			fmt.Println("failed to send message: ", err)
		}
	}
}

//func receivePackets(c net.PacketConn, mtu int, ipdst_1sfpsla_str string, counter int) {
func (test *testThr) receivePackets(c net.PacketConn, mtu int, ipdst_1sfpsla_str string, quit chan int, t_type uint16) { //, counter chan<- int) {
	var f ethernet.Frame
	b := make([]byte, mtu)
	//var count int
	// Keep receiving messages forever.
	for {
		select {
		case <-quit:
			return
		default:
		}
		n, _, err := c.ReadFrom(b)
		if err != nil {
			log.Fatalf("failed to receive message: %v", err)
		}

		// Unpack Ethernet II frame into Go representation.
		if err := (&f).UnmarshalBinary(b[:n]); err != nil {
			log.Fatalf("failed to unmarshal ethernet frame: %v", err)
		}
		//fmt.Printf("\n\n--=Test %x - \n", f.Payload[12:16])
		var ips [4]byte
		copy(ips[:], (net.ParseIP(ipdst_1sfpsla_str)).To4())

		var t_ips [2]byte
		t_ips[1] = byte(t_type & 0xFF)
		t_ips[0] = byte((t_type >> 8) & 0xFF)
		//fmt.Printf("\n\n--=T_so %x - \n", ips)

		// Display source of message and message itself.
		if (f.Payload[20] == 0xFC) && (bytes.Equal(f.Payload[12:16], ips[:]) == true) && (bytes.Equal(f.Payload[50:52], t_ips[:]) == true) {
			//count++
			//	fmt.Printf("-->>Detect")
			//	counter <-count
			//(*test).numberCounter = uint32(count)
			(*test).numberCounter++
		}
	}
}

/*
// sendMessages continuously sends a message over a connection at regular intervals,
// sourced from specified hardware address.
func sendMessages(c net.PacketConn, source net.HardwareAddr) {

	t := time.NewTicker(1 * time.Microsecond)
	for range t.C {

		ipsrc := net.ParseIP(ipsrcstr)
		ipdst1 := net.ParseIP(ipdst_1sfpsla_str)
		ipdst2 := net.ParseIP(ipdst_2sfpsla_str)

		// Default message to system's hostname if empty.
		ip := iphdr{
			vhl:   0x45,
			tos:   0,
			id:    0x0000, // the kernel overwrites id if it is zero
			off:   0,
			ttl:   0xFF,
			proto: 0x5E,
		}
		copy(ip.src[:], ipsrc.To4())
		copy(ip.dst[:], ipdst1.To4())
		sfpdat := sfpsla{
			id: 0xFC,
		}
		copy(sfpdat.dst[:], ipdst2.To4())
		numberTx++
		sfpdat.number = numberTx
		ip.iplen = uint16(20 + 26 + 4)
		ip.checksum()

		var bin_buf bytes.Buffer
		binary.Write(&bin_buf, binary.BigEndian, ip)
		binary.Write(&bin_buf, binary.BigEndian, sfpdat)

		msg := bin_buf.Bytes()
		f := &ethernet.Frame{
			//Destination: ethernet.Broadcast,
			Destination: []byte{0x5A, 0x11, 0x22, 0x33, 0x44, 0x00},
			//Destination: []byte{0x64, 0xD1, 0x54, 0x17, 0xF6, 0x82},
			Source:    source,
			EtherType: 0x0800,
			Payload:   []byte(msg),
		}

		b, err := f.MarshalBinary()
		if err != nil {
			log.Fatalf("failed to marshal ethernet frame: %v", err)
		}

		// Required by Linux, even though the Ethernet frame has a destination.
		// Unused by BSD.
		addr := &raw.Addr{
			HardwareAddr: ethernet.Broadcast,
		}

		fmt.Printf("raw:  %x \n", b)
		fmt.Println(" --== Packet send ==--")
		fmt.Printf("mac dst  %x \n", b[0:6])
		fmt.Printf("mac src  %x \n", b[6:12])
		fmt.Printf("type eth %x \n", b[12:14])
		fmt.Printf("size     %v \n", b[16:18])

		fmt.Printf("ip sourse %v.%v.%v.%v \n", b[26], b[27], b[28], b[29])
		fmt.Printf("ip dst    %v.%v.%v.%v \n", b[30], b[31], b[32], b[33])
		fmt.Println(" --== End Packet ==--")

		if _, err := c.WriteTo(b, addr); err != nil {
			log.Fatalf("failed to send message: %v", err)
		}
	}
}

// receiveMessages continuously receives messages over a connection. The messages
// may be up to the interface's MTU in size.
*/
func (test *testSLA) receiveMessages(id int, c net.PacketConn, ipdst_1sfpsla_str string, node_zabbix string, host_zabbix string, port_zabbix int, mtu int, test_id testReal, t_type uint16, tMax testRealMax, packetSize int) {
	var f ethernet.Frame
	b := make([]byte, mtu)

	var t_ips [2]byte
	t_ips[1] = byte(t_type & 0xFF)
	t_ips[0] = byte((t_type >> 8) & 0xFF)

	start := time.Now()
	// Keep receiving messages forever.
	for {
		n, _, err := c.ReadFrom(b)
		if err != nil {
			fmt.Printf("failed to receive message: %v", err)
			log.Fatalf("failed to receive message: %v", err)

		}

		if time.Since(start) > (time.Millisecond * 500) {
			break
		}

		if (n) != packetSize {
			continue
		}

		t_time := int64(float64(time.Now().UnixNano())*float64(math.Pow(2, 32)/1000000000)) - 0x55817800000000
		t_time = t_time & int64(0xFFFFFFFFFFFFFF)

		//n, addr, err := c.ReadFrom(b)
		// Unpack Ethernet II frame into Go representation.
		if err := (&f).UnmarshalBinary(b[:n]); err != nil {
			fmt.Printf("failed to unmarshal ethernet frame: %v", err)
			log.Fatalf("failed to unmarshal ethernet frame: %v", err)
		}

		//fmt.Println("\n\n--=Test ==-- - ")
		//fmt.Printf("\n\rEthernet source: [%s]\n", addr.String())
		//fmt.Printf("ip sourse %v.%v.%v.%v \n", f.Payload[12], f.Payload[13], f.Payload[14], f.Payload[15])
		//fmt.Printf("ip dst    %v.%v.%v.%v \n", f.Payload[16], f.Payload[17], f.Payload[18], f.Payload[19])
		var ips [4]byte
		copy(ips[:], (net.ParseIP(ipdst_1sfpsla_str)).To4())
		//fmt.Printf("\n--=T_so ip dst    %v.%v.%v.%v \n", ips[0], ips[1], ips[2], ips[3])
		//fmt.Printf("ip SFP2   %v.%v.%v.%v \n", f.Payload[21], f.Payload[22], f.Payload[23], f.Payload[24])

		//fmt.Printf("time marker_SFP1_1 :   %x \n", f.Payload[25:32])
		//fmt.Printf("time marker_SFP2   :   %x \n", f.Payload[32:39])
		//fmt.Printf("time marker_SFP1_2 :   %x \n", f.Payload[39:46])

		//fmt.Println(" --== End Test ==--")

		// Display source of message and message itself.
		if (len(f.Payload) >= 52) && (f.Payload[20] == 0xFC) && (bytes.Equal(f.Payload[12:16], ips[:]) == true) && (bytes.Equal(f.Payload[50:52], t_ips[:]) == true) {

			(*test).number++
			/*
				fmt.Printf("\n\n--=Packet DETECT!!!=--\n")
				//fmt.Printf("\n\n--=Test %x - \n -== %x\n",f.Payload[12:15],net.ParseIP(ipdst_1sfpsla_str))
				fmt.Printf("size: %v raw:  %x \n", len(f.Payload), f.Payload)
				fmt.Printf("\n\rEthernet source: [%s]\n", addr.String())

				fmt.Printf("size     %x \n", b[2:4])

				fmt.Printf("ip sourse %v.%v.%v.%v \n", f.Payload[12], f.Payload[13], f.Payload[14], f.Payload[15])
				fmt.Printf("ip dst    %v.%v.%v.%v \n", f.Payload[16], f.Payload[17], f.Payload[18], f.Payload[19])

				fmt.Printf("ip SFP2   %v.%v.%v.%v \n", f.Payload[21], f.Payload[22], f.Payload[23], f.Payload[24])

				fmt.Printf("time marker_SFP1_1 :   %x \n", f.Payload[25:32])
				fmt.Printf("time marker_SFP2   :   %x \n", f.Payload[32:39])
				fmt.Printf("time marker_SFP1_2 :   %x \n", f.Payload[39:46])
				fmt.Printf("Number marker      :   %x \n", f.Payload[46:50])
				fmt.Println(" --== End Packet ==--")
			*/
			var markerSFP11, markerSFP12, markerSFP2 int64
			var ind uint

			for ind = 0; ind < 7; ind++ {
				markerSFP11 = markerSFP11 + int64(f.Payload[31-ind])<<(8*ind)
				markerSFP2 = markerSFP2 + int64(f.Payload[38-ind])<<(8*ind)
				markerSFP12 = markerSFP12 + int64(f.Payload[45-ind])<<(8*ind)
			}

			var numberR uint32
			for ind = 0; ind < 4; ind++ {
				numberR += uint32(f.Payload[49-ind]) << (8 * ind)
			}

			var delay, delay1, delay2, jitter, jitter1, jitter2, loss float32

			if test_id.test_delay == true {
				if test_id.test_type == 1 {
					delay = zabbix_delay(node_zabbix, markerSFP12-markerSFP11, host_zabbix, port_zabbix)
					if test_id.test_delay_jitter == true {
						jitter = zabbix_jitter(node_zabbix, (*test).getJitter(markerSFP12-markerSFP11), host_zabbix, port_zabbix)
					}
				}
				if test_id.test_type == 2 {
					//t_time := int64(float32(time.Now().Nanosecond())) //* 1000000 / float32(math.Pow(2, 32)))

					delay = zabbix_delay(node_zabbix, t_time-markerSFP2, host_zabbix, port_zabbix)
					if test_id.test_delay_jitter == true {
						jitter = zabbix_jitter(node_zabbix, (*test).getJitter(t_time-markerSFP2), host_zabbix, port_zabbix)
					}
				}
			}
			fmt.Println("==>> number_pack - ", numberR)
			fmt.Println("==>> number_test - ", test.number)
			if test_id.test_loss == true {
				loss = zabbix_error(node_zabbix, float32(numberR-test.number)/float32(numberR), host_zabbix, port_zabbix)
			}
			//*
			if test_id.test_type == 1 {
				if test_id.test_delay_1 == true {
					rez_delay_to, rez_delay_un := (*test).getOneDelay(markerSFP2-markerSFP11, markerSFP12-markerSFP2)
					delay1 = zabbix_delay_to(node_zabbix, rez_delay_to, host_zabbix, port_zabbix)
					delay2 = zabbix_delay_un(node_zabbix, rez_delay_un, host_zabbix, port_zabbix)
					if test_id.test_delay_1_jitter == true {
						jitter1 = zabbix_jitter_to(node_zabbix, (*test).getJitterto(rez_delay_to), host_zabbix, port_zabbix)
						jitter2 = zabbix_jitter_un(node_zabbix, (*test).getJitterun(rez_delay_un), host_zabbix, port_zabbix)

					}
				}
			}
			if test_id.test_type == 2 {
				//t_time := int64(time.Now().Nanosecond())
				//t_time := int64(float64(time.Now().UnixNano())*float64(math.Pow(2, 32)/1000000000)) - 0x55817800000000
				if test_id.test_delay_1 == true {
					rez_delay_to, rez_delay_un := (*test).getOneDelay(markerSFP12-markerSFP2, markerSFP12-t_time)
					delay1 = zabbix_delay_to(node_zabbix, rez_delay_to, host_zabbix, port_zabbix)
					delay2 = zabbix_delay_un(node_zabbix, rez_delay_un, host_zabbix, port_zabbix)
					if test_id.test_delay_1_jitter == true {
						jitter1 = zabbix_jitter_to(node_zabbix, (*test).getJitterto(rez_delay_to), host_zabbix, port_zabbix)
						jitter2 = zabbix_jitter_un(node_zabbix, (*test).getJitterun(rez_delay_un), host_zabbix, port_zabbix)

					}
				}
			}
			//*/
			db, err := sql.Open("mysql", db_user+":"+db_user_pass+"@/"+db_database)
			if err != nil {
				db.Close()
				fmt.Println(" -!! Error !!-")
				fmt.Println(err)
				fmt.Println(" ----=====----")
				return
			}
			var dt = time.Now()
			//dt.Format(time.RFC3339)
			_, err = db.Exec("INSERT INTO test_sla_real_rez (datetime, test_id, delay_rez, delay_to_rez, delay_un_rez, jitter_delay_rez, jitter_delay_to, jitter_delay_un, packet_loss) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)", dt, id, delay, delay1, delay2, jitter, jitter1, jitter2, loss)
			if err != nil {
				db.Close()
				fmt.Println(" -!! Error !!-")
				fmt.Println(err)
				fmt.Println(" ----=====----")
				return
			}
			if (tMax.delayMax != 0) && (tMax.delayMax < delay) {
				msg := fmt.Sprintf("!! Превышение порогового значения времени двусторонней задержки на %.4f мкс", delay-tMax.delayMax)
				db.Exec("INSERT INTO test_sla_real_alarm (id_test, datatime, message) VALUES(?, ?, ?)", id, dt, msg)
			}
			if (tMax.jitterMax != 0) && (tMax.jitterMax < float32(math.Abs(float64(jitter)))) {
				msg := fmt.Sprintf("!! Превышение порогового значения джиттера времени двусторонней задержки на %.4f мкс", float32(math.Abs(float64(jitter)))-tMax.jitterMax)
				db.Exec("INSERT INTO test_sla_real_alarm (id_test, datatime, message) VALUES(?, ?, ?)", id, dt, msg)
			}
			if (tMax.lossMax != 0) && (tMax.lossMax < loss) {
				msg := fmt.Sprintf("!! Превышение порогового значения вероятности ошибки на %.6f", loss-tMax.lossMax)
				db.Exec("INSERT INTO test_sla_real_alarm (id_test, datatime, message) VALUES(?, ?, ?)", id, dt, msg)

			}

			if (tMax.delayOneMax != 0) && ((tMax.delayOneMax < delay1) || (tMax.delayOneMax < delay2)) {
				msg := fmt.Sprintf("!! Превышение порогового значения времени односторонней задержки на %.4f мкс", float32(math.Max(float64(delay1), float64(delay2)))-tMax.delayOneMax)
				db.Exec("INSERT INTO test_sla_real_alarm (id_test, datatime, message) VALUES(?, ?, ?)", id, dt, msg)
			}
			if (tMax.jitterOneMax != 0) && (tMax.jitterOneMax < float32(math.Abs(float64(jitter1)))) {
				msg := fmt.Sprintf("!! Превышение порогового значения джиттера времени односторонней задержки на %.4f мкс", float32(math.Abs(float64(jitter1)))-tMax.jitterOneMax)
				db.Exec("INSERT INTO test_sla_real_alarm (id_test, datatime, message) VALUES(?, ?, ?)", id, dt, msg)
			}
			if (tMax.jitterOneMax != 0) && (tMax.jitterOneMax < float32(math.Abs(float64(jitter2)))) {
				msg := fmt.Sprintf("!! Превышение порогового значения джиттера времени односторонней задержки на %.4f мкс", float32(math.Abs(float64(jitter2)))-tMax.jitterOneMax)
				db.Exec("INSERT INTO test_sla_real_alarm (id_test, datatime, message) VALUES(?, ?, ?)", id, dt, msg)
			}

			db.Close()
			fmt.Printf("  ==>> %v  --> %s\n", id, time.Since(start))
			break
		} else {
			//fmt.Printf("\n\n\r[%s] %v %x", addr.String(), len(f.Payload), f.Payload[:25])

		}
	}
}

//var mass_solve []int64
func (test *testSLA) getOneDelay(in_delay_to int64, in_delay_un int64) (int64, int64) {

	if len((*test).delay_solve_to) < 2 {

		(*test).delay_solve_to = append((*test).delay_solve_to, in_delay_to)
		(*test).delay_solve_to = append((*test).delay_solve_to, 1)

		(*test).delay_solve_un = append((*test).delay_solve_un, in_delay_un)
		(*test).delay_solve_un = append((*test).delay_solve_un, 1)

		return in_delay_to, in_delay_un

	}

	mean_to := float32((*test).delay_solve_to[0]) * float32((*test).delay_solve_to[1])
	(*test).delay_solve_to[1] = (*test).delay_solve_to[1] + 1
	mean_to = (mean_to + float32(in_delay_to)) / float32((*test).delay_solve_to[1])
	(*test).delay_solve_to[0] = int64(mean_to)

	mean_un := float32((*test).delay_solve_un[0]) * float32((*test).delay_solve_un[1])
	(*test).delay_solve_un[1] = (*test).delay_solve_un[1] + 1
	mean_un = (mean_un + float32(in_delay_un)) / float32((*test).delay_solve_un[1])
	(*test).delay_solve_un[0] = int64(mean_un)

	/*
		fmt.Printf(" --== Jitter debug ==-- \n")
		fmt.Printf(" --== Slice: %x \n", (*test).delay_solve)
		fmt.Printf(" --== Max = %x \n", max)
		fmt.Printf(" --== Min = %x \n", min)
		fmt.Printf(" --== Mean = %f \n", mean)
		fmt.Printf(" --== Jitter = %f \n", jitter)
		fmt.Printf(" --== End Jitter debug ==-- \n")
	*/
	return int64(mean_to), int64(mean_un)
}

func (test *testSLA) getJitter(in_solve int64) float32 {
	var jitter float32
	if len((*test).delay_solve) == 0 {
		(*test).delay_solve = append((*test).delay_solve, in_solve)
	}
	jitter = float32((*test).delay_solve[0] - in_solve)
	(*test).delay_solve[0] = in_solve
	/*
		var jitter, mean float32
		var size_s int
		var max, min int64

		size_s = 100
		(*test).delay_solve = append((*test).delay_solve, in_solve)
		if len((*test).delay_solve) < (size_s + 1) {
			return 0
		}
		test.delay_solve = (*test).delay_solve[1:(size_s + 1)]
		max = (*test).delay_solve[0]
		min = max
		mean = float32((*test).delay_solve[0]) / float32(size_s)

		for ind := 1; ind < size_s; ind++ {
			if max < (*test).delay_solve[ind] {
				max = (*test).delay_solve[ind]
			}
			if min > (*test).delay_solve[ind] {
				min = (*test).delay_solve[ind]
			}
			mean = mean + (float32((*test).delay_solve[ind]) / float32(size_s))
		}
		if (float32(max) - mean) > (mean - float32(min)) {
			jitter = float32(max) - mean
		} else {
			jitter = mean - float32(min)
		}
		/*
			fmt.Printf(" --== Jitter debug ==-- \n")
			fmt.Printf(" --== Slice: %x \n", (*test).delay_solve)
			fmt.Printf(" --== Max = %x \n", max)
			fmt.Printf(" --== Min = %x \n", min)
			fmt.Printf(" --== Mean = %f \n", mean)
			fmt.Printf(" --== Jitter = %f \n", jitter)
			fmt.Printf(" --== End Jitter debug ==-- \n")
	*/
	return jitter
}

func (test *testSLA) getJitterto(in_solve int64) float32 {
	var jitter float32
	if len((*test).delay_solve_to) < 3 {
		(*test).delay_solve_to = append((*test).delay_solve_to, in_solve)
	}
	jitter = float32((*test).delay_solve_to[2] - in_solve)
	(*test).delay_solve_to[2] = in_solve
	/*
		var jitter, mean float32
		var size_s int
		var max, min int64

		size_s = 100
		(*test).delay_solve_to = append((*test).delay_solve_to, in_solve)
		if len((*test).delay_solve_to) < (size_s + 1) {
			return 0
		}
		test.delay_solve_to = (*test).delay_solve_to[1:(size_s + 1)]

		max = (*test).delay_solve_to[0]
		min = max
		mean = float32((*test).delay_solve_to[0]) / float32(size_s)

		for ind := 1; ind < size_s; ind++ {
			if max < (*test).delay_solve_to[ind] {
				max = (*test).delay_solve_to[ind]
			}
			if min > (*test).delay_solve_to[ind] {
				min = (*test).delay_solve_to[ind]
			}
			mean = mean + (float32((*test).delay_solve_to[ind]) / float32(size_s))
		}
		if (float32(max) - mean) > (mean - float32(min)) {
			jitter = float32(max) - mean
		} else {
			jitter = mean - float32(min)
		}*/
	return jitter
}

func (test *testSLA) getJitterun(in_solve int64) float32 {
	var jitter float32
	if len((*test).delay_solve_un) < 3 {
		(*test).delay_solve_un = append((*test).delay_solve_un, in_solve)
	}
	jitter = float32((*test).delay_solve_un[2] - in_solve)
	(*test).delay_solve_un[2] = in_solve
	/*
		var jitter, mean float32
		var size_s int
		var max, min int64

		size_s = 100
		(*test).delay_solve_un = append((*test).delay_solve_un, in_solve)
		if len((*test).delay_solve_un) < (size_s + 1) {
			return 0
		}
		test.delay_solve_un = (*test).delay_solve_un[1:(size_s + 1)]

		max = (*test).delay_solve_un[0]
		min = max
		mean = float32((*test).delay_solve_un[0]) / float32(size_s)

		for ind := 1; ind < size_s; ind++ {
			if max < (*test).delay_solve_un[ind] {
				max = (*test).delay_solve_un[ind]
			}
			if min > (*test).delay_solve_un[ind] {
				min = (*test).delay_solve_un[ind]
			}
			mean = mean + (float32((*test).delay_solve_un[ind]) / float32(size_s))
		}
		if (float32(max) - mean) > (mean - float32(min)) {
			jitter = float32(max) - mean
		} else {
			jitter = mean - float32(min)
		}*/
	return jitter
}

func zabbixHello(host string, defaultHost string, defaultPort int) {
	var delay int
	for {
		delay = rand.Intn(1500)
		//delay:=i*100
		var metrics []*Metric
		metrics = append(metrics, NewMetric(host, "delay", fmt.Sprint(delay), time.Now().Unix()))

		// Create instance of Packet class
		packet := NewPacket(metrics)
		//fmt.Println(packet);
		// Send packet to zabbix
		z := NewSender(defaultHost, defaultPort)
		z.Send(packet)
		time.Sleep(5 * time.Second)
	}
}

func zabbix_delay(host string, delay int64, defaultHost string, defaultPort int) float32 {

	//delay = delay * 8 // [mks] 125 MGz - clock, => T = 8 mks
	//delay = int64( float64(delay) * 1000000 / (math.Pow(2, 32))) // [mks] 125 MGz - clock, => T = 8 mks
	var delfloat float32
	delfloat = float32(delay) * 1000000 / float32(math.Pow(2, 32))
	var metrics []*Metric
	metrics = append(metrics, NewMetric(host, "delay", fmt.Sprint(delfloat), time.Now().Unix()))

	// Create instance of Packet class
	packet := NewPacket(metrics)
	//fmt.Println(packet);
	// Send packet to zabbix
	z := NewSender(defaultHost, defaultPort)
	z.Send(packet)
	return delfloat
}

func zabbix_delay_to(host string, delay int64, defaultHost string, defaultPort int) float32 {

	//delay = delay * 8 // [mks] 125 MGz - clock, => T = 8 mks
	//delay = int64( float64(delay) * 1000000 / (math.Pow(2, 32))) // [mks] 125 MGz - clock, => T = 8 mks
	var delfloat float32
	delfloat = float32(delay) * 1000000 / float32(math.Pow(2, 32))
	var metrics []*Metric
	metrics = append(metrics, NewMetric(host, "delay_SFP1_SFP2", fmt.Sprint(delfloat), time.Now().Unix()))

	// Create instance of Packet class
	packet := NewPacket(metrics)
	//fmt.Println(packet);
	// Send packet to zabbix
	z := NewSender(defaultHost, defaultPort)
	z.Send(packet)
	return delfloat
}

func zabbix_delay_un(host string, delay int64, defaultHost string, defaultPort int) float32 {

	//delay = delay * 8 // [mks] 125 MGz - clock, => T = 8 mks
	//delay = int64( float64(delay) * 1000000 / (math.Pow(2, 32))) // [mks] 125 MGz - clock, => T = 8 mks
	var delfloat float32
	delfloat = float32(delay) * 1000000 / float32(math.Pow(2, 32))
	var metrics []*Metric
	metrics = append(metrics, NewMetric(host, "delay_SFP2_SFP1", fmt.Sprint(delfloat), time.Now().Unix()))

	// Create instance of Packet class
	packet := NewPacket(metrics)
	//fmt.Println(packet);
	// Send packet to zabbix
	z := NewSender(defaultHost, defaultPort)
	z.Send(packet)
	return delfloat
}

func zabbix_jitter(host string, jitter float32, defaultHost string, defaultPort int) float32 {

	//delay = delay * 8 // [mks] 125 MGz - clock, => T = 8 mks
	if jitter != 0 {
		jitter = jitter * 1000000 / float32(math.Pow(2, 32)) // [mks] 125 MGz - clock, => T = 8 mks
	}
	var metrics []*Metric
	metrics = append(metrics, NewMetric(host, "jitter", fmt.Sprint(math.Abs(float64(jitter))), time.Now().Unix()))

	// Create instance of Packet class
	packet := NewPacket(metrics)
	//fmt.Println(packet);
	// Send packet to zabbix
	z := NewSender(defaultHost, defaultPort)
	z.Send(packet)
	return jitter
}

func zabbix_jitter_to(host string, jitter float32, defaultHost string, defaultPort int) float32 {

	if jitter != 0 {
		jitter = jitter * 1000000 / float32(math.Pow(2, 32)) // [mks] 125 MGz - clock, => T = 8 mks
	}
	var metrics []*Metric
	metrics = append(metrics, NewMetric(host, "jitter_delay_SFP1_SFP2", fmt.Sprint(math.Abs(float64(jitter))), time.Now().Unix()))
	packet := NewPacket(metrics)
	z := NewSender(defaultHost, defaultPort)
	z.Send(packet)
	return jitter
}

func zabbix_jitter_un(host string, jitter float32, defaultHost string, defaultPort int) float32 {

	if jitter != 0 {
		jitter = jitter * 1000000 / float32(math.Pow(2, 32)) // [mks] 125 MGz - clock, => T = 8 mks
	}
	var metrics []*Metric
	metrics = append(metrics, NewMetric(host, "jitter_delay_SFP2_SFP1", fmt.Sprint(math.Abs(float64(jitter))), time.Now().Unix()))
	packet := NewPacket(metrics)
	z := NewSender(defaultHost, defaultPort)
	z.Send(packet)
	return jitter
}

func zabbix_error(host string, err float32, defaultHost string, defaultPort int) float32 {

	var metrics []*Metric
	metrics = append(metrics, NewMetric(host, "error_probability", fmt.Sprint(err), time.Now().Unix()))

	// Create instance of Packet class
	packet := NewPacket(metrics)
	//fmt.Println(packet);
	// Send packet to zabbix
	z := NewSender(defaultHost, defaultPort)
	z.Send(packet)
	return err
}

/*import (
	"encoding/json"
	"log"
	"net/http"

)

func main() {
	http.HandleFunc("/", mainPage)
	port := ":313"
	println("Server Listen on port", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	user := User{"Alex", "Net"}
	js, _ := json.Marshal(user)
	w.Write(js)
}
*/
