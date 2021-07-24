package langdetect

import (
	"net/http"
	"regexp"
	"strings"
)

var LangVarName = `lang`

type Languages map[string]bool

func (a Languages) DetectFromURIPath(ctx Context) string {
	p := strings.TrimPrefix(ctx.Path(), `/`)
	s := strings.Index(p, `/`)
	var lang string
	if s != -1 {
		lang = p[0:s]
	} else {
		lang = p
	}
	if len(lang) > 0 {
		on, ok := a[lang]
		if !ok {
			return ``
		}
		ctx.SetPath(strings.TrimPrefix(p, lang))
		if !on {
			return ``
		}
	}
	return lang
}

func (a Languages) IsSupported(lang string) bool {
	if len(lang) > 0 {
		if on, ok := a[lang]; ok {
			return on
		}
	}
	return false
}

var headerAcceptRemove = regexp.MustCompile(`;[\s]*q=[0-9.]+`)

func (a Languages) DetectFromHeader(ctx Context) string {
	al := ctx.Header(`Accept-Language`)
	al = headerAcceptRemove.ReplaceAllString(al, ``)
	lg := strings.SplitN(al, `,`, 5)
	for _, lang := range lg {
		lang = strings.ToLower(lang)
		if a.IsSupported(lang) {
			return lang
		}
	}
	return ``
}

func (a Languages) Detect(ctx Context) string {
	lang := ctx.Query(LangVarName)
	var hasCookie bool
	defer func() {
		if !hasCookie {
			ctx.SetCookie(LangVarName, lang)
		}
	}()
	if a.IsSupported(lang) {
		return lang
	}
	lang = a.DetectFromURIPath(ctx)
	if a.IsSupported(lang) {
		return lang
	}
	lang = ctx.Cookie(LangVarName)
	if a.IsSupported(lang) {
		hasCookie = true
	}
	lang = a.DetectFromHeader(ctx)
	return lang
}

func (a Languages) DetectFromRequest(r *http.Request, w http.ResponseWriter) string {
	return a.Detect(NewContextStandard(r, w))
}
