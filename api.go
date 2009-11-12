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
const kErrFormat = "Error: "

const (
  QUERY_GETSTATUS = "%sstatuses/show/%d.%s";
)

type Api struct {
  pass int;
}

func (api *Api) Init() {
  return
}

func (api *Api) GetStatusAsync(id int64, response chan string) os.Error {
  go wrapGetStatus(id, response);
  return nil;
}

func (api *Api) GetStatus(id int64) (status string, e os.Error) {
  c := make(chan string);
  err := wrapGetStatus(id, c);

  return <-c, err;
}

func wrapGetStatus(id int64, response chan string) os.Error {
  url := fmt.Sprintf(QUERY_GETSTATUS, kTwitterUrl, id, kFormat);

  r, _, err := http.Get(url);
  if err != nil {
    response <- kErrFormat + err.String();
  }

  j, err := parseResponse(r);
  if err != nil {
    response <- kErrFormat + err.String();
  }

  response <- j.Get("text").String();
  return nil;
}

func parseResponse(response *http.Response) (*json.Json, os.Error) {
  var b []byte;
  b, _ = io.ReadAll(response.Body);
  response.Body.Close();

  j, ok, _ := json.StringToJson(string(b));
  if !ok {
    return nil, nil
  }

  return &j, nil;
}


