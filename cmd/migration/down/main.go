package main

import (
	"Project/internal/config"
	"Project/internal/domain/entity"
	"Project/internal/infrastructure/database"
	"log"
)

func main() {
	cg := config.LoadConfig()

	db, err := database.Connect(cg)
	if err != nil {
		log.Fatal("Gagal Terhubung Kedalam Database")
	}
	err = db.Migrator().DropTable(
		&entity.User{},
		// tabel lainnya
	)
}
