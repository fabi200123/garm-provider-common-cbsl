// Copyright 2023 Cloudbase Solutions SRL
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package util

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"runtime"
	"testing"

	runnerErrors "github.com/cloudbase/garm-provider-common/errors"
	params "github.com/cloudbase/garm-provider-common/params/v0.1.1"
	"github.com/stretchr/testify/require"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func TestResolveToGithubArch(t *testing.T) {
	ghArch, err := ResolveToGithubArch("amd64")
	require.NoError(t, err)
	require.Equal(t, "x64", ghArch)
}

func TestResolveToGithubArchUnknown(t *testing.T) {
	arch := "some-unknown-arch"

	_, err := ResolveToGithubArch(arch)
	require.Error(t, err)
	require.EqualError(t, runnerErrors.NewNotFoundError("arch %s is unknown", arch), err.Error())
}

func TestResolveToGithubOSType(t *testing.T) {
	ghOSType, err := ResolveToGithubOSType("linux")
	require.NoError(t, err)
	require.Equal(t, "linux", ghOSType)
}

func TestResolveToGithubOSTypeUnknown(t *testing.T) {
	osType := "some-unknown-os"

	_, err := ResolveToGithubOSType(osType)
	require.Error(t, err)
	require.EqualError(t, runnerErrors.NewNotFoundError("os %s is unknown", osType), err.Error())
}

func TestResolveToGithubTag(t *testing.T) {
	ghOSTag, err := ResolveToGithubTag("linux")
	require.NoError(t, err)
	require.Equal(t, "Linux", ghOSTag)
}

func TestResolveToGithubTagUnknown(t *testing.T) {
	osTag := params.OSType("some-unknown-os")

	_, err := ResolveToGithubTag(osTag)
	require.Error(t, err)
	require.EqualError(t, runnerErrors.NewNotFoundError("os %s is unknown", osTag), err.Error())
}

func TestIsValidEmail(t *testing.T) {
	validEmail := "test@example.com"

	isValid := IsValidEmail(validEmail)
	require.True(t, isValid)
	require.Equal(t, true, isValid)
}

func TestIsInvalidEmail(t *testing.T) {
	validEmail := "invalid-email"

	isValid := IsValidEmail(validEmail)
	require.False(t, isValid)
	require.Equal(t, false, isValid)
}

func TestIsAlphanumeric(t *testing.T) {
	validString := "test123"

	isValid := IsAlphanumeric(validString)
	require.True(t, isValid)
	require.Equal(t, true, isValid)
}

func TestIsAlphanumericInvalid(t *testing.T) {
	validString := "test@123"

	isValid := IsAlphanumeric(validString)
	require.False(t, isValid)
	require.Equal(t, false, isValid)
}

func TestGetLoggingWriterEmptyLogFile(t *testing.T) {
	writer, err := GetLoggingWriter("")
	require.NoError(t, err)
	require.Equal(t, os.Stdout, writer)
}

func TestGetLoggingWriterValidLogFile(t *testing.T) {
	// Create a temporary directory for log files.
	dir, err := os.MkdirTemp("", "garm-log")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %s", err)
	}
	// Remove the temporary directory when the test finishes.
	t.Cleanup(func() { os.RemoveAll(dir) })

	// Create a log file path within the temporary directory.
	logFile := path.Join(dir, "test.log")

	writer, err := GetLoggingWriter(logFile)
	require.NoError(t, err)

	// Assert that the writer is of type *lumberjack.Logger and it has the
	// expected settings.
	logger, ok := writer.(*lumberjack.Logger)
	require.True(t, ok)
	require.Equal(t, logFile, logger.Filename)
	require.Equal(t, 500, logger.MaxSize)
	require.Equal(t, 3, logger.MaxBackups)
	require.Equal(t, 28, logger.MaxAge)
	require.Equal(t, true, logger.Compress)
}

func TestGetLoggingWriterFailedToCreateLogFolder(t *testing.T) {
	// Add a log file path that includes a directory that does not exist.
	var logFile string
	if runtime.GOOS == "windows" {
		logFile = "T:/test.log"
	} else {
		logFile = "/non-existent-dir/test.log"
	}

	writer, err := GetLoggingWriter(logFile)
	require.Equal(t, nil, writer)
	require.Error(t, err)
	require.EqualError(t, err, "failed to create log folder")
}

func TestGetLoggingWriterPermisionDenied(t *testing.T) {
	// Create a temporary directory for testing
	dir, err := os.MkdirTemp("", "test-dir")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %s", err)
	}
	// Remove the temporary directory when the test finishes.
	t.Cleanup(func() { os.RemoveAll(dir) })

	// Remove execute permission from the temporary directory
	err = os.Chmod(dir, 0644)
	if err != nil {
		t.Fatalf("failed to remove execute permission from temporary directory: %s", err)
	}

	var dirPath string
	if runtime.GOOS == "windows" {
		dirPath = "T:/non-existent-dir"
	} else {
		dirPath = "/non-existent-dir/"
	}

	_, err = GetLoggingWriter(path.Join(dir, dirPath, "test.log"))
	require.Error(t, err)
	require.EqualError(t, err, "failed to create log folder")
}

func TestConvertFileToBase64(t *testing.T) {
	// Create a temporary file with some test data to be converted to base64.
	err := os.WriteFile("file.txt", []byte("test"), 0o644)
	if err != nil {
		t.Fatalf("failed to create temporary file: %s", err)
	}
	// Remove the temporary file when the test finishes.
	defer os.Remove("file.txt")

	base64Data, err := ConvertFileToBase64("file.txt")
	require.NoError(t, err)
	require.Equal(t, "dGVzdA==", base64Data)
}

