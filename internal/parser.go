package internal

import (
	"credit-risk-mvp/internal/domain"
)

func ParseApplication(payload map[string]any) (domain.Application, error) {
	rawName, err := getString(payload, "name")
	if err != nil {
		return domain.Application{}, err
	}
	name, err := domain.NewFullName(rawName)
	if err != nil {
		return domain.Application{}, err
	}
	rawBirthdate, err := getString(payload, "birthdate")
	if err != nil {
		return domain.Application{}, err
	}
	birthdate, err := domain.NewBirthdate(rawBirthdate)
	if err != nil {
		return domain.Application{}, err
	}
	rawPhone, err := getString(payload, "phone")
	if err != nil {
		return domain.Application{}, err
	}
	phone, err := domain.NewPhone(rawPhone)
	if err != nil {
		return domain.Application{}, err
	}
	rawPassport, err := getString(payload, "passport")
	if err != nil {
		return domain.Application{}, err
	}
	passport, err := domain.NewPassport(rawPassport)
	if err != nil {
		return domain.Application{}, err
	}
	residency, err := getString(payload, "residency")
	if err != nil {
		return domain.Application{}, err
	}
	firstTime, err := getBool(payload, "first_time")
	if err != nil {
		return domain.Application{}, err
	}
	rawAmount, err := getInt(payload, "requested_amount")
	if err != nil {
		return domain.Application{}, err
	}
	amount, err := domain.NewAmount(rawAmount)
	if err != nil {
		return domain.Application{}, err
	}
	return domain.Application{
		Payload:         payload,
		Name:            name,
		Birthdate:       birthdate,
		Phone:           phone,
		Passport:        passport,
		RequestedAmount: amount,
		Residency:       residency,
		FirstTime:       firstTime,
	}, nil
}
