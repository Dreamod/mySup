package gSheets

import (
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/api/sheets/v4"
)

var service = GetService()

// GetUsers получает список участников
func GetUsers() ([]string, error) {
	sheetRange := os.Getenv("DICTIONARY_LIST_RANGE") + "!A2:A"
	data, err := getData(sheetRange)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, value := range data {
		result = append(result, fmt.Sprintf("%s", value[0]))
	}
	return result, nil
}

// GetRivers получает список участников
func GetRivers() ([]string, error) {
	sheetRange := os.Getenv("DICTIONARY_LIST_RANGE") + "!C2:C"
	data, err := getData(sheetRange)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, value := range data {
		result = append(result, fmt.Sprintf("%s", value[0]))
	}
	return result, nil
}

// GetStatistic получает статистику
func GetStatistic() ([][]interface{}, error) {
	sheetRange := os.Getenv("RESULT_LIST_NAME") + "!A2:E"
	return getData(sheetRange)
}

// получает информацию из таблицы
func getData(readRange string) ([][]interface{}, error) {
	response, err := service.Spreadsheets.Values.Get(getSheetId(), readRange).Do()
	if err != nil {
		if strings.Contains(err.Error(), "Token has been expired") || strings.Contains(err.Error(), "Error 401: Request had invalid authentication credentials") {
			DeleteTokenFile()
			log.Fatalf("Токен доступа истек или невалиден, приложение необходимо перезапустить")
		} else {
			log.Fatalf("Не удалось получить данные: %v", err)
		}
	}
	return response.Values, nil
}

// получает ID таблицы
func getSheetId() string {
	sheetId := os.Getenv("SHEET_ID")
	if sheetId == "" {
		log.Fatalf("Не получен ID таблицы")
	}
	return sheetId
}

// WriteData записывает данные в таблицу
func WriteData(writeRange string, values [][]interface{}) {
	var vr sheets.ValueRange
	vr.Values = values
	_, err := service.Spreadsheets.Values.Append(getSheetId(), writeRange, &vr).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Fatalf("Не удалось записать данные: %v", err)
	}
	fmt.Println("Данные успешно записаны.")
}
