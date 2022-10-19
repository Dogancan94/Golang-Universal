package toolkit

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

func TestTools_PushJSONToRemote(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString("ok")),
			Header:     make(http.Header),
		}
	})

	var testTools Tools
	var foo struct {
		Bar string `json:"bar"`
	}
	foo.Bar = "bar"
	_, _, err := testTools.PushJSONToRemote("http://example.com/some/path", foo, client)
	if err != nil {
		t.Error("failed to call remote url:", err)
	}
}

func TestToolkit_CreateRandomString(t *testing.T) {
	var testTools Tools
	randomString := testTools.CreateRandomString(10)
	if len(randomString) != 10 {
		t.Error("wrong length random string return")
	}
}

var uploadTests = []struct {
	name          string
	allowedType   []string
	renameFile    bool
	errorExpected bool
}{
	{name: "allowed no rename", allowedType: []string{"image/jpeg", "image/png", "image/gif", "image/jfif"}, renameFile: false, errorExpected: false},
	{name: "allowed rename", allowedType: []string{"image/jpeg", "image/png", "image/gif", "image/jfif"}, renameFile: true, errorExpected: false},
	{name: "not allowed", allowedType: []string{"image/jpeg", "image/gif", "image/jfif"}, renameFile: false, errorExpected: true},
}

func TestTools_UploadFiles(t *testing.T) {
	for _, e := range uploadTests {
		pipeReader, pipeWriter := io.Pipe()
		writer := multipart.NewWriter(pipeWriter)
		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer writer.Close()
			defer wg.Done()

			part, err := writer.CreateFormFile("file", "./test/img.png")
			if err != nil {
				t.Error(err)
			}

			file, err := os.Open("./test/img.png")
			if err != nil {
				t.Error(err)
			}

			defer file.Close()
			img, _, err := image.Decode(file)
			if err != nil {
				t.Error("error decoding image", err)
			}

			err = png.Encode(part, img)
			if err != nil {
				t.Error("error decoding image", err)
			}
		}()

		request := httptest.NewRequest("POST", "/", pipeReader)
		request.Header.Add("Content-Type", writer.FormDataContentType())

		var testTools Tools
		testTools.AllowedFileTypes = e.allowedType
		uploadedFiles, err := testTools.UploadFiles(request, "./test/uploads", e.renameFile)
		if err != nil && !e.errorExpected {
			t.Error(err)
		}

		if !e.errorExpected {
			if _, err := os.Stat(fmt.Sprintf("./test/uploads/%s", uploadedFiles[0].NewFileName)); os.IsNotExist(err) {
				t.Errorf("%s: expected file to exist: %s", e.name, err.Error())
			}

			//clean
			_ = os.Remove(fmt.Sprintf("./test/uploads/%s", uploadedFiles[0].NewFileName))
		}

		if !e.errorExpected && err != nil {
			t.Errorf("%s: error expected but none received", e.name)
		}

		wg.Wait()
	}
}

func TestTools_UploadFile(t *testing.T) {
	pipeReader, pipeWriter := io.Pipe()
	writer := multipart.NewWriter(pipeWriter)

	go func() {
		defer writer.Close()

		part, err := writer.CreateFormFile("file", "./test/img.png")
		if err != nil {
			t.Error(err)
		}

		file, err := os.Open("./test/img.png")
		if err != nil {
			t.Error(err)
		}

		defer file.Close()
		img, _, err := image.Decode(file)
		if err != nil {
			t.Error("error decoding image", err)
		}

		err = png.Encode(part, img)
		if err != nil {
			t.Error("error decoding image", err)
		}
	}()

	request := httptest.NewRequest("POST", "/", pipeReader)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	var testTools Tools

	uploadedFile, err := testTools.UploadOneFile(request, "./test/uploads", true)
	if err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(fmt.Sprintf("./test/uploads/%s", uploadedFile.NewFileName)); os.IsNotExist(err) {
		t.Errorf("expected file to exist: %s", err.Error())
	}

	//clean
	_ = os.Remove(fmt.Sprintf("./test/uploads/%s", uploadedFile.NewFileName))

}

func TestTools_CreateDirIfNotExist(t *testing.T) {
	var testTools Tools

	err := testTools.CreateDirectoryIfNotExist("./test/myDir")
	if err != nil {
		t.Error(err)
	}

	err = testTools.CreateDirectoryIfNotExist("./test/myDir")
	if err != nil {
		t.Error(err)
	}

	_ = os.Remove("./test/myDir")
}

