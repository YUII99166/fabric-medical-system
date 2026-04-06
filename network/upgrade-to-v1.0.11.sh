#!/bin/bash

echo "=========================================="
echo "升级链码到 v1.0.11（包含访问溯源功能）"
echo "=========================================="

# 检查操作系统类型并设置PATH
if [[ `uname` == 'Darwin' ]]; then
  export PATH=${PWD}/hyperledger-fabric-darwin-amd64-1.4.12/bin:$PATH
elif [[ `uname` == 'Linux' ]]; then
  export PATH=${PWD}/hyperledger-fabric-linux-amd64-1.4.12/bin:$PATH
fi

TaobaoPeer0Cli="CORE_PEER_ADDRESS=peer0.taobao.com:7051 CORE_PEER_LOCALMSPID=TaobaoMSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/taobao.com/users/Admin@taobao.com/msp"
JDPeer0Cli="CORE_PEER_ADDRESS=peer0.jd.com:7051 CORE_PEER_LOCALMSPID=JDMSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/jd.com/users/Admin@jd.com/msp"
WenjinPeer0Cli="CORE_PEER_ADDRESS=peer0.wenjin.com:7051 CORE_PEER_LOCALMSPID=WenjinMSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/wenjin.com/users/Admin@wenjin.com/msp"
RegCenterPeer0Cli="CORE_PEER_ADDRESS=peer0.regcenter.com:7051 CORE_PEER_LOCALMSPID=RegCenterMSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/regcenter.com/users/Admin@regcenter.com/msp"

echo ""
echo "1. 清理旧的链码容器和镜像..."
echo "=========================================="
docker rm -f $(docker ps -aq --filter name=dev-peer) 2>/dev/null || true
docker rmi -f $(docker images | grep dev-peer | awk '{print $3}') 2>/dev/null || true

echo ""
echo "2. 在所有peer0节点上安装新版本链码 (v1.0.11)..."
echo "=========================================="
docker exec cli bash -c "$TaobaoPeer0Cli peer chaincode install -n fabric-mims -v 1.0.11 -l golang -p chaincode"
docker exec cli bash -c "$JDPeer0Cli peer chaincode install -n fabric-mims -v 1.0.11 -l golang -p chaincode"
docker exec cli bash -c "$WenjinPeer0Cli peer chaincode install -n fabric-mims -v 1.0.11 -l golang -p chaincode"
docker exec cli bash -c "$RegCenterPeer0Cli peer chaincode install -n fabric-mims -v 1.0.11 -l golang -p chaincode"

echo ""
echo "3. 升级链码..."
echo "=========================================="
docker exec cli bash -c "$TaobaoPeer0Cli peer chaincode upgrade -o orderer.qq.com:7050 -C appchannel -n fabric-mims -l golang -v 1.0.11 -c '{\"Args\":[\"init\"]}' -P \"OR ('TaobaoMSP.member','JDMSP.member','WenjinMSP.member','RegCenterMSP.member')\""

echo ""
echo "等待链码升级完成..."
sleep 10

echo ""
echo "4. 测试基本功能..."
echo "=========================================="
docker exec cli bash -c "$TaobaoPeer0Cli peer chaincode invoke -C appchannel -n fabric-mims -c '{\"Args\":[\"hello\"]}'"

echo ""
echo "5. 测试访问记录功能..."
echo "=========================================="
docker exec cli bash -c "$TaobaoPeer0Cli peer chaincode invoke -C appchannel -n fabric-mims -c '{\"Args\":[\"recordPrescriptionAccess\",\"test-log-001\",\"test-prescription-001\",\"PRE-001\",\"patient-001\",\"科比\",\"doctor-001\",\"华佗\",\"doctor\",\"TaobaoMSP\",\"华西医院\",\"view\",\"测试访问\"]}'"

echo ""
echo "=========================================="
echo "✅ 链码升级完成！"
echo "=========================================="
echo ""
echo "新版本功能："
echo "- recordPrescriptionAccess: 记录病历访问"
echo "- queryAccessLogsByPatient: 查询患者的访问日志"
echo ""
