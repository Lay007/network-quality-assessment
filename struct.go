package main

//import "time"

type testSLA struct {
	delay_solve    []int64
	delay_solve_to []int64
	delay_solve_un []int64
	number         uint32
}

type testThr struct {
	testID        int
	numberTx      uint32
	numberCounter uint32
}

type iphdr struct {
	vhl   uint8
	tos   uint8
	iplen uint16
	id    uint16
	off   uint16
	ttl   uint8
	proto uint8
	csum  uint16
	src   [4]byte
	dst   [4]byte
}

type sfpsla struct {
	id          uint8
	dst         [4]byte
	merkertime1 [7]byte
	merkertime2 [7]byte
	merkertime3 [7]byte
	number      uint32
	test_type   uint16
	//						     /- test type -\  /----           test ID          ----\
	//                            15    14    13   12 11 10 09 08 07 06 05 04 03 02 01 00
	//                   test type:   0     0     0  - test Real SLA

	//          2000     test type:   0     0     1  - test RFC 2544 - Пропускная способность (нагрузка)
	//          4000     test type:   0     1     0  - test RFC 2544 - Пропускная способность (мониторинг)
	//          2000     test type:   0     0     1  - test RFC 2544 - Задержка (нагрузка)
	//          6000     test type:   0     1     1  - test RFC 2544 - Задержка (мониторинг)
	//          8000     test type:   1     0     0  - test RFC 2544 - Потеря пакетов (нагрузка)
	//          A000     test type:   1     0     1  - test RFC 2544 - Потеря пакетов (мониторинг)
	//          C000     test type:   1     1     0  - test RFC 2544 - Берстность
	//          F000     test type:   1     1     1  - test Y.1544
}

type global_config struct {
	server_ip          string
	net_interface_name string
	zabbix_server_name string
	zabbix_server_port int
	vlan               int
	vlan_number        int
	QinQ               int
	QinQ_number        int
}
type module_sfp struct {
	id         int
	addres_mac int64
	name       string
	address_ip string
	version    string
	location   string
}

type testThroughput struct {
	id            int
	test_type     int
	module_first  int
	module_second int
	thr_begin     int
	count         int
	ch_type       int
	max_loss      int
	rez_64        float32
	rez_128       float32
	rez_256       float32
	rez_512       float32
	rez_1024      float32
	rez_1280      float32
	rez_1518      float32
	rez_4096      float32
	rez_9000      float32
	status        int
}

type testDelay struct {
	id            int
	test_type     int
	module_first  int
	module_second int
	thr_begin     int
	count_packs   int
	count_tests   int

	rez_64       float32
	rez_64_max   float32
	rez_64_min   float32
	rez_128      float32
	rez_128_max  float32
	rez_128_min  float32
	rez_256      float32
	rez_256_max  float32
	rez_256_min  float32
	rez_512      float32
	rez_512_max  float32
	rez_512_min  float32
	rez_1024     float32
	rez_1024_max float32
	rez_1024_min float32
	rez_1280     float32
	rez_1280_max float32
	rez_1280_min float32
	rez_1518     float32
	rez_1518_max float32
	rez_1518_min float32
	status       int
}

type testReal struct {
	id        int
	test_type int //	1 - "SFP-SLA1 - SFP-SLA2"
	//  2 - "Сервер - SFP-SLA1"
	//  3 - "Сервер - SFP-SLA2"
	name                string
	module_first        int
	module_second       int
	block_size          int
	clock               int
	count               int
	node_zabbix         string
	test_delay          bool
	test_delay_jitter   bool
	test_loss           bool
	test_delay_1        bool
	test_delay_1_jitter bool
	datetime            string
	status              int
}
type testRealMax struct {
	delayMax     float32
	jitterMax    float32
	delayOneMax  float32
	jitterOneMax float32
	lossMax      float32
}
