build-win:
	go build -o broadcast.exe .

run-win: build-win
	java -jar ..\..\..\maelstrom\lib\maelstrom.jar test -w broadcast --bin .\broadcast.exe --node-count 1 --time-limit 20 --rate 10

build-linux:
	go build -o broadcast .

run-linux: build-linux
	java -Djava.awt.headless=true -jar ../../../maelstrom/lib/maelstrom.jar test -w broadcast --bin ./broadcast --node-count 1 --time-limit 20 --rate 10
