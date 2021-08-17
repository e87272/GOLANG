package common

import (
	"strconv"
	"strings"
	"sync"
	"time"
)

// 因為snowFlake目的是解決分散式下生成唯一id 所以ID中是包含叢集和節點編號在內的
// |<--                      total 64bits                      -->|
// |<--     time     -->|<--   workerId   -->|<-- serialNumber -->|

const (
	numberBits  uint8 = 4                       // 每毫秒可生成的流水號的bit數
	workerBits  uint8 = 16                      // 每臺機器(節點)的識別碼的bit數
	numberMax   int64 = -1 ^ (-1 << numberBits) // 每毫秒可生成的流水號的最大值
	workerMax   int64 = -1 ^ (-1 << workerBits) // 每臺機器(節點)的識別碼的最大值
	timeShift   uint8 = workerBits + numberBits // 時間戳向左的偏移量
	workerShift uint8 = numberBits              // 節點ID向左的偏移量
	epoch       int64 = 1577836800000           // 起始時間戳(2020/01/01)
)

var timestamp int64
var number int64
var mutex = new(sync.Mutex)

// 生成方法一定要掛載在某個worker下，這樣邏輯會比較清晰 指定某個節點生成id

func GetUuid() string {

	workerIdStr := MyIp()
	ipAry := strings.Split(workerIdStr, ".")
	ip2, _ := strconv.ParseInt(ipAry[2], 10, 64)
	ip3, _ := strconv.ParseInt(ipAry[3], 10, 64)
	workerId := ip2*256 + ip3

	// 獲取id最關鍵的一點 加鎖 加鎖 加鎖
	mutex.Lock()
	defer mutex.Unlock() // 生成完成後記得 解鎖 解鎖 解鎖

	// 獲取生成時的時間戳(毫秒)
	now := time.Now().UnixNano() / int64(time.Millisecond)

	if timestamp == now {
		// 這裡要判斷，當前工作節點是否在1毫秒內已經生成numberMax個ID
		// 如果當前工作節點在1毫秒內生成的ID已經超過上限，則要等到下一毫秒再繼續生成
		if number < numberMax {
			number++
		} else {
			for now <= timestamp {
				now = time.Now().UnixNano() / int64(time.Millisecond)
			}
			// 重置工作節點生成ID的流水號，並將機器上一次生成ID的時間更新為當前時間
			number = 0
			timestamp = now
		}
	} else {
		// 重置工作節點生成ID的流水號，並將機器上一次生成ID的時間更新為當前時間
		number = 0
		timestamp = now
	}

	uuidInt := int64(((now - epoch) << timeShift) | (workerId << workerShift) | number)
	uuid := strconv.FormatInt(uuidInt, 16)
	return strings.Repeat("0", 16-len(uuid)) + uuid
}
