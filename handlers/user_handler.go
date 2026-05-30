package handlers

import (
	"encoding/json"
	"net/http"
	"project/storage"
	"project/models" 
	"strconv"
	"os"
)

type UserHandler struct {
	Storage *storage.UserStorage
}
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Метод не поддерживается. Используйте GET"))
		return
	}

	users, err := h.Storage.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка при чтении данных"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}



func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Метод не поддерживается. Используйте POST"))
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Некорректный JSON в теле запроса"))
		return
	}

	if user.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Поле 'name' не может быть пустым"))
		return
	}

	err = h.Storage.Create(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка при сохранении пользователя"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Пользователь успешно создан!"}`))
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Метод не поддерживается. Используйте GET"))
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Некорректный ID пользователя"))
		return
	}
	user, err := h.Storage.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка при поиске пользователя"))
		return
	}

	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 Not Found: Пользователь не найден"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Метод не поддерживается. Используйте PUT"))
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Некорректный ID пользователя"))
		return
	}

	var updatedUser models.User
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Некорректный JSON"))
		return
	}
	err = h.Storage.Update(id, updatedUser)
	if err != nil {
		if err == os.ErrNotExist {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 Not Found: Пользователь для обновления не найден"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка при обновлении данных"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Данные пользователя успешно обновлены!"}`))
}

