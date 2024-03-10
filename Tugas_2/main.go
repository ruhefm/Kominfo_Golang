package main

import (
	"tugas2/database"
	"tugas2/routers"
)

var PORT = ":8000"

func main() {
	database.StartDB()
	routers.StartServer().Run(PORT)

}
