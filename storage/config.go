package storage

type (
	Type   string
	Config struct {
		Type     Type   `json:"type"`
		Endpoint string `json:"endpoint"`
		AccessId string `json:"access_id"`
		Secret   string `json:"secret"`
		SSL      bool   `json:"ssl"`
		Bucket   string `json:"bucket"`
	}
)

const (
	TypeMinio Type = "minio"
)
