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
		log.Fatal("Failed Connection Posgres")
	}
	err = db.AutoMigrate(
		&entity.Role{},
		&entity.Transaction{},
		&entity.Film{},
		&entity.Genre{},
		&entity.Room{},

		&entity.User{},
		&entity.Seat{},
		&entity.Source{},

		&entity.Ratting{},
		&entity.Bookmark{},
		&entity.Schedule{},
		&entity.ScheduleSeat{},

		&entity.Ticket{},
	)
	if err != nil {
		log.Fatal("Failed Migration")
	}
	log.Println("Success Migration")
}
