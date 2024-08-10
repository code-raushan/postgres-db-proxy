## Postgres Proxy Server

- Written in Go
- Proxies the client connection to raw db connection and vice versa.
- Operates on raw TCP level without any database level overheads.
- Utilizes goroutines for better concurrency
- Utilizes only standard library utils for minimum dependencies.