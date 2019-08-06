package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func routes(ctx context.Context, e *echo.Echo) {
	// Swagger OpenAPI spec
	e.File("/swagger.json", "/swagger.json")
	// GET endpoint
	e.GET("/", func(c echo.Context) error {
		// cache this one
		helloworld := "hello world"
		c.Response().Header().Set("Cache-Control", "max-age=3600") // 1 hour
		c.Response().Header().Set("Content-Type", "application/json; charset=UTF-8")

		return c.JSON(http.StatusOK, helloworld)
	})

	e.GET("/still", func(c echo.Context) error {

		source := c.QueryParam("s")

		task := NewExtractImageTask(source)

		// if err := task.Run(); err != nil {
		// 	return c.JSON(http.StatusInternalServerError, fmt.Sprintf("task failed: %s", err.Error()))
		// }
		stream, err := task.StderrPipe()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("task output failed, error pipe failed: %s", err.Error()))
		}

		results, err := task.Output()
		if err != nil {
			return c.Stream(http.StatusInternalServerError, "text/plain", stream)
		}

		return c.Blob(http.StatusOK, "image/png", results)
	})

	e.GET("/scale", func(c echo.Context) error {

		source := c.QueryParam("s")
		height := c.QueryParam("h")
		width := c.QueryParam("w")

		h, err := strconv.ParseUint(height, 10, 16)
		if err != nil {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("error parsing height param: %s", err.Error()))
		}
		w, err := strconv.ParseUint(width, 10, 16)
		if err != nil {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("error parsing width param: %s", err.Error()))
		}

		task := NewScaleImageTask(source, uint(w), uint(h))

		// if err := task.Run(); err != nil {
		// 	return c.JSON(http.StatusInternalServerError, fmt.Sprintf("task failed: %s", err.Error()))
		// }

		results, err := task.Output()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("task output failed: %s", err.Error()))
		}

		return c.Blob(http.StatusOK, "image/png", results)
	})
}
