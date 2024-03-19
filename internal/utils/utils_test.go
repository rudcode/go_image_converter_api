package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertImageSizeEqual(t *testing.T, inBuf io.Reader, width uint16, height uint16) {
	assert := assert.New(t)

	ansWidth, ansHeight, _err := GetImageSize(inBuf)
	assert.NoError(_err, "Failed to get image size")
	assert.Equal(width, ansWidth, fmt.Sprintf("got width %d, want %d", ansWidth, width))
	assert.Equal(height, ansHeight, fmt.Sprintf("got height %d, want %d", ansHeight, height))
}

func AssertImageFormatEqual(t *testing.T, inBuf io.Reader, format string) {
	assert := assert.New(t)

	ansFormat, err := GetImageFormat(inBuf)
	assert.NoError(err, "Failed to get image format")
	assert.Equal(format, ansFormat, fmt.Sprintf("got %s, want %s", ansFormat, format))
}

func TestMapfloat64(t *testing.T) {
	assert := assert.New(t)

	var tests = []struct {
		x      float64
		inMin  float64
		inMax  float64
		outMin float64
		outMax float64
		want   float64
	}{
		{0, 0, 1, 0, 100, 0},
		{0.25, 0, 1, 0, 100, 25},
		{0.50, 0, 1, 0, 100, 50},
		{0.75, 0, 1, 0, 100, 75},
		{1, 0, 1, 0, 100, 100},
		{0, 0, 1, 0, -100, 0},
		{0.25, 0, 1, 0, -100, -25},
		{0.50, 0, 1, 0, -100, -50},
		{0.75, 0, 1, 0, -100, -75},
		{1, 0, 1, 0, 100, 100},
		{0, 0, 100, 0, 1, 0},
		{25, 0, 100, 0, 1, 0.25},
		{50, 0, 100, 0, 1, 0.50},
		{75, 0, 100, 0, 1, 0.75},
		{100, 0, 100, 0, 1, 1},
		{1, 0, 1, 0, 100, 100},
		{0, 0, -100, 0, 1, 0},
		{25, 0, -100, 0, 1, -0.25},
		{50, 0, -100, 0, 1, -0.50},
		{75, 0, -100, 0, 1, -0.75},
		{100, 0, -100, 0, 1, -1},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf(
			"TestMapFloat64:x=%.2f,inMin=%.2f,inMax=%.2f,outMin=%.2f,outMax=%.2f,want=%.2f",
			tt.x, tt.inMin, tt.inMax, tt.outMin, tt.outMax, tt.want,
		), func(t *testing.T) {
			ans := Mapfloat64(tt.x, tt.inMin, tt.inMax, tt.outMin, tt.outMax)
			assert.Equal(tt.want, ans, fmt.Sprintf("got %.2f, want %.2f", ans, tt.want))
		})
	}
}

func TestGetImageFormat(t *testing.T) {
	assert := assert.New(t)

	var tests = []struct {
		fileName   string
		wantFormat string
	}{
		{"../../test/data/test_1000x1000.bmp", "bmp"},
		{"../../test/data/test_1000x1000.jpg", "mjpeg"},
		{"../../test/data/test_1000x1000.png", "png"},
		{"../../test/data/test_1000x1000.webp", "webp"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf(
			"TestGetImageFormat %s",
			tt.wantFormat,
		), func(t *testing.T) {
			inBuf, err := os.Open(tt.fileName)
			assert.NoError(err, fmt.Sprintf("Failed to open file: %s", tt.fileName))
			defer inBuf.Close()
			AssertImageFormatEqual(t, inBuf, tt.wantFormat)
		})
	}
}

func TestGetImageSize(t *testing.T) {
	assert := assert.New(t)

	var tests = []struct {
		fileName   string
		wantWidth  uint16
		wantHeight uint16
	}{
		{"../../test/data/test_1000x1000.png", 1000, 1000},
		{"../../test/data/test_1000x625.png", 1000, 625},
		{"../../test/data/test_625x1000.png", 625, 1000},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf(
			"TestGetImageSize %dx%d",
			tt.wantWidth, tt.wantHeight,
		), func(t *testing.T) {
			inBuf, err := os.Open(tt.fileName)
			assert.NoError(err, fmt.Sprintf("Failed to open file: %s", tt.fileName))
			defer inBuf.Close()
			AssertImageSizeEqual(t, inBuf, tt.wantWidth, tt.wantHeight)
		})
	}
}

func TestConvertPngToJpeg(t *testing.T) {
	assert := assert.New(t)

	var tests = []struct {
		fileName  string
		width     uint16
		height    uint16
		wantError bool
	}{
		{"../../test/data/test_1000x1000.bmp", 1000, 1000, true},
		{"../../test/data/test_1000x1000.jpg", 1000, 1000, true},
		{"../../test/data/test_1000x1000.png", 1000, 1000, false},
		{"../../test/data/test_1000x1000.webp", 1000, 1000, true},
		{"../../test/data/test_1000x625.png", 1000, 625, false},
		{"../../test/data/test_625x1000.png", 625, 1000, false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf(
			"TestConvertPngToJpeg %s",
			tt.fileName,
		), func(t *testing.T) {
			inBuf, err := os.Open(tt.fileName)
			assert.NoError(err, fmt.Sprintf("Failed to open file: %s", tt.fileName))
			defer inBuf.Close()

			outBuf := bytes.NewBuffer(nil)
			err = ConvertPngToJpeg(inBuf, outBuf)
			assert.Equal(err != nil, tt.wantError, fmt.Sprintf("got %s, want error %t", err, tt.wantError))

			if err == nil {
				outBufReader := bytes.NewReader(outBuf.Bytes())
				AssertImageSizeEqual(t, outBufReader, tt.width, tt.height)

				outBufReader.Seek(0, 0)
				AssertImageFormatEqual(t, outBufReader, "mjpeg")
			}
		})
	}
}

