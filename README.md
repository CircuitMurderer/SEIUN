# SEIUN
Safe, Efficient, Immutable Unified Network

---

Usage:

+ Start the Hyperledger Fabric network:

```shell
./setDNS.sh
source ./envPeerCompany.sh
./network.sh -u
./registerUser.sh
./enrollUser.sh
./peerStartUp.sh
./instCC.sh
```

Then modify the `CHAINCODE_ID` of `testCC.sh` to the chaincode ID that `instCC.sh` outputed.

```shell
./testCC.sh
```

+ Start the restAPI:

```shell
cd ./restAPI
go run .
```
