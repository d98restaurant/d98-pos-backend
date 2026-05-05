package repository

import (
	"encoding/json"
	"fmt"
	"time"

	"pos-backend/internal/models"

	"github.com/dgraph-io/badger/v4"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(user *models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Active = true
	
	// Check if user already exists
	existing, _ := r.FindByUsername(user.Username)
	if existing != nil {
		return fmt.Errorf("user already exists")
	}
	
	seq, _ := GetNextSequence("user_id")
	user.ID = fmt.Sprintf("%d", seq)
	
	// Log the password hash being stored
	fmt.Printf("📝 Storing user: %s, hash length: %d\n", user.Username, len(user.PasswordHash))
	
	// Save by ID
	if err := SaveJSON("user:"+user.ID, user); err != nil {
		return err
	}
	// Save by username for lookup
	if err := SaveJSON("user:username:"+user.Username, user.ID); err != nil {
		return err
	}
	// Save by email for lookup
	if err := SaveJSON("user:email:"+user.Email, user.ID); err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var userID string
	err := GetJSON("user:username:"+username, &userID)
	if err == badger.ErrKeyNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return r.FindByID(userID)
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var userID string
	err := GetJSON("user:email:"+email, &userID)
	if err == badger.ErrKeyNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return r.FindByID(userID)
}

func (r *UserRepository) FindByID(id string) (*models.User, error) {
	var user models.User
	err := GetJSON("user:"+id, &user)
	if err == badger.ErrKeyNotFound {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		
		prefix := []byte("user:")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			key := string(it.Item().Key())
			// Skip index keys
			if len(key) > 13 && (key[5:13] == "username:" || key[5:10] == "email:") {
				continue
			}
			if len(key) > 10 && key[5:10] == "email:" {
				continue
			}
			
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var user models.User
				if err := json.Unmarshal(val, &user); err != nil {
					return err
				}
				users = append(users, user)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return users, err
}

func (r *UserRepository) Update(id string, updates map[string]interface{}) error {
	user, err := r.FindByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}
	
	if passwordHash, ok := updates["passwordHash"]; ok {
		user.PasswordHash = passwordHash.(string)
		fmt.Printf("📝 Updating password hash for user %s, new hash length: %d\n", user.Username, len(user.PasswordHash))
	}
	if role, ok := updates["role"]; ok {
		user.Role = models.UserRole(role.(string))
	}
	if active, ok := updates["active"]; ok {
		user.Active = active.(bool)
	}
	if lastLogin, ok := updates["lastLogin"]; ok {
		user.LastLogin = lastLogin.(*time.Time)
	}
	user.UpdatedAt = time.Now()
	
	return SaveJSON("user:"+id, user)
}

func (r *UserRepository) Delete(id string) error {
	user, err := r.FindByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return nil
	}
	
	// Delete indices
	DeleteKey("user:username:" + user.Username)
	DeleteKey("user:email:" + user.Email)
	// Delete main record
	return DeleteKey("user:" + id)
}

func (r *UserRepository) UpdateLastLogin(id string) error {
	now := time.Now()
	return r.Update(id, map[string]interface{}{"lastLogin": &now})
}
