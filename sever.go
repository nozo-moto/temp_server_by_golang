package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Data struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Timestamp   string  `json:"timestamp"`
}

type DB struct {
	sql.DB
}

func (db *DB) show(w http.ResponseWriter, r *http.Request) {
	datas, err := db.query()
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(datas); err != nil {
		// too late to return errors at least
		log.Println(err)
	}
}

func (db *DB) post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var datas Data
	if err := json.NewDecoder(r.Body).Decode(&datas); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err := db.insert(datas); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (db *DB) query() ([]Data, error) {
	rows, err := db.Query(`
	select temperature, humidity, timestamp from sensor order by timestamp desc limit 5
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Data
	for rows.Next() {
		var data Data
		err = rows.Scan(&data.Temperature, &data.Humidity, &data.Timestamp)
		if err != nil {
			return nil, err
		}
		res = append(res, data)
	}
	return res, nil
}

func (db *DB) insert(datas Data) error {
	_, err := db.Exec(`
	insert into sensor(temperature, humidity, timestamp)
	values(?, ?, ?)
	`, datas.Temperature, datas.Humidity, time.Now().String())
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) prepare() error {
	_, err := db.Exec(`
	create table if not exists sensor(id integer primary key, temperature REAL , humidity REAL, timestamp text)
	`)
	return err
}

// how to post
//  curl http://localhost:8080/post -X POST -d "{"temperature":1.2,"humidity":3.3}"

func main() {
	db, err := sql.Open("sqlite3", "sensor.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	conn := &DB{*db}

	if err := conn.prepare(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", conn.show)
	http.HandleFunc("/post", conn.post)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
