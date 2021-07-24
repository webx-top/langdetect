package langdetect

import "net/http"

func NewContextStandard(r *http.Request, w http.ResponseWriter) Context {
	return &contextStandard{r: r, w: w}
}

type contextStandard struct {
	r *http.Request
	w http.ResponseWriter
}

func (c *contextStandard) Path() string {
	return c.r.URL.Path
}

func (c *contextStandard) SetPath(path string) {
	c.r.URL.Path = path
}

func (c *contextStandard) Header(name string) string {
	return c.r.Header.Get(name)
}

func (c *contextStandard) Query(name string) string {
	return c.r.URL.Query().Get(name)
}

func (c *contextStandard) Cookie(name string) string {
	cookie, err := c.r.Cookie(name)
	if err != nil {
		return ``
	}
	return cookie.Value
}

func (c *contextStandard) SetCookie(name string, value string) {
	cookie := &http.Cookie{Name: name, Value: value, Path: `/`}
	c.w.Header().Add(`Set-Cookie`, cookie.String())
}
