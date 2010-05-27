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
  "http"
  "fmt"
  "os"
  "json"
  "strconv"
  "time"
  "regexp"
)

const (
  kDefaultClient        = "go-twitter"
  kDefaultClientURL     = "http://jb55.github.com/go-twitter"
  kDefaultClientVersion = "0.1"
  kDefaultUserAgent     = "go-twitter"
  kErr                  = "GoTwitter Error: "
  kWarn                 = "GoTwitter Warning: "
  kDefaultTimelineAlloc = 20

  _QUERY_GETSTATUS       = "http://www.twitter.com/statuses/show/%d.json"
  _QUERY_UPDATESTATUS    = "http://www.twitter.com/statuses/update/update.json"
  _QUERY_PUBLICTIMELINE  = "http://www.twitter.com/statuses/public_timeline.json"
  _QUERY_USERTIMELINE    = "http://www.twitter.com/statuses/user_timeline.json"
  _QUERY_REPLIES         = "http://www.twitter.com/statuses/mentions.json"
  _QUERY_FRIENDSTIMELINE = "http://www.twitter.com/statuses/friends_timeline.json"
  _QUERY_USER_NAME       = "http://www.twitter.com/%s.json?screen_name=%s"
  _QUERY_USER_ID         = "http://www.twitter.com/%s.json?user_id=%d"
  _QUERY_USER_DEFAULT    = "http://www.twitter.com/%s.json"
  _QUERY_SEARCH          = "http://search.twitter.com/search.json"
  _QUERY_RATELIMIT       = "http://twitter.com/account/rate_limit_status.json"
)

const (
  _STATUS = iota
  _SLICESTATUS
  _SLICESEARCH
  _USER
  _SLICEUSER
  _BOOL
  _ERROR
  _RATELIMIT
)

type TwitterError struct {
  error string
}

type Api struct {
  user           string
  pass           string
  errors         chan os.Error
  lastError      os.Error
  client         string
  clientURL      string
  clientVersion  string
  userAgent      string
  receiveChannel interface{}
}

// type that satisfies the os.Error interface
func (self TwitterError) String() string { return self.error }

// Creates and initializes new Api objec
func NewApi() *Api {
  api := new(Api)
  api.init()
  return api
}

func (self *Api) isAuthed() bool {
  // TODO: validate user and pass
  return self.user != "" && self.pass != ""
}

// Returns the last error sent to the error channel.
// Calling this function pops the last error, subsequent calls will be nil
// unless another error has occured.
func (self *Api) GetLastError() os.Error {
  last := self.lastError
  self.lastError = nil
  return last
}

// Gets the followers for a given user represented by a slice
// of twitter.User instances
//
// user:
//  A user id or name to fetch the followers from. If this argument
//  is nil, then the followers are fetched from the authenticated user.
//  This paramater must be an int, int64, or string.
//
// page:
//  Not yet implemented
func (self *Api) GetFollowers(user interface{}, page int) <-chan []User {
  return self.getUsersByType(user, page, "statuses/followers")
}

// Gets the friends for a given user represented by a slice
// of twitter.User instances
//
// user:
//  A user id or name to fetch the friends from. If this argument
//  is nil, then the friends are fetched from the authenticated user.
//  This paramater must be an int, int64, or string.
//
// page:
//  Not yet implemented
func (self *Api) GetFriends(user interface{}, page int) <-chan []User {
  return self.getUsersByType(user, page, "statuses/friends")
}

func (self *Api) getUsersByType(user interface{}, page int, typ string) <-chan []User {
  var url string
  var ok bool
  responseChannel := self.buildRespChannel(_SLICEUSER).(chan []User)

  if url, ok = self.buildUserUrl(typ, user, page); !ok {
    responseChannel <- nil
    return responseChannel
  }

  go self.goGetUsers(url, responseChannel)
  return responseChannel
}

// Performs a simple Twitter search. Returns a slice of twitter.SearchResult
// instances
//
// query:
//  The string of text to search for. This is URL encoded automatically.
func (self *Api) SearchSimple(query string) <-chan []SearchResult {
  return self.Search(query, 0, 0, 0, "", "")
}

