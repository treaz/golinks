package models

import (
	"database/sql"
	"errors"
	"regexp"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type LinkInput struct {
	Keyword         string `json:"keyword"`
	Description     string `json:"description"`
	Destination     string `json:"destination"`
	Views           int    `json:"views"`
	Tags            []Tag  `json:"tags"`
	IsParameterized bool   `json:"is_parameterized"`
}

type Link struct {
	ID              uint         `json:"id" gorm:"primarykey"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
	DeletedAt       sql.NullTime `json:"-" gorm:"index"`
	Keyword         string       `json:"keyword" gorm:"unique;index"`
	Description     string       `json:"description"`
	Destination     string       `json:"destination"`
	Views           int          `json:"views"`
	Tags            []Tag        `json:"tags" gorm:"many2many:link_tags;"`
	IsParameterized bool         `json:"is_parameterized"`
}

func hasScheme(value interface{}) error {
	s, _ := value.(string)
	if !(strings.HasPrefix(s, "http:") || strings.HasPrefix(s, "https:")) {
		return errors.New("requires scheme http or https")
	}
	return nil
}

func validDestination(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return errors.New("must be a string")
	}
	// Replace {*} with a safe placeholder for validation
	s = strings.ReplaceAll(s, "{*}", "placeholder")

	// Use ozzo-validation's URL validator on the modified string
	return is.URL.Validate(s)
}

// TODO Put in validation for programtic links gh/%s
func (l LinkInput) Validate() error {
	return validation.ValidateStruct(&l,
		// Name cannot be empty, and the length must be between 3 and 20.
		// Regex for keyword https://regex101.com/r/HetwqX/1
		validation.Field(&l.Keyword,
			validation.Required,
			validation.Length(1, 100),
			validation.NotIn("api", "static", "directory", "healthz", "favicon.ico"),
			validation.Match(regexp.MustCompile(`^([a-zA-Z0-9\-\/]+)`))),
		validation.Field(&l.Destination,
			validation.Required,
			validation.By(validDestination),
			validation.By(hasScheme),
		),
	)
}

func (li *LinkInput) ToNative() Link {

	isParameterized := li.IsParameterized
	if strings.Contains(li.Destination, "{*}") {
		isParameterized = true
	}

	return Link{
		Keyword:         strings.ToLower(li.Keyword),
		Destination:     li.Destination,
		Description:     li.Description,
		IsParameterized: isParameterized,
	}
}

type QueryString struct {
	Sort  string
	Order string
}

func (qs QueryString) Validate() error {
	return validation.ValidateStruct(&qs,
		validation.Field(&qs.Sort, validation.Match(regexp.MustCompile(`^([a-zA-Z]+)`))),
		validation.Field(&qs.Order, validation.Match(regexp.MustCompile(`^([a-zA-Z]+)`))),
	)
}

type Links []Link
