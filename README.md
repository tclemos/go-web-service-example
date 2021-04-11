# go-web-service-example

> THIS REPOSITORY IS A WORKING IN PROGRESS

This repository aims to provide a fully implemented service to be used as an example helping people creating services professionally.

This project handles a sample scenario where `Things` need to be managed.

---

## API

This project uses `oapi-codegen` to generate the HTTP layer based on our documentation.

**WHY?**

This removes a lot o manual work needed to create the HTTP layer and avoid HTTP documentation to get out of sync with the code.

---

## Persistence

This project uses `sqlc` to generate the persistence layer based on our `SQL code` which represent our databse `queries` and `models`

**WHY?**

`sqlc` integrates very well with `Postgres` and also checks the `SQL commands and entities` during the code generation step, this feature allows us to identify `syntax mistakes` during the CI process instead of while executing the application in PROD, just because we have missed creating a test for that specific scenario.

> TODO: Why Postgres?

Also:

- the standard library is fast but requires a lot of repeated boring code which can lead us to mistakes
- `gorm` helps a lot with CRUD operations, but it requires specific knowledge in the library to advanced usage and it can be 5 times slower than the standard library when dealing with a high volume of data.
- `sqlx` is good, keep things simply and gave control over the SQL code we want to write, but still requires a lot of code to be written repeatedly, this is our recommendation in case you have to work a database engine different than `Postgres`

> TODO: link documentation with gorm analysis

---

## Tests

This project uses `dockertest` from `ORY` to spin up containers during the integration tests, simulating as much as possible a `PRODUCTION` environment as much as possible during the `DEVELOP` phase and also the `CI` process.

To help with `dockertest` usage this project also use `goit` that's a productivity library over `dockertest` to provide pre-defined containers ready to be used and also a container lifecycle.

**WHY?**

First of all `dockertest` allows integrated tests to be written and executed as regular unit tests, this provide us an easy way to maintain and debug integrated tests. 

`goit` take care about running the tests only when the container are ready to be consumed, avoiding a lot of code needed to control the infrastructure like "is my postgres already running?", "is the aws cloud available?", etc.

Another cool stuff from `goit` is that this lib provides a lot of auxiliary tools, for example, you do not need to write your own `SQS Service` in order to test a integrated flow that receives messages from `SQS Queue`, `goit` provides as built-in `SQS Service` already pointing to the `AWS Cloud container` you asked to be created that allows you to send messages to a specific queue in a simple and easy way.

Once `API` and `Persistence` layer is being automatically generated, we can assume these layers will be implemented with no syntax mistakes, but even though the code will compile and have the best practices in its implementation, we still need to guarantee that the domain is being used accordingly to our expectations and that everything is integrated properly.

In other words, if a call is made to `CREATE` something through the `API` that is supposed to call an `internal service` and then `persist data` into the `database`, we MUST test this integration scenario avoid mocking as much as possible, so in this case we spin up a `Postgres` container to provide a real database and also start the application to run the `HTTP server`, after that we can call the `HTTP API` in the test via an `HTTP Client` to check the responses and also execute queries in the database to make sure everything is working as expected.
