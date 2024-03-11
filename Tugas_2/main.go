//   Product Api:
//    version: 0.1
//    title: Product Api
//   Schemes: http, https
//   Host:
//   BasePath: /orders
//   BasePath: /items
//      Consumes:
//      - application/json
//   Produces:
//   - application/json
//   swagger:meta

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
