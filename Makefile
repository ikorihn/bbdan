verup:
	echo Makefile cmd/version.go | xargs -n 1 sed -i "s/v1.0.0/$(VER)/g"
	echo Makefile cmd/version.go | xargs -n 1 git add
	git commit -m "update $(VER)"
	git tag $(VER)

