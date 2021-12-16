package server

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"time"
)

type Response struct {
	Status string      `json:"status"`
	Code   int         `json:"code"`
	Body   interface{} `json:"body"`
	Msg    string      `json:"msg"`
	Time   time.Time   `json:"time"`
}

func Result(code int, data interface{}, msg string, c *fiber.Ctx) error {
	return c.Status(code).JSON(Response{
		Status: http.StatusText(code),
		Code:   code,
		Body:   data,
		Msg:    msg,
		Time:   time.Now(),
	})
}

func Ok(c *fiber.Ctx) error {
	return Result(http.StatusOK, nil, "Success", c)
}

func OkWithMessage(message string, c *fiber.Ctx) error {
	return Result(http.StatusOK, nil, message, c)
}

func OkWithData(data interface{}, c *fiber.Ctx) error {
	return Result(http.StatusOK, data, "Success", c)
}

func OkWithDetailed(data interface{}, message string, c *fiber.Ctx) error {
	return Result(http.StatusOK, data, message, c)
}

func OkWithRaw(contentType string, data []byte, c *fiber.Ctx) error {
	return c.Status(200).Type(contentType, "utf-8").Send(data)
}

func OkWithEmptyRaw(contentType string, c *fiber.Ctx) error {
	return c.Status(204).Type(contentType, "utf-8").Send(nil)
}

func OkWithNotFound(contentType string, c *fiber.Ctx) error {
	return c.Status(404).Type(contentType, "utf-8").Send(nil)
}

func Fail(code int, c *fiber.Ctx) error {
	return Result(code, nil, "操作失败", c)
}

func FailWithMessage(code int, message string, c *fiber.Ctx) error {
	return Result(code, nil, message, c)
}

func FailWithDetailed(code int, data interface{}, message string, c *fiber.Ctx) error {
	return Result(code, data, message, c)
}
