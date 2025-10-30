package auth

import (
	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/infrastructure/logger/sl"
	"edtech/internal/service"

	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	userService service.AuthService
}

func NewAuthHandler(usecase service.AuthService) *AuthHandler {
	return &AuthHandler{
		userService: usecase,
	}
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.auth.Register"

	log := logger.GetLogger(r.Context(), op)

	log.Info("register request started")

	var req dto.RegisterRequest

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed decode json", sl.Error(err))
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	resp, err := a.userService.Register(r.Context(), &req)
	if err != nil {
		log.Error("failed service ", sl.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError) // сделать обработчик ошибок
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, resp)
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	const op = "http.handlers.auth.Login"

	log := logger.GetLogger(r.Context(), op)

	var req dto.LoginRequest

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode json", sl.Error(err))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	//TODO сделать правильную валидацию и обработку ошибок
	if err := validate.Struct(req); err != nil {
		log.Error("failed to validate request", sl.Error(err))
		http.Error(w, "Параметры неверно указаны", http.StatusBadRequest)
		return
	}
	// TODO сделать правильную обработку всех ошибок
	response, err := a.userService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		log.Error("failed to login", sl.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// TODO сделать функцию для ответа клиенту
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)

}
