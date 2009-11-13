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
  "http";
  "fmt";
  "os";
  "io";
  "json";
)

const (
  kTwitterUrl = "http://www.twitter.com/";
  kDefaultClient = "go-twitter";
  kFormat = "json";
  kErr = "GoTwitter Error: ";
  kWarn = "GoTwitter Warning: ";l
)

const (
  QUERY_GETSTATUS = "%sstatuses/show/%d.%s";
  QUERY_UPDATESTATUS = "%sstatuses/update/update.%s";
)

type TwitterError struct {
  error string;
}

type Api struct {
  user string;
  pass string;
  errors chan os.Error;
  lastError os.Error;
  client string;
}

// type that satisfies the os.Error interface
func (self TwitterError) String() string {
  return self.error;
}

// Creates and initializes new Api objec
func NewApi() *Api {
  api := new(Api);
  api.init();
  return api;
}

func (self *Api) isAuthed() bool {
  // TODO: validate user and pass
  return self.user != "" && self.pass != "";
}

// Returns the last error sent to the error channel.
// Calling this function pops the last error, subsequent calls will be nil
// unless another error has occured.
func (self *Api) GetLastError() os.Error {
  last := self.lastError;
  self.lastError = nil;
  return last;
}

// Sets the Twitter client header, aka the X-Twitter-Client http header on 
// all POST operations
func (self *Api) SetClientString(client string) {
  self.client = client;
}

// Initializes a new Api object, called by NewApi()
func (self *Api) init() {
  self.errors = make(chan os.Error, 16);
  self.client = kDefaultClient;
}

// Sets the username and password string for all subsequent authorized
// HTTP requests
func (self *Api) Authenticate(username, password string) {
  self.user = username;
  self.pass = password;
}

// Disable Twitter authentication, subsequent REST calls will not use
// Authentication
func (self *Api) Logout() {
  self.user = "";
  self.pass = "";
}

// Returns a channel which receives API errors. Can be used for logging
// errors or pseudo-exception handling. Eg.
//
//    monitorErrors - listens to api errors and logs them    
//
//    func monitorErrors(quit chan bool, errors chan os.Error) {
//      for ;; {
//        select {
//        case err := <-errors:
//          fmt.Fprintf(os.Stderr, err.String());
//          break;
//        case <-quit:
//          return;
//        }
//      }
//    }
//
func (self *Api) GetErrorChannel() chan os.Error {
  return self.errors;
}

// Post a Twitter status message to the authenticated user
//
// The twitter.Api instance must be authenticated
func (self *Api) PostUpdate(status string, inReplyToId int64) {
  url := fmt.Sprintf(QUERY_UPDATESTATUS, kTwitterUrl, kFormat);
  var data string;

  data = "status=" + http.URLEscape(status);
  if inReplyToId != 0 {
    reply_data := fmt.Sprintf("&in_reply_to_status_id=%d", inReplyToId);
    data += reply_data;
  }

  _, err := httpPost(url, self.user, self.pass, self.client, data);
  if err != nil {
    self.reportError(kErr + err.String());
    return;
  }

  return;
}

// Gets a Twitter status given a status id
//
// The call is made asyncronously and returns instantly
// returns a channel that receives the Status interface when the request
// is completed
//
// The twitter.Api instance must be authenticated if the status message
// is private
func (self *Api) GetStatusAsync(id int64) chan Status {
  // make it a one-sized buffered channel so our goroutine doesn't sit
  // blocking waiting for the client to receive the data
  response := make(chan Status, 1);
  go self.wrapGetStatus(id, response);
  return response;
}

// Gets a Twitter status given a status id
func (self *Api) GetStatus(id int64) Status {
  status := self.wrapGetStatus(id, nil);

  return status;
}

func (self *Api) wrapGetStatus(id int64, response chan Status) Status {
  url := fmt.Sprintf(QUERY_GETSTATUS, kTwitterUrl, id, kFormat);

  r, _, err := httpGet(url, self.user, self.pass);
  if err != nil {
    self.reportError(kErr + err.String());
    return nil;
  }

  j, raw, err := parseResponse(r);
  if err != nil {
    self.reportError(kErr + err.String());
    return nil;
  }

  status := jsonToStatus(raw, j, self.errors);
  if response != nil {
    response <- status;
    return status;
  }

  return status;
}

func (self *Api) reportError(error string) {
  error += "\n";
  err := &TwitterError{error};
  self.lastError = err;
  ok := self.errors <- err;
  if !ok {
    // The error buffer is full, make room for one
    <-self.errors;
    ok := self.errors <- err;
    if !ok {
      fmt.Fprintf(os.Stderr, kErr + "Error adding error to error buffer\n");
    }
  }
}

func parseResponse(response *http.Response) (*json.Json, string, os.Error) {
  var b []byte;
  b, _ = io.ReadAll(response.Body);
  response.Body.Close();
  bStr := string(b);

  j, ok, err := json.StringToJson(bStr);
  if !ok {
    return nil, bStr, &TwitterError{err};
  }

  return &j, bStr, nil;
}

