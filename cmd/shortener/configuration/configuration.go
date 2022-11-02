package configuration

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/caarlos0/env"
)

const (
	DefaultFileName      = "urls.log"
	DefaultFilePerm      = 0755
	DefaultServerAddress = "localhost:8080"
	DefaultBaseURL       = "http://localhost:8080/"
)

type config struct {
	ServerAdress string `env:"SERVER_ADDRESS"`
	BaseURL      string `env:"BASE_URL"`
	FilePath     string `env:"FILE_STORAGE_PATH"`
}

func New() *config {
	cfg := config{
		ServerAdress: DefaultServerAddress,
		FilePath:     DefaultFileName,
		BaseURL:      DefaultBaseURL,
	}
	cfg.BaseURL = fmt.Sprintf("http://%s/", cfg.ServerAdress)

	err := env.Parse(&cfg)

	if err != nil {
		log.Fatal(err)
	}
	flagServerAdress := flag.String("a", DefaultServerAddress, "server adress")
	flagBaseURL := flag.String("b", DefaultBaseURL, "base url")
	flagFilePath := flag.String("f", DefaultFileName, "file path")
	flag.Parse()

	if *flagServerAdress != DefaultServerAddress {
		cfg.ServerAdress = *flagServerAdress
	}
	if *flagBaseURL != DefaultBaseURL {
		cfg.BaseURL = *flagBaseURL
	}
	if *flagFilePath != DefaultFileName {
		cfg.FilePath = *flagFilePath
	}

	if cfg.FilePath != DefaultFileName {
		if _, err := os.Stat(filepath.Dir(cfg.FilePath)); os.IsNotExist(err) {
			log.Println("Creating folder")
			err := os.Mkdir(filepath.Dir(cfg.FilePath), DefaultFilePerm)
			if err != nil {
				log.Printf("Error: %v \n", err)
			}
		}
	}

	if string(cfg.BaseURL[len(cfg.BaseURL)-1]) != "/" {
		cfg.BaseURL += "/"
	}

	return &cfg
}
