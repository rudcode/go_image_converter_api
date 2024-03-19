package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/rudcode/go_image_converter_api/internal/utils"
	"github.com/stretchr/testify/assert"
)

func AssertImageSizeEqual(t *testing.T, inBuf io.Reader, width uint16, height uint16) {
	assert := assert.New(t)

	ansWidth, ansHeight, _err := utils.GetImageSize(inBuf)
	assert.NoError(_err, "Failed to get image size")
	assert.Equal(width, ansWidth, fmt.Sprintf("got width %d, want %d", ansWidth, width))
	assert.Equal(height, ansHeight, fmt.Sprintf("got height %d, want %d", ansHeight, height))
}

func AssertImageFormatEqual(t *testing.T, inBuf io.Reader, format string) {
	assert := assert.New(t)

	ansFormat, err := utils.GetImageFormat(inBuf)
	assert.NoError(err, "Failed to get image format")
	assert.Equal(format, ansFormat, fmt.Sprintf("got %s, want %s", ansFormat, format))
}

func TestConvertPngToJpeg(t *testing.T) {
	assert := assert.New(t)
	router := setupRouter()

	var tests = []struct {
		fileName string
		width    uint16
		height   uint16
		wantCode int
	}{
		{"../../test/data/test_1000x1000.bmp", 1000, 1000, http.StatusBadRequest},
		{"../../test/data/test_1000x1000.jpg", 1000, 1000, http.StatusBadRequest},
		{"../../test/data/test_1000x1000.png", 1000, 1000, http.StatusOK},
		{"../../test/data/test_1000x1000.webp", 1000, 1000, http.StatusBadRequest},
		{"../../test/data/test_1000x625.png", 1000, 625, http.StatusOK},
		{"../../test/data/test_625x1000.png", 625, 1000, http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf(
			"TestConvertPngToJpeg %s",
			tt.fileName,
		), func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			multipartWriter := multipart.NewWriter(body)

			// Create file form
			fileForm, err := multipartWriter.CreateFormFile("file", tt.fileName)
			assert.NoError(err)
			fileBuf, err := os.Open(tt.fileName)
			assert.NoError(err, fmt.Sprintf("Failed to open file: %s", tt.fileName))
			defer fileBuf.Close()
			_, err = io.Copy(fileForm, fileBuf)
			assert.NoError(err)

			assert.NoError(multipartWriter.Close())

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/convert_png_to_jpeg", body)
			assert.NoError(err)
			req.Header.Add("Content-Type", multipartWriter.FormDataContentType())
			router.ServeHTTP(res, req)

			assert.Equal(tt.wantCode, res.Code, res.Body.String())
			if res.Code == http.StatusOK {
				outBufReader := bytes.NewReader(res.Body.Bytes())
				AssertImageSizeEqual(t, outBufReader, tt.width, tt.height)

				outBufReader.Seek(0, 0)
				AssertImageFormatEqual(t, outBufReader, "mjpeg")
			}
		})
	}
}

