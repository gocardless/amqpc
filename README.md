amqpc
=====

amqpc is a tool to test your AMQP broker.  
It has been developed and tested with RabbitMQ but any AMQP broker should work.  
We put it together very quickly to benchmark and test our RabbitMQ HA Cluster, so it might be unstable.  

## Usage

```bash
$ amqpc -h

amqpc is CLI tool for testing AMQP brokers
Usage:
  Consumer : amqpc [options] -c exchange routingkey queue
  Producer : amqpc [options] -p exchange routingkey [message]

Options:
  -c=true: Act as a consumer
  -ct="amqpc-consumer": AMQP consumer tag (should not be blank)
  -g=1: Concurrency
  -i=500: (Producer only) Interval at which send messages (in ms)
  -n=0: (Producer only) Number of messages to send
  -p=false: Act as a producer
  -r=true: Wait for the publisher confirmation before exiting
  -t="direct": Exchange type - direct|fanout|topic|x-custom
  -u="amqp://guest:guest@localhost:5672/": AMQP URI
  ```

## Installation

Right now it's not available in any package manager.
To build it your need to have a Go compiler.

```bash
$ go get
$ go build
```

## Examples

```bash
$ # Start one consumer
$ amqpc -u amqp://guest:guest@localhost:5672/ -c your-exchange routing-key your-queue
$ # Start one producer
$ amqpc -u amqp://guest:guest@localhost:5672/ -p your-exchange routing-key your-message
$ # Start 10 producers
$ amqpc -u amqp://guest:guest@localhost:5672/ -g 10 -p your-exchange routing-key your-message
$ # Start 10 producers, each one will be sending 100 at a rate of 1msg/s
$ amqpc -u amqp://guest:guest@localhost:5672/ -g 10 -i 1000 -n 100 -p your-exchange routing-key your-message
```

## TODO

* Package management
* Tests
* Publisher confirms
