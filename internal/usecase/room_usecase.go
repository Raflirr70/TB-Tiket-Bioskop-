package usecase

import (
	"Project/internal/domain/entity"
	"Project/internal/domain/repository"
	du "Project/internal/domain/usecase"
	"fmt"
)

type RoomUsecaseImpl struct {
	roomRepo repository.RoomRepository
}

func NewRoomUsecase(roomRepo repository.RoomRepository) du.RoomUseCase {
	return &RoomUsecaseImpl{roomRepo: roomRepo}
}

func (u *RoomUsecaseImpl) GetAllRooms() ([]du.RoomResponse, error) {
	rooms, err := u.roomRepo.GetAll()
	if err != nil {
		return nil, err
	}

	result := make([]du.RoomResponse, 0)
	for _, room := range rooms {
		active, err := u.roomRepo.CountActiveSchedules(room.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, toRoomResponse(room, active))
	}
	return result, nil
}

func (u *RoomUsecaseImpl) CreateRoom(req du.CreateRoomRequest) (*du.RoomResponse, error) {
	var seats []entity.Seat
	total := 0

	for _, row := range req.Seats {
		for n := 1; n <= row.Count; n++ {
			seats = append(seats, entity.Seat{Name: fmt.Sprintf("%s%d", row.Label, n)})
		}
		total += row.Count
	}

	status := req.Status
	if status == "" {
		status = "ready"
	}

	room := &entity.Room{
		Name:     req.Name,
		Capasity: total,
		Status:   status,
		Seats:    seats,
	}

	if err := u.roomRepo.Create(room); err != nil {
		return nil, err
	}

	resp := toRoomResponse(*room, 0)
	return &resp, nil
}

func (u *RoomUsecaseImpl) UpdateRoom(id uint, req du.UpdateRoomRequest) (*du.RoomResponse, error) {
	room, err := u.roomRepo.FindById(id)
	if err != nil {
		return nil, err
	}
	if req.Name == "" {
		return nil, fmt.Errorf("nama ruangan wajib diisi")
	}

	status := req.Status
	if status == "" {
		status = room.Status
	} else if status != "ready" && status != "not ready" {
		return nil, fmt.Errorf("status tidak valid")
	}

	var requestedSeatNames []string
	total := 0
	for _, row := range req.Seats {
		for n := 1; n <= row.Count; n++ {
			requestedSeatNames = append(requestedSeatNames, fmt.Sprintf("%s%d", row.Label, n))
		}
		total += row.Count
	}

	room.Name = req.Name
	room.Status = status
	room.Capasity = total

	if err := u.roomRepo.SyncSeats(id, requestedSeatNames); err != nil {
		return nil, err
	}
	room.Seats = nil
	if err := u.roomRepo.Update(room); err != nil {
		return nil, err
	}

	updatedRoom, err := u.roomRepo.FindById(id)
	if err != nil {
		return nil, err
	}

	active, _ := u.roomRepo.CountActiveSchedules(id)
	resp := toRoomResponse(*updatedRoom, active)
	return &resp, nil
}

func (u *RoomUsecaseImpl) DeleteRoom(id uint) error {
	ok, err := u.roomRepo.HasActiveSchedules(id)
	if err != nil {
		return err
	}
	if ok {
		return fmt.Errorf("ruangan tidak dapat dihapus, masih ada jadwal aktif")
	}
	return u.roomRepo.SoftDelete(id)
}

func toRoomResponse(room entity.Room, active int64) du.RoomResponse {
	rows := []string{}
	seen := map[string]bool{}
	for _, s := range room.Seats {
		if len(s.Name) == 0 {
			continue
		}
		label := string(s.Name[0])
		if !seen[label] {
			seen[label] = true
			rows = append(rows, label)
		}
	}

	counts := make([]int, len(rows))
	for _, s := range room.Seats {
		if len(s.Name) == 0 {
			continue
		}
		label := string(s.Name[0])
		for i, r := range rows {
			if r == label {
				counts[i]++
				break
			}
		}
	}

	status := room.Status
	if status == "" {
		status = "ready"
	}

	return du.RoomResponse{
		ID:              room.ID,
		Name:            room.Name,
		Status:          status,
		Capacity:        room.Capasity,
		TotalSeats:      len(room.Seats),
		SeatRows:        rows,
		SeatCounts:      counts,
		ActiveSchedules: active,
	}
}
