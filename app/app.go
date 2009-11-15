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
package main

import "twitter"
import "fmt"

func main() {
  const nIds = 10;
  var startId int64 = 5641609144;

  api := twitter.NewApi();
  errors := api.GetErrorChannel();
  receiveChannel := make(chan twitter.Status, nIds);
  api.SetReceiveChannel(receiveChannel);

  for i := 0; i < nIds; i++ {
    api.GetStatus(startId);
    startId++;
  }

  for i := 0; i < nIds; i++ {
    // reads in status messages as they come in
    status := <-receiveChannel;
    fmt.Printf("Status #%d %s: %s\n", i,
                status.GetUser().GetScreenName(),
                status.GetText());
  }

  showFollowers(api, "jb55");

  for i := 0; api.HasErrors(); i++ {
    fmt.Printf("Error #%d: %s\n", i, <-errors);
  }

  //api.PostUpdate("Testing my Go twitter library", 0);
}

func showFollowers(api *twitter.Api, user interface{}) {
  followers := <-api.GetFollowers(user, 0);

  for _, follower := range followers {
    fmt.Printf("%v\n", follower.GetName());
  }
}