func TestResizeImage(t *testing.T) {
	assert := assert.New(t)
	router := setupRouter()

	var tests = []struct {
		fileName string
		format   string
		width    uint16
		height   uint16
		wantCode int
	}{
		{"../../test/data/test_1000x1000.bmp", "bmp", 100, 200, http.StatusOK},
		{"../../test/data/test_1000x1000.bmp", "bmp", 200, 100, http.StatusOK},
		{"../../test/data/test_1000x1000.bmp", "bmp", 0, 100, http.StatusBadRequest},
		{"../../test/data/test_1000x1000.bmp", "bmp", 100, 0, http.StatusBadRequest},
		{"../../test/data/test_1000x1000.bmp", "bmp", utils.ResizeMaxWidth + 1, 100, http.StatusBadRequest},
		{"../../test/data/test_1000x1000.bmp", "bmp", 100, utils.ResizeMaxHeight + 1, http.StatusBadRequest},
		{"../../test/data/test_1000x1000.jpg", "mjpeg", 100, 200, http.StatusOK},
		{"../../test/data/test_1000x1000.jpg", "mjpeg", 200, 100, http.StatusOK},
		{"../../test/data/test_1000x1000.jpg", "mjpeg", 0, 100, http.StatusBadRequest},
		{"../../test/data/test_1000x1000.jpg", "mjpeg", 100, 0, http.StatusBadRequest},
		{"../../test/data/test_1000x1000.jpg", "mjpeg", utils.ResizeMaxWidth + 1, 100, http.StatusBadRequest},
		{"../../test/data/test_1000x1000.jpg", "mjpeg", 100, utils.ResizeMaxHeight + 1, http.StatusBadRequest},
		{"../../test/data/test_1000x1000.png", "png", 100, 200, http.StatusOK},
		{"../../test/data/test_1000x1000.png", "png", 200, 100, http.StatusOK},
		{"../../test/data/test_1000x1000.png", "png", 0, 100, http.StatusBadRequest},
		{"../../test/data/test_1000x1000.png", "png", 100, 0, http.StatusBadRequest},
		{"../../test/data/test_1000x1000.png", "png", utils.ResizeMaxWidth + 1, 100, http.StatusBadRequest},
		{"../../test/data/test_1000x1000.png", "png", 100, utils.ResizeMaxHeight + 1, http.StatusBadRequest},
		{"../../test/data/test_1000x1000.webp", "webp", 200, 200, http.StatusOK},
		{"../../test/data/test_1000x1000.webp", "webp", 100, 100, http.StatusOK},
		{"../../test/data/test_1000x1000.webp", "webp", 0, 100, http.StatusBadRequest},
		{"../../test/data/test_1000x1000.webp", "webp", 100, 0, http.StatusBadRequest},
		{"../../test/data/test_1000x1000.webp", "webp", utils.ResizeMaxWidth + 1, 100, http.StatusBadRequest},
		{"../../test/data/test_1000x1000.webp", "webp", 100, utils.ResizeMaxHeight + 1, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf(
			"TestResizeImage %dx%d",
			tt.width, tt.height,
		), func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			multipartWriter := multipart.NewWriter(body)

			// Create file form
			formFile, err := multipartWriter.CreateFormFile("file", tt.fileName)
			assert.NoError(err)
			fileBuf, err := os.Open(tt.fileName)
			assert.NoError(err, fmt.Sprintf("Failed to open file: %s", tt.fileName))
			defer fileBuf.Close()
			_, err = io.Copy(formFile, fileBuf)
			assert.NoError(err)

			// Create width field form
			formWidth, err := multipartWriter.CreateFormField("width")
			assert.NoError(err)
			formWidth.Write([]byte(fmt.Sprintf("%d", tt.width)))

			// Create height field form
			formHeight, err := multipartWriter.CreateFormField("height")
			assert.NoError(err)
			formHeight.Write([]byte(fmt.Sprintf("%d", tt.height)))

			assert.NoError(multipartWriter.Close())

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/resize_image", body)
			assert.NoError(err)
			req.Header.Add("Content-Type", multipartWriter.FormDataContentType())
			router.ServeHTTP(res, req)

			assert.Equal(tt.wantCode, res.Code, res.Body.String())
			if res.Code == http.StatusOK {
				outBufReader := bytes.NewReader(res.Body.Bytes())
				AssertImageSizeEqual(t, outBufReader, tt.width, tt.height)

				outBufReader.Seek(0, 0)
				AssertImageFormatEqual(t, outBufReader, tt.format)
			}
		})
	}
}

