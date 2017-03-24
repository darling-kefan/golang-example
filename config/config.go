package config

import (
	"log"
	"strings"
	"io/ioutil"
)

var cfg map[string]string

func init() {
	cfg = make(map[string]string)

	// How to get project root directory?
	envFile := "/home/shouqiang/codes/Go/publish-profile/src/.env"
	content, err := ioutil.ReadFile(envFile)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		lineParts := strings.Split(line, "=")
		if len(lineParts) != 2 {
			continue
		}
		cfg[strings.TrimSpace(lineParts[0])] = strings.TrimSpace(lineParts[1])
	}
}

func Get(key string, defVal string) string {
	if v, ok := cfg[key]; ok {
		return v
	}
	return defVal
}

func GetAll() map[string]string {
	return cfg
}
