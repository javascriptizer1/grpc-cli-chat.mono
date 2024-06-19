package manager

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type FileTokenManager struct {
	once     sync.Once
	mu       sync.Mutex
	token    Token
	filePath string
}

func NewFileTokenManager(filePath string) *FileTokenManager {
	return &FileTokenManager{
		filePath: filePath,
	}
}

func (tm *FileTokenManager) Load() error {
	tm.once.Do(func() {
		tm.mu.Lock()

		defer tm.mu.Unlock()

		file, _ := os.Open(tm.filePath)

		defer func() {
			_ = file.Close()
		}()

		_ = json.NewDecoder(file).Decode(&tm.token)

	})

	return nil
}

func (tm *FileTokenManager) Save() error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if err := tm.ensureDirExists(); err != nil {
		return err
	}

	file, err := os.Create(tm.filePath)

	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}

	defer func() {
		_ = file.Close()
	}()

	err = json.NewEncoder(file).Encode(&tm.token)

	if err != nil {
		fmt.Println("Error encoding tokens:", err)
		return err
	}

	return nil
}

func (tm *FileTokenManager) SetTokens(accessToken, refreshToken string) error {
	if accessToken != "" {
		tm.token.AccessToken = accessToken
	}

	if refreshToken != "" {
		tm.token.RefreshToken = refreshToken
	}

	return tm.Save()
}

func (tm *FileTokenManager) AccessToken() string {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	return tm.token.AccessToken
}

func (tm *FileTokenManager) RefreshToken() string {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	return tm.token.RefreshToken
}

func (tm *FileTokenManager) ensureDirExists() error {

	dir := filepath.Dir(tm.filePath)

	if _, err := os.Stat(dir); os.IsNotExist(err) {

		err := os.MkdirAll(dir, 0750)

		if err != nil {
			return err
		}
	}

	return nil
}
