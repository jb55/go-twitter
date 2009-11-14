
include Make.$(GOARCH)

TARG=twitter
GOFILES=\
	api.go\
	status.go\
	user.go\
	util.go\
	http_auth.go

include Make.pkg

.PHONY: doc

doc: 
	godoc -html=true twitter \
		| sed -e 's/\/src\/pkg\/twitter\//\.\.\//g' > doc/go-twitter.htm
	godoc -html=false twitter > doc/go-twitter.txt
