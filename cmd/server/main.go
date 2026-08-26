package main

import (
	"Project/internal/config"
	"Project/internal/delivery/http/handler"
	"Project/internal/delivery/http/router"
	"log"
)

func main() {
	//1. load config
	cg := config.LoadConfig()

	//2. inisialisasi database
	// db, err := database.Connect(cg)
	// if err != nil {
	// 	log.Fatal("Fail Connection Database")
	// }

	//3. repository
	// userRepository := repository.NewUserRepository(db)

	//4. usecase
	// userUsecase := usecase.NewUserUsecase(userRepository)

	//5. handler
	// userHandler := handler.NewUserHandler(userUsecase)
	pageHandler := handler.NewPageHandler()

	//7. routes
	r := router.NewRouter(cg, pageHandler)

	//8. run server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed Run Server : ", err)
	}
}
