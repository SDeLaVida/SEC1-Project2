```bash

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative MPC/MPC.proto


```

# Run the go program while logging the std.out and std.in
```bash
go run . 0 | tee -filepath "logs/text0.txt"
go run . 1 | tee -filepath "logs/text1.txt"
go run . 2 | tee -filepath "logs/text2.txt"
go run . 3 | tee -filepath "logs/text3.txt"

```