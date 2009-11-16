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

  _QUERY_GETSTATUS = "%sstatuses/show/%d.%s";
  _QUERY_UPDATESTATUS = "%sstatuses/update/update.%s";
  _QUERY_PUBLICTIMELINE = "%sstatuses/public_timeline.%s";
  _QUERY_USERTIMELINE = "%sstatuses/user_timeline.%s";
  _QUERY_REPLIES = "%sstatuses/mentions.%s";
  _QUERY_FRIENDSTIMELINE = "%sstatuses/friends_timeline.%s";
  _QUERY_USER_NAME = "%s%s.%s?screen_name=%s";
  _QUERY_USER_ID = "%s%s.%s?user_id=%d";
  _QUERY_USER_DEFAULT = "%s%s.%s";
)

const (
  _STATUS = iota;
  _SLICESTATUS;
  _USER;
  _SLICEUSER;
  _BOOL;
  _ERROR;
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
  cacheBackend *CacheBackend;
  receiveChannel interface{};
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

func (self *Api) SetCacheBackend(backend *CacheBackend) {
  self.cacheBackend = backend;
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
  return self.getUsersByType(user, page, "statuses/followers");
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
  return self.getUsersByType(user, page, "statuses/friends");
}

func (self *Api) getUsersByType(user interface{}, page int, typ string)
                               (<-chan []User) {
  var url string;
  var ok bool;
  responseChannel := self.buildRespChannel(_SLICEUSER).(chan []User);

  if url, ok = self.buildUserUrl(typ, user, page); !ok {
    responseChannel <- nil;
    return responseChannel;
  }

  go self.goGetUsers(url, responseChannel);
  return responseChannel;
}

func (self *Api) GetUser(user interface{}) <-chan User {
  var url string;
  var ok bool;
  var userId int64 = 0;
  responseChannel := self.buildRespChannel(_USER).(chan User);

  // TODO: use username as the cache key instead of id
  switch user.(type) {
  case int:
    userId = int64(user.(int));
    break;
  case int64:
    userId = user.(int64);
  }

  if userId != 0 {
    if cached, hasCached := self.getCachedUser(userId); hasCached {
      responseChannel <- cached;
      return responseChannel;
    }
  }

  if url, ok = self.buildUserUrl("users/show", user, 0); !ok {
    responseChannel <- nil;
    return responseChannel;
  }

  go self.goGetUser(url, responseChannel);
  return responseChannel;
}

// Checks to see if there are any errors in the error channel
func (self *Api) HasErrors() bool {
  return len(self.errors) > 0;
}

// Retrieves the public timeline as a slice of Status objects
func (self *Api) GetPublicTimeline() <-chan []Status {
  responseChannel := self.buildRespChannel(_SLICESTATUS).(chan []Status);

  url := fmt.Sprintf(_QUERY_PUBLICTIMELINE, kTwitterUrl, kFormat);
  go self.goGetStatuses(url, responseChannel);

  return responseChannel;
}

// Retrieves the currently authorized user's 
// timeline as a slice of Status objects
func (self *Api) GetUserTimeline() <-chan []Status {
  responseChannel := self.buildRespChannel(_SLICESTATUS).(chan []Status);

  url := fmt.Sprintf(_QUERY_USERTIMELINE, kTwitterUrl, kFormat);
  go self.goGetStatuses(url, responseChannel);

  return responseChannel;
}

// Returns the 20 most recent statuses posted by the authenticating user and 
// that user's friends. This is the equivalent of /timeline/home on the Web.
// Returns the statuses as a slice of Status objects
func (self *Api) GetFriendsTimeline() <-chan []Status {

  responseChannel := self.buildRespChannel(_SLICESTATUS).(chan []Status);

  url := fmt.Sprintf(_QUERY_FRIENDSTIMELINE, kTwitterUrl, kFormat);
  go self.goGetStatuses(url, responseChannel);

  return responseChannel;
}

// Returns the 20 most recent mentions for the authenticated user
// Returns the statuses as a slice of Status objects
func (self *Api) GetReplies() <-chan []Status {
  responseChannel := self.buildRespChannel(_SLICESTATUS).(chan []Status);

  url := fmt.Sprintf(_QUERY_REPLIES, kTwitterUrl, kFormat);
  go self.goGetStatuses(url, responseChannel);
  return responseChannel;
}

// Builds a response channel for async function calls
func (self *Api) buildRespChannel(channelType int) interface {} {
  const size = 1;

  // TODO: I think it's time to learn the reflect package...
  // this switch statement is to protect the client from
  // using a wrong receive channel
  if self.receiveChannel != nil {
    switch(channelType) {
    case _STATUS:
      if _, ok := self.receiveChannel.(chan Status); ok {
        return self.receiveChannel;
      }
      break;
    case _SLICESTATUS:
      if _, ok := self.receiveChannel.(chan []Status); ok {
        return self.receiveChannel;
      }
      break;
    case _USER:
      if _, ok := self.receiveChannel.(chan User); ok {
        return self.receiveChannel;
      }
      break;
    case _SLICEUSER:
      if _, ok := self.receiveChannel.(chan []User); ok {
        return self.receiveChannel;
      }
      break;
    case _BOOL:
      if _, ok := self.receiveChannel.(chan bool); ok {
        return self.receiveChannel;
      }
      break;
    case _ERROR:
      if _, ok := self.receiveChannel.(chan os.Error); ok {
        return self.receiveChannel;
      }
    }
  }

  switch(channelType) {
  case _STATUS:
    return make(chan Status, size);
  case _SLICESTATUS:
    return make(chan []Status, size);
  case _USER:
    return make(chan User, size);
  case _SLICEUSER:
    return make(chan []User, size);
  case _BOOL:
    return make(chan bool, size);
  case _ERROR:
    return make(chan os.Error, size);
  }

  self.reportError("Invalid channel type");
  return nil;
}

func (self *Api) goGetStatuses(url string, responseChannel chan []Status) {
  responseChannel <- self.getStatuses(url);
}

func (self *Api) goGetUsers(url string, responseChannel chan []User) {
  responseChannel <- self.getUsers(url);
}

func (self *Api) getStatuses(url string) []Status {
  var timelineDummy tTwitterTimelineDummy;
  var timeline []Status;

  jsonString := self.getJsonFromUrl(url);
  json.Unmarshal(jsonString, &timelineDummy);

  dummyLen := len(timelineDummy.Object);
  timeline = make([]Status, dummyLen);

  for i := 0; i < dummyLen; i++ {
    status := &timelineDummy.Object[i];
    timeline[i] = status;
    if err := status.GetError(); err != "" {
      self.reportError(err);
    } else {
      self.cacheBackend.StoreStatus(status);
      self.cacheBackend.StoreUser(status.GetUser());
    }
  }

  return timeline;
}


// TODO: consolidate getStatuses/getUsers when we get generics or when someone
//       submits a patch of reflect wizardry which I can't seem to wrap my head
//       around
func (self *Api) getUsers(url string) []User {
  var usersDummy tTwitterUserListDummy;
  var users []User;

  jsonString := self.getJsonFromUrl(url);
  json.Unmarshal(jsonString, &usersDummy);

  dummyLen := len(usersDummy.Object);
  users = make([]User, dummyLen);

  for i := 0; i < dummyLen; i++ {
    user := &usersDummy.Object[i];
    users[i] = user;
    if err := user.GetError(); err != "" {
      self.reportError(err);
    } else  {
      self.cacheBackend.StoreUser(user);
      self.cacheBackend.StoreStatus(user.GetStatus());
    }
  }

  return users;
}

// Sets the Twitter client header, aka the X-Twitter-Client http header on 
// all POST operations
func (self *Api) SetClientString(client string) {
  self.client = client;
}

// Initializes a new Api object, called by NewApi()
func (self *Api) init() {
  self.errors = make(chan os.Error, 16);
  self.receiveChannel = nil;
  self.client = kDefaultClient;

  // default cache
  userCache := NewMemoryCache();
  statusCache := NewMemoryCache();
  self.cacheBackend = NewCacheBackend(userCache, statusCache, kExpireTime);
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
  return self.errors;
}

// Post a Twitter status message to the authenticated user
//
// The twitter.Api instance must be authenticated
func (self *Api) PostUpdate(status string, inReplyToId int64) <-chan bool  {
  responseChannel := self.buildRespChannel(_BOOL).(chan bool);

  go self.goPostUpdate(status, inReplyToId, responseChannel);
  return responseChannel;
}

func (self *Api) goPostUpdate(status string, inReplyToId int64,
                              response chan bool) {
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
    response <- false;
  }

  response <- true;
}

