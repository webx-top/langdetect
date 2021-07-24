package langdetect

type Context interface {
	Path() string
	SetPath(path string)
	Header(name string) string
	Query(name string) string
	Cookie(name string) string
	SetCookie(name string, value string)
}
