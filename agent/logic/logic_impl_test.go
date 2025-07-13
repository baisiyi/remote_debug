package logic

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestImpl_UploadHandler(t *testing.T) {
	impl := New()
	uploadDir := "./uploads"

	// Clean up the uploads directory before and after the test
	os.RemoveAll(uploadDir)
	defer os.RemoveAll(uploadDir)

	// Create a test file
	fileContent := "test file content"
	fileName := "test.txt"
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	fw, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(fw, strings.NewReader(fileContent))
	if err != nil {
		t.Fatal(err)
	}
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", &b)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(impl.UploadHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	var resp map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	expectedPath := filepath.Join(uploadDir, fileName)
	if resp["path"] != expectedPath {
		t.Errorf("handler returned unexpected path: got %v want %v",
			resp["path"], expectedPath)
	}

	// Check if the file was created
	_, err = os.Stat(expectedPath)
	if os.IsNotExist(err) {
		t.Errorf("expected file to be created, but it was not")
	}
}

func TestImpl_UploadHandler_WrongMethod(t *testing.T) {
	impl := New()
	req := httptest.NewRequest("GET", "/upload", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(impl.UploadHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}
