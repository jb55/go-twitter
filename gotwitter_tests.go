package main

import "twitter"
import "fmt"

func main() {
  api := &twitter.Api{0};
  //c := make(chan string);

  status, _ := api.GetStatus(5611676860);

  fmt.Printf("%s\n", status);
}
