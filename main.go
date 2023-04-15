package main

import (
	"bufio"
	"fmt"
	"github.com/go-redis/redis"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var n = 6 // number of searchs
var api_url = "https://www.googleapis.com/books/v1/volumes?q="
var timeAPI = make([]int64, n)
var timeRedis = make([]int64, n)

func formatBook(book string) string {
	return strings.ReplaceAll(book, " ", "")
}
func query(book string) string {
	start := time.Now()
	data := getApi(book)
	end := time.Now()
	delta := end.Sub(start)
	println("Tiempo de consulta API en ms: ", delta.Milliseconds())
	timeAPI = append(timeAPI, delta.Milliseconds())
	return data
}

func getApi(book string) string {
	response, err := http.Get(api_url + book)
	if err != nil {
		return err.Error()
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err.Error()
	}
	return string(responseData)
}
func verifyConnection(client *redis.Client) {
	pong, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pong)
}
func findInCache(book string, client1 *redis.Client, client2 *redis.Client, client3 *redis.Client) bool {
	result1 := client1.Get(book).Val()
	result2 := client2.Get(book).Val()
	result3 := client3.Get(book).Val()
	if result1 != "" || result2 != "" || result3 != "" {
		return true
	}
	return false
}
func getBook(book string, client1 *redis.Client, client2 *redis.Client, client3 *redis.Client) {
	result1 := client1.Get(book).Val()
	result2 := client2.Get(book).Val()
	result3 := client3.Get(book).Val()
	if result1 != "" {
		println("-------- REDIS 1 ----------")
		println(result1)

	} else if result2 != "" {
		println("-------- REDIS 2 ----------")
		println(result2)
	} else {
		println("-------- REDIS 3 ----------")
		println(result3)
	}
}
func readInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	fmt.Print("Ingrese un libro: ")
	scanner.Scan()
	return scanner.Text()
}
func main() {
	// Client connections
	client1 := redis.NewClient(&redis.Options{
		Addr:     "172.19.0.4:6379",
		Password: "",
		DB:       0,
	})
	client2 := redis.NewClient(&redis.Options{
		Addr:     "172.19.0.3:6379",
		Password: "",
		DB:       0,
	})
	client3 := redis.NewClient(&redis.Options{
		Addr:     "172.19.0.2:6379",
		Password: "",
		DB:       0,
	})
	verifyConnection(client1)
	verifyConnection(client2)
	verifyConnection(client3)
	//Configure Policy
	client1.ConfigSet("maxmemory-policy", "allkeys-lru")
	client2.ConfigSet("maxmemory-policy", "allkeys-lru")
	client3.ConfigSet("maxmemory-policy", "allkeys-lru")
	// API queries
	var book string
	id := 1          //id of search
	var capacity int // capacity range of index
	var range_ int   // current range
	if n > 3 {
		capacity = (n / 3) + 1
	} else {
		capacity = n
	}

	fmt.Println("Bienvenido al buscador de libros. Puedes buscar ", n, " veces.")
	for true {
		if n == 0 {
			fmt.Println("Resumen de tiempos de busqueda en ms: ")
			fmt.Println(" 	API: ", timeAPI)
			fmt.Println(" 	Redis: ", timeRedis)
			fmt.Println("La sesion ha terminado.")
			break
		} else {
			book = readInput()
			book = formatBook(book)
			println("..........................................................")
			if !findInCache(book, client1, client2, client3) { // create query and save in redis cache (if it doesn't exist)
				range_ = id / capacity
				println(range_)
				if range_ == 0 { //redis 1
					client1.Set(book, query(book), time.Duration(60)*time.Second)
				} else if range_ == 1 { //redis 2
					client2.Set(book, query(book), time.Duration(60)*time.Second)
				} else { //redis 3
					client3.Set(book, query(book), time.Duration(60)*time.Second)
				}
				getBook(book, client1, client2, client3)
			} else {
				// get query from redis cache
				start := time.Now()
				getBook(book, client1, client2, client3)
				end := time.Now()
				delta := end.Sub(start)
				println("Tiempo de consulta REDIS en ms: ", delta.Milliseconds())
				timeRedis = append(timeRedis, delta.Milliseconds()) // save the time of the query
			}

			println("..........................................................")
			id++
		}
		n--
	}
}
