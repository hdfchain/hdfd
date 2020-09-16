# hdfd v1.1.0

This release of hdfd primarily introduces a new consensus vote agenda which
allows the stakeholders to decide whether or not to activate the features needed
for providing full support for Lightning Network.  For those unfamiliar with the
voting process in Decred, this means that all code in order to support these
features is already included in this release, however its enforcement will
remain dormant until the stakeholders vote to activate it.

The following Decred Change Proposals (DCPs) describe the proposed changes in detail:
- [DCP0002](https://github.com/hdfchain/dcps/blob/master/dcp-0002/dcp-0002.mediawiki)
- [DCP0003](https://github.com/hdfchain/dcps/blob/master/dcp-0003/dcp-0003.mediawiki)

It is important for everyone to upgrade their software to this latest release
even if you don't intend to vote in favor of the agenda.

## Notable Changes

### Lightning Network Features Vote

In order to fully support many of the benefits that the Lightning Network will
bring, there are some primitives that involve changes to the current consensus
that need to be enabled.  A new vote with the id `lnfeatures` is now available
as of this release.  After upgrading, stakeholders may set their preferences
through their wallet or stake pool's website.

### Transaction Finality Policy

The standard policy for transaction relay has been changed to use the median
time of the past several blocks instead of the current network adjusted time
when examining lock times to determine if a transaction is final.  This provides
a more deterministic check across all peers and prevents the possibility of
miners attempting to game the timestamps in order to include more transactions.

Consensus enforcement of this change relies on the result of the aforementioned
`lnfeatures` vote.

### Relative Time Locks Policy

The standard policy for transaction relay has been modified to enforce relative
lock times for version 2 transactions via their sequence numbers and a new
`OP_CHECKSEQUENCEVERIFY` opcode.

Consensus enforcement of this change relies on the result of the aforementioned
`lnfeatures` vote.

### OP_SHA256 Opcode

In order to better support cross-chain interoperability, a new opcode to compute
the SHA-256 hash is being proposed.  Since this opcode is implemented as a hard
fork, it will not be available for use in scripts unless the aforementioned
`lnfeatures` vote passes.

## Changelog

All commits since the last release may be viewed on GitHub [here](https://github.com/hdfchain/hdfd/compare/v1.0.7...v1.1.0).

### Protocol and network:
- chaincfg: update checkpoints for 1.1.0 release [hdfchain/hdfd#850](https://github.com/hdfchain/hdfd/pull/850)
- chaincfg: Introduce agenda for v5 lnfeatures vote [hdfchain/hdfd#848](https://github.com/hdfchain/hdfd/pull/848)
- txscript: Introduce OP_SHA256 [hdfchain/hdfd#851](https://github.com/hdfchain/hdfd/pull/851)
- wire: Decrease num allocs when decoding headers [hdfchain/hdfd#861](https://github.com/hdfchain/hdfd/pull/861)
- blockchain: Implement enforced relative seq locks [hdfchain/hdfd#864](https://github.com/hdfchain/hdfd/pull/864)
- txscript: Implement CheckSequenceVerify [hdfchain/hdfd#864](https://github.com/hdfchain/hdfd/pull/864)
- multi: Enable vote for DCP0002 and DCP0003 [hdfchain/hdfd#855](https://github.com/hdfchain/hdfd/pull/855)

### Transaction relay (memory pool):
- mempool: Use median time for tx finality checks [hdfchain/hdfd#860](https://github.com/hdfchain/hdfd/pull/860)
- mempool: Enforce relative sequence locks [hdfchain/hdfd#864](https://github.com/hdfchain/hdfd/pull/864)
- policy/mempool: Enforce CheckSequenceVerify opcode [hdfchain/hdfd#864](https://github.com/hdfchain/hdfd/pull/864)

### RPC:
- rpcserver: check whether ticketUtx was found [hdfchain/hdfd#824](https://github.com/hdfchain/hdfd/pull/824)
- rpcserver: return rule error on rejected raw tx [hdfchain/hdfd#808](https://github.com/hdfchain/hdfd/pull/808)

### hdfd command-line flags:
- config: Extend --profile cmd line option to allow interface to be specified [hdfchain/hdfd#838](https://github.com/hdfchain/hdfd/pull/838)

### Documentation
- docs: rpcapi format update [hdfchain/hdfd#807](https://github.com/hdfchain/hdfd/pull/807)
- config: export sampleconfig for use by dcrinstall [hdfchain/hdfd#834](https://github.com/hdfchain/hdfd/pull/834)
- sampleconfig: Add package README and doc.go [hdfchain/hdfd#835](https://github.com/hdfchain/hdfd/pull/835)
- docs: create entry for getstakeversions in rpcapi [hdfchain/hdfd#819](https://github.com/hdfchain/hdfd/pull/819)
- docs: crosscheck and update all rpc doc entries [hdfchain/hdfd#847](https://github.com/hdfchain/hdfd/pull/847)
- docs: update git commit messages section heading [hdfchain/hdfd#863](https://github.com/hdfchain/hdfd/pull/863)

### Developer-related package changes:
- Fix and regenerate precomputed secp256k1 curve [hdfchain/hdfd#823](https://github.com/hdfchain/hdfd/pull/823)
- dcrec: use hardcoded datasets in tests [hdfchain/hdfd#822](https://github.com/hdfchain/hdfd/pull/822)
- Use dchest/blake256  [hdfchain/hdfd#827](https://github.com/hdfchain/hdfd/pull/827)
- glide: use jessevdk/go-flags for consistency [hdfchain/hdfd#833](https://github.com/hdfchain/hdfd/pull/833)
- multi: Error descriptions are in lower case [hdfchain/hdfd#842](https://github.com/hdfchain/hdfd/pull/842)
- txscript: Rename OP_SHA256 to OP_BLAKE256 [hdfchain/hdfd#840](https://github.com/hdfchain/hdfd/pull/840)
- multi: Abstract standard verification flags [hdfchain/hdfd#852](https://github.com/hdfchain/hdfd/pull/852)
- chain: Remove memory block node pruning [hdfchain/hdfd#858](https://github.com/hdfchain/hdfd/pull/858)
- txscript: Add API to parse atomic swap contracts [hdfchain/hdfd#862](https://github.com/hdfchain/hdfd/pull/862)

### Testing and Quality Assurance:
- Test against go 1.9 [hdfchain/hdfd#836](https://github.com/hdfchain/hdfd/pull/836)
- dcrec: remove testify dependency [hdfchain/hdfd#829](https://github.com/hdfchain/hdfd/pull/829)
- mining_test: add edge conditions from btcd [hdfchain/hdfd#831](https://github.com/hdfchain/hdfd/pull/831)
- stake: Modify ticket tests to use chaincfg params [hdfchain/hdfd#844](https://github.com/hdfchain/hdfd/pull/844)
- blockchain: Modify tests to use chaincfg params [hdfchain/hdfd#845](https://github.com/hdfchain/hdfd/pull/845)
- blockchain: Cleanup various tests [hdfchain/hdfd#843](https://github.com/hdfchain/hdfd/pull/843)
- Ensure run_tests.sh local fails correctly when gometalinter errors [hdfchain/hdfd#846](https://github.com/hdfchain/hdfd/pull/846)
- peer: fix logic race in peer connection test [hdfchain/hdfd#865](https://github.com/hdfchain/hdfd/pull/865)

### Misc:
- glide: sync deps [hdfchain/hdfd#837](https://github.com/hdfchain/hdfd/pull/837)
- Update hdfchain deps for v1.1.0 [hdfchain/hdfd#868](https://github.com/hdfchain/hdfd/pull/868)
- Bump for v1.1.0 [hdfchain/hdfd#867](https://github.com/hdfchain/hdfd/pull/867)

### Code Contributors (alphabetical order):

- Alex Yocom-Piatt
- Dave Collins
- David Hill
- Donald Adu-Poku
- Jason Zavaglia
- Jean-Christophe Mincke
- Jolan Luff
- Josh Rickmar
