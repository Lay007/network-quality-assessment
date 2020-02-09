package main

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
}

type global_config struct {
	server_ip          string
	net_interface_name string
	zabbix_server_name string
	vlan               int
	vlan_number        int
}
type module_sfp struct {
	id          int
	name        string
	address_mac string
	address_ip  string
	version     string
	location    string
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

type testReal struct {
	id                int
	test_type         int
	name              string
	module_first      int
	module_second     int
	block_size        int
	clock             int
	count             int
	node_zabbix       string
	test_delay        bool
	test_delay_jitter bool
	test_loss         bool
	test_delay_1      bool
	test_delay_2      bool
	status            int
}
