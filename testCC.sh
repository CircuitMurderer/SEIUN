#!/bin/bash -eu
export CHAINCODE_ID=basic_1:03a3033b7f316aa0846b1e8c7bfce88be4d39fc0f61a0abd9b92c5c03ac55673

source envPeerCompany.sh
peer lifecycle chaincode approveformyorg -o orderer1.council.seiun.net:7051 --tls --cafile $ORDERER_CA --channelID testchannel --name basic --version 1.0 --sequence 1 --waitForEvent --init-required --package-id $CHAINCODE_ID --signature-policy "OR('groupMSP.peer')"
peer lifecycle chaincode queryapproved -C testchannel -n basic --sequence 1
source envPeerSchool.sh
peer lifecycle chaincode approveformyorg -o orderer3.council.seiun.net:7057 --tls --cafile $ORDERER_CA --channelID testchannel --name basic --version 1.0 --sequence 1 --waitForEvent --init-required --package-id $CHAINCODE_ID --signature-policy "OR('groupMSP.peer')"
peer lifecycle chaincode queryapproved -C testchannel -n basic --sequence 1
source envPeerGroup.sh
peer lifecycle chaincode approveformyorg -o orderer2.council.seiun.net:7054 --tls --cafile $ORDERER_CA --channelID testchannel --name basic --version 1.0 --sequence 1 --waitForEvent --init-required --package-id $CHAINCODE_ID --signature-policy "OR('groupMSP.peer')"
peer lifecycle chaincode queryapproved -C testchannel -n basic --sequence 1

peer lifecycle chaincode commit -o orderer2.council.seiun.net:7054 --tls --cafile $ORDERER_CA --channelID testchannel --name basic --init-required --version 1.0 --sequence 1 --peerAddresses peer1.company.seiun.net:7251 --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE --peerAddresses peer1.school.seiun.net:7351 --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE --peerAddresses peer1.group.seiun.net:7451 --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE --signature-policy "OR('groupMSP.peer')"
sleep 5
peer chaincode invoke --isInit -o orderer1.council.seiun.net:7051 --tls --cafile $ORDERER_CA --channelID testchannel --name basic --peerAddresses peer1.group.seiun.net:7451 --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE -c '{"Args":["InitLedger"]}' # --peerAddresses peer1.school.seiun.net:7351 --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE 
sleep 5
#  peer chaincode invoke -o orderer1.council.seiun.net:7051 --tls --cafile $ORDERER_CA --channelID testchannel --name basic --peerAddresses peer1.company.seiun.net:7251 --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE --peerAddresses peer1.school.seiun.net:7351 --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE -c '{"Args":["GetAllItems"]}'
peer chaincode query -C testchannel -n basic -c '{"Args":["GetAllCerts"]}'
