#!/bin/bash

echo "升级链码到v1.0.8..."

# 检查操作系统类型并设置PATH
if [[ `uname` == 'Darwin' ]]; then
  export PATH=${PWD}/hyperledger-fabric-darwin-amd64-1.4.12/bin:$PATH
elif [[ `uname` == 'Linux' ]]; then
  export PATH=${PWD}/hyperledger-fabric-linux-amd64-1.4.12/bin:$PATH
fi

TaobaoPeer0Cli="CORE_PEER_ADDRESS=peer0.taobao.com:7051 CORE_PEER_LOCALMSPID=TaobaoMSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/taobao.com/users/Admin@taobao.com/msp"
JDPeer0Cli="CORE_PEER_ADDRESS=peer0.jd.com:7051 CORE_PEER_LOCALMSPID=JDMSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/jd.com/users/Admin@jd.com/msp"

echo "1. 清理旧的链码容器和镜像..."
docker rm -f $(docker ps -aq --filter name=dev-peer) 2>/dev/null || true
docker rmi -f $(docker images | grep dev-peer | awk '{print $3}') 2>/dev/null || true

echo "2. 清理peer节点中的旧链码缓存..."
docker exec cli bash -c "rm -rf /var/hyperledger/production/chaincodes/fabric-mims.1.0.7" 2>/dev/null || true
docker exec cli bash -c "rm -rf /var/hyperledger/production/chaincodes/fabric-mims.1.0.8" 2>/dev/null || true
docker exec peer0.taobao.com bash -c "rm -rf /var/hyperledger/production/chaincodes/fabric-mims.1.0.7" 2>/dev/null || true
docker exec peer0.taobao.com bash -c "rm -rf /var/hyperledger/production/chaincodes/fabric-mims.1.0.8" 2>/dev/null || true
docker exec peer0.jd.com bash -c "rm -rf /var/hyperledger/production/chaincodes/fabric-mims.1.0.7" 2>/dev/null || true
docker exec peer0.jd.com bash -c "rm -rf /var/hyperledger/production/chaincodes/fabric-mims.1.0.8" 2>/dev/null || true

echo "3. 安装新版本链码 (v1.0.8)..."
docker exec cli bash -c "$TaobaoPeer0Cli peer chaincode install -n fabric-mims -v 1.0.8 -l golang -p chaincode"
docker exec cli bash -c "$JDPeer0Cli peer chaincode install -n fabric-mims -v 1.0.8 -l golang -p chaincode"

echo "4. 升级链码..."
docker exec cli bash -c "$TaobaoPeer0Cli peer chaincode upgrade -o orderer.qq.com:7050 -C appchannel -n fabric-mims -l golang -v 1.0.8 -c '{\"Args\":[\"init\"]}' -P \"AND ('TaobaoMSP.member','JDMSP.member')\""

echo "等待链码升级完成..."
sleep 10

echo "5. 测试链码..."
docker exec cli bash -c "$TaobaoPeer0Cli peer chaincode invoke -C appchannel -n fabric-mims -c '{\"Args\":[\"hello\"]}'"

echo "链码升级完成！"

