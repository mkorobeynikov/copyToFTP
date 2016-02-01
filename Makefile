BINNAME=cptoftp
BUILDDIR=release

all:
	go get github.com/mitchellh/gox
	go get
	chmod +x ./build/*
	mkdir ./$(BUILDDIR)/
	gox --build-toolchain 
	gox --output=$(BUILDDIR)"/{{.OS}}/{{.Arch}}/"$(BINNAME)
	BUILDDIR=$(BUILDDIR) ./build/copy_config.sh

clean: 
	rm -rf ./$(BUILDDIR)/
