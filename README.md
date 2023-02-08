# calculator
Microservice based equation calculator

### Info
- `server` is a service that exposes a GRPC API to solve an equation.
- `cli` is a tool to call this API to solve some equations.

The protos are placed outside the project to simulate a realistic setup.

### Example
Make a query against the running server using grpcurl
```
grpcurl -plaintext -d '{"equation":"1+1"}' localhost:8000 calculator.CalculationService/Run
```

List all grpc end points using grpcurl
```
grpcurl --plaintext localhost:8000 list
```

### Dev iterations
**OK** grpc/protobuf service for calculation + client calling service
everything in docker + message broker to communicate indirectly with calculation service
