package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
	"tugas3/routers"
)

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

func acakStatus() Status {
	return Status{
		Water: rand.Intn(100) + 1, // Dari 1 ke 100
		Wind:  rand.Intn(100) + 1, // Dari 1 ke 100
	}
}

func jsonSimpan(status Status, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777) //path file, write read only, create if tidak ada, trunc if ada, granted power privileges.
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(map[string]Status{"status": status})
	if err != nil {
		return err
	}

	fmt.Println("Updated status:", status)
	return nil
}

func rutinitas15(filePath string) {

	for {
		status := acakStatus()
		err := jsonSimpan(status, filePath)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		time.Sleep(15 * time.Second)
	}
}

var PORT = ":8080"

func main() {
	filePath := "./views/status.json"
	go rutinitas15(filePath)
	routers.StartServer().Run(PORT)
}
