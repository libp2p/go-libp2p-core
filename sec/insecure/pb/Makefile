PB = $(wildcard *.proto)
GO = $(PB:.proto=.pb.go)

OUTDIR = ./out

all: $(GO)

$(OUTDIR):
	mkdir -p $(OUTDIR)

%.pb.go: %.proto $(OUTDIR)
		protoc --proto_path=$(GOPATH)/src:../../..:. --gogofaster_out=$(OUTDIR) $< && \
		find $(OUTDIR) -name $@ -exec mv '{}' . \; && \
		rm -rf $(OUTDIR)

clean:
		rm -f *.pb.go
		rm -f *.go
