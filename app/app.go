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

  api := twitter.NewApi();
  errors := api.GetErrorChannel();

  //showSearch(api, "@jb55")
  rateLimitInfo := <-api.GetRateLimitInfo()

  fmt.Printf("Remaining hits this hour: %d/%d\n",
    rateLimitInfo.GetRemainingHits(),
    rateLimitInfo.GetHourlyLimit())

  fmt.Printf("other: %d, %s\n",
    rateLimitInfo.GetResetTimeInSeconds(),
    rateLimitInfo.GetResetTime())

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

func showFriends(api *twitter.Api, user interface{}) {
  friends := <-api.GetFriends(user, 0);

  for _, friend := range friends {
    fmt.Printf("%v\n", friend.GetName());
  }
}

func showSearch(api *twitter.Api, query string) {
  results := <-api.SearchSimple(query);

  for _, result := range results {
    fmt.Printf("%s: %s\n", result.GetFromUser(), result.GetText());
  }
}
