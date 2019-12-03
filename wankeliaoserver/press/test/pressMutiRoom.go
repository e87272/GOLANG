package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var domainName = "beta-mml-lb.zanxingbctv.com"
var addr = flag.String("addr", domainName, "http service address")

var mutex sync.Mutex

var roomInfoArr [100]map[string]string

func main() {

	roomInfoArr[0] = map[string]string{}
	roomInfoArr[1] = map[string]string{}
	roomInfoArr[2] = map[string]string{}
	roomInfoArr[3] = map[string]string{}
	roomInfoArr[4] = map[string]string{}
	roomInfoArr[5] = map[string]string{}
	roomInfoArr[6] = map[string]string{}
	roomInfoArr[7] = map[string]string{}
	roomInfoArr[8] = map[string]string{}
	roomInfoArr[9] = map[string]string{}
	roomInfoArr[10] = map[string]string{}
	roomInfoArr[11] = map[string]string{}
	roomInfoArr[12] = map[string]string{}
	roomInfoArr[13] = map[string]string{}
	roomInfoArr[14] = map[string]string{}
	roomInfoArr[15] = map[string]string{}
	roomInfoArr[16] = map[string]string{}
	roomInfoArr[17] = map[string]string{}
	roomInfoArr[18] = map[string]string{}
	roomInfoArr[19] = map[string]string{}
	roomInfoArr[20] = map[string]string{}
	roomInfoArr[21] = map[string]string{}
	roomInfoArr[22] = map[string]string{}
	roomInfoArr[23] = map[string]string{}
	roomInfoArr[24] = map[string]string{}
	roomInfoArr[25] = map[string]string{}
	roomInfoArr[26] = map[string]string{}
	roomInfoArr[27] = map[string]string{}
	roomInfoArr[28] = map[string]string{}
	roomInfoArr[29] = map[string]string{}
	roomInfoArr[30] = map[string]string{}
	roomInfoArr[31] = map[string]string{}
	roomInfoArr[32] = map[string]string{}
	roomInfoArr[33] = map[string]string{}
	roomInfoArr[34] = map[string]string{}
	roomInfoArr[35] = map[string]string{}
	roomInfoArr[36] = map[string]string{}
	roomInfoArr[37] = map[string]string{}
	roomInfoArr[38] = map[string]string{}
	roomInfoArr[39] = map[string]string{}
	roomInfoArr[40] = map[string]string{}
	roomInfoArr[41] = map[string]string{}
	roomInfoArr[42] = map[string]string{}
	roomInfoArr[43] = map[string]string{}
	roomInfoArr[44] = map[string]string{}
	roomInfoArr[45] = map[string]string{}
	roomInfoArr[46] = map[string]string{}
	roomInfoArr[47] = map[string]string{}
	roomInfoArr[48] = map[string]string{}
	roomInfoArr[49] = map[string]string{}
	roomInfoArr[50] = map[string]string{}
	roomInfoArr[51] = map[string]string{}
	roomInfoArr[52] = map[string]string{}
	roomInfoArr[53] = map[string]string{}
	roomInfoArr[54] = map[string]string{}
	roomInfoArr[55] = map[string]string{}
	roomInfoArr[56] = map[string]string{}
	roomInfoArr[57] = map[string]string{}
	roomInfoArr[58] = map[string]string{}
	roomInfoArr[59] = map[string]string{}
	roomInfoArr[60] = map[string]string{}
	roomInfoArr[61] = map[string]string{}
	roomInfoArr[62] = map[string]string{}
	roomInfoArr[63] = map[string]string{}
	roomInfoArr[64] = map[string]string{}
	roomInfoArr[65] = map[string]string{}
	roomInfoArr[66] = map[string]string{}
	roomInfoArr[67] = map[string]string{}
	roomInfoArr[68] = map[string]string{}
	roomInfoArr[69] = map[string]string{}
	roomInfoArr[70] = map[string]string{}
	roomInfoArr[71] = map[string]string{}
	roomInfoArr[72] = map[string]string{}
	roomInfoArr[73] = map[string]string{}
	roomInfoArr[74] = map[string]string{}
	roomInfoArr[75] = map[string]string{}
	roomInfoArr[76] = map[string]string{}
	roomInfoArr[77] = map[string]string{}
	roomInfoArr[78] = map[string]string{}
	roomInfoArr[79] = map[string]string{}
	roomInfoArr[80] = map[string]string{}
	roomInfoArr[81] = map[string]string{}
	roomInfoArr[82] = map[string]string{}
	roomInfoArr[83] = map[string]string{}
	roomInfoArr[84] = map[string]string{}
	roomInfoArr[85] = map[string]string{}
	roomInfoArr[86] = map[string]string{}
	roomInfoArr[87] = map[string]string{}
	roomInfoArr[88] = map[string]string{}
	roomInfoArr[89] = map[string]string{}
	roomInfoArr[90] = map[string]string{}
	roomInfoArr[91] = map[string]string{}
	roomInfoArr[92] = map[string]string{}
	roomInfoArr[93] = map[string]string{}
	roomInfoArr[94] = map[string]string{}
	roomInfoArr[95] = map[string]string{}
	roomInfoArr[96] = map[string]string{}
	roomInfoArr[97] = map[string]string{}
	roomInfoArr[98] = map[string]string{}
	roomInfoArr[99] = map[string]string{}

	roomInfoArr[0]["roomType"] = "liveGroup"
	roomInfoArr[1]["roomType"] = "liveGroup"
	roomInfoArr[2]["roomType"] = "liveGroup"
	roomInfoArr[3]["roomType"] = "liveGroup"
	roomInfoArr[4]["roomType"] = "liveGroup"
	roomInfoArr[5]["roomType"] = "liveGroup"
	roomInfoArr[6]["roomType"] = "liveGroup"
	roomInfoArr[7]["roomType"] = "liveGroup"
	roomInfoArr[8]["roomType"] = "liveGroup"
	roomInfoArr[9]["roomType"] = "liveGroup"
	roomInfoArr[10]["roomType"] = "liveGroup"
	roomInfoArr[11]["roomType"] = "liveGroup"
	roomInfoArr[12]["roomType"] = "liveGroup"
	roomInfoArr[13]["roomType"] = "liveGroup"
	roomInfoArr[14]["roomType"] = "liveGroup"
	roomInfoArr[15]["roomType"] = "liveGroup"
	roomInfoArr[16]["roomType"] = "liveGroup"
	roomInfoArr[17]["roomType"] = "liveGroup"
	roomInfoArr[18]["roomType"] = "liveGroup"
	roomInfoArr[19]["roomType"] = "liveGroup"
	roomInfoArr[20]["roomType"] = "liveGroup"
	roomInfoArr[21]["roomType"] = "liveGroup"
	roomInfoArr[22]["roomType"] = "liveGroup"
	roomInfoArr[23]["roomType"] = "liveGroup"
	roomInfoArr[24]["roomType"] = "liveGroup"
	roomInfoArr[25]["roomType"] = "liveGroup"
	roomInfoArr[26]["roomType"] = "liveGroup"
	roomInfoArr[27]["roomType"] = "liveGroup"
	roomInfoArr[28]["roomType"] = "liveGroup"
	roomInfoArr[29]["roomType"] = "liveGroup"
	roomInfoArr[30]["roomType"] = "liveGroup"
	roomInfoArr[31]["roomType"] = "liveGroup"
	roomInfoArr[32]["roomType"] = "liveGroup"
	roomInfoArr[33]["roomType"] = "liveGroup"
	roomInfoArr[34]["roomType"] = "liveGroup"
	roomInfoArr[35]["roomType"] = "liveGroup"
	roomInfoArr[36]["roomType"] = "liveGroup"
	roomInfoArr[37]["roomType"] = "liveGroup"
	roomInfoArr[38]["roomType"] = "liveGroup"
	roomInfoArr[39]["roomType"] = "liveGroup"
	roomInfoArr[40]["roomType"] = "liveGroup"
	roomInfoArr[41]["roomType"] = "liveGroup"
	roomInfoArr[42]["roomType"] = "liveGroup"
	roomInfoArr[43]["roomType"] = "liveGroup"
	roomInfoArr[44]["roomType"] = "liveGroup"
	roomInfoArr[45]["roomType"] = "liveGroup"
	roomInfoArr[46]["roomType"] = "liveGroup"
	roomInfoArr[47]["roomType"] = "liveGroup"
	roomInfoArr[48]["roomType"] = "liveGroup"
	roomInfoArr[49]["roomType"] = "liveGroup"
	roomInfoArr[50]["roomType"] = "liveGroup"
	roomInfoArr[51]["roomType"] = "liveGroup"
	roomInfoArr[52]["roomType"] = "liveGroup"
	roomInfoArr[53]["roomType"] = "liveGroup"
	roomInfoArr[54]["roomType"] = "liveGroup"
	roomInfoArr[55]["roomType"] = "liveGroup"
	roomInfoArr[56]["roomType"] = "liveGroup"
	roomInfoArr[57]["roomType"] = "liveGroup"
	roomInfoArr[58]["roomType"] = "liveGroup"
	roomInfoArr[59]["roomType"] = "liveGroup"
	roomInfoArr[60]["roomType"] = "liveGroup"
	roomInfoArr[61]["roomType"] = "liveGroup"
	roomInfoArr[62]["roomType"] = "liveGroup"
	roomInfoArr[63]["roomType"] = "liveGroup"
	roomInfoArr[64]["roomType"] = "liveGroup"
	roomInfoArr[65]["roomType"] = "liveGroup"
	roomInfoArr[66]["roomType"] = "liveGroup"
	roomInfoArr[67]["roomType"] = "liveGroup"
	roomInfoArr[68]["roomType"] = "liveGroup"
	roomInfoArr[69]["roomType"] = "liveGroup"
	roomInfoArr[70]["roomType"] = "liveGroup"
	roomInfoArr[71]["roomType"] = "liveGroup"
	roomInfoArr[72]["roomType"] = "liveGroup"
	roomInfoArr[73]["roomType"] = "liveGroup"
	roomInfoArr[74]["roomType"] = "liveGroup"
	roomInfoArr[75]["roomType"] = "liveGroup"
	roomInfoArr[76]["roomType"] = "liveGroup"
	roomInfoArr[77]["roomType"] = "liveGroup"
	roomInfoArr[78]["roomType"] = "liveGroup"
	roomInfoArr[79]["roomType"] = "liveGroup"
	roomInfoArr[80]["roomType"] = "liveGroup"
	roomInfoArr[81]["roomType"] = "liveGroup"
	roomInfoArr[82]["roomType"] = "liveGroup"
	roomInfoArr[83]["roomType"] = "liveGroup"
	roomInfoArr[84]["roomType"] = "liveGroup"
	roomInfoArr[85]["roomType"] = "liveGroup"
	roomInfoArr[86]["roomType"] = "liveGroup"
	roomInfoArr[87]["roomType"] = "liveGroup"
	roomInfoArr[88]["roomType"] = "liveGroup"
	roomInfoArr[89]["roomType"] = "liveGroup"
	roomInfoArr[90]["roomType"] = "liveGroup"
	roomInfoArr[91]["roomType"] = "liveGroup"
	roomInfoArr[92]["roomType"] = "liveGroup"
	roomInfoArr[93]["roomType"] = "liveGroup"
	roomInfoArr[94]["roomType"] = "liveGroup"
	roomInfoArr[95]["roomType"] = "liveGroup"
	roomInfoArr[96]["roomType"] = "liveGroup"
	roomInfoArr[97]["roomType"] = "liveGroup"
	roomInfoArr[98]["roomType"] = "liveGroup"
	roomInfoArr[99]["roomType"] = "liveGroup"

	roomInfoArr[0]["roomUuid"] = "000d7cf0e0003000"
	roomInfoArr[1]["roomUuid"] = "000d7cf0e6503000"
	roomInfoArr[2]["roomUuid"] = "000d7cf0ec703000"
	roomInfoArr[3]["roomUuid"] = "000d7cf0f2903000"
	roomInfoArr[4]["roomUuid"] = "000d7cf0f8f03000"
	roomInfoArr[5]["roomUuid"] = "000d7cf0ff303000"
	roomInfoArr[6]["roomUuid"] = "000d7cf105603000"
	roomInfoArr[7]["roomUuid"] = "000d7cf10ba03000"
	roomInfoArr[8]["roomUuid"] = "000d7cf111e03000"
	roomInfoArr[9]["roomUuid"] = "000d7cf118203000"
	roomInfoArr[10]["roomUuid"] = "000d7cf11e603000"
	roomInfoArr[11]["roomUuid"] = "000d7cf124a03000"
	roomInfoArr[12]["roomUuid"] = "000d7cf12ae03000"
	roomInfoArr[13]["roomUuid"] = "000d7cf131103000"
	roomInfoArr[14]["roomUuid"] = "000d7cf137703000"
	roomInfoArr[15]["roomUuid"] = "000d7cf13d603000"
	roomInfoArr[16]["roomUuid"] = "000d7cf143a03000"
	roomInfoArr[17]["roomUuid"] = "000d7cf149e03000"
	roomInfoArr[18]["roomUuid"] = "000d7cf150203000"
	roomInfoArr[19]["roomUuid"] = "000d7cf156703000"
	roomInfoArr[20]["roomUuid"] = "000d7cf15cb03000"
	roomInfoArr[21]["roomUuid"] = "000d7cf162d03000"
	roomInfoArr[22]["roomUuid"] = "000d7cf169303000"
	roomInfoArr[23]["roomUuid"] = "000d7cf16f803000"
	roomInfoArr[24]["roomUuid"] = "000d7cf175e03000"
	roomInfoArr[25]["roomUuid"] = "000d7cf17c403000"
	roomInfoArr[26]["roomUuid"] = "000d7cf182803000"
	roomInfoArr[27]["roomUuid"] = "000d7cf188903000"
	roomInfoArr[28]["roomUuid"] = "000d7cf18ed03000"
	roomInfoArr[29]["roomUuid"] = "000d7cf195103000"
	roomInfoArr[30]["roomUuid"] = "000d7cf19b503000"
	roomInfoArr[31]["roomUuid"] = "000d7cf1a1903000"
	roomInfoArr[32]["roomUuid"] = "000d7cf1a7c03000"
	roomInfoArr[33]["roomUuid"] = "000d7cf1ae203000"
	roomInfoArr[34]["roomUuid"] = "000d7cf1b4603000"
	roomInfoArr[35]["roomUuid"] = "000d7cf1baa03000"
	roomInfoArr[36]["roomUuid"] = "000d7cf1c0c03000"
	roomInfoArr[37]["roomUuid"] = "000d7cf1c7203000"
	roomInfoArr[38]["roomUuid"] = "000d7cf1cd503000"
	roomInfoArr[39]["roomUuid"] = "000d7cf1d3903000"
	roomInfoArr[40]["roomUuid"] = "000d7cf1d9d03000"
	roomInfoArr[41]["roomUuid"] = "000d7cf1e0103000"
	roomInfoArr[42]["roomUuid"] = "000d7cf1e6503000"
	roomInfoArr[43]["roomUuid"] = "000d7cf1ec903000"
	roomInfoArr[44]["roomUuid"] = "000d7cf1f2d03000"
	roomInfoArr[45]["roomUuid"] = "000d7cf1f9203000"
	roomInfoArr[46]["roomUuid"] = "000d7cf1ff603000"
	roomInfoArr[47]["roomUuid"] = "000d7cf205a03000"
	roomInfoArr[48]["roomUuid"] = "000d7cf20bd03000"
	roomInfoArr[49]["roomUuid"] = "000d7cf212103000"
	roomInfoArr[50]["roomUuid"] = "000d7cf218503000"
	roomInfoArr[51]["roomUuid"] = "000d7cf21e903000"
	roomInfoArr[52]["roomUuid"] = "000d7cf224d03000"
	roomInfoArr[53]["roomUuid"] = "000d7cf22b203000"
	roomInfoArr[54]["roomUuid"] = "000d7cf231603000"
	roomInfoArr[55]["roomUuid"] = "000d7cf237a03000"
	roomInfoArr[56]["roomUuid"] = "000d7cf23df03000"
	roomInfoArr[57]["roomUuid"] = "000d7cf244203000"
	roomInfoArr[58]["roomUuid"] = "000d7cf24a403000"
	roomInfoArr[59]["roomUuid"] = "000d7cf250903000"
	roomInfoArr[60]["roomUuid"] = "000d7cf256d03000"
	roomInfoArr[61]["roomUuid"] = "000d7cf25d103000"
	roomInfoArr[62]["roomUuid"] = "000d7cf263503000"
	roomInfoArr[63]["roomUuid"] = "000d7cf269903000"
	roomInfoArr[64]["roomUuid"] = "000d7cf26fe03000"
	roomInfoArr[65]["roomUuid"] = "000d7cf276103000"
	roomInfoArr[66]["roomUuid"] = "000d7cf27c603000"
	roomInfoArr[67]["roomUuid"] = "000d7cf282903000"
	roomInfoArr[68]["roomUuid"] = "000d7cf288e03000"
	roomInfoArr[69]["roomUuid"] = "000d7cf28f103000"
	roomInfoArr[70]["roomUuid"] = "000d7cf295503000"
	roomInfoArr[71]["roomUuid"] = "000d7cf29b903000"
	roomInfoArr[72]["roomUuid"] = "000d7cf2a1d03000"
	roomInfoArr[73]["roomUuid"] = "000d7cf2a8203000"
	roomInfoArr[74]["roomUuid"] = "000d7cf2ae503000"
	roomInfoArr[75]["roomUuid"] = "000d7cf2b4a03000"
	roomInfoArr[76]["roomUuid"] = "000d7cf2bae03000"
	roomInfoArr[77]["roomUuid"] = "000d7cf2c1003000"
	roomInfoArr[78]["roomUuid"] = "000d7cf2c7403000"
	roomInfoArr[79]["roomUuid"] = "000d7cf2cd803000"
	roomInfoArr[80]["roomUuid"] = "000d7cf2d3d03000"
	roomInfoArr[81]["roomUuid"] = "000d7cf2da103000"
	roomInfoArr[82]["roomUuid"] = "000d7cf2e0503000"
	roomInfoArr[83]["roomUuid"] = "000d7cf2e6903000"
	roomInfoArr[84]["roomUuid"] = "000d7cf2ece03000"
	roomInfoArr[85]["roomUuid"] = "000d7cf2f3103000"
	roomInfoArr[86]["roomUuid"] = "000d7cf2f9603000"
	roomInfoArr[87]["roomUuid"] = "000d7cf2ffa03000"
	roomInfoArr[88]["roomUuid"] = "000d7cf305d03000"
	roomInfoArr[89]["roomUuid"] = "000d7cf30c203000"
	roomInfoArr[90]["roomUuid"] = "000d7cf312403000"
	roomInfoArr[91]["roomUuid"] = "000d7cf318a03000"
	roomInfoArr[92]["roomUuid"] = "000d7cf31ec03000"
	roomInfoArr[93]["roomUuid"] = "000d7cf325103000"
	roomInfoArr[94]["roomUuid"] = "000d7cf32b503000"
	roomInfoArr[95]["roomUuid"] = "000d7cf331a03000"
	roomInfoArr[96]["roomUuid"] = "000d7cf337d03000"
	roomInfoArr[97]["roomUuid"] = "000d7cf33e203000"
	roomInfoArr[98]["roomUuid"] = "000d7cf344603000"
	roomInfoArr[99]["roomUuid"] = "000d7cf34a703000"

	roomInfoArr[0]["roomName"] = "test0000"
	roomInfoArr[1]["roomName"] = "test0001"
	roomInfoArr[2]["roomName"] = "test0002"
	roomInfoArr[3]["roomName"] = "test0003"
	roomInfoArr[4]["roomName"] = "test0004"
	roomInfoArr[5]["roomName"] = "test0005"
	roomInfoArr[6]["roomName"] = "test0006"
	roomInfoArr[7]["roomName"] = "test0007"
	roomInfoArr[8]["roomName"] = "test0008"
	roomInfoArr[9]["roomName"] = "test0009"
	roomInfoArr[10]["roomName"] = "test0010"
	roomInfoArr[11]["roomName"] = "test0011"
	roomInfoArr[12]["roomName"] = "test0012"
	roomInfoArr[13]["roomName"] = "test0013"
	roomInfoArr[14]["roomName"] = "test0014"
	roomInfoArr[15]["roomName"] = "test0015"
	roomInfoArr[16]["roomName"] = "test0016"
	roomInfoArr[17]["roomName"] = "test0017"
	roomInfoArr[18]["roomName"] = "test0018"
	roomInfoArr[19]["roomName"] = "test0019"
	roomInfoArr[20]["roomName"] = "test0020"
	roomInfoArr[21]["roomName"] = "test0021"
	roomInfoArr[22]["roomName"] = "test0022"
	roomInfoArr[23]["roomName"] = "test0023"
	roomInfoArr[24]["roomName"] = "test0024"
	roomInfoArr[25]["roomName"] = "test0025"
	roomInfoArr[26]["roomName"] = "test0026"
	roomInfoArr[27]["roomName"] = "test0027"
	roomInfoArr[28]["roomName"] = "test0028"
	roomInfoArr[29]["roomName"] = "test0029"
	roomInfoArr[30]["roomName"] = "test0030"
	roomInfoArr[31]["roomName"] = "test0031"
	roomInfoArr[32]["roomName"] = "test0032"
	roomInfoArr[33]["roomName"] = "test0033"
	roomInfoArr[34]["roomName"] = "test0034"
	roomInfoArr[35]["roomName"] = "test0035"
	roomInfoArr[36]["roomName"] = "test0036"
	roomInfoArr[37]["roomName"] = "test0037"
	roomInfoArr[38]["roomName"] = "test0038"
	roomInfoArr[39]["roomName"] = "test0039"
	roomInfoArr[40]["roomName"] = "test0040"
	roomInfoArr[41]["roomName"] = "test0041"
	roomInfoArr[42]["roomName"] = "test0042"
	roomInfoArr[43]["roomName"] = "test0043"
	roomInfoArr[44]["roomName"] = "test0044"
	roomInfoArr[45]["roomName"] = "test0045"
	roomInfoArr[46]["roomName"] = "test0046"
	roomInfoArr[47]["roomName"] = "test0047"
	roomInfoArr[48]["roomName"] = "test0048"
	roomInfoArr[49]["roomName"] = "test0049"
	roomInfoArr[50]["roomName"] = "test0050"
	roomInfoArr[51]["roomName"] = "test0051"
	roomInfoArr[52]["roomName"] = "test0052"
	roomInfoArr[53]["roomName"] = "test0053"
	roomInfoArr[54]["roomName"] = "test0054"
	roomInfoArr[55]["roomName"] = "test0055"
	roomInfoArr[56]["roomName"] = "test0056"
	roomInfoArr[57]["roomName"] = "test0057"
	roomInfoArr[58]["roomName"] = "test0058"
	roomInfoArr[59]["roomName"] = "test0059"
	roomInfoArr[60]["roomName"] = "test0060"
	roomInfoArr[61]["roomName"] = "test0061"
	roomInfoArr[62]["roomName"] = "test0062"
	roomInfoArr[63]["roomName"] = "test0063"
	roomInfoArr[64]["roomName"] = "test0064"
	roomInfoArr[65]["roomName"] = "test0065"
	roomInfoArr[66]["roomName"] = "test0066"
	roomInfoArr[67]["roomName"] = "test0067"
	roomInfoArr[68]["roomName"] = "test0068"
	roomInfoArr[69]["roomName"] = "test0069"
	roomInfoArr[70]["roomName"] = "test0070"
	roomInfoArr[71]["roomName"] = "test0071"
	roomInfoArr[72]["roomName"] = "test0072"
	roomInfoArr[73]["roomName"] = "test0073"
	roomInfoArr[74]["roomName"] = "test0074"
	roomInfoArr[75]["roomName"] = "test0075"
	roomInfoArr[76]["roomName"] = "test0076"
	roomInfoArr[77]["roomName"] = "test0077"
	roomInfoArr[78]["roomName"] = "test0078"
	roomInfoArr[79]["roomName"] = "test0079"
	roomInfoArr[80]["roomName"] = "test0080"
	roomInfoArr[81]["roomName"] = "test0081"
	roomInfoArr[82]["roomName"] = "test0082"
	roomInfoArr[83]["roomName"] = "test0083"
	roomInfoArr[84]["roomName"] = "test0084"
	roomInfoArr[85]["roomName"] = "test0085"
	roomInfoArr[86]["roomName"] = "test0086"
	roomInfoArr[87]["roomName"] = "test0087"
	roomInfoArr[88]["roomName"] = "test0088"
	roomInfoArr[89]["roomName"] = "test0089"
	roomInfoArr[90]["roomName"] = "test0090"
	roomInfoArr[91]["roomName"] = "test0091"
	roomInfoArr[92]["roomName"] = "test0092"
	roomInfoArr[93]["roomName"] = "test0093"
	roomInfoArr[94]["roomName"] = "test0094"
	roomInfoArr[95]["roomName"] = "test0095"
	roomInfoArr[96]["roomName"] = "test0096"
	roomInfoArr[97]["roomName"] = "test0097"
	roomInfoArr[98]["roomName"] = "test0098"
	roomInfoArr[99]["roomName"] = "test0099"

	roomInfoArr[0]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[1]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[2]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[3]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[4]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[5]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[6]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[7]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[8]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[9]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[10]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[11]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[12]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[13]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[14]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[15]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[16]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[17]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[18]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[19]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[20]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[21]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[22]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[23]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[24]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[25]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[26]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[27]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[28]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[29]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[30]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[31]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[32]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[33]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[34]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[35]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[36]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[37]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[38]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[39]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[40]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[41]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[42]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[43]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[44]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[45]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[46]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[47]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[48]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[49]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[50]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[51]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[52]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[53]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[54]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[55]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[56]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[57]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[58]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[59]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[60]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[61]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[62]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[63]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[64]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[65]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[66]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[67]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[68]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[69]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[70]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[71]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[72]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[73]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[74]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[75]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[76]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[77]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[78]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[79]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[80]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[81]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[82]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[83]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[84]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[85]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[86]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[87]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[88]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[89]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[90]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[91]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[92]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[93]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[94]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[95]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[96]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[97]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[98]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"
	roomInfoArr[99]["adminSet"] = "{\"0009f37a15d14000\": \"admin\"}"

	for i := 0; i < 10; i++ {
		for j := 0; j < 37; j++ {
			log.Println("main:", i)
			go wsConnect(i, 0, 1)
			time.Sleep(time.Duration(30) * time.Millisecond)
			// go wsConnect(i, i*250+250000, 1000/2)
			// time.Sleep(time.Duration(250) * time.Millisecond)
		}
	}

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for {
		select {
		case t := <-ticker.C:
			log.Println("t:", t)
		}
	}
}
func wsConnect(i int, delayTime int, loopTime int) {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	//log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			receivePacketHandle(c, i)
		}
	}()

	tokenChange(c)

	time.Sleep(time.Duration(delayTime) * time.Millisecond)

	sendChatMessage(c, i)
	ticker := time.NewTicker(time.Second * time.Duration(loopTime))

	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			// log.Println("t:", t)
			//log.Printf("recv roomInfo: %+v", roomInfo)
			sendChatMessage(c, i)
		case <-interrupt:
			log.Println("interrupt")
			select {
			case <-done:
			case c := <-time.After(time.Second):
				log.Println("c", c)
			}
			return
		}
	}
}

