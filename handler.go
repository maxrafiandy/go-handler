package handler

import (
	"bytes"
	"encoding/json"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"net/url"

	"os"
	"strings"

	"github.com/gorilla/schema"
)

// errorImage render an error image and send it as response
func errorImage(w http.ResponseWriter) {

	_, err := os.Open(os.Getenv("ERROR_IMAGE"))
	if err == nil {
		WriteImage(noImagePath, w)
		return
	}

	response(w, MessagePageNotFound, nil, http.StatusNotFound)
}

// WriteImage send response as an image
func WriteImage(path string, w http.ResponseWriter) error {
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
		return err
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
		return err
	}

	w.Header().Set("Strict-Transport-Security", "max-age=31536000")
	w.Header().Set("Content-Security-Policy", "default-src 'self'")

	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(bimg.Bytes()); err != nil {
		fail(err)
		return err
	}
	Logger(path)
	return nil
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
func DecodeURLQuery(w http.ResponseWriter, v url.Values) (args URLQuery, err error) {
	// get url's query
	decoder := schema.NewDecoder()
	decoder.Decode(&args, v)
	if err := args.Validate(); err != nil {
		return args, err
	}

	return args, nil
}

// Write custom message response
func Write(w http.ResponseWriter, message string, data interface{}, status int) interface{} {
	return response(w, message, data, status)
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
