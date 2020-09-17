# hdfd v1.3.0

This release of hdfd contains significant performance enhancements for startup
speed, validation, and network operations that directly benefit lightweight
clients, such as SPV (Simplified Payment Verification) wallets, a policy change
to reduce the default minimum transaction fee rate, a new public test network
version, removal of bloom filter support, infrastructure improvements, and other
quality assurance changes.

**It is highly recommended that everyone upgrade to this latest release as it
contains many important scalability improvements and is required to be able to
use the new public test network.**

## Downgrade Warning

The database format in v1.3.0 is not compatible with previous versions of the
software.  This only affects downgrades as users upgrading from previous
versions will see a one time database migration.

Once this migration has been completed, it will no longer be possible to
downgrade to a previous version of the software without having to delete the
database and redownload the chain.

## Notable Changes

### Reduction of Default Minimum Transaction Fee Rate Policy

The default setting for the policy which specifies the minimum transaction fee
rate that will be accepted and relayed to the rest of the network has been
reduced to 0.0001 HDF/kB (10,000 atoms/kB) from the previous value of 0.001
HDF/kB (100,000 atoms/kB).

Transactions should not attempt to use the reduced fee rate until the majority
of the network has upgraded to this release as otherwise the transactions will
likely have issues relaying through the network since old nodes that have not
updated their policy will reject them due to not paying a high enough fee.

### Several Speed Optimizations

This release contains several enhancements to improve speed for startup,
the initial sync process, validation, and network operations.

In order to achieve these speedups, there is a one time database migration, as
previously mentioned, that typically only takes a few seconds to complete on
most hardware.

#### Further Improved Startup Speed

The startup time has been improved by roughly 2x on both slower hard disk drives
(HDDs) and solid state drives (SSDs) as compared to v1.2.0.

#### Significantly Faster Network Operations

The ability to serve information to other peers on the network has received
several optimizations which, in addition to generally improving the overall
scalability and throughput of the network, also directly benefits SPV
(Simplified Payment Verification) clients by delivering the block headers they
require roughly 3x to 4x faster.

#### Signature Hash Calculation Optimization

Part of validating that transactions are only spending coins that the owner has
authorized involves ensuring the validity of cryptographic signatures.  This
release provides a speedup of about 75% to a key portion of that validation
which results in a roughly 20% faster initial sync process.

### Bloom Filters Removal

Bloom filters were deprecated as of the last release in favor of the more recent
privacy-preserving GCS committed filters.  Consequently, this release removes
support for bloom filters completely.  There are no known clients which use
bloom filters, however, if there are any unknown clients which use them, those
clients will need to be updated to use the GCS committed filters accordingly.

### Public Test Network Version 3

The public test network has been reset and bumped to version 3.  All of the new
consensus rules voted in by version 2 of the public test network have been
retained and are therefore active on the new version 3 test network without
having to vote them in again.

## Changelog

