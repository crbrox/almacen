package almacen

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	Address  string
	MongoURL string
}

func LoadConfig(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return Load(f)
}

func Load(reader io.Reader) (*Config, error) {
	var c = &Config{}
	err := json.NewDecoder(reader).Decode(c)
	if err != nil {
		return nil, err
	}

	// additional validation

	return c, nil
}
