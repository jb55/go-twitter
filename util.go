//
// Copyright 2009 Bill Casarin <billcasarin@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package twitter

import (
	"io/ioutil"
	"net/http"
	"fmt"
	"net/url"
)

func fixBrokenJson(j string) string { return `{"object":` + j + "}" }

func parseResponse(response *http.Response) (string, error) {
	var b []byte
	b, _ = ioutil.ReadAll(response.Body)
	response.Body.Close()
	bStr := string(b)

	return bStr, nil
}

func addQueryVariables(url_ string, variables map[string]string) string {
	var addition string
	newUrl := url_

	i := 0
	for key, value := range variables {
		if i == 0 {
			addition = fmt.Sprintf("?%s=%s", key, url.QueryEscape(value))
		} else {
			addition = fmt.Sprintf("&%s=%s", key, url.QueryEscape(value))
		}
		newUrl += addition
		i++
	}

	return newUrl
}
