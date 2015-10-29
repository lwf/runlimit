clean:
	rm -f target/*

target:
	mkdir -p target

target/runlimit: target
	go build -tags netgo -installsuffix netgo -o $@ .
