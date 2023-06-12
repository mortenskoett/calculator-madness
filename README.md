# Calculator Madness
Microservice playground with a calculator theme.
**This project does nothing and is not used for anyting other than learning.**

https://github.com/mortenskoett/calculator-madness/assets/17837870/8c066d91-523f-4b28-8727-fbf2b9fda362

## So what is it?
This is a small application that takes text input (an equation) and spends an inordinate amount of time counting how many chars the equation consists of. The complexity lies in the backend which has a microservice-based architecture that facilitates async message handling between `calculator` and `viewer`.
Each connected client (browser) is connected using websocket making it possible to receive continous progress messages in the UI.

The system consists of the following components:
1. `viewer` is the front facing UI server which establishes websocket connections to the clients, does GRPC calls to the `calculator` and receives results and progress messages over `nsq`. Because a single `viewer` instance might serve multiple websocket connections, it is necessary to pass along a connection id through the backend calls.
2. `nsq` is an instance of the NSQ(New Simple Queue) message queue implemented as a shared Go module. It contains a producer and a consumer that is reused across services. The consumer is generic on the message interface type. Due to the uniqueness of each websocket connection, a unique topic and consumer is used per `viewer` instance.
3. `calculator` is the core business logic which amounts to `return len(equation)`. However, the service features a GRPC interface and will publish equation processing results to `nsq`.

## Usage
### Run all services in docker
```
make build run
```

### Make a query against the calculator
```
grpcurl -plaintext -d '{"client_id":"<cid>", "equation":{"id":"<eid>", "value":"1+1"}, "result_topic":"<topic>"}' localhost:8000 calculator.CalculationService/Run
```

### List all grpc endpoints of the calculator
```
grpcurl --plaintext localhost:8000 list
```

### Debugging lookupd / running locally
- If using nsqlookupd it can be necessary to set broadcast address of nsqd to be able to run nsq in docker next to other services running on the host e.g. for debugging. See: https://github.com/nsqio/go-nsq/issues/69

```
command: /nsqd --lookupd-tcp-address=nsqlookupd:4160 #--broadcast-address=127.0.0.1 #example for debugging locally
```

# Details
## Overall system architecture
Component diagram showing the essential modules and their general interaction.
```mermaid
flowchart LR
    subgraph viewer
        http-->static
        http-->manager
        http-->templates
        subgraph websocket
            manager-->clients
            clients-->router
        end
    end

    subgraph calculator
        grpc--create calculation-->calc
    end

    subgraph nsq
        producer-.enqueue.->queue
        consumer-.dequeue.->queue
    end

    browser-->http
    calc--async-->producer
    manager--async-->consumer
    router--sync-->grpc
```

## Interaction between calculator and nsq
Sequence diagram showing interaction between `calculator` and `nsq` when calculations are created. The `calculator` service takes requests on a GRPC endpoint but returns results asynchronously over nsq.
```mermaid
sequenceDiagram
    participant cli as Client
    participant cal as Calculator
    participant nsq as NSQ

    cli->>cal: RunCalculation(id)
    cal-->>cli: CalcRecieved(id)
    Note over cal: Calculator maintains state of calculations

    %%cal--)nsq: SendCalcStarted(id)
    loop calc in progress
        cal-)nsq: SendCalcProgress(id)
    end
    cal-)nsq: SendCalcEnded(id)

    Note over nsq, cli: Client listens and receives messages on topics

    nsq--)cli: CalcProgress()
    nsq--)cli: CalcEnded()
```

## Creating a new calculation in the viewer
Sequence diagram of how `viewer` interacts with the backend when a new calculation is created.
```mermaid
sequenceDiagram
    participant B as Browser
    participant W as Viewer
    participant C as Calculator
    participant Q as NSQ

    B->>W: New equation

    Note over W: Keep state of calculations

    W-->>B: Create new calc <ID>

    W->>+C: Start calculation

    loop while in progress
        C-)Q: Enqueue calc progress <ID>
        Q--)W: Return calc progress <ID>
        W-->>B: Send calc progress <ID> event
    end

    C->>-Q: Enqueue calc done <ID>
    Q--)W: Return calc done <ID>
    W-->>B: Send calc done <ID> event
```

## Handling calculation progress and results concurrently
Sequence diagram showing the interaction of the concurrent handling of equations in the CalculatorService.
```mermaid
sequenceDiagram
    participant CL as Viewer
    participant GR as GRPCServer
    participant CS as CalculatorService
    participant EP as EquationProcessor
    participant RN as ResultNotifier

    CL->>GR: Request to solve equation
    GR->>CS: Call Solve(eq)
    CS->>EP: Send equation to intake chan
    CS-->>CL: Return OK grpc response

    Note over CL, RN: Equation results are received on queue

    loop for N=len(equation) seconds
        EP-->>CS: New progress msg
        CS->>RN: Post progress msg to notifier
        RN-->>CL: Return progress response to client
    end


    EP-->>CS: Ended equation msg
    CS->>RN: Post ended equation msg to notifier
    RN-->>CL: Return ended equation response to client

```
