# hdfd v1.2.0

This release of hdfd contains significant performance enhancements,
infrastructure improvements, improved access to chain-related information for
providing better SPV (Simplified Payment Verification) support, and other
quality assurance changes.

A significant amount of infrastructure work has also been done this release
cycle towards being able to support several planned scalability optimizations.

## Downgrade Warning

The database format in v1.2.0 is not compatible with previous versions of the
software.  This only affects downgrades as users upgrading from previous
versions will see a one time database migration.

Once this migration has been completed, it will no longer be possible to
downgrade to a previous version of the software without having to delete the
database and redownload the chain.

## Notable Changes

### Significantly Faster Startup

The startup time has been improved by roughly 17x on slower hard disk drives
(HDDs) and 8x on solid state drives (SSDs).

In order to achieve these speedups, there is a one time database migration, as
previously mentioned, that will likely take a while to complete (typically
around 5 to 6 minutes on HDDs and 2 to 3 minutes on SSDs).

### Support For DNS Seed Filtering

In order to better support the forthcoming SPV wallets, support for finding
other peers based upon their enabled services has been added.  This is useful
for both SPV wallets and full nodes since SPV wallets will require access to
full nodes in order to retrieve the necessary proofs and full nodes are
generally not interested in making outgoing connections to SPV wallets.

### Committed Filters

With the intention of supporting light clients, such as SPV wallets, in a
privacy-preserving way while still minimizing the amount of data that needs to
be downloaded, this release adds support for committed filters.  A committed
filter is a combination of a probalistic data structure that is used to test
whether an element is a member of a set with a predetermined collision
probability along with a commitment by consensus-validating full nodes to that
data.

A committed filter is created for every block which allows light clients to
download the filters and match against them locally rather than uploading
personal data to other nodes.

A new service flag is also provided to allow clients to discover nodes that
provide access to filters.

There is a one time database update to build and store the filters for all
existing historical blocks which will likely take a while to complete (typically
around 2 to 3 minutes on HDDs and 1 to 1.5 minutes on SSDs).

### Updated Atomic Swap Contracts

The standard checks for atomic swap contracts have been updated to ensure the
contracts enforce the secret size for safer support between chains with
disparate script rules.

### RPC Server Changes

#### New `getchaintips` RPC

A new RPC named `getchaintips` has been added which allows callers to query
information about the status of known side chains and their branch lengths.
It currently only provides support for side chains that have been seen while the
current instance of the process is running.  This will be further improved in
future releases.

## Changelog

