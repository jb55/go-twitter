package twitter

import (
  "http";
  "fmt";
  "os";
  "io";
  "json";
)

const kTwitterUrl = "http://www.twitter.com/"
const kFormat = "json"
const kErrFormat = "GoTwitter Error: "

const (
  QUERY_GETSTATUS = "%sstatuses/show/%d.%s";
)

type TwitterError struct {
  error string;
}

type Api struct {
  user string;
  pass string;
  errors chan os.Error;
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

func (self *Api) Init() {
  self.errors = make(chan os.Error, 16);
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

func (self *Api) GetStatusAsync(id int64) (response chan Status, err os.Error) {
  // make it a one-sized buffered channel so our goroutine doesn't sit
  // blocking waiting for the client to receive the data
  c := make(chan Status, 1);
  go self.wrapGetStatus(id, c);
  return c, nil;
}

func (self *Api) GetStatus(id int64) (status Status, e os.Error) {
  status, err := self.wrapGetStatus(id, nil);

  return status, err;
}

func (self *Api) wrapGetStatus(id int64, response chan Status)
                              (status Status, err os.Error) {
  url := fmt.Sprintf(QUERY_GETSTATUS, kTwitterUrl, id, kFormat);

  r, _, err := self.httpGet(url, self.user, self.pass);
  if err != nil {
    err := &TwitterError{kErrFormat + err.String()};
    self.reportError(err);
    return nil, err;
  }

  j, raw, err := parseResponse(r);
  if err != nil {
    err := &TwitterError{kErrFormat + err.String()};
    self.reportError(err);
    return nil, err;
  }

  status = jsonToStatus(raw, j, self.errors);
  if response != nil {
    response <- status;
    return status, nil;
  }

  return status, nil
}

func (self *Api) reportError(error *TwitterError) {
  error.error += "\n";
  ok := self.errors <- error;
  if !ok {
    // The error buffer is full, make room for one
    <-self.errors;
    ok := self.errors <- error;
    if !ok {
      fmt.Fprintf(os.Stderr, "Error buffer error\n");
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

