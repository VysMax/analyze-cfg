package analysis

import (
	"fmt"
	"os"
	"strings"

	"github.com/VysMax/analyze-cfg/models"
)

type Analyzer struct {
	problems []models.Problem
}

func NewAnalyzer() *Analyzer {
	p := make([]models.Problem, 0)
	return &Analyzer{problems: p}
}

func (a *Analyzer) AnalyzeCfg(cfg *models.Config) ([]models.Problem, error) {
	a.CheckHost(cfg)
	a.CheckPassword(cfg)
	a.CheckTLS(cfg)
	a.CheckLogLevel(cfg)
	a.CheckAlgorithm(cfg)
	a.CheckPermissionsToSet(cfg)

	if cfg.Storage.Path != "" {
		if err := a.CheckCurrentPermissions(cfg); err != nil {
			return nil, err
		}
	}

	if cfg.File != "" {
		if err := a.CheckConfigFilePermissions(cfg); err != nil {
			return nil, err
		}
	}

	return a.problems, nil
}

func (a *Analyzer) CheckHost(cfg *models.Config) {
	if cfg.Host == "0.0.0.0" {
		a.problems = append(a.problems, models.Problem{
			Filename:       cfg.File,
			Path:           "server.host",
			Description:    "Сервер слушает на всех хостах (0.0.0.0).",
			Recommendation: "Ограничьте доступ. Для локального доступа используйте 127.0.0.1 .",
			Severity:       "HIGH",
		})
	}
}

func (a *Analyzer) CheckPassword(cfg *models.Config) {
	hiddenPasswordSymbols := map[string]string{
		"${": "}",
		"$(": ")",
	}

	if cfg.Database.Password != "" {
		seemsSecure := false
		for key, symbol := range hiddenPasswordSymbols {
			if strings.HasPrefix(cfg.Database.Password, key) && strings.HasSuffix(cfg.Database.Password, symbol) {
				seemsSecure = true
				break
			}
		}

		if !seemsSecure {
			a.problems = append(a.problems, models.Problem{
				Filename:       cfg.File,
				Path:           "database.password",
				Description:    "Пароль в открытом виде.",
				Recommendation: "Используйте переменные окружения.",
				Severity:       "HIGH",
			})
		}
	}
}

func (a *Analyzer) CheckTLS(cfg *models.Config) {
	if cfg.Server.TlsVerify != nil && *cfg.Server.TlsVerify == false {
		a.problems = append(a.problems, models.Problem{
			Filename:       cfg.File,
			Path:           "server.tls_verify",
			Description:    "TLS проверка выключена.",
			Recommendation: "Включите TLS-проверку.",
			Severity:       "HIGH",
		})
	}
}

func (a *Analyzer) CheckLogLevel(cfg *models.Config) {
	if cfg.Log.Level == "debug" {
		a.problems = append(a.problems, models.Problem{
			Filename:       cfg.File,
			Path:           "log.level",
			Description:    "логирование в debug-режиме.",
			Recommendation: "Поменяйте режим на более избирательный (info+).",
			Severity:       "LOW",
		})
	}
}

func (a *Analyzer) CheckAlgorithm(cfg *models.Config) {
	insecureAlgorithms := map[string]string{
		"MD5":      "слишком слабый алгоритм -",
		"sha1":     "устаревший алгоритм -",
		"des":      "устаревший алгоритм -",
		"3des":     "устаревший алгоритм -",
		"rc4":      "небезопасный алгоритм -",
		"blowfish": "устаревший алгоритм -",
		"tlsv1.0":  "устаревший алгоритм -",
		"tlsv1.1":  "устаревший алгоритм -",
		"sslv3":    "небезопасный алгоритм -",
	}

	if cfg.DigestAlgorithm != "" {
		reason, ok := insecureAlgorithms[cfg.DigestAlgorithm]
		if ok {
			a.problems = append(a.problems, models.Problem{
				Filename:       cfg.File,
				Path:           "cfg.digest-algorithm",
				Description:    fmt.Sprintf("%s %s.", reason, cfg.DigestAlgorithm),
				Recommendation: "Замените его на более безопасный.",
				Severity:       "HIGH",
			})
		}
	}
}

func (a *Analyzer) CheckPermissionsToSet(cfg *models.Config) {
	if cfg.Permissions == "0777" || cfg.Permissions == "777" {
		a.problems = append(a.problems, models.Problem{
			Filename:       cfg.File,
			Path:           "storage.permissions",
			Description:    "слишком широкие права доступа.",
			Recommendation: "Поменяйте на менее широкие (например, 750 или 700).",
			Severity:       "HIGH",
		})
	}
}

func (a *Analyzer) CheckCurrentPermissions(cfg *models.Config) error {

	info, err := os.Stat(cfg.Storage.Path)
	if err != nil {
		return fmt.Errorf("ошибка открытия файла. Возможно, файл не существует")
	}

	mode := info.Mode().Perm()

	if mode&0004 != 0 {
		a.problems = append(a.problems, models.Problem{
			Filename:       cfg.Storage.Path,
			Path:           "storage.path",
			Description:    fmt.Sprintf("слишком широкие права доступа для файла %s: другие могут читать файл.", cfg.Storage.Path),
			Recommendation: "Поменяйте права доступа на менее широкие (например, 750 или 700).",
			Severity:       "MEDIUM",
		})
	}

	if mode&0002 != 0 {
		a.problems = append(a.problems, models.Problem{
			Filename:       cfg.Storage.Path,
			Path:           "storage.path",
			Description:    fmt.Sprintf("слишком широкие права доступа для файла %s: другие могут писать в файл.", cfg.Storage.Path),
			Recommendation: "Поменяйте права доступа на менее широкие (например, 750 или 700).",
			Severity:       "HIGH",
		})
	}

	if mode&0020 != 0 {
		a.problems = append(a.problems, models.Problem{
			Filename:       cfg.File,
			Path:           "storage.path",
			Description:    fmt.Sprintf("слишком широкие права доступа для файла %s: группа может писать в файл.", cfg.Storage.Path),
			Recommendation: "Поменяйте права доступа на менее широкие (например, 750 или 700).",
			Severity:       "MEDIUM",
		})
	}

	return nil
}

func (a *Analyzer) CheckConfigFilePermissions(cfg *models.Config) error {
	info, err := os.Stat(cfg.File)
	if err != nil {
		return fmt.Errorf("ошибка открытия файла. Возможно, файл не существует")
	}

	mode := info.Mode().Perm()

	if mode&0002 != 0 {
		a.problems = append(a.problems, models.Problem{
			Filename:       cfg.File,
			Path:           "entire config",
			Description:    "слишком широкие права доступа: другие могут писать в файл конфигурации.",
			Recommendation: "Поменяйте права доступа на менее широкие (например, 644 или 600).",
			Severity:       "HIGH",
		})
	}

	if mode&0020 != 0 {
		a.problems = append(a.problems, models.Problem{
			Filename:       cfg.File,
			Path:           "entire config",
			Description:    "слишком широкие права доступа: группа может писать в файл конфигурации.",
			Recommendation: "Поменяйте права доступа на менее широкие (например, 644 или 600).",
			Severity:       "MEDIUM",
		})
	}

	return nil
}
