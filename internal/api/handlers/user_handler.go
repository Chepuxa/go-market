package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/mail"
	"os"
	"training/proj/internal/api/models"
	"training/proj/internal/customerrors"
	"training/proj/internal/db/repositories"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgerrcode"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	UserRepository *repositories.UserRepository
}

func NewUserHandler(ur *repositories.UserRepository) *UserHandler {
	return &UserHandler{
		UserRepository: ur,
	}
}

type credentials struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (h *UserHandler) PostUser(w http.ResponseWriter, r *http.Request) {

	var userReq models.User

	decodeErr := json.NewDecoder(r.Body).Decode(&userReq)

	if decodeErr != nil {
		customerrors.BadRequestResponse(w, r, decodeErr)
		return
	}

	validate := validator.New()
	validErr := validate.Struct(userReq)

	if validErr != nil {
		customerrors.BadRequestResponse(w, r, validErr)
		return
	}

	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(userReq.Password), 12)

	if hashErr != nil {
		customerrors.BadRequestResponse(w, r, hashErr)
		return
	}

	userResp, crudErr := h.UserRepository.Create(&userReq, hashedPassword)

	if crudErr != nil {
		switch crudErr.Code {
		case pgerrcode.UniqueViolation:
			customerrors.EditConflictResponse(w, r)
		default:
			customerrors.ServerErrorResponse(w, r, crudErr)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userResp)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials credentials

	decodeErr := json.NewDecoder(r.Body).Decode(&credentials)

	if decodeErr != nil {
		customerrors.BadRequestResponse(w, r, decodeErr)
		return
	}

	validate := validator.New()
	validErr := validate.Struct(credentials)

	if validErr != nil {
		customerrors.BadRequestResponse(w, r, validErr)
		return
	}

	var user models.User
	var getUserErr error

	_, emailErr := mail.ParseAddress(credentials.Login)

	if emailErr != nil {
		user, getUserErr = h.UserRepository.GetByUsername(credentials.Login)
	} else {
		user, getUserErr = h.UserRepository.GetByEmail(credentials.Login)
	}

	if getUserErr == sql.ErrNoRows {
		customerrors.InvalidCredentialsResponse(w, r)
		return
	}

	if getUserErr != nil {
		customerrors.ServerErrorResponse(w, r, getUserErr)
		return
	}

	match, passErr := credentials.passwordMatches(user.Password)

	if passErr != nil {
		customerrors.ServerErrorResponse(w, r, passErr)
		return
	}

	if !match {
		customerrors.InvalidCredentialsResponse(w, r)
		return
	}

	token := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET_KEY")), nil)
	claims := map[string]interface{}{"user_id": user.UserID, "email": user.Email}
	_, tokenString, err := token.Encode(claims)

	if err != nil {
		customerrors.ServerErrorResponse(w, r, err)
		return
	}

	resp := struct {
		Token string `json:"token"`
	}{Token: tokenString}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (c *credentials) passwordMatches(plainText string) (bool, error) {
	log.Println(c.Password)
	log.Println(plainText)
	err := bcrypt.CompareHashAndPassword([]byte(plainText), []byte(c.Password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}
