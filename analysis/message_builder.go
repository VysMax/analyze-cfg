package analysis

import (
	"fmt"
	"strings"
)

func MessageBuilder(filepath string, problems Problems) string {
	var sb strings.Builder

	if filepath != "" {
		sb.WriteString(fmt.Sprintf("Файл %s:\n\n", filepath))
	}

	if len(problems) == 0 {
		sb.WriteString("Проблем не обнаружено\n")
		return sb.String()
	}

	switch {
	case len(problems)%10 == 1 && len(problems)%100 != 11:
		sb.WriteString(fmt.Sprintf("Обнаружена %d потенциально опасная настройка:\n\n", len(problems)))
	case len(problems)%10 >= 2 && len(problems)%10 <= 4 && (len(problems)%100)/10 != 1:
		sb.WriteString(fmt.Sprintf("Обнаружено %d потенциально опасные настройки:\n\n", len(problems)))
	default:
		sb.WriteString(fmt.Sprintf("Обнаружено %d потенциально опасных настроек:\n\n", len(problems)))
	}

	for i, problem := range problems {
		sb.WriteString(fmt.Sprintf("%d. %s: %s %s\n\n", i+1, problem.Severity, problem.Description, problem.Recommendation))
	}

	return sb.String()
}
