# GRPC HTTP Dialer

GRPC HTTP Dialer lets you make a GRPC connection over HTTP.

## Example

* Server side

  ```go
  http.Handle(grpchttpdialer.ProxyPath, grpchttpdialer.Handler())
  if err := http.ListenAndServe(":8080", nil); err != nil {
  	panic(err)
  }
  ```

* Client side

  ```go
  conn, err := grpc.Dial(
  	`localhost:43210`,
  	grpc.WithInsecure(),
  	grpc.WithContextDialer(grpchttpdialer.Dialer(`localhost:8080`)),
  )
  ```

For more details, see the [example](./example) folder.
