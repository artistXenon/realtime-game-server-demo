# realtime-game-server-demo

## Goal

- creating connection between [realtime-game-client-demo](https://github.com/artistXenon/realtime-game-client-demo) to communicate through UDP and TCP channel to exchange game data.
- creating connection between [real-time-match-server-demo](https://github.com/artistXenon/realtime-game-match-server-demo) to communicate through http to exchange user authentication, match-making data.

## Build

1. Clone the repository.

> git clone https://github.com/artistXenon/realtime-game-server-demo.git game-server <br>
cd game-server
> 
1. Make sure to configure database address in db/index.go (* preparing a config file is planned for later) 
2. Run 

> go run .
>
