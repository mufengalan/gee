package main

import (
	"database/sql"
	"fmt"
	"gee_demo/geecache"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *geecache.Group {
	return geecache.NewGroup("scores", 2<<10, geecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startCacheServer(addr string, addrs []string, gee *geecache.Group) {
	peers := geecache.NewHTTPPool(addr)
	peers.Set(addrs...)
	gee.RegisterPeers(peers)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, gee *geecache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := gee.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlices())

		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))

}

func main() {
	//var port int
	//var api bool
	//flag.IntVar(&port, "port", 8001, "Geecache server port")
	//flag.BoolVar(&api, "api", false, "Start a api server?")
	//flag.Parse()
	//
	//apiAddr := "http://localhost:9999"
	//addrMap := map[int]string{
	//	8001: "http://localhost:8001",
	//	8002: "http://localhost:8002",
	//	8003: "http://localhost:8003",
	//}
	//
	//var addrs []string
	//for _, v := range addrMap {
	//	addrs = append(addrs, v)
	//}
	//
	//gee := createGroup()
	//if api {
	//	go startAPIServer(apiAddr, gee)
	//}
	//startCacheServer(addrMap[port], []string(addrs), gee)

	////proto测试
	//test := &hello.Student{
	//	Name:   "geektutu",
	//	Male:   true,
	//	Scores: []int32{98, 85, 88},
	//}
	//data, err := proto.Marshal(test)
	//if err != nil {
	//	log.Fatal("marshaling error:", err)
	//}
	//newTest := &hello.Student{}
	//err = proto.Unmarshal(data, newTest)
	//if err != nil {
	//	log.Fatal("unmarshal error:", err)
	//}
	//
	//if test.GetName() != newTest.GetName(){
	//	log.Fatal("data mismatch %q != %q", test.GetName(), newTest.GetName())
	//}

	db, _ := sql.Open("sqlite3", "gee.db")
	defer func() { _ = db.Close() }()
	_, _ = db.Exec("DROP TABLE IF EXISTS User;")
	_, _ = db.Exec("CREATE TABLE User(Name text);")
	result, err := db.Exec("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam")
	if err == nil {
		affected, _ := result.RowsAffected()
		log.Println(affected)
	}
	row := db.QueryRow("SELECT Name FROM User LIMIT 1")
	var name string
	if err := row.Scan(&name); err == nil {
		log.Println(name)
	}
}
