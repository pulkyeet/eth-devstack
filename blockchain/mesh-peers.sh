#!/bin/bash

ENODES=(
  "enode://23ce3880ab80f67a7f0fc28769429f713dd72e26a2f63b95b7bcbb6f33c152a463e0069aef8eee5b47ebd716905745d49757e3be24f29a26155945997148151d@signer1:30303"
  "enode://943295d904c637f075e9026d93d86ee15cdae35d995180161619fa99a4d7c2d700219bb8573dfe23f206d8136cecbbfcc294892a302b64c489808f12320c8fe4@signer2:30303"
  "enode://fbf8735860cf2226526c2fcfaca44937deef9d4ead1adddb6e56dbe971c81ce9ee7ed31d6a4cb18fdae7cf347f0422fa0cb46a73ba27de6ede675e9aade959e6@signer3:30303"
  "enode://6712b8f49722cf06cd71f65fee2f76cd20f6820911df9266ac189da24e0ce89dc5580733fa7a28405c62a45cc9ec76efc9d00318f4d271b3c867bd77b117a071@signer4:30303"
  "enode://e1a4d37bf1f751c507f07bb8a135245d8d71240b9cac57fe1de225f38d774946a5e5f5da1d2a88bfb2a4bee12a88ee72ed82eb534a21f8afa9c63fdf195667ff@signer5:30303"
  "enode://f4f4542cf64b55af1c6f33157ea40dac2b54a9ae40b7878db904dd2bfc6bc8a93d34d4976dbeb5982d3b26566d191d1647e04dd0461a52a651d15bf009f15c5d@signer6:30303"
  "enode://05c45d794a10c6e134c9721c14f4dd8a0c42a5119b1d7e4ddd47e741101a9a72f74461a9d841d21f87c5a7bea637fbb3498df3ed1fe9e0d2422f8912801bac5d@signer7:30303"
)

for i in {1..7}; do
  echo "Adding peers to signer$i..."
  for enode in "${ENODES[@]}"; do
    docker exec eth-signer$i geth attach --exec "admin.addPeer(\"$enode\")" /data/geth.ipc 2>/dev/null
  done
done

echo "Done. Checking peer counts..."
for i in {1..7}; do
  count=$(docker exec eth-signer$i geth attach --exec 'admin.peers.length' /data/geth.ipc 2>/dev/null)
  echo "Signer$i: $count peers"
done
