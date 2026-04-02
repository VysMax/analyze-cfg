package analysis

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/VysMax/analyze-cfg/models"
	"gopkg.in/yaml.v3"
)

func AnalyseFile(cfg *models.Config) (string, error) {

	file, err := os.Open(cfg.ConfigPath)
	if err != nil {
		return "", fmt.Errorf("ошибка открытия файла %s: %v\n", cfg.ConfigPath, err)
	}
	defer file.Close()

	var r io.Reader
	r = file
	err = SetReader(&r)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения из стандартного ввода: %v\n", err)
	}

	format := path.Ext(cfg.ConfigPath)

	switch format {
	case ".json":
		err = json.NewDecoder(r).Decode(&cfg)
	case ".yaml", ".yml":
		err = yaml.NewDecoder(r).Decode(&cfg)
	default:
		return "", fmt.Errorf("неподдерживаемый формат:%s\n", format)
	}

	if err != nil && err != io.EOF {
		return "", fmt.Errorf("ошибка парсинга конфигурации: %v\n", err)
	}

	var problems Problems
	if err = problems.AnalyseCfg(cfg); err != nil {
		return "", fmt.Errorf("ошибка проверки конфигурации: %v\n", err)
	}

	message := MessageBuilder(cfg.ConfigPath, problems)

	return message, nil
}
