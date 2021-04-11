# ADR: Hexagonal software architecture

## CONTEXT

It's common to see developers having difficult organizing their `Go` components trying to import organization from other languages.

> Why Is Software Architecture Important?
>
> Chapter 1 covered the importance of architecture to an enterprise. In this chapter, we focus on why architecture matters from a technical perspective. In that context, there are fundamentally three reasons for software architecture's importance.
>
> Communication among stakeholders. `Software architecture represents a common abstraction of a system that most if not all of the system's stakeholders can use as a basis for mutual understanding, negotiation, consensus, and communication.`
>
> Early design decisions. Software architecture manifests the earliest design decisions about a system, and these early bindings carry weight far out of proportion to their individual gravity with respect to the system's remaining development, its deployment, and its maintenance life. `It is also the earliest point at which design decisions governing the system to be built can be analyzed.`
>
> Transferable abstraction of a system. Software architecture constitutes a relatively small, intellectually graspable model for how a system is structured and how its elements work together, and this model is transferable across systems. `In particular, it can be applied to other systems exhibiting similar quality attribute and functional requirements and can promote large-scale re-use.`

source: http://www.ece.ubc.ca/~matei/EECE417/BASS/ch02lev1sec4.html#:~:text=Software%20architecture%20represents%20a%20common,negotiation%2C%20consensus%2C%20and%20communication.&text=It%20is%20also%20the%20earliest,be%20built%20can%20be%20analyzed.

## DECISION

We will use the `Hexagonal` software architecture, because `it is easy to learn, concise, expressive and readable` which matches very well with the `Go` programming language purpose.

> It aims at creating loosely coupled application components that can be easily connected to their software environment by means of ports and adapters.

source: https://en.wikipedia.org/wiki/Hexagonal_architecture_(software)

### Directory structure

The project directory structure will be as simple as this:

```text
├── adapters
├── core
│   ├── domain
│   ├── ports
│   └── services
└── main.go
```

Here is an example with more context: 

```text
├── adapters
│   ├── postgres
│   │   └── thing_repository.go
│   ├── kafka
│   │   └── thing_notifier.go
│   └── ...
├── core
│   ├── domain
│   │   ├── thing.go
│   │   ├── another_thing.go
│   │   └── ...
│   ├── ports
│   │   ├── thing_repository.go
│   │   ├── thing_notifier.go
│   │   └── ...
│   └── services
│       ├── thing_service.go
│       └── ...
└── main.go
```

Ok, so let's explain this structure a little better.

Let's start talking about the `core` directory.

The `core` directory must have the software implementation, all the business logic must be contained inside of this directory and they will be separated in 3 parts. Also, no external dependency must be implemented here, we will see more about this below.

The `domain` directory must contain the definition to the core objects, it can be anything your business understand as an object.

The `ports` directory must be used to specify all the contracts between our core and the external world, they are the adapter our core will drive, also known in the `Hexagonal` architecture as the `driven adapters`. Basically they are resources our core requires to work properly, for example a persistence or notification layer. Be aware that this is not the place to implement the communication to a specific database like `postgres` or to a message broker like `kafka`, it must have only the contracts that are agnostic to any specific implementation and depend only to the `domain` objects

The `services` directory must contain all the business logic and will connect the `domain` objects with the `ports` to finally achieve the software goal.

**The `core` implementation must not depend on the `adapters` implementation directory, it must be through `ports` contracts.**

Now, let's talk about the `adapters` directory.

The `adapters` directory must contain all the implementation needed to fulfil the external dependencies needed by the `core`, matching the `ports` contracts.

First of all, the `adapters` implementation must depend on `core`.

So, let's say our software requires a place to store data and we decided to use `postgres` as our database engine. The `ports` directory will have a definition of a `contract` like `thing_repository` specifying the `core` persistence requirements. The `adapters` directory will have a sub directory called `postgres` which will have an implementation of the `thing_requirement` contract to use `postgres` as the database engine and fulfil all the `core` requirements.

The `adapters` directory can also have implementation that are not defined in the `ports` contracts, these are the adapters that will call our `core`, in the `Hexagonal` architecture these adapters are also known as the `driver adapters`.

Let's say you want to expose your `core` via the `HTTP` protocol with `REST APIs`, so you should create a sub directory inside of the `adapters` directory to have all the `HTTP` implementation, defining all the `HTTP specific objects` and parsing them to the respective `domain` objects and then calling the `core services`, once you get the response from the `core services`, you parse it again to a `HTTP response` and the returns it.

The name of the `adapters` sub directories must have a `straight to the point` name, for example, if you want to use `postgres` as a database engine, the directory must have the name `postgres`, if you want to have a `HTTP` exposition of the core using the `Echo` framework, it should be called `http` or `echo`, if you want to send or receive messages from `kafka`, it must be called `kafka`.

Finally we have the `main.go` file at the root of the project. It basically bootstraps the whole application, it has all the `dependency injection` decisions and tie all the `adapters` implementations with the `core`.

## STATUS

Accepted.

## CONSEQUENCES

The `Hexagonal` architecture will help us having a very well organized project.

This will also help new joiners to easily figure out the technologies being used by this project by simply taking a look inside of the `adapters` directory.

This architecture protects our `core` from external interferences, providing us a safe environment to change it as much as needed accordingly to the business instead of requirement to an specific technology or tool.
