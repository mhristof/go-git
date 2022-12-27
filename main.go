package git

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Root(pwd string) string {
	var root string

	err := filepath.Walk(pwd,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info, err := os.Stat(path); err == nil && !info.IsDir() {
				return nil
			}

			dotGit := filepath.Join(path, ".git")
			if _, err := os.Stat(dotGit); !os.IsNotExist(err) {
				root = path
			}

			return nil
		})
	if err != nil {
		log.Println(err)
	}

	return root
}

func Origin(pwd string) string {
	config := filepath.Join(Root(pwd), ".git/config")

	viper.SetConfigFile(config)
	viper.SetConfigType("ini")

	err := viper.MergeInConfig()
	if err != nil {
		log.WithFields(log.Fields{
			"err":    err,
			"config": config,
		}).Panic("cannot read git config")
	}

	return viper.GetStringMapString(`remote "origin"`)["url"]
}
