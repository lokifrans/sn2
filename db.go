package main

import (
	"fmt"
	"log"
)

func CheckUser(cfg *apiConfig, id string, pas string) (bool, error) {
	var passwordForCheck string
	log.Printf("id = %v, password = %v", id, pas)

	queryPassword := "SELECT password FROM public.user WHERE id = $1"

	cfg.DB.Get(&passwordForCheck, queryPassword, id)

	if pas != passwordForCheck {
		log.Printf("ID: %v or Password: \"%v\" incorrect", id, pas)
		return false, fmt.Errorf("ID or Password incorrect")
	}
	log.Println("ID and Password correct")
	return true, nil

	//
}
