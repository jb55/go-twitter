package twitter

import "io";
import "os";
import "http";

func fixBrokenJson(j string) string {
  return `{"object":` + j + "}";
}

func parseResponse(response *http.Response) (string, os.Error) {
  var b []byte;
  b, _ = io.ReadAll(response.Body);
  response.Body.Close();
  bStr := string(b);

  return bStr, nil;
}
