/*
sudo yum install net-snmp-utils


example change description for port
interface GigabitEthernet6/0/18

1) find if_id (in GigabitEthernet6/0/18)
snmpwalk -v2c -c public 10.0.0.65 1.3.6.1.2.1.2.2.1.2
IF-MIB::ifDescr.33 = STRING: GigabitEthernet6/0/18
return id, type string


2) get current description (in interface_id)
snmpget -v2c -c public 10.0.0.65 1.3.6.1.2.1.31.1.1.1.18.33
IF-MIB::ifAlias.33 = STRING:
-it is empty


3) set description (in interface_id)
snmpset -v2c -c publicrw 10.0.0.65 1.3.6.1.2.1.31.1.1.1.18.33 s "Emilius_test_description_22"
IF-MIB::ifAlias.33 = STRING: Emilius_test_description_22


4) optional get status
snmpget -v2c -c public 10.0.0.65 1.3.6.1.2.1.2.2.1.7.33
IF-MIB::ifAdminStatus.33 = INTEGER: up(1)



*/

package main

import (
	"fmt"
	"log"

	"strings"
	"time"

	g "github.com/soniah/gosnmp"
)

func main() {
	log.Println("start")
	communityRead := "public"
	communityWrite := "publicrw"
	device := &g.GoSNMP{
		Target:    "10.0.0.65",
		Port:      161,
		Version:   g.Version2c,
		Timeout:   time.Duration(7 * time.Second),
		Community: communityRead,
	}
	err := device.Connect()
	if err != nil {
		log.Panic(err)
	}
	defer device.Conn.Close()
	//1) find if_id (in GigabitEthernet6/0/18)
	list, err := device.BulkWalkAll("1.3.6.1.2.1.2.2.1.2")
	if err != nil {
		log.Panic(err)
	}
	printSnmpPdu(list)
	id := findText(list, "GigabitEthernet6/0/18")
	log.Println("id=", id)
	lid := getLastDigit(id)
	log.Println("last id", lid)
	//2) get current description (in interface_id)
	result, err := device.Get([]string{"1.3.6.1.2.1.31.1.1.1.18." + lid})
	if err != nil {
		log.Panic(err)
	}
	log.Println(result.Variables)
	printSnmpPdu(result.Variables)
	//3) set description (in interface_id)
	device.Community = communityWrite
	forSet := g.SnmpPDU{}
	forSet.Name = "1.3.6.1.2.1.31.1.1.1.18." + lid
	forSet.Type = g.OctetString
	forSet.Value = "Emilius description from go"
	res, err := device.Set([]g.SnmpPDU{forSet})
	if err != nil {
		log.Panic(err)
	}
	log.Println("##", res)
	// print current description
	result, err = device.Get([]string{"1.3.6.1.2.1.31.1.1.1.18." + lid})
	if err != nil {
		log.Panic(err)
	}
	//log.Println(result.Variables)
	printSnmpPdu(result.Variables)
}

func printSnmpPdu(list []g.SnmpPDU) {
	for _, l := range list {
		//fmt.Printf("%v\r\n", l)
		switch l.Type {
		case g.OctetString:
			fmt.Println(l.Name, string(l.Value.([]byte)))
		default:
			fmt.Printf("number: %d\n", g.ToBigInt(l.Value))
		}
	}
}

func findText(list []g.SnmpPDU, findStr string) (id string) {
	for _, l := range list {
		if string(l.Value.([]byte)) == findStr {
			return l.Name
		}
	}
	return ""
}

func getLastDigit(x string) string {
	return x[strings.LastIndex(x, ".")+1:]
}
