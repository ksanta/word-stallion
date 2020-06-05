# Functional design

This section describes what happens when a player connections, disconnects and sends messages.

## OnConnect

This happens when a player makes their initial websocket connection. There is nothing to do here.
To reduce the number of DB writes, the server will accept a new player once they send their
initial details (name, horse preference).

## OnNewPlayer

`MessageType: "NewPlayer"`

If a game is in progress, reject the new player request. This behaviour will change in the future.

Create a new Player item in Dynamo. Associate the player with a hardcoded game ID for now.

Send a welcome message to the player. The welcome message contains rules of the upcoming game.

Send a "round summary" message to all players waiting for this game to start. This will show already
waiting players that a new player has joined the game and adds a new track on the screen.

If the number of players waiting to start hits a pre-defined limit, then a game will automatically start.

## DoStartGame

## DoRound

## OnPlayerResponse

## OnDisconnect

## DoWordScrape