// Gets a Twitter status given a status id
//
// The twitter.Api instance must be authenticated if the status message
// is private
//
// Returns: a channel which receives a twitter.Status object when
//          the request is completed
func (self *Api) GetStatus(id int64) <-chan Status {
  responseChannel := self.buildRespChannel(_STATUS).(chan Status);

  // grab from cache if we have it
  if cached, hasCached := self.getCachedStatus(id); hasCached {
    responseChannel <- cached;
    return responseChannel;
  }

  go self.goGetStatus(id, responseChannel);
  return responseChannel;
}

func (self *Api) getCachedStatus(id int64) (Status, bool) {
  if self.cacheBackend.HasStatusExpired(id) {
    return nil, false;
  }

  return self.cacheBackend.GetStatus(id), true;
}

func (self *Api) getCachedUser(id int64) (User, bool) {
  if self.cacheBackend.HasUserExpired(id) {
    return nil, false;
  }

  return self.cacheBackend.GetUser(id), true;
}

func (self *Api) SetReceiveChannel(receiveChannel interface{}) {
  self.receiveChannel = receiveChannel;
}

func (self *Api) goGetUser(url string, response chan User) {
  var user tTwitterUserDummy;
  jsonString := self.getJsonFromUrl(url);
  json.Unmarshal(jsonString, &user);

  u := &(user.Object);
  if err := u.GetError(); err != "" {
    self.reportError(err);
  } else {
    self.cacheBackend.StoreUser(u);
    self.cacheBackend.StoreStatus(u.GetStatus());
  }

  response <- u;
}

