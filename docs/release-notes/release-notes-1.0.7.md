## hdfd v1.0.7

This release of hdfd primarily contains improvements to the infrastructure and
other quality assurance changes that are bringing us closer to providing full
support for Lightning Network.

A lot of work required for Lightning Network support went into getting the
required code merged into the upstream project, btcd, which now fully supports
it.  These changes also must be synced and integrated with hdfd as well and
therefore many of the changes in this release are related to that process.

## Notable Changes

### Dust check removed from stake transactions

The standard policy for regular transactions is to reject any transactions that
have outputs so small that they cost more to the network than their value.  This
behavior is desirable for regular transactions, however it was also being
applied to vote and revocation transactions which could lead to a situation
where stake pools with low fees could result in votes and revocations having
difficulty being mined.

This check has been changed to only apply to regular transactions now in order
to prevent any issues.  Stake transactions have several other checks that make
this one unnecessary for them.

### New `feefilter` peer-to-peer message

A new optional peer-to-peer message named `feefilter` has been added that allows
peers to inform others about the minimum transaction fee rate they are willing
to accept.  This will enable peers to avoid notifying others about transactions
they will not accept anyways and therefore can result in a significant bandwidth
savings.

### Bloom filter service bit enforcement

Peers that are configured to disable bloom filter support will now disconnect
remote peers that send bloom filter related commands rather than simply ignoring
them.  This allows any light clients that do not observe the service bit to
potentially find another peer that provides the service.  Additionally, remote
peers that have negotiated a high enough protocol version to observe the service
bit and still send bloom filter related commands anyways will now be banned.


## Changelog

