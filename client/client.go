package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
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

	rowNo := 0

	for {

		rowNo++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal()
		}
		//fmt.Println(row[full_name], row[age], row[city])

		names := strings.Fields(string(row[full_name]))

		First_name := string(names[1])
		Second_name := string(names[0])
		age, err := strconv.Atoi(row[age])
		if err != nil {
			fmt.Fprintf(os.Stderr, "RowNo %d err %+v\n", rowNo, err)
			continue
		}
		city := string(row[city])

		fmt.Println(First_name, Second_name, age, city)

		httpPostUrl := "http://127.0.0.1:8080/user/registre"
		log.Println("Http POST JSON URL", httpPostUrl)

		var jsonDate = []byte({
			"first_name" : First_name,
			"second_name" : Second_name,
			"age" : age,
			"biography" : "test bio",
			"city" : city,
			"password" : "testPassword"
		})

		req, err := http.NewRequest("POST", httpPostUrl, bytes.NewBuffer(jsonDate))
		req.Header.Set("CContent-Type", "application/json; charset=UTF-8")
		if err != nil {
			log.Println(err)
		}

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			log.Println()
			panic(err)
		}
		defer res.Body.Close()

		log.Println("\nResponse Status", res.Status)
		log.Println("\nResponse Header", res.Header)
		body, _ := io.ReadAll(res.Body)
		log.Println("\nResponse body", string(body))
	}

}
