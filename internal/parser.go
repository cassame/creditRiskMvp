package internal

func ParseApplication(payload map[string]any) (Application, error) {
	name, err := getString(payload, "name")
	if err != nil {
		return Application{}, err
	}
	birthdate, err := getString(payload, "birthdate")
	if err != nil {
		return Application{}, err
	}
	phone, err := getString(payload, "phone")
	if err != nil {
		return Application{}, err
	}
	passport, err := getString(payload, "passport")
	if err != nil {
		return Application{}, err
	}
	residency, err := getString(payload, "residency")
	if err != nil {
		return Application{}, err
	}
	firstTime, err := getBool(payload, "first_time")
	if err != nil {
		return Application{}, err
	}
	requestedAmount, err := getInt(payload, "requested_amount")
	if err != nil {
		return Application{}, err
	}
	return Application{
		Payload:         payload,
		Name:            name,
		Birthdate:       birthdate,
		Phone:           phone,
		Passport:        passport,
		RequestedAmount: requestedAmount,
		Residency:       residency,
		FirstTime:       firstTime,
	}, nil
}
