package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var validate *validator.Validate
var trans ut.Translator

type FieldError struct {
	Key    string
	Errors []string
}

func (fE FieldError) PrefixAll(prefix string) FieldError {
	prefixed := make([]string, len(fE.Errors))

	for i, msg := range fE.Errors {
		prefixed[i] = fmt.Sprintf("%s%s", prefix, msg)
	}

	return FieldError{
		Key: fE.Key,
		Errors: prefixed,
	}
}

type ValidationError struct {
	Errs []FieldError
}

func (vE ValidationError) PrefixAll(prefix string) ValidationError {
	prefixed := make([]FieldError, len(vE.Errs))

	for i, err := range vE.Errs {
		prefixed[i] = err.PrefixAll(prefix)
	}

	return ValidationError{
		Errs: prefixed,
	}
}

func (err ValidationError) Error() string {
	return "invalid data"
}

type Validator struct {
	validate *validator.Validate
	trans    ut.Translator
}

type Translation struct {
	Rule          string
	RegisterFunc  func(ut ut.Translator) error
	TranslateFunc func(ut ut.Translator, fe validator.FieldError) string
}

type CustomRule struct {
	Rule    string
	Handler func(fl validator.FieldLevel) bool
}

type StructValidator struct {
	Struct    any
	Validator validator.StructLevelFunc
}

type Extension interface {
	Translations() []Translation
	Rules() []CustomRule
	StructValidators() []StructValidator
}

func (v *Validator) Validate(in any) error {

	err := validate.Struct(in)

	if err == nil {
		return nil
	}

	errs, ok := err.(validator.ValidationErrors)

	if !ok {
		return err
	}

	// return err
	fieldErrors := []FieldError{}

	for _, validationErr := range errs {
		// The StructNamespace includes the type of the struct we're validating
		// we don't actually care about the top-level so we remove it
		// e.g. MyStructType.MyField.SubField becomes MyField.SubField
		field := strings.Join(
			strings.Split(validationErr.StructNamespace(), ".")[1:],
			".",
		)

		fieldErrors = append(fieldErrors, FieldError{
			Key:    field,
			Errors: []string{validationErr.Translate(trans)},
		})
	}

	return ValidationError{
		Errs: fieldErrors,
	}
}

func NewValidator(extensions... Extension) *Validator {
	en := en.New()
	uni := ut.New(en, en)

	trans, _ = uni.GetTranslator("en")

	validate = validator.New()
	en_translations.RegisterDefaultTranslations(validate, trans)

	for _, ext := range extensions {
		for _, rule := range ext.Rules() {
			validate.RegisterValidation(rule.Rule, rule.Handler)
		}

		for _, translation := range ext.Translations() {
			validate.RegisterTranslation(
				translation.Rule,
				trans,
				translation.RegisterFunc,
				translation.TranslateFunc,
			)
		}

		for _, sv := range ext.StructValidators() {
			validate.RegisterStructValidation(sv.Validator, sv.Struct)
		}
	}

	return &Validator{
		validate: validate,
		trans:    trans,
	}
}
