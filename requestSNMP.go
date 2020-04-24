package main

import (
	"fmt"
//	"log"
	"time"

	g "github.com/soniah/gosnmp"
)

func findSFP(ip_1sfpsla_str string, ip_2sfpsla_str string) bool {
	fmt.Println(" ==> TEST SFP way ==")
	g.Default.Target = ip_1sfpsla_str
	err := g.Default.Connect()
	if err != nil {
		fmt.Printf("Connect() err: %v", err)
		//	log.Fatalf("Connect() err: %v", err)
	}
	defer g.Default.Conn.Close()

	oids := []string{".1.3.6.1.4.1.2010.1.13.0", ".1.3.6.1.4.1.2010.1.14.0"}

	timer_SNMP := time.NewTicker(100 * time.Millisecond)
	for range timer_SNMP.C {
		result, err2 := g.Default.Get(oids) // Get() accepts up to g.MAX_OIDS
		if err2 != nil {
			fmt.Printf("Get() err: %v", err2)

		}

		for i, variable := range result.Variables {
			fmt.Printf("%d: oid: %s ", i, variable.Name)

			// the Value of each variable returned by Get() implements
			// interface{}. You could do a type switch...
			switch variable.Type {
			case g.OctetString:
				fmt.Printf("string: %s\n", string(variable.Value.([]byte)))
			default:
				// ... or often you're just interested in numeric values.
				// ToBigInt() will return the Value as a BigInt, for plugging
				// into your calculations.
				fmt.Printf("number: %v   Byte/s\n", g.ToBigInt(variable.Value)) //*8)/1000000)
			}
		}
	}
	return true
}
