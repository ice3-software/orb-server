#Orb Server

	Under heavy development

###What will the Orb Server be?

- A scalable TCP/UDP server
- Perfect for realtime backends (e.g. simple multiplayer games)
- Support for grouping, aggregating and querying connected clients
- Uses RPC to achieve distribution
- Authentication and authorization

###First Release

######Functional Requirements

- `Client`s can connect to a `Room`.
- On connection a `Client` receives a connection identifier (UUID).
- `Room`s can contain a configurable number of `Client`s.
- `Client`s can broadcast BSON messages to the other `Client`s in their `Room`.
- `Client`s cannot broadcast messages until their parent `Room` is full.
- The `Server` maintains a stack of `Room`s called the `World`.
- The topmost `Room` in the `World` is the `Waiting Room`, that incoming `Client` connections are assigned to.
- When the `Waiting Room` is full, `Client`s are notified and messages are broadcasted within the `Room`. A new `Waiting Room` is pushed to the `World`'s stack.
- `Room`s have a configurable lifetime.
- `Room`s close when all but 1 of its `Client`s have disconnected or its lifetime has expired.
- When a `Room` is closed, it is sliced out of the `World`'s stack, stops broadcasting message and notifies its remaining clients.

######Non-Functional Requirements

- Distributed: use RPC to maintain a shared world that can run across n nodes.
- Benchmarking: messages need to be broadcasted at blazing speed. `Client`s may broadcast a continuous stream of updates to their `Room`, for example, to update their player's position on-screen ever few milliseconds.
- Load testing: see [this awesome talk](https://www.youtube.com/watch?v=2-pPAvqyluI) for performance monitoring ideas.
- By the end of the first release, I want to have determined acceptable application response times and other performance mertics.
