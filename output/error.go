package output

import "github.com/fatih/color"

func PrintError(value any) {

	switch valueWithType := value.(type) {

	case string:
		color.Red(valueWithType)
	case error:
		color.Red(valueWithType.Error())
	case int:
		color.Yellow("Код ошибки: %d", valueWithType)
	default:
		color.Cyan("Неизвестный тип ошибки")
	}
}
