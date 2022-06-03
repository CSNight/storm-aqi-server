package server

import (
	"github.com/csnight/storm-aqi-server/tools"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"time"
)

type Response struct {
	Status string      `json:"status"`
	Code   int         `json:"code"`
	Body   interface{} `json:"body"`
	Msg    string      `json:"msg"`
	Time   int64       `json:"time"`
}

func Result(code int, data interface{}, msg string, c *fiber.Ctx) error {
	return c.Status(code).JSON(Response{
		Status: http.StatusText(code),
		Code:   code,
		Body:   data,
		Msg:    msg,
		Time:   time.Now().UnixMilli(),
	})
}

func Ok(c *fiber.Ctx) error {
	return Result(http.StatusOK, nil, "Success", c)
}

func OkWithMessage(message string, c *fiber.Ctx) error {
	return Result(http.StatusOK, nil, message, c)
}

func OkWithData(data interface{}, c *fiber.Ctx) error {
	etag, err := tools.NewNanoId()
	if err != nil {
		return err
	}
	c.Set("Cache-Control", "max-age=600")
	c.Set("ETag", etag)
	c.Set("Last-Modified", time.Now().Format(time.RFC1123))
	return Result(http.StatusOK, data, "Success", c)
}

func OkWithDetailed(data interface{}, message string, c *fiber.Ctx) error {
	return Result(http.StatusOK, data, message, c)
}

func OkWithRaw(contentType string, data []byte, c *fiber.Ctx) error {
	return c.Status(200).Type(contentType).Send(data)
}

func OkWithEmptyRaw(contentType string, c *fiber.Ctx) error {
	return c.Status(204).Type(contentType).Send(nil)
}

func OkWithNotFound(contentType string, c *fiber.Ctx) error {
	return c.Status(404).Type(contentType).Send(nil)
}

func Fail(code int, c *fiber.Ctx) error {
	return Result(code, nil, "Unknown error", c)
}

func FailWithMessage(code int, message string, c *fiber.Ctx) error {
	return Result(code, nil, message, c)
}

func FailWithDetailed(code int, data interface{}, message string, c *fiber.Ctx) error {
	return Result(code, data, message, c)
}
