# LRU Cache Visualizer

This Project is for the visualization of LRU (Lease recently used) Cache. The server is written in GOLANG and client is in React.Js

## Installation

### Client

make sure `node` is installed, v18.16.0 or higher, after that in client directory run:

```bash
npm install
```

to run client locally:

```bash
npm run start
```

### server

install `golang` after that in server directory, run:

```bash
go run server.go
```

## Usage

in the client you can add, get or delete keys, also key will be deleted after duration sopecified experies.
client makes `socket` connection with the server to get updates on cache.