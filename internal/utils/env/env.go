package env

import "os"

func Get(k string, dfault ...string) string {
	if os.Getenv(k) == "" {
		if len(dfault) > 0 {
			return dfault[0]
		}
		return ""
	}
	return os.Getenv(k)
}
