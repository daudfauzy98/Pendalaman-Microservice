package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/daudfauzy98/Pendalaman-Microservice/auth-service/database"
	"github.com/daudfauzy98/Pendalaman-Microservice/utils"
	"gorm.io/gorm"
)

type AuthDB struct {
	Db *gorm.DB
}

func (db *AuthDB) ValidateAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WrapAPIError(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	authToken := r.Header.Get("Authorization")

	res, err := database.ValidateAuth(authToken, db.Db)
	if err != nil {
		utils.WrapAPIError(w, r, err.Error(), http.StatusForbidden)
		return
	}

	utils.WrapAPIData(w, r, database.Auth{
		Username: res.Username,
		Token:    res.Token,
	}, http.StatusOK, "Success!")

	return
}

// Menjalankan fungsi Sign Up setelah menerima request dari endpoint /auth/signup
func (db *AuthDB) SignUp(w http.ResponseWriter, r *http.Request) {
	// Jika method request yang diterima bukan POST maka blok program
	// menampilkan error method not allowed
	if r.Method != "POST" {
		utils.WrapAPIError(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Inisialisasi variabel body berupa data-data yang di-request
	// Inisialisasi variabel err berupa pesan error yang diterima dari library ioutil
	// Kedua inisialisasi tersebut menggunakan bantuan library ioutil
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close() // Menghentikan proses pembacaan Body request
	if err != nil {      // Jika variabel err berisi pesan error (tidak kosong)
		utils.WrapAPIError(w, r, "Can't read body", http.StatusBadRequest)
		return
	}

	var signup database.Auth

	err = json.Unmarshal(body, &signup)
	if err != nil {
		utils.WrapAPIError(w, r, "Error unmarshal : "+err.Error(), http.StatusInternalServerError)
		return
	}

	signup.Token = utils.IDGenerator()

	err = signup.SignUp(db.Db)
	if err != nil {
		utils.WrapAPIError(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	utils.WrapAPISuccess(w, r, "Success!", http.StatusOK)
}

// Menjalankan fungsi Login setelah menerima request dari endpoint /auth/login
func (db *AuthDB) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		utils.WrapAPIError(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		utils.WrapAPIError(w, r, "Can't read body", http.StatusBadRequest)
		return
	}

	var login database.Auth

	err = json.Unmarshal(body, &login)
	if err != nil {
		utils.WrapAPIError(w, r, "Error unmarshal : "+err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := login.Login(db.Db)
	if err != nil {
		utils.WrapAPIError(w, r, "Error unmarshal : "+err.Error(), http.StatusBadRequest)
		return
	}

	utils.WrapAPIData(w, r, database.Auth{
		Username: res.Username,
		Token:    res.Token,
	}, http.StatusOK, "Success!")
	return
}
