PACKAGES=dbm httpsrv ident log uuid
BASE_URL=sour.is/x/toolbox/

test: $(PACKAGES)
	go test $(addprefix $(BASE_URL), $(PACKAGES))

dep: $(PACKAGES)
	go get $(addprefix $(BASE_URL), $(PACKAGES))

