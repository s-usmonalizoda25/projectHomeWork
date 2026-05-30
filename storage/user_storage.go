package storage

import (
	"encoding/json"
	"io"
	"os"
	"sync"
	"project/models"
)

type UserStorage struct {
	Mu sync.Mutex
	FileName string
}

func (s *UserStorage) GetAll() ([]models.User, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock() 
	if _, err := os.Stat(s.FileName); os.IsNotExist(err) {
		return []models.User{}, nil
	}
	file, err := os.Open(s.FileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if len(byteValue) == 0 {
		return []models.User{}, nil
	}
	var users []models.User
	err = json.Unmarshal(byteValue, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}


func (s *UserStorage) save(users []models.User) error {
	byteValue, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(s.FileName, byteValue, 0644)
	if err != nil {
		return err
	}
	return nil
}




func (s *UserStorage) Create(user models.User) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	var users []models.User
	if _, err := os.Stat(s.FileName); err == nil {
		byteValue, err := os.ReadFile(s.FileName)
		if err != nil {
			return err
		}
		if len(byteValue) > 0 {
			if err := json.Unmarshal(byteValue, &users); err != nil {
				return err
			}
		}
	}
	newID := 1
	if len(users) > 0 {
		newID = users[len(users)-1].ID + 1
	}
	user.ID = newID
	users = append(users, user)
	return s.save(users)
}



func (s *UserStorage) GetByID(id int) (*models.User, error) {
	users, err := s.GetAll()
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, nil
}


func (s *UserStorage) Update(id int, updatedUser models.User) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	var users []models.User
	byteValue, err := os.ReadFile(s.FileName)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if len(byteValue) > 0 {
		if err := json.Unmarshal(byteValue, &users); err != nil {
			return err
		}
	}
	found := false
	for i, user := range users {
		if user.ID == id {
			users[i].Name = updatedUser.Name
			found = true
			break
		}
	}
	if !found {
		return os.ErrNotExist
	}
	return s.save(users)
}
