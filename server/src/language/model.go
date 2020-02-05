package language

import (
	"qoj/server/config"
	"strings"
)

type Language struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Extension string `json:"-"`
	Command   string `json:"-"`
}

func (lan *Language) normaliseLanguage() {
	lan.Name = strings.TrimSpace(lan.Name)
	lan.Extension = strings.TrimSpace(lan.Extension)
}

func fetchAllLanguages() ([]Language, error) {
	rows, err := config.DB.Query("SELECT * FROM languages")
	if err != nil {
		return nil, err
	}

	languageList := make([]Language, 0)
	for rows.Next() {
		var language Language
		if err := rows.Scan(&language.Id, &language.Name, &language.Extension, &language.Command); err == nil {
			language.normaliseLanguage()
			languageList = append(languageList, language)
		}
	}

	return languageList, nil
}

func FetchLanguageById(languageId int) (Language, error) {
	var language Language
	if err := config.DB.
		QueryRow("SELECT * FROM languages WHERE id = $1", languageId).
		Scan(&language.Id, &language.Name, &language.Extension, &language.Command); err != nil {
		return Language{}, err
	}
	language.normaliseLanguage()
	return language, nil
}
