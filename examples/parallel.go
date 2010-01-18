package main

import "twitter"
import "flag"
import "fmt"

const kStart = "jb55"

var api *twitter.Api

func main() {
  api = twitter.NewApi()
  isParallel := flag.Bool("p", false, "parallel downloads with go channels")
  flag.Parse()

  if (isParallel != nil && *isParallel) {
    fmt.Printf("Parallel...\n");

    pub_timeline_chan := api.GetPublicTimeline()
    search_chan := api.SearchSimple("@jb55")
    search_chan2 := api.SearchSimple("Google")
    search_chan3 := api.SearchSimple("Hi there")

    <-pub_timeline_chan
    <-search_chan
    <-search_chan2
    <-search_chan3

  } else {
    fmt.Printf("Not Parallel... (-p for parallel)\n");

    <-api.GetPublicTimeline()
    <-api.SearchSimple("@jb55")
    <-api.SearchSimple("Google")
    <-api.SearchSimple("Hi there")
  }
}