All commits since the last release may be viewed on GitHub [here](https://github.com/hdfchain/hdfd/compare/v1.1.2...v1.2.0).

### Protocol and network:

- chaincfg: Add checkpoints for 1.2.0 release ([hdfchain/hdfd#1139](https://github.com/hdfchain/hdfd/pull/1139))
- chaincfg: Introduce new type DNSSeed ([hdfchain/hdfd#961](https://github.com/hdfchain/hdfd/pull/961))
- blockmanager: sync with the most updated peer ([hdfchain/hdfd#984](https://github.com/hdfchain/hdfd/pull/984))
- multi: remove MsgAlert ([hdfchain/hdfd#1161](https://github.com/hdfchain/hdfd/pull/1161))
- multi: Add initial committed filter (CF) support ([hdfchain/hdfd#1151](https://github.com/hdfchain/hdfd/pull/1151))

### Transaction relay (memory pool):

- txscript: Correct nulldata standardness check ([hdfchain/hdfd#935](https://github.com/hdfchain/hdfd/pull/935))
- mempool: Optimize orphan map limiting ([hdfchain/hdfd#1117](https://github.com/hdfchain/hdfd/pull/1117))
- mining: Fix duplicate txns in the prio heap ([hdfchain/hdfd#1108](https://github.com/hdfchain/hdfd/pull/1108))
- mining: Stop transactions losing their dependants ([hdfchain/hdfd#1109](https://github.com/hdfchain/hdfd/pull/1109))

### RPC:

- rpcserver: skip cert create when RPC is disabled ([hdfchain/hdfd#949](https://github.com/hdfchain/hdfd/pull/949))
- rpcserver: remove redundant checks in blockTemplateResult ([hdfchain/hdfd#826](https://github.com/hdfchain/hdfd/pull/826))
- rpcserver: assert network for validateaddress rpc ([hdfchain/hdfd#963](https://github.com/hdfchain/hdfd/pull/963))
- rpcserver: Do not rebroadcast stake transactions ([hdfchain/hdfd#973](https://github.com/hdfchain/hdfd/pull/973))
- dcrjson: add ticket fee field to PurchaseTicketCmd ([hdfchain/hdfd#902](https://github.com/hdfchain/hdfd/pull/902))
- hdfwalletextcmds: remove getseed ([hdfchain/hdfd#985](https://github.com/hdfchain/hdfd/pull/985))
- dcrjson: Add SweepAccountCmd & SweepAccountResult ([hdfchain/hdfd#1027](https://github.com/hdfchain/hdfd/pull/1027))
- rpcserver: add sweepaccount to the wallet list of commands ([hdfchain/hdfd#1028](https://github.com/hdfchain/hdfd/pull/1028))
- rpcserver: add batched request support (json 2.0) ([hdfchain/hdfd#841](https://github.com/hdfchain/hdfd/pull/841))
- dcrjson: include summary totals in GetBalanceResult ([hdfchain/hdfd#1062](https://github.com/hdfchain/hdfd/pull/1062))
- multi: Implement getchaintips JSON-RPC ([hdfchain/hdfd#1098](https://github.com/hdfchain/hdfd/pull/1098))
- rpcserver: Add hdfd version info to getversion RPC ([hdfchain/hdfd#1097](https://github.com/hdfchain/hdfd/pull/1097))
- rpcserver: Correct getblockheader result text ([hdfchain/hdfd#1104](https://github.com/hdfchain/hdfd/pull/1104))
- dcrjson: add StartAutoBuyerCmd & StopAutoBuyerCmd ([hdfchain/hdfd#903](https://github.com/hdfchain/hdfd/pull/903))
- dcrjson: fix typo for StartAutoBuyerCmd ([hdfchain/hdfd#1146](https://github.com/hdfchain/hdfd/pull/1146))
- dcrjson: require passphrase for StartAutoBuyerCmd ([hdfchain/hdfd#1147](https://github.com/hdfchain/hdfd/pull/1147))
- dcrjson: fix StopAutoBuyerCmd registration bug ([hdfchain/hdfd#1148](https://github.com/hdfchain/hdfd/pull/1148))
- blockchain: Support testnet stake diff estimation ([hdfchain/hdfd#1115](https://github.com/hdfchain/hdfd/pull/1115))
- rpcserver: fix jsonRPCRead data race ([hdfchain/hdfd#1157](https://github.com/hdfchain/hdfd/pull/1157))
- dcrjson: Add VerifySeedCmd ([hdfchain/hdfd#1160](https://github.com/hdfchain/hdfd/pull/1160))

### hdfd command-line flags and configuration:

- mempool: Rename RelayNonStd config option ([hdfchain/hdfd#1024](https://github.com/hdfchain/hdfd/pull/1024))
- sampleconfig: Update min relay fee ([hdfchain/hdfd#959](https://github.com/hdfchain/hdfd/pull/959))
- sampleconfig: Correct comment ([hdfchain/hdfd#1063](https://github.com/hdfchain/hdfd/pull/1063))
- multi: Expand ~ to correct home directory on all OSes ([hdfchain/hdfd#1041](https://github.com/hdfchain/hdfd/pull/1041))

### checkdevpremine utility changes:

- checkdevpremine: Remove --skipverify option ([hdfchain/hdfd#969](https://github.com/hdfchain/hdfd/pull/969))
- checkdevpremine: Implement --notls option ([hdfchain/hdfd#969](https://github.com/hdfchain/hdfd/pull/969))
- checkdevpremine: Make file naming consistent ([hdfchain/hdfd#969](https://github.com/hdfchain/hdfd/pull/969))
- checkdevpremine: Fix comment ([hdfchain/hdfd#969](https://github.com/hdfchain/hdfd/pull/969))
- checkdevpremine: Remove utility ([hdfchain/hdfd#1068](https://github.com/hdfchain/hdfd/pull/1068))

### Documentation:

- fullblocktests: Add missing doc.go file ([hdfchain/hdfd#956](https://github.com/hdfchain/hdfd/pull/956))
- docs: Add fullblocktests entry and make consistent ([hdfchain/hdfd#956](https://github.com/hdfchain/hdfd/pull/956))
- docs: Add mempool entry to developer tools section ([hdfchain/hdfd#1058](https://github.com/hdfchain/hdfd/pull/1058))
- mempool: Add docs.go and flesh out README.md ([hdfchain/hdfd#1058](https://github.com/hdfchain/hdfd/pull/1058))
- docs: document packages and fix typo  ([hdfchain/hdfd#965](https://github.com/hdfchain/hdfd/pull/965))
- docs: rpcclient is now part of the main hdfd repo ([hdfchain/hdfd#970](https://github.com/hdfchain/hdfd/pull/970))
- dcrjson: Update README.md ([hdfchain/hdfd#982](https://github.com/hdfchain/hdfd/pull/982))
- docs: Remove carriage return ([hdfchain/hdfd#1106](https://github.com/hdfchain/hdfd/pull/1106))
- Adjust README.md for new Go version ([hdfchain/hdfd#1105](https://github.com/hdfchain/hdfd/pull/1105))
- docs: document how to use go test -coverprofile ([hdfchain/hdfd#1107](https://github.com/hdfchain/hdfd/pull/1107))
- addrmgr: Improve documentation ([hdfchain/hdfd#1125](https://github.com/hdfchain/hdfd/pull/1125))
- docs: Fix links for internal packages ([hdfchain/hdfd#1144](https://github.com/hdfchain/hdfd/pull/1144))

### Developer-related package changes:

- chaingen: Add revocation generation infrastructure ([hdfchain/hdfd#1120](https://github.com/hdfchain/hdfd/pull/1120))
- txscript: Add null data script creator ([hdfchain/hdfd#943](https://github.com/hdfchain/hdfd/pull/943))
- txscript: Cleanup and improve NullDataScript tests ([hdfchain/hdfd#943](https://github.com/hdfchain/hdfd/pull/943))
- txscript: Allow external signature hash calc ([hdfchain/hdfd#951](https://github.com/hdfchain/hdfd/pull/951))
- secp256k1: update func signatures ([hdfchain/hdfd#934](https://github.com/hdfchain/hdfd/pull/934))
- txscript: enforce MaxDataCarrierSize for GenerateProvablyPruneableOut ([hdfchain/hdfd#953](https://github.com/hdfchain/hdfd/pull/953))
- txscript: Remove OP_SMALLDATA ([hdfchain/hdfd#954](https://github.com/hdfchain/hdfd/pull/954))
- blockchain: Accept header in CheckProofOfWork ([hdfchain/hdfd#977](https://github.com/hdfchain/hdfd/pull/977))
- blockchain: Make func definition style consistent ([hdfchain/hdfd#983](https://github.com/hdfchain/hdfd/pull/983))
- blockchain: only fetch the parent block in BFFastAdd ([hdfchain/hdfd#972](https://github.com/hdfchain/hdfd/pull/972))
- blockchain: Switch to FindSpentTicketsInBlock ([hdfchain/hdfd#915](https://github.com/hdfchain/hdfd/pull/915))
- stake: Add Hash256PRNG init vector support ([hdfchain/hdfd#986](https://github.com/hdfchain/hdfd/pull/986))
- blockchain/stake: Use Hash256PRNG init vector ([hdfchain/hdfd#987](https://github.com/hdfchain/hdfd/pull/987))
- blockchain: Don't store full header in block node ([hdfchain/hdfd#988](https://github.com/hdfchain/hdfd/pull/988))
- blockchain: Reconstruct headers from block nodes ([hdfchain/hdfd#989](https://github.com/hdfchain/hdfd/pull/989))
- stake/multi: Don't return errors for IsX functions ([hdfchain/hdfd#995](https://github.com/hdfchain/hdfd/pull/995))
- blockchain: Rename block index to main chain index ([hdfchain/hdfd#996](https://github.com/hdfchain/hdfd/pull/996))
- blockchain: Refactor main block index logic ([hdfchain/hdfd#990](https://github.com/hdfchain/hdfd/pull/990))
- blockchain: Use hash values in structs ([hdfchain/hdfd#992](https://github.com/hdfchain/hdfd/pull/992))
- blockchain: Remove unused dump function ([hdfchain/hdfd#1001](https://github.com/hdfchain/hdfd/pull/1001))
- blockchain: Generalize and optimize chain reorg ([hdfchain/hdfd#997](https://github.com/hdfchain/hdfd/pull/997))
- blockchain: Pass parent block in connection code ([hdfchain/hdfd#998](https://github.com/hdfchain/hdfd/pull/998))
- blockchain: Explicit block fetch semanticss ([hdfchain/hdfd#999](https://github.com/hdfchain/hdfd/pull/999))
- blockchain: Use next detach block in reorg chain ([hdfchain/hdfd#1002](https://github.com/hdfchain/hdfd/pull/1002))
- blockchain: Limit header sanity check to header ([hdfchain/hdfd#1003](https://github.com/hdfchain/hdfd/pull/1003))
- blockchain: Validate num votes in header sanity ([hdfchain/hdfd#1005](https://github.com/hdfchain/hdfd/pull/1005))
- blockchain: Validate max votes in header sanity ([hdfchain/hdfd#1006](https://github.com/hdfchain/hdfd/pull/1006))
- blockchain: Validate stake diff in header context ([hdfchain/hdfd#1004](https://github.com/hdfchain/hdfd/pull/1004))
- blockchain: No votes/revocations in header sanity ([hdfchain/hdfd#1007](https://github.com/hdfchain/hdfd/pull/1007))
- blockchain: Validate max purchases in header sanity ([hdfchain/hdfd#1008](https://github.com/hdfchain/hdfd/pull/1008))
- blockchain: Validate early votebits in header sanity ([hdfchain/hdfd#1009](https://github.com/hdfchain/hdfd/pull/1009))
- blockchain: Validate block height in header context ([hdfchain/hdfd#1010](https://github.com/hdfchain/hdfd/pull/1010))
- blockchain: Move check block context func ([hdfchain/hdfd#1011](https://github.com/hdfchain/hdfd/pull/1011))
- blockchain: Block sanity cleanup and consistency ([hdfchain/hdfd#1012](https://github.com/hdfchain/hdfd/pull/1012))
- blockchain: Remove dup ticket purchase value check ([hdfchain/hdfd#1013](https://github.com/hdfchain/hdfd/pull/1013))
- blockchain: Only tickets before SVH in block sanity ([hdfchain/hdfd#1014](https://github.com/hdfchain/hdfd/pull/1014))
- blockchain: Remove unused vote bits function ([hdfchain/hdfd#1015](https://github.com/hdfchain/hdfd/pull/1015))
- blockchain: Move upgrade-only code to upgrade.go ([hdfchain/hdfd#1016](https://github.com/hdfchain/hdfd/pull/1016))
- stake: Static assert of vote commitment ([hdfchain/hdfd#1020](https://github.com/hdfchain/hdfd/pull/1020))
- blockchain: Remove unused error code ([hdfchain/hdfd#1021](https://github.com/hdfchain/hdfd/pull/1021))
- blockchain: Improve readability of parent approval ([hdfchain/hdfd#1022](https://github.com/hdfchain/hdfd/pull/1022))
- peer: rename mruinvmap, mrunoncemap to lruinvmap, lrunoncemap ([hdfchain/hdfd#976](https://github.com/hdfchain/hdfd/pull/976))
- peer: rename noncemap to noncecache ([hdfchain/hdfd#976](https://github.com/hdfchain/hdfd/pull/976))
- peer: rename inventorymap to inventorycache ([hdfchain/hdfd#976](https://github.com/hdfchain/hdfd/pull/976))
- connmgr: convert state to atomic ([hdfchain/hdfd#1025](https://github.com/hdfchain/hdfd/pull/1025))
- blockchain/mining: Full checks in CCB ([hdfchain/hdfd#1017](https://github.com/hdfchain/hdfd/pull/1017))
- blockchain: Validate pool size in header context ([hdfchain/hdfd#1018](https://github.com/hdfchain/hdfd/pull/1018))
- blockchain: Vote commitments in block sanity ([hdfchain/hdfd#1023](https://github.com/hdfchain/hdfd/pull/1023))
- blockchain: Validate early final state is zero ([hdfchain/hdfd#1031](https://github.com/hdfchain/hdfd/pull/1031))
- blockchain: Validate final state in header context ([hdfchain/hdfd#1034](https://github.com/hdfchain/hdfd/pull/1033))
- blockchain: Max revocations in block sanity ([hdfchain/hdfd#1034](https://github.com/hdfchain/hdfd/pull/1034))
- blockchain: Allowed stake txns in block sanity ([hdfchain/hdfd#1035](https://github.com/hdfchain/hdfd/pull/1035))
- blockchain: Validate allowed votes in block context ([hdfchain/hdfd#1036](https://github.com/hdfchain/hdfd/pull/1036))
- blockchain: Validate allowed revokes in blk contxt ([hdfchain/hdfd#1037](https://github.com/hdfchain/hdfd/pull/1037))
- blockchain/stake: Rename tix spent to tix voted ([hdfchain/hdfd#1038](https://github.com/hdfchain/hdfd/pull/1038))
- txscript: Require atomic swap contracts to specify the secret size ([hdfchain/hdfd#1039](https://github.com/hdfchain/hdfd/pull/1039))
- blockchain: Remove unused struct ([hdfchain/hdfd#1043](https://github.com/hdfchain/hdfd/pull/1043))
- blockchain: Store side chain blocks in database ([hdfchain/hdfd#1000](https://github.com/hdfchain/hdfd/pull/1000))
- blockchain: Simplify initial chain state ([hdfchain/hdfd#1045](https://github.com/hdfchain/hdfd/pull/1045))
- blockchain: Rework database versioning ([hdfchain/hdfd#1047](https://github.com/hdfchain/hdfd/pull/1047))
- blockchain: Don't require chain for db upgrades ([hdfchain/hdfd#1051](https://github.com/hdfchain/hdfd/pull/1051))
- blockchain/indexers: Allow interrupts ([hdfchain/hdfd#1052](https://github.com/hdfchain/hdfd/pull/1052))
- blockchain: Remove old version information ([hdfchain/hdfd#1055](https://github.com/hdfchain/hdfd/pull/1055))
- stake: optimize FindSpentTicketsInBlock slightly ([hdfchain/hdfd#1049](https://github.com/hdfchain/hdfd/pull/1049))
- blockchain: Do not accept orphans/genesis block ([hdfchain/hdfd#1057](https://github.com/hdfchain/hdfd/pull/1057))
- blockchain: Separate node ticket info population ([hdfchain/hdfd#1056](https://github.com/hdfchain/hdfd/pull/1056))
- blockchain: Accept parent in blockNode constructor ([hdfchain/hdfd#1056](https://github.com/hdfchain/hdfd/pull/1056))
- blockchain: Combine ErrDoubleSpend & ErrMissingTx ([hdfchain/hdfd#1064](https://github.com/hdfchain/hdfd/pull/1064))
- blockchain: Calculate the lottery IV on demand ([hdfchain/hdfd#1065](https://github.com/hdfchain/hdfd/pull/1065))
- blockchain: Simplify add/remove node logic ([hdfchain/hdfd#1067](https://github.com/hdfchain/hdfd/pull/1067))
- blockchain: Infrastructure to manage block index ([hdfchain/hdfd#1044](https://github.com/hdfchain/hdfd/pull/1044))
- blockchain: Add block validation status to index ([hdfchain/hdfd#1044](https://github.com/hdfchain/hdfd/pull/1044))
- blockchain: Migrate to new block indexuse it ([hdfchain/hdfd#1044](https://github.com/hdfchain/hdfd/pull/1044))
- blockchain: Lookup child in force head reorg ([hdfchain/hdfd#1070](https://github.com/hdfchain/hdfd/pull/1070))
- blockchain: Refactor block idx entry serialization ([hdfchain/hdfd#1069](https://github.com/hdfchain/hdfd/pull/1069))
- blockchain: Limit GetStakeVersions count ([hdfchain/hdfd#1071](https://github.com/hdfchain/hdfd/pull/1071))
- blockchain: Remove dry run flag ([hdfchain/hdfd#1073](https://github.com/hdfchain/hdfd/pull/1073))
- blockchain: Remove redundant stake ver calc func ([hdfchain/hdfd#1087](https://github.com/hdfchain/hdfd/pull/1087))
- blockchain: Reduce GetGeneration to TipGeneration ([hdfchain/hdfd#1083](https://github.com/hdfchain/hdfd/pull/1083))
- blockchain: Add chain tip tracking ([hdfchain/hdfd#1084](https://github.com/hdfchain/hdfd/pull/1084))
- blockchain: Switch tip generation to chain tips ([hdfchain/hdfd#1085](https://github.com/hdfchain/hdfd/pull/1085))
- blockchain: Simplify voter version calculation ([hdfchain/hdfd#1088](https://github.com/hdfchain/hdfd/pull/1088))
- blockchain: Remove unused threshold serialization ([hdfchain/hdfd#1092](https://github.com/hdfchain/hdfd/pull/1092))
- blockchain: Simplify chain tip tracking ([hdfchain/hdfd#1092](https://github.com/hdfchain/hdfd/pull/1092))
- blockchain: Cache tip and parent at init ([hdfchain/hdfd#1100](https://github.com/hdfchain/hdfd/pull/1100))
- mining: Obtain block by hash instead of top block ([hdfchain/hdfd#1094](https://github.com/hdfchain/hdfd/pull/1094))
- blockchain: Remove unused GetTopBlock function ([hdfchain/hdfd#1094](https://github.com/hdfchain/hdfd/pull/1094))
- multi: Rename BIP0111Version to NodeBloomVersion ([hdfchain/hdfd#1112](https://github.com/hdfchain/hdfd/pull/1112))
- mining/mempool: Move priority code to mining pkg ([hdfchain/hdfd#1110](https://github.com/hdfchain/hdfd/pull/1110))
- mining: Use single uint64 coinbase extra nonce ([hdfchain/hdfd#1116](https://github.com/hdfchain/hdfd/pull/1116))
- mempool/mining: Clarify tree validity semantics ([hdfchain/hdfd#1118](https://github.com/hdfchain/hdfd/pull/1118))
- mempool/mining: TxSource separation ([hdfchain/hdfd#1119](https://github.com/hdfchain/hdfd/pull/1119))
- connmgr: Use same Dial func signature as net.Dial ([hdfchain/hdfd#1113](https://github.com/hdfchain/hdfd/pull/1113))
- addrmgr: Declutter package API ([hdfchain/hdfd#1124](https://github.com/hdfchain/hdfd/pull/1124))
- mining: Correct initial template generation ([hdfchain/hdfd#1122](https://github.com/hdfchain/hdfd/pull/1122))
- cpuminer: Use header for extra nonce ([hdfchain/hdfd#1123](https://github.com/hdfchain/hdfd/pull/1123))
- addrmgr: Make writing of peers file safer ([hdfchain/hdfd#1126](https://github.com/hdfchain/hdfd/pull/1126))
- addrmgr: Save peers file only if necessary ([hdfchain/hdfd#1127](https://github.com/hdfchain/hdfd/pull/1127))
- addrmgr: Factor out common code ([hdfchain/hdfd#1138](https://github.com/hdfchain/hdfd/pull/1138))
- addrmgr: Improve isBad() performance ([hdfchain/hdfd#1134](https://github.com/hdfchain/hdfd/pull/1134))
- dcrutil: Disallow creation of hybrid P2PK addrs ([hdfchain/hdfd#1154](https://github.com/hdfchain/hdfd/pull/1154))
- chainec/dcrec: Remove hybrid pubkey support ([hdfchain/hdfd#1155](https://github.com/hdfchain/hdfd/pull/1155))
- blockchain: Only fetch inputs once in connect txns ([hdfchain/hdfd#1152](https://github.com/hdfchain/hdfd/pull/1152))
- indexers: Provide interface for index removal ([hdfchain/hdfd#1158](https://github.com/hdfchain/hdfd/pull/1158))

### Testing and Quality Assurance:

- travis: set GOVERSION environment properly ([hdfchain/hdfd#958](https://github.com/hdfchain/hdfd/pull/958))
- stake: Override false positive vet error ([hdfchain/hdfd#960](https://github.com/hdfchain/hdfd/pull/960))
- docs: make example code compile ([hdfchain/hdfd#970](https://github.com/hdfchain/hdfd/pull/970))
- blockchain: Add median time tests ([hdfchain/hdfd#991](https://github.com/hdfchain/hdfd/pull/991))
- chaingen: Update vote commitments on hdr updates ([hdfchain/hdfd#1023](https://github.com/hdfchain/hdfd/pull/1023))
- fullblocktests: Add tests for early final state ([hdfchain/hdfd#1031](https://github.com/hdfchain/hdfd/pull/1031))
- travis: test in docker container ([hdfchain/hdfd#1053](https://github.com/hdfchain/hdfd/pull/1053))
- blockchain: Correct error stringer tests ([hdfchain/hdfd#1066](https://github.com/hdfchain/hdfd/pull/1066))
- blockchain: Remove superfluous reorg tests ([hdfchain/hdfd#1072](https://github.com/hdfchain/hdfd/pull/1072))
- blockchain: Use chaingen for forced reorg tests ([hdfchain/hdfd#1074](https://github.com/hdfchain/hdfd/pull/1074))
- blockchain: Remove superfluous test checks ([hdfchain/hdfd#1075](https://github.com/hdfchain/hdfd/pull/1075))
- blockchain: move block validation rule tests into fullblocktests ([hdfchain/hdfd#1060](https://github.com/hdfchain/hdfd/pull/1060))
- fullblocktests: Cleanup after refactor ([hdfchain/hdfd#1080](https://github.com/hdfchain/hdfd/pull/1080))
- chaingen: Prevent dup block names in NextBlock ([hdfchain/hdfd#1079](https://github.com/hdfchain/hdfd/pull/1079))
- blockchain: Remove duplicate val tests ([hdfchain/hdfd#1082](https://github.com/hdfchain/hdfd/pull/1082))
- chaingen: Break dependency on blockchain ([hdfchain/hdfd#1076](https://github.com/hdfchain/hdfd/pull/1076))
- blockchain: Consolidate tests into the main package ([hdfchain/hdfd#1077](https://github.com/hdfchain/hdfd/pull/1077))
- chaingen: Export vote commitment script function ([hdfchain/hdfd#1081](https://github.com/hdfchain/hdfd/pull/1081))
- fullblocktests: Improve vote on wrong block tests ([hdfchain/hdfd#1081](https://github.com/hdfchain/hdfd/pull/1081))
- chaingen: Export func to check if block is solved ([hdfchain/hdfd#1089](https://github.com/hdfchain/hdfd/pull/1089))
- fullblocktests: Use new exported IsSolved func ([hdfchain/hdfd#1089](https://github.com/hdfchain/hdfd/pull/1089))
- chaingen: Accept mungers for create premine block ([hdfchain/hdfd#1090](https://github.com/hdfchain/hdfd/pull/1090))
- blockchain: Add tests for chain tip tracking ([hdfchain/hdfd#1096](https://github.com/hdfchain/hdfd/pull/1096))
- blockchain: move block validation rule tests into fullblocktests (2/x) ([hdfchain/hdfd#1095](https://github.com/hdfchain/hdfd/pull/1095))
- addrmgr: Remove obsolete coverage script ([hdfchain/hdfd#1103](https://github.com/hdfchain/hdfd/pull/1103))
- chaingen: Track expected blk heights separately ([hdfchain/hdfd#1101](https://github.com/hdfchain/hdfd/pull/1101))
- addrmgr: Improve test coverage ([hdfchain/hdfd#1111](https://github.com/hdfchain/hdfd/pull/1111))
- chaingen: Add revocation generation infrastructure ([hdfchain/hdfd#1120](https://github.com/hdfchain/hdfd/pull/1120))
- fullblocktests: Add some basic revocation tests ([hdfchain/hdfd#1121](https://github.com/hdfchain/hdfd/pull/1121))
- addrmgr: Test removal of corrupt peers file ([hdfchain/hdfd#1129](https://github.com/hdfchain/hdfd/pull/1129))

### Misc:

- release: Bump for v1.2.0 ([hdfchain/hdfd#1140](https://github.com/hdfchain/hdfd/pull/1140))
- goimports -w . ([hdfchain/hdfd#968](https://github.com/hdfchain/hdfd/pull/968))
- dep: sync ([hdfchain/hdfd#980](https://github.com/hdfchain/hdfd/pull/980))
- multi: Simplify code per gosimple linter ([hdfchain/hdfd#993](https://github.com/hdfchain/hdfd/pull/993))
- multi: various cleanups ([hdfchain/hdfd#1019](https://github.com/hdfchain/hdfd/pull/1019))
- multi: release the mutex earlier ([hdfchain/hdfd#1026](https://github.com/hdfchain/hdfd/pull/1026))
- multi: fix some maligned linter warnings ([hdfchain/hdfd#1025](https://github.com/hdfchain/hdfd/pull/1025))
- blockchain: Correct a few log statements ([hdfchain/hdfd#1042](https://github.com/hdfchain/hdfd/pull/1042))
- mempool: cleaner ([hdfchain/hdfd#1050](https://github.com/hdfchain/hdfd/pull/1050))
- multi: fix misspell linter warnings ([hdfchain/hdfd#1054](https://github.com/hdfchain/hdfd/pull/1054))
- dep: sync ([hdfchain/hdfd#1091](https://github.com/hdfchain/hdfd/pull/1091))
- multi: Properly capitalize Hdfchain ([hdfchain/hdfd#1102](https://github.com/hdfchain/hdfd/pull/1102))
- build: Correct semver build handling ([hdfchain/hdfd#1097](https://github.com/hdfchain/hdfd/pull/1097))
- main: Make func definition style consistent ([hdfchain/hdfd#1114](https://github.com/hdfchain/hdfd/pull/1114))
- main: Allow semver prerel via linker flags ([hdfchain/hdfd#1128](https://github.com/hdfchain/hdfd/pull/1128))

### Code Contributors (alphabetical order):

- Andrew Chiw
- Daniel Krawsiz
- Dave Collins
- David Hill
- Donald Adu-Poku
- Javed Khan
- Jolan Luff
- Jon Gillham
- Josh Rickmar
- Markus Richter
- Matheus Degiovani
- Ryan Vacek
