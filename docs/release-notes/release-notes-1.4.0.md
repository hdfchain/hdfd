# hdfd v1.4.0

This release of hdfd introduces a new consensus vote agenda which allows the
stakeholders to decide whether or not to activate changes needed to modify the
sequence lock handling which is required for providing full support for the
Lightning Network.  For those unfamiliar with the voting process in Decred, this
means that all code in order to make the necessary changes is already included
in this release, however its enforcement will remain dormant until the
stakeholders vote to activate it.

It also contains smart fee estimation, performance enhancements for block relay
and processing, a major internal restructuring of how unspent transaction
outputs are handled, support for whitelisting inbound peers to ensure service
for your own SPV (Simplified Payment Verification) wallets, various updates to
the RPC server such as a new method to query the state of the chain and more
easily supporting external RPC connections over TLS, infrastructure
improvements, and other quality assurance changes.

The following Decred Change Proposals (DCP) describes the proposed changes in detail:
- [DCP0004](https://github.com/hdfchain/dcps/blob/master/dcp-0004/dcp-0004.mediawiki)

**It is important for everyone to upgrade their software to this latest release
even if you don't intend to vote in favor of the agenda.**

## Downgrade Warning

The database format in v1.4.0 is not compatible with previous versions of the
software.  This only affects downgrades as users upgrading from previous
versions will see a lengthy one time database migration.

Once this migration has been completed, it will no longer be possible to
downgrade to a previous version of the software without having to delete the
database and redownload the chain.

## Notable Changes

### Fix Lightning Network Sequence Locks Vote

In order to fully support the Lightning Network, the current sequence lock
consensus rules need to be modified.  A new vote with the id `fixlnseqlocks` is
now available as of this release.  After upgrading, stakeholders may set their
preferences through their wallet or Voting Service Provider's (VSP) website.

### Smart Fee Estimation (`estimatesmartfee`)

A new RPC named `estimatesmartfee` is now available which returns a suitable
fee rate for transactions to use in order to have a high probability of them
being mined within a specified number of confirmations.  The estimation is based
on actual network usage and thus varies according to supply and demand.

This is important in the context of the Lightning Network (LN) and, more
generally, it provides services and users with a mechanism to choose how to
handle network congestion.  For example, payments that are high priority might
be willing to pay a higher fee to help ensure the transaction is mined more
quickly, while lower priority payments might be willing to wait longer in
exchange for paying a lower fee.  This estimation capability provides a way to
obtain a fee that will achieve the desired result with a high probability.

### Support for Whitelisting Inbound Peers

When peers are whitelisted via the `--whitelist` option, they will now be
allowed to connect even when they would otherwise exceed the maximum number of
peers.  This is highly useful in cases where users have configured their wallet
to use SPV mode and only connect to hdfd instances that they control for
increased privacy and guaranteed service.

### Several Speed Optimizations

Similar to previous releases, this release also contains several enhancements to
improve speed for the initial sync process, validation, and network operations.

In order to achieve these speedups, there is a lengthy one time database
migration, as previously mentioned, that typically takes anywhere from 30
minutes to an hour to complete depending on hardware.

#### Faster Tip Block Relay

Blocks that extend the current best chain are now relayed to the network
immediately after they pass the initial sanity and contextual checks, most
notably valid proof of work.  This allows blocks to propagate more quickly
throughout the network, which in turn improves vote times.

#### UTXO Set Restructuring

The way the unspent transaction outputs are handled internally has been
overhauled to significantly decrease the time it takes to validate blocks and
transactions.  While this has many benefits, probably the most important one
for most stakeholders is that votes can be cast more quickly which helps reduce
the number of missed votes.

### RPC Server Changes

#### New Chain State Query RPC (`getblockchaininfo`)

A new RPC named `getblockchaininfo` is now available which can be used to query
the state of the chain including details such as its overall verification
progress during initial sync, the maximum supported block size, and that status
of consensus changes (deployments) which require stakeholder votes.  See the
[JSON-RPC API Documentation](https://github.com/hdfchain/hdfd/blob/master/docs/json_rpc_api.mediawiki#getblockchaininfo)
for API details.

#### Removal of Vote Creation RPC (`createrawssgen`)

The deprecated `createrawssgen`, which was previously used to allow creating a
vote via RPC is no longer available.  Votes are time sensitive and thus it does
not make sense to create them offline.

#### Updates to Block and Transaction RPCs

The `getblock`, `getblockheader`, `getrawtransaction`, and
`searchrawtransactions` RPCs now contain additional information such as the
`extradata` field in the header, the `expiry` field in transactions, and the
`blockheight` and `blockindex` of  the block that contains a transaction if it
has been mined.  See the [JSON-RPC API Documentation](https://github.com/hdfchain/hdfd/blob/master/docs/json_rpc_api.md)
for API details.

#### Built-in Support for Enabling External TLS RPC Connections

A new command line parameter (`--altdnsnames`) and environment variable
(`HDFD_ALT_DNSNAMES`) can now be used before the first launch of drcd to specify
additional external IP addresses and DNS names to add during the certificate
creation that are permitted to connect to the RPC server via TLS.  Previously,
a separate tool was required to accomplish this configuration.

## Changelog

All commits since the last release may be viewed on GitHub [here](https://github.com/hdfchain/hdfd/compare/release-v1.3.0...release-v1.4.0).

### Protocol and network:

- chaincfg: Add checkpoints for 1.4.0 release ([hdfchain/hdfd#1547](https://github.com/hdfchain/hdfd/pull/1547))
- chaincfg: Introduce agenda for fixlnseqlocks vote ([hdfchain/hdfd#1578](https://github.com/hdfchain/hdfd/pull/1578))
- multi: Enable vote for DCP0004 ([hdfchain/hdfd#1579](https://github.com/hdfchain/hdfd/pull/1579))
- peer: Add support for specifying ua comments ([hdfchain/hdfd#1413](https://github.com/hdfchain/hdfd/pull/1413))
- blockmanager: Fast relay checked tip blocks ([hdfchain/hdfd#1443](https://github.com/hdfchain/hdfd/pull/1443))
- multi: Latest consensus active from simnet genesis ([hdfchain/hdfd#1482](https://github.com/hdfchain/hdfd/pull/1482))
- server: Always allow whitelisted inbound peers ([hdfchain/hdfd#1516](https://github.com/hdfchain/hdfd/pull/1516))

### Transaction relay (memory pool):

- blockmanager: handle txs in invalid blocks ([hdfchain/hdfd#1430](https://github.com/hdfchain/hdfd/pull/1430))
- mempool: Remove potential negative locktime check ([hdfchain/hdfd#1455](https://github.com/hdfchain/hdfd/pull/1455))
- mempool: Stake-related readability improvements ([hdfchain/hdfd#1456](https://github.com/hdfchain/hdfd/pull/1456))

### RPC:

- multi: Include additional fields on RPC tx results ([hdfchain/hdfd#1441](https://github.com/hdfchain/hdfd/pull/1441))
- rpcserver: Allow scripthash addrs in createrawsstx ([hdfchain/hdfd#1444](https://github.com/hdfchain/hdfd/pull/1444))
- rpcserver: Remove createrawssgen RPC ([hdfchain/hdfd#1448](https://github.com/hdfchain/hdfd/pull/1448))
- rpcclient: support getchaintips RPC ([hdfchain/hdfd#1469](https://github.com/hdfchain/hdfd/pull/1469))
- multi: Add getblockchaininfo rpc ([hdfchain/hdfd#1479](https://github.com/hdfchain/hdfd/pull/1479))
- rpcserver: Adds ability to allow alternative dns names for TLS ([hdfchain/hdfd#1476](https://github.com/hdfchain/hdfd/pull/1476))
- multi: Cleanup recent alt DNS names additions ([hdfchain/hdfd#1493](https://github.com/hdfchain/hdfd/pull/1493))
- multi: Cleanup getblock and getblockheader RPCs ([hdfchain/hdfd#1497](https://github.com/hdfchain/hdfd/pull/1497))
- multi: Return total chain work in RPC results ([hdfchain/hdfd#1498](https://github.com/hdfchain/hdfd/pull/1498))
- rpcserver: Improve GenerateNBlocks error message ([hdfchain/hdfd#1507](https://github.com/hdfchain/hdfd/pull/1507))
- rpcserver: Fix verify progress calculation ([hdfchain/hdfd#1508](https://github.com/hdfchain/hdfd/pull/1508))
- rpcserver: Fix sendrawtransaction error code ([hdfchain/hdfd#1512](https://github.com/hdfchain/hdfd/pull/1512))
- blockchain: Notify stake states after connected block ([hdfchain/hdfd#1515](https://github.com/hdfchain/hdfd/pull/1515))
- rpcserver: bump version to 5.0. ([hdfchain/hdfd#1531](https://github.com/hdfchain/hdfd/pull/1531))
- rpcclient: support getblockchaininfo RPC ([hdfchain/hdfd#1539](https://github.com/hdfchain/hdfd/pull/1539))
- rpcserver: update block template reconstruction ([hdfchain/hdfd#1567](https://github.com/hdfchain/hdfd/pull/1567))

### hdfd command-line flags and configuration:

- config: add --maxsameip to limit # of conns to same IP ([hdfchain/hdfd#1517](https://github.com/hdfchain/hdfd/pull/1517))

### Documentation:

- docs: Update docs for versioned modules ([hdfchain/hdfd#1391](https://github.com/hdfchain/hdfd/pull/1391))
- docs: Update for fees package ([hdfchain/hdfd#1540](https://github.com/hdfchain/hdfd/pull/1540))
- docs: Revamp main README.md and update docs ([hdfchain/hdfd#1447](https://github.com/hdfchain/hdfd/pull/1447))
- docs: Use relative versions in contrib checklist ([hdfchain/hdfd#1451](https://github.com/hdfchain/hdfd/pull/1451))
- docs: Use the correct binary name ([hdfchain/hdfd#1461](https://github.com/hdfchain/hdfd/pull/1461))
- docs: Add github pull request template ([hdfchain/hdfd#1474](https://github.com/hdfchain/hdfd/pull/1474))
- docs: Use unix line ending in mod hierarchy gv ([hdfchain/hdfd#1487](https://github.com/hdfchain/hdfd/pull/1487))
- docs: Add README badge and link for goreportcard ([hdfchain/hdfd#1492](https://github.com/hdfchain/hdfd/pull/1492))
- sampleconfig: Fix proxy typo ([hdfchain/hdfd#1513](https://github.com/hdfchain/hdfd/pull/1513))

### Developer-related package and module changes:

- release: Bump module versions and deps ([hdfchain/hdfd#1541](https://github.com/hdfchain/hdfd/pull/1541))
- build: Tidy module sums with go mod tidy ([hdfchain/hdfd#1408](https://github.com/hdfchain/hdfd/pull/1408))
- blockchain: update BestState ([hdfchain/hdfd#1416](https://github.com/hdfchain/hdfd/pull/1416))
- mempool: tweak trace logs ([hdfchain/hdfd#1429](https://github.com/hdfchain/hdfd/pull/1429))
- blockchain: Correct best pool size on disconnect ([hdfchain/hdfd#1431](https://github.com/hdfchain/hdfd/pull/1431))
- multi: Make use of new internal version package ([hdfchain/hdfd#1435](https://github.com/hdfchain/hdfd/pull/1435))
- peer: Protect handlePongMsg with p.statsMtx ([hdfchain/hdfd#1438](https://github.com/hdfchain/hdfd/pull/1438))
- limits: Make limits package internal ([hdfchain/hdfd#1436](https://github.com/hdfchain/hdfd/pull/1436))
- indexers: Remove unneeded existsaddrindex iface ([hdfchain/hdfd#1439](https://github.com/hdfchain/hdfd/pull/1439))
- blockchain: Reduce block availability assumptions ([hdfchain/hdfd#1442](https://github.com/hdfchain/hdfd/pull/1442))
- peer: Provide immediate queue inventory func ([hdfchain/hdfd#1443](https://github.com/hdfchain/hdfd/pull/1443))
- server: Add infrastruct for immediate inv relay ([hdfchain/hdfd#1443](https://github.com/hdfchain/hdfd/pull/1443))
- blockchain: Add new tip block checked notification ([hdfchain/hdfd#1443](https://github.com/hdfchain/hdfd/pull/1443))
- multi: remove chainState dependency in rpcserver ([hdfchain/hdfd#1417](https://github.com/hdfchain/hdfd/pull/1417))
- mining: remove chainState dependency ([hdfchain/hdfd#1418](https://github.com/hdfchain/hdfd/pull/1418))
- multi: remove chainState deps in server & cpuminer ([hdfchain/hdfd#1419](https://github.com/hdfchain/hdfd/pull/1419))
- blockmanager: remove block manager chain state ([hdfchain/hdfd#1420](https://github.com/hdfchain/hdfd/pull/1420))
- multi: move MinHighPriority into mining package ([hdfchain/hdfd#1421](https://github.com/hdfchain/hdfd/pull/1421))
- multi: add BlkTmplGenerator ([hdfchain/hdfd#1422](https://github.com/hdfchain/hdfd/pull/1422))
- multi: add cpuminerConfig ([hdfchain/hdfd#1423](https://github.com/hdfchain/hdfd/pull/1423))
- multi: Move update blk time to blk templ generator ([hdfchain/hdfd#1454](https://github.com/hdfchain/hdfd/pull/1454))
- multi: No stake height checks in check tx inputs ([hdfchain/hdfd#1457](https://github.com/hdfchain/hdfd/pull/1457))
- blockchain: Separate tx input stake checks ([hdfchain/hdfd#1452](https://github.com/hdfchain/hdfd/pull/1452))
- blockchain: Ensure no stake opcodes in tx sanity ([hdfchain/hdfd#1453](https://github.com/hdfchain/hdfd/pull/1453))
- blockchain: Move finalized tx func to validation ([hdfchain/hdfd#1465](https://github.com/hdfchain/hdfd/pull/1465))
- blockchain: Move unique coinbase func to validate ([hdfchain/hdfd#1466](https://github.com/hdfchain/hdfd/pull/1466))
- blockchain: Store interrupt channel with state ([hdfchain/hdfd#1467](https://github.com/hdfchain/hdfd/pull/1467))
- multi: Cleanup and optimize tx input check code ([hdfchain/hdfd#1468](https://github.com/hdfchain/hdfd/pull/1468))
- blockmanager: Avoid duplicate header announcements ([hdfchain/hdfd#1473](https://github.com/hdfchain/hdfd/pull/1473))
- dcrjson: additions for pay to contract hash ([hdfchain/hdfd#1260](https://github.com/hdfchain/hdfd/pull/1260))
- multi: Break blockchain dependency on dcrjson ([hdfchain/hdfd#1488](https://github.com/hdfchain/hdfd/pull/1488))
- chaincfg: Unexport internal errors ([hdfchain/hdfd#1489](https://github.com/hdfchain/hdfd/pull/1489))
- multi: Cleanup the unsupported hdfwallet commands ([hdfchain/hdfd#1478](https://github.com/hdfchain/hdfd/pull/1478))
- multi: Rename ThresholdState to NextThresholdState ([hdfchain/hdfd#1494](https://github.com/hdfchain/hdfd/pull/1494))
- dcrjson: Add listtickets command ([hdfchain/hdfd#1267](https://github.com/hdfchain/hdfd/pull/1267))
- multi: Add started and done reorg notifications ([hdfchain/hdfd#1495](https://github.com/hdfchain/hdfd/pull/1495))
- blockchain: Remove unused CheckWorklessBlockSanity ([hdfchain/hdfd#1496](https://github.com/hdfchain/hdfd/pull/1496))
- blockchain: Simplify block template checking ([hdfchain/hdfd#1499](https://github.com/hdfchain/hdfd/pull/1499))
- blockchain: Only mark nodes modified when modified ([hdfchain/hdfd#1503](https://github.com/hdfchain/hdfd/pull/1503))
- blockchain: Cleanup and optimize stake node logic ([hdfchain/hdfd#1504](https://github.com/hdfchain/hdfd/pull/1504))
- blockchain: Separate full data context checks ([hdfchain/hdfd#1509](https://github.com/hdfchain/hdfd/pull/1509))
- blockchain: Reverse utxo set semantics ([hdfchain/hdfd#1471](https://github.com/hdfchain/hdfd/pull/1471))
- blockchain: Convert to direct single-step reorgs ([hdfchain/hdfd#1500](https://github.com/hdfchain/hdfd/pull/1500))
- multi: Migration for utxo set semantics reversal ([hdfchain/hdfd#1520](https://github.com/hdfchain/hdfd/pull/1520))
- blockchain: Make version 5 update atomic ([hdfchain/hdfd#1529](https://github.com/hdfchain/hdfd/pull/1529))
- blockchain: Simplify force head reorgs ([hdfchain/hdfd#1526](https://github.com/hdfchain/hdfd/pull/1526))
- secp256k1: Correct edge case in deterministic sign ([hdfchain/hdfd#1533](https://github.com/hdfchain/hdfd/pull/1533))
- dcrjson: Add gettransaction txtype/ticketstatus ([hdfchain/hdfd#1276](https://github.com/hdfchain/hdfd/pull/1276))
- txscript: Use ScriptBuilder more ([hdfchain/hdfd#1519](https://github.com/hdfchain/hdfd/pull/1519))
- fees: Add estimator package ([hdfchain/hdfd#1434](https://github.com/hdfchain/hdfd/pull/1434))
- multi: Integrate fee estimation ([hdfchain/hdfd#1434](https://github.com/hdfchain/hdfd/pull/1434))

### Testing and Quality Assurance:

- multi: Use temp directories for database tests ([hdfchain/hdfd#1404](https://github.com/hdfchain/hdfd/pull/1404))
- multi: Only use module-scoped data in tests ([hdfchain/hdfd#1405](https://github.com/hdfchain/hdfd/pull/1405))
- blockchain: Use temp dirs for fullblocks test ([hdfchain/hdfd#1406](https://github.com/hdfchain/hdfd/pull/1406))
- database: Use module-scoped data in iface tests ([hdfchain/hdfd#1407](https://github.com/hdfchain/hdfd/pull/1407))
- travis: Update for Go1.11 and module builds ([hdfchain/hdfd#1415](https://github.com/hdfchain/hdfd/pull/1415))
- indexers: Use testable bucket for existsaddrindex ([hdfchain/hdfd#1440](https://github.com/hdfchain/hdfd/pull/1440))
- txscript: group numeric encoding tests with their opcodes ([hdfchain/hdfd#1382](https://github.com/hdfchain/hdfd/pull/1382))
- txscript: add p2sh opcode tests ([hdfchain/hdfd#1381](https://github.com/hdfchain/hdfd/pull/1381))
- txscript: add stake opcode tests ([hdfchain/hdfd#1383](https://github.com/hdfchain/hdfd/pull/1383))
- main: add address encoding magic constants test ([hdfchain/hdfd#1458](https://github.com/hdfchain/hdfd/pull/1458))
- chaingen: Only revoke missed tickets once ([hdfchain/hdfd#1484](https://github.com/hdfchain/hdfd/pull/1484))
- chaingen/fullblocktests: Add disapproval tests ([hdfchain/hdfd#1485](https://github.com/hdfchain/hdfd/pull/1485))
- multi: Resurrect regression network ([hdfchain/hdfd#1480](https://github.com/hdfchain/hdfd/pull/1480))
- multi: Use regression test network in unit tests ([hdfchain/hdfd#1481](https://github.com/hdfchain/hdfd/pull/1481))
- main: move cert tests to a separated file ([hdfchain/hdfd#1502](https://github.com/hdfchain/hdfd/pull/1502))
- mempool: Accept test mungers for create signed tx ([hdfchain/hdfd#1576](https://github.com/hdfchain/hdfd/pull/1576))
- mempool: Implement test harness seq lock calc ([hdfchain/hdfd#1577](https://github.com/hdfchain/hdfd/pull/1577))

### Misc:

- release: Bump for 1.4 release cycle ([hdfchain/hdfd#1414](https://github.com/hdfchain/hdfd/pull/1414))
- multi: Make changes suggested by Go 1.11 gofmt -s ([hdfchain/hdfd#1415](https://github.com/hdfchain/hdfd/pull/1415))
- build: Remove dep toml and lock file ([hdfchain/hdfd#1460](https://github.com/hdfchain/hdfd/pull/1460))
- docker: Update to go 1.11 ([hdfchain/hdfd#1463](https://github.com/hdfchain/hdfd/pull/1463))
- build: Support MacOS sed for obtaining module list ([hdfchain/hdfd#1483](https://github.com/hdfchain/hdfd/pull/1483))
- multi: Correct a few typos found by misspell ([hdfchain/hdfd#1490](https://github.com/hdfchain/hdfd/pull/1490))
- multi: Address some golint complaints ([hdfchain/hdfd#1491](https://github.com/hdfchain/hdfd/pull/1491))
- multi: Remove unused code ([hdfchain/hdfd#1505](https://github.com/hdfchain/hdfd/pull/1505))
- release: Bump siphash version to v1.2.1 ([hdfchain/hdfd#1538](https://github.com/hdfchain/hdfd/pull/1538))
- release: Bump module versions and deps ([hdfchain/hdfd#1541](https://github.com/hdfchain/hdfd/pull/1541))
- Fix required version of stake module ([hdfchain/hdfd#1549](https://github.com/hdfchain/hdfd/pull/1549))
- release: Tidy module files with published versions ([hdfchain/hdfd#1543](https://github.com/hdfchain/hdfd/pull/1543))
- mempool: Fix required version of mining module ([hdfchain/hdfd#1551](https://github.com/hdfchain/hdfd/pull/1551))

### Code Contributors (alphabetical order):

- Corey Osman
- Dave Collins
- David Hill
- Dmitry Fedorov
- Donald Adu-Poku
- ggoranov
- githubsands
- J Fixby
- Jonathan Chappelow
- Josh Rickmar
- Matheus Degiovani
- Sarlor
- zhizhongzhiwai