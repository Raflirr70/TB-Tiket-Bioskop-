package handler

import (
	du "Project/internal/domain/usecase"
	"Project/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	roomUsecase du.RoomUseCase
}

func NewRoomHandler(roomUsecase du.RoomUseCase) *RoomHandler {
	return &RoomHandler{roomUsecase: roomUsecase}
}

func (h *RoomHandler) GetAllRooms(c *gin.Context) {
	rooms, err := h.roomUsecase.GetAllRooms()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, 200, rooms)
}

func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var req du.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	if req.Name == "" {
		response.Error(c, 400, "nama ruangan wajib diisi")
		return
	}

	room, err := h.roomUsecase.CreateRoom(req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, 201, room)
}

func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "id tidak valid")
		return
	}

	var req du.UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	room, err := h.roomUsecase.UpdateRoom(uint(id64), req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, 200, room)
}

func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "id tidak valid")
		return
	}

	err = h.roomUsecase.DeleteRoom(uint(id64))
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, 200, "ruangan berhasil dihapus")
}
