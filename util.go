package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"net/url"

	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

// RestHandlers interface
type RestHandlers interface {
	Post() interface{}
	Get() interface{}
	GetID(string) interface{}
	Put() interface{}
	PutID(string) interface{}
	Patch() interface{}
	PatchID(string) interface{}
	Delete() interface{}
	DeleteID(string) interface{}
	reset(http.ResponseWriter, *http.Request)
}

// REST maps router to appropriate methods
func REST(rest RestHandlers, ctx *Context) interface{} {
	id, withid := mux.Vars(ctx.Request)[id]

	rest.reset(ctx.Writer, ctx.Request)

	if !withid { // route to /
		switch ctx.Request.Method {
		case get:
			return rest.Get()
		case post:
			return rest.Post()
		case put:
			return rest.Put()
		case patch:
			return rest.Patch()
		case delete:
			return rest.Delete()
		}
	} else { // route to /{id:[0-9]+}
		switch ctx.Request.Method {
		case get:
			return rest.GetID(id)
		case put:
			return rest.PutID(id)
		case patch:
			return rest.PatchID(id)
		case delete:
			return rest.DeleteID(id)
		}
	}

	return ctx.MethodNotAllowed()
}

// errorImage render an error image and send it as response
func errorImage(w http.ResponseWriter) interface{} {
	// check for ERROR_IMAGE. If it exist return default
	// no image file insted of 404-not found
	file, err := os.Open(os.Getenv("ERROR_IMAGE"))
	if err != nil {
		return response(w, "No image", nil, http.StatusNotFound)
	}
	defer file.Close()
	return WriteImage(os.Getenv("ERROR_IMAGE"), w)
}

// WriteImage send response as an image
func WriteImage(path string, w http.ResponseWriter) error {
	// inner function for failure action
	fail := func(err error) error {
		errorImage(w)
		Logger(err)
		return DescError(err)
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
		return fail(err)
	}

	w.Header().Set("Strict-Transport-Security", "max-age=31536000")
	w.Header().Set("Content-Security-Policy", "default-src 'self'")

	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(bimg.Bytes()); err != nil {
		return fail(err)
	}
	Logger(path)
	return nil
}

// FormData parse the incoming POST body into struct
// only handle application/json and application/x-www-form-urlencoded
func FormData(form interface{}, r *http.Request) error {
	var (
		err           error
		contentType   string
		jsonDecoder   *json.Decoder
		schemaDecoder *schema.Decoder
	)
	contentType = r.Header.Get(contentType)

	if strings.Contains(contentType, "application/json") {
		jsonDecoder = json.NewDecoder(r.Body)
		jsonDecoder.Decode(form)
	} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		r.ParseForm()
		schemaDecoder = schema.NewDecoder()
		schemaDecoder.Decode(form, r.Form)
	} else if strings.Contains(contentType, "multipart/form-data") {
		r.ParseMultipartForm(10 << 20)
		schemaDecoder = schema.NewDecoder()
		schemaDecoder.Decode(form, r.Form)
	} else {
		err = &Error{
			Description: invalidContentType,
		}
		Logger(err)
	}
	return err
}

// DecodeURLQuery parse the incoming URL query into struct Urlq.
// Returns true if everyting went well, otherwise false.
func DecodeURLQuery(w http.ResponseWriter, v url.Values) (args URLQuery, err error) {
	// get url's query
	decoder := schema.NewDecoder()
	decoder.Decode(&args, v)
	if err := args.Validate(); err != nil {
		return args, DescError(err)
	}

	return args, nil
}

// Write custom message response
func Write(w http.ResponseWriter, message string, data interface{}, status int) interface{} {
	return response(w, message, data, status)
}

// DescError returns handler.Error struct with generated
// string err.Error() as its description
func DescError(err error) *Error {
	return &Error{
		Description: err.Error(),
		Errors:      err,
	}
}

// response returns JSON encoded data
func response(w http.ResponseWriter, message string, data interface{}, status int) interface{} {
	var (
		response            Response
		responseContentType string
		jsonEncoder         *json.Encoder
	)

	responseContentType = w.Header().Get(contentType)

	if strings.Contains(responseContentType, "application/json") {
		// write collected headers with status
		w.WriteHeader(status)

		// encode data to json format
		jsonEncoder = json.NewEncoder(w)
		response = Response{Message: message, Data: data}
		jsonEncoder.Encode(response)
	} else {
		var (
			buffer bytes.Buffer
			err    error
		)
		_, err = fmt.Fprintf(&buffer, "%+v", data)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong!"))
		} else {
			w.WriteHeader(status)
			w.Write(buffer.Bytes())
		}
	}
	// write log
	Logger(response)

	// return data
	return data
}
