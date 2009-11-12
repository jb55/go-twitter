
include Make.$(GOARCH)

TARG=twitter
GOFILES=\
	api.go\
	status.go\
	http_auth.go

include Make.pkg

