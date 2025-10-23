package main

import (
	"MySup/gSheets"
	"MySup/output"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
)

//TODO - написать юнит тесты

var menu = map[string]func(){
	"1": addResult,
	"2": getUserData,
	"3": getYearData,
	"4": addUser,
	"5": addRiver,
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Не удалось загрузить .env файл")
	}
	fmt.Println("<-- Свой в доску -->")
menu:
	for {
		action := promptData(
			"Меню:",
			"1: Добавить запись",
			"-----------------------",
			"2: Данные по участнику",
			"3: Данные за год",
			"-----------------------",
			"4: Добавить участника",
			"5: Добавить реку",
			"-----------------------",
			"0: Выход",
			"Ввод",
		)
		menuFunc := menu[action]
		if action == "0" {
			break menu
		}
		if menuFunc == nil && action != "0" {
			output.PrintError("Неверный ввод, выбери что-то из меню")
			continue
		}
		menuFunc()
		if !needMenu() {
			break menu
		}
	}
	fmt.Println("\nДо встречи =)")
}

// запрашивает необходимость продолжения работы с приложением
func needMenu() bool {
	input := promptData("Продолжить? (Y:n)")
	return input == "Y" || input == "y"
}

// запрашивает ввод данных
func promptData(prompt ...any) string {
	input := ""
	for i, value := range prompt {
		if i+1 != len(prompt) {
			fmt.Println(value)
		} else {
			fmt.Printf("%s: ", any(value))
		}
	}
	fmt.Scanln(&input)
	return input
}

// добавляет результат в таблицу
func addResult() {
	values := [][]interface{}{
		{
			time.Now().Format("02.01.2006 15:04:05"),
			getUser(),
			getRiver(),
			getDate(),
			getDistance(),
		},
	}
	sheetRange := os.Getenv("RESULT_LIST_NAME") + "!A:E"
	gSheets.WriteData(sheetRange, values)
}

// получает данные по году
func getYearData() {
	year := getYear()
	data, err := gSheets.GetStatistic()
	if err != nil {
		output.PrintError("Чет пошло не так")
	}
	var resultMap map[string]int
	resultMap = make(map[string]int)
	// соберем данные за год
	// участник - километраж
	for _, value := range data {
		date, _ := time.Parse("02.01.2006", fmt.Sprint(value[3]))
		if date.Year() == year {
			userName := fmt.Sprint(value[1])
			distance, _ := strconv.Atoi(fmt.Sprint(value[4]))
			if resultMap[userName] > 0 {
				resultMap[userName] = resultMap[userName] + distance
			} else {
				resultMap[userName] = distance
			}
		}
	}

	// отсортируем данные
	// > километраж
	if len(resultMap) > 0 {
		type resultType struct {
			Name     string
			Distance int
		}
		var result []resultType
		for k, v := range resultMap {
			result = append(result, resultType{k, v})
		}
		sort.Slice(result, func(i, j int) bool {
			return result[i].Distance > result[j].Distance
		})

		// вывести данные участник - километраж
		fmt.Println("Статистика за", year, "год")
		for _, value := range result {
			fmt.Println(value.Name, ":", value.Distance, "км")
		}
	} else {
		output.PrintError(fmt.Sprintf("Статистики за %d еще нет", year))
	}
}

// получает данные по юзеру
func getUserData() {
	user := getUser()
	data, err := gSheets.GetStatistic()
	if err != nil {
		output.PrintError("Чет пошло не так")
	}
	var resultMap map[int]int
	resultMap = make(map[int]int)
	// соберем данные участника
	// год - километраж
	for _, value := range data {
		date, _ := time.Parse("02.01.2006", fmt.Sprint(value[3]))
		year := date.Year()
		if fmt.Sprint(value[1]) == user {
			distance, _ := strconv.Atoi(fmt.Sprint(value[4]))
			if resultMap[year] > 0 {
				resultMap[year] = resultMap[year] + distance
			} else {
				resultMap[year] = distance
			}
		}
	}

	// отсортируем данные
	// год по убыванию
	if len(resultMap) > 0 {
		type resultType struct {
			Year     int
			Distance int
		}
		var result []resultType
		for y, d := range resultMap {
			result = append(result, resultType{y, d})
		}
		sort.Slice(result, func(i, j int) bool {
			return result[i].Year > result[j].Year
		})

		// вывести данные год - километраж
		fmt.Println("Статистика", user)
		for _, value := range result {
			fmt.Println(value.Year, ":", value.Distance, "км")
		}
	} else {
		output.PrintError(fmt.Sprintf("Статистики для %s еще нет", user))
	}
}

