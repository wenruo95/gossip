all: 
	go build -v -o gossip && \
		rm -rf release && \
		mkdir -p release/conf && \
		mv gossip release/ && cp -r conf/* release/conf/
