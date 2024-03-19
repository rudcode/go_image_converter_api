package main

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/rudcode/go_image_converter_api/docs"
	"github.com/rudcode/go_image_converter_api/internal/utils"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type convertToJpegInputParameter struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type resizeImageInputParameter struct {
	Width  *uint16               `form:"width" binding:"required"`
	Height *uint16               `form:"height" binding:"required"`
	File   *multipart.FileHeader `form:"file" binding:"required"`
}

type compressImageInputParameter struct {
	CompressionLevel *uint8                `form:"compression_level" binding:"required"`
	File             *multipart.FileHeader `form:"file" binding:"required"`
}

type ErrorResponse struct {
	Detail string `json:"detail"`
}

// @Summary		Convert PNG to JPEG
// @Description	Convert image from PNG format to JPEG format
// @ID			convert_png_to_jpeg
// @Accept		multipart/form-data
// @Produce		json
// @Param		file	formData	file	true	"image file"
//
// @Router		/convert_png_to_jpeg [post]
func convertPngToJpeg(c *gin.Context) {
	var input convertToJpegInputParameter

	// Input verification
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Detail: err.Error(),
		})
		return
	}

	if input.File == nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Detail: "File is missing",
		})
		return
	}

	// Get file buffer
	inBuf, err := input.File.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Detail: err.Error(),
		})
		return
	}
	defer inBuf.Close()

	// Convert PNG to JPG
	outBuf := bytes.NewBuffer(nil)
	err = utils.ConvertPngToJpeg(inBuf, outBuf)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Detail: fmt.Sprintf("Error while converting: %s", err.Error()),
		})
		return
	}
	c.Data(http.StatusOK, "image/jpeg", outBuf.Bytes())
}

// @Summary		Resize image
// @Description	Resize image to the specified width and height
// @ID			resize_image
// @Accept		multipart/form-data
// @Produce		json
// @Param		file	formData	file	true	"image file"
// @Param		width	formData	uint16	true	"width"
// @Param		height	formData	uint16	true	"height"
//
// @Router		/resize_image [post]
func resizeImage(c *gin.Context) {
	var input resizeImageInputParameter

	// Input verification
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Detail: err.Error()})
		return
	}

	if input.File == nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Detail: "File is missing"})
		return
	}

	// Get file buffer
	inBuf, err := input.File.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Detail: err.Error(),
		})
		return
	}
	defer inBuf.Close()

	// Get format
	format, err := utils.GetImageFormat(inBuf)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Detail: fmt.Sprintf("Can't probe file, make sure the file is valid image: %s", err.Error()),
		})
		return
	}

	inBuf.Seek(0, 0)

	// Resize
	outBuf := bytes.NewBuffer(nil)
	err = utils.ResizeImage(inBuf, format, *input.Width, *input.Height, outBuf)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Detail: fmt.Sprintf("Error while resizing: %s", err.Error()),
		})
		return
	}

	if format == "mjpeg" {
		format = "jpeg" // return image/jpeg mimetype
	}
	c.Data(http.StatusOK, fmt.Sprintf("image/%s", format), outBuf.Bytes())
}

// @Summary		Compress image
// @Description	Compress image with specified compression level (1-5)
// @ID			compress_image
// @Accept		multipart/form-data
// @Produce		json
// @Param		file				formData	file	true	"image file"
// @Param		compression_level	formData	uint8	true	"compression level (1-5)"
//
// @Router		/compress_image [post]
func compressImage(c *gin.Context) {
	var input compressImageInputParameter

	// Input verification
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Detail: err.Error(),
		})
		return
	}

	if input.File == nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Detail: "File is missing",
		})
		return
	}

	if *input.CompressionLevel < 1 || *input.CompressionLevel > 5 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Detail: "Compression Level must be 1 <= level <= 5",
		})
		return
	}

	// Get file buffer
	inBuf, err := input.File.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Detail: err.Error(),
		})
		return
	}
	defer inBuf.Close()

	// Get format
	format, err := utils.GetImageFormat(inBuf)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Detail: fmt.Sprintf("Can't probe file, make sure the file is valid image: %s", err.Error()),
		})
		return
	}
	inBuf.Seek(0, 0)

	// Compress
	outBuf := bytes.NewBuffer(nil)
	err = utils.CompressImage(inBuf, format, *input.CompressionLevel, outBuf)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Detail: fmt.Sprintf("Error while compressing: %s", err.Error()),
		})
		return
	}

	if format == "mjpeg" {
		format = "jpeg" // return image/jpeg mimetype
	}
	c.Data(http.StatusOK, fmt.Sprintf("image/%s", format), outBuf.Bytes())
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.StaticFile("/favicon.ico", "./favicon.ico")
	r.POST("/convert_png_to_jpeg", convertPngToJpeg)
	r.POST("/resize_image", resizeImage)
	r.POST("/compress_image", compressImage)

	// swagger
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}

// @title	Go Image Converter API
func main() {
	r := setupRouter()
	r.Run("localhost:8000")
	fmt.Println("Starting service")
}
