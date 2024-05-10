build_static:
	go build --ldflags '-linkmode external -extldflags "-static"'

build_static_for_interactive: build_static
	mv interpreter ./interactive