package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/raw"

	"golang.org/x/net/bpf"
)

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
		//	log.Fatalf("failed to find interface %q: %v", net_interface_name, err)
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
	test_type = 0x2000 + (uint16(id) & 0x1FFF)

	var test_c testThr
	test_c.numberCounter = uint32(test.count)
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

	c, err := raw.ListenPacket(ifi, etherType, netConf)

	if err != nil {
		//	log.Fatalf("failed to listen: %v", err)
		db.Close()

		fmt.Println(" -!! Error !!-")
		fmt.Println(err)
		fmt.Println(" ----=====----")
		return
	}

	period_test := test.count // период теста - 10 секунд
	size := 64
	fmt.Println("->> test.thr_begin = ", test.thr_begin)
	period_nano := int64(size*8*1000000000) / (int64(test.thr_begin * 1000 * 1000))
	packet_count := (int64(period_test * 1000000000)) / period_nano

	fmt.Println("->> period_nano  = ", period_nano)
	fmt.Println("->> packet_count = ", packet_count)

	b := packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, 64, number, test_type, test.test_type)
	//count_rez, per := test_c.testMax(b, c, addr, ifi.MTU, ipdst_1sfpsla_str, counter, test_type)
	count_rez, per := test_c.testThrGen(net_interface_name, b, c, addr, ifi.MTU, ipdst_1sfpsla_str, packet_count, period_nano, test_type)
	test.rez_64 = (float32)(64.0 * 8.0 * (float32)(count_rez) * 1000000 / (float32)(per))
	fmt.Println("->> rez_64 = ", count_rez, " period = ", per, " !!!  rez=", test.rez_64)

	size = 128
	period_nano = int64(size * 8 * 1000000000 / (test.thr_begin * 1000 * 1000))
	packet_count = (int64(period_test * 1000000000)) / period_nano
	test_c.numberCounter = 0
	b = packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, 128, number, test_type, test.test_type)

	//count_rez, per = test_c.testMax(b, c, addr, ifi.MTU, ipdst_1sfpsla_str, counter, test_type)
	count_rez, per = test_c.testThrGen(net_interface_name, b, c, addr, ifi.MTU, ipdst_1sfpsla_str, packet_count, period_nano, test_type)
	test.rez_128 = (float32)(128.0 * 8.0 * (float32)(count_rez) * 1000000 / (float32)(per))
	fmt.Println("->> rez_128 = ", count_rez, " period = ", per, " !!!  rez=", test.rez_128)

	size = 256
	period_nano = int64(size * 8 * 1000000000 / (test.thr_begin * 1000 * 1000))
	packet_count = (int64(period_test * 1000000000)) / period_nano
	test_c.numberCounter = 0
	b = packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, 256, number, test_type, test.test_type)

	//count_rez, per = test_c.testMax(b, c, addr, ifi.MTU, ipdst_1sfpsla_str, counter, test_type)
	count_rez, per = test_c.testThrGen(net_interface_name, b, c, addr, ifi.MTU, ipdst_1sfpsla_str, packet_count, period_nano, test_type)
	test.rez_256 = (float32)(256.0 * 8.0 * (float32)(count_rez) * 1000000 / (float32)(per))
	fmt.Println("->> rez_256 = ", count_rez, " period = ", per, " !!!  rez=", test.rez_256)

	size = 512
	period_nano = int64(size * 8 * 1000000000 / (test.thr_begin * 1000 * 1000))
	packet_count = (int64(period_test * 1000000000)) / period_nano
	test_c.numberCounter = 0
	b = packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, 512, number, test_type, test.test_type)

	//count_rez, per = test_c.testMax(b, c, addr, ifi.MTU, ipdst_1sfpsla_str, counter, test_type)
	count_rez, per = test_c.testThrGen(net_interface_name, b, c, addr, ifi.MTU, ipdst_1sfpsla_str, packet_count, period_nano, test_type)
	test.rez_512 = (float32)(512.0 * 8.0 * (float32)(count_rez) * 1000000 / (float32)(per))
	fmt.Println("->> rez_512 = ", count_rez, " period = ", per, " !!!  rez=", test.rez_512)

	size = 1024
	period_nano = int64(size * 8 * 1000000000 / (test.thr_begin * 1000 * 1000))
	packet_count = (int64(period_test * 1000000000)) / period_nano
	test_c.numberCounter = 0
	b = packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, size, number, test_type, test.test_type)
	count_rez, per = test_c.testThrGen(net_interface_name, b, c, addr, ifi.MTU, ipdst_1sfpsla_str, packet_count, period_nano, test_type)
	test.rez_1024 = (float32)(float32(size) * 8.0 * (float32)(count_rez) * 1000000 / (float32)(per))
	fmt.Println("->> rez_1024 = ", count_rez, " period = ", per, " !!!  rez=", test.rez_1024)

	size = 1280
	period_nano = int64(size * 8 * 1000000000 / (test.thr_begin * 1000 * 1000))
	packet_count = (int64(period_test * 1000000000)) / period_nano
	test_c.numberCounter = 0
	b = packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, size, number, test_type, test.test_type)
	count_rez, per = test_c.testThrGen(net_interface_name, b, c, addr, ifi.MTU, ipdst_1sfpsla_str, packet_count, period_nano, test_type)
	test.rez_1280 = (float32)(float32(size) * 8.0 * (float32)(count_rez) * 1000000 / (float32)(per))
	fmt.Println("->> rez_1280 = ", count_rez, " period = ", per, " !!!  rez=", test.rez_1280)

	size = 1500
	period_nano = int64(size * 8 * 1000000000 / (test.thr_begin * 1000 * 1000))
	packet_count = (int64(period_test * 1000000000)) / period_nano
	test_c.numberCounter = 0
	b = packetForm(ipsrc, ipdst1, ipdst2, ifi.HardwareAddr, mac_dst, size, number, test_type, test.test_type)
	count_rez, per = test_c.testThrGen(net_interface_name, b, c, addr, ifi.MTU, ipdst_1sfpsla_str, packet_count, period_nano, test_type)
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
	db.Exec("UPDATE test_throughput SET rez_64=?,rez_128=?,rez_256=?,rez_512=?,rez_1024=?,rez_1280=?, rez_1518=?,datetime=?, status=? WHERE id=?", test.rez_64, test.rez_128, test.rez_256, test.rez_512, test.rez_1024, test.rez_1280, test.rez_1518, time.Now(), test.status, id)

	db.Close()
}

