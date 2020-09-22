# hdfd v1.1.2

This release of hdfd primarily contains performance enhancements, infrastructure
improvements, and other quality assurance changes.

While it is not visible in this release, significant infrastructure work has
also been done this release cycle towards porting the Lightning Network (LN)
daemon which will ultimately allow LN payments to be backed by Hdfchain.

## Notable Changes

### Faster Block Validation

A significant portion of block validation involves handling the stake tickets
which form an integral part of Hdfchain's hybrid proof-of-work and proof-of-stake
system.  The code which handles this portion of validation has been
significantly optimized in this release such that overall block validation is
up to approximately 3 times faster depending on the specific underlying hardware
configuration.  This also has a noticeable impact on the speed of the initial
block download process as well as how quickly votes for winning tickets are
submitted to the network.

### Data Carrier Transaction Standardness Policy

The standard policy for transaction relay of data carrier transaction outputs
has been modified to support canonically-encoded small data pushes.  These
outputs are also known as `OP_RETURN` or `nulldata` outputs.  In particular,
single byte small integers data pushes (0-16) are now supported.

## Changelog

All commits since the last release may be viewed on GitHub [here](https://github.com/hdfchain/hdfd/compare/v1.1.0...v1.1.2).

### Protocol and network:
- chaincfg: update checkpoints for 1.1.2 release [hdfchain/hdfd#946](https://github.com/hdfchain/hdfd/pull/946)
- chaincfg: Rename one of the testnet seeders [hdfchain/hdfd#873](https://github.com/hdfchain/hdfd/pull/873)
- stake: treap index perf improvement [hdfchain/hdfd#853](https://github.com/hdfchain/hdfd/pull/853)
- stake: ticket expiry perf improvement [hdfchain/hdfd#853](https://github.com/hdfchain/hdfd/pull/853)

### Transaction relay (memory pool):

- txscript: Correct nulldata standardness check [hdfchain/hdfd#935](https://github.com/hdfchain/hdfd/pull/935)

### RPC:

- rpcserver: searchrawtransactions skip first input for vote tx [hdfchain/hdfd#859](https://github.com/hdfchain/hdfd/pull/859)
- multi: update stakebase tx vin[0] structure [hdfchain/hdfd#859](https://github.com/hdfchain/hdfd/pull/859)
- rpcserver: Fix empty ssgen verbose results [hdfchain/hdfd#871](https://github.com/hdfchain/hdfd/pull/871)
- rpcserver: check for error in getwork request [hdfchain/hdfd#898](https://github.com/hdfchain/hdfd/pull/898)
- multi: Add NoSplitTransaction to purchaseticket [hdfchain/hdfd#904](https://github.com/hdfchain/hdfd/pull/904)
- rpcserver: avoid nested decodescript p2sh addrs [hdfchain/hdfd#929](https://github.com/hdfchain/hdfd/pull/929)
- rpcserver: skip generating certs when nolisten set [hdfchain/hdfd#932](https://github.com/hdfchain/hdfd/pull/932)
- rpc: Add localaddr and relaytxes to getpeerinfo [hdfchain/hdfd#933](https://github.com/hdfchain/hdfd/pull/933)
- rpcserver: update handleSendRawTransaction error handling [hdfchain/hdfd#939](https://github.com/hdfchain/hdfd/pull/939)

### hdfd command-line flags:

- config: add --nofilelogging option [hdfchain/hdfd#872](https://github.com/hdfchain/hdfd/pull/872)

### Documentation:

- rpcclient: Remove docker info from README.md [hdfchain/hdfd#886](https://github.com/hdfchain/hdfd/pull/886)
- bloom: Fix link in README [hdfchain/hdfd#922](https://github.com/hdfchain/hdfd/pull/922)
- doc: tiny fix url [hdfchain/hdfd#928](https://github.com/hdfchain/hdfd/pull/928)
- doc: update go version for example test run in readme [hdfchain/hdfd#936](https://github.com/hdfchain/hdfd/pull/936)

### Developer-related package changes:

- multi: Drop glide, use dep [hdfchain/hdfd#818](https://github.com/hdfchain/hdfd/pull/818)
- txsort: Implement stable tx sorting package  [hdfchain/hdfd#940](https://github.com/hdfchain/hdfd/pull/940)
- coinset: Remove package [hdfchain/hdfd#888](https://github.com/hdfchain/hdfd/pull/888)
- base58: Use new github.com/hdfchain/base58 package [hdfchain/hdfd#888](https://github.com/hdfchain/hdfd/pull/888)
- certgen: Move self signed certificate code into package [hdfchain/hdfd#879](https://github.com/hdfchain/hdfd/pull/879)
- certgen: Add doc.go and README.md [hdfchain/hdfd#883](https://github.com/hdfchain/hdfd/pull/883)
- rpcclient: Allow request-scoped cancellation during Connect [hdfchain/hdfd#880](https://github.com/hdfchain/hdfd/pull/880)
- rpcclient: Import dcrrpcclient repo into rpcclient directory [hdfchain/hdfd#880](https://github.com/hdfchain/hdfd/pull/880)
- rpcclient: json unmarshal into unexported embedded pointer  [hdfchain/hdfd#941](https://github.com/hdfchain/hdfd/pull/941)
- bloom: Copy github.com/hdfchain/dcrutil/bloom to bloom package [hdfchain/hdfd#881](https://github.com/hdfchain/hdfd/pull/881)
- Improve gitignore [hdfchain/hdfd#887](https://github.com/hdfchain/hdfd/pull/887)
- dcrutil: Import dcrutil repo under dcrutil directory [hdfchain/hdfd#888](https://github.com/hdfchain/hdfd/pull/888)
- hdkeychain: Move to github.com/hdfchain/hdfd/hdkeychain [hdfchain/hdfd#892](https://github.com/hdfchain/hdfd/pull/892)
- stake: Add IsStakeSubmission [hdfchain/hdfd#907](https://github.com/hdfchain/hdfd/pull/907)
- txscript: Require SHA256 secret hashes for atomic swaps [hdfchain/hdfd#930](https://github.com/hdfchain/hdfd/pull/930)

### Testing and Quality Assurance:

- gometalinter: run on subpkgs too [hdfchain/hdfd#878](https://github.com/hdfchain/hdfd/pull/878)
- travis: test Gopkg.lock [hdfchain/hdfd#889](https://github.com/hdfchain/hdfd/pull/889)
- hdkeychain: Work around go vet issue with examples [hdfchain/hdfd#890](https://github.com/hdfchain/hdfd/pull/890)
- bloom: Add missing import to examples [hdfchain/hdfd#891](https://github.com/hdfchain/hdfd/pull/891)
- bloom: workaround go vet issue in example [hdfchain/hdfd#895](https://github.com/hdfchain/hdfd/pull/895)
- tests: make lockfile test work locally [hdfchain/hdfd#894](https://github.com/hdfchain/hdfd/pull/894)
- peer: Avoid goroutine leaking during handshake timeout [hdfchain/hdfd#909](https://github.com/hdfchain/hdfd/pull/909)
- travis: add gosimple linter [hdfchain/hdfd#897](https://github.com/hdfchain/hdfd/pull/897)
- multi: Handle detected data race conditions [hdfchain/hdfd#920](https://github.com/hdfchain/hdfd/pull/920)
- travis: add ineffassign linter [hdfchain/hdfd#896](https://github.com/hdfchain/hdfd/pull/896)
- rpctest: Choose flags based on provided params [hdfchain/hdfd#937](https://github.com/hdfchain/hdfd/pull/937)

### Misc:

- gofmt [hdfchain/hdfd#876](https://github.com/hdfchain/hdfd/pull/876)
- dep: sync third-party deps [hdfchain/hdfd#877](https://github.com/hdfchain/hdfd/pull/877)
- Bump for v1.1.2 [hdfchain/hdfd#916](https://github.com/hdfchain/hdfd/pull/916)
- dep: Use upstream jrick/bitset [hdfchain/hdfd#899](https://github.com/hdfchain/hdfd/pull/899)
- blockchain: removed unused funcs and vars [hdfchain/hdfd#900](https://github.com/hdfchain/hdfd/pull/900)
- blockchain: remove unused file [hdfchain/hdfd#900](https://github.com/hdfchain/hdfd/pull/900)
- rpcserver: nil pointer dereference when submit orphan block [hdfchain/hdfd#906](https://github.com/hdfchain/hdfd/pull/906)
- multi: remove unused funcs and vars [hdfchain/hdfd#901](https://github.com/hdfchain/hdfd/pull/901)

### Code Contributors (alphabetical order):

- Alex Yocom-Piatt
- Dave Collins
- David Hill
- detailyang
- Donald Adu-Poku
- Federico Gimenez
- Jason Zavaglia
- John C. Vernaleo
- Jonathan Chappelow
- Jolan Luff
- Josh Rickmar
- Maninder Lall
- Matheus Degiovani
- Nicola Larosa
- Samarth Hattangady
- Ugwueze Onyekachi Michael
