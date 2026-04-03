package analysis

import (
	"fmt"
	"os"
	"strings"

	"github.com/VysMax/analyze-cfg/models"
)

type Problems []models.Problem

func (p *Problems) AnalyseCfg(cfg *models.Config) error {
	p.CheckHost(cfg)
	p.CheckPassword(cfg)
	p.CheckTLS(cfg)
	p.CheckLogLevel(cfg)
	p.CheckAlgorithm(cfg)
	p.CheckPermissionsToSet(cfg)

	if cfg.Storage.Path != "" {
		if err := p.CheckCurrentPermissions(cfg); err != nil {
			return err
		}
	}

	if cfg.Path != "" {
		if err := p.CheckConfigFilePermissions(cfg); err != nil {
			return err
		}
	}

	return nil
}

func (p *Problems) CheckHost(cfg *models.Config) {
	if cfg.Server.Host == "0.0.0.0" {
		*p = append(*p, models.Problem{
			Filename:       cfg.Path,
			Path:           "server.host",
			Description:    "Сервер слушает на всех хостах (0.0.0.0).",
			Recommendation: "Ограничьте доступ. Для локального доступа используйте 127.0.0.1 .",
			Severity:       "HIGH",
		})
	}
}

func (p *Problems) CheckPassword(cfg *models.Config) {
	passwordMarkers := []string{"${", "$("}

	if cfg.Database.Password != "" {
		isSecure := false
		for _, marker := range passwordMarkers {
			if strings.HasPrefix(cfg.Database.Password, marker) {
				isSecure = true
				break
			}
		}

		if !isSecure {
			*p = append(*p, models.Problem{
				Filename:       cfg.Path,
				Path:           "database.password",
				Description:    "Пароль в открытом виде.",
				Recommendation: "Используйте переменные окружения.",
				Severity:       "HIGH",
			})
		}
	}
}

func (p *Problems) CheckTLS(cfg *models.Config) {
	if !cfg.Server.TlsVerify {
		*p = append(*p, models.Problem{
			Filename:       cfg.Path,
			Path:           "server.tls_verify",
			Description:    "TLS проверка выключена.",
			Recommendation: "Включите TLS-проверку.",
			Severity:       "HIGH",
		})
	}
}

func (p *Problems) CheckLogLevel(cfg *models.Config) {
	if cfg.Log.Level == "debug" {
		*p = append(*p, models.Problem{
			Filename:       cfg.Path,
			Path:           "log.level",
			Description:    "логирование в debug-режиме.",
			Recommendation: "Поменяйте режим на более избирательный (info+).",
			Severity:       "LOW",
		})
	}
}

func (p *Problems) CheckAlgorithm(cfg *models.Config) {
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
			*p = append(*p, models.Problem{
				Filename:       cfg.Path,
				Path:           "cfg.digest-algorithm",
				Description:    fmt.Sprintf("%s %s.", reason, cfg.DigestAlgorithm),
				Recommendation: "Замените его на более безопасный.",
				Severity:       "HIGH",
			})
		}
	}
}

func (p *Problems) CheckPermissionsToSet(cfg *models.Config) {
	if cfg.Permissions == "0777" || cfg.Permissions == "777" {
		*p = append(*p, models.Problem{
			Filename:       cfg.Path,
			Path:           "storage.permissions",
			Description:    "слишком широкие права доступа.",
			Recommendation: "Поменяйте на менее широкие (например, 750 или 700).",
			Severity:       "HIGH",
		})
	}
}

func (p *Problems) CheckCurrentPermissions(cfg *models.Config) error {

	info, err := os.Stat(cfg.Storage.Path)
	if err != nil {
		return fmt.Errorf("ошибка открытия файла. Возможно, файл не существует")
	}

	mode := info.Mode().Perm()

	if mode&0004 != 0 {
		*p = append(*p, models.Problem{
			Filename:       cfg.Storage.Path,
			Path:           "storage.path",
			Description:    fmt.Sprintf("слишком широкие права доступа для файла %s: другие могут читать файл.", cfg.Storage.Path),
			Recommendation: "Поменяйте права доступа на менее широкие (например, 750 или 700).",
			Severity:       "MEDIUM",
		})
	}

	if mode&0002 != 0 {
		*p = append(*p, models.Problem{
			Filename:       cfg.Path,
			Path:           "storage.path",
			Description:    fmt.Sprintf("слишком широкие права доступа для файла %s: другие могут писать в файл.", cfg.Storage.Path),
			Recommendation: "Поменяйте права доступа на менее широкие (например, 750 или 700).",
			Severity:       "HIGH",
		})
	}

	if mode&0020 != 0 {
		*p = append(*p, models.Problem{
			Filename:       cfg.Path,
			Path:           "storage.path",
			Description:    fmt.Sprintf("слишком широкие права доступа для файла %s: группа может писать в файл.", cfg.Storage.Path),
			Recommendation: "Поменяйте права доступа на менее широкие (например, 750 или 700).",
			Severity:       "MEDIUM",
		})
	}

	return nil
}

func (p *Problems) CheckConfigFilePermissions(cfg *models.Config) error {
	info, err := os.Stat(cfg.Path)
	if err != nil {
		return fmt.Errorf("ошибка открытия файла. Возможно, файл не существует")
	}

	mode := info.Mode().Perm()

	if mode&0002 != 0 {
		*p = append(*p, models.Problem{
			Filename:       cfg.Path,
			Path:           "entire config",
			Description:    "слишком широкие права доступа: другие могут писать в файл конфигурации.",
			Recommendation: "Поменяйте права доступа на менее широкие (например, 644 или 600).",
			Severity:       "HIGH",
		})
	}

	if mode&0020 != 0 {
		*p = append(*p, models.Problem{
			Filename:       cfg.Path,
			Path:           "entire config",
			Description:    "слишком широкие права доступа: группа может писать в файл конфигурации.",
			Recommendation: "Поменяйте права доступа на менее широкие (например, 644 или 600).",
			Severity:       "MEDIUM",
		})
	}

	return nil
}
