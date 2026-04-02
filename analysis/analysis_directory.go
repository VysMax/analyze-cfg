package analysis

import (
	"io/fs"
	"path/filepath"

	"github.com/VysMax/analyze-cfg/models"
)

func AnalyseDirectory(root string, cfg *models.Config) ([]Problems, error) {
	var allowedExts = map[string]struct{}{
		".json": {},
		".yaml": {},
		".yml":  {},
	}

	// messages := make([]string, 1)
	// messages[0] = fmt.Sprintf("Анализ папки %s:\n", root)

	var multFileProblems []Problems

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

		problems, err := AnalyseFile(cfg)
		if err != nil {
			return err
		}

		multFileProblems = append(multFileProblems, problems)

		return nil
	})

	return multFileProblems, err
	// if len(messages) == 1 {
	// 	return fmt.Sprintf("Папка %s не содержит файлов конфигурации\n", root), nil
	// }

	// return strings.Join(messages, "\n"), err
}
