package core

import "github.com/gosimple/slug"

func ValidateSlug(input string) bool {
	return input == slug.Make(input)
}