func TestConvertFileToBase64FileNotFound(t *testing.T) {
	_, err := ConvertFileToBase64("")
	require.Error(t, err)
	require.ErrorContains(t, err, "reading file")
}

func TestOSToOSType(t *testing.T) {
	osType, err := OSToOSType("windows")
	require.NoError(t, err)
	require.Equal(t, params.Windows, osType)
}

func TestOSToOSTypeUnknown(t *testing.T) {
	os := "some-unknown-os"

	osType, err := OSToOSType(os)
	require.Error(t, err)
	require.Equal(t, params.Unknown, osType)
	require.EqualError(t, err, fmt.Sprintf("no OS to OS type mapping for %s", os))
}

func TestGetTools(t *testing.T) {
	ghArch, err := ResolveToGithubArch("amd64")
	if err != nil {
		t.Fatalf("failed to resolve to github arch: %s", err)
	}

	ghOS, err := ResolveToGithubOSType("linux")
	if err != nil {
		t.Fatalf("failed to resolve to github os type: %s", err)
	}

	tools := []params.RunnerApplicationDownload{
		{
			OS:           &ghOS,
			Architecture: &ghArch,
		},
	}

	ghTools, err := GetTools("linux", "amd64", tools)
	require.NoError(t, err)
	require.Equal(t, "linux", *ghTools.OS)
	require.Equal(t, "x64", *ghTools.Architecture)
}

func TestGetToolsUnsupportedOSType(t *testing.T) {
	osType := params.OSType("some-unknown-os")

	_, err := GetTools(osType, "amd64", nil)
	require.Error(t, err)
	require.EqualError(t, err, fmt.Sprintf("unsupported OS type: %s", osType))
}

func TestGetToolsUnsupportedOSArch(t *testing.T) {
	osArch := params.OSArch("some-unknown-arch")

	_, err := GetTools("linux", osArch, nil)
	require.Error(t, err)
	require.EqualError(t, err, fmt.Sprintf("unsupported OS arch: %s", osArch))
}

func TestGetToolsFailed(t *testing.T) {
	osType := params.OSType("linux")
	osArch := params.OSArch("amd64")

	_, err := GetTools(osType, osArch, nil)
	require.Error(t, err)
	require.EqualError(t, err, fmt.Sprintf("failed to find tools for OS %s and arch %s", osType, osArch))
}

func TestGetRandomString(t *testing.T) {
	randomString, err := GetRandomString(10)
	require.NoError(t, err)
	require.Equal(t, 10, len(randomString))
}

func TestPaswsordToBcrypt(t *testing.T) {
	hash, err := PaswsordToBcrypt("random-password")
	require.NoError(t, err)
	require.Equal(t, 60, len(hash))
}

func TestPaswsordToBcryptFailed(t *testing.T) {
	// We define a long password that exceeds the maximum allowed length for bcrypt
	password := "we-pass-a-password-that-is-more-than-72-bytes-long-which-is-the-maximum-allowed"

	hash, err := PaswsordToBcrypt(password)
	require.Error(t, err)
	require.Equal(t, "", hash)
	require.EqualError(t, err, "failed to hash password")
}

func TestNewLoggingMiddleware(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create a new logging middleware using the test buffer
	loggingMiddleware := NewLoggingMiddleware(&buf)

	// Create a test request and response recorder
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Serve the request through the logging middleware
	loggingMiddleware(testHandler).ServeHTTP(w, req)

	// Assert that the request was served successfully and the log output
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, buf.String(), "GET / HTTP/1.1")
}

func TestSanitizeLogEntry(t *testing.T) {
	sanitizedEntry := SanitizeLogEntry("test\n")
	require.Equal(t, "test", sanitizedEntry)
}

func TestNewID(t *testing.T) {
	// Create a new ID
	id := NewID()

	// Assert that the ID is 12 characters long and is alphanumeric
	require.Equal(t, 12, len(id))
	require.True(t, IsAlphanumeric(id))
}

func TestUTF16FromString(t *testing.T) {
	utf16String, err := UTF16FromString("test")
	require.NoError(t, err)
	require.Equal(t, []uint16{116, 101, 115, 116, 0}, utf16String)
}

func TestUTF16ToString(t *testing.T) {
	string := UTF16ToString([]uint16{116, 101, 115, 116, 0})
	require.Equal(t, "test", string)
}

func TestUint16ToByteArray(t *testing.T) {
	byteArray := Uint16ToByteArray([]uint16{116, 101, 115, 116, 0, 0})
	require.Equal(t, []byte{116, 0, 101, 0, 115, 0, 116, 0, 0, 0}, byteArray)
}

func TestUTF16EncodedByteArrayFromString(t *testing.T) {
	utf16EncodedByteArray, err := UTF16EncodedByteArrayFromString("test")
	require.NoError(t, err)
	require.Equal(t, []byte{116, 0, 101, 0, 115, 0, 116, 0}, utf16EncodedByteArray)
}

func TestCompressData(t *testing.T) {
	compressedData, err := CompressData([]byte("test"))
	require.NoError(t, err)
	require.Equal(t, []byte{31, 139, 8, 0, 0, 0, 0, 0, 0, 255, 42, 73, 45, 46, 1, 0, 0, 0, 255, 255, 1, 0, 0, 255, 255, 12, 126, 127, 216, 4, 0, 0, 0}, compressedData)
}
