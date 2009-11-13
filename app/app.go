package main

import "twitter"
import "fmt"

func main() {
  api := twitter.NewApi();

  api.Authenticate("jb55", "");
  text := api.GetStatus(5641609144).GetText();
  api.PostUpdate("Testing my Go twitter library", 0);

  fmt.Printf(text);
}
