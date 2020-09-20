hdfd
====

[![Build Status](http://img.shields.io/travis/hdfchain/hdfd.svg)]
(https://travis-ci.org/hdfchain/hdfd) [![ISC License]
(http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)

hdfd is a Hdfchain full node implementation written in Go (golang).

This acts as a chain daemon for the [Hdfchain](https://clkj.ltd) cryptocurrency.
hdfd maintains the entire past transactional ledger of Hdfchain and allows
 relaying of transactions to other Hdfchain nodes across the world.  To read more 
about Hdfchain please see the 
[project documentation](https://docs.clkj.ltd/#overview).

Note: To send or receive funds and join Proof-of-Stake mining, you will also need
[hdfwallet](https://github.com/hdfchain/hdfwallet).

This project is currently under active development and is in a Beta state.  It
is extremely stable and has been in production use since February 2016.  

It is forked from [btcd](https://github.com/btcsuite/btcd) which is a bitcoin
full node implementation written in Go.  btcd is a ongoing project under active 
development.  Because hdfd is constantly synced with btcd codebase, it will 
get the benefit of btcd's ongoing upgrades to peer and connection handling, 
database optimization and other blockchain related technology improvements.

## Requirements

[Go](http://golang.org) 1.7 or newer.

## Getting Started

- hdfd (and utilities) will now be installed in either ```$GOROOT/bin``` or
  ```$GOPATH/bin``` depending on your configuration.  If you did not already
  add the bin directory to your system path during Go installation, we
  recommend you do so now.

## Updating

#### Windows

Install a newer MSI

#### Linux/BSD/MacOSX/POSIX - Build from Source

- **Glide**

  Glide is used to manage project dependencies and provide reproducible builds.
  To install:

  `go get -u github.com/Masterminds/glide`

Unfortunately, the use of `glide` prevents a handy tool such as `go get` from
automatically downloading, building, and installing the source in a single
command.  Instead, the latest project and dependency sources must be first
obtained manually with `git` and `glide`, and then `go` is used to build and
install the project.

**Getting the source**:

For a first time installation, the project and dependency sources can be
obtained manually with `git` and `glide` (create directories as needed):

```
git clone https://github.com/hdfchain/hdfd $GOPATH/src/github.com/hdfchain/hdfd
cd $GOPATH/src/github.com/hdfchain/hdfd
glide install
go install $(glide nv)
```

To update an existing source tree, pull the latest changes and install the
matching dependencies:

```
cd $GOPATH/src/github.com/hdfchain/hdfd
git pull
glide install
go install $(glide nv)
```

For more information about hdfchain and how to set up your software please go to
our docs page at [docs.clkj.ltd](https://docs.clkj.ltd/getting-started/beginner-guide/).  

## Contact 

If you have any further questions you can find us at:

- irc.freenode.net (channel #hdfchain)
- [webchat](https://webchat.freenode.net/?channels=hdfchain)
- forum.clkj.ltd
- hdfchain.slack.com

## Issue Tracker

The [integrated github issue tracker](https://github.com/hdfchain/hdfd/issues)
is used for this project.

## Documentation

The documentation is a work-in-progress.  It is located in the [docs](https://github.com/hdfchain/hdfd/tree/master/docs) folder.

## License

hdfd is licensed under the [copyfree](http://copyfree.org) ISC License.