type mutexCounter struct {
	mu sync.Mutex
	x  int64
}

func (c *mutexCounter) Inc() {
	c.mu.Lock()
	c.x -= 1
	c.mu.Unlock()
}

func (c *mutexCounter) Value() (x int64) {
	c.mu.Lock()
	x = c.x
	c.mu.Unlock()
	return
}

func (c *mutexCounter) Set(val int64) {
	c.mu.Lock()
	c.x = val
	c.mu.Unlock()
}

func (test *testThr) testThrGen(net_interface_name string, b []byte, c *raw.Conn, addr *raw.Addr, mtu int, ipdst_1sfpsla_str string, cnt int64, period_nano int64, t_type uint16) (int, int64) {
	time_to_gen := ((cnt * period_nano) * 150) / 100
	time_gen := time.Duration(cnt * period_nano)
	fmt.Println("		-- period_to_generate [ms] = ", (cnt*period_nano)/1000000)
	fmt.Println("		-- time_to_gen [ms]        = ", time_to_gen/1000000)
	fmt.Println("		-- time gen ", time_gen)
	fmt.Println("		-- cnt start= ", cnt)

	ifi, err := net.InterfaceByName(net_interface_name)
	if err != nil {
		return 0, 0
	}

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

	var addDelay bool
	addDelay = false

	if ((cnt*period_nano)/1000000000)*int64(pps_rez) > cnt {
		addDelay = true
	}

	var timerReal TimerR
	_ = timerReal.InitTimer()

	findSFP(ipdst_1sfpsla_str, ipdst_1sfpsla_str)

	//min_per_rez := int64(time.Since(gen_test_min_period_start)) / (1000 * 1000)
	//fmt.Println("		 -*- min period [mks] = ", min_per_rez)

	g_start := time.Now()
	fmt.Println(" == gStart ", g_start)
	rez_time := make(chan int64)
	test.numberCounter = 0

	K := 1
	quit := make(chan int, K+1)

	counter := mutexCounter{}

	if period_nano > 15000 {
		counter.Set(cnt)
		//for i := 0; i < 32; i++ {

		var test_type uint16
		test_type = 0x2000 + (uint16((*test).testID) & 0x1FFF)

		var netConfRecive *raw.Config = new(raw.Config)

		(*netConfRecive).Filter, _ = bpf.Assemble([]bpf.Instruction{
			// Проверка идентификатора пакета (34 бит) (xFA-от 1 ко 2, xFB – от 2 к 1, xFC – от 1 к Серверу)
			bpf.LoadAbsolute{Off: 34, Size: 1},
			bpf.JumpIf{Cond: bpf.JumpNotEqual, Val: 0xFC, SkipTrue: 5},
			// Проверка идентификатора теста
			bpf.LoadAbsolute{Off: 64, Size: 2},
			bpf.JumpIf{Cond: bpf.JumpNotEqual, Val: uint32(test_type), SkipTrue: 3},
			// Выбор одного из 1000
			bpf.LoadExtension{Num: bpf.ExtRand},
			//	bpf.JumpIf{Cond: bpf.JumpLessThan, Val: 0xFF, SkipFalse: 1},
			bpf.JumpIf{Cond: bpf.JumpGreaterThan, Val: 0x03FFFFFF, SkipTrue: 1},
			// Verdict is "send up to 4k of the packet to userspace."
			bpf.RetConstant{Val: 4096},
			// Verdict is "ignore packet."
			bpf.RetConstant{Val: 0},
		})

		conRcv, _ := raw.ListenPacket(ifi, etherType, netConfRecive)

		go (*test).receivePackets(conRcv, mtu, ipdst_1sfpsla_str, quit, t_type)

		go func() {
			ticker := time.NewTicker(time.Duration(period_nano))
			//timer := time.NewTimer(time.Microsecond * 10)
			//period := time.Duration(period_nano)
			//fmt.Println("Start")
			//for range ticker.C {
			for {
				select {
				case <-ticker.C:
					//time.Sleep(period)
					counter.Inc()
					c.WriteTo(b, addr)
					if counter.Value() <= 0 {
						//	rez_time <- (int64)(time.Since(time.Time(g_start)))
						break
					}
					if counter.Value()%10000 == 0 {
						if time.Since(g_start) >= time_gen {
							//rez_time <- (int64)(time.Since(time.Time(g_start)))
							break
						}
						//		}
					}
				}
			}
			rez_time <- (int64)(time.Since(time.Time(g_start)))
			ticker.Stop()
			//fmt.Println("		 -*- rez_time = ", rez_t)
		}()
		//}
	} else {

		blen := len(b)

		var test_type uint16
		test_type = 0x2000 + (uint16((*test).testID) & 0x1FFF)

		var netConfRecive *raw.Config = new(raw.Config)

		(*netConfRecive).Filter, _ = bpf.Assemble([]bpf.Instruction{
			// Проверка идентификатора пакета (34 бит) (xFA-от 1 ко 2, xFB – от 2 к 1, xFC – от 1 к Серверу)
			bpf.LoadAbsolute{Off: 34, Size: 1},
			bpf.JumpIf{Cond: bpf.JumpNotEqual, Val: 0xFC, SkipTrue: 5},
			// Проверка идентификатора теста
			bpf.LoadAbsolute{Off: 64, Size: 2},
			bpf.JumpIf{Cond: bpf.JumpNotEqual, Val: uint32(test_type), SkipTrue: 3},
			// Выбор одного из 1000
			bpf.LoadExtension{Num: bpf.ExtRand},
			//	bpf.JumpIf{Cond: bpf.JumpLessThan, Val: 0xFF, SkipFalse: 1},
			bpf.JumpIf{Cond: bpf.JumpGreaterThan, Val: 0x03FFFFFF, SkipTrue: 1},
			// Verdict is "send up to 4k of the packet to userspace."
			bpf.RetConstant{Val: 4096},
			// Verdict is "ignore packet."
			bpf.RetConstant{Val: 0},
		})

		conRcv, _ := raw.ListenPacket(ifi, etherType, netConfRecive)

		go (*test).receivePackets(conRcv, mtu, ipdst_1sfpsla_str, quit, t_type)

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

		counter.Set(cnt)

		for i := 0; i < K; i++ {
			go func() {

				con, err := raw.ListenPacket(ifi, etherType, netConf)
				defer con.Close()

				if err != nil {

					fmt.Println(" -!! Error !!-")
					fmt.Println(err)
					fmt.Println(" ----=====----")
					return
				}

			ExitLoop:
				for {
					select {
					case <-quit:
						fmt.Println(" == Quit ", (int64)(time.Since(g_start)))
						rez_time <- (int64)(time.Since(g_start))
						break ExitLoop
					default:
						n, err := con.WriteTo(b, addr)
						if addDelay {
							timerReal.timerDelayNano(period_nano)
						}
						if err != nil {
							fmt.Printf("%v", err)
							continue
						}
						if n < blen {
							fmt.Printf("Partial write: %d", n)
							continue
						}
						counter.Inc()
						if counter.Value() <= 0 {
							fmt.Println(" == cnt<0 ", (int64)(time.Since(g_start)))
							rez_time <- (int64)(time.Since(time.Time(g_start)))
							break ExitLoop
						}
						if counter.Value()%1000 == 0 {
							if time.Since(g_start) >= time_gen {
								fmt.Println(" == time out ", (int64)(time.Since(g_start)))
								rez_time <- (int64)(time.Since(time.Time(g_start)))
								break ExitLoop
							}
						}
					}
				}

			}()
		}

	}

	time.Sleep(time.Duration(time_to_gen))
	rez_count := test.numberCounter
	fmt.Println("		 --->> rez_count= ", rez_count*64)
	for i := 0; i < K; i++ {
		quit <- 1
	}
	time.Sleep(time.Millisecond * 10)
	rez := <-rez_time
	fmt.Println("		 --->> rez_counterRez= ", cnt-counter.Value())
	return int(cnt - counter.Value()), rez
}

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
			//log.Fatalf("failed to receive message: %v", err)
			fmt.Println(" -****- ")
			continue
		}

		// Unpack Ethernet II frame into Go representation.
		if err := (&f).UnmarshalBinary(b[:n]); err != nil {
			//log.Fatalf("failed to unmarshal ethernet frame: %v", err)
		}
		//fmt.Printf("\n\n--=Test %x - \n", f.Payload[12:16])
		var ips [4]byte
		copy(ips[:], (net.ParseIP(ipdst_1sfpsla_str)).To4())

		var t_ips [2]byte
		t_ips[1] = byte(t_type & 0xFF)
		t_ips[0] = byte((t_type >> 8) & 0xFF)
		//fmt.Printf("\n\n--=T_so %x - \n", ips)

		// Display source of message and message itself.
		if (len(f.Payload) >= 52) && (f.Payload[20] == 0xFC) && (bytes.Equal(f.Payload[12:16], ips[:]) == true) && (bytes.Equal(f.Payload[50:52], t_ips[:]) == true) {
			//count++
			//	fmt.Printf("-->>Detect")
			//	counter <-count
			//(*test).numberCounter = uint32(count)
			(*test).numberCounter++
		}
	}
}
