# Distributed arithmetic expression evaluator

## Technology stack

* Golang
* Vue js
* PostgreSQL
* RabbitMQ
* Docker
* Docker-compose

## Description

The user wants to calculate arithmetic expressions. He enters the string 2 + 2 * 2 and wants the answer to be 6. But our addition and multiplication operations (also division and subtraction) take a “very, very” long time to complete. Therefore, the option in which the user makes an http request and receives the result as a response is impossible. Moreover: the calculation of each such operation in our “alternative reality” takes “giant” computing power. Accordingly, we must be able to perform each action separately and we can scale this system by adding computing power to our system in the form of new “machines”. Therefore, when a user sends an expression, he receives an expression identifier in response and can, at some periodicity, check with the server whether the expression has been counted? If the expression is finally evaluated, he will get the result. Remember that some parts of an arphimetic expression can be evaluated in parallel.


## Requirements

### Front-end part

* Arithmetic expression input form. The user enters an arithmetic expression and sends a POST http request with this expression to the back-end. Note: Requests must be idempotent. A unique identifier is added to requests. If a user sends a request with an identifier that has already been sent and accepted for processing, the response is 200. Possible response options:

    * [200]. The expression was successfully accepted, parsed and accepted for processing
    * [400]. The expression is invalid
    * [500]. Something is wrong on the back-end. As a response, you need to return the id of the expression accepted for execution.

* Page with a list of expressions in the form of a list with expressions. Each entry on the page contains a status, an expression, the date it was created, and the date the calculation was completed. The page receives data with a GET http request from the back-end
* A page with a list of operations in the form of pairs: operation name + execution time (editable field). As already stated in the problem statement, our operations take “as if for a very long time.” The page receives data with a GET http request from the back-end. The user can configure the operation execution time and save the changes.
* Page with a list of computing capabilities. The page receives data with a GET http request from the server in the form of pairs: the name of the computing resource + the operation performed on it.

Requirements:

* The orchestrator can be restarted without losing state. We store all expressions in the DBMS.
* The orchestrator must keep track of tasks that take too long to complete (the computer may also go offline) and make them available for computation again.

### Back-end part

Consists of 2 elements:

* The server, which receives an arithmetic expression, translates it into a set of sequential tasks and ensures the order in which they are executed. From now on we will call it Orchestrator.
* A computer that can receive orchestrator task, execute it and return the result to the server. In what follows we will call it Agent.

Orchestrator
Server that has the following endpoints:

* Adding an arithmetic expression evaluation.
* Getting a list of expressions with statuses.
* Getting the value of an expression by its identifier.
* Obtaining a list of available operations with their execution time.
* Receiving a task for execution.
* Receiving the result of data processing.

The Daemon Agent
Receives an expression to be evaluated from the server, evaluates it and sends the result of the expression to the server. When starting, the daemon launches several goroutines, each of which acts as an independent computer. The number of goroutines is controlled by an environment variable.

## Getting Started

### Prerequisites

To start project you need:
* [Docker](https://www.docker.com/)

### Building from source

```shell
git clone https://github.com/ExcaliBBur/yandexProject.git
cd yandexProject

# build project
docker-compose build

# After running the command below, check that the server starts correctly. 
# If this is not the case - repeat, this is due to the nature of docker.
docker-compose up server -d

# $ - quantity of workers
docker-compose up --scale agent=$ agent -d

docker-compose up frontend -d
```

## How to work

Go to the [localhost:8081](http://localhost:8081). Here you go.

Endpoints you can see at [localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

## How it works

![Scheme](/scheme.drawio.svg)
