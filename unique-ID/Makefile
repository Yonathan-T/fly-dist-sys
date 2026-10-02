build-win:
	go build -o unique-ids.exe .

run-win: build-win
	java -jar ..\..\maelstrom\lib\maelstrom.jar test -w unique-ids --bin .\unique-ids.exe --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition

build-linux:
	go build -o unique-ids .

run-linux: build-linux
	java -Djava.awt.headless=true -jar ../../maelstrom/lib/maelstrom.jar test -w unique-ids --bin ./unique-ids --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition
