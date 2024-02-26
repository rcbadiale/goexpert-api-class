package handlers

import (
	"encoding/json"
	"goexpert-api/internal/dto"
	"goexpert-api/internal/entity"
	"goexpert-api/internal/infra/database"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
)

type UserHandler struct {
	UserService  database.UserInterface
	TokenAuth    *jwtauth.JWTAuth
	JWTExpiresIn int
}

func NewUserHandler(service database.UserInterface, tokenAuth *jwtauth.JWTAuth, jwtExpiresIn int) *UserHandler {
	return &UserHandler{
		UserService:  service,
		TokenAuth:    tokenAuth,
		JWTExpiresIn: jwtExpiresIn,
	}
}

// Get JWT godoc
// @Summary      Get a user JWT
// @Description  Get a user JWT
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body      dto.GetJWTInput true "user credentials"
// @Success      200      {object}  dto.GetJWTOutput
// @Failure      400      {object}  dto.ErrorOutput
// @Failure      403      {object}  dto.ErrorOutput
// @Failure      404      {object}  dto.ErrorOutput
// @Failure      500      {object}  dto.ErrorOutput
// @Router       /user/generate_token [post]
func (h *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	var userInput dto.GetJWTInput
	err := json.NewDecoder(r.Body).Decode(&userInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := dto.ErrorOutput{Message: "invalid format"}
		json.NewEncoder(w).Encode(error)
		return
	}

	user, err := h.UserService.FindByEmail(userInput.Email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error := dto.ErrorOutput{Message: "not found"}
		json.NewEncoder(w).Encode(error)
		return
	}
	if !user.ValidatePassword(userInput.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		error := dto.ErrorOutput{Message: "unauthorized"}
		json.NewEncoder(w).Encode(error)
		return
	}

	_, token, _ := h.TokenAuth.Encode(map[string]interface{}{
		"sub": user.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(h.JWTExpiresIn)).Unix(),
	})
	accessToken := dto.GetJWTOutput{AccessToken: token}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

// Create user godoc
// @Summary      Create user
// @Description  Create user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateUserInput true "user request"
// @Success      201
// @Failure      400      {object}  dto.ErrorOutput
// @Failure      500      {object}  dto.ErrorOutput
// @Router       /user [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := dto.ErrorOutput{Message: "invalid format"}
		json.NewEncoder(w).Encode(error)
		return
	}
	u, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := dto.ErrorOutput{Message: "invalid user data"}
		json.NewEncoder(w).Encode(error)
		return
	}
	err = h.UserService.Create(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := dto.ErrorOutput{Message: "error creating user"}
		json.NewEncoder(w).Encode(error)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
