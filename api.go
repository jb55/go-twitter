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

func NewApi() *Api {
  api := new(Api);
  api.Init();
  return api;
}

func (self *Api) isAuthed() bool {
  // TODO: validate user and pass
  return self.user != "" && self.pass != "";
}

func (self *Api) GetLastError() os.Error {
  last := self.lastError;
  self.lastError = nil;
  return last;
}

func (self *Api) SetClientString(client string) {
  self.client = client;
}

func (self *Api) Init() {
  self.errors = make(chan os.Error, 16);
  self.client = kDefaultClient;
}

func (self *Api) Authenticate(username, password string) {
  self.user = username;
  self.pass = password;
}

func (self *Api) Logout() {
  self.user = "";
  self.pass = "";
}

func (self *Api) GetErrorChannel() chan os.Error {
  return self.errors;
}

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

func (self *Api) GetStatusAsync(id int64) chan Status {
  // make it a one-sized buffered channel so our goroutine doesn't sit
  // blocking waiting for the client to receive the data
  response := make(chan Status, 1);
  go self.wrapGetStatus(id, response);
  return response;
}

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

