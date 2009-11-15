============
go-twitter
============

go-twitter is a Twitter library package for Go. The interface is similar to 
python-twitter.


Installation
============

#. Make sure you have Go installed and have your environment is set up
   correctly: $GOROOT, $GOARCH, $GOBIN, etc.

#. Checkout the code from the repository or extract the source code.

#. cd go-twitter && make && make install


Quick Start
===========

::

  import "twitter"

  // Prints the public timeline
  func main() {
    api := twitter.NewApi();
    pubTimeline := <-api.GetPublicTimeline();

    for i, status := range pubTimeline {
      fmt.Printf("#%d %s: %s", i,
                status.GetUser().GetScreenName(),
                status.GetText());
    }
  }


Documentation
=============

doc/ - godoc generated files, site coming soon