func TestResizeImage(t *testing.T) {
	assert := assert.New(t)

	var tests = []struct {
		fileName  string
		format    string
		width     uint16
		height    uint16
		wantError bool
	}{
		{"../../test/data/test_1000x1000.bmp", "bmp", 100, 200, false},
		{"../../test/data/test_1000x1000.bmp", "bmp", 200, 100, false},
		{"../../test/data/test_1000x1000.bmp", "bmp", 0, 100, true},
		{"../../test/data/test_1000x1000.bmp", "bmp", 100, 0, true},
		{"../../test/data/test_1000x1000.bmp", "bmp", ResizeMaxWidth + 1, 100, true},
		{"../../test/data/test_1000x1000.bmp", "bmp", 100, ResizeMaxHeight + 1, true},
		{"../../test/data/test_1000x1000.jpg", "mjpeg", 100, 200, false},
		{"../../test/data/test_1000x1000.jpg", "mjpeg", 200, 100, false},
		{"../../test/data/test_1000x1000.jpg", "mjpeg", 0, 100, true},
		{"../../test/data/test_1000x1000.jpg", "mjpeg", 100, 0, true},
		{"../../test/data/test_1000x1000.jpg", "mjpeg", ResizeMaxWidth + 1, 100, true},
		{"../../test/data/test_1000x1000.jpg", "mjpeg", 100, ResizeMaxHeight + 1, true},
		{"../../test/data/test_1000x1000.png", "png", 100, 200, false},
		{"../../test/data/test_1000x1000.png", "png", 200, 100, false},
		{"../../test/data/test_1000x1000.png", "png", 0, 100, true},
		{"../../test/data/test_1000x1000.png", "png", 100, 0, true},
		{"../../test/data/test_1000x1000.png", "png", ResizeMaxWidth + 1, 100, true},
		{"../../test/data/test_1000x1000.png", "png", 100, ResizeMaxHeight + 1, true},
		{"../../test/data/test_1000x1000.webp", "webp", 200, 200, false},
		{"../../test/data/test_1000x1000.webp", "webp", 100, 100, false},
		{"../../test/data/test_1000x1000.webp", "webp", 0, 100, true},
		{"../../test/data/test_1000x1000.webp", "webp", 100, 0, true},
		{"../../test/data/test_1000x1000.webp", "webp", ResizeMaxWidth + 1, 100, true},
		{"../../test/data/test_1000x1000.webp", "webp", 100, ResizeMaxHeight + 1, true},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf(
			"TestResizeImage %dx%d",
			tt.width, tt.height,
		), func(t *testing.T) {
			inBuf, err := os.Open(tt.fileName)
			assert.NoError(err, fmt.Sprintf("Failed to open file: %s", tt.fileName))
			defer inBuf.Close()

			outBuf := bytes.NewBuffer(nil)
			err = ResizeImage(inBuf, tt.format, tt.width, tt.height, outBuf)
			assert.Equal(err != nil, tt.wantError, fmt.Sprintf("got %s, want error %t", err, tt.wantError))

			if err == nil {
				outBufReader := bytes.NewReader(outBuf.Bytes())
				AssertImageSizeEqual(t, outBufReader, tt.width, tt.height)

				outBufReader.Seek(0, 0)
				AssertImageFormatEqual(t, outBufReader, tt.format)
			}
		})
	}
}

func TestCompressImage(t *testing.T) {
	assert := assert.New(t)

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
			inBuf, err := os.Open(tt.fileName)
			assert.NoError(err, fmt.Sprintf("Failed to open file: %s", tt.fileName))
			defer inBuf.Close()

			outBuf := bytes.NewBuffer(nil)
			err = CompressImage(inBuf, tt.format, tt.compressionLevel, outBuf)
			assert.Error(err, fmt.Sprintf(
				"want error, format:%s, level:%d",
				tt.format, tt.compressionLevel,
			))
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

			inBuf, err := os.Open(tt.fileName)
			assert.NoError(err, fmt.Sprintf("Failed to open file: %s", tt.fileName))
			defer inBuf.Close()

			for i := 0; i < 3; i++ {
				inBuf.Seek(0, 0)
				outBuf := bytes.NewBuffer(nil)
				err = CompressImage(inBuf, tt.format, compressionLevel[i], outBuf)
				assert.NoError(err, fmt.Sprintf("Failed to compress image: %s", tt.fileName))

				outBufReader := bytes.NewReader(outBuf.Bytes())
				AssertImageSizeEqual(t, outBufReader, tt.width, tt.height)

				outBufReader.Seek(0, 0)
				AssertImageFormatEqual(t, outBufReader, tt.format)
				outBufSizes[i] = outBuf.Len()
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
