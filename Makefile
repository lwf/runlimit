target/runlimit: target
	go build -tags netgo -installsuffix netgo -o $@ .

clean:
	rm -f target/*

target:
	mkdir -p target
