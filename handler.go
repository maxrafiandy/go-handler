package handler

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"net/http"
	"net/url"

	"os"
	"strconv"
	"strings"

	"github.com/gorilla/schema"
)

// errorImage render an error image and send it as response
func errorImage(w http.ResponseWriter) {

	_, err := os.Open(os.Getenv("ERROR_IMAGE"))
	if err == nil {
		WriteImage(w, noImagePath)
		return
	}

	// Error image is a JPEG draw
	w.Header().Set(contentType, imageJPG)

	// create new image object
	img := image.NewRGBA(image.Rect(0, 0, 240, 240))
	blue := color.RGBA{255, 255, 255, 0}

	// // draw the image
	draw.Draw(img, img.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)

	// image buffer
	buffer := new(bytes.Buffer)

	// encode the buffer
	jpeg.Encode(buffer, img, nil)

	// write response
	w.Header().Set(contentLength, strconv.Itoa(len(buffer.Bytes())))
	w.WriteHeader(http.StatusOK)
	w.Write(buffer.Bytes())
}

// WriteImage send response as an image
func WriteImage(w http.ResponseWriter, path string) {
	// inner function for failure action
	fail := func(data interface{}) {
		errorImage(w)
		Logger(data)
	}
	var fimg image.Image
	var img *os.File
	var err error

	img, err = os.Open(path)
	if err != nil {
		fail(err)
		return
	}
	defer img.Close()

	// get image data and type (extension)
	_, itype, _ := image.Decode(img)

	// We only need this because we already read from the file
	// We have to reset the file pointer back to beginning
	img.Seek(0, 0)

	// buffer image
	bimg := new(bytes.Buffer)

	// the image should decode in
	// corresponding format
	switch {
	// if image is a JPG/JPEEG
	case itype == "jpg", itype == "jpeg":
		w.Header().Set(contentType, imageJPG)
		fimg, err = jpeg.Decode(img)
		jpeg.Encode(bimg, fimg, nil)
	// if image is a PNG
	case itype == "png":
		w.Header().Set(contentType, imagePNG)
		fimg, err = png.Decode(img)
		png.Encode(bimg, fimg)
	// no match, raise error
	default:
		err = &Error{Description: decodeFail}
	}

	if err != nil {
		fail(err)
		return
	}

	w.Header().Set("Strict-Transport-Security", "max-age=31536000")
	w.Header().Set("Content-Security-Policy", "default-src 'self'")

	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(bimg.Bytes()); err != nil {
		fail(err)
		return
	}
	Logger(path)
}

// FormData parse the incoming POST body into struct
// only handle application/json and application/x-www-form-urlencoded
func FormData(form interface{}, r *http.Request) error {
	contentType := r.Header.Get(contentType)

	if strings.Contains(contentType, "application/json") {
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(form)
		return nil
	} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		r.ParseForm()
		decoder := schema.NewDecoder()
		decoder.Decode(form, r.Form)
		return nil
	} else if strings.Contains(contentType, "multipart/form-data") {
		r.ParseMultipartForm(10 << 20)
		decoder := schema.NewDecoder()
		decoder.Decode(form, r.Form)
		return nil
	} else {
		custError := &Error{
			Description: invalidContentType,
		}
		Logger(custError)
		return custError
	}
}

// DecodeURLQuery parse the incoming URL query into struct Urlq.
// Returns true if everyting went well, otherwise false.
func DecodeURLQuery(w http.ResponseWriter, v url.Values) (URLQuery, bool) {
	// get url's query
	var args URLQuery
	decoder := schema.NewDecoder()
	decoder.Decode(&args, v)
	if err := args.Validate(); err != nil {
		response(w, MessageBadRequest, err, http.StatusBadRequest)
		return args, false
	}

	return args, true
}

func response(w http.ResponseWriter, message string, data interface{}, status int) interface{} {
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	response := Response{Message: message, Data: data}
	encoder.Encode(response)
	Logger(response)

	// return data
	return data
}