func tokenChange(c *websocket.Conn) {

	mutex.Lock()
	msg := map[string]interface{}{}
	msg["cmd"] = "2"
	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	msg["idem"] = timeUnix
	payload := map[string]interface{}{}
	payload["platform"] = "MM"
	payload["platformUuid"] = "5d831e15-4808-5ec7-99f1-3b37-a3022700"
	msg["payload"] = payload

	packetMsg, _ := json.Marshal(msg)

	c.WriteMessage(websocket.TextMessage, packetMsg)
	mutex.Unlock()
}

func enterRoom(c *websocket.Conn, i int) {

	mutex.Lock()
	msg := map[string]interface{}{}
	msg["cmd"] = "10"
	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	msg["idem"] = timeUnix
	msg["payload"] = roomInfoArr[i]

	packetMsg, _ := json.Marshal(msg)

	c.WriteMessage(websocket.TextMessage, packetMsg)
	mutex.Unlock()
}

func sendChatMessage(c *websocket.Conn, i int) {
	mutex.Lock()
	// log.Println("Lock")

	msg := map[string]interface{}{}
	msg["cmd"] = "80"
	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	msg["idem"] = timeUnix
	payload := map[string]interface{}{}
	payload["roomInfo"] = roomInfoArr[i]
	payload["message"] = timeUnix
	payload["style"] = "1"
	msg["payload"] = payload

	packetMsg, err := json.Marshal(msg)
	if err != nil {
		log.Println("sendChatMessage json err:", err)
		return
	}

	err = c.WriteMessage(websocket.TextMessage, packetMsg)
	if err != nil {
		log.Println("sendChatMessage WriteMessage err:", err)
		return
	}

	// log.Println("Unlock")
	mutex.Unlock()
}
func receivePacketHandle(connect *websocket.Conn, i int) {

	_, msg, err := connect.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return
	}
	//log.Printf("recv mtype: %s", mtype)
	//log.Printf("recv msg: %s", msg)

	if err != nil {
		log.Println("receivePacketHandle Readpacket:", err)
	}

	//timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)

	////log.Printf("timeUnix : [%s] ", timeUnix)

	var mapResult map[string]interface{}
	//使用 json.Unmarshal(data []byte, v interface{})进行转换,返回 error 信息
	if err := json.Unmarshal([]byte(msg), &mapResult); err != nil {
		log.Println("receivePacketHandle Unmarshal:", err)
	}

	// log.Printf("mapResult : %+v\n", mapResult)

	switch mapResult["cmd"] {
	case "3":
		enterRoom(connect, i)
		break
	case "11":
		// sendChatMessage(connect, i)
		break
	case "81":
		// log.Printf("mapResult : %+v\n", mapResult)
		break
	}
}
