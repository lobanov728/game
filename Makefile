build:
	go build
	go build -o ./server ./server

clean:
	rm ./mub
	rm ./server

serve:
	go run ./server/...

play:
	go run ./main.go ./game/...
	