package valid

import "regexp"

var (
	RgxOneLower        = regexp.MustCompile(`[a-z]+`)
	RgxOneUpper        = regexp.MustCompile(`[A-Z]+`)
	RgxOneNumeric      = regexp.MustCompile(`\d+`)
	RgxPasswordLength  = regexp.MustCompile(`^.{7,15}$`)
	RgxSnipName        = regexp.MustCompile(`^[a-zA-Z0-9\-\._]{1,50}$`)
	RgxExtension       = regexp.MustCompile(`^[a-zA-Z0-9]{1,12}$`)
	RgxPouchName       = regexp.MustCompile(`^[a-zA-Z0-9\-_]{1,30}$`)
	RgxSnipDescription = regexp.MustCompile(`.*`)
	RgxUsername        = regexp.MustCompile(`^[a-zA-Z0-9\-]{3,15}$`)
	RgxEmail           = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
)

func Test(input string, rgx ...*regexp.Regexp) bool {
	for _, v := range rgx {
		if match := v.Find([]byte(input)); len(match) == 0 {
			return false
		}
	}
	return true
}
