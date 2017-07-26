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

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var datas Data
	if err := json.NewDecoder(r.Body).Decode(&datas); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	db(datas)
}

func db(datas Data) {
	db, err := sql.Open("sqlite3", "sensor.sqlite3")
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println("insert into sensor values(" +  fmt.Sprint(datas.Temperature) + ", " +  fmt.Sprint(datas.Humidity) + ", \"" +  fmt.Sprint(time.Now()) + "\");")
	_, err = db.Exec("insert into sensor(temperature, humidity, timestanp) values(" + fmt.Sprint(datas.Temperature) + ", " + fmt.Sprint(datas.Humidity) + ", \"" + fmt.Sprint(time.Now()) + "\");")
	if err != nil {
		fmt.Println(err)
	}
	db.Close()
}

func prepare() error {
	db, err := sql.Open("sqlite3", "sensor.sqlite3")
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = db.Exec(
		`create table if not exists sensor(id integer primary key, temperature REAL , humidity REAL, timestanp text)`,
	)
	return err
}

// how to post
//  curl http://localhost:8080/post -X POST -d "{"temperature":1.2,"humidity":3.3}"

func main() {
	if err := prepare(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", hello)
	http.HandleFunc("/post", post)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
