package main

import (
	"Project/internal/config"
	"Project/internal/delivery/http/handler"
	"Project/internal/delivery/http/router"
	"Project/internal/infrastructure/database"
	"Project/internal/repository"
	"Project/internal/usecase"
	"log"
)

func main() {
	//1. load config
	cg := config.LoadConfig()

	//2. inisialisasi database
	db, err := database.Connect(cg)
	if err != nil {
		log.Fatal("Fail Connection Database")
	}

	//3. repository
	userRepository := repository.NewUserRepository(db)
	roleRepository := repository.NewRoleRepository(db)

	//4. usecase
	// userUsecase := usecase.NewUserUsecase(userRepository)
	authUseCase := usecase.NewAuthUsecase(roleRepository, userRepository, cg)

	//5. handler
	authHandler := handler.NewAuthHandler(authUseCase)
	pageHandler := handler.NewPageHandler()

	//7. routes
	r := router.NewRouter(cg, pageHandler, authHandler)

	//8. run server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed Run Server : ", err)
	}
}
