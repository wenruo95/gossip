all: 
	go build -v -o gossip && \
		rm -rf bin && \
		mkdir -p bin/conf && \
		mv gossip bin/ && cp -r conf/* bin/conf/
