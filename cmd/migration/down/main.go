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
		&entity.Role{},
		&entity.Transaction{},
		&entity.Film{},
		&entity.Genre{},
		&entity.Room{},

		&entity.User{},
		&entity.Seat{},
		&entity.Source{},

		&entity.Comment{},
		&entity.Bookmark{},
		&entity.Schedule{},
		&entity.ScheduleSeat{},

		&entity.Ticket{},
		"genre_films",
	)
	if err != nil {
		log.Fatal("Failed Rollback")
	}
	log.Println("Success Rollback")
}
