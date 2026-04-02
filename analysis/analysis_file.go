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

func AnalyseFile(cfg *models.Config) (Problems, error) {

	file, err := os.Open(cfg.ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия файла %s: %v\n", cfg.ConfigPath, err)
	}
	defer file.Close()

	var r io.Reader
	r = file
	err = SetReader(&r)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения из стандартного ввода: %v\n", err)
	}

	format := path.Ext(cfg.ConfigPath)

	switch format {
	case ".json":
		err = json.NewDecoder(r).Decode(&cfg)
	case ".yaml", ".yml":
		err = yaml.NewDecoder(r).Decode(&cfg)
	default:
		return nil, fmt.Errorf("неподдерживаемый формат:%s\n", format)
	}

	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("ошибка парсинга конфигурации: %v\n", err)
	}

	var problems Problems
	if err = problems.AnalyseCfg(cfg); err != nil {
		return nil, fmt.Errorf("ошибка проверки конфигурации: %v\n", err)
	}

	return problems, nil
}
