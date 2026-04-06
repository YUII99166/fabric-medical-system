#!/bin/bash

echo "=== 在所有 peer 节点上安装链码 v1.0.8 ==="

# 检查操作系统类型并设置PATH
if [[ `uname` == 'Darwin' ]]; then
  export PATH=${PWD}/hyperledger-fabric-darwin-amd64-1.4.12/bin:$PATH
elif [[ `uname` == 'Linux' ]]; then
  export PATH=${PWD}/hyperledger-fabric-linux-amd64-1.4.12/bin:$PATH
fi

# 定义所有 peer 节点
declare -a PEERS=(
  "peer0.taobao.com:7051:TaobaoMSP:/etc/hyperledger/peer/taobao.com/users/Admin@taobao.com/msp"
  "peer1.taobao.com:7051:TaobaoMSP:/etc/hyperledger/peer/taobao.com/users/Admin@taobao.com/msp"
  "peer0.jd.com:7051:JDMSP:/etc/hyperledger/peer/jd.com/users/Admin@jd.com/msp"
  "peer1.jd.com:7051:JDMSP:/etc/hyperledger/peer/jd.com/users/Admin@jd.com/msp"
  "peer0.wenjin.com:7051:WenjinMSP:/etc/hyperledger/peer/wenjin.com/users/Admin@wenjin.com/msp"
  "peer1.wenjin.com:7051:WenjinMSP:/etc/hyperledger/peer/wenjin.com/users/Admin@wenjin.com/msp"
  "peer0.regcenter.com:7051:RegCenterMSP:/etc/hyperledger/peer/regcenter.com/users/Admin@regcenter.com/msp"
  "peer1.regcenter.com:7051:RegCenterMSP:/etc/hyperledger/peer/regcenter.com/users/Admin@regcenter.com/msp"
)

echo "1. 清理所有旧的链码容器和镜像..."
docker rm -f $(docker ps -aq --filter name=dev-peer) 2>/dev/null || true
docker rmi -f $(docker images | grep dev-peer | awk '{print $3}') 2>/dev/null || true

echo ""
echo "2. 在所有 peer 节点上安装链码..."
for peer_info in "${PEERS[@]}"; do
  IFS=':' read -r peer_addr peer_port msp_id msp_path <<< "$peer_info"
  
  echo "   安装到 $peer_addr ..."
  docker exec cli bash -c "CORE_PEER_ADDRESS=$peer_addr:$peer_port CORE_PEER_LOCALMSPID=$msp_id CORE_PEER_MSPCONFIGPATH=$msp_path peer chaincode install -n fabric-mims -v 1.0.8 -l golang -p chaincode"
  
  if [ $? -eq 0 ]; then
    echo "   ✓ $peer_addr 安装成功"
  else
    echo "   ✗ $peer_addr 安装失败"
  fi
done

echo ""
echo "3. 在每个 peer 上触发链码容器启动..."
for peer_info in "${PEERS[@]}"; do
  IFS=':' read -r peer_addr peer_port msp_id msp_path <<< "$peer_info"
  
  echo "   触发 $peer_addr ..."
  docker exec cli bash -c "CORE_PEER_ADDRESS=$peer_addr:$peer_port CORE_PEER_LOCALMSPID=$msp_id CORE_PEER_MSPCONFIGPATH=$msp_path peer chaincode invoke -C appchannel -n fabric-mims -c '{\"Args\":[\"hello\"]}'" 2>/dev/null &
done

echo "   等待链码容器构建..."
sleep 30

echo ""
echo "4. 检查链码容器状态..."
docker ps | grep dev-peer

echo ""
echo "=== 安装完成 ==="
echo "请检查上方是否有多个 dev-peer 容器在运行"
