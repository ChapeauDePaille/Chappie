package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

//IP !
type IP struct {
	IDip    int
	Address string
}

//Device !
type Device struct {
	IDdevice int
	IPv4     string
	MAC      string
	Hostname string
	State    bool
	IDos     int
	IDport   int
}

//OS !
type OS struct {
	IDos         int
	Name         string
	Type         string
	Manufacturer string
	Family       string
	Generation   string
}

//Version !
type Version struct {
	IDversion int
	Name      string
}

//StatesPort !
type StatesPort struct {
	IDstateport int
	State       string
}

//Protocol !
type Protocol struct {
	IDprotocol int
	Name       string
}

//Service !
type Service struct {
	IDservice int
	Name      string
}

//Ports !
type Ports struct {
	IDport      int
	Number      int
	Used        bool
	IDversion   int
	IDprotocol  int
	IDstateport int
	IDservice   int
}

//InitDB : initialise la DB
func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("db nil")
	}
	return db
}

//CreateTable : CREATE TABLE IF NOT EXISTS (création des tables de la DB)
/*func CreateTable(db *sql.DB) {
	sqlTable := `
	CREATE TABLE IF NOT EXISTS cars(
		ID INT,
		Marque TEXT,
		Modele TEXT,
		InsertedDatetime DATETIME
		);
		`

	_, err := db.Exec(sqlTable)
	if err != nil {
		panic(err)
	}
}*/

//StoreDevice : INSERT (permet d'insérer un device avec plus de caractéristiques)
func StoreDevice(db *sql.DB, x Device) {
	sqlAdd := `
	INSERT OR REPLACE INTO devices(
		IPv4,
		MAC,
		Hostname,
		Up
		) values (?, ?, ?, ?)
		`
	stmt, err := db.Prepare(sqlAdd)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err2 := stmt.Exec(x.IPv4, x.MAC, x.Hostname, x.State)
	if err2 != nil {
		panic(err2)
	}
}

//StoreIP : INSERT (a servi à insérer les plus de 6000 IPs du tableau excel)
func StoreIP(db *sql.DB, x Device) {
	sqlAdd := `
	INSERT OR REPLACE INTO devices(
		IPv4
		) values (?)
		`
	stmt, err := db.Prepare(sqlAdd)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err2 := stmt.Exec(x.IPv4)
	if err2 != nil {
		panic(err2)
	}
}

//ReadIP : SELECT dans DB de test contenant uniquement des IPs
func ReadIP(db *sql.DB, a []string) ([]string, int) {
	sqlReader := `
	SELECT ID_device, IPv4 FROM devices
	WHERE ID_device=10
	`
	lecture, err := db.Query(sqlReader)
	if err != nil {
		panic(err)
	}
	defer lecture.Close()

	i := 0
	for lecture.Next() {
		d := Device{}
		err2 := lecture.Scan(&d.IDdevice, &d.IPv4)
		if err2 != nil {
			panic(err2)
		}
		a[i] = d.IPv4
		i++
	}
	return a, i
}
