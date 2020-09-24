# Plumber
Plumber is a framework for creating data pipelines and stream processing tools.

# Goals
- Simple 
- extensible
- stateful ( both computation state and input stream position )
- cluster support (runners manage execution and data partioning)
- fault tolerance (multiple checkpoint strategies)
- multiple strategies for handling failures ( at most once, at least once, exactly once (actually exactly once affect state))

# Terminology
## Stream
Streams are the way we move data around.
## Lambda 
Lambdas are processing and actual program logic
## System 
System is the collection of lambdas plus Input and output streams.
System does all plumbing between lambdas.