// Performs a Twitter search. Returns a slice of twitter.SearchResult instances
// string fields are automatically URL Encoded
// query:
//  The string of text to search for.
// page:
//  The page of results to return. Set to 0 to use the default value.
// perPage:
//  The number of results per page. Set to 0 to use the default value.
// sinceId:
//  Return tweets with status ids greater than the given id. Set to 0
//  to use the default value.
// locale:
//  Specify the language of the query you are sending (only ja is currently
//  effective). This is intended for language-specific clients and the default
//  should work in the majority of cases. Set to an empty string to use
//  the default value.
// lang:
//  Restricts tweets to the given language, given by an ISO 639-1 code.
//  Set to an empty string to use the default value.
func (self *Api) Search(query string, page int, perPage int, sinceId int, locale string, lang string) <-chan []SearchResult {
  variables := make(map[string]string)
  url := _QUERY_SEARCH
  responseChannel := self.buildRespChannel(_SLICESEARCH).(chan []SearchResult)

  variables["q"] = query

  if page >= 2 {
    variables["page"] = strconv.Itoa(page)
  }

  if perPage > 0 {
    variables["rpp"] = strconv.Itoa(perPage)
  }

  if sinceId > 0 {
    variables["since_id"] = strconv.Itoa(sinceId)
  }

  if locale != "" {
    variables["locale"] = locale
  }

  if lang != "" {
    variables["lang"] = lang
  }

  url = addQueryVariables(url, variables)
  go self.goGetSearchResults(url, responseChannel)

  return responseChannel
}

// Returns a channel which receives a twitter.User instance for the given
// username.
//
// id:
//  A twiter user id
func (self *Api) GetUserById(id int64) <-chan User {
  var url string
  var ok bool
  responseChannel := self.buildRespChannel(_USER).(chan User)

  if url, ok = self.buildUserUrl("users/show", id, 0); !ok {
    responseChannel <- nil
    return responseChannel
  }

  go self.goGetUser(url, responseChannel)
  return responseChannel
}

// Returns a channel which receives a twitter.User instance for the given
// username.
//
// name:
//  The screenname of the user
func (self *Api) GetUser(name string) <-chan User {
  var url string
  var ok bool
  responseChannel := self.buildRespChannel(_USER).(chan User)

  if url, ok = self.buildUserUrl("users/show", name, 0); !ok {
    responseChannel <- nil
    return responseChannel
  }

  go self.goGetUser(url, responseChannel)
  return responseChannel
}

// Checks to see if there are any errors in the error channel
func (self *Api) HasErrors() bool { return len(self.errors) > 0 }

// Retrieves the public timeline as a slice of Status objects
func (self *Api) GetPublicTimeline() <-chan []Status {
  responseChannel := self.buildRespChannel(_SLICESTATUS).(chan []Status)
  go self.goGetStatuses(_QUERY_PUBLICTIMELINE, responseChannel)
  return responseChannel
}

// Retrieves the currently authorized user's
// timeline as a slice of Status objects
func (self *Api) GetUserTimeline() <-chan []Status {
  responseChannel := self.buildRespChannel(_SLICESTATUS).(chan []Status)
  go self.goGetStatuses(_QUERY_USERTIMELINE, responseChannel)
  return responseChannel
}

// Returns the 20 most recent statuses posted by the authenticating user and
// that user's friends. This is the equivalent of /timeline/home on the Web.
// Returns the statuses as a slice of Status objects
func (self *Api) GetFriendsTimeline() <-chan []Status {
  responseChannel := self.buildRespChannel(_SLICESTATUS).(chan []Status)
  go self.goGetStatuses(_QUERY_FRIENDSTIMELINE, responseChannel)
  return responseChannel
}

// Returns the 20 most recent mentions for the authenticated user
// Returns the statuses as a slice of Status objects
func (self *Api) GetReplies() <-chan []Status {
  responseChannel := self.buildRespChannel(_SLICESTATUS).(chan []Status)
  go self.goGetStatuses(_QUERY_REPLIES, responseChannel)
  return responseChannel
}

// Returns rate limiting information
func (self *Api) GetRateLimitInfo() <-chan RateLimit {
  responseChannel := self.buildRespChannel(_RATELIMIT).(chan RateLimit)
  go self.goGetRateLimit(_QUERY_RATELIMIT, responseChannel)
  return responseChannel
}

// Set the X-Twitter HTTP headers that will be sent to the server.
//
// client:
//   The client name as a string.  Will be sent to the server as
//   the 'X-Twitter-Client' header.
// url:
//   The URL of the meta.xml as a string.  Will be sent to the server
//   as the 'X-Twitter-Client-URL' header.
// version:
//   The client version as a string.  Will be sent to the server
//   as the 'X-Twitter-Client-Version' header.
func (self *Api) SetXTwitterHeaders(client, url, version string) {
  self.client = client
  self.clientURL = url
  self.clientVersion = version
}

