.PHONY: test-compose

test-compose:
	for t in scripts/* ; do bash $$t ; done
