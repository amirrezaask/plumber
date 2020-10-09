# Plumber
Plumber is a framework for creating data pipelines and stream processing tools.

# Goals
- Simple [x] 
- extensible [x]
- Event buckets [x]
- stateful ( both computation state and input stream position ) [DONE]
- Config-based approach - create stream processor using just JSON config. [DONE]
- cluster support (runners manage execution and data partioning) [TBA]
- fault tolerance (multiple checkpoint strategies) [WIP]
- multiple strategies for handling failures ( at most once, at least once, exactly once (actually exactly once affect state)) [WIP]

# Terminology
## Checkpoint
Checkpoints run under special circumstances and backup current state of system. 
- TimeBased
## State
Backends for our stateful processor.
- Redis
- Map
- Bolt
## Stream
Streams are the way we move data around. Streams are the input and output of our application. Streams are stateful and their state is just a part of System state.
- Nats
- Nats-Streaming
- Channel
- HTTP
- File
- Printer
- Array
## Pipe 
Pipes are pure functions that get the state and an input and return some output. Remember that since Pipes get runned using Goroutiens you can block in them so you can do any 
kind of event buckets in them. ( Similar to Windows in ApacheFlink)
## Pipeline 
Pipeline is where our pipes are glued together and state is being handled as a single application with input and output.

# Usage
## As a library
You can use plumber as a simple library for creating fast scalable data piplines and stream processing tools in Golang.
Example:
```go
package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/amirrezaask/plumber"
	"github.com/amirrezaask/plumber/checkpoint"
	"github.com/amirrezaask/plumber/state"
	"github.com/amirrezaask/plumber/stream"
	"github.com/amirrezaask/plumber/pipeline"
)

func toLower(state plumber.State, value interface{}) (interface{}, error) {
	word := value.(string)
	word = strings.ToLower(word)
	return word, nil
}

func toUpper(state plumber.State, value interface{}) (interface{}, error) {
	word := value.(string)
	word = strings.ToUpper(word)
	return word, nil
}

func count(state plumber.State, input interface{}) (interface{}, error) {
	word := input.(string)
	counter, err := state.GetInt(string(word))
	if err != nil {
		return nil, err
	}
	counter = counter + 1
	err = state.Set(string(word), counter)
	if err != nil {
		return nil, err
	}
	return word, nil
}
func main() {
	// input, err := stream.NewNatsStreaming("localhost:4222", "plumber", "clusterID", "thisclient")
	// input, err := stream.NewNats("localhost:4222", "plumber")
	// if err != nil {
	// 	panic(err)
	// }
	input := stream.NewChanStream()
	// feed some data into our input stream
	go func() {
		for {
			input.Write("This Is tHe eNd")
		}
	}()
	output := stream.NewChanStream()
	//consume our output data
	go func() {
		for v := range output.ReadChan() {
			if v != nil {
				fmt.Println(v)
			}
		}
	}()
	r, err := state.NewRedis(context.Background(), "localhost", "6379", "", "", 0)
	if err != nil {
		panic(err)
	}
	//create our plumber pipeline
	errs, err := pipeline.
		NewDefaultSystem().
		SetCheckpoint(checkpoint.WithInterval(time.Second * 1)).
		SetState(r).
		//SetState(state.NewBolt())
		// SetState(state.NewMapState()).
		From(input).
		Then(toLower).
		Then(toUpper).
		Then(count).
		To(output).
		Initiate()
	if err != nil {
		panic(err)
	}
	for err := range errs {
		fmt.Println(err)
	}
}
```
## As a standalone program
Plumber can also be used as a standalone binary that you feed configuration into it. It has all the benefits of plumber but you can write you processing logic in any language even in Bash.
Example configuration:
```json
{
    "from": {
        "type": "array",
        "args": {
            "words": ["amirreza"]
        }
    },
    "to": {
        "type": "printer",
        "args": {
        }
    },
    "checkpoint": {
        "type": "time-based",
        "args": {
            "interval": 2
        }
    },
    "state": {
        "type": "map",
        "args": {
        }
    },
    "pipeline": [
        {
            "path": "echo",
            "needs_state": false
        },
        {
            "path": "cowsay",
            "needs_state": false
        }
    ]
}
```


## TODO
- Standardize Stream constructors so we don't need to explicily name them in our command line tool.
