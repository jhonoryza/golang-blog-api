package helper

func CastNilString(value *string) string {
	var summary string
	if value != nil {
		summary = *value
	}
	return summary
}
