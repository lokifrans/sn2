package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

type account_add struct {
	First_name  string `json:"first_name"`
	Second_name string `json:"second_name"`
	Age         int    `json:"age"`
	Biography   string `json:"biography"`
	City        string `json:"city"`
	Password    string `json:"password"`
}

const (
	full_name = iota
	age
	city
)

func main() {

	f, err := os.Open("people.csv")
	if err != nil {
		log.Fatal(err)
	}

	reader := csv.NewReader(f)

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal()
		}
		fmt.Println(row[full_name], row[age], row[city])
	}

	// httpPostUrl := "http://127.0.0.1:8080/user/registre"
	// log.Println("Http POST JSON URL", httpPostUrl)

	// var jsonDate = []byte(`
	// {
	// 	"first_name" : "test user",
	// 	"second_name" : "test user",
	// 	"age" : 100,
	// 	"biography" : "test user",
	// 	"city" : "test user",
	// 	"password" : "test user"
	// }`)

	// req, err := http.NewRequest("POST", httpPostUrl, bytes.NewBuffer(jsonDate))
	// req.Header.Set("CContent-Type", "application/json; charset=UTF-8")
	// if err != nil {
	// 	log.Println(err)
	// }

	// client := &http.Client{}
	// res, err := client.Do(req)
	// if err != nil {
	// 	log.Println()
	// 	panic(err)
	// }
	// defer res.Body.Close()

	// log.Println("\nResponse Status", res.Status)
	// log.Println("\nResponse Header", res.Header)
	// body, _ := io.ReadAll(res.Body)
	// log.Println("\nResponse body", string(body))

}
