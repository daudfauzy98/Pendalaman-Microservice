package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/daudfauzy98/Pendalaman-Microservice/menu-service/entity"
	"github.com/gorilla/context"

	"github.com/daudfauzy98/Pendalaman-Microservice/menu-service/config"
	"github.com/daudfauzy98/Pendalaman-Microservice/utils"
)

// Untuk menghandler autentifikasi
type AuthHandler struct {
	Config config.AuthService
}

// Menjalankan validasi terlebih dahulu -> baru next handler
func (handler *AuthHandler) ValidateAdmin(nextHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := http.NewRequest(http.MethodPost, handler.Config.Host+"/auth/validate", nil)
		if err != nil {
			utils.WrapAPIError(w, r, err.Error(), http.StatusInternalServerError)
			return
		}

		request.Header = r.Header
		authResponse, err := http.DefaultClient.Do(request)

		if err != nil {
			log.Println("ERROR DISINI 1", handler.Config.Host)
			utils.WrapAPIError(w, r, "Validate auth failed : "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer authResponse.Body.Close()

		responseBody, err := ioutil.ReadAll(authResponse.Body)

		if err != nil {
			log.Println("ERROR DISINI 2")
			utils.WrapAPIError(w, r, err.Error(), http.StatusInternalServerError)
			return
		}

		var responseData entity.AuthResponse
		err = json.Unmarshal(responseBody, &responseBody)

		if err != nil {
			utils.WrapAPIError(w, r, err.Error(), http.StatusInternalServerError)
			return
		}

		if authResponse.StatusCode != http.StatusOK { // http.statusOK = kode status 200
			utils.WrapAPIError(w, r, responseData.ErrorDetails, authResponse.StatusCode)
			return
		}

		context.Set(r, "user", responseData.Data.Username)
		nextHandler(w, r)
	}
}
