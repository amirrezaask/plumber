# Plumber
Plumber is a framework for creating data pipelines and stream processing tools.

# Features 
- Stateful
- Fault Tolerance
- Extensible

# Goals
- Simple 
- extensible
- stateful ( both computation state and input stream position )
- cluster support (runners manage execution and data partioning)
- fault tolerance (multiple checkpoint strategies)
- multiple strategies for handling failures ( at most once, at least once, exactly once (actually exactly once affect state))

# Terminology
## Stream
Streams are the way we move data around. Streams are the input and output of our application.
## Lambda 
Lambdas are processing and actual program logic
## System 
System is where our lambdas are glued together as a single application with input and output.
## Runner
Runners are used if you want cluster support, Runners are similar to systems but they can partition data.
