# The "Word of Wisdom" TCP server with POW protection from DDOS attacks

Implementation:
- TCP server on Golang with DDOS attacks protection using [Proof Of Work](https://en.wikipedia.org/wiki/Proof_of_work)
- Hashcash algorithm as the implementation of POW
- Server uses the Word of Wisdom book service with in memory loaded quotes
- Client to solve the POW challenge for the demo
- Docker files provided for both the server and the client, managed with the Docker Compose

## Proof Of Work
Proof Of Work is a popular approach to protect systems from DDOS attacks.
The POW interface is designed here to use the challenge–response protocol with a direct interactive link between the requester (client) and the provider (server).

### Hashcash algorithm
The application uses the [Hashcash](https://en.wikipedia.org/wiki/Hashcash) POW algorithm. The Hashcash is known as standard implementation of POW system, it's clear and good enough to simply show the concept of POW approach.
This algorithm is used in the applications like - Bitcoin mining, Spam filters email client, etc.

The idea is that Client must Do some Work to solve the Challenge by hashing serialized string (Header) of Hashcash data (generated on Server side) and checking if it meets leading zero bits requirement, if not increment the Counter of the data itself and try again in a loop.
Ones the Work is succeeded Client should send the result data back to Server for the validation. If server POW system proves the Work done, Client could access resources.

The standard Header formant is `version:bits:date:source::random:counter` for the Data object with fields:
- *version*: Hashcash format version, should be `1`
- *bits*: number of leading zero bits required in the hashed Header (`20` bits by default)
- *date*: a sting with timestamp of sending the challenge, originally used to ensure the data is not expired
- *resource*: data string being transmitted, ie IP address of the client
- *random*: base-64 encoded random sequence to increase the data complexity
- *counter*: base-64 encoded value, represents the number of Work done iterations on client side (initialized as zero on the server side)

## Application Demo

Run cmd to see the demo
```
make start
```

The process flow:
- Docker Composer builds and runs Server and 2 Clients services in separate containers
- Server starts listening and serving client connections
- Each Client connects to the Server and makes a request for the Challenge
- Server receives the Challenge request and respond with newly generated Hashcash Data object
- Client receive Challenge response and Does some work (compute Hashcash Data) till it meets the Hashcash requirement 
- Client makes second request with computed Hashcash Data to get the Quote
- Server validates the result Hashcash Data, if it can prove the Work Done, return random Quote the Quotes Book
- Client prints Quote, waits a few seconds and start requesting new Challenge and Quote again

Use `Ctrl+C` to exit the demo.

## Testing

Run cmd run unit tests
```
make test
```

Some packages are covered with Unit Tests to show the testing approach.

- Mocks are generated for Interfaces using the [minimock](https://github.com/gojuno/minimock) tool.
Any change in Interface require the mock structures to be regenerated with cmd `make gen_code`.
The `minimock` tool must be installed for the mock code generation (see installation details in the tool source link).


- The [ginkgo](https://github.com/onsi/ginkgo) and [gomega](https://github.com/onsi/gomega) frameworks are used to write unit tests.

## TODO
- Implement workers pool for handling client connections on server side
- The TTL config for the work to be done should be provided from server side as well
- Complete tests for the rest of the packages
