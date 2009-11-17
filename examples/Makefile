include $(GOROOT)/src/Make.$(GOARCH)

all: random-crawler

clean:
	rm -f random-crawler random-crawler.$O

random-crawler: random-crawler.$O
	$(LD) -o $@ $^

random-crawler.$O: random-crawler.go
	$(GC) -o $@ $^
