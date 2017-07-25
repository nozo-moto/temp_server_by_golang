package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"net/http"
	"unsafe"
	// "reflect"
	"time"
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
	t, _ := ioutil.ReadAll(r.Body)
	text := string(t)
	v := *(*[]byte)(unsafe.Pointer(&text))
	if err := json.Unmarshal(v, &datas); err != nil {
		fmt.Println(err)
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
	_, err = db.Exec("insert into sensor(temperature, humidity, timestanp) values(" +  fmt.Sprint(datas.Temperature) + ", " +  fmt.Sprint(datas.Humidity) + ", \"" +  fmt.Sprint(time.Now()) + "\");")
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
