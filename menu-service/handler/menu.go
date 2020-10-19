package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/daudfauzy98/Pendalaman-Microservice/menu-service/database"
	"github.com/daudfauzy98/Pendalaman-Microservice/utils"

	"gorm.io/gorm"
)

type MenuHandler struct {
	Db *gorm.DB
}

func (handler *MenuHandler) AddMenu(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // Relevan dengan -> r.Method != "POST"
		utils.WrapAPIError(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.WrapAPIError(w, r, "Can't read body", http.StatusInternalServerError)
		return
	}

	var menu database.Menu // Mereferensi ke folder database di file menu.go
	err = json.Unmarshal(body, &menu)

	if err != nil {
		utils.WrapAPIError(w, r, "Error unmarshal : "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = menu.Insert(handler.Db)
	if err != nil {
		utils.WrapAPIError(w, r, "Insert menu error : "+err.Error(), http.StatusInternalServerError)
	}

	utils.WrapAPISuccess(w, r, "Success add menu!", http.StatusOK)
}

func (handler *MenuHandler) GetMenu(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet { // Relevan dengan -> r.Method != "GET"
		utils.WrapAPIError(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	menu := database.Menu{}
	menus, err := menu.GetAll(handler.Db)

	if err != nil {
		utils.WrapAPIError(w, r, "Failed get menu : "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WrapAPIData(w, r, menus, 200, "Success get all menu!")
}
