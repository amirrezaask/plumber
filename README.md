# Plumber
Plumber is a framework for creating data pipelines and stream processing tools.

# Goals
- Simple 
- extensible
- stateful ( both computation state and input stream position )
- cluster support (runners manage execution and data partioning)
- fault tolerance (multiple checkpoint strategies)
- multiple strategies for handling failures ( at most once, at least once, exactly once (actually exactly once affect state))
# Binary protocol
Plumber can use binaries as part of it's processing units. Binraies should get two inputs and return two outputs.
## Example
```shell
$ ./bin1 <state json encoded> input
<updated fields of state json encoded> output

```
# Terminology
## Checkpoint
Checkpoints run under special circumstances and backup current state of system. 
## State
Backends for our stateful processor.
### Redis
### Map
### Bolt
## Stream
Streams are the way we move data around. Streams are the input and output of our application. Streams are stateful and their state is just a part of System state.
### Nats
### Nats-Streaming
### Channel
## Lambda 
Lambdas are pure functions that get the state and an input and return some output.
## System 
System is where our lambdas are glued together and state is being handled as a single application with input and output.
