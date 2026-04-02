package analysis

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/VysMax/analyze-cfg/models"
)

func AnalyseDirectory(root string, cfg *models.Config) (string, error) {
	var allowedExts = map[string]struct{}{
		".json": {},
		".yaml": {},
		".yml":  {},
	}

	messages := make([]string, 0)

	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {

		if err != nil {
			return nil
		}

		ext := filepath.Ext(path)
		_, isAllowed := allowedExts[ext]
		if !isAllowed {
			return nil
		}

		cfg.ConfigPath = path

		message, err := AnalyseFile(cfg)
		if err != nil {
			return err
		}

		messages = append(messages, message)

		return nil
	})

	return strings.Join(messages, "\n"), err
}
