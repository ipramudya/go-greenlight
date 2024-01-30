# Project Structure

```
.
├── bin
├── cmd
│ └── api
│ └── main.go
├── internal
├── migrations
├── remote
├── go.mod
└── Makefile
```

## Breakdown

- The `bin` directory will contain our compiled application binaries, ready for deployment to a production server.
- The `cmd/api` directory will contain the application-specific code for our Greenlight API application. This will include the code for running the server, reading and writing HTTP requests, and managing authentication.
- The `internal` directory will contain various ancillary packages used by our API. It will contain the code for interacting with our database, doing data validation, sending emails and so on. Basically, any code which isn’t application-specific and can potentially be reused will live in here. Our Go code under cmd/api will import the packages in the internal directory (but never the other way around).
- The `migrations` directory will contain the `SQL` migration files for our database.
- The `remote` directory will contain the `configuration` files and setup `scripts` for our production server.
- The `Makefile` will contain recipes for automating common administrative tasks — like auditing our Go code, building binaries, and executing database migrations.

Any packages which live under `internal` can only be imported by code inside the parent of the internal directory. In our case, this means that any packages which live in internal can only be imported by code inside our greenlight project