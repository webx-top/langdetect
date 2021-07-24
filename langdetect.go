package langdetect

import (
	"net/http"
	"regexp"
	"strings"
)

var LangVarName = `lang`

func New(fallback string, supported map[string]bool) *Languages {
	if supported == nil {
		supported = make(map[string]bool)
	}
	return &Languages{
		Fallback:  fallback,
		Supported: supported,
		VarName:   LangVarName,
	}
}

type Languages struct {
	Supported map[string]bool
	Fallback  string
	VarName   string
}

func (a *Languages) DetectFromURIPath(ctx Context) string {
	p := strings.TrimPrefix(ctx.Path(), `/`)
	s := strings.Index(p, `/`)
	var lang string
	if s != -1 {
		lang = p[0:s]
	} else {
		lang = p
	}
	if len(lang) > 0 {
		on, ok := a.Supported[lang]
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

func (a *Languages) IsSupported(lang string) bool {
	if len(lang) > 0 {
		if on, ok := a.Supported[lang]; ok {
			return on
		}
	}
	return false
}

var headerAcceptRemove = regexp.MustCompile(`;[\s]*q=[0-9.]+`)

func (a *Languages) DetectFromHeader(ctx Context) string {
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

func (a *Languages) Detect(ctx Context) string {
	lang := ctx.Query(a.VarName)
	var hasCookie bool
	defer func() {
		if !hasCookie {
			ctx.SetCookie(a.VarName, lang)
		}
	}()
	if a.IsSupported(lang) {
		return lang
	}
	lang = a.DetectFromURIPath(ctx)
	if a.IsSupported(lang) {
		return lang
	}
	lang = ctx.Cookie(a.VarName)
	if a.IsSupported(lang) {
		hasCookie = true
		return lang
	}
	lang = a.DetectFromHeader(ctx)
	if len(lang) == 0 {
		lang = a.Fallback
	}
	return lang
}

func (a *Languages) DetectFromRequest(r *http.Request, w http.ResponseWriter) string {
	return a.Detect(NewContextStandard(r, w))
}
