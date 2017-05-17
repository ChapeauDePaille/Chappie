package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os/exec"
)

//Ensemble de structures permettant la lecture dans un fichier XML

//NmapRun !
type NmapRun struct {
	Scanner string `xml:"scanner,attr"`
	Version string `xml:"version,attr"`
	Hosts   []Host `xml:"host" `
}

//Host !
type Host struct {
	Addresses   []Address  `xml:"address"`
	Hostnames   []Hostname `xml:"hostnames>hostname"`
	Os          Os         `xml:"os"`
	HostScripts []Script   `xml:"hostscript>script"`
	Trace       Trace      `xml:"trace"`
}

//Address !
type Address struct {
	Addr     string `xml:"addr,attr"`
	AddrType string `xml:"addrtype,attr"`
	Vendor   string `xml:"vendor,attr"`
}

//Hostname !
type Hostname struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

//Os !
type Os struct {
	PortsUsed []PortUsed `xml:"portused"`
}

//PortUsed !
type PortUsed struct {
	State  string `xml:"state,attr"`
	Proto  string `xml:"proto,attr"`
	PortID int    `xml:"portid,attr"`
}

//Script contains information from Nmap Scripting Engine.
type Script struct {
	ID     string  `xml:"id,attr"`
	Output string  `xml:"output,attr"`
	Tables []Table `xml:"table"`
}

//Table contains the output of the script in a more parse-able form.
//ToDo: This should be a map[string][]string
type Table struct {
	Key      string   `xml:"key,attr"`
	Elements []string `xml:"elem"`
}

//Trace !
type Trace struct {
	Proto string `xml:"proto,attr"`
	Port  int    `xml:"port,attr"`
	Hops  []Hop  `xml:"hop"`
}

//Hop is a ip hop to a Host.
type Hop struct {
	TTL    float32 `xml:"ttl,attr"`
	RTT    float32 `xml:"rtt,attr"`
	IPAddr string  `xml:"ipaddr,attr"`
	Host   string  `xml:"host,attr"`
}

func main() {

	//Initialisation de la DB
	const dbpath = "test.db"

	db := InitDB(dbpath)
	defer db.Close()

	//Stockage des adresses dans la DB de test

	//Lecture fichier excel (xlsx)
	var a = []string{}
	a = ReadXlsx()
	//fmt.Println(len(a))
	//fmt.Println(a[0])

	/*var d Device
	for i := 0; i < len(a); i++ {
		d.IPv4 = a[i]
		StoreIP(db, d)
	}*/

	//Lecture des IP dans la DB de test et s'en servir pour un scan
	var address []string
	var count int
	address = make([]string, len(a))
	address, count = ReadIP(db, address)
	fmt.Println("On souhaite scanner la (ou les) adresse(s) suivante(s) :")
	for i := 0; i < count; i++ {
		fmt.Println(address[i], "; ")
	}

	//Adresse(s) IP que l'on veut scanner
	var h []string
	h = make([]string, count)
	//fmt.Println(cap(h))
	for i := 0; i < count; /*len(h)*/ i++ {
		h[i] = address[i]
	}

	//Fichier(s) XML où l'on veut sauvegarder l'output
	var f []string
	f = make([]string, len(a))
	//fmt.Println(cap(f))
	for i := 0; i < count; /*len(h)*/ i++ {
		numero := address[i]
		f[i] = numero + ".xml"
	}

	//Lancement de scan
	n := NmapRun{}
	for i := 0; i < count; /*len(h)*/ i++ {

		RunScan(h[i], f[i])

		//Récupération des données du fichier XML
		nmaprun, _ := ioutil.ReadFile(f[i])
		err := xml.Unmarshal(nmaprun, &n)
		if err != nil {
			panic(err)
		}

		//Affichage de certaines données
		Affichage(n)
	}

	//Recup traceroute du XML
	var trace = []string{}
	var traceCount int
	trace = make([]string, len(a))
	trace, traceCount = RecupTrace(n, trace)
	for i := 0; i < traceCount; i++ {
		fmt.Println(trace[i])
	}

	//Insertion des IPs traceroute dans la DB de test
	/*var d Device
	for i := 0; i < traceCount; i++ {
		d.IPv4 = trace[i]
		StoreIP(db, d)
	}*/

}

