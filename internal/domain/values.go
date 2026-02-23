package domain

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

var passportRegex = regexp.MustCompile(`^\d{4}\s?\d{6}$`)
var phoneRegex = regexp.MustCompile(`^\+?\d{10,15}$`)

type Passport string

func NewPassport(s string) (Passport, error) {
	//1234 567890 or 1234567890
	if !passportRegex.MatchString(s) {
		return "", errors.New("invalid passport format, expected '1234 567890'")
	}
	return Passport(s), nil
}

type Phone string

func NewPhone(phone string) (Phone, error) {
	if !phoneRegex.MatchString(phone) {
		return "", errors.New("invalid phone format")
	}
	return Phone(phone), nil
}

type Amount int

func NewAmount(val int) (Amount, error) {
	if val <= 0 {
		return 0, errors.New("amount must be greater than zero")
	}
	return Amount(val), nil
}

type Birthdate time.Time

func NewBirthdate(s string) (Birthdate, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return Birthdate{}, errors.New("invalid date format")
	}
	now := time.Now()
	age := now.Year() - t.Year()
	//if there's no birthdate in this year then decrease age by 1
	if now.Month() < t.Month() || (now.Month() == t.Month() && now.Day() < t.Day()) {
		age--
	}
	if age < 18 {
		return Birthdate{}, errors.New("client is under 18")
	}
	return Birthdate(t), nil
}
func (b Birthdate) Age() int {
	t := time.Time(b)
	now := time.Now()
	age := now.Year() - t.Year()
	if now.Month() < t.Month() || (now.Month() == t.Month() && now.Day() < t.Day()) {
		age--
	}
	return age
}

type FullName string

func NewFullName(fullName string) (FullName, error) {
	if len(strings.TrimSpace(fullName)) == 0 {
		return "", errors.New("name cannot be empty")
	}
	return FullName(fullName), nil
}
func (n FullName) HasPatronymic() bool {
	return len(strings.Fields(string(n))) >= 3
}
