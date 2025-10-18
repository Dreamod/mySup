package gSheets

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
	"log"
	"os"

	"golang.org/x/oauth2"
)

var tokenFileName string = "gSheets/token.json"

// GetService получает клиент сервис
func GetService() *sheets.Service {
	// получаем конфиг
	config, err := getConfig()
	if err != nil {
		log.Fatalf("Не удалось создать конфигурацию: %v", err)
	}
	// получаем токен доступа
	token, err := getTokenFromFile()
	if err != nil {
		token = getTokenFromWeb(config)
		saveToken(token)
	}
	// получаем клиент
	client := config.Client(context.Background(), token)
	// получаем сервис
	service, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Не удалось создать сервис: %v", err)
	}
	return service
}

// получает креды доступа из файла
func getCredentials() ([]byte, error) {
	data, err := os.ReadFile("gSheets/credentials.json")
	if err != nil {
		return nil, err
	}
	return data, err
}

// получает конфиг
func getConfig() (*oauth2.Config, error) {
	credentials, err := getCredentials()
	if err != nil {
		log.Fatalf("Не удалось загрузить файл учетных данных: %v", err)
	}
	config, err := google.ConfigFromJSON(credentials, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, err
	}
	return config, err
}

// получает токен из файла
func getTokenFromFile() (*oauth2.Token, error) {
	f, err := os.Open(tokenFileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token, err
}

// получает токен у гугла
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	// Создаем URL для авторизации
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Перейдите по следующему URL для авторизации: %v\n", authURL)

	// Получаем код авторизации от пользователя
	var code string
	fmt.Print("Введите код авторизации: ")
	fmt.Scan(&code)

	// Обмениваем код на токен
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Не удалось обменять код: %v", err)
	}
	return token
}

// сохраняет токен в файл
func saveToken(token *oauth2.Token) {
	file, err := os.OpenFile(tokenFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Не удалось открыть файл для записи токена: %v", err)
	}
	defer file.Close()
	json.NewEncoder(file).Encode(token)
}
