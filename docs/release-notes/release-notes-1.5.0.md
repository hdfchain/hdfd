# hdfd v1.5.0

This release of hdfd introduces a large number of updates.  Some of the key highlights are:

* A new consensus vote agenda which allows the stakeholders to decide whether or not to activate support for block header commitments
* More efficient block filters
* Significant improvements to the mining infrastructure including asynchronous work notifications
* Major performance enhancements for transaction script validation
* Automatic external IP address discovery
* Support for IPv6 over Tor
* Various updates to the RPC server such as:
  * A new method to query information about the network
  * A method to retrieve the new version 2 block filters
  * More calls available to limited access users
* Infrastructure improvements
* Quality assurance changes

For those unfamiliar with the voting process in Hdfchain, all code in order to support block header commitments is already included in this release, however its enforcement will remain dormant until the stakeholders vote to activate it.

For reference, block header commitments were originally proposed and approved for initial implementation via the following Politeia proposal:
- [Block Header Commitments Consensus Change](https://proposals.clkj.ltd/proposals/0a1ff846ec271184ea4e3a921a3ccd8d478f69948b984445ee1852f272d54c58)


The following Hdfchain Change Proposal (DCP) describes the proposed changes in detail and provides a full technical specification:
- [DCP0005](https://github.com/hdfchain/dcps/blob/master/dcp-0005/dcp-0005.mediawiki)

**It is important for everyone to upgrade their software to this latest release even if you don't intend to vote in favor of the agenda.**

## Downgrade Warning

The database format in v1.5.0 is not compatible with previous versions of the software.  This only affects downgrades as users upgrading from previous versions will see a one time database migration.

Once this migration has been completed, it will no longer be possible to downgrade to a previous version of the software without having to delete the database and redownload the chain.

## Notable Changes

### Block Header Commitments Vote

A new vote with the id `headercommitments` is now available as of this release.  After upgrading, stakeholders may set their preferences through their wallet or Voting Service Provider's (VSP) website.

The primary goal of this change is to increase the security and efficiency of lightweight clients, such as Hdfchainiton in its lightweight mode and the dcrandroid/dcrios mobile wallets, as well as add infrastructure that paves the
way for several future scalability enhancements.

A high level overview aimed at a general audience including a cost benefit analysis can be found in the  [Politeia proposal](https://proposals.clkj.ltd/proposals/0a1ff846ec271184ea4e3a921a3ccd8d478f69948b984445ee1852f272d54c58).

In addition, a much more in-depth treatment can be found in the [motivation section of DCP0005](https://github.com/hdfchain/dcps/blob/master/dcp-0005/dcp-0005.mediawiki#motivation).

### Version 2 Block Filters

The block filters used by lightweight clients, such as SPV (Simplified Payment Verification) wallets, have been updated to improve their efficiency, ergonomics, and include additional information such as the full ticket
commitment script.  The new block filters are version 2.  The older version 1 filters are now deprecated and scheduled to be removed in the next release, so consumers should update to the new filters as soon as possible.

An overview of block filters can be found in the [block filters section of DCP0005](https://github.com/hdfchain/dcps/blob/master/dcp-0005/dcp-0005.mediawiki#block-filters).

Also, the specific contents and technical specification of the new version 2 block filters is available in the
[version 2 block filters section of DCP0005](https://github.com/hdfchain/dcps/blob/master/dcp-0005/dcp-0005.mediawiki#version-2-block-filters).

Finally, there is a one time database update to build and store the new filters for all existing historical blocks which will likely take a while to complete (typically around 8 to 10 minutes on HDDs and 4 to 5 minutes on SSDs).

### Mining Infrastructure Overhaul

The mining infrastructure for building block templates and delivering the work to miners has been significantly overhauled to improve several aspects as follows:

* Support asynchronous background template generation with intelligent vote propagation handling
* Improved handling of chain reorganizations necessary when the current tip is unable to obtain enough votes
* Current state synchronization
* Near elimination of stale templates when new blocks and votes are received
* Subscriptions for streaming template updates

The standard [getwork RPC](https://github.com/hdfchain/hdfd/blob/master/docs/json_rpc_api.mediawiki#getwork) that PoW miners currently use to perform the mining process has been updated to make use of this new infrastructure, so existing PoW miners will seamlessly get the vast majority of benefits without requiring any updates.

However, in addition, a new [notifywork RPC](https://github.com/hdfchain/hdfd/blob/master/docs/json_rpc_api.mediawiki#notifywork) is now available that allows miners to register for work to be delivered
asynchronously as it becomes available via a WebSockets [work notification](https://github.com/hdfchain/hdfd/blob/master/docs/json_rpc_api.mediawiki#work).  These notifications include the same information that `getwork` provides along with an additional `reason` parameter which allows the miners to make better decisions about when they should instruct workers to discard the current template immediately or should be allowed to finish their current round before being provided with the new template.

Miners are highly encouraged to update their software to make use of the new asynchronous notification infrastructure since it is more robust, efficient, and faster than polling `getwork` to manually determine the aforementioned conditions.

The following is a non-exhaustive overview that highlights the major benefits of the changes for both cases:

- Requests for updated templates during the normal mining process in between tip   changes will now be nearly instant instead of potentially taking several seconds to build the new template on the spot
- When the chain tip changes, requesting a template will now attempt to wait until either all votes have been received or a timeout occurs prior to handing out a template which is beneficial for PoW miners, PoS miners, and the network as a whole
- PoW miners are much less likely to end up with template with less than the max number of votes which means they are less likely to receive a reduced subsidy
- PoW miners will be much less likely to receive stale templates during chain tip changes due to vote propagation
- PoS voters whose votes end up arriving to the miner slightly slower than the minimum number required are much less likely to have their votes excluded despite having voted simply due to propagation delay

PoW miners who choose to update their software, pool or otherwise, to make use of the asynchronous work notifications will receive additional benefits such as:

- Ability to start mining a new block sooner due to receiving updated work as soon as it becomes available
- Immediate notification with new work that includes any votes that arrive late
- Periodic notifications with new work that include new transactions only when there have actually been new transaction
- Simplified interface code due to removal of the need for polling and manually checking the work bytes for special cases such as the number of votes

**NOTE: Miners that are not rolling the timestamp field as they mine should ensure their software is upgraded to roll the timestamp to the latest timestamp each time they hand work out to a miner.  This helps ensure the block timestamps are as accurate as possible.**

### Transaction Script Validation Optimizations

Transaction script validation has been almost completely rewritten to significantly improve its speed and reduce the number of memory allocations. While this has many more benefits than enumerated here, probably the most
important ones for most stakeholders are:

- Votes can be cast more quickly which helps reduce the number of missed votes
- Blocks are able to propagate more quickly throughout the network, which in turn further improves votes times
- The initial sync process is around 20-25% faster

### Automatic External IP Address Discovery

In order for nodes to fully participate in the peer-to-peer network, they must be publicly accessible and made discoverable by advertising their external IP address.  This is typically made slightly more complicated since most users run their nodes on networks behind Network Address Translation (NAT).

Previously, in addition to configuring the network firewall and/or router to allow inbound connections to port 9108 and forwarding the port to the internal IP address running hdfd, it was also required to manually set the public external IP address via the `--externalip` CLI option.

This release will now make use of other nodes on the network in a decentralized fashion to automatically discover the external IP address, so it is no longer necessary to manually set CLI option for the vast majority of users.

### Tor IPv6 Support

It is now possible to resolve and connect to IPv6 peers over Tor in addition to the existing IPv4 support.

### RPC Server Changes

#### New Version 2 Block Filter Query RPC (`getcfilterv2`)

A new RPC named `getcfilterv2` is now available which can be used to retrieve the version 2 [block filter](https://github.com/hdfchain/dcps/blob/master/dcp-0005/dcp-0005.mediawiki#Block_Filters)
for a given block along with its associated inclusion proof.  See the [getcfilterv2 JSON-RPC API Documentation](https://github.com/hdfchain/hdfd/blob/master/docs/json_rpc_api.mediawiki#getcfilterv2)
for API details.

#### New Network Information Query RPC (`getnetworkinfo`)

A new RPC named `getnetworkinfo` is now available which can be used to query information related to the peer-to-peer network such as the protocol version, the local time offset, the number of current connections, the supported network protocols, the current transaction relay fee, and the external IP addresses for
the local interfaces.  See the [getnetworkinfo JSON-RPC API Documentation](https://github.com/hdfchain/hdfd/blob/master/docs/json_rpc_api.mediawiki#getnetworkinfo) for API details.

#### Updates to Chain State Query RPC (`getblockchaininfo`)

The `difficulty` field of the `getblockchaininfo` RPC is now deprecated in favor of a new field named `difficultyratio` which matches the result returned by the `getdifficulty` RPC.

See the [getblockchaininfo JSON-RPC API Documentation](https://github.com/hdfchain/hdfd/blob/master/docs/json_rpc_api.mediawiki#getblockchaininfo) for API details.

#### New Optional Version Parameter on Script Decode RPC (`decodescript`)

The `decodescript` RPC now accepts an additional optional parameter to specify the script version.  The only currently supported script version in Hdfchain is version 0 which means decoding scripts with versions other than 0 will be seen as non standard.

#### Removal of Deprecated Block Template RPC (`getblocktemplate`)

The previously deprecated `getblocktemplate` RPC is no longer available.  All known miners are already using the preferred `getwork` RPC since Hdfchain's block header supports more than enough nonce space to keep mining hardware busy without needing to resort to building custom templates with less efficient extra nonce coinbase workarounds.

#### Additional RPCs Available To Limited Access Users

The following RPCs that were previously unavailable to the limited access RPC user are now available to it:

- `estimatefee`
- `estimatesmartfee`
- `estimatestakediff`
- `existsaddress`
- `existsaddresses`
- `existsexpiredtickets`
- `existsliveticket`
- `existslivetickets`
- `existsmempoltxs`
- `existsmissedtickets`
- `getblocksubsidy`
- `getcfilter`
- `getcoinsupply`
- `getheaders`
- `getstakedifficulty`
- `getstakeversioninfo`
- `getstakeversions`
- `getvoteinfo`
- `livetickets`
- `missedtickets`
- `rebroadcastmissed`
- `rebroadcastwinners`
- `ticketfeeinfo`
- `ticketsforaddress`
- `ticketvwap`
- `txfeeinfo`

### Single Mining State Request

The peer-to-peer protocol message to request the current mining state (`getminings`) is used when peers first connect to retrieve all known votes for the current tip block.  This is only useful when the peer first connects because all future votes will be relayed once the connection has been established.  Consequently, nodes will now only respond to a single mining state request.  Subsequent requests are ignored.

### Developer Go Modules

A full suite of versioned Go modules (essentially code libraries) are now available for use by applications written in Go that wish to create robust software with reproducible, verifiable, and verified builds.

These modules are used to build hdfd itself and are therefore well maintained, tested, documented, and relatively efficient.

## Changelog

This release consists of 600 commits from 17 contributors which total to 537 files changed, 41494 additional lines of code, and 29215 deleted lines of code.

All commits since the last release may be viewed on GitHub [here](https://github.com/hdfchain/hdfd/compare/release-v1.4.0...release-v1.5.0).

### Protocol and network:

- chaincfg: Add checkpoints for 1.5.0 release ([hdfchain/hdfd#1924](https://github.com/hdfchain/hdfd/pull/1924))
- chaincfg: Introduce agenda for header cmtmts vote ([hdfchain/hdfd#1904](https://github.com/hdfchain/hdfd/pull/1904))
- multi: Implement combined merkle root and vote ([hdfchain/hdfd#1906](https://github.com/hdfchain/hdfd/pull/1906))
- blockchain: Implement v2 block filter storage ([hdfchain/hdfd#1906](https://github.com/hdfchain/hdfd/pull/1906))
- gcs/blockcf2: Implement v2 block filter creation ([hdfchain/hdfd#1906](https://github.com/hdfchain/hdfd/pull/1906))
- wire: Implement getcfilterv2/cfilterv2 messages ([hdfchain/hdfd#1906](https://github.com/hdfchain/hdfd/pull/1906))
- peer: Implement getcfilterv2/cfilterv2 listeners ([hdfchain/hdfd#1906](https://github.com/hdfchain/hdfd/pull/1906))
- server: Implement getcfilterv2 ([hdfchain/hdfd#1906](https://github.com/hdfchain/hdfd/pull/1906))
- multi: Implement header commitments and vote ([hdfchain/hdfd#1906](https://github.com/hdfchain/hdfd/pull/1906))
- server: Remove instead of disconnect node ([hdfchain/hdfd#1644](https://github.com/hdfchain/hdfd/pull/1644))
- server: limit getminingstate requests ([hdfchain/hdfd#1678](https://github.com/hdfchain/hdfd/pull/1678))
- peer: Prevent last block height going backwards ([hdfchain/hdfd#1747](https://github.com/hdfchain/hdfd/pull/1747))
- connmgr: Add ability to remove pending connections ([hdfchain/hdfd#1724](https://github.com/hdfchain/hdfd/pull/1724))
- connmgr: Add cancellation of pending requests ([hdfchain/hdfd#1724](https://github.com/hdfchain/hdfd/pull/1724))
- connmgr: Check for canceled connection before connect ([hdfchain/hdfd#1724](https://github.com/hdfchain/hdfd/pull/1724))
- multi: add automatic network address discovery ([hdfchain/hdfd#1522](https://github.com/hdfchain/hdfd/pull/1522))
- connmgr: add TorLookupIPContext, deprecate TorLookupIP ([hdfchain/hdfd#1849](https://github.com/hdfchain/hdfd/pull/1849))
- connmgr: support resolving ipv6 hosts over Tor ([hdfchain/hdfd#1908](https://github.com/hdfchain/hdfd/pull/1908))

### Transaction relay (memory pool):

- mempool: Reject same block vote double spends ([hdfchain/hdfd#1597](https://github.com/hdfchain/hdfd/pull/1597))
- mempool: Limit max vote double spends exactly ([hdfchain/hdfd#1596](https://github.com/hdfchain/hdfd/pull/1596))
- mempool: Optimize pool double spend check ([hdfchain/hdfd#1561](https://github.com/hdfchain/hdfd/pull/1561))
- txscript: Tighten standardness pubkey checks ([hdfchain/hdfd#1649](https://github.com/hdfchain/hdfd/pull/1649))
- mempool: drop container/list for simple FIFO ([hdfchain/hdfd#1681](https://github.com/hdfchain/hdfd/pull/1681))
- mempool: remove unused error return value ([hdfchain/hdfd#1785](https://github.com/hdfchain/hdfd/pull/1785))
- mempool: Add ErrorCode to returned TxRuleErrors ([hdfchain/hdfd#1901](https://github.com/hdfchain/hdfd/pull/1901))

### Mining:

- mining: Optimize get the block's votes tx ([hdfchain/hdfd#1563](https://github.com/hdfchain/hdfd/pull/1563))
- multi: add BgBlkTmplGenerator ([hdfchain/hdfd#1424](https://github.com/hdfchain/hdfd/pull/1424))
- mining: Remove unnecessary notify goroutine ([hdfchain/hdfd#1708](https://github.com/hdfchain/hdfd/pull/1708))
- mining: Improve template key handling ([hdfchain/hdfd#1709](https://github.com/hdfchain/hdfd/pull/1709))
- mining:  fix scheduled template regen ([hdfchain/hdfd#1717](https://github.com/hdfchain/hdfd/pull/1717))
- miner: Improve background generator lifecycle ([hdfchain/hdfd#1715](https://github.com/hdfchain/hdfd/pull/1715))
- cpuminer: No speed monitor on discrete mining ([hdfchain/hdfd#1716](https://github.com/hdfchain/hdfd/pull/1716))
- mining: Run vote ntfn in a separate goroutine ([hdfchain/hdfd#1718](https://github.com/hdfchain/hdfd/pull/1718))
- mining: Overhaul background template generator ([hdfchain/hdfd#1748](https://github.com/hdfchain/hdfd/pull/1748))
- mining: Remove unused error return value ([hdfchain/hdfd#1859](https://github.com/hdfchain/hdfd/pull/1859))
- cpuminer: Fix off-by-one issues in nonce handling ([hdfchain/hdfd#1865](https://github.com/hdfchain/hdfd/pull/1865))
- mining: Remove dead code ([hdfchain/hdfd#1882](https://github.com/hdfchain/hdfd/pull/1882))
- mining: Remove unused extra nonce update code ([hdfchain/hdfd#1883](https://github.com/hdfchain/hdfd/pull/1883))
- mining: Minor cleanup of aggressive mining path ([hdfchain/hdfd#1888](https://github.com/hdfchain/hdfd/pull/1888))
- mining: Remove unused error codes ([hdfchain/hdfd#1889](https://github.com/hdfchain/hdfd/pull/1889))
- mining: fix data race ([hdfchain/hdfd#1894](https://github.com/hdfchain/hdfd/pull/1894))
- mining: fix data race ([hdfchain/hdfd#1896](https://github.com/hdfchain/hdfd/pull/1896))
- cpuminer: fix race ([hdfchain/hdfd#1899](https://github.com/hdfchain/hdfd/pull/1899))
- cpuminer: Improve speed stat tracking ([hdfchain/hdfd#1921](https://github.com/hdfchain/hdfd/pull/1921))
- rpcserver/mining: Use bg tpl generator for getwork ([hdfchain/hdfd#1922](https://github.com/hdfchain/hdfd/pull/1922))
- mining: Export TemplateUpdateReason ([hdfchain/hdfd#1923](https://github.com/hdfchain/hdfd/pull/1923))
- multi: Add tpl update reason to work ntfns ([hdfchain/hdfd#1923](https://github.com/hdfchain/hdfd/pull/1923))
- mining: Store block templates given by notifywork ([hdfchain/hdfd#1949](https://github.com/hdfchain/hdfd/pull/1949))

### RPC:

- dcrjson: add cointype to WalletInfoResult ([hdfchain/hdfd#1606](https://github.com/hdfchain/hdfd/pull/1606))
- rpcclient: Introduce v2 module using wallet types ([hdfchain/hdfd#1608](https://github.com/hdfchain/hdfd/pull/1608))
- rpcserver: Update for dcrjson/v2 ([hdfchain/hdfd#1612](https://github.com/hdfchain/hdfd/pull/1612))
- rpcclient: Add EstimateSmartFee ([hdfchain/hdfd#1641](https://github.com/hdfchain/hdfd/pull/1641))
- rpcserver: remove unused quit chan ([hdfchain/hdfd#1629](https://github.com/hdfchain/hdfd/pull/1629))
- rpcserver: Undeprecate getwork ([hdfchain/hdfd#1635](https://github.com/hdfchain/hdfd/pull/1635))
- rpcserver: Add difficultyratio to getblockchaininfo ([hdfchain/hdfd#1630](https://github.com/hdfchain/hdfd/pull/1630))
- multi:  add version arg to decodescript rpc ([hdfchain/hdfd#1731](https://github.com/hdfchain/hdfd/pull/1731))
- dcrjson: Remove API breaking change ([hdfchain/hdfd#1778](https://github.com/hdfchain/hdfd/pull/1778))
- rpcclient: Add GetMasterPubkey ([hdfchain/hdfd#1777](https://github.com/hdfchain/hdfd/pull/1777))
- multi: add getnetworkinfo rpc ([hdfchain/hdfd#1536](https://github.com/hdfchain/hdfd/pull/1536))
- rpcserver: Better error message ([hdfchain/hdfd#1861](https://github.com/hdfchain/hdfd/pull/1861))
- multi: update limited user rpcs ([hdfchain/hdfd#1870](https://github.com/hdfchain/hdfd/pull/1870))
- multi: make rebroadcast winners & missed ws only ([hdfchain/hdfd#1872](https://github.com/hdfchain/hdfd/pull/1872))
- multi: remove getblocktemplate ([hdfchain/hdfd#1736](https://github.com/hdfchain/hdfd/pull/1736))
- rpcserver: Match tx filter on ticket commitments ([hdfchain/hdfd#1881](https://github.com/hdfchain/hdfd/pull/1881))
- rpcserver: don't use activeNetParams ([hdfchain/hdfd#1733](https://github.com/hdfchain/hdfd/pull/1733))
- rpcserver: update rpcAskWallet rpc set ([hdfchain/hdfd#1892](https://github.com/hdfchain/hdfd/pull/1892))
- rpcclient: close the unused response body ([hdfchain/hdfd#1905](https://github.com/hdfchain/hdfd/pull/1905))
- rpcclient: Support getcfilterv2 JSON-RPC ([hdfchain/hdfd#1906](https://github.com/hdfchain/hdfd/pull/1906))
- multi: add notifywork rpc ([hdfchain/hdfd#1410](https://github.com/hdfchain/hdfd/pull/1410))
- rpcserver: Cleanup getvoteinfo RPC ([hdfchain/hdfd#2005](https://github.com/hdfchain/hdfd/pull/2005))

### hdfd command-line flags and configuration:

- config: Remove deprecated getworkkey option ([hdfchain/hdfd#1594](https://github.com/hdfchain/hdfd/pull/1594))

### certgen utility changes:

- certgen: Support Ed25519 cert generation on Go 1.13 ([hdfchain/hdfd#1757](https://github.com/hdfchain/hdfd/pull/1757))

### hdfctl utility changes:

- hdfctl: Make version string consistent ([hdfchain/hdfd#1598](https://github.com/hdfchain/hdfd/pull/1598))
- hdfctl: Update for dcrjson/v2 and wallet types ([hdfchain/hdfd#1609](https://github.com/hdfchain/hdfd/pull/1609))
- sampleconfig: add export hdfctl sample config ([hdfchain/hdfd#2006](https://github.com/hdfchain/hdfd/pull/2006))

### promptsecret utility changes:

- promptsecret: Add -n flag to prompt multiple times ([hdfchain/hdfd#1705](https://github.com/hdfchain/hdfd/pull/1705))

### Documentation:

- docs: Update for secp256k1 v2 module ([hdfchain/hdfd#1919](https://github.com/hdfchain/hdfd/pull/1919))
- docs: document module breaking changes process ([hdfchain/hdfd#1891](https://github.com/hdfchain/hdfd/pull/1891))
- docs: Link to btc whitepaper on clkj.ltd ([hdfchain/hdfd#1885](https://github.com/hdfchain/hdfd/pull/1885))
- docs: Update for mempool v3 module ([hdfchain/hdfd#1835](https://github.com/hdfchain/hdfd/pull/1835))
- docs: Update for peer v2 module ([hdfchain/hdfd#1834](https://github.com/hdfchain/hdfd/pull/1834))
- docs: Update for connmgr v2 module ([hdfchain/hdfd#1833](https://github.com/hdfchain/hdfd/pull/1833))
- docs: Update for mining v2 module ([hdfchain/hdfd#1831](https://github.com/hdfchain/hdfd/pull/1831))
- docs: Update for blockchain v2 module ([hdfchain/hdfd#1823](https://github.com/hdfchain/hdfd/pull/1823))
- docs: Update for rpcclient v4 module ([hdfchain/hdfd#1807](https://github.com/hdfchain/hdfd/pull/1807))
- docs: Update for blockchain/stake v2 module ([hdfchain/hdfd#1803](https://github.com/hdfchain/hdfd/pull/1803))
- docs: Update for database v2 module ([hdfchain/hdfd#1799](https://github.com/hdfchain/hdfd/pull/1799))
- docs: Update for rpcclient v3 module ([hdfchain/hdfd#1793](https://github.com/hdfchain/hdfd/pull/1793))
- docs: Update for dcrjson/v3 module ([hdfchain/hdfd#1792](https://github.com/hdfchain/hdfd/pull/1792))
- docs: Update for txscript v2 module ([hdfchain/hdfd#1774](https://github.com/hdfchain/hdfd/pull/1774))
- docs: Update for dcrutil v2 module ([hdfchain/hdfd#1770](https://github.com/hdfchain/hdfd/pull/1770))
- docs: Update for dcrec/edwards v2 module ([hdfchain/hdfd#1765](https://github.com/hdfchain/hdfd/pull/1765))
- docs: Update for chaincfg v2 module ([hdfchain/hdfd#1698](https://github.com/hdfchain/hdfd/pull/1698))
- docs: Update for hdkeychain v2 module ([hdfchain/hdfd#1696](https://github.com/hdfchain/hdfd/pull/1696))
- hdkeychain: Correct docs key examples ([hdfchain/hdfd#1696](https://github.com/hdfchain/hdfd/pull/1696))
- docs: allowHighFees arg has been implemented ([hdfchain/hdfd#1695](https://github.com/hdfchain/hdfd/pull/1695))
- docs: move json rpc docs to mediawiki ([hdfchain/hdfd#1687](https://github.com/hdfchain/hdfd/pull/1687))
- docs: Update for lru module ([hdfchain/hdfd#1683](https://github.com/hdfchain/hdfd/pull/1683))
- docs: fix formatting in json rpc doc ([hdfchain/hdfd#1633](https://github.com/hdfchain/hdfd/pull/1633))
- docs: Update for mempool v2 module ([hdfchain/hdfd#1613](https://github.com/hdfchain/hdfd/pull/1613))
- docs: Update for rpcclient v2 module ([hdfchain/hdfd#1608](https://github.com/hdfchain/hdfd/pull/1608))
- docs: Update for dcrjson v2 module ([hdfchain/hdfd#1607](https://github.com/hdfchain/hdfd/pull/1607))
- jsonrpc/types: Add README.md and doc.go ([hdfchain/hdfd#1794](https://github.com/hdfchain/hdfd/pull/1794))
- dcrjson: Update README.md ([hdfchain/hdfd#1607](https://github.com/hdfchain/hdfd/pull/1607))
- dcrec/secp256k1: Update README.md broken link ([hdfchain/hdfd#1631](https://github.com/hdfchain/hdfd/pull/1631))
- bech32: Correct README build badge reference ([hdfchain/hdfd#1689](https://github.com/hdfchain/hdfd/pull/1689))
- hdkeychain: Update README.md ([hdfchain/hdfd#1686](https://github.com/hdfchain/hdfd/pull/1686))
- bech32: Correct README links ([hdfchain/hdfd#1691](https://github.com/hdfchain/hdfd/pull/1691))
- stake: Remove unnecessary language in comment ([hdfchain/hdfd#1752](https://github.com/hdfchain/hdfd/pull/1752))
- multi: Use https links where available ([hdfchain/hdfd#1771](https://github.com/hdfchain/hdfd/pull/1771))
- stake: Make doc.go formatting consistent ([hdfchain/hdfd#1803](https://github.com/hdfchain/hdfd/pull/1803))
- blockchain: Update doc.go to reflect reality ([hdfchain/hdfd#1823](https://github.com/hdfchain/hdfd/pull/1823))
- multi: update rpc documentation ([hdfchain/hdfd#1867](https://github.com/hdfchain/hdfd/pull/1867))
- dcrec: fix examples links ([hdfchain/hdfd#1914](https://github.com/hdfchain/hdfd/pull/1914))
- gcs: Improve package documentation ([hdfchain/hdfd#1915](https://github.com/hdfchain/hdfd/pull/1915))

### Developer-related package and module changes:

- dcrutil: Return deep copied tx in NewTxDeepTxIns ([hdfchain/hdfd#1545](https://github.com/hdfchain/hdfd/pull/1545))
- mining: Remove superfluous error check ([hdfchain/hdfd#1552](https://github.com/hdfchain/hdfd/pull/1552))
- dcrutil: Block does not cache the header bytes ([hdfchain/hdfd#1571](https://github.com/hdfchain/hdfd/pull/1571))
- blockchain: Remove superfluous GetVoteInfo check ([hdfchain/hdfd#1574](https://github.com/hdfchain/hdfd/pull/1574))
- blockchain: Make consensus votes network agnostic ([hdfchain/hdfd#1590](https://github.com/hdfchain/hdfd/pull/1590))
- blockchain: Optimize skip stakebase input ([hdfchain/hdfd#1565](https://github.com/hdfchain/hdfd/pull/1565))
- txscript: code cleanup ([hdfchain/hdfd#1591](https://github.com/hdfchain/hdfd/pull/1591))
- dcrjson: Move estimate fee test to matching file ([hdfchain/hdfd#1607](https://github.com/hdfchain/hdfd/pull/1607))
- dcrjson: Move raw stake tx cmds to correct file ([hdfchain/hdfd#1607](https://github.com/hdfchain/hdfd/pull/1607))
- dcrjson: Move best block result to correct file ([hdfchain/hdfd#1607](https://github.com/hdfchain/hdfd/pull/1607))
- dcrjson: Move winning tickets ntfn to correct file ([hdfchain/hdfd#1607](https://github.com/hdfchain/hdfd/pull/1607))
- dcrjson: Move spent tickets ntfn to correct file ([hdfchain/hdfd#1607](https://github.com/hdfchain/hdfd/pull/1607))
- dcrjson: Move stake diff ntfn to correct file ([hdfchain/hdfd#1607](https://github.com/hdfchain/hdfd/pull/1607))
- dcrjson: Move new tickets ntfn to correct file ([hdfchain/hdfd#1607](https://github.com/hdfchain/hdfd/pull/1607))
- txscript: Rename p2sh indicator to isP2SH ([hdfchain/hdfd#1605](https://github.com/hdfchain/hdfd/pull/1605))
- mempool: Remove deprecated min high prio constant ([hdfchain/hdfd#1613](https://github.com/hdfchain/hdfd/pull/1613))
- mempool: Remove tight coupling with dcrjson ([hdfchain/hdfd#1613](https://github.com/hdfchain/hdfd/pull/1613))
- blockmanager: only check if current once handling inv's ([hdfchain/hdfd#1621](https://github.com/hdfchain/hdfd/pull/1621))
- connmngr: Add DialAddr config option ([hdfchain/hdfd#1642](https://github.com/hdfchain/hdfd/pull/1642))
- txscript: Consistent checksigaltverify handling ([hdfchain/hdfd#1647](https://github.com/hdfchain/hdfd/pull/1647))
- multi: preallocate memory ([hdfchain/hdfd#1646](https://github.com/hdfchain/hdfd/pull/1646))
- wire: Fix maximum payload length of MsgAddr ([hdfchain/hdfd#1638](https://github.com/hdfchain/hdfd/pull/1638))
- blockmanager: remove unused requestedEverTxns ([hdfchain/hdfd#1624](https://github.com/hdfchain/hdfd/pull/1624))
- blockmanager: remove useless requestedEverBlocks ([hdfchain/hdfd#1624](https://github.com/hdfchain/hdfd/pull/1624))
- txscript: Introduce constant for max CLTV bytes ([hdfchain/hdfd#1650](https://github.com/hdfchain/hdfd/pull/1650))
- txscript: Introduce constant for max CSV bytes ([hdfchain/hdfd#1651](https://github.com/hdfchain/hdfd/pull/1651))
- chaincfg: Remove unused definition ([hdfchain/hdfd#1661](https://github.com/hdfchain/hdfd/pull/1661))
- chaincfg: Use expected regnet merkle root var ([hdfchain/hdfd#1662](https://github.com/hdfchain/hdfd/pull/1662))
- blockchain: Deprecate BlockOneCoinbasePaysTokens ([hdfchain/hdfd#1657](https://github.com/hdfchain/hdfd/pull/1657))
- blockchain: Explicit script ver in coinbase checks ([hdfchain/hdfd#1658](https://github.com/hdfchain/hdfd/pull/1658))
- chaincfg: Explicit unique net addr prefix ([hdfchain/hdfd#1663](https://github.com/hdfchain/hdfd/pull/1663))
- chaincfg: Introduce params lookup by addr prefix ([hdfchain/hdfd#1664](https://github.com/hdfchain/hdfd/pull/1664))
- dcrutil: Lookup params by addr prefix in chaincfg ([hdfchain/hdfd#1665](https://github.com/hdfchain/hdfd/pull/1665))
- peer: Deprecate dependency on chaincfg ([hdfchain/hdfd#1671](https://github.com/hdfchain/hdfd/pull/1671))
- server: Update for deprecated peer chaincfg ([hdfchain/hdfd#1671](https://github.com/hdfchain/hdfd/pull/1671))
- fees: drop unused chaincfg ([hdfchain/hdfd#1675](https://github.com/hdfchain/hdfd/pull/1675))
- lru: Implement a new module with generic LRU cache ([hdfchain/hdfd#1683](https://github.com/hdfchain/hdfd/pull/1683))
- peer: Use lru cache module for inventory ([hdfchain/hdfd#1683](https://github.com/hdfchain/hdfd/pull/1683))
- peer: Use lru cache module for nonces ([hdfchain/hdfd#1683](https://github.com/hdfchain/hdfd/pull/1683))
- server: Use lru cache module for addresses ([hdfchain/hdfd#1683](https://github.com/hdfchain/hdfd/pull/1683))
- multi: drop init and just set default log ([hdfchain/hdfd#1676](https://github.com/hdfchain/hdfd/pull/1676))
- multi: deprecate DisableLog ([hdfchain/hdfd#1676](https://github.com/hdfchain/hdfd/pull/1676))
- blockchain: Remove unused params from block index ([hdfchain/hdfd#1674](https://github.com/hdfchain/hdfd/pull/1674))
- bech32: Initial Version ([hdfchain/hdfd#1646](https://github.com/hdfchain/hdfd/pull/1646))
- chaincfg: Add extended key accessor funcs ([hdfchain/hdfd#1694](https://github.com/hdfchain/hdfd/pull/1694))
- chaincfg: Rename extended key accessor funcs ([hdfchain/hdfd#1699](https://github.com/hdfchain/hdfd/pull/1699))
- wire: Accurate calculations of maximum length ([hdfchain/hdfd#1672](https://github.com/hdfchain/hdfd/pull/1672))
- wire: Fix MsgCFTypes maximum payload length ([hdfchain/hdfd#1673](https://github.com/hdfchain/hdfd/pull/1673))
- txscript: Deprecate HasP2SHScriptSigStakeOpCodes ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Deprecate IsStakeOutput ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Deprecate GetMultisigMandN ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Introduce zero-alloc script tokenizer ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize script disasm ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Introduce raw script sighash calc func ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize CalcSignatureHash ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Make isSmallInt accept raw opcode ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Make asSmallInt accept raw opcode ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Make isStakeOpcode accept raw opcode ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize IsPayToScriptHash ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize IsMultisigScript ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize IsMultisigSigScript ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize GetSigOpCount ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize isAnyKindOfScriptHash ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize IsPushOnlyScript ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize new engine push only script ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Check p2sh push before parsing scripts ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize GetPreciseSigOpCount ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Make typeOfScript accept raw script ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize typeOfScript pay-to-script-hash ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Remove unused isScriptHash function ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize typeOfScript multisig ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Remove unused isMultiSig function ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize typeOfScript pay-to-pubkey ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Remove unused isPubkey function ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize typeOfScript pay-to-alt-pubkey ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize typeOfScript pay-to-pubkey-hash ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Remove unused isPubkeyHash function ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize typeOfScript pay-to-alt-pk-hash ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize typeOfScript nulldata detection ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Remove unused isNullData function ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize typeOfScript stakesub detection ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Remove unused isStakeSubmission function ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize typeOfScript stakegen detection ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Remove unused isStakeGen function ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize typeOfScript stakerev detection ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Remove unused isStakeRevocation function ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize typeOfScript stakechange detect ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Remove unused isSStxChange function ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize ContainsStakeOpCodes ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize ExtractCoinbaseNullData ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Convert CalcScriptInfo ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Remove unused isPushOnly function ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Remove unused getSigOpCount function ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize CalcMultiSigStats ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize multi sig redeem script func ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Convert GetScriptHashFromP2SHScript ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize PushedData ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize IsUnspendable ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Make canonicalPush accept raw opcode ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize ExtractAtomicSwapDataPushes ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize ExtractPkScriptAddrs scripthash ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize ExtractPkScriptAddrs pubkeyhash ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize ExtractPkScriptAddrs altpubkeyhash ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize ExtractPkScriptAddrs pubkey ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize ExtractPkScriptAddrs altpubkey ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize ExtractPkScriptAddrs multisig ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize ExtractPkScriptAddrs stakesub ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize ExtractPkScriptAddrs stakegen ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize ExtractPkScriptAddrs stakerev ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize ExtractPkScriptAddrs stakechange ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize ExtractPkScriptAddrs nulldata ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Optimize ExtractPkScriptAltSigType ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Remove unused extractOneBytePush func ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Remove unused isPubkeyAlt function ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Remove unused isPubkeyHashAlt function ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Remove unused isOneByteMaxDataPush func ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: mergeMultiSig function def order cleanup ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Use raw scripts in RawTxInSignature ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Use raw scripts in RawTxInSignatureAlt ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Correct p2pkSignatureScriptAlt comment ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Use raw scripts in SignTxOutput ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Implement efficient opcode data removal ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Make isDisabled accept raw opcode ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Make alwaysIllegal accept raw opcode ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Make isConditional accept raw opcode ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Make min push accept raw opcode and data ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Convert to use non-parsed opcode disasm ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Refactor engine to use raw scripts ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Remove unused removeOpcodeByData func ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Rename removeOpcodeByDataRaw func ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Remove unused calcSignatureHash func ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Rename calcSignatureHashRaw func ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Remove unused parseScript func ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Remove unused unparseScript func ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Remove unused parsedOpcode.bytes func ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Remove unused parseScriptTemplate func ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Make executeOpcode take opcode and data ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Make op callbacks take opcode and data ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- dcrutil: Fix NewTxDeepTxIns implementation ([hdfchain/hdfd#1685](https://github.com/hdfchain/hdfd/pull/1685))
- stake: drop txscript.DefaultScriptVersion usage ([hdfchain/hdfd#1704](https://github.com/hdfchain/hdfd/pull/1704))
- peer: invSendQueue is a FIFO ([hdfchain/hdfd#1680](https://github.com/hdfchain/hdfd/pull/1680))
- peer: pendingMsgs is a FIFO ([hdfchain/hdfd#1680](https://github.com/hdfchain/hdfd/pull/1680))
- blockchain: drop container/list ([hdfchain/hdfd#1682](https://github.com/hdfchain/hdfd/pull/1682))
- blockmanager: use local var for the request queue ([hdfchain/hdfd#1622](https://github.com/hdfchain/hdfd/pull/1622))
- server: return on outbound peer creation error ([hdfchain/hdfd#1637](https://github.com/hdfchain/hdfd/pull/1637))
- hdkeychain: Remove Address method ([hdfchain/hdfd#1696](https://github.com/hdfchain/hdfd/pull/1696))
- hdkeychain: Remove SetNet method ([hdfchain/hdfd#1696](https://github.com/hdfchain/hdfd/pull/1696))
- hdkeychain: Require network on decode extended key ([hdfchain/hdfd#1696](https://github.com/hdfchain/hdfd/pull/1696))
- hdkeychain: Don't rely on global state ([hdfchain/hdfd#1696](https://github.com/hdfchain/hdfd/pull/1696))
- hdkeychain: Introduce NetworkParams interface ([hdfchain/hdfd#1696](https://github.com/hdfchain/hdfd/pull/1696))
- server: Remove unused ScheduleShutdown func ([hdfchain/hdfd#1711](https://github.com/hdfchain/hdfd/pull/1711))
- server: Remove unused dynamicTickDuration func ([hdfchain/hdfd#1711](https://github.com/hdfchain/hdfd/pull/1711))
- main: Convert signal handling to use context ([hdfchain/hdfd#1712](https://github.com/hdfchain/hdfd/pull/1712))
- txscript: Remove checks for impossible conditions ([hdfchain/hdfd#1713](https://github.com/hdfchain/hdfd/pull/1713))
- indexers: Remove unused func ([hdfchain/hdfd#1714](https://github.com/hdfchain/hdfd/pull/1714))
- multi: fix onVoteReceivedHandler shutdown ([hdfchain/hdfd#1721](https://github.com/hdfchain/hdfd/pull/1721))
- wire: Rename extended errors to malformed errors ([hdfchain/hdfd#1742](https://github.com/hdfchain/hdfd/pull/1742))
- rpcwebsocket: convert from list to simple FIFO ([hdfchain/hdfd#1726](https://github.com/hdfchain/hdfd/pull/1726))
- dcrec: implement GenerateKey ([hdfchain/hdfd#1652](https://github.com/hdfchain/hdfd/pull/1652))
- txscript: Remove SigHashOptimization constant ([hdfchain/hdfd#1698](https://github.com/hdfchain/hdfd/pull/1698))
- txscript: Remove CheckForDuplicateHashes constant ([hdfchain/hdfd#1698](https://github.com/hdfchain/hdfd/pull/1698))
- txscript: Remove CPUMinerThreads constant ([hdfchain/hdfd#1698](https://github.com/hdfchain/hdfd/pull/1698))
- chaincfg: Move DNSSeed stringer next to type def ([hdfchain/hdfd#1698](https://github.com/hdfchain/hdfd/pull/1698))
- chaincfg: Remove all registration capabilities ([hdfchain/hdfd#1698](https://github.com/hdfchain/hdfd/pull/1698))
- chaincfg: Move mainnet code to mainnet files ([hdfchain/hdfd#1698](https://github.com/hdfchain/hdfd/pull/1698))
- chaincfg: Move testnet3 code to testnet files ([hdfchain/hdfd#1698](https://github.com/hdfchain/hdfd/pull/1698))
- chaincfg: Move simnet code to testnet files ([hdfchain/hdfd#1698](https://github.com/hdfchain/hdfd/pull/1698))
- chaincfg: Move regnet code to regnet files ([hdfchain/hdfd#1698](https://github.com/hdfchain/hdfd/pull/1698))
- chaincfg: Concrete genesis hash in Params struct ([hdfchain/hdfd#1698](https://github.com/hdfchain/hdfd/pull/1698))
- chaincfg: Use scripts in block one token payouts ([hdfchain/hdfd#1698](https://github.com/hdfchain/hdfd/pull/1698))
- chaincfg: Convert global param defs to funcs ([hdfchain/hdfd#1698](https://github.com/hdfchain/hdfd/pull/1698))
- edwards: remove curve param ([hdfchain/hdfd#1762](https://github.com/hdfchain/hdfd/pull/1762))
- edwards: unexport EncodedBytesToBigIntPoint ([hdfchain/hdfd#1762](https://github.com/hdfchain/hdfd/pull/1762))
- edwards: unexport a slew of funcs ([hdfchain/hdfd#1762](https://github.com/hdfchain/hdfd/pull/1762))
- edwards: add signature IsEqual and Verify methods ([hdfchain/hdfd#1762](https://github.com/hdfchain/hdfd/pull/1762))
- edwards: add Sign method to PrivateKey ([hdfchain/hdfd#1762](https://github.com/hdfchain/hdfd/pull/1762))
- chaincfg: Add addr params accessor funcs ([hdfchain/hdfd#1766](https://github.com/hdfchain/hdfd/pull/1766))
- schnorr: remove curve param ([hdfchain/hdfd#1764](https://github.com/hdfchain/hdfd/pull/1764))
- schnorr: unexport functions ([hdfchain/hdfd#1764](https://github.com/hdfchain/hdfd/pull/1764))
- schnorr: add signature IsEqual and Verify methods ([hdfchain/hdfd#1764](https://github.com/hdfchain/hdfd/pull/1764))
- secp256k1: unexport NAF ([hdfchain/hdfd#1764](https://github.com/hdfchain/hdfd/pull/1764))
- addrmgr: drop container/list ([hdfchain/hdfd#1679](https://github.com/hdfchain/hdfd/pull/1679))
- dcrutil: Remove unused ErrAddressCollision ([hdfchain/hdfd#1767](https://github.com/hdfchain/hdfd/pull/1767))
- dcurtil: Remove unused ErrMissingDefaultNet ([hdfchain/hdfd#1767](https://github.com/hdfchain/hdfd/pull/1767))
- dcrutil: Require network on address decode ([hdfchain/hdfd#1767](https://github.com/hdfchain/hdfd/pull/1767))
- dcrutil: Remove IsForNet from Address interface ([hdfchain/hdfd#1767](https://github.com/hdfchain/hdfd/pull/1767))
- dcrutil: Remove DSA from Address interface ([hdfchain/hdfd#1767](https://github.com/hdfchain/hdfd/pull/1767))
- dcrutil: Remove Net from Address interface ([hdfchain/hdfd#1767](https://github.com/hdfchain/hdfd/pull/1767))
- dcrutil: Rename EncodeAddress to Address ([hdfchain/hdfd#1767](https://github.com/hdfchain/hdfd/pull/1767))
- dcrutil: Don't store net ref in addr impls ([hdfchain/hdfd#1767](https://github.com/hdfchain/hdfd/pull/1767))
- dcrutil: Require network on WIF decode ([hdfchain/hdfd#1767](https://github.com/hdfchain/hdfd/pull/1767))
- dcrutil: Accept magic bytes directly in NewWIF ([hdfchain/hdfd#1767](https://github.com/hdfchain/hdfd/pull/1767))
- dcrutil: Introduce AddressParams interface ([hdfchain/hdfd#1767](https://github.com/hdfchain/hdfd/pull/1767))
- blockchain: Do coinbase nulldata check locally ([hdfchain/hdfd#1770](https://github.com/hdfchain/hdfd/pull/1770))
- blockchain: update CalcBlockSubsidy ([hdfchain/hdfd#1750](https://github.com/hdfchain/hdfd/pull/1750))
- txscript: Use const for sighashall optimization ([hdfchain/hdfd#1774](https://github.com/hdfchain/hdfd/pull/1774))
- txscript: Remove DisableLog ([hdfchain/hdfd#1774](https://github.com/hdfchain/hdfd/pull/1774))
- txscript: Unexport HasP2SHScriptSigStakeOpCodes ([hdfchain/hdfd#1774](https://github.com/hdfchain/hdfd/pull/1774))
- txscript: Remove third GetPreciseSigOpCount param ([hdfchain/hdfd#1774](https://github.com/hdfchain/hdfd/pull/1774))
- txscript: Remove IsMultisigScript err return ([hdfchain/hdfd#1774](https://github.com/hdfchain/hdfd/pull/1774))
- txscript: Unexport IsStakeOutput ([hdfchain/hdfd#1774](https://github.com/hdfchain/hdfd/pull/1774))
- txscript: Remove CalcScriptInfo ([hdfchain/hdfd#1774](https://github.com/hdfchain/hdfd/pull/1774))
- txscript: Remove multisig redeem script err return ([hdfchain/hdfd#1774](https://github.com/hdfchain/hdfd/pull/1774))
- txscript: Remove GetScriptHashFromP2SHScript ([hdfchain/hdfd#1774](https://github.com/hdfchain/hdfd/pull/1774))
- txscript: Remove GetMultisigMandN ([hdfchain/hdfd#1774](https://github.com/hdfchain/hdfd/pull/1774))
- txscript: Remove DefaultScriptVersion ([hdfchain/hdfd#1774](https://github.com/hdfchain/hdfd/pull/1774))
- txscript: Use secp256k1 types in sig cache ([hdfchain/hdfd#1774](https://github.com/hdfchain/hdfd/pull/1774))
- multi: decouple BlockManager from server ([hdfchain/hdfd#1728](https://github.com/hdfchain/hdfd/pull/1728))
- database: Introduce BlockSerializer interface ([hdfchain/hdfd#1799](https://github.com/hdfchain/hdfd/pull/1799))
- hdkeychain: Add ChildNum and Depth methods ([hdfchain/hdfd#1800](https://github.com/hdfchain/hdfd/pull/1800))
- chaincfg: Avoid block 1 subsidy codegen explosion ([hdfchain/hdfd#1801](https://github.com/hdfchain/hdfd/pull/1801))
- chaincfg: Add stake params accessor funcs ([hdfchain/hdfd#1802](https://github.com/hdfchain/hdfd/pull/1802))
- stake: Remove DisableLog ([hdfchain/hdfd#1803](https://github.com/hdfchain/hdfd/pull/1803))
- stake: Remove unused TxSSGenStakeOutputInfo ([hdfchain/hdfd#1803](https://github.com/hdfchain/hdfd/pull/1803))
- stake: Remove unused TxSSRtxStakeOutputInfo ([hdfchain/hdfd#1803](https://github.com/hdfchain/hdfd/pull/1803))
- stake: Remove unused SetTxTree ([hdfchain/hdfd#1803](https://github.com/hdfchain/hdfd/pull/1803))
- stake: Introduce StakeParams interface ([hdfchain/hdfd#1803](https://github.com/hdfchain/hdfd/pull/1803))
- stake: Accept AddressParams for ticket commit addr ([hdfchain/hdfd#1803](https://github.com/hdfchain/hdfd/pull/1803))
- gcs: Optimize AddSigScript ([hdfchain/hdfd#1804](https://github.com/hdfchain/hdfd/pull/1804))
- chaincfg: Add subsidy params accessor funcs ([hdfchain/hdfd#1813](https://github.com/hdfchain/hdfd/pull/1813))
- blockchain/standalone: Implement a new module ([hdfchain/hdfd#1808](https://github.com/hdfchain/hdfd/pull/1808))
- blockchain/standalone: Add merkle root calc funcs ([hdfchain/hdfd#1809](https://github.com/hdfchain/hdfd/pull/1809))
- blockchain/standalone: Add subsidy calc funcs ([hdfchain/hdfd#1812](https://github.com/hdfchain/hdfd/pull/1812))
- blockchain/standalone: Add IsCoinBaseTx ([hdfchain/hdfd#1815](https://github.com/hdfchain/hdfd/pull/1815))
- crypto/blake256: Add module with zero alloc funcs ([hdfchain/hdfd#1811](https://github.com/hdfchain/hdfd/pull/1811))
- stake: Check minimum req outputs for votes earlier ([hdfchain/hdfd#1819](https://github.com/hdfchain/hdfd/pull/1819))
- blockchain: Use standalone module for merkle calcs ([hdfchain/hdfd#1816](https://github.com/hdfchain/hdfd/pull/1816))
- blockchain: Use standalone for coinbase checks ([hdfchain/hdfd#1816](https://github.com/hdfchain/hdfd/pull/1816))
- blockchain: Use standalone module subsidy calcs ([hdfchain/hdfd#1816](https://github.com/hdfchain/hdfd/pull/1816))
- blockchain: Use standalone module for work funcs ([hdfchain/hdfd#1816](https://github.com/hdfchain/hdfd/pull/1816))
- blockchain: Remove deprecated code ([hdfchain/hdfd#1823](https://github.com/hdfchain/hdfd/pull/1823))
- blockchain: Accept subsidy cache in config ([hdfchain/hdfd#1823](https://github.com/hdfchain/hdfd/pull/1823))
- mining: Use lastest major version deps ([hdfchain/hdfd#1831](https://github.com/hdfchain/hdfd/pull/1831))
- connmgr: Accept DNS seeds as string slice ([hdfchain/hdfd#1833](https://github.com/hdfchain/hdfd/pull/1833))
- peer: Remove deprecated Config.ChainParams field ([hdfchain/hdfd#1834](https://github.com/hdfchain/hdfd/pull/1834))
- peer: Accept hash slice for block locators ([hdfchain/hdfd#1834](https://github.com/hdfchain/hdfd/pull/1834))
- peer: Use latest major version deps ([hdfchain/hdfd#1834](https://github.com/hdfchain/hdfd/pull/1834))
- mempool: Use latest major version deps ([hdfchain/hdfd#1835](https://github.com/hdfchain/hdfd/pull/1835))
- main: Update to use all new major module versions ([hdfchain/hdfd#1837](https://github.com/hdfchain/hdfd/pull/1837))
- blockchain: Implement stricter bounds checking ([hdfchain/hdfd#1825](https://github.com/hdfchain/hdfd/pull/1825))
- gcs: Start v2 module dev cycle ([hdfchain/hdfd#1843](https://github.com/hdfchain/hdfd/pull/1843))
- gcs: Support empty filters ([hdfchain/hdfd#1844](https://github.com/hdfchain/hdfd/pull/1844))
- gcs: Make error consistent with rest of codebase ([hdfchain/hdfd#1846](https://github.com/hdfchain/hdfd/pull/1846))
- gcs: Add filter version support ([hdfchain/hdfd#1848](https://github.com/hdfchain/hdfd/pull/1848))
- gcs: Correct zero hash filter matches ([hdfchain/hdfd#1857](https://github.com/hdfchain/hdfd/pull/1857))
- gcs: Standardize serialization on a single format ([hdfchain/hdfd#1851](https://github.com/hdfchain/hdfd/pull/1851))
- gcs: Optimize Hash ([hdfchain/hdfd#1853](https://github.com/hdfchain/hdfd/pull/1853))
- gcs: Group V1 filter funcs after filter defs ([hdfchain/hdfd#1854](https://github.com/hdfchain/hdfd/pull/1854))
- gcs: Support independent fp rate and bin size ([hdfchain/hdfd#1854](https://github.com/hdfchain/hdfd/pull/1854))
- blockchain: Refactor best chain state init ([hdfchain/hdfd#1871](https://github.com/hdfchain/hdfd/pull/1871))
- gcs: Implement version 2 filters ([hdfchain/hdfd#1856](https://github.com/hdfchain/hdfd/pull/1856))
- blockchain: Cleanup subsidy cache init order ([hdfchain/hdfd#1873](https://github.com/hdfchain/hdfd/pull/1873))
- multi: use chain ref. from blockmanager config ([hdfchain/hdfd#1879](https://github.com/hdfchain/hdfd/pull/1879))
- multi: remove unused funcs and vars ([hdfchain/hdfd#1880](https://github.com/hdfchain/hdfd/pull/1880))
- gcs: Prevent empty data elements in v2 filters ([hdfchain/hdfd#1911](https://github.com/hdfchain/hdfd/pull/1911))
- crypto: import ripemd160 ([hdfchain/hdfd#1907](https://github.com/hdfchain/hdfd/pull/1907))
- multi: Use secp256k1/v2 module ([hdfchain/hdfd#1919](https://github.com/hdfchain/hdfd/pull/1919))
- multi: Use crypto/ripemd160 module ([hdfchain/hdfd#1918](https://github.com/hdfchain/hdfd/pull/1918))
- multi: Use dcrec/edwards/v2 module ([hdfchain/hdfd#1920](https://github.com/hdfchain/hdfd/pull/1920))
- gcs: Prevent empty data elements fp matches ([hdfchain/hdfd#1940](https://github.com/hdfchain/hdfd/pull/1940))
- main: Update to use all new module versions ([hdfchain/hdfd#1946](https://github.com/hdfchain/hdfd/pull/1946))
- blockchain/standalone: Add inclusion proof funcs ([hdfchain/hdfd#1906](https://github.com/hdfchain/hdfd/pull/1906))

### Developer-related module management:

- build: Require dcrjson v1.2.0 ([hdfchain/hdfd#1607](https://github.com/hdfchain/hdfd/pull/1607))
- multi: Remove non-root module replacements ([hdfchain/hdfd#1599](https://github.com/hdfchain/hdfd/pull/1599))
- dcrjson: Introduce v2 module without wallet types ([hdfchain/hdfd#1607](https://github.com/hdfchain/hdfd/pull/1607))
- release: Freeze version 1 mempool module ([hdfchain/hdfd#1613](https://github.com/hdfchain/hdfd/pull/1613))
- release: Introduce mempool v2 module ([hdfchain/hdfd#1613](https://github.com/hdfchain/hdfd/pull/1613))
- main: Tidy module to latest ([hdfchain/hdfd#1613](https://github.com/hdfchain/hdfd/pull/1613))
- main: Update for mempool/v2 ([hdfchain/hdfd#1616](https://github.com/hdfchain/hdfd/pull/1616))
- multi: Add go 1.11 directive to all modules ([hdfchain/hdfd#1677](https://github.com/hdfchain/hdfd/pull/1677))
- build: Tidy module sums (go mod tidy) ([hdfchain/hdfd#1692](https://github.com/hdfchain/hdfd/pull/1692))
- release: Freeze version 1 hdkeychain module ([hdfchain/hdfd#1696](https://github.com/hdfchain/hdfd/pull/1696))
- release: Introduce hdkeychain v2 module ([hdfchain/hdfd#1696](https://github.com/hdfchain/hdfd/pull/1696))
- release: Freeze version 1 chaincfg module ([hdfchain/hdfd#1698](https://github.com/hdfchain/hdfd/pull/1698))
- chaincfg: Introduce chaincfg v2 module ([hdfchain/hdfd#1698](https://github.com/hdfchain/hdfd/pull/1698))
- chaincfg: Use dcrec/edwards/v1.0.0 ([hdfchain/hdfd#1758](https://github.com/hdfchain/hdfd/pull/1758))
- dcrutil: Prepare v1.3.0 ([hdfchain/hdfd#1761](https://github.com/hdfchain/hdfd/pull/1761))
- release: freeze version 1 dcrec/edwards module ([hdfchain/hdfd#1762](https://github.com/hdfchain/hdfd/pull/1762))
- edwards: Introduce v2 module ([hdfchain/hdfd#1762](https://github.com/hdfchain/hdfd/pull/1762))
- release: freeze version 1 dcrec/secp256k1 module ([hdfchain/hdfd#1764](https://github.com/hdfchain/hdfd/pull/1764))
- secp256k1: Introduce v2 module ([hdfchain/hdfd#1764](https://github.com/hdfchain/hdfd/pull/1764))
- multi: Update all modules for chaincfg v1.5.1 ([hdfchain/hdfd#1768](https://github.com/hdfchain/hdfd/pull/1768))
- release: Freeze version 1 dcrutil module ([hdfchain/hdfd#1767](https://github.com/hdfchain/hdfd/pull/1767))
- dcrutil: Update to use chaincfg/v2 module ([hdfchain/hdfd#1767](https://github.com/hdfchain/hdfd/pull/1767))
- release: Introduce dcrutil v2 module ([hdfchain/hdfd#1767](https://github.com/hdfchain/hdfd/pull/1767))
- database: Use chaincfg/v2 ([hdfchain/hdfd#1772](https://github.com/hdfchain/hdfd/pull/1772))
- txscript: Prepare v1.1.0 ([hdfchain/hdfd#1773](https://github.com/hdfchain/hdfd/pull/1773))
- stake: Prepare v1.2.0 ([hdfchain/hdfd#1775](https://github.com/hdfchain/hdfd/pull/1775))
- release: Freeze version 1 txscript module ([hdfchain/hdfd#1774](https://github.com/hdfchain/hdfd/pull/1774))
- txscript: Use dcrutil/v2 ([hdfchain/hdfd#1774](https://github.com/hdfchain/hdfd/pull/1774))
- release: Introduce txscript v2 module ([hdfchain/hdfd#1774](https://github.com/hdfchain/hdfd/pull/1774))
- main: Add requires for new version modules ([hdfchain/hdfd#1776](https://github.com/hdfchain/hdfd/pull/1776))
- dcrjson: Introduce v3 and move types to module ([hdfchain/hdfd#1779](https://github.com/hdfchain/hdfd/pull/1779))
- jsonrpc/types: Prepare 1.0.0 ([hdfchain/hdfd#1787](https://github.com/hdfchain/hdfd/pull/1787))
- main: Use latest JSON-RPC types ([hdfchain/hdfd#1789](https://github.com/hdfchain/hdfd/pull/1789))
- multi: Use hdfchain fork of go-socks ([hdfchain/hdfd#1790](https://github.com/hdfchain/hdfd/pull/1790))
- rpcclient: Prepare v2.1.0 ([hdfchain/hdfd#1791](https://github.com/hdfchain/hdfd/pull/1791))
- release: Freeze version 2 rpcclient module ([hdfchain/hdfd#1793](https://github.com/hdfchain/hdfd/pull/1793))
- rpcclient: Use dcrjson/v3 ([hdfchain/hdfd#1793](https://github.com/hdfchain/hdfd/pull/1793))
- release: Introduce rpcclient v3 module ([hdfchain/hdfd#1793](https://github.com/hdfchain/hdfd/pull/1793))
- main: Use rpcclient/v3 ([hdfchain/hdfd#1795](https://github.com/hdfchain/hdfd/pull/1795))
- hdkeychain: Prepare v2.0.1 ([hdfchain/hdfd#1798](https://github.com/hdfchain/hdfd/pull/1798))
- release: Freeze version 1 database module ([hdfchain/hdfd#1799](https://github.com/hdfchain/hdfd/pull/1799))
- database: Use dcrutil/v2 ([hdfchain/hdfd#1799](https://github.com/hdfchain/hdfd/pull/1799))
- release: Introduce database v2 module ([hdfchain/hdfd#1799](https://github.com/hdfchain/hdfd/pull/1799))
- release: Freeze version 1 blockchain/stake module ([hdfchain/hdfd#1803](https://github.com/hdfchain/hdfd/pull/1803))
- stake: Use dcrutil/v2 and chaincfg/v2 ([hdfchain/hdfd#1803](https://github.com/hdfchain/hdfd/pull/1803))
- Use txscript/v2 ([hdfchain/hdfd#1803](https://github.com/hdfchain/hdfd/pull/1803))
- stake: Use database/v2 ([hdfchain/hdfd#1803](https://github.com/hdfchain/hdfd/pull/1803))
- release: Introduce blockchain/stake v2 module ([hdfchain/hdfd#1803](https://github.com/hdfchain/hdfd/pull/1803))
- gcs: Use txscript/v2 ([hdfchain/hdfd#1804](https://github.com/hdfchain/hdfd/pull/1804))
- gcs: Prepare v1.1.0 ([hdfchain/hdfd#1804](https://github.com/hdfchain/hdfd/pull/1804))
- release: Freeze version 3 rpcclient module ([hdfchain/hdfd#1807](https://github.com/hdfchain/hdfd/pull/1807))
- rpcclient: Use dcrutil/v2 and chaincfg/v2 ([hdfchain/hdfd#1807](https://github.com/hdfchain/hdfd/pull/1807))
- release: Introduce rpcclient v4 module ([hdfchain/hdfd#1807](https://github.com/hdfchain/hdfd/pull/1807))
- blockchain/standalone: Prepare v1.0.0 ([hdfchain/hdfd#1817](https://github.com/hdfchain/hdfd/pull/1817))
- multi: Use crypto/blake256 ([hdfchain/hdfd#1818](https://github.com/hdfchain/hdfd/pull/1818))
- main: Consume latest module minors and patches ([hdfchain/hdfd#1822](https://github.com/hdfchain/hdfd/pull/1822))
- blockchain: Prepare v1.2.0 ([hdfchain/hdfd#1820](https://github.com/hdfchain/hdfd/pull/1820))
- mining: Prepare v1.1.1 ([hdfchain/hdfd#1826](https://github.com/hdfchain/hdfd/pull/1826))
- release: Freeze version 1 blockchain module use ([hdfchain/hdfd#1823](https://github.com/hdfchain/hdfd/pull/1823))
- blockchain: Use lastest major version deps ([hdfchain/hdfd#1823](https://github.com/hdfchain/hdfd/pull/1823))
- release: Introduce blockchain v2 module ([hdfchain/hdfd#1823](https://github.com/hdfchain/hdfd/pull/1823))
- connmgr: Prepare v1.1.0 ([hdfchain/hdfd#1828](https://github.com/hdfchain/hdfd/pull/1828))
- peer: Prepare v1.2.0 ([hdfchain/hdfd#1830](https://github.com/hdfchain/hdfd/pull/1830))
- release: Freeze version 1 mining module use ([hdfchain/hdfd#1831](https://github.com/hdfchain/hdfd/pull/1831))
- release: Introduce mining v2 module ([hdfchain/hdfd#1831](https://github.com/hdfchain/hdfd/pull/1831))
- mempool: Prepare v2.1.0 ([hdfchain/hdfd#1832](https://github.com/hdfchain/hdfd/pull/1832))
- release: Freeze version 1 connmgr module use ([hdfchain/hdfd#1833](https://github.com/hdfchain/hdfd/pull/1833))
- release: Introduce connmgr v2 module ([hdfchain/hdfd#1833](https://github.com/hdfchain/hdfd/pull/1833))
- release: Freeze version 1 peer module use ([hdfchain/hdfd#1834](https://github.com/hdfchain/hdfd/pull/1834))
- release: Introduce peer v2 module ([hdfchain/hdfd#1834](https://github.com/hdfchain/hdfd/pull/1834))
- blockchain: Prepare v2.0.1 ([hdfchain/hdfd#1836](https://github.com/hdfchain/hdfd/pull/1836))
- release: Freeze version 2 mempool module use ([hdfchain/hdfd#1835](https://github.com/hdfchain/hdfd/pull/1835))
- release: Introduce mempool v3 module ([hdfchain/hdfd#1835](https://github.com/hdfchain/hdfd/pull/1835))
- go.mod: sync ([hdfchain/hdfd#1913](https://github.com/hdfchain/hdfd/pull/1913))
- secp256k1: Prepare v2.0.0 ([hdfchain/hdfd#1916](https://github.com/hdfchain/hdfd/pull/1916))
- wire: Prepare v1.3.0 ([hdfchain/hdfd#1925](https://github.com/hdfchain/hdfd/pull/1925))
- chaincfg: Prepare v2.3.0 ([hdfchain/hdfd#1926](https://github.com/hdfchain/hdfd/pull/1926))
- dcrjson: Prepare v3.0.1 ([hdfchain/hdfd#1927](https://github.com/hdfchain/hdfd/pull/1927))
- rpc/jsonrpc/types: Prepare v2.0.0 ([hdfchain/hdfd#1928](https://github.com/hdfchain/hdfd/pull/1928))
- dcrutil: Prepare v2.0.1 ([hdfchain/hdfd#1929](https://github.com/hdfchain/hdfd/pull/1929))
- blockchain/standalone: Prepare v1.1.0 ([hdfchain/hdfd#1930](https://github.com/hdfchain/hdfd/pull/1930))
- txscript: Prepare v2.1.0 ([hdfchain/hdfd#1931](https://github.com/hdfchain/hdfd/pull/1931))
- database: Prepare v2.0.1 ([hdfchain/hdfd#1932](https://github.com/hdfchain/hdfd/pull/1932))
- blockchain/stake: Prepare v2.0.2 ([hdfchain/hdfd#1933](https://github.com/hdfchain/hdfd/pull/1933))
- gcs: Prepare v2.0.0 ([hdfchain/hdfd#1934](https://github.com/hdfchain/hdfd/pull/1934))
- blockchain: Prepare v2.1.0 ([hdfchain/hdfd#1935](https://github.com/hdfchain/hdfd/pull/1935))
- addrmgr: Prepare v1.1.0 ([hdfchain/hdfd#1936](https://github.com/hdfchain/hdfd/pull/1936))
- connmgr: Prepare v2.1.0 ([hdfchain/hdfd#1937](https://github.com/hdfchain/hdfd/pull/1937))
- hdkeychain: Prepare v2.1.0 ([hdfchain/hdfd#1938](https://github.com/hdfchain/hdfd/pull/1938))
- peer: Prepare v2.1.0 ([hdfchain/hdfd#1939](https://github.com/hdfchain/hdfd/pull/1939))
- fees: Prepare v2.0.0 ([hdfchain/hdfd#1941](https://github.com/hdfchain/hdfd/pull/1941))
- rpcclient: Prepare v4.1.0 ([hdfchain/hdfd#1943](https://github.com/hdfchain/hdfd/pull/1943))
- mining: Prepare v2.0.1 ([hdfchain/hdfd#1944](https://github.com/hdfchain/hdfd/pull/1944))
- mempool: Prepare v3.1.0 ([hdfchain/hdfd#1945](https://github.com/hdfchain/hdfd/pull/1945))

### Testing and Quality Assurance:

- mempool: Accept test mungers for vote tx ([hdfchain/hdfd#1595](https://github.com/hdfchain/hdfd/pull/1595))
- build: Replace TravisCI with CI via Github actions ([hdfchain/hdfd#1903](https://github.com/hdfchain/hdfd/pull/1903))
- build: Setup github actions for CI ([hdfchain/hdfd#1902](https://github.com/hdfchain/hdfd/pull/1902))
- TravisCI: Recommended install for golangci-lint ([hdfchain/hdfd#1808](https://github.com/hdfchain/hdfd/pull/1808))
- TravisCI: Use more portable module ver stripping ([hdfchain/hdfd#1784](https://github.com/hdfchain/hdfd/pull/1784))
- TravisCI: Test and lint latest version modules ([hdfchain/hdfd#1776](https://github.com/hdfchain/hdfd/pull/1776))
- TravisCI: Disable race detector ([hdfchain/hdfd#1749](https://github.com/hdfchain/hdfd/pull/1749))
- TravisCI: Set ./run_tests.sh executable perms ([hdfchain/hdfd#1648](https://github.com/hdfchain/hdfd/pull/1648))
- travis: bump golangci-lint to v1.18.0 ([hdfchain/hdfd#1890](https://github.com/hdfchain/hdfd/pull/1890))
- travis: Test go1.13 and drop go1.11 ([hdfchain/hdfd#1875](https://github.com/hdfchain/hdfd/pull/1875))
- travis: Allow staged builds with build cache ([hdfchain/hdfd#1797](https://github.com/hdfchain/hdfd/pull/1797))
- travis: drop docker and test directly ([hdfchain/hdfd#1783](https://github.com/hdfchain/hdfd/pull/1783))
- travis: test go1.12 ([hdfchain/hdfd#1627](https://github.com/hdfchain/hdfd/pull/1627))
- travis: Add misspell linter ([hdfchain/hdfd#1618](https://github.com/hdfchain/hdfd/pull/1618))
- travis: run linters in each module ([hdfchain/hdfd#1601](https://github.com/hdfchain/hdfd/pull/1601))
- multi: switch to golangci-lint ([hdfchain/hdfd#1575](https://github.com/hdfchain/hdfd/pull/1575))
- blockchain: Consistent legacy seq lock tests ([hdfchain/hdfd#1580](https://github.com/hdfchain/hdfd/pull/1580))
- blockchain: Add test logic to find deployments ([hdfchain/hdfd#1581](https://github.com/hdfchain/hdfd/pull/1581))
- blockchain: Introduce chaingen test harness ([hdfchain/hdfd#1583](https://github.com/hdfchain/hdfd/pull/1583))
- blockchain: Use harness in force head reorg tests ([hdfchain/hdfd#1584](https://github.com/hdfchain/hdfd/pull/1584))
- blockchain: Use harness in stake version tests ([hdfchain/hdfd#1585](https://github.com/hdfchain/hdfd/pull/1585))
- blockchain: Use harness in checkblktemplate tests ([hdfchain/hdfd#1586](https://github.com/hdfchain/hdfd/pull/1586))
- blockchain: Use harness in threshold state tests ([hdfchain/hdfd#1587](https://github.com/hdfchain/hdfd/pull/1587))
- blockchain: Use harness in legacy seqlock tests ([hdfchain/hdfd#1588](https://github.com/hdfchain/hdfd/pull/1588))
- blockchain: Use harness in fixed seqlock tests ([hdfchain/hdfd#1589](https://github.com/hdfchain/hdfd/pull/1589))
- multi: cleanup linter warnings ([hdfchain/hdfd#1601](https://github.com/hdfchain/hdfd/pull/1601))
- txscript: Add remove signature reference test ([hdfchain/hdfd#1604](https://github.com/hdfchain/hdfd/pull/1604))
- rpctest: Update for rpccclient/v2 and dcrjson/v2 ([hdfchain/hdfd#1610](https://github.com/hdfchain/hdfd/pull/1610))
- wire: Add tests for MsgCFTypes ([hdfchain/hdfd#1619](https://github.com/hdfchain/hdfd/pull/1619))
- chaincfg: Move a test to chainhash package ([hdfchain/hdfd#1632](https://github.com/hdfchain/hdfd/pull/1632))
- rpctest: Add RemoveNode ([hdfchain/hdfd#1643](https://github.com/hdfchain/hdfd/pull/1643))
- rpctest: Add NodesConnected ([hdfchain/hdfd#1643](https://github.com/hdfchain/hdfd/pull/1643))
- dcrutil: Reduce global refs in addr unit tests ([hdfchain/hdfd#1666](https://github.com/hdfchain/hdfd/pull/1666))
- dcrutil: Consolidate tests into package ([hdfchain/hdfd#1669](https://github.com/hdfchain/hdfd/pull/1669))
- peer: Consolidate tests into package ([hdfchain/hdfd#1670](https://github.com/hdfchain/hdfd/pull/1670))
- wire: Add tests for BlockHeader (From)Bytes ([hdfchain/hdfd#1600](https://github.com/hdfchain/hdfd/pull/1600))
- wire: Add tests for MsgGetCFilter ([hdfchain/hdfd#1628](https://github.com/hdfchain/hdfd/pull/1628))
- dcrutil: Add tests for NewTxDeep ([hdfchain/hdfd#1684](https://github.com/hdfchain/hdfd/pull/1684))
- rpctest: Introduce VotingWallet ([hdfchain/hdfd#1668](https://github.com/hdfchain/hdfd/pull/1668))
- txscript: Add stake tx remove opcode tests ([hdfchain/hdfd#1210](https://github.com/hdfchain/hdfd/pull/1210))
- txscript: Move init func in benchmarks to top ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Add benchmark for script parsing ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Add benchmark for DisasmString ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Convert sighash calc tests ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Add benchmark for IsPayToScriptHash ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Add benchmarks for IsMutlsigScript ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Add benchmarks for IsMutlsigSigScript ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Add benchmark for GetSigOpCount ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Add tests for stake-tagged script hash ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Add benchmark for isAnyKindOfScriptHash ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Add benchmark for IsPushOnlyScript ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Add benchmark for GetPreciseSigOpCount ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Add benchmark for GetScriptClass ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Add benchmark for pay-to-pubkey scripts ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Add bench for pay-to-alt-pubkey scripts ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Add bench for pay-to-pubkey-hash scripts ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Add bench for pay-to-alt-pubkey-hash scripts ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Add bench for null scripts ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Add bench for stake submission scripts ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Add bench for stake generation scripts ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Add bench for stake revocation scripts ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Add bench for stake change scripts ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Add benchmark for ContainsStakeOpCodes ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Add benchmark for ExtractCoinbaseNullData ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Add CalcMultiSigStats benchmark ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Add multisig redeem script extract bench ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Add benchmark for PushedData ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Add benchmark for IsUnspendable ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Add tests for atomic swap extraction ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Add ExtractAtomicSwapDataPushes benches ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Add ExtractPkScriptAddrs benchmarks ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- txscript: Add ExtractPkScriptAltSigType benchmark ([hdfchain/hdfd#1656](https://github.com/hdfchain/hdfd/pull/1656))
- wire: Add tests for MsgGetCFTypes ([hdfchain/hdfd#1703](https://github.com/hdfchain/hdfd/pull/1703))
- blockchain: Allow named blocks in chaingen harness ([hdfchain/hdfd#1701](https://github.com/hdfchain/hdfd/pull/1701))
- txscript: Cleanup opcode removal by data tests ([hdfchain/hdfd#1702](https://github.com/hdfchain/hdfd/pull/1702))
- hdkeychain: Correct benchmark extended key ([hdfchain/hdfd#1696](https://github.com/hdfchain/hdfd/pull/1696))
- hdkeychain: Consolidate tests into package ([hdfchain/hdfd#1696](https://github.com/hdfchain/hdfd/pull/1696))
- hdkeychain: Use locally-scoped netparams in tests ([hdfchain/hdfd#1696](https://github.com/hdfchain/hdfd/pull/1696))
- hdkeychain: Use mock net params in tests ([hdfchain/hdfd#1696](https://github.com/hdfchain/hdfd/pull/1696))
- wire: Add tests for MsgGetCFHeaders ([hdfchain/hdfd#1720](https://github.com/hdfchain/hdfd/pull/1720))
- wire: Add tests for MsgCFHeaders ([hdfchain/hdfd#1732](https://github.com/hdfchain/hdfd/pull/1732))
- main/rpctest: Update for hdkeychain/v2 ([hdfchain/hdfd#1707](https://github.com/hdfchain/hdfd/pull/1707))
- rpctest: Allow custom miner on voting wallet ([hdfchain/hdfd#1751](https://github.com/hdfchain/hdfd/pull/1751))
- wire: Add tests for MsgCFilter ([hdfchain/hdfd#1741](https://github.com/hdfchain/hdfd/pull/1741))
- chaincfg; Add tests for required unique fields ([hdfchain/hdfd#1698](https://github.com/hdfchain/hdfd/pull/1698))
- fullblocktests: Add coinbase nulldata tests ([hdfchain/hdfd#1769](https://github.com/hdfchain/hdfd/pull/1769))
- dcrutil: Make docs example testable and correct it ([hdfchain/hdfd#1767](https://github.com/hdfchain/hdfd/pull/1767))
- dcrutil: Use mock addr params in tests ([hdfchain/hdfd#1767](https://github.com/hdfchain/hdfd/pull/1767))
- wire: assert MaxMessagePayload limit in tests ([hdfchain/hdfd#1755](https://github.com/hdfchain/hdfd/pull/1755))
- docker: use go 1.12 ([hdfchain/hdfd#1782](https://github.com/hdfchain/hdfd/pull/1782))
- docker: update alpine and include notes ([hdfchain/hdfd#1786](https://github.com/hdfchain/hdfd/pull/1786))
- hdkeychain: Correct a few comment typos ([hdfchain/hdfd#1796](https://github.com/hdfchain/hdfd/pull/1796))
- database: Use unique test db names for v2 module ([hdfchain/hdfd#1806](https://github.com/hdfchain/hdfd/pull/1806))
- main: Add database/v2 override for tests ([hdfchain/hdfd#1806](https://github.com/hdfchain/hdfd/pull/1806))
- gcs: Add benchmark for AddSigScript ([hdfchain/hdfd#1804](https://github.com/hdfchain/hdfd/pull/1804))
- txscript: Fix typo in script test data ([hdfchain/hdfd#1821](https://github.com/hdfchain/hdfd/pull/1821))
- database: Separate dbs for concurrent db tests ([hdfchain/hdfd#1824](https://github.com/hdfchain/hdfd/pull/1824))
- gcs: Overhaul tests and benchmarks ([hdfchain/hdfd#1845](https://github.com/hdfchain/hdfd/pull/1845))
- rpctest: Remove leftover debug print ([hdfchain/hdfd#1862](https://github.com/hdfchain/hdfd/pull/1862))
- txscript: Fix duplicate test name ([hdfchain/hdfd#1863](https://github.com/hdfchain/hdfd/pull/1863))
- gcs: Add benchmark for filter hashing ([hdfchain/hdfd#1853](https://github.com/hdfchain/hdfd/pull/1853))
- gcs: Add tests for bit reader/writer ([hdfchain/hdfd#1855](https://github.com/hdfchain/hdfd/pull/1855))
- peer: Ensure listener tests sync with messages ([hdfchain/hdfd#1874](https://github.com/hdfchain/hdfd/pull/1874))
- rpctest: remove always-nil error ([hdfchain/hdfd#1913](https://github.com/hdfchain/hdfd/pull/1913))
- rpctest: use errgroup to catch errors from go routines ([hdfchain/hdfd#1913](https://github.com/hdfchain/hdfd/pull/1913))

### Misc:

- release: Bump for 1.5 release cycle ([hdfchain/hdfd#1546](https://github.com/hdfchain/hdfd/pull/1546))
- mempool: Fix typo in fetchInputUtxos comment ([hdfchain/hdfd#1562](https://github.com/hdfchain/hdfd/pull/1562))
- blockchain: Fix typos found by misspell ([hdfchain/hdfd#1617](https://github.com/hdfchain/hdfd/pull/1617))
- dcrutil: Fix typos found by misspell ([hdfchain/hdfd#1617](https://github.com/hdfchain/hdfd/pull/1617))
- main: Write memprofile on shutdown ([hdfchain/hdfd#1655](https://github.com/hdfchain/hdfd/pull/1655))
- config: Parse network interfaces ([hdfchain/hdfd#1514](https://github.com/hdfchain/hdfd/pull/1514))
- config: Cleanup and simplify network info parsing ([hdfchain/hdfd#1706](https://github.com/hdfchain/hdfd/pull/1706))
- main: Rework windows service sod notification ([hdfchain/hdfd#1710](https://github.com/hdfchain/hdfd/pull/1710))
- multi: fix recent govet findings ([hdfchain/hdfd#1727](https://github.com/hdfchain/hdfd/pull/1727))
- rpcserver: Fix misspelling ([hdfchain/hdfd#1763](https://github.com/hdfchain/hdfd/pull/1763))
- chaincfg: Run gofmt -s ([hdfchain/hdfd#1776](https://github.com/hdfchain/hdfd/pull/1776))
- jsonrpc/types: Update copyright years ([hdfchain/hdfd#1794](https://github.com/hdfchain/hdfd/pull/1794))
- stake: Correct comment typo on Hash256PRNG ([hdfchain/hdfd#1803](https://github.com/hdfchain/hdfd/pull/1803))
- multi: Correct typos ([hdfchain/hdfd#1839](https://github.com/hdfchain/hdfd/pull/1839))
- wire: Fix a few messageError string typos ([hdfchain/hdfd#1840](https://github.com/hdfchain/hdfd/pull/1840))
- miningerror: Remove duplicate copyright ([hdfchain/hdfd#1860](https://github.com/hdfchain/hdfd/pull/1860))
- multi: Correct typos ([hdfchain/hdfd#1864](https://github.com/hdfchain/hdfd/pull/1864))

### Code Contributors (alphabetical order):

- Aaron Campbell
- Conner Fromknecht
- Dave Collins
- David Hill
- Donald Adu-Poku
- Hamid
- J Fixby
- Jamie Holdstock
- JoeGruffins
- Jonathan Chappelow
- Josh Rickmar
- Matheus Degiovani
- Nicola Larosa
- Olaoluwa Osuntokun
- Roei Erez
- Sarlor
- Victor Oliveira