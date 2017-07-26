package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	// "reflect"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Data struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
}

type DB struct {
	sql.DB
}

func (db *DB) show(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
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
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func (db *DB) insert(datas Data) error {
	_, err := db.Exec("insert into sensor(temperature, humidity, timestanp) values(" + fmt.Sprint(datas.Temperature) + ", " + fmt.Sprint(datas.Humidity) + ", \"" + fmt.Sprint(time.Now()) + "\");")
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) prepare() error {
	_, err := db.Exec(
		`create table if not exists sensor(id integer primary key, temperature REAL , humidity REAL, timestanp text)`,
	)
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
