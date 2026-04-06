#!/bin/bash

echo "验证链码升级状态..."

echo ""
echo "1. 检查已安装的链码版本："
docker exec cli peer chaincode list --installed

echo ""
echo "2. 检查已实例化的链码版本："
docker exec cli peer chaincode list --instantiated -C appchannel

echo ""
echo "3. 检查链码容器："
docker ps --filter "name=dev-peer" --format "table {{.Names}}\t{{.Status}}\t{{.Image}}"

echo ""
echo "4. 测试 hello 函数："
docker exec cli bash -c "CORE_PEER_ADDRESS=peer0.taobao.com:7051 CORE_PEER_LOCALMSPID=TaobaoMSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/taobao.com/users/Admin@taobao.com/msp peer chaincode invoke -C appchannel -n fabric-mims -c '{\"Args\":[\"hello\"]}'"

echo ""
echo "5. 测试 queryPrescriptionsByPatient 函数："
docker exec cli bash -c "CORE_PEER_ADDRESS=peer0.taobao.com:7051 CORE_PEER_LOCALMSPID=TaobaoMSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/taobao.com/users/Admin@taobao.com/msp peer chaincode invoke -C appchannel -n fabric-mims -c '{\"Args\":[\"queryPrescriptionsByPatient\",\"patient1\",\"56ab46ae3600\"]}'"

echo ""
echo "验证完成！"
