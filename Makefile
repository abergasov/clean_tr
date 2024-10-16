build: # build and install app
	go build -o bin/torrent_cleaner ./cmd
	sudo cp bin/torrent_cleaner /usr/local/bin/torrent_cleaner