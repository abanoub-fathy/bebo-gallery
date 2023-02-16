package utils

// Must is used to panic an error if exist
func Must(err error) {
	if err != nil {
		panic(err)
	}
}