// добавляет участника
func addUser() {
	values := [][]interface{}{
		{
			getUserName(),
		},
	}
	sheetRange := os.Getenv("DICTIONARY_LIST_RANGE") + "!A:A"
	gSheets.WriteData(sheetRange, values)
}

// добавляет реку
func addRiver() {
	values := [][]interface{}{
		{
			getRiverName(),
		},
	}
	sheetRange := os.Getenv("DICTIONARY_LIST_RANGE") + "!C:C"
	gSheets.WriteData(sheetRange, values)
}

// запрашивает год
func getYear() int {
	re := regexp.MustCompile(`\d\d\d\d`)
	result := 0
	for {
		year := promptData("Год [YYYY]")
		if year == "" {
			result = time.Now().Year()
			break
		} else {
			if re.MatchString(year) {
				intYear, _ := strconv.Atoi(year)
				if intYear < 2023 {
					output.PrintError("Статистики до 2023 нет")
					continue
				}
				if intYear > time.Now().Year() {
					output.PrintError("Статистики из будущего нет")
					continue
				}
				result = intYear
				break
			} else {
				output.PrintError("Неверный формат года")
				continue
			}
		}
	}
	return result
}

// запрашивает участника
func getUser() string {
	// получим список участников
	users, err := gSheets.GetUsers()
	if err != nil {
		output.PrintError("Чет пошло не так")
	}
	result := ""
	// преобразуем список участников
	var promptUsers = []any{"Участник"}
	for i, value := range users {
		promptUsers = append(promptUsers, fmt.Sprintf("%d: %s", i+1, value))
	}
	// запросим участника
	for {
		user := promptData(promptUsers...)
		userIndex, _ := strconv.Atoi(user)
		if userIndex > len(users) || userIndex < 1 {
			output.PrintError("Неверный ввод, выбери из списка")
			continue
		} else {
			result = users[userIndex-1]
			break
		}
	}
	return result
}

// запрашивает реку
func getRiver() string {
	// получим список рек
	rivers, err := gSheets.GetRivers()
	if err != nil {
		output.PrintError("Чет пошло не так")
	}
	result := ""
	// преобразуем список рек
	var promptRivers = []any{"Река"}
	for i, value := range rivers {
		promptRivers = append(promptRivers, fmt.Sprintf("%d: %s", i+1, value))
	}
	// запросим реку
	for {
		river := promptData(promptRivers...)
		riverIndex, _ := strconv.Atoi(river)
		if riverIndex > len(rivers) {
			output.PrintError("Неверный ввод, выбери из списка")
			continue
		} else {
			result = rivers[riverIndex-1]
			break
		}
	}
	return result
}

// запрашивает дату
func getDate() string {
	re := regexp.MustCompile(`\d\d.\d\d.\d\d\d\d`)
	result := ""
	for {
		date := promptData("Дата [dd.mm.YYYY]")
		if date == "" {
			result = time.Now().Format("02.01.2006")
			break
		} else {
			if re.MatchString(date) {
				result = date
				break
			} else {
				output.PrintError("Неверный формат даты")
				continue
			}
		}
	}
	return result
}

// запрашивает дистанцию
func getDistance() int {
	result := 0
	for {
		value := promptData("Километры")
		if intValue, err := strconv.Atoi(value); err == nil {
			result = intValue
			break
		} else {
			output.PrintError("Нужно ввести целое число")
			continue
		}
	}
	return result
}

// запрашивает имя нового участника
func getUserName() string {
	users, err := gSheets.GetUsers()
	if err != nil {
		output.PrintError("Чет пошло не так")
	}
	result := ""
	for {
		value := promptData("Имя участника [Дима И.]")
		if value == "" {
			output.PrintError("Поле обязательно для заполнения")
			continue
		} else {
			if !itemExistInSlice(value, users) {
				result = value
				break
			} else {
				output.PrintError("Такой участник уже есть")
				continue
			}
		}
	}
	return result
}

// запрашивает имя новой реки
func getRiverName() string {
	rivers, err := gSheets.GetRivers()
	if err != nil {
		output.PrintError("Чет пошло не так")
	}
	result := ""
	for {
		value := promptData("Имя реки [Клязьма]")
		if value == "" {
			output.PrintError("Поле обязательно для заполнения")
			continue
		} else {
			if !itemExistInSlice(value, rivers) {
				result = value
				break
			} else {
				output.PrintError("Такая река уже есть")
				continue
			}
		}
	}
	return result
}

// проверяет есть ли в слайсе строк
func itemExistInSlice(string string, slice []string) bool {
	for _, value := range slice {
		if value == string {
			return true
		}
	}
	return false
}
