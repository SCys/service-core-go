package core

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/xid"
)

// FiberJSONError 输出JSON异常
func FiberJSONError(ctx *fiber.Ctx, code int, err error) error {
	ctx.Set("content_type", "application/json; charset=utf-8")

	content, err := JSONError(code, err)
	if err != nil {
		return err
	}

	ctx.Set("Content-Type", "application/json")
	_, err = ctx.Write(content)
	return err
}

// FiberJSON 输出JSON
func FiberJSON(ctx *fiber.Ctx, data H) error {
	content, err := jsoniter.Marshal(H{"data": data})

	if err != nil {
		return err
	}

	ctx.Set("Content-Type", "application/json")
	_, err = ctx.Write(content)
	return err
}

// FiberIP 获取请求IP， X-Real-IP > X-Forward-Forr > fiber.Ctx.IP
func FiberIP(ctx *fiber.Ctx) string {
	addr := ctx.Get("X-Real-IP")
	if addr != "" {
		return addr
	}

	addr = ctx.Get("X-Forward-For")
	if addr != "" {
		items := strings.Split(addr, ",")
		if len(items) > 1 {
			return items[0]
		}

		return addr
	}

	return ctx.IP()
}

// FiberBasicInfo 增加额外的基础信息到 Context
func FiberBasicInfo() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		ctx := c.Context()

		ctx.SetUserValue("_id", xid.New().String())
		ctx.SetUserValue("client_ip", FiberIP(c))

		return c.Next()
	}
}
