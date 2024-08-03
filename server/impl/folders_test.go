/* Copyright 2023 Take Control - Software & Infrastructure

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package impl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-errors/errors"
	"github.com/takecontrolsoft/go_multi_log/logger"
	"github.com/takecontrolsoft/sync_server/server/config"
	"github.com/takecontrolsoft/sync_server/server/utils"
)

func init() {
	config.InitFromEnvVariables()
}

func TestGetFolders(t *testing.T) {
	userName := "Desi"
	deviceId := utils.GenerateRandomString(5)
	body, err := postForm(userName, deviceId)
	if err != nil {
		t.Fatal(err)
	}
	print(body)
}
func jsonReaderFactory(in interface{}) (io.Reader, error) {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	err := enc.Encode(in)
	if err != nil {
		return nil, fmt.Errorf("creating reader: error encoding data: %s", err)
	}
	return buf, nil
}

type data struct {
	User     string
	DeviceId string
}

func postForm(userName string, deviceId string) (string, error) {
	body := data{User: userName, DeviceId: deviceId}
	r, err := jsonReaderFactory(body)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	req, err := http.NewRequest("POST", "/folders", r)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetFoldersHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code == http.StatusOK {
		result := []string{}
		if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
			logger.Error(err)
			return "", err
		}
		return result[0], nil
	} else {
		return "", &RequestError{
			StatusCode: rr.Code,
			Err:        errors.New(rr.Body),
		}
	}
}
