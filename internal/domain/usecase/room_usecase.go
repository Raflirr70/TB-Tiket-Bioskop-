package usecase

type RoomUseCase interface {
	GetAllRooms() ([]RoomResponse, error)
	CreateRoom(req CreateRoomRequest) (*RoomResponse, error)
	UpdateRoom(id uint, req UpdateRoomRequest) (*RoomResponse, error)
	DeleteRoom(id uint) error
}

type RoomResponse struct {
	ID              uint     `json:"id"`
	Name            string   `json:"name"`
	Status          string   `json:"status"`
	Capacity        int      `json:"capacity"`
	TotalSeats      int      `json:"total_seats"`
	SeatRows        []string `json:"seat_rows"`
	SeatCounts      []int    `json:"seat_counts"`
	ActiveSchedules int64    `json:"active_schedules"`
}

type SeatRowRequest struct {
	Label string `json:"label"`
	Count int    `json:"count"`
}

type CreateRoomRequest struct {
	Name   string           `json:"name"`
	Status string           `json:"status"`
	Seats  []SeatRowRequest `json:"seats"`
}

type UpdateRoomRequest struct {
	Name   string           `json:"name"`
	Status string           `json:"status"`
	Seats  []SeatRowRequest `json:"seats"`
}
