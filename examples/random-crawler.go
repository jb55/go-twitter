package main

import "twitter"
import "fmt"
import "time"
import "rand"

const kMaxDepth = 5
const kStart = "jb55"

var api *twitter.Api;
var r *rand.Rand;

func main() {
  api = twitter.NewApi();
  r = rand.New(rand.NewSource(time.Seconds()));
  crawl(kStart, 0);
}

func crawl(userName string, level int) {
  // Get the user's status
  text := (<-api.GetUser(userName)).GetStatus().GetText();

  for i := 0; i < level; i++ {
    fmt.Printf("  ");
  }

  fmt.Printf("%s: %s\n", userName, text);
  // Get the user's friends
  friends := <-api.GetFriends(userName, 1);
  length := len(friends);

  if length == 0 {
    return;
  }

  rVal := r.Intn(length-1);
  // Choose a random friend for the next user
  nextUser := friends[rVal].GetScreenName();

  level++;
  if level > kMaxDepth {
    return;
  }
  crawl(nextUser, level);
}
