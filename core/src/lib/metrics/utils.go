package myLogger

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetReqMessage(c *fiber.Ctx) string {
	reqBody := MaskSensibleInformation(string(c.Request().Body()))
	reqHeaders := MaskSensibleInformation(string(c.Request().Header.Header()))
	LabelReqStr := fmt.Sprintf(
		"\n\nRequest Headers:\n%s\nRequest Body:\n%s",
		TruncateUTF8(reqHeaders, 3000),
		TruncateUTF8(reqBody, 3000),
	)
	return LabelReqStr
}

func GetResMessage(c *fiber.Ctx) string {
	resHeaders := MaskSensibleInformation(string(c.Response().Header.Header()))
	resBody := MaskSensibleInformation(string(c.Response().Body()))
	LabelResStr := fmt.Sprintf(
		"\n\nResponse Headers:\n%s\nResponse Body:\n%s",
		TruncateUTF8(resHeaders, 3000),
		TruncateUTF8(resBody, 3000),
	)
	return LabelResStr
}

type Labels struct {
	App    string `json:"app"`
	Method string `json:"method"`
	Path   string `json:"path"`
	IP     string `json:"ip"`
	Host   string `json:"host"`
}

func (l *Labels) GetReqLabels() map[string]string {
	lokiLabels := map[string]string{
		"app":    l.App,
		"level":  "info",
		"title":  "Incoming Request",
		"type":   "request",
		"method": l.Method,
		"path":   l.Path,
		"ip":     l.IP,
		"host":   l.Host,
	}
	return lokiLabels
}

func (l *Labels) GetResLabels(startTime time.Time, status int) map[string]string {
	lokiLabels := l.GetReqLabels()
	duration := time.Since(startTime)

	title := fmt.Sprintf("%d - %s", status, http.StatusText(status))

	lokiLabels["type"] = "response"
	lokiLabels["status_code"] = fmt.Sprintf("%d", status)
	lokiLabels["duration"] = duration.String()
	lokiLabels["title"] = title

	if status >= 100 && status < 200 {
		lokiLabels["level"] = "info"
	} else if status >= 200 && status < 300 {
		lokiLabels["level"] = "success"
	} else if status >= 300 && status < 400 {
		lokiLabels["level"] = "info"
	} else if status >= 400 && status < 500 {
		lokiLabels["level"] = "warning"
	} else if status >= 500 && status < 600 {
		lokiLabels["level"] = "error"
	} else {
		lokiLabels["level"] = "strange"
	}
	return lokiLabels
}

func MaskSensibleInformation(body string) string {
	patterns := []string{
		`("password"\s*:\s*")([^"]*)(")`,
		`("token"\s*:\s*")([^"]*)(")`,
		`("secret"\s*:\s*")([^"]*)(")`,
		`("access_token"\s*:\s*")([^"]*)(")`,
		`("refresh_token"\s*:\s*")([^"]*)(")`,
		`(namespace.HeadersKey.Auth\s*:\s*")([^"]*)(")`,
	}

	masked := body
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		masked = re.ReplaceAllString(masked, `$1***masked***$3`)
	}

	return masked
}

func TruncateUTF8(s string, limit int) string {
	runes := []rune(s)
	if len(runes) > limit {
		return string(runes[:limit]) + "...(truncated)"
	}
	return s
}

