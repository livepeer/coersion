package main

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"sync"
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

func imageMatch(wg *sync.WaitGroup, pixelsA, pixelsB [][]Pixel, v uint64, pass, fail *int) {
	defer wg.Done()
	*pass = 0
	*fail = 0
	for x, a := range pixelsA {
		for y, b := range a {
			if math.Abs(float64(b.R-pixelsB[x][y].R)) < float64(v) &&
				math.Abs(float64(b.G-pixelsB[x][y].G)) < float64(v) &&
				math.Abs(float64(b.B-pixelsB[x][y].B)) < float64(v) &&
				math.Abs(float64(b.A-pixelsB[x][y].A)) < float64(v) {
				*pass++
			} else {
				*fail++
			}
		}
	}
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
		// dont error - this arg is optional
		o, _ := strconv.ParseUint(offset, 10, 16)

		task := NewScaleImageTask(source, uint(w), uint(h), time.Duration(o)*time.Second)

		results, err := StdoutPipe(task)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("task output failed: %s", err.Error()))
		}

		if err := Start(task); err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("task failed: %s", err.Error()))
		}

		return c.Stream(http.StatusOK, "image/png", results)
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
		// dont error - this arg is optional
		o, _ := strconv.ParseUint(offset, 10, 16)

		task := NewBitmapTask(source, uint(w), uint(h), time.Duration(o)*time.Second)

		results, err := StdoutPipe(task)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("task output failed: %s", err.Error()))
		}

		if err := Start(task); err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("task failed: %s", err.Error()))
		}

		return c.Stream(http.StatusOK, "image/bmp", results)
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
		// dont error - this arg is optional
		o, _ := strconv.ParseUint(offset, 10, 16)

		task := NewScaleImageTask(source, uint(w), uint(h), time.Duration(o)*time.Second)

		results, err := StdoutPipe(task)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("task output failed: %s", err.Error()))
		}

		if err := Start(task); err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("task failed: %s", err.Error()))
		}

		pixels, err := getPixels(results)
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

	e.GET("/match2", func(c echo.Context) error {

		source0 := c.QueryParam("s0")
		source1 := c.QueryParam("s1")
		height := c.QueryParam("h")
		width := c.QueryParam("w")
		offset := c.QueryParam("o")
		variation := c.QueryParam("v")

		h, err := strconv.ParseUint(height, 10, 16)
		if err != nil {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("error parsing height param: %s", err.Error()))
		}
		w, err := strconv.ParseUint(width, 10, 16)
		if err != nil {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("error parsing width param: %s", err.Error()))
		}
		// dont error - this arg is optional
		o, _ := strconv.ParseUint(offset, 10, 16)

		v, err := strconv.ParseUint(variation, 10, 16)
		if err != nil {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("error parsing variation param: %s", err.Error()))
		}

		task0 := NewScaleImageTask(source0, uint(w), uint(h), time.Duration(o)*time.Second)
		task1 := NewScaleImageTask(source1, uint(w), uint(h), time.Duration(o)*time.Second)

		img0, err := StdoutPipe(task0)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("task0 output failed: %s", err.Error()))
		}

		img1, err := StdoutPipe(task1)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("task1 output failed: %s", err.Error()))
		}

		if err := Start(task0); err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("task0 failed: %s", err.Error()))
		}

		if err := Start(task1); err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("task1 failed: %s", err.Error()))
		}

		pixels0, err := getPixels(img0)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("pixel parsing for task0 failed: %s", err.Error()))
		}

		pixels1, err := getPixels(img1)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("pixel parsing for task1 failed: %s", err.Error()))
		}

		pass := 0
		fail := 0
		// we don't need a waitgroup here, but we want to use a common function, so just pass one in anyway
		var wg sync.WaitGroup
		wg.Add(1)
		go imageMatch(&wg, pixels0, pixels1, v, &pass, &fail)
		wg.Wait()
		blob := fmt.Sprintf("%v%% match", 100*float64(pass)/float64(pass+fail))

		return c.Blob(http.StatusOK, "text/plain", []byte(blob))
	})

	e.GET("/match3", func(c echo.Context) error {

		source0 := c.QueryParam("s0")
		source1 := c.QueryParam("s1")
		source2 := c.QueryParam("s2")
		height := c.QueryParam("h")
		width := c.QueryParam("w")
		offset := c.QueryParam("o")
		variation := c.QueryParam("v")

		h, err := strconv.ParseUint(height, 10, 16)
		if err != nil {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("error parsing height param: %s", err.Error()))
		}
		w, err := strconv.ParseUint(width, 10, 16)
		if err != nil {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("error parsing width param: %s", err.Error()))
		}
		// dont error - this arg is optional
		o, _ := strconv.ParseUint(offset, 10, 16)

		v, err := strconv.ParseUint(variation, 10, 16)
		if err != nil {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("error parsing variation param: %s", err.Error()))
		}

		task0 := NewScaleImageTask(source0, uint(w), uint(h), time.Duration(o)*time.Second)
		task1 := NewScaleImageTask(source1, uint(w), uint(h), time.Duration(o)*time.Second)
		task2 := NewScaleImageTask(source2, uint(w), uint(h), time.Duration(o)*time.Second)

		img0, err := StdoutPipe(task0)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("task0 output failed: %s", err.Error()))
		}

		img1, err := StdoutPipe(task1)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("task1 output failed: %s", err.Error()))
		}

		img2, err := StdoutPipe(task2)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("task2 output failed: %s", err.Error()))
		}

		if err := Start(task0); err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("task0 failed: %s", err.Error()))
		}

		if err := Start(task1); err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("task1 failed: %s", err.Error()))
		}

		if err := Start(task2); err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("task2 failed: %s", err.Error()))
		}

		pixels0, err := getPixels(img0)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("pixel parsing for task0 failed: %s", err.Error()))
		}

		pixels1, err := getPixels(img1)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("pixel parsing for task1 failed: %s", err.Error()))
		}

		pass0 := 0
		fail0 := 0
		pass1 := 0
		fail1 := 0
		pass2 := 0
		fail2 := 0

		var wg sync.WaitGroup
		wg.Add(3)

		go imageMatch(&wg, pixels0, pixels1, v, &pass0, &fail0)

		// this is a blocking call, so start the above async job before this one.
		pixels2, err := getPixels(img2)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("pixel parsing for task2 failed: %s", err.Error()))
		}

		// these next two depend on pixels2, above
		go imageMatch(&wg, pixels1, pixels2, v, &pass1, &fail1)
		go imageMatch(&wg, pixels2, pixels0, v, &pass2, &fail2)

		wg.Wait()
		blob0 := fmt.Sprintf("s0 & s1: %v%% match\n", 100*float64(pass0)/float64(pass0+fail0))
		blob1 := fmt.Sprintf("s1 & s2: %v%% match\n", 100*float64(pass1)/float64(pass1+fail1))
		blob2 := fmt.Sprintf("s2 & s0: %v%% match\n", 100*float64(pass2)/float64(pass2+fail2))

		return c.Blob(http.StatusOK, "text/plain", []byte(blob0+blob1+blob2))
	})
}
