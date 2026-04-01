package analysis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/VysMax/analyze-cfg/models"
	"gopkg.in/yaml.v3"
)

func DecodeFromStdin(r io.Reader, cfg *models.Config) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("ошибка чтения данных из стандартного ввода: %w", err)
	}

	if len(data) == 0 {
		return fmt.Errorf("отсутствуют данные в стандартном вводе")
	}

	var eachAttemptErrInfo []string

	if err = json.NewDecoder(bytes.NewReader(data)).Decode(&cfg); err != nil {
		eachAttemptErrInfo = append(eachAttemptErrInfo, fmt.Sprintf("json: %v", err))
	} else {
		return nil
	}

	if err = yaml.NewDecoder(bytes.NewReader(data)).Decode(&cfg); err != nil {
		eachAttemptErrInfo = append(eachAttemptErrInfo, fmt.Sprintf("; %v", err))
	} else {
		return nil
	}

	return fmt.Errorf("%s", eachAttemptErrInfo)
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
