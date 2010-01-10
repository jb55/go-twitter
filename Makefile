
include $(GOROOT)/src/Make.$(GOARCH)

TARG=twitter
GOFILES=\
	api.go\
	status.go\
	user.go\
	hacks.go\
	cache.go\
	search.go\
	util.go\
	rate_limit.go\
	http_auth.go

include $(GOROOT)/src/Make.pkg

.PHONY: doc

doc: 
	godoc -html=true twitter \
		| sed -e 's/\/src\/pkg\/twitter\//\.\.\//g' > doc/go-twitter.htm
	godoc -html=false twitter > doc/go-twitter.txt
