package forms

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// Form creates a custom form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// Valid returns true if there are no errors, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required checks for required field
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		val := f.Get(field)
		if strings.TrimSpace(val) == "" {
			log.Println("form error 1")
			f.Errors.Add(field, "The field cannot be blank")
		}
	}
}

// Has checks if form field is not empty
func (f *Form) Has(field string, r *http.Request) bool {
	// deprecated since this will fail the test
	// x := r.Form.Get(field)
	x := f.Get(field)
	if x == "" {
		// f.Errors.Add(field, "This field cannot be blank")
		return false
	}
	return true
}

// MinLength checks for string minimum length
func (f *Form) MinLength(field string, length int, r *http.Request) bool {
	// deprecated since this will fail the test
	// x := r.Form.Get(field)
	x := f.Get(field)

	if len(x) < length {
		log.Println("form error 2")

		f.Errors.Add(field, fmt.Sprintf("this field must be at least %d characters long", length))
		return false
	}
	return true
}

// IsEmail checks for valid email address
func (f *Form) IsEmail(field string) bool {

	email := f.Get(field)

	emailRegex := regexp.MustCompile(`^[A-Za-z0-9._+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

	if !emailRegex.MatchString(email) {
		f.Errors.Add(field, "invalid email address")
		return false
	}

	domain := strings.Split(email, "@")
	nonBussinessEmail := []string{"gmail", "outlook", "live", "yahoo"}

	for _, e := range nonBussinessEmail {
		if strings.Contains(domain[1], e) {
			f.Errors.Add(field, "invalid bussiness email address")
			return false
		}
	}

	return true
}
