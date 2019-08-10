package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

const swaggerJSON = "swagger.json"

func routes(ctx context.Context, e *echo.Echo) {
	// Swagger OpenAPI spec
	if fileExists(swaggerJSON) {
		e.File("/swagger.json", swaggerJSON)
	} else if fileExists("/" + swaggerJSON) {
		e.File("/swagger.json", "/"+swaggerJSON)
	}
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
		offset := c.QueryParam("o")

		h, err := strconv.ParseUint(height, 10, 16)
		if err != nil {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("error parsing height param: %s", err.Error()))
		}
		w, err := strconv.ParseUint(width, 10, 16)
		if err != nil {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("error parsing width param: %s", err.Error()))
		}
		o, err := strconv.ParseUint(offset, 10, 16)
		if err != nil {
			// dont error - this arg is optional
			// return c.JSON(http.StatusBadRequest, fmt.Sprintf("error parsing offset param: %s", err.Error()))
		}

		task := NewScaleImageTask(source, uint(w), uint(h), time.Duration(o)*time.Second)

		// if err := task.Run(); err != nil {
		// 	return c.JSON(http.StatusInternalServerError, fmt.Sprintf("task failed: %s", err.Error()))
		// }

		results, err := task.Output()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("task output failed: %s", err.Error()))
		}

		return c.Blob(http.StatusOK, "image/png", results)
	})

	e.GET("/bitmap", func(c echo.Context) error {

		source := c.QueryParam("s")
		height := c.QueryParam("h")
		width := c.QueryParam("w")
		offset := c.QueryParam("o")

		h, err := strconv.ParseUint(height, 10, 16)
		if err != nil {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("error parsing height param: %s", err.Error()))
		}
		w, err := strconv.ParseUint(width, 10, 16)
		if err != nil {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("error parsing width param: %s", err.Error()))
		}
		o, err := strconv.ParseUint(offset, 10, 16)
		if err != nil {
			// dont error - this arg is optional
			// return c.JSON(http.StatusBadRequest, fmt.Sprintf("error parsing offset param: %s", err.Error()))
		}

		task := NewBitmapTask(source, uint(w), uint(h), time.Duration(o)*time.Second)

		// if err := task.Run(); err != nil {
		// 	return c.JSON(http.StatusInternalServerError, fmt.Sprintf("task failed: %s", err.Error()))
		// }

		results, err := task.Output()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("task output failed: %s", err.Error()))
		}

		return c.Blob(http.StatusOK, "image/png", results)
	})

	e.GET("/pixels", func(c echo.Context) error {

		source := c.QueryParam("s")
		height := c.QueryParam("h")
		width := c.QueryParam("w")
		offset := c.QueryParam("o")

		h, err := strconv.ParseUint(height, 10, 16)
		if err != nil {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("error parsing height param: %s", err.Error()))
		}
		w, err := strconv.ParseUint(width, 10, 16)
		if err != nil {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("error parsing width param: %s", err.Error()))
		}
		o, err := strconv.ParseUint(offset, 10, 16)
		if err != nil {
			// dont error - this arg is optional
			// return c.JSON(http.StatusBadRequest, fmt.Sprintf("error parsing offset param: %s", err.Error()))
		}

		task := NewScaleImageTask(source, uint(w), uint(h), time.Duration(o)*time.Second)

		// if err := task.Run(); err != nil {
		// 	return c.JSON(http.StatusInternalServerError, fmt.Sprintf("task failed: %s", err.Error()))
		// }

		results, err := task.Output()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("task output failed: %s", err.Error()))
		}

		pixels, err := getPixels(bytes.NewReader(results))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("pixel parsing failed: %s", err.Error()))
		}

		blob := ""
		for _, a := range pixels {
			line := ""
			for _, b := range a {
				if line != "" {
					line += ", "
				}
				line = line + b.String()
			}
			blob = blob + line + "\n"
		}

		return c.Blob(http.StatusOK, "text/plain", []byte(blob))
	})
}
