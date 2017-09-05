package helper

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

var translation string

// T : translation
func T(key string) string {
	lang := "en"
	loadTranslations(lang)
	key = lang + "." + key
	val := viper.GetString(key)
	val = strings.Replace(val, "\n", "\n   ", -1)

	return val
}

func loadTranslations(lang string) {
	if translation == "" {
		file := "lang/" + lang + ".yml"
		data, err := Asset(file)
		if err != nil {
			fmt.Println("error: ", err)
			return
		}

		viper.SetConfigType("yaml")
		_ = viper.ReadConfig(bytes.NewBuffer(data))
		translation = lang
	}
}
