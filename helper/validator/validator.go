package validator

import (
	"regexp"
	"time"
)

func StringCheck(input string, minLength, maxLength int) bool {
	// minLength += 1
	length := len(input)
	// fmt.Println(length >= minLength && length <= maxLength)
	return length >= minLength && length <= maxLength
}

func IsURL(s string) bool {

	regex := `https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,4}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)$`
	pattern := regexp.MustCompile(regex)
	return pattern.MatchString(s)
}
func IsDateISO860(s string)bool {
	_, err := time.Parse("2006-01-02T15:04:05Z07:00", s)
	return err == nil
}

func IsEmail(email string) bool {
    // Define a regular expression for validating an email
    var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    
    // Return true if the email matches the regex
    return emailRegex.MatchString(email)
}