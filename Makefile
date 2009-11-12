
include Make.$(GOARCH)

TARG=twitter
PREREQ+=http_auth
GOFILES=\
	api.go\
	status.go\
	http_auth.go

include Make.pkg