func TestCompressImage(t *testing.T) {
	assert := assert.New(t)
	router := setupRouter()

	var failTests = []struct {
		fileName         string
		format           string
		compressionLevel uint8
	}{
		{"../../test/data/test_1000x1000.jpg", "mjpeg", 0},
		{"../../test/data/test_1000x1000.jpg", "mjpeg", 6},
		{"../../test/data/test_1000x1000.png", "png", 0},
		{"../../test/data/test_1000x1000.png", "png", 6},
		{"../../test/data/test_1000x1000.webp", "webp", 0},
		{"../../test/data/test_1000x1000.webp", "webp", 6},
		{"../../test/data/test_1000x1000.bmp", "bmp", 5},
	}

	for _, tt := range failTests {
		t.Run(fmt.Sprintf(
			"TestCompressImage case fail %s,level:%d",
			tt.format, tt.compressionLevel,
		), func(t *testing.T) {
			body := bytes.NewBuffer(nil)
			multipartWriter := multipart.NewWriter(body)

			// Create file form
			formFile, err := multipartWriter.CreateFormFile("file", tt.fileName)
			assert.NoError(err)
			fileBuf, err := os.Open(tt.fileName)
			assert.NoError(err, fmt.Sprintf("Failed to open file: %s", tt.fileName))
			defer fileBuf.Close()
			_, err = io.Copy(formFile, fileBuf)
			assert.NoError(err)

			// Create compression_level field form
			formCompressionLevel, err := multipartWriter.CreateFormField("compression_level")
			assert.NoError(err)
			formCompressionLevel.Write([]byte(fmt.Sprintf("%d", tt.compressionLevel)))

			assert.NoError(multipartWriter.Close())

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/compress_image", body)
			assert.NoError(err)
			req.Header.Add("Content-Type", multipartWriter.FormDataContentType())
			router.ServeHTTP(res, req)

			assert.Equal(http.StatusBadRequest, res.Code, res.Body.String())
		})
	}

	var tests = []struct {
		fileName string
		width    uint16
		height   uint16
		format   string
	}{
		{"../../test/data/test_1000x1000.jpg", 1000, 1000, "mjpeg"},
		{"../../test/data/test_1000x1000.png", 1000, 1000, "png"},
		{"../../test/data/test_1000x1000.webp", 1000, 1000, "webp"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf(
			"TestCompressImage compare compressionLevel %s",
			tt.fileName,
		), func(t *testing.T) {
			var outBufSizes [3]int
			compressionLevel := []uint8{1, 3, 5}

			fileBuf, err := os.Open(tt.fileName)
			assert.NoError(err, fmt.Sprintf("Failed to open file: %s", tt.fileName))
			defer fileBuf.Close()

			for i := 0; i < 3; i++ {
				fileBuf.Seek(0, 0)
				body := bytes.NewBuffer(nil)
				multipartWriter := multipart.NewWriter(body)

				// Create file form
				formFile, err := multipartWriter.CreateFormFile("file", tt.fileName)
				assert.NoError(err)
				_, err = io.Copy(formFile, fileBuf)
				assert.NoError(err)

				// Create compression_level field form
				formCompressionLevel, err := multipartWriter.CreateFormField("compression_level")
				assert.NoError(err)
				formCompressionLevel.Write([]byte(fmt.Sprintf("%d", compressionLevel[i])))

				assert.NoError(multipartWriter.Close())

				res := httptest.NewRecorder()
				req, err := http.NewRequest(http.MethodPost, "/compress_image", body)
				assert.NoError(err)
				req.Header.Add("Content-Type", multipartWriter.FormDataContentType())
				router.ServeHTTP(res, req)

				assert.Equal(http.StatusOK, res.Code, res.Body.String())
				outBufReader := bytes.NewReader(res.Body.Bytes())
				AssertImageSizeEqual(t, outBufReader, tt.width, tt.height)

				outBufReader.Seek(0, 0)
				AssertImageFormatEqual(t, outBufReader, tt.format)

				outBufSizes[i] = res.Body.Len()
			}

			// Check compression level size
			assert.Greater(outBufSizes[0], outBufSizes[1], fmt.Sprintf(
				"Outbuf size not align with compression level: %d=%d, %d=%d",
				compressionLevel[0], outBufSizes[0],
				compressionLevel[1], outBufSizes[1],
			))
			assert.Greater(outBufSizes[1], outBufSizes[2], fmt.Sprintf(
				"Outbuf size not align with compression level: %d=%d, %d=%d",
				compressionLevel[1], outBufSizes[1],
				compressionLevel[2], outBufSizes[2],
			))
		})
	}
}