All commits since the last release may be viewed on GitHub [here](https://github.com/hdfchain/hdfd/compare/release-v1.2.0...release-v1.3.0).

### Protocol and network:

- chaincfg: Add checkpoints for 1.3.0 release ([hdfchain/hdfd#1385](https://github.com/hdfchain/hdfd/pull/1385))
- multi: Remove everything to do about bloom filters ([hdfchain/hdfd#1162](https://github.com/hdfchain/hdfd/pull/1162))
- wire: Remove TxSerializeWitnessSigning ([hdfchain/hdfd#1180](https://github.com/hdfchain/hdfd/pull/1180))
- addrmgr: Skip low quality addresses for getaddr ([hdfchain/hdfd#1135](https://github.com/hdfchain/hdfd/pull/1135))
- addrmgr: Fix race in save peers ([hdfchain/hdfd#1259](https://github.com/hdfchain/hdfd/pull/1259))
- server: Only respond to getaddr once per conn ([hdfchain/hdfd#1257](https://github.com/hdfchain/hdfd/pull/1257))
- peer: Rework version negotiation ([hdfchain/hdfd#1250](https://github.com/hdfchain/hdfd/pull/1250))
- peer: Allow OnVersion callback to reject peer ([hdfchain/hdfd#1251](https://github.com/hdfchain/hdfd/pull/1251))
- server: Reject outbound conns to non-full nodes ([hdfchain/hdfd#1252](https://github.com/hdfchain/hdfd/pull/1252))
- peer: Improve net address service adverts ([hdfchain/hdfd#1253](https://github.com/hdfchain/hdfd/pull/1253))
- addrmgr: Expose method to update services ([hdfchain/hdfd#1254](https://github.com/hdfchain/hdfd/pull/1254))
- server: Update addrmgr services on outbound conns ([hdfchain/hdfd#1254](https://github.com/hdfchain/hdfd/pull/1254))
- server: Use local inbound var in version handler ([hdfchain/hdfd#1255](https://github.com/hdfchain/hdfd/pull/1255))
- server: Only advertise local addr when current ([hdfchain/hdfd#1256](https://github.com/hdfchain/hdfd/pull/1256))
- server: Use local addr var in version handler ([hdfchain/hdfd#1258](https://github.com/hdfchain/hdfd/pull/1258))
- chaincfg: split params into per-network files ([hdfchain/hdfd#1265](https://github.com/hdfchain/hdfd/pull/1265))
- server: Always reply to getheaders with headers ([hdfchain/hdfd#1295](https://github.com/hdfchain/hdfd/pull/1295))
- addrmgr: skip never-successful addresses ([hdfchain/hdfd#1313](https://github.com/hdfchain/hdfd/pull/1313))
- multi: Introduce default coin type for SLIP0044 ([hdfchain/hdfd#1293](https://github.com/hdfchain/hdfd/pull/1293))
- blockchain: Modify diff redux logic for testnet ([hdfchain/hdfd#1387](https://github.com/hdfchain/hdfd/pull/1387))
- multi: Reset testnet and bump to version 3 ([hdfchain/hdfd#1387](https://github.com/hdfchain/hdfd/pull/1387))
- multi: Remove testnet version 2 defs and refs ([hdfchain/hdfd#1387](https://github.com/hdfchain/hdfd/pull/1387))

### Transaction relay (memory pool):

- policy: Lower default relay fee to 0.0001/kB ([hdfchain/hdfd#1202](https://github.com/hdfchain/hdfd/pull/1202))
- mempool: Use blockchain for tx expiry check ([hdfchain/hdfd#1199](https://github.com/hdfchain/hdfd/pull/1199))
- mempool: use secp256k1 functions directly ([hdfchain/hdfd#1213](https://github.com/hdfchain/hdfd/pull/1213))
- mempool: Make expiry pruning self contained ([hdfchain/hdfd#1378](https://github.com/hdfchain/hdfd/pull/1378))
- mempool: Stricter orphan evaluation and eviction ([hdfchain/hdfd#1207](https://github.com/hdfchain/hdfd/pull/1207))
- mempool: use secp256k1 functions directly ([hdfchain/hdfd#1213](https://github.com/hdfchain/hdfd/pull/1213))
- multi: add specialized rebroadcast handling for stake txs ([hdfchain/hdfd#979](https://github.com/hdfchain/hdfd/pull/979))
- mempool: Make expiry pruning self contained ([hdfchain/hdfd#1378](https://github.com/hdfchain/hdfd/pull/1378))

### RPC:

- rpcserver: Improve JSON-RPC compatibility ([hdfchain/hdfd#1150](https://github.com/hdfchain/hdfd/pull/1150))
- rpcserver: Correct rebroadcastwinners handler ([hdfchain/hdfd#1234](https://github.com/hdfchain/hdfd/pull/1234))
- dcrjson: Add Expiry field to CreateRawTransactionCmd ([hdfchain/hdfd#1149](https://github.com/hdfchain/hdfd/pull/1149))
- dcrjson: add estimatesmartfee ([hdfchain/hdfd#1201](https://github.com/hdfchain/hdfd/pull/1201))
- rpc: Use upstream gorilla/websocket ([hdfchain/hdfd#1218](https://github.com/hdfchain/hdfd/pull/1218))
- dcrjson: add createvotingaccount and dropvotingaccount rpc methods ([hdfchain/hdfd#1217](https://github.com/hdfchain/hdfd/pull/1217))
- multi: Change NoSplitTransaction param to SplitTx ([hdfchain/hdfd#1231](https://github.com/hdfchain/hdfd/pull/1231))
- rpcclient: pass default value for NewPurchaseTicketCmd's comment param ([hdfchain/hdfd#1232](https://github.com/hdfchain/hdfd/pull/1232))
- multi: No winning ticket ntfns for big reorg depth ([hdfchain/hdfd#1235](https://github.com/hdfchain/hdfd/pull/1235))
- multi: modify PurchaseTicketCmd ([hdfchain/hdfd#1241](https://github.com/hdfchain/hdfd/pull/1241))
- multi: move extension commands into associated normal command files ([hdfchain/hdfd#1238](https://github.com/hdfchain/hdfd/pull/1238))
- dcrjson: Fix NewCreateRawTransactionCmd comment ([hdfchain/hdfd#1262](https://github.com/hdfchain/hdfd/pull/1262))
- multi: revert TicketChange addition to PurchaseTicketCmd ([hdfchain/hdfd#1278](https://github.com/hdfchain/hdfd/pull/1278))
- rpcclient: Implement fmt.Stringer for Client ([hdfchain/hdfd#1298](https://github.com/hdfchain/hdfd/pull/1298))
- multi: add amount field to TransactionInput ([hdfchain/hdfd#1297](https://github.com/hdfchain/hdfd/pull/1297))
- dcrjson: Ready GetStakeInfoResult for SPV wallets ([hdfchain/hdfd#1333](https://github.com/hdfchain/hdfd/pull/1333))
- dcrjson: add fundrawtransaction command ([hdfchain/hdfd#1316](https://github.com/hdfchain/hdfd/pull/1316))
- dcrjson: Make linter happy by renaming Id to ID ([hdfchain/hdfd#1374](https://github.com/hdfchain/hdfd/pull/1374))
- dcrjson: Remove unused vote bit concat codec funcs ([hdfchain/hdfd#1384](https://github.com/hdfchain/hdfd/pull/1384))
- rpcserver: Cleanup cfilter handling ([hdfchain/hdfd#1398](https://github.com/hdfchain/hdfd/pull/1398))

### hdfd command-line flags and configuration:

- multi: Correct clean and expand path handling ([hdfchain/hdfd#1186](https://github.com/hdfchain/hdfd/pull/1186))

### hdfctl utility changes:

- hdfctl: Fix --skipverify failing if rpc.cert not found ([hdfchain/hdfd#1163](https://github.com/hdfchain/hdfd/pull/1163))

### Documentation:

- hdkeychain: Correct hash algorithm in comment ([hdfchain/hdfd#1171](https://github.com/hdfchain/hdfd/pull/1171))
- Fix typo in blockchain ([hdfchain/hdfd#1185](https://github.com/hdfchain/hdfd/pull/1185))
- docs: Update node.js example for v8.11.1 LTS ([hdfchain/hdfd#1209](https://github.com/hdfchain/hdfd/pull/1209))
- docs: Update txaccepted method in json_rpc_api.md ([hdfchain/hdfd#1242](https://github.com/hdfchain/hdfd/pull/1242))
- docs: Correct blockmaxsize and blockprioritysize ([hdfchain/hdfd#1339](https://github.com/hdfchain/hdfd/pull/1339))
- server: Correct comment in getblocks handler ([hdfchain/hdfd#1269](https://github.com/hdfchain/hdfd/pull/1269))
- config: Fix typo ([hdfchain/hdfd#1274](https://github.com/hdfchain/hdfd/pull/1274))
- multi: Fix badges in README ([hdfchain/hdfd#1277](https://github.com/hdfchain/hdfd/pull/1277))
- stake: Correct comment in connectNode ([hdfchain/hdfd#1325](https://github.com/hdfchain/hdfd/pull/1325))
- txscript: Update comments for removal of flags ([hdfchain/hdfd#1336](https://github.com/hdfchain/hdfd/pull/1336))
- docs: Update docs for versioned modules ([hdfchain/hdfd#1391](https://github.com/hdfchain/hdfd/pull/1391))
- mempool: Correct min relay tx fee comment to HDF ([hdfchain/hdfd#1396](https://github.com/hdfchain/hdfd/pull/1396))

### Developer-related package and module changes:

- blockchain: CheckConnectBlockTemplate with tests ([hdfchain/hdfd#1086](https://github.com/hdfchain/hdfd/pull/1086))
- addrmgr: Simplify package API ([hdfchain/hdfd#1136](https://github.com/hdfchain/hdfd/pull/1136))
- txscript: Remove unused strict multisig flag ([hdfchain/hdfd#1203](https://github.com/hdfchain/hdfd/pull/1203))
- txscript: Move sig hash logic to separate file ([hdfchain/hdfd#1174](https://github.com/hdfchain/hdfd/pull/1174))
- txscript: Remove SigHashAllValue ([hdfchain/hdfd#1175](https://github.com/hdfchain/hdfd/pull/1175))
- txscript: Decouple and optimize sighash calc ([hdfchain/hdfd#1179](https://github.com/hdfchain/hdfd/pull/1179))
- wire: Remove TxSerializeWitnessValueSigning ([hdfchain/hdfd#1176](https://github.com/hdfchain/hdfd/pull/1176))
- hdkeychain: Satisfy fmt.Stringer interface ([hdfchain/hdfd#1168](https://github.com/hdfchain/hdfd/pull/1168))
- blockchain: Validate tx expiry in block context ([hdfchain/hdfd#1187](https://github.com/hdfchain/hdfd/pull/1187))
- blockchain: rename ErrRegTxSpendStakeOut to ErrRegTxCreateStakeOut ([hdfchain/hdfd#1195](https://github.com/hdfchain/hdfd/pull/1195))
- multi: Break coinbase dep on standardness rules ([hdfchain/hdfd#1196](https://github.com/hdfchain/hdfd/pull/1196))
- txscript: Cleanup code for the substr opcode ([hdfchain/hdfd#1206](https://github.com/hdfchain/hdfd/pull/1206))
- multi: use secp256k1 types and fields directly ([hdfchain/hdfd#1211](https://github.com/hdfchain/hdfd/pull/1211))
- dcrec: add Pubkey func to secp256k1 and edwards elliptic curves ([hdfchain/hdfd#1214](https://github.com/hdfchain/hdfd/pull/1214))
- blockchain: use secp256k1 functions directly ([hdfchain/hdfd#1212](https://github.com/hdfchain/hdfd/pull/1212))
- multi: Replace btclog with slog ([hdfchain/hdfd#1216](https://github.com/hdfchain/hdfd/pull/1216))
- multi: Define vgo modules ([hdfchain/hdfd#1223](https://github.com/hdfchain/hdfd/pull/1223))
- chainhash: Define vgo module ([hdfchain/hdfd#1224](https://github.com/hdfchain/hdfd/pull/1224))
- wire: Refine vgo deps ([hdfchain/hdfd#1225](https://github.com/hdfchain/hdfd/pull/1225))
- addrmrg: Refine vgo deps ([hdfchain/hdfd#1226](https://github.com/hdfchain/hdfd/pull/1226))
- chaincfg: Refine vgo deps ([hdfchain/hdfd#1227](https://github.com/hdfchain/hdfd/pull/1227))
- multi: Return fork len from ProcessBlock ([hdfchain/hdfd#1233](https://github.com/hdfchain/hdfd/pull/1233))
- blockchain: Panic on fatal assertions ([hdfchain/hdfd#1243](https://github.com/hdfchain/hdfd/pull/1243))
- blockchain: Convert to full block index in mem ([hdfchain/hdfd#1229](https://github.com/hdfchain/hdfd/pull/1229))
- blockchain: Optimize checkpoint handling ([hdfchain/hdfd#1230](https://github.com/hdfchain/hdfd/pull/1230))
- blockchain: Optimize block locator generation ([hdfchain/hdfd#1237](https://github.com/hdfchain/hdfd/pull/1237))
- multi: Refactor and optimize inv discovery ([hdfchain/hdfd#1239](https://github.com/hdfchain/hdfd/pull/1239))
- peer: Minor function definition order cleanup ([hdfchain/hdfd#1247](https://github.com/hdfchain/hdfd/pull/1247))
- peer: Remove superfluous dup version check ([hdfchain/hdfd#1248](https://github.com/hdfchain/hdfd/pull/1248))
- txscript: export canonicalDataSize ([hdfchain/hdfd#1266](https://github.com/hdfchain/hdfd/pull/1266))
- blockchain: Add BuildMerkleTreeStore alternative for MsgTx ([hdfchain/hdfd#1268](https://github.com/hdfchain/hdfd/pull/1268))
- blockchain: Optimize exported header access ([hdfchain/hdfd#1273](https://github.com/hdfchain/hdfd/pull/1273))
- txscript: Cleanup P2SH and stake opcode handling ([hdfchain/hdfd#1318](https://github.com/hdfchain/hdfd/pull/1318))
- txscript: Significantly improve errors ([hdfchain/hdfd#1319](https://github.com/hdfchain/hdfd/pull/1319))
- txscript: Remove pay-to-script-hash flag ([hdfchain/hdfd#1321](https://github.com/hdfchain/hdfd/pull/1321))
- txscript: Remove DER signature verification flag ([hdfchain/hdfd#1323](https://github.com/hdfchain/hdfd/pull/1323))
- txscript: Remove verify minimal data flag ([hdfchain/hdfd#1326](https://github.com/hdfchain/hdfd/pull/1326))
- txscript: Remove script num require minimal flag ([hdfchain/hdfd#1328](https://github.com/hdfchain/hdfd/pull/1328))
- txscript: Make PeekInt consistent with PopInt ([hdfchain/hdfd#1329](https://github.com/hdfchain/hdfd/pull/1329))
- build: Add experimental support for vgo ([hdfchain/hdfd#1215](https://github.com/hdfchain/hdfd/pull/1215))
- build: Update some vgo dependencies to use tags ([hdfchain/hdfd#1219](https://github.com/hdfchain/hdfd/pull/1219))
- stake: add ExpiredByBlock to stake.Node ([hdfchain/hdfd#1221](https://github.com/hdfchain/hdfd/pull/1221))
- server: Minor function definition order cleanup ([hdfchain/hdfd#1271](https://github.com/hdfchain/hdfd/pull/1271))
- server: Convert CF code to use new inv discovery ([hdfchain/hdfd#1272](https://github.com/hdfchain/hdfd/pull/1272))
- multi: add valueIn parameter to wire.NewTxIn ([hdfchain/hdfd#1287](https://github.com/hdfchain/hdfd/pull/1287))
- txscript: Remove low S verification flag ([hdfchain/hdfd#1308](https://github.com/hdfchain/hdfd/pull/1308))
- txscript: Remove unused old sig hash type ([hdfchain/hdfd#1309](https://github.com/hdfchain/hdfd/pull/1309))
- txscript: Remove strict encoding verification flag ([hdfchain/hdfd#1310](https://github.com/hdfchain/hdfd/pull/1310))
- blockchain: Combine block by hash functions ([hdfchain/hdfd#1330](https://github.com/hdfchain/hdfd/pull/1330))
- multi: Continue conversion from chainec to dcrec ([hdfchain/hdfd#1304](https://github.com/hdfchain/hdfd/pull/1304))
- multi: Remove unused secp256k1 sig parse parameter ([hdfchain/hdfd#1335](https://github.com/hdfchain/hdfd/pull/1335))
- blockchain: Refactor db main chain idx to blk idx ([hdfchain/hdfd#1332](https://github.com/hdfchain/hdfd/pull/1332))
- blockchain: Remove main chain index from db ([hdfchain/hdfd#1334](https://github.com/hdfchain/hdfd/pull/1334))
- blockchain: Implement new chain view ([hdfchain/hdfd#1337](https://github.com/hdfchain/hdfd/pull/1337))
- blockmanager: remove unused Pause() API ([hdfchain/hdfd#1340](https://github.com/hdfchain/hdfd/pull/1340))
- chainhash: Remove dup code from hash funcs ([hdfchain/hdfd#1342](https://github.com/hdfchain/hdfd/pull/1342))
- connmgr: Fix the ConnReq print out causing panic ([hdfchain/hdfd#1345](https://github.com/hdfchain/hdfd/pull/1345))
- gcs: Pool MatchAny data allocations ([hdfchain/hdfd#1348](https://github.com/hdfchain/hdfd/pull/1348))
- blockchain: Faster chain view block locator ([hdfchain/hdfd#1338](https://github.com/hdfchain/hdfd/pull/1338))
- blockchain: Refactor to use new chain view ([hdfchain/hdfd#1344](https://github.com/hdfchain/hdfd/pull/1344))
- blockchain: Remove unnecessary genesis block check ([hdfchain/hdfd#1368](https://github.com/hdfchain/hdfd/pull/1368))
- chainhash: Update go build module support ([hdfchain/hdfd#1358](https://github.com/hdfchain/hdfd/pull/1358))
- wire: Update go build module support ([hdfchain/hdfd#1359](https://github.com/hdfchain/hdfd/pull/1359))
- addrmgr: Update go build module support ([hdfchain/hdfd#1360](https://github.com/hdfchain/hdfd/pull/1360))
- chaincfg: Update go build module support ([hdfchain/hdfd#1361](https://github.com/hdfchain/hdfd/pull/1361))
- connmgr: Refine go build module support ([hdfchain/hdfd#1363](https://github.com/hdfchain/hdfd/pull/1363))
- secp256k1: Refine go build module support ([hdfchain/hdfd#1362](https://github.com/hdfchain/hdfd/pull/1362))
- dcrec: Refine go build module support ([hdfchain/hdfd#1364](https://github.com/hdfchain/hdfd/pull/1364))
- certgen: Update go build module support ([hdfchain/hdfd#1365](https://github.com/hdfchain/hdfd/pull/1365))
- dcrutil: Refine go build module support ([hdfchain/hdfd#1366](https://github.com/hdfchain/hdfd/pull/1366))
- hdkeychain: Refine go build module support ([hdfchain/hdfd#1369](https://github.com/hdfchain/hdfd/pull/1369))
- txscript: Refine go build module support ([hdfchain/hdfd#1370](https://github.com/hdfchain/hdfd/pull/1370))
- multi: Remove go modules that do not build ([hdfchain/hdfd#1371](https://github.com/hdfchain/hdfd/pull/1371))
- database: Refine go build module support ([hdfchain/hdfd#1372](https://github.com/hdfchain/hdfd/pull/1372))
- build: Refine build module support ([hdfchain/hdfd#1384](https://github.com/hdfchain/hdfd/pull/1384))
- blockmanager: make pruning transactions consistent ([hdfchain/hdfd#1376](https://github.com/hdfchain/hdfd/pull/1376))
- blockchain: Optimize reorg to use known status ([hdfchain/hdfd#1367](https://github.com/hdfchain/hdfd/pull/1367))
- blockchain: Make block index flushable ([hdfchain/hdfd#1375](https://github.com/hdfchain/hdfd/pull/1375))
- blockchain: Mark fastadd block valid ([hdfchain/hdfd#1392](https://github.com/hdfchain/hdfd/pull/1392))
- release: Bump module versions and deps ([hdfchain/hdfd#1390](https://github.com/hdfchain/hdfd/pull/1390))
- blockchain: Mark fastadd block valid ([hdfchain/hdfd#1392](https://github.com/hdfchain/hdfd/pull/1392))
- gcs: use dchest/siphash ([hdfchain/hdfd#1395](https://github.com/hdfchain/hdfd/pull/1395))
- dcrec: Make function defs more consistent ([hdfchain/hdfd#1432](https://github.com/hdfchain/hdfd/pull/1432))

### Testing and Quality Assurance:

- addrmgr: Simplify tests for KnownAddress ([hdfchain/hdfd#1133](https://github.com/hdfchain/hdfd/pull/1133))
- blockchain: move block validation rule tests into fullblocktests ([hdfchain/hdfd#1141](https://github.com/hdfchain/hdfd/pull/1141))
- addrmgr: Test timestamp update during AddAddress ([hdfchain/hdfd#1137](https://github.com/hdfchain/hdfd/pull/1137))
- txscript: Consolidate tests into txscript package ([hdfchain/hdfd#1177](https://github.com/hdfchain/hdfd/pull/1177))
- txscript: Add JSON-based signature hash tests ([hdfchain/hdfd#1178](https://github.com/hdfchain/hdfd/pull/1178))
- txscript: Correct JSON-based signature hash tests ([hdfchain/hdfd#1181](https://github.com/hdfchain/hdfd/pull/1181))
- txscript: Add benchmark for sighash calculation ([hdfchain/hdfd#1179](https://github.com/hdfchain/hdfd/pull/1179))
- mempool: Refactor pool membership test logic ([hdfchain/hdfd#1188](https://github.com/hdfchain/hdfd/pull/1188))
- blockchain: utilize CalcNextReqStakeDifficulty in fullblocktests ([hdfchain/hdfd#1189](https://github.com/hdfchain/hdfd/pull/1189))
- fullblocktests: add additional premine and malformed tests ([hdfchain/hdfd#1190](https://github.com/hdfchain/hdfd/pull/1190))
- txscript: Improve substr opcode test coverage ([hdfchain/hdfd#1205](https://github.com/hdfchain/hdfd/pull/1205))
- txscript: Convert reference tests to new format ([hdfchain/hdfd#1320](https://github.com/hdfchain/hdfd/pull/1320))
- txscript: Remove P2SH flag from test data ([hdfchain/hdfd#1322](https://github.com/hdfchain/hdfd/pull/1322))
- txscript: Remove DERSIG flag from test data ([hdfchain/hdfd#1324](https://github.com/hdfchain/hdfd/pull/1324))
- txscript: Remove MINIMALDATA flag from test data ([hdfchain/hdfd#1327](https://github.com/hdfchain/hdfd/pull/1327))
- fullblocktests: Add expired stake tx test ([hdfchain/hdfd#1184](https://github.com/hdfchain/hdfd/pull/1184))
- travis: simplify Docker files ([hdfchain/hdfd#1275](https://github.com/hdfchain/hdfd/pull/1275))
- docker: Add dockerfiles for running hdfd nodes ([hdfchain/hdfd#1317](https://github.com/hdfchain/hdfd/pull/1317))
- blockchain: Improve spend journal tests ([hdfchain/hdfd#1246](https://github.com/hdfchain/hdfd/pull/1246))
- txscript: Cleanup and add tests for left opcode ([hdfchain/hdfd#1281](https://github.com/hdfchain/hdfd/pull/1281))
- txscript: Cleanup and add tests for right opcode ([hdfchain/hdfd#1282](https://github.com/hdfchain/hdfd/pull/1282))
- txscript: Cleanup and add tests for the cat opcode ([hdfchain/hdfd#1283](https://github.com/hdfchain/hdfd/pull/1283))
- txscript: Cleanup and add tests for rotr opcode ([hdfchain/hdfd#1285](https://github.com/hdfchain/hdfd/pull/1285))
- txscript: Cleanup and add tests for rotl opcode ([hdfchain/hdfd#1286](https://github.com/hdfchain/hdfd/pull/1286))
- txscript: Cleanup and add tests for lshift opcode ([hdfchain/hdfd#1288](https://github.com/hdfchain/hdfd/pull/1288))
- txscript: Cleanup and add tests for rshift opcode ([hdfchain/hdfd#1289](https://github.com/hdfchain/hdfd/pull/1289))
- txscript: Cleanup and add tests for div opcode ([hdfchain/hdfd#1290](https://github.com/hdfchain/hdfd/pull/1290))
- txscript: Cleanup and add tests for mod opcode ([hdfchain/hdfd#1291](https://github.com/hdfchain/hdfd/pull/1291))
- txscript: Update CSV to match tests in DCP0003 ([hdfchain/hdfd#1292](https://github.com/hdfchain/hdfd/pull/1292))
- txscript: Introduce repeated syntax to test data ([hdfchain/hdfd#1299](https://github.com/hdfchain/hdfd/pull/1299))
- txscript: Allow multi opcode test data repeat ([hdfchain/hdfd#1300](https://github.com/hdfchain/hdfd/pull/1300))
- txscript: Improve and correct some script tests ([hdfchain/hdfd#1303](https://github.com/hdfchain/hdfd/pull/1303))
- main: verify network pow limits ([hdfchain/hdfd#1302](https://github.com/hdfchain/hdfd/pull/1302))
- txscript: Remove STRICTENC flag from test data ([hdfchain/hdfd#1311](https://github.com/hdfchain/hdfd/pull/1311))
- txscript: Cleanup plus tests for checksig opcodes ([hdfchain/hdfd#1315](https://github.com/hdfchain/hdfd/pull/1315))
- blockchain: Add negative tests for forced reorg ([hdfchain/hdfd#1341](https://github.com/hdfchain/hdfd/pull/1341))
- dcrjson: Consolidate tests into dcrjson package ([hdfchain/hdfd#1373](https://github.com/hdfchain/hdfd/pull/1373))
- txscript: add additional data push op code tests ([hdfchain/hdfd#1346](https://github.com/hdfchain/hdfd/pull/1346))
- txscript: add/group control op code tests ([hdfchain/hdfd#1349](https://github.com/hdfchain/hdfd/pull/1349))
- txscript: add/group stack op code tests ([hdfchain/hdfd#1350](https://github.com/hdfchain/hdfd/pull/1350))
- txscript: group splice opcode tests ([hdfchain/hdfd#1351](https://github.com/hdfchain/hdfd/pull/1351))
- txscript: add/group bitwise logic, comparison & rotation op code tests ([hdfchain/hdfd#1352](https://github.com/hdfchain/hdfd/pull/1352))
- txscript: add/group numeric related opcode tests ([hdfchain/hdfd#1353](https://github.com/hdfchain/hdfd/pull/1353))
- txscript: group reserved op code tests ([hdfchain/hdfd#1355](https://github.com/hdfchain/hdfd/pull/1355))
- txscript: add/group crypto related op code tests ([hdfchain/hdfd#1354](https://github.com/hdfchain/hdfd/pull/1354))
- multi: Reduce testnet2 refs in unit tests ([hdfchain/hdfd#1387](https://github.com/hdfchain/hdfd/pull/1387))
- blockchain: Avoid deployment expiration in tests ([hdfchain/hdfd#1450](https://github.com/hdfchain/hdfd/pull/1450))

### Misc:

- release: Bump for v1.3.0 ([hdfchain/hdfd#1388](https://github.com/hdfchain/hdfd/pull/1388))
- multi: Correct typos found by misspell ([hdfchain/hdfd#1197](https://github.com/hdfchain/hdfd/pull/1197))
- main: Correct mem profile error message ([hdfchain/hdfd#1183](https://github.com/hdfchain/hdfd/pull/1183))
- multi: Use saner permissions saving certs ([hdfchain/hdfd#1263](https://github.com/hdfchain/hdfd/pull/1263))
- server: only call time.Now() once ([hdfchain/hdfd#1313](https://github.com/hdfchain/hdfd/pull/1313))
- multi: linter cleanup ([hdfchain/hdfd#1305](https://github.com/hdfchain/hdfd/pull/1305))
- multi: Remove unnecessary network name funcs ([hdfchain/hdfd#1387](https://github.com/hdfchain/hdfd/pull/1387))
- config: Warn if testnet2 database exists ([hdfchain/hdfd#1389](https://github.com/hdfchain/hdfd/pull/1389))

### Code Contributors (alphabetical order):

- Dave Collins
- David Hill
- Dmitry Fedorov
- Donald Adu-Poku
- harzo
- hypernoob
- J Fixby
- Jonathan Chappelow
- Josh Rickmar
- Markus Richter
- matadormel
- Matheus Degiovani
- Michael Eze
- Orthomind
- Shuai Qi
- Tibor Bősze
- Victor Oliveira
