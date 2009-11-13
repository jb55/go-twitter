include ../Make.$(GOARCH)

all: app

clean:
	rm -f app app.$O

app: app.$O
	$(LD) -o $@ $^

app.$O: app.go
	$(GC) -o $@ $^
