package main

import (
	"net"
)

type testSLA struct {
	delay_solve    []int64
	delay_solve_to []int64
	delay_solve_un [4]int64

	delay_to_sum int64
	delay_un_sum int64

	number uint32
}

type testThr struct {
	testID        int
	numberTx      uint64
	numberCounter uint64

	period int16
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

	//          E000     test type:   1     1     1  - test Y.1544
}

type global_config struct {
	server_ip          string
	net_interface_name string
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
	chan_stop  chan int
}

type testThroughput struct {
	id             int
	test_type      int
	miss_init_test int
	module_first   int
	module_second  int
	thr_begin      int
	count          int
	ch_type        int
	max_loss       int
	rez_64         float32
	rez_128        float32
	rez_256        float32
	rez_512        float32
	rez_1024       float32
	rez_1280       float32
	rez_1518       float32
	rez_4096       float32
	rez_9000       float32
	status         int
}

type testDelay struct {
	id             int
	test_type      int
	miss_init_test int
	module_first   int
	module_second  int
	thr_begin      int
	count_packs    int
	count_tests    int

	net_interface_name string

	ipsrc  net.IP
	ipdst1 net.IP
	ipdst2 net.IP

	mac_src  []byte
	mac_dst  []byte
	mac_dst2 []byte

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

	id_test_type uint16
}

type testReal struct {
	id                  int
	test_type           int //	1 - "SFP-SLA1 - SFP-SLA2"
	name                string
	module_first        int
	module_second       int
	block_size          int
	clock               int
	count               int
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

type testBert struct {
	id               int
	test_type        int
	miss_init_test   int
	module_first     int
	module_second    int
	thr_begin        int
	count_prob_packs int
	count_probs      int
	rez_64           float32
	rez_128          float32
	rez_256          float32
	rez_512          float32
	rez_1024         float32
	rez_1280         float32
	rez_1518         float32
	rez_4096         float32
	rez_9000         float32
	status           int
}

type testY1564 struct {
	id            int
	test_type     int
	module_first  int
	module_second int
	block_size    int
	ToS           uint8
	VLAN_priority int

	CIR int
	EIR int
	TP  int

	period     int
	step_count int

	max_FTD float32
	max_FVD float32
	max_FLR float32

	net_interface_name string

	ipsrc  net.IP
	ipdst1 net.IP
	ipdst2 net.IP

	mac_src  []byte
	mac_dst  []byte
	mac_dst2 []byte

	numberTx    uint64
	numberRx    uint64
	numberRxRes uint64

	rez_IR_s1  float32
	rez_FTD_s1 float32
	rez_FVD_s1 float32
	rez_FLR_s1 float32

	rez_IR_s2  float32
	rez_FTD_s2 float32
	rez_FVD_s2 float32
	rez_FLR_s2 float32

	rez_IR_s3  float32
	rez_FTD_s3 float32
	rez_FVD_s3 float32
	rez_FLR_s3 float32

	rez_IR_s4  float32
	rez_FTD_s4 float32
	rez_FVD_s4 float32
	rez_FLR_s4 float32

	rez_IR_eir  float32
	rez_FTD_eir float32
	rez_FVD_eir float32
	rez_FLR_eir float32

	rez_IR_tp  float32
	rez_FTD_tp float32
	rez_FVD_tp float32
	rez_FLR_tp float32

	status int

	id_test_type      uint16
	id_test_type_temp uint16
}

type testLoss struct {
	id             int
	test_type      int
	miss_init_test int
	module_first   int
	module_second  int
	thr_begin      int
	step           int

	count_frames int
	count_steps  int

	net_interface_name string

	ipsrc  net.IP
	ipdst1 net.IP
	ipdst2 net.IP

	mac_src  []byte
	mac_dst  []byte
	mac_dst2 []byte

	numberTx    int64
	numberRx    uint64
	numberRxRes uint64

	status int

	id_test_type uint16
}

type testLossRez struct {
	id          int
	id_test     int
	step_number int

	rez_64   float32
	rez_128  float32
	rez_256  float32
	rez_512  float32
	rez_1024 float32
	rez_1280 float32
	rez_1518 float32
	rez_4096 float32
	rez_9000 float32
}

type testWaySFP struct {
	test_type     int
	module_first  int
	module_second int
	thr_begin     int

	packet_size int

	net_interface_name string
	addr               net.Addr
	conn               net.PacketConn

	ipsrc  net.IP
	ipdst1 net.IP
	ipdst2 net.IP

	period_gen_ms int64
	pause_ms      int64

	mac_src  []byte
	mac_dst  []byte
	mac_dst2 []byte

	status int

	id_test_type uint16

	SFP1_com_min float32
	SFP1_laz_min float32

	SFP2_com_min float32
	SFP2_laz_min float32
}
