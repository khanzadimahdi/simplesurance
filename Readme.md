# How to Run:


### Docker

if you have `docker compose` setuped, run the below command in the project's directory:


```bash
docker compose up
```

then open `localhost:8080` to see the result.

### Go

you need Golang 1.22+ to be able to run the application.

the application can be ran by Golang like below:

```bash
go run main.go
```

then open `localhost:8080` to see the result.

### store / load previous state

on termination, the state of application will be stored in the `./storage/counter.json` and will be loaded on next runs.