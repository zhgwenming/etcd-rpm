NAME ?= etcd
VERSION ?= 0.4.5
NV = $(NAME)-v$(VERSION)
SPEC_FILE := $(NAME).spec
TARBALL := $(NV).tar.gz
SRCDIR ?= ./
DESTDIR ?= ./

.PHONY: srpm

srpm: $(TARBALL) $(SPEC_FILE)
	rpmbuild --define '_sourcedir $(SRCDIR)' --define '_srcrpmdir $(DESTDIR)' -bs $(SPEC_FILE)

$(TARBALL): $(NAME)
	tar --xform='s/etcd/etcd-$(VERSION)/' -zcf $@ $^

clean:
	rm -fv *gz *rpm

build: srpm
	rpmbuild --rebuild $(NAME)*$(VERSION)*.src.rpm
