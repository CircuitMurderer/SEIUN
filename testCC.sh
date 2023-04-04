#!/bin/bash -eu
source envPeerCompany.sh
peer lifecycle chaincode package basic.tar.gz --path chaincode --lang golang --label basic_1

source envPeerCompany.sh
peer lifecycle chaincode install basic.tar.gz
peer lifecycle chaincode queryinstalled

source envPeerSchool.sh
peer lifecycle chaincode install basic.tar.gz
peer lifecycle chaincode queryinstalled

source envPeerGroup.sh
peer lifecycle chaincode install basic.tar.gz
peer lifecycle chaincode queryinstalled

export CHAINCODE_ID=basic_1:ec7aea7e807ead8328badfe02885439b9425f2f32229a2f8eaf28b69e2e6e36a

source envPeerCompany.sh
peer lifecycle chaincode approveformyorg -o orderer1.council.seiun.net:7051 --tls --cafile $ORDERER_CA --channelID testchannel --name basic --version 1.0 --sequence 1 --waitForEvent --init-required --package-id $CHAINCODE_ID
peer lifecycle chaincode queryapproved -C testchannel -n basic --sequence 1
source envPeerSchool.sh
peer lifecycle chaincode approveformyorg -o orderer3.council.seiun.net:7057 --tls --cafile $ORDERER_CA --channelID testchannel --name basic --version 1.0 --sequence 1 --waitForEvent --init-required --package-id $CHAINCODE_ID
peer lifecycle chaincode queryapproved -C testchannel -n basic --sequence 1
source envPeerGroup.sh
peer lifecycle chaincode approveformyorg -o orderer2.council.seiun.net:7054 --tls --cafile $ORDERER_CA --channelID testchannel --name basic --version 1.0 --sequence 1 --waitForEvent --init-required --package-id $CHAINCODE_ID
peer lifecycle chaincode queryapproved -C testchannel -n basic --sequence 1

peer lifecycle chaincode commit -o orderer2.council.seiun.net:7054 --tls --cafile $ORDERER_CA --channelID testchannel --name basic --init-required --version 1.0 --sequence 1 --peerAddresses peer1.company.seiun.net:7251 --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE --peerAddresses peer1.school.seiun.net:7351 --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE
sleep 2
peer chaincode invoke --isInit -o orderer1.council.seiun.net:7051 --tls --cafile $ORDERER_CA --channelID testchannel --name basic --peerAddresses peer1.company.seiun.net:7251 --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE --peerAddresses peer1.school.seiun.net:7351 --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE -c '{"Args":["InitLedger"]}'
sleep 2
# peer chaincode invoke -o orderer1.council.seiun.net:7051 --tls --cafile $ORDERER_CA --channelID testchannel --name basic --peerAddresses peer1.company.seiun.net:7251 --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE --peerAddresses peer1.school.seiun.net:7351 --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE -c '{"Args":["GetAllItems"]}'
peer chaincode query -C testchannel -n basic -c '{"Args":["GetAllItems"]}'
