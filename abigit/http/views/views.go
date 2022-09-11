package views

import (
	"github.com/codemicro/abigit/abigit/config"
	"github.com/codemicro/abigit/abigit/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

//go:generate neontc --extension ntc.html

type RenderContext struct {
	// set by the page in question:
	pageTitle      string
	reducedContent bool // for transitionary pages only

	// set from environment when instantiating
	externalURL  string
	platformName string

	// set from arguments when instantiating
	user *models.User
}

// NewRenderContext creates a new context for rendering pages.
//
// User is the currently authenticated user. May be `nil` if no use is
// currently authenticated.
func NewRenderContext(user *models.User) *RenderContext {
	rctx := &RenderContext{
		externalURL:  config.HTTP.ExternalURL,
		platformName: config.Platform.Name,
		user:         user,
	}

	return rctx
}

func (rctx *RenderContext) isAuthed() bool {
	return rctx.user != nil
}

func SendPage(ctx *fiber.Ctx, content string) error {
	ctx.Type("html")
	return ctx.SendString(content)
}

func formatFileSize(s int64) string {
	unit := "B"

	c := float64(s)
	for _, item := range []string{"KB", "MB", "GB"} {
		if c >= 1000 {
			c /= 1000
			unit = item
		} else {
			break
		}
	}

	precision := 2
	if unit == "B" {
		precision = 0
	}

	return strconv.FormatFloat(c, 'f', precision, 64) + unit
}

var statusMessages = map[int]string{
	100: "Continue",
	101: "Switching Protocols",
	102: "Processing",
	103: "Early Hints",
	200: "OK",
	201: "Created",
	202: "Accepted",
	203: "Non-Authoritative Information",
	204: "No Content",
	205: "Reset Content",
	206: "Partial Content",
	207: "Multi-Status",
	208: "Already Reported",
	226: "IM Used",
	300: "Multiple Choices",
	301: "Moved Permanently",
	302: "Found",
	303: "See Other",
	304: "Not Modified",
	305: "Use Proxy",
	306: "Switch Proxy",
	307: "Temporary Redirect",
	308: "Permanent Redirect",
	400: "Bad Request",
	401: "Unauthorized",
	402: "Payment Required",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	406: "Not Acceptable",
	407: "Proxy Authentication Required",
	408: "Request Timeout",
	409: "Conflict",
	410: "Gone",
	411: "Length Required",
	412: "Precondition Failed",
	413: "Request Entity Too Large",
	414: "Request URI Too Long",
	415: "Unsupported Media Type",
	416: "Requested Range Not Satisfiable",
	417: "Expectation Failed",
	418: "I'm a teapot",
	421: "Misdirected Request",
	422: "Unprocessable Entity",
	423: "Locked",
	424: "Failed Dependency",
	426: "Upgrade Required",
	428: "Precondition Required",
	429: "Too Many Requests",
	431: "Request Header Fields Too Large",
	451: "Unavailable For Legal Reasons",
	500: "Internal Server Error",
	501: "Not Implemented",
	502: "Bad Gateway",
	503: "Service Unavailable",
	504: "Gateway Timeout",
	505: "HTTP Version Not Supported",
	506: "Variant Also Negotiates",
	507: "Insufficient Storage",
	508: "Loop Detected",
	510: "Not Extended",
	511: "Network Authentication Required",
}
