package rod_test

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-rod/rod"
	"github.com/stretchr/testify/suite"
	"github.com/ysmood/kit"
	"go.uber.org/goleak"
)

var slash = filepath.FromSlash

// S test suite
type S struct {
	suite.Suite
	browser *rod.Browser
	page    *rod.Page
}

func TestMain(m *testing.M) {
	// to prevent false positive of goleak
	http.DefaultClient = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	goleak.VerifyTestMain(m)
}

func Test(t *testing.T) {
	s := new(S)
	s.browser = rod.New().Client(nil).Connect()

	defer s.browser.Close()

	s.page = s.browser.Page("")
	s.page.Viewport(800, 600, 1, false)

	suite.Run(t, s)
}

// get abs file path from fixtures folder, return sample "file:///a/b/click.html"
func srcFile(path string) string {
	return "file://" + file(path)
}

// get abs file path from fixtures folder, return sample "/a/b/click.html"
func file(path string) string {
	f, err := filepath.Abs(slash(path))
	kit.E(err)
	return f
}

func ginHTML(body string) gin.HandlerFunc {
	return func(ctx kit.GinContext) {
		ctx.Header("Content-Type", "text/html;")
		kit.E(ctx.Writer.WriteString(body))
	}
}

func ginString(body string) gin.HandlerFunc {
	return func(ctx kit.GinContext) {
		kit.E(ctx.Writer.WriteString(body))
	}
}

func ginHTMLFile(path string) gin.HandlerFunc {
	body, err := kit.ReadString(path)
	kit.E(err)
	return ginHTML(body)
}

// returns url prefix, engin, close
func serve() (string, *gin.Engine, func()) {
	srv := kit.MustServer("127.0.0.1:0")
	opt := &http.Server{}
	opt.SetKeepAlivesEnabled(false)
	srv.Set(opt)
	go func() { kit.Noop(srv.Do()) }()

	url := "http://" + srv.Listener.Addr().String()

	return url, srv.Engine, func() { kit.E(srv.Listener.Close()) }
}
