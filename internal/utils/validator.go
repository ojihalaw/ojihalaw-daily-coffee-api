package utils

import (
	"github.com/go-playground/locales/en"              // locale bahasa
	ut "github.com/go-playground/universal-translator" // universal translator
	"github.com/go-playground/validator/v10"           // validator
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/spf13/viper"
)

type Validator struct {
	Validate   *validator.Validate
	Translator ut.Translator
}

func NewValidator(viper *viper.Viper) *Validator {
	validate := validator.New()

	// setup locale English
	locale := en.New()
	uni := ut.New(locale, locale)
	trans, _ := uni.GetTranslator("en")

	// register default translation
	_ = en_translations.RegisterDefaultTranslations(validate, trans)

	return &Validator{
		Validate:   validate,
		Translator: trans,
	}
}

// Helper untuk translate error ke slice string
func (v *Validator) TranslateError(err error) []string {
	var errs []string
	if err == nil {
		return errs
	}

	for _, e := range err.(validator.ValidationErrors) {
		errs = append(errs, e.Translate(v.Translator))
	}
	return errs
}
