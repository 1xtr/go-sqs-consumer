## [2.1.1](https://github.com/1xtr/go-sqs-consumer/compare/v2.1.0...v2.1.1) (2024-12-26)


### Bug Fixes

* throw args as objects pointer instead ([aae9e98](https://github.com/1xtr/go-sqs-consumer/commit/aae9e9828fbca7950068e6b716e1ad23906a94ec))

# [2.1.0](https://github.com/1xtr/go-sqs-consumer/compare/v2.0.0...v2.1.0) (2024-12-25)


### Features

* added batch handler ([cc0ef53](https://github.com/1xtr/go-sqs-consumer/commit/cc0ef5382e4871ee16183a643dd11692db4743f7))

# [2.0.0](https://github.com/1xtr/go-sqs-consumer/compare/v1.2.1...v2.0.0) (2024-12-25)


### Code Refactoring

* remove ctx from handler args ([772f351](https://github.com/1xtr/go-sqs-consumer/commit/772f35132e88a6502c18bb5b6d645a89d1e5aed1))


### Features

* add new logger ([59db623](https://github.com/1xtr/go-sqs-consumer/commit/59db6230f086b1a3b79fab538a68c4bcf4275ce9))


### BREAKING CHANGES

* added new logger

## [1.2.1](https://github.com/1xtr/go-sqs-consumer/compare/v1.2.0...v1.2.1) (2024-12-25)


### Bug Fixes

* added default value 100ms for poll delay ([a67988a](https://github.com/1xtr/go-sqs-consumer/commit/a67988ada3cbe143bb036b28eebbdb524ab4771d))

# [1.2.0](https://github.com/1xtr/go-sqs-consumer/compare/v1.1.0...v1.2.0) (2024-11-19)


### Features

* updated waitForProcessing method for use PollDelayInMs ([0c1be5b](https://github.com/1xtr/go-sqs-consumer/commit/0c1be5b5410bc3529353f5b4a3ccf2bdd18c19a3))

# [1.1.0](https://github.com/1xtr/go-sqs-consumer/compare/v1.0.1...v1.1.0) (2024-11-19)


### Features

* rename consumer log level environment "LOG_LEVEL" => "CONSUMER_LOG_LEVEL" ([3664912](https://github.com/1xtr/go-sqs-consumer/commit/3664912ab903391550d6b96cb690cc54960e36e7))
* set consumer default log level to `warning` ([73c6cc4](https://github.com/1xtr/go-sqs-consumer/commit/73c6cc4f1ec51c832c2c0ebf5a21cd0c26b9de3e))

## [1.0.1](https://github.com/1xtr/go-sqs-consumer/compare/v1.0.0...v1.0.1) (2024-11-19)


### Bug Fixes

* remove default value for visibilityTimeout ([3e48e63](https://github.com/1xtr/go-sqs-consumer/commit/3e48e63df5e3f9547a5c6f78278d8f098dd8d056))

# 1.0.0 (2024-11-19)


### Features

* added consumer ([a1c2b1b](https://github.com/1xtr/go-sqs-consumer/commit/a1c2b1beeceaaf1cfb5fd1330b2128cefafe3468))