func (self *Api) goGetStatus(id int64, response chan Status) {
  url := fmt.Sprintf(_QUERY_GETSTATUS, kTwitterUrl, id, kFormat);
  var status tTwitterStatusDummy;
  jsonString := self.getJsonFromUrl(url);
  json.Unmarshal(jsonString, &status);

  s := &(status.Object);
  if err := s.GetError(); err != "" {
    self.reportError(err);
  } else {
    self.cacheBackend.StoreStatus(s);
    self.cacheBackend.StoreUser(s.GetUser());
  }

  response <- s;
}

func (self *Api) reportError(error string) {
  err := &TwitterError{error};
  self.lastError = err;
  ok := self.errors <- err;
  if !ok {
    // The error buffer is full, make room for one
    <-self.errors;
    ok := self.errors <- err;
    if !ok {
      // Yo dawg
      fmt.Fprintf(os.Stderr, "Error adding error to error buffer\n");
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

func (self *Api) buildUserUrl(typ string, user interface{}, page int)
                             (string, bool) {
  var url string;

  if user == nil {
    url = fmt.Sprintf(_QUERY_USER_DEFAULT, kTwitterUrl, typ, kFormat);
    return url, true;
  }

  switch(user.(type)) {
  case string:
    url = fmt.Sprintf(_QUERY_USER_NAME, kTwitterUrl, typ, kFormat, user.(string));
    break;
  case int64:
    url = fmt.Sprintf(_QUERY_USER_ID, kTwitterUrl, typ, kFormat, user.(int64));
    break;
  case int:
    url = fmt.Sprintf(_QUERY_USER_ID, kTwitterUrl, typ, kFormat, user.(int));
    break;
  default:
    self.reportError("User parameter must be a string, int, or int64");
    return "", false;
  }

  return url, true;
}
