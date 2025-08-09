package api

import (
	"effective-mobile/logger"
	"effective-mobile/store"
	"effective-mobile/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type UsersHandler struct {
	store store.UsersStore
}

func NewUsersHandler(store store.UsersStore) *UsersHandler {
	return &UsersHandler{
		store: store,
	}
}

func (h *UsersHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var user store.User

	logger.LogInfo("CreateUser handler: creating new user...")

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logger.LogError("CreateUser handler: failed to decode body", err)
		utils.WriteResponse(w, http.StatusBadRequest, utils.Envelope{"error": http.StatusText(http.StatusBadRequest)})
		return
	}

	if err := h.store.CreateUser(&user); err != nil {
		logger.LogError("CreateUser handler: failed to create user", err)
		utils.WriteResponse(w, http.StatusInternalServerError, utils.Envelope{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	logger.LogInfo(fmt.Sprintf("CreateUser handler: user created successfully, uuid %s", user.ID.String()))

	utils.WriteResponse(w, http.StatusCreated, utils.Envelope{"user": user})
}

func (h *UsersHandler) HandleGetUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	logger.LogInfo(fmt.Sprintf("GetUserById handler: looking up user by id %s", id))

	if err := uuid.Validate(id); err != nil {
		logger.LogError("GetUserByID handler: failed to validate id", err)
		utils.WriteResponse(w, http.StatusBadRequest, utils.Envelope{"error": http.StatusText(http.StatusBadRequest)})
		return
	}

	userId := uuid.MustParse(id)
	user, err := h.store.GetUserByID(userId)

	if err != nil {
		logger.LogError("GetUserByID handler: failed to get user by id", err)
		utils.WriteResponse(w, http.StatusInternalServerError, utils.Envelope{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	if user == nil {
		logger.LogInfo("GetUserByID handler: user not found")
		utils.WriteResponse(w, http.StatusNotFound, utils.Envelope{"error": http.StatusText(http.StatusNotFound)})
		return
	}

	logger.LogInfo(fmt.Sprintf("GetUserById handler: user found, uuid %s", user.ID.String()))

	utils.WriteResponse(w, http.StatusOK, utils.Envelope{"user": user})
}
