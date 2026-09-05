package repository

import (
	"Project/internal/domain/entity"
	dr "Project/internal/domain/repository"

	"gorm.io/gorm"
)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) dr.RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) GetAll() ([]entity.Room, error) {
	var rooms []entity.Room
	err := r.db.Preload("Seats").Order("id ASC").Find(&rooms).Error
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *RoomRepository) Create(room *entity.Room) error {
	return r.db.Create(room).Error
}

func (r *RoomRepository) FindById(id uint) (*entity.Room, error) {
	var room entity.Room
	err := r.db.Preload("Seats").First(&room, id).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *RoomRepository) Update(room *entity.Room) error {
	return r.db.Save(room).Error
}

func (r *RoomRepository) SoftDelete(id uint) error {
	return r.db.Delete(&entity.Room{}, id).Error
}

func (r *RoomRepository) ReplaceSeats(roomID uint, seats []entity.Seat) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("room_id = ?", roomID).Delete(&entity.Seat{}).Error; err != nil {
			return err
		}
		if len(seats) > 0 {
			if err := tx.Create(&seats).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *RoomRepository) SyncSeats(roomID uint, requestedSeatNames []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var existingSeats []entity.Seat
		if err := tx.Unscoped().Where("room_id = ?", roomID).Find(&existingSeats).Error; err != nil {
			return err
		}

		existingMap := make(map[string]entity.Seat)
		for _, seat := range existingSeats {
			existingMap[seat.Name] = seat
		}

		requestedMap := make(map[string]bool)
		for _, name := range requestedSeatNames {
			requestedMap[name] = true
		}

		var newSeats []entity.Seat
		for _, name := range requestedSeatNames {
			if existingSeat, exists := existingMap[name]; exists {
				if existingSeat.DeletedAt.Valid {
					if err := tx.Unscoped().Model(&existingSeat).Update("deleted_at", nil).Error; err != nil {
						return err
					}
				}
			} else {
				newSeats = append(newSeats, entity.Seat{
					RoomID: roomID,
					Name:   name,
				})
			}
		}

		if len(newSeats) > 0 {
			if err := tx.Create(&newSeats).Error; err != nil {
				return err
			}
		}

		for _, existingSeat := range existingSeats {
			if !requestedMap[existingSeat.Name] && !existingSeat.DeletedAt.Valid {
				if err := tx.Delete(&existingSeat).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *RoomRepository) HasActiveSchedules(roomID uint) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Schedule{}).
		Where("room_id = ? AND status <> ?", roomID, "cancelled").
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *RoomRepository) CountActiveSchedules(roomID uint) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Schedule{}).
		Where("room_id = ? AND status <> ?", roomID, "cancelled").
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
