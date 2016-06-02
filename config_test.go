package almacen

import "testing"

func TestLoadConfig(t *testing.T) {
	const (
		mongoURL = "mongo url"
		address  = "testing address"
	)
	c, err := LoadConfig("testdata/config.json")
	if err != nil {
		t.Fatal(err)
	}
	if c.Address != address {
		t.Errorf("config address: wanted %v, got %v", address, c.Address)
	}
	if c.MongoURL != mongoURL {
		t.Errorf("config mongo url: wanted %v, got %v", mongoURL, c.MongoURL)
	}
}
func TestLoadConfigErrOpen(t *testing.T) {
	_, err := LoadConfig("testdata/not_existing_file")
	if err == nil {
		t.Error("not existent file: wanted error, got nil")
	}
}

func TestLoadConfigErrDecode(t *testing.T) {
	_, err := LoadConfig("testdata/invalid_config.txt")
	if err == nil {
		t.Error("not valid json: wanted error, got nil")
	}
}
