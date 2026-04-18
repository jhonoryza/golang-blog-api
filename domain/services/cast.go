package services

func CastNilString(value *string) string {
	if value != nil {
		return *value
	}
	return ""
}