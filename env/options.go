package env

type (
	options struct {
		TagName string
	}
	OptionsFunc func(o *options)
)

func TagName(name string) OptionsFunc {
	return func(o *options) { o.TagName = name }
}
