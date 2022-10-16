package toolkit

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
)

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
