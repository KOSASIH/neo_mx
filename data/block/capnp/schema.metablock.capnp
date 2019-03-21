@0xddffc3d7f7f36183;
using Go = import "/go.capnp";
$Go.package("capnp");
$Go.import("_");

struct PeerDataCapn {
    publicKey @0: Data;
    action    @1: UInt8;
    timestamp @2: UInt64;
    value     @3: Data;
}

struct ShardMiniBlockHeaderCapn {
   hash            @0: Data;
   receiverShardId @1: UInt32;
   senderShardId   @2: UInt32;
}

struct ShardDataCapn {
    shardId          @0: UInt32;
    headerHash       @1: Data;
    shardMiniBlockHeaders @2: List(ShardMiniBlockHeaderCapn);
}

struct MetaBlockCapn {
    nonce         @0:  UInt64;
    epoch         @1:  UInt32;
    round         @2:  UInt32;
    timeStamp     @3:  UInt64;
    shardInfo     @4:  List(ShardDataCapn);
    peerInfo      @5:  List(PeerDataCapn);
    signature     @6:  Data;
    pubKeysBitmap @7:  Data;
    previousHash  @8:  Data;
    prevRandSeed  @9:  Data;
    randSeed      @10: Data;
    stateRootHash @11: Data;
}

##compile with:

##
##
##   capnpc  -I$GOPATH/src/github.com/glycerine/go-capnproto -ogo $GOPATH/src/github.com/ElrondNetwork/elrond-go-sandbox/data/block/capnp/schema.metablock.capnp