// Builds a response channel for async function calls
func (self *Api) buildRespChannel(channelType int) interface{} {
  const size = 1

  // TODO: I think it's time to learn the reflect package...
  // this switch statement is to protect the client from
  // using a wrong receive channel
  if self.receiveChannel != nil {
    switch channelType {
    case _STATUS:
      if _, ok := self.receiveChannel.(chan Status); ok {
        return self.receiveChannel
      }
      break
    case _SLICESTATUS:
      if _, ok := self.receiveChannel.(chan []Status); ok {
        return self.receiveChannel
      }
      break
    case _SLICESEARCH:
      if _, ok := self.receiveChannel.(chan []SearchResult); ok {
        return self.receiveChannel
      }
      break
    case _USER:
      if _, ok := self.receiveChannel.(chan User); ok {
        return self.receiveChannel
      }
      break
    case _RATELIMIT:
      if _, ok := self.receiveChannel.(chan RateLimit); ok {
        return self.receiveChannel
      }
      break
    case _SLICEUSER:
      if _, ok := self.receiveChannel.(chan []User); ok {
        return self.receiveChannel
      }
      break
    case _BOOL:
      if _, ok := self.receiveChannel.(chan bool); ok {
        return self.receiveChannel
      }
      break
    case _ERROR:
      if _, ok := self.receiveChannel.(chan os.Error); ok {
        return self.receiveChannel
      }
    }
  }

  switch channelType {
  case _STATUS:
    return make(chan Status, size)
  case _SLICESTATUS:
    return make(chan []Status, size)
  case _SLICESEARCH:
    return make(chan []SearchResult, size)
  case _USER:
    return make(chan User, size)
  case _RATELIMIT:
    return make(chan RateLimit, size)
  case _SLICEUSER:
    return make(chan []User, size)
  case _BOOL:
    return make(chan bool, size)
  case _ERROR:
    return make(chan os.Error, size)
  }

  self.reportError("Invalid channel type")
  return nil
}

func (self *Api) goGetStatuses(url string, responseChannel chan []Status) {
  responseChannel <- self.getStatuses(url)
}

func (self *Api) goGetUsers(url string, responseChannel chan []User) {
  responseChannel <- self.getUsers(url)
}

func (self *Api) goGetRateLimit(url string, responseChannel chan RateLimit) {
  var rateLimitDummy tTwitterRateLimitDummy
  jsonString := self.getJsonFromUrl(url)
  json.Unmarshal([]uint8(jsonString), &rateLimitDummy)

  rateLimit := &(rateLimitDummy.Object)

  responseChannel <- rateLimit
}

func (self *Api) goGetSearchResults(url string, responseChannel chan []SearchResult) {
  var searchDummy tTwitterSearchDummy
  var results []SearchResult

  jsonString := self.getJsonFromUrl(url)
  json.Unmarshal([]uint8(jsonString), &searchDummy)

  dummyLen := len(searchDummy.Object.Results)
  results = make([]SearchResult, dummyLen)

  for i := 0; i < dummyLen; i++ {
    result := &searchDummy.Object.Results[i]
    results[i] = result
    if err := result.GetError(); err != "" {
      self.reportError(err)
    }
  }

  responseChannel <- results
}

func (self *Api) getStatuses(url string) []Status {
  var timelineDummy tTwitterTimelineDummy
  var timeline []Status

  jsonString := self.getJsonFromUrl(url)
  json.Unmarshal([]uint8(jsonString), &timelineDummy)

  dummyLen := len(timelineDummy.Object)
  timeline = make([]Status, dummyLen)

  for i := 0; i < dummyLen; i++ {
    status := &timelineDummy.Object[i]
    timeline[i] = status
    if err := status.GetError(); err != "" {
      self.reportError(err)
    } else {
    }
  }

  return timeline
}

func parseTwitterDate(date string) *time.Time {
  r, err := regexp.Compile("\\+0000")

  if err != nil {
    fmt.Fprintf(os.Stderr, err.String() + "\n")
  }

  newStr := r.ReplaceAllString(date, "-0000")
  parsedTime, err := time.Parse(time.RubyDate, newStr)

  if err != nil {
    fmt.Fprintf(os.Stderr, err.String() + "\n")
    return time.LocalTime()
  }

  return parsedTime
}


// TODO: consolidate getStatuses/getUsers when we get generics or when someone
//       submits a patch of reflect wizardry which I can't seem to wrap my head
//       around
func (self *Api) getUsers(url string) []User {
  var usersDummy tTwitterUserListDummy
  var users []User

  jsonString := self.getJsonFromUrl(url)
  json.Unmarshal([]uint8(jsonString), &usersDummy)

  dummyLen := len(usersDummy.Object)
  users = make([]User, dummyLen)

  for i := 0; i < dummyLen; i++ {
    user := &usersDummy.Object[i]
    users[i] = user
    if err := user.GetError(); err != "" {
      self.reportError(err)
    }
  }

  return users
}

// Sets the Twitter client header, aka the X-Twitter-Client http header on
// all POST operations
func (self *Api) SetClientString(client string) {
  self.client = client
}

// Initializes a new Api object, called by NewApi()
func (self *Api) init() {
  self.errors = make(chan os.Error, 16)
  self.receiveChannel = nil
  self.client = kDefaultClient
  self.clientURL = kDefaultClientURL
  self.clientVersion = kDefaultClientVersion
  self.userAgent = kDefaultUserAgent
}

// Overrides the default user agent (go-twitter)
func (self *Api) SetUserAgent(agent string) { self.userAgent = agent }

// Sets the username and password string for all subsequent authorized
// HTTP requests
func (self *Api) SetCredentials(username, password string) {
  self.user = username
  self.pass = password
}

// Disable Twitter authentication, subsequent REST calls will not use
// Authentication
func (self *Api) ClearCredentials() {
  self.user = ""
  self.pass = ""
}

// Returns a channel which receives API errors. Can be used for logging
// errors.
//
//    monitorErrors - listens to api errors and logs them
//
//    func monitorErrors(quit chan bool, errors <-chan os.Error) {
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
func (self *Api) GetErrorChannel() <-chan os.Error {
  return self.errors
}

// Post a Twitter status message to the authenticated user
//
// The twitter.Api instance must be authenticated
func (self *Api) PostUpdate(status string, inReplyToId int64) <-chan bool {
  responseChannel := self.buildRespChannel(_BOOL).(chan bool)

  go self.goPostUpdate(status, inReplyToId, responseChannel)
  return responseChannel
}

func (self *Api) goPostUpdate(status string, inReplyToId int64, response chan bool) {
  url := _QUERY_UPDATESTATUS
  var data string

  data = "status=" + http.URLEscape(status)
  if inReplyToId != 0 {
    reply_data := fmt.Sprintf("&in_reply_to_status_id=%d", inReplyToId)
    data += reply_data
  }

  _, err := httpPost(url, self.user, self.pass, self.client, self.clientURL,
    self.clientVersion, self.userAgent, data)
  if err != nil {
    self.reportError(kErr + err.String())
    response <- false
  }

  response <- true
}

// Gets a Twitter status given a status id
//
// The twitter.Api instance must be authenticated if the status message
// is private
//
// Returns: a channel which receives a twitter.Status object when
//          the request is completed
func (self *Api) GetStatus(id int64) <-chan Status {
  responseChannel := self.buildRespChannel(_STATUS).(chan Status)

  go self.goGetStatus(id, responseChannel)
  return responseChannel
}

func (self *Api) SetReceiveChannel(receiveChannel interface{}) {
  self.receiveChannel = receiveChannel
}

func (self *Api) goGetUser(url string, response chan User) {
  var user tTwitterUserDummy
  jsonString := self.getJsonFromUrl(url)
  json.Unmarshal([]uint8(jsonString), &user)

  u := &(user.Object)
  if err := u.GetError(); err != "" {
    self.reportError(err)
  }

  response <- u
}

func (self *Api) goGetStatus(id int64, response chan Status) {
  url := fmt.Sprintf(_QUERY_GETSTATUS, id)
  var status tTwitterStatusDummy
  jsonString := self.getJsonFromUrl(url)
  json.Unmarshal([]uint8(jsonString), &status)

  s := &(status.Object)
  if err := s.GetError(); err != "" {
    self.reportError(err)
  }

  response <- s
}

func (self *Api) reportError(error string) {
  err := &TwitterError{error}
  self.lastError = err
  ok := self.errors <- err
  if !ok {
    // The error buffer is full, make room for one
    <-self.errors
    ok := self.errors <- err
    if !ok {
      // Yo dawg
      fmt.Fprintf(os.Stderr, "Error adding error to error buffer\n")
    }
  }
}

func (self *Api) getJsonFromUrl(url string) string {
  r, _, error := httpGet(url, self.user, self.pass)
  if error != nil {
    self.reportError(kErr + error.String())
    return ""
  }

  data, err := parseResponse(r)
  data = fixBrokenJson(data)
  if err != nil {
    self.reportError(kErr + err.String())
    return ""
  }

  return data
}

func (self *Api) buildUserUrl(typ string, user interface{}, page int) (string, bool) {
  var url string

  if user == nil {
    url = fmt.Sprintf(_QUERY_USER_DEFAULT, typ)
    return url, true
  }

  switch user.(type) {
  case string:
    url = fmt.Sprintf(_QUERY_USER_NAME, typ, user.(string))
    break
  case int64:
    url = fmt.Sprintf(_QUERY_USER_ID, typ, user.(int64))
    break
  case int:
    url = fmt.Sprintf(_QUERY_USER_ID, typ, user.(int))
    break
  default:
    self.reportError("User parameter must be a string, int, or int64")
    return "", false
  }

  return url, true
}
