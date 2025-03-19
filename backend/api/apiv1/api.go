package apiv1;

import (
	"log"
	"fmt"
	"sync/atomic"
	"encoding/json"
	"encoding/binary"
	gnet  "github.com/panjf2000/gnet/v2"
);



type ApiV1 struct {
	gnet.BuiltinEventEngine
	totalPlayers uint64
	TotalPackets  uint64
	TotalFailures uint64
	TotalPlayers  uint64
}

func New() *ApiV1 {
	print("\n");
	api := ApiV1{
		totalPlayers: 0,
	};
	return &api;
}



func (api *ApiV1) OnTraffic(conn gnet.Conn) gnet.Action {
	dataIn, _ := conn.Next(-1);
	jsonIn := ProtocolParse(dataIn);
//TODO
numPlayers := 1;
atomic.AddUint64(&api.TotalPackets, 1);
atomic.AddUint64(&api.TotalPlayers, uint64(numPlayers));
fmt.Printf("> %s\n", jsonIn);
	// send reply
	jsonOut, _ := json.Marshal(Result_Submit {
		Uptime: 1,
		Rank:  11,
	});
fmt.Printf("< %s\n\n", jsonOut);
	out := ProtocolEncode(jsonOut);
	fmt.Fprintf(conn, "%s", out);
	return gnet.None;
}



//TODO: add compression and checksums
func ProtocolParse(data []byte) string {
	sizeData := len(data);
	if data[sizeData-1] == 0x0 { sizeData--; }
	if sizeData <= 6 { log.Print("Short packet"); return ""; }
	sizePack := binary.BigEndian.Uint16(data[:2]);
	if uint64(sizeData) != uint64(sizePack + 5) {
		log.Printf("Invalid packet length %d expected: %d", sizeData, sizePack+5); return ""; }
	sumPack := binary.BigEndian.Uint16(data[2:4]);
	jsonPack := data[4 : sizePack+4];
	sizeJson := len(jsonPack);
	if uint64(sizeJson) != uint64(sizePack) {
		log.Printf("Invalid packet actual length %d expected: %d", sizeJson, sizePack); return ""; }
	sumData := ChecksumEncode(jsonPack);
	if sumData != sumPack {
		log.Printf("Invalid packet checksum: %d expected: %d", sumData, sumPack); return ""; }
	return string(jsonPack);
}

func ProtocolEncode(data []byte) []byte {
	size := len(data);
	sizePack := uint16(size + 6);
	sumPack := ChecksumEncode(data);
	out := make([]byte, sizePack);
	copy(out[4:], data[0:size]);
	bb := make([]byte, 2);
	binary.BigEndian.PutUint16(bb, sizePack); copy(out[0:], bb[0:2]); // data len
	binary.BigEndian.PutUint16(bb, sumPack ); copy(out[2:], bb[0:2]); // checksum
	out[sizePack-2] = 0xa; // \n
	out[sizePack-1] = 0x0; // \0
	return out;
}



func ChecksumEncode(data []byte) uint16 {
	return 0;
}
