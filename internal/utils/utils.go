package utils

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

const ResizeMaxWidth uint16 = 4096
const ResizeMaxHeight uint16 = 4096

var FfmpegResizeImageFormats = [...]string{
	"mjpeg",
	"png",
	"webp",
	"bmp",
}

var FfmpegCompressImageFormats = [...]string{
	"mjpeg",
	"png",
	"webp",
}

func Mapfloat64(x float64, inMin float64, inMax float64, outMin float64, outMax float64) float64 {
	return (x-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}

func GetImageFormat(inBuf io.Reader) (string, error) {
	format, err := ffmpeg.ProbeReaderWithTimeoutExec(inBuf, 0, ffmpeg.KwArgs{
		"v":            "error",
		"show_entries": "stream=codec_name",
		"of":           "default=noprint_wrappers=1:nokey=1",
	})
	return strings.TrimSuffix(format, "\r\n"), err
}

func GetImageSize(inBuf io.Reader) (uint16, uint16, error) {
	size, err := ffmpeg.ProbeReaderWithTimeoutExec(inBuf, 0, ffmpeg.KwArgs{
		"v":            "error",
		"show_entries": "stream=width,height",
		"of":           "default=noprint_wrappers=1:nokey=1",
	})
	if err != nil {
		return 0, 0, err
	}
	sizeSlice := strings.Split(strings.TrimSuffix(size, "\r\n"), "\r\n")
	width, err := strconv.ParseInt(sizeSlice[0], 10, 32)
	if err != nil {
		return 0, 0, err
	}
	height, err := strconv.ParseInt(sizeSlice[1], 10, 32)
	if err != nil {
		return 0, 0, err
	}
	return uint16(width), uint16(height), err
}

func ConvertPngToJpeg(inBuf io.ReadSeeker, outBuf io.Writer) error {
	// Check if PNG
	format, err := GetImageFormat(inBuf)
	if err != nil {
		return fmt.Errorf("can't probe file, make sure the file is valid png image: %s", err.Error())
	}
	if format != "png" {
		return fmt.Errorf("image format must be PNG")
	}
	inBuf.Seek(0, 0)

	// Convert to JPG
	err = ffmpeg.
		Input("pipe:").
		WithInput(inBuf).
		Output("pipe:", ffmpeg.KwArgs{
			"vcodec": "mjpeg",
			"f":      "image2",
		}).
		WithOutput(outBuf). //, os.Stdout).
		Silent(true).
		Run()
	if err != nil {
		return fmt.Errorf("error while transcoding: %s", err.Error())
	}
	return err
}

// ResizeImage function resize the image stored in inBuf and write the output to outBuf
// format is one of the following ("mjpeg", "png", "webp", "bmp")
// width is value between 1 to 4096 (ResizeMaxWidth)
// height is value between 1 to 4096 (ResizeMaxHeight)
func ResizeImage(inBuf io.Reader, format string, width uint16, height uint16, outBuf io.Writer) error {
	// Check width and height
	if width < 1 || width > ResizeMaxWidth {
		return fmt.Errorf("width must be positive and < %d", ResizeMaxWidth)
	}

	if height < 1 || height > ResizeMaxHeight {
		return fmt.Errorf("height must be positive and < %d", ResizeMaxHeight)
	}

	// Check format
	formatFound := false
	for _, ffmpegFormat := range FfmpegResizeImageFormats {
		if format == ffmpegFormat {
			formatFound = true
			break
		}
	}
	if !formatFound {
		return fmt.Errorf("file format %s is not supported", format)
	}

	// Resize
	err := ffmpeg.
		Input("pipe:").
		WithInput(inBuf).
		Output("pipe:", ffmpeg.KwArgs{
			"vf":     fmt.Sprintf("scale=%d:%d", width, height),
			"vcodec": format,
			"f":      "image2",
		}).
		WithOutput(outBuf). //, os.Stdout)
		Silent(true).
		Run()
	if err != nil {
		return fmt.Errorf("error while transcoding: %s", err.Error())
	}
	return err
}

// CompressImage function compress the image stored in inBuf and write the output to outBuf
// format is one of the following ("mjpeg", "png", "webp")
// compressionLevel is value between 1-5 where 1 means largest file size and 5 means smallest file size
func CompressImage(inBuf io.Reader, format string, compressionLevel uint8, outBuf io.Writer) error {
	// Check format
	formatFound := false
	for _, ffmpegFormat := range FfmpegCompressImageFormats {
		if format == ffmpegFormat {
			formatFound = true
			break
		}
	}
	if !formatFound {
		return fmt.Errorf("file format %s is not supported", format)
	}

	// Check compressionLevel
	if compressionLevel < 1 || compressionLevel > 5 {
		return fmt.Errorf("compression level must between 1 <= level <= 5")
	}

	// Compress
	outKwargs := ffmpeg.KwArgs{
		"f":      "image2",
		"vcodec": format,
	}

	if format == "png" {
		// For PNG we use compression_level that range from 1-9
		outKwargs["compression_level"] = Mapfloat64(float64(compressionLevel), 1, 5, 1, 9)
	} else if format == "webp" {
		// For PNG we use compression_level that range from 1-6
		outKwargs["compression_level"] = Mapfloat64(float64(compressionLevel), 1, 5, 1, 6)
	} else if format == "mjpeg" {
		// For JPEG we use q for quality that range from 1-31
		outKwargs["q"] = Mapfloat64(float64(compressionLevel), 1, 5, 1, 31)
	}

	err := ffmpeg.
		Input("pipe:").
		WithInput(inBuf).
		Output("pipe:", outKwargs).
		WithOutput(outBuf). //, os.Stdout).
		Silent(true).
		Run()
	if err != nil {
		return fmt.Errorf("error while transcoding: %s", err.Error())
	}
	return err
}
