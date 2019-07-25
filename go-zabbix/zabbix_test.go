package zabbix

import "testing"

const (
	hostname   = `localhost`
	zabbixhost = `127.0.0.1`
	zabbixport = 10051
)

func TestSend(t *testing.T) {
	sender := NewSender(zabbixhost, zabbixport)

	metrics := []*Metric{NewMetric(hostname, `key`, `value`)}
	_, err := sender.Send(NewPacket(metrics))

	if err == nil {
		t.Error("sending should have failed")
	}

	t.Logf("error: %v", err.Error())
}