//RunScan : lance un scan Nmap
func RunScan(hosts string, fichier string) ([]Host, error) {

	var n NmapRun

	out := "-oX"
	path := "/home/anton/Documents/Projets/Go/src/github.com/ChapeauDePaille/Chappie/" + fichier
	param := "-T4"
	param2 := "-A"
	param3 := "-v"
	args := []string{out, path, param, param2, param3}
	args = append(args, hosts)

	output, err := exec.Command("nmap", args...).Output()
	fmt.Printf("%s", output)

	if err != nil {
		return n.Hosts, err
	}

	if err := xml.Unmarshal(output, &n); err != nil {
		return n.Hosts, err
	}

	return n.Hosts, nil
}

//Affichage des données
func Affichage(n NmapRun) {

	fmt.Println("Scanner :", n.Scanner)
	fmt.Println("Version :", n.Version)

	//Recup hostname/type
	for i := 0; i < len(n.Hosts); i++ {
		for j := 0; j < len(n.Hosts[i].Hostnames); j++ {
			fmt.Println("Name :", n.Hosts[i].Hostnames[j].Name)
			fmt.Println("Type :", n.Hosts[i].Hostnames[j].Type)
		}
	}

	//Recup adresse/typeadresse/vendeur
	for i := 0; i < len(n.Hosts); i++ {
		for j := 0; j < len(n.Hosts[i].Addresses); j++ {
			fmt.Println("Adresse :", n.Hosts[i].Addresses[j].Addr)
			fmt.Println("Type d'adresse :", n.Hosts[i].Addresses[j].AddrType)
			if n.Hosts[i].Addresses[j].AddrType == "mac" {
				fmt.Println("Vendeur :", n.Hosts[i].Addresses[j].Vendor)
			}
		}
	}

	//Recup ports/protocoles/states
	for i := 0; i < len(n.Hosts); i++ {
		for j := 0; j < len(n.Hosts[i].Os.PortsUsed); j++ {
			fmt.Println("Numéro :", n.Hosts[i].Os.PortsUsed[j].PortID)
			fmt.Println("Protocole :", n.Hosts[i].Os.PortsUsed[j].Proto)
			fmt.Println("Etat :", n.Hosts[i].Os.PortsUsed[j].State)
		}
	}

	//Recup MAC
	for i := 0; i < len(n.Hosts); i++ {
		for j := 0; j < len(n.Hosts[i].HostScripts); j++ {
			for k := 0; k < len(n.Hosts[i].HostScripts[j].Tables); k++ {
				if n.Hosts[i].HostScripts[j].Tables[k].Key == "mac" {
					fmt.Print(n.Hosts[i].HostScripts[j].Tables[k].Key, ": ")
					fmt.Println(n.Hosts[i].HostScripts[j].Tables[k].Elements[0])

				}
			}
		}
	}

	//Recup IP traceroute
	for i := 0; i < len(n.Hosts); i++ {
		for j := 0; j < len(n.Hosts[i].Trace.Hops)-1; j++ { //La dernière adresse du traceroute est forcément l'IP scanné
			fmt.Println(n.Hosts[i].Trace.Hops[j].IPAddr)
		}
	}
}

//RecupTrace : permet de récupérer l'adresses IP de la traceroute pour les mettre en DB en vue d'un scan
func RecupTrace(n NmapRun, t []string) ([]string, int) {
	var c int
	for i := 0; i < len(n.Hosts); i++ {
		c = len(n.Hosts[i].Trace.Hops) - 1
		for j := 0; j < len(n.Hosts[i].Trace.Hops)-1; j++ { //La dernière adresse du traceroute est forcément l'IP scanné
			t[j] = n.Hosts[i].Trace.Hops[j].IPAddr
		}
	}
	return t, c
}
