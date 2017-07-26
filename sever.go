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
	rows, _ := db.Query(" select count(*) from sqlite_master where type='table' and name= 'sensor';", 1)
	defer rows.Close()
	var n int
	_ = rows.Scan(&n)
	if n == 0 {
		_, _ = db.Exec(
			`create table sensor (id integer primary key, temperature REAL , humidity REAL, timestanp text)`,
		)
	}
	// fmt.Println("insert into sensor values(" +  fmt.Sprint(datas.Temperature) + ", " +  fmt.Sprint(datas.Humidity) + ", \"" +  fmt.Sprint(time.Now()) + "\");")
	_, err = db.Exec("insert into sensor(temperature, humidity, timestanp) values(" + fmt.Sprint(datas.Temperature) + ", " + fmt.Sprint(datas.Humidity) + ", \"" + fmt.Sprint(time.Now()) + "\");")
	if err != nil {
		fmt.Println(err)
	}
	db.Close()
}

// how to post
//  curl http://localhost:8080/post -X POST -d "{"temperature":1.2,"humidity":3.3}"

func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/post", post)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
