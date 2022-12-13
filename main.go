package main

import (
	"fmt"
	"gee/geecache"
	"gee/geeorm"
	_ "github.com/mattn/go-sqlite3"
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

	//test
	engine, _ := geeorm.NewEngine("sqlite3","gee.db")
	defer engine.Close()
	s := engine.NewSession()
	_,_ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	result, _ := s.Raw("INSERT INTO User(`Name`) VALUES (?),(?)","tom","sam").Exec()
	count, _ := result.RowsAffected()
	fmt.Printf("Exec success, %d affected\n", count)
}
