package models

type Config struct {
	File     string `json:"path" yaml:"path"`
	Server   `json:"server" yaml:"server"`
	Database `json:"database" yaml:"database"`
	Storage  `json:"storage" yaml:"storage"`
	Log      `json:"log" yaml:"log"`
}

type Server struct {
	Host      string `json:"host" yaml:"host"`
	TlsVerify bool   `json:"tls_verify,omitempty" yaml:"tls_verify"`
}

type Database struct {
	Password string `json:"password" yaml:"password"`
}

type Storage struct {
	Path            string `json:"path" yaml:"path"`
	Permissions     string `json:"permissions" yaml:"permissions"`
	DigestAlgorithm string `json:"digest_algorithm" yaml:"digest_algorithm"`
}

type Log struct {
	// Output string `json:"output" yaml:"output"`
	Level string `json:"level" yaml:"level"`
}

type Problem struct {
	Filename       string `json:"filename,omitempty"`
	Path           string `json:"path"`
	Description    string `json:"description"`
	Recommendation string `json:"recommendation"`
	Severity       string `json:"severity"`
}
