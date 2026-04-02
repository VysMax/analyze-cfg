package analysis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/VysMax/analyze-cfg/models"
	"gopkg.in/yaml.v3"
)

func AnalysisStdin(input *os.File, cfg *models.Config) (string, error) {
	var r io.Reader
	r = input
	err := SetReader(&r)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения из стандартного ввода: %v\n", err)
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения данных из стандартного ввода: %w", err)
	}

	if len(data) == 0 {
		return "", fmt.Errorf("отсутствуют данные в стандартном вводе")
	}

	var (
		eachAttemptErrInfo []string
		isParsed           bool
	)

	err = json.NewDecoder(bytes.NewReader(data)).Decode(&cfg)
	if err != nil {
		eachAttemptErrInfo = append(eachAttemptErrInfo, fmt.Sprintf("json: %v", err))
	} else {
		isParsed = true
	}

	if !isParsed {
		err = yaml.NewDecoder(bytes.NewReader(data)).Decode(&cfg)
		if err != nil {
			eachAttemptErrInfo = append(eachAttemptErrInfo, fmt.Sprintf("; %v", err))
		} else {
			isParsed = true
		}
	}

	if !isParsed {
		return "", fmt.Errorf("%s", eachAttemptErrInfo)
	}

	var problems Problems
	if err = problems.AnalyseCfg(cfg); err != nil {
		return "", fmt.Errorf("ошибка проверки конфигурации: %v\n", err)
	}

	message := MessageBuilder("", problems)

	return message, nil
}

func SetReader(r *io.Reader) error {
	data := make([]byte, 1024)
	n, err := (*r).Read(data)
	if err != nil && err != io.EOF {
		return err
	}

	*r = io.MultiReader(strings.NewReader(string(data[:n])), *r)
	return nil
}
