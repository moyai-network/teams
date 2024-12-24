package lang

import (
	"fmt"
	"github.com/moyai-network/teams/internal/ports/model"

	"github.com/restartfu/gophig"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"golang.org/x/text/language"
)

// translationData contains the properties and translations of a language.
type translationData struct {
	Properties struct {
		Name  string `json:"name"`
		Image string `json:"image"`
	} `json:"properties"`
	Translations map[string]string `json:"translations"`
}

// translations stores a mapping between the language and the translation data.
var translations = make(map[language.Tag]translationData)

// Register registers a translation file and adds the decoded data to the translations list.
func Register(lang language.Tag) {
	var dat translationData
	if err := gophig.GetConfComplex(fmt.Sprintf("assets/translations/%v.json", lang.String()), gophig.JSONMarshaler{}, &dat); err != nil {
		panic(err)
	}
	translations[lang] = dat
}

// Properties returns the name and image of a language.
func Properties(lang model.Language) (string, string, bool) {
	dat, ok := translations[lang.Tag]
	return dat.Properties.Name, dat.Properties.Image, ok
}

// Translatef returns the translated version of a string.
func Translatef(lang model.Language, key string, a ...interface{}) string {
	return text.Colourf(Translate(lang, key), a...)
}

// Translate returns the translated version of a string.
func Translate(lang model.Language, key string) string {
	t, ok := translations[lang.Tag]
	if !ok {
		t = translations[language.English]
	}
	if r, ok := t.Translations[key]; ok {
		return r
	} else {
		return translations[language.English].Translations[key]
	}
}
