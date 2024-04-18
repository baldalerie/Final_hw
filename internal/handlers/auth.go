package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"FinalTaskAppGoBasic/internal/logs"
	"FinalTaskAppGoBasic/internal/models"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const jwtKey = "FinalTaskAppGoBasic"

func (h *Handlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.Users
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		logs.HandleMessage(w, r, http.StatusBadRequest, "error decoding body")

		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logs.HandleMessage(w, r, http.StatusInternalServerError, "crypt password failed")

		return
	}

	user.Password = string(hashedPassword)

	result := h.gorm.Create(&user)
	if result.Error != nil {
		logs.HandleMessage(w, r, http.StatusBadRequest, "create user failed")

		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	user.Password = ""
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		logs.Log.WithError(err).Error("write response failed")

		return
	}

	logs.Log.
		WithField("user_id", user.ID).
		WithField("email", user.Email).
		WithField("user_name", user.Username).
		Info("user successfully added")
}

func (h *Handlers) LoginUser(w http.ResponseWriter, r *http.Request) {
	var rawUser models.Users
	err := json.NewDecoder(r.Body).Decode(&rawUser)
	if err != nil {
		logs.HandleMessage(w, r, http.StatusBadRequest, "error decoding body")

		return
	}

	var user models.Users
	result := h.gorm.Where("email = ?", user.Email).First(&user)
	if result.Error != nil {
		logs.HandleMessage(w, r, http.StatusBadRequest, "finding user by email failed")

		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(rawUser.Password))
	if err != nil {
		logs.HandleMessage(w, r, http.StatusUnauthorized, "bad credentials")

		return
	}

	expirationTime := time.Now().Add(30 * time.Minute)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     expirationTime.Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		logs.HandleMessage(w, r, http.StatusInternalServerError, "signing token failed")

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Token", tokenString)

	logs.Log.
		WithField("user_id", user.ID).
		WithField("email", user.Email).
		Info("user successfully logged in")
}