All commits since the last release may be viewed on GitHub [here](https://github.com/hdfchain/hdfd/compare/v1.0.5...v1.0.7).

### Protocol and network:
- Allow reorg of block one [hdfchain/hdfd#745](https://github.com/hdfchain/hdfd/pull/745)
- blockchain: use the time source [hdfchain/hdfd#747](https://github.com/hdfchain/hdfd/pull/747)
- peer: Strictly enforce bloom filter service bit [hdfchain/hdfd#768](https://github.com/hdfchain/hdfd/pull/768)
- wire/peer: Implement feefilter p2p message [hdfchain/hdfd#779](https://github.com/hdfchain/hdfd/pull/779)
- chaincfg: update checkpoints for 1.0.7 release  [hdfchain/hdfd#816](https://github.com/hdfchain/hdfd/pull/816)

### Transaction relay (memory pool):
- mempool: Break dependency on chain instance [hdfchain/hdfd#754](https://github.com/hdfchain/hdfd/pull/754)
- mempool: unexport the mutex [hdfchain/hdfd#755](https://github.com/hdfchain/hdfd/pull/755)
- mempool: Add basic test harness infrastructure [hdfchain/hdfd#756](https://github.com/hdfchain/hdfd/pull/756)
- mempool: Improve tx input standard checks [hdfchain/hdfd#758](https://github.com/hdfchain/hdfd/pull/758)
- mempool: Update comments for dust calcs [hdfchain/hdfd#764](https://github.com/hdfchain/hdfd/pull/764)
- mempool: Only perform standard dust checks on regular transactions  [hdfchain/hdfd#806](https://github.com/hdfchain/hdfd/pull/806)

### RPC:
- Fix gettxout includemempool handling [hdfchain/hdfd#738](https://github.com/hdfchain/hdfd/pull/738)
- Improve help text for getmininginfo [hdfchain/hdfd#748](https://github.com/hdfchain/hdfd/pull/748)
- rpcserverhelp: update TicketFeeInfo help [hdfchain/hdfd#801](https://github.com/hdfchain/hdfd/pull/801)
- blockchain: Improve getstakeversions efficiency [hdfchain/hdfd#81](https://github.com/hdfchain/hdfd/pull/813)

### hdfd command-line flags:
- config: introduce new flags to accept/reject non-std transactions [hdfchain/hdfd#757](https://github.com/hdfchain/hdfd/pull/757)
- config: Add --whitelist option [hdfchain/hdfd#352](https://github.com/hdfchain/hdfd/pull/352)
- config: Improve config file handling [hdfchain/hdfd#802](https://github.com/hdfchain/hdfd/pull/802)
- config: Improve blockmaxsize check [hdfchain/hdfd#810](https://github.com/hdfchain/hdfd/pull/810)

### hdfctl:
- Add --walletrpcserver option [hdfchain/hdfd#736](https://github.com/hdfchain/hdfd/pull/736)

### Documentation
- docs: add commit prefix notes  [hdfchain/hdfd#788](https://github.com/hdfchain/hdfd/pull/788)

### Developer-related package changes:
- blockchain: check errors and remove ineffectual assignments [hdfchain/hdfd#689](https://github.com/hdfchain/hdfd/pull/689)
- stake: less casting [hdfchain/hdfd#705](https://github.com/hdfchain/hdfd/pull/705)
- blockchain: chainstate only needs prev block hash [hdfchain/hdfd#706](https://github.com/hdfchain/hdfd/pull/706)
- remove dead code [hdfchain/hdfd#715](https://github.com/hdfchain/hdfd/pull/715)
- Use btclog for determining valid log levels [hdfchain/hdfd#738](https://github.com/hdfchain/hdfd/pull/738)
- indexers: Minimize differences with upstream code [hdfchain/hdfd#742](https://github.com/hdfchain/hdfd/pull/742)
- blockchain: Add median time to state snapshot [hdfchain/hdfd#753](https://github.com/hdfchain/hdfd/pull/753)
- blockmanager: remove unused GetBlockFromHash function [hdfchain/hdfd#761](https://github.com/hdfchain/hdfd/pull/761)
- mining: call CheckConnectBlock directly [hdfchain/hdfd#762](https://github.com/hdfchain/hdfd/pull/762)
- blockchain: add missing error code entries [hdfchain/hdfd#763](https://github.com/hdfchain/hdfd/pull/763)
- blockchain: Sync main chain flag on ProcessBlock [hdfchain/hdfd#767](https://github.com/hdfchain/hdfd/pull/767)
- blockchain: Remove exported CalcPastTimeMedian func [hdfchain/hdfd#770](https://github.com/hdfchain/hdfd/pull/770)
- blockchain: check for error [hdfchain/hdfd#772](https://github.com/hdfchain/hdfd/pull/772)
- multi: Optimize by removing defers [hdfchain/hdfd#782](https://github.com/hdfchain/hdfd/pull/782)
- blockmanager: remove unused logBlockHeight [hdfchain/hdfd#787](https://github.com/hdfchain/hdfd/pull/787)
- dcrutil: Replace DecodeNetworkAddress with DecodeAddress [hdfchain/hdfd#746](https://github.com/hdfchain/hdfd/pull/746)
- txscript: Force extracted addrs to compressed [hdfchain/hdfd#775](https://github.com/hdfchain/hdfd/pull/775)
- wire: Remove legacy transaction decoding [hdfchain/hdfd#794](https://github.com/hdfchain/hdfd/pull/794)
- wire: Remove dead legacy tx decoding code [hdfchain/hdfd#796](https://github.com/hdfchain/hdfd/pull/796)
- mempool/wire: Don't make policy decisions in wire [hdfchain/hdfd#797](https://github.com/hdfchain/hdfd/pull/797)
- dcrjson: Remove unused cmds & types [hdfchain/hdfd#795](https://github.com/hdfchain/hdfd/pull/795)
- dcrjson: move cmd types [hdfchain/hdfd#799](https://github.com/hdfchain/hdfd/pull/799)
- multi: Separate tx serialization type from version [hdfchain/hdfd#798](https://github.com/hdfchain/hdfd/pull/798)
- dcrjson: add Unconfirmed field to dcrjson.GetAccountBalanceResult [hdfchain/hdfd#812](https://github.com/hdfchain/hdfd/pull/812)
- multi: Error descriptions should be lowercase [hdfchain/hdfd#771](https://github.com/hdfchain/hdfd/pull/771)
- blockchain: cast to int64  [hdfchain/hdfd#817](https://github.com/hdfchain/hdfd/pull/817)

### Testing and Quality Assurance:
- rpcserver: Upstream sync to add basic RPC tests [hdfchain/hdfd#750](https://github.com/hdfchain/hdfd/pull/750)
- rpctest: Correct several issues tests and joins [hdfchain/hdfd#751](https://github.com/hdfchain/hdfd/pull/751)
- rpctest: prevent process leak due to panics [hdfchain/hdfd#752](https://github.com/hdfchain/hdfd/pull/752)
- rpctest: Cleanup resources on failed setup [hdfchain/hdfd#759](https://github.com/hdfchain/hdfd/pull/759)
- rpctest: Use ports based on the process id [hdfchain/hdfd#760](https://github.com/hdfchain/hdfd/pull/760)
- rpctest/deps: Update dependencies and API [hdfchain/hdfd#765](https://github.com/hdfchain/hdfd/pull/765)
- rpctest: Gate rpctest-based behind a build tag [hdfchain/hdfd#766](https://github.com/hdfchain/hdfd/pull/766)
- mempool: Add test for max orphan entry eviction [hdfchain/hdfd#769](https://github.com/hdfchain/hdfd/pull/769)
- fullblocktests: Add more consensus tests [hdfchain/hdfd#77](https://github.com/hdfchain/hdfd/pull/773)
- fullblocktests: Sync upstream block validation [hdfchain/hdfd#774](https://github.com/hdfchain/hdfd/pull/774)
- rpctest: fix a harness range bug in syncMempools [hdfchain/hdfd#778](https://github.com/hdfchain/hdfd/pull/778)
- secp256k1: Add regression tests for field.go [hdfchain/hdfd#781](https://github.com/hdfchain/hdfd/pull/781)
- secp256k1: Sync upstream test consolidation [hdfchain/hdfd#783](https://github.com/hdfchain/hdfd/pull/783)
- txscript: Correct p2sh hashes in json test data  [hdfchain/hdfd#785](https://github.com/hdfchain/hdfd/pull/785)
- txscript: Replace CODESEPARATOR json test data [hdfchain/hdfd#786](https://github.com/hdfchain/hdfd/pull/786)
- txscript: Remove multisigdummy from json test data [hdfchain/hdfd#789](https://github.com/hdfchain/hdfd/pull/789)
- txscript: Remove max money from json test data [hdfchain/hdfd#790](https://github.com/hdfchain/hdfd/pull/790)
- txscript: Update signatures in json test data [hdfchain/hdfd#791](https://github.com/hdfchain/hdfd/pull/791)
- txscript: Use native encoding in json test data [hdfchain/hdfd#792](https://github.com/hdfchain/hdfd/pull/792)
- rpctest: Store logs and data in same path [hdfchain/hdfd#780](https://github.com/hdfchain/hdfd/pull/780)
- txscript: Cleanup reference test code  [hdfchain/hdfd#793](https://github.com/hdfchain/hdfd/pull/793)

### Misc:
- Update deps to pull in additional logging changes [hdfchain/hdfd#734](https://github.com/hdfchain/hdfd/pull/734)
- Update markdown files for GFM changes [hdfchain/hdfd#744](https://github.com/hdfchain/hdfd/pull/744)
- blocklogger: Show votes, tickets, & revocations [hdfchain/hdfd#784](https://github.com/hdfchain/hdfd/pull/784)
- blocklogger: Remove STransactions from transactions calculation [hdfchain/hdfd#811](https://github.com/hdfchain/hdfd/pull/811)

### Contributors (alphabetical order):

- Alex Yocomm-Piatt
- Atri Viss
- Chris Martin
- Dave Collins
- David Hill
- Donald Adu-Poku
- Jimmy Song
- John C. Vernaleo
- Jolan Luff
- Josh Rickmar
- Olaoluwa Osuntokun
- Marco Peereboom
