# hdfd v1.5.1

This is a patch release of hdfd to address a minor memory leak with authenticated RPC websocket clients on intermittent connections.   It also updates the `hdfctl` utility to include the new `auditreuse` hdfwallet command.

## Changelog

This patch release consists of 4 commits from 3 contributors which total to 4 files changed, 27 additional lines of code, and 6 deleted lines of code.

All commits since the last release may be viewed on GitHub [here](https://github.com/hdfchain/hdfd/compare/release-v1.5.0...release-v1.5.1).

### RPC:

- rpcwebsocket: Remove client from missed maps ([hdfchain/hdfd#2049](https://github.com/hdfchain/hdfd/pull/2049))
- rpcwebsocket: Use nonblocking messages and ntfns ([hdfchain/hdfd#2050](https://github.com/hdfchain/hdfd/pull/2050))

### hdfctl utility changes:

- hdfctl: Update hdfwallet RPC types package ([hdfchain/hdfd#2051](https://github.com/hdfchain/hdfd/pull/2051))

### Misc:

- release: Bump for 1.5.1([hdfchain/hdfd#2052](https://github.com/hdfchain/hdfd/pull/2052))

### Code Contributors (alphabetical order):

- Dave Collins
- Josh Rickmar
- Matheus Degiovani