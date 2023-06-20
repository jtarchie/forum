package templates

import (
	"github.com/gosimple/slug"
)

func Slug(name string) string {
	return slug.Make(name)
}
