# Ernest CLI

## Synopsis

Ernest helps you deploy your vcloud infrastructure with a few shell commands.

## Build status

* Master [![CircleCI](https://circleci.com/gh/ernestio/ernest-cli/tree/master.svg?style=svg)](https://circleci.com/gh/ernestio/ernest-cli/tree/master)
* Develop [![CircleCI](https://circleci.com/gh/ernestio/ernest-cli/tree/develop.svg?style=svg)](https://circleci.com/gh/ernestio/ernest-cli/tree/develop)

## Installation

### Master Branch

Follow the instructions [here](http://docs.ernest.io/downloads/) to install the appropriate master branch binary of the cli.

### Develop Branch

In order to install the develop branch and have it running, you just need to run the following commands:

```
$ git clone git@github.com:ernestio/ernest-cli.git
$ cd ernest-cli
$ make
$ make install
```

## Set up

You need to set up your remote ernest instance
```
$ ernest-cli target "http://my.ernest.io"
```

## Run it

You can get help by running:
```
$ ernest-cli
```

And read our documentation about [how to use the CLI](http://docs.ernest.io/getting-started/)

## Running Tests

```
make test
```

## Contributing

Please read through our
[contributing guidelines](CONTRIBUTING.md).
Included are directions for opening issues, coding standards, and notes on
development.

Moreover, if your pull request contains patches or features, you must include
relevant unit tests.

## Versioning

For transparency into our release cycle and in striving to maintain backward
compatibility, this project is maintained under [the Semantic Versioning guidelines](http://semver.org/).

## Copyright and License

Code and documentation copyright since 2015 r3labs.io authors.

Code released under
[the Mozilla Public License Version 2.0](LICENSE).
