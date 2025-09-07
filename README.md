# Tx Parser
## Initial data
The service comes pre-configured with some seeded addresses and an initial block to start parsing from, both of the addresses should have some transactions after the blocks start getting parsed on startup.  
```go
func (e *Ethereum) SetDefaults() {
	e.RPCEndpoint = "https://ethereum-rpc.publicnode.com/"
	e.InitialHeight = "21732451"
	e.SeedAddresses = []string{"0x2527d2ed1dd0e7de193cf121f1630caefc23ac70", "0xf70da97812cb96acdf810712aa562db8dfa3dbef"}
}
```

## Linting the service 
- `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest` will install the linter
- `golangci-lint run --verbose` will run the linter
## Generate code
- run `make generate` to generate code (enums & mocks), the command also installs the necessary tools.

## Starting the service
- duplicate the `.env.example` file and rename it to `.env`
- `make up` will start the service, after that you should be able to access the service on `localhost:8181`

## Tech stack
- HTTP for the API
- Docker
- golangci-lint for the linter 
- go-fumpt for stricter formatting
- enumer for generating string methods for enums
- testify for unit tests

## API 
- GET /v1/eth/blocks/current - returns the last block the service has parsed
- GET /v1/eth//addresses/:address/transactions - returns all transactions for the given address from all blocks parsed by the service 
- POST /v1/eth/addresses/subscribe - subscribes an address to the service, the service will start parsing all transactions for the given address (if any) from the last parsed block and onwards


## Structure
- cmd/api - contains the main package (entry point for the service) this includes both the api and a task background worker, ideally those would be two separate binaries, but because we are dealing with in memory storage that is not easily done with the current setup.
- specification/ - contains a postman collection and a basic openapi spec
- internal/
    - app  - contains entity/domain types and business logic
       - contract - contains all interfaces that are used in the app
       - entity - contains all domain types 
       - dto - contains all request and response types
    - pkg - contains packages and patterns that I typically use in my projects (those were not built for this project, excluding the `scheduler` and `ethclient` packages)
## Postman collection
A postman collection and a basic openapi spec is included in the specification dir


## Decisions 
- I decided to use a simple in memory storage for the sake of simplicity, this is not ideal for a production system, but it is good enough for this task, the way the service is structured allows for easy swapping of the storage layer with minimal changes.
- The transactions parsing logic is very basic, it doesn't handle the `withdraws` object in the block payload, as I wasn't sure if it counts as a transaction in the context of this task, if it does, it can be easily added.
  - The block parsing is done serial in a single go routine and a single block takes around 100ms (there is a basic benchmark) to parse (mostly spent on network roundtrips/api calls), this can be improved by parsing blocks in parallel, but that would require a more complex setup and would require a more complex storage and coordination layer. So I've opted with the simplest solution. The rpc calls can be batched to improve performance as well.
- I've used external libarires and a bunch of small packages that I've created over the years to speed up the boring parts of the project, e.g. configuration, logging, request decoding and etc. 
- I've built the project in such a way that another blockchain can be plugged in with minimal changes.
- The task scheduler is dead simple, in a production system I would opt for a ready made solution, or stick with this one but introduce some sort of leader election and have the scheduling be done in serial by the leader.
- I've introduced a simple Unit of Work (transactionmanger) pattern that has a noop implementation atm, with the idea that when a proper storage layer is used the bussness logic will gain atomicity and consistency guarantees.
- In the spirit of keeping things simple, I've opted out of implementing the kafka event publishing part of the task, but I've implemented a simple Publisher interface whenever a transaction is persisted, so that part can be easily added later on.
## This that can be improved 
- Track block confirmations (e.g. wait for 12 confirmations before considering a block as final and implement block state based on that)
- Include more information about individual transactions and the type of interaction, e.g. is it a transfer, a contract call and etc. Atm the service just stores the raw transaction data as is, token transfers and contract calls are not parsed properly.
- A batch transaction endpoints could be useful depending on usage patterns/backfills 
- Setup resilience patterns (cb, retries, fallbacks and etc.) to make the service more resilient and degradation more gracefully 
- e2e/contract testing 
- Improve code coverage in the service domain
- Instrument the service with Tracing/profiling/metrics
- Caching could be introduced to improve read performance (benchmark based on the expected read patters as cache invalidation can be tricky and caching might just complicate things for no real gain)


## Packages in the pkg directory
All packages in this dir are copied over and stripped from things that are not directly used in this project. Here is the list of packages and their purpose:
- configuration: Configuration package that reads from env vars and provides a typed config struct
- handler: A generic http handler that makes handlers have stricter input/output types
- render: A generic render package that can render json responses and handle errors in a consistent way
- requestdecoder: A wrapper around github.com/ggicci/httpin which provides a way to decode various http request params into a struct field, it also provides some json tag validation to prevent common mistakes (not including a json tag and potential field value overwrite exploits)
- server: A small wrapper around http.Server that provides a way to start/stop the server and handle shutdown gracefully 
- sloglog: slog context propagation & helpers
- syncx: provides a generic sync.Map wrapper
- requestvalidator: A wrapper around github.com/go-playground/validator/v10 that provides a way to validate structs and return a list of errors in a consistent way
- transactionmanger: A simple transaction manager that provides a way to rollback transactions in a consistent way.