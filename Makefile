export VERSION := $(shell cat VERSION)

inc-patch:
	./inc_version.sh -p $(VERSION) > VERSION
inc-minor:
	./inc_version.sh -m $(VERSION) > VERSION
inc-major:
	./inc_version.sh -M $(VERSION) > VERSION

REPO_HOST="kapha"
REPO_PATH="/opt/web/pub/sour.is/debian/"
ANSIBLE_HOST="phoenix"

release: inc-patch
	git commit -am "release version $(VERSION)"
	git tag -a "v$(VERSION)" -m "tag version $(VERSION)"
	git push --tags
