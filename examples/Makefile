include $(GOROOT)/src/Make.$(GOARCH)

all: random-crawler parallel

clean:
	rm -f random-crawler random-crawler.$O

random-crawler: random-crawler.$O
	$(LD) -o $@ $^

random-crawler.$O: random-crawler.go
	$(GC) -o $@ $^

parallel: parallel.$O
	$(LD) -o $@ $^

parallel.$O: parallel.go
	$(GC) -o $@ $^

