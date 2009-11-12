
include $(GOROOT)/src/Make.$(GOARCH)

TARG=twitter
GOFILES=\
	api.go\
	status.go\

tests: all twitter_tests

twitter_tests: tests.$O
	$(LD) -o $@ $^

tests.$O: gotwitter_tests.go
	$(GC) -o $@ $^
 
include $(GOROOT)/src/Make.pkg