var slugTests = []struct {
	name          string
	slugifyString string
	expected      string
	errorExpected bool
}{
	{name: "valid string", slugifyString: "now is the time", expected: "now-is-the-time", errorExpected: false},
	{name: "empty string", slugifyString: "", expected: "", errorExpected: true},
	{name: "complex string", slugifyString: "Now is the time for all Good men! + fish & such &^123",
		expected: "now-is-the-time-for-all-good-men-fish-such-123", errorExpected: false},
	{name: "russian string", slugifyString: "Привет, мир",
		expected: "", errorExpected: true},
	{name: "russian string and roman", slugifyString: "Привет, мир hello world",
		expected: "hello-world", errorExpected: false},
}

func TestTools_Slugify(t *testing.T) {
	var testTools Tools

	for _, e := range slugTests {
		slug, err := testTools.Slugify(e.slugifyString)
		if err != nil && !e.errorExpected {
			t.Errorf("%s: error received when none expected: %s", e.name, err.Error())
		}

		if !e.errorExpected && slug != e.expected {
			t.Errorf("%s: wrong slug returned; expected %s but got %s", e.name, e.expected, slug)
		}
	}
}

func TestTools_DownloadStaticFile(t *testing.T) {
	recorderRequest := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)

	var testTools Tools
	testTools.DownloadStaticFile(recorderRequest, request, "./test/pic.jpg", "puppy.jpg")

	res := recorderRequest.Result()
	defer res.Body.Close()

	if res.Header["Content-Length"][0] != "98827" {
		t.Error("wrong content length of", res.Header["Content-Length"][0])
	}

	if res.Header["Content-Disposition"][0] != `attachment; filename="puppy.jpg"` {
		t.Error("wrong content disposition")
	}

	_, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
}

var jsonTests = []struct {
	name          string
	json          string
	errorExpected bool
	maxSize       int
	allowUnknown  bool
}{
	{name: "good json", json: `{"foo": "bar"}`, errorExpected: false, maxSize: 1024, allowUnknown: false},
	{name: "badly formatted json", json: `{"foo":}`, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "incorrect type json", json: `{"foo": 1}`, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "two json files", json: `{"foo": "1"}{"alpha":"beta"}`, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "empty body", json: ``, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "syntax error in json", json: `{"foo": 1"`, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "unknown field in json", json: `{"fooo": "1"}`, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "allow unkowns fields in json", json: `{"fooo": "1"}`, errorExpected: false, maxSize: 1024, allowUnknown: true},
	{name: "missing field in json", json: `{jack: "1"}`, errorExpected: true, maxSize: 1024, allowUnknown: true},
	{name: "file too large", json: `{"foo": "bar"}`, errorExpected: true, maxSize: 5, allowUnknown: true},
	{name: "not json", json: `hello world`, errorExpected: true, maxSize: 1024, allowUnknown: true},
}

func TestTools_ReadJSON(t *testing.T) {
	var testTool Tools
	for _, e := range jsonTests {
		testTool.MaxJSONSize = e.maxSize
		testTool.AllowUnknownFields = e.allowUnknown
		var decodedJSON struct {
			Foo string `json:"foo"`
		}

		request, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(e.json)))
		if err != nil {
			t.Log("Error: ", err)
		}

		requestRecorder := httptest.NewRecorder()
		err = testTool.ReadJson(requestRecorder, request, &decodedJSON)

		if e.errorExpected && err == nil {
			t.Errorf("%s error expected, but none received", e.name)
		}

		if !e.errorExpected && err != nil {
			t.Errorf("%s error not expected, but one received: %s", e.name, err.Error())
		}

		request.Body.Close()
	}
}

func TestTools_WriteJson(t *testing.T) {
	var testTools Tools
	responseRecorder := httptest.NewRecorder()

	payload := JSONResponse{
		Error:   false,
		Message: "foo",
	}

	headers := make(http.Header)
	headers.Add("FOO", "BAR")
	err := testTools.WriteJSON(responseRecorder, http.StatusOK, payload, headers)

	if err != nil {
		t.Errorf("Failed to werite json: %v", err)
	}
}

func TestTools_ErrorJson(t *testing.T) {
	var testTools Tools
	responseRecorder := httptest.NewRecorder()

	err := testTools.ErrorJSON(responseRecorder, errors.New("some error"), http.StatusServiceUnavailable)
	if err != nil {
		t.Error(err)
	}

	var payload JSONResponse
	decoder := json.NewDecoder(responseRecorder.Body)
	err = decoder.Decode(&payload)

	if err != nil {
		t.Error("Received error when decoding JSON", err)
	}

	if !payload.Error {
		t.Error("error set to false in JSON, and it should be true")
	}

	if responseRecorder.Code != http.StatusServiceUnavailable {
		t.Errorf("wrong status code returned; expected 503 but got %d", responseRecorder.Code)
	}
}
