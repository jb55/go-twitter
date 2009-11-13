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
  "json";
)

const (
  kTwitterUrl = "http://www.twitter.com/";
  kDefaultClient = "go-twitter";
  kFormat = "json";
  kErr = "GoTwitter Error: ";
  kWarn = "GoTwitter Warning: ";
  kDefaultTimelineAlloc = 20;
)

const (
  _QUERY_GETSTATUS = "%sstatuses/show/%d.%s";
  _QUERY_UPDATESTATUS = "%sstatuses/update/update.%s";
  _QUERY_PUBLICTIMELINE = "%sstatuses/public_timeline.%s";
  _QUERY_USERTIMELINE = "%sstatuses/user_timeline.%s";
  _QUERY_REPLIES = "%sstatuses/mentions.%s";
  _QUERY_FRIENDSTIMELINE = "%sstatuses/friends_timeline.%s";
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

// Checks to see if there are any errors in the error channel
func (self *Api) HasErrors() bool {
  return len(self.errors) > 0;
}

// Retrieves the public timeline as a slice of Status objects
func (self *Api) GetPublicTimeline() []Status {
  url := fmt.Sprintf(_QUERY_PUBLICTIMELINE, kTwitterUrl, kFormat);
  statuses := self.getStatuses(url);

  return statuses;
}

// Retrieves the currently authorized user's 
// timeline as a slice of Status objects
func (self *Api) GetUserTimeline() []Status {
  url := fmt.Sprintf(_QUERY_USERTIMELINE, kTwitterUrl, kFormat);
  statuses := self.getStatuses(url);

  return statuses;
}

// Returns the 20 most recent statuses posted by the authenticating user and 
// that user's friends. This is the equivalent of /timeline/home on the Web.
// Returns the statuses as a slice of Status objects
func (self *Api) GetFriendsTimeline() []Status {
  url := fmt.Sprintf(_QUERY_FRIENDSTIMELINE, kTwitterUrl, kFormat);
  statuses := self.getStatuses(url);

  return statuses;
}

// Returns the 20 most recent mentions for the authenticated user
// Returns the statuses as a slice of Status objects
func (self *Api) GetReplies() []Status {
  url := fmt.Sprintf(_QUERY_REPLIES, kTwitterUrl, kFormat);
  statuses := self.getStatuses(url);
  return statuses;
}

func (self *Api) getStatuses(url string) []Status {
  var timelineDummy tTwitterTimelineDummy;
  var timeline []Status;

  jsonString := self.getJsonFromUrl(url);
  json.Unmarshal(jsonString, &timelineDummy);

  dummyLen := len(timelineDummy.Object);
  timeline = make([]Status, dummyLen);

  for i := 0; i < dummyLen; i++ {
    timeline[i] = &timelineDummy.Object[i];
  }

  return timeline;
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
  url := fmt.Sprintf(_QUERY_UPDATESTATUS, kTwitterUrl, kFormat);
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
  url := fmt.Sprintf(_QUERY_GETSTATUS, kTwitterUrl, id, kFormat);
  var status tTwitterStatusDummy;
  jsonString := self.getJsonFromUrl(url);
  json.Unmarshal(jsonString, &status);

  s := &(status.Object);

  if response != nil {
    response <- s;
    return s;
  }

  return s;
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

func (self *Api) getJsonFromUrl(url string) string {
  r, _, error := httpGet(url, self.user, self.pass);
  if error != nil {
    self.reportError(kErr + error.String());
    return "";
  }

  data, err := parseResponse(r);
  data = fixBrokenJson(data);
  if err != nil {
    self.reportError(kErr + err.String());
    return "";
  }

  return data;
}
