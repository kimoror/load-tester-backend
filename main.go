package main

import (
	"load-tester-backend/modules/api"
	"load-tester-backend/modules/db"
	"load-tester-backend/modules/env"
)

func main() {
	env.Load()
	db.Init()
	api.Init()
}
