package utils

const na = "N/A"

func SetDefaultBuildInfo(v *string) {
	if *v == "" {
		*v = na
	}
}
