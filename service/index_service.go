package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"redis_queue/structs"
)

type IndexService struct {}

func (i *IndexService) LPushRedis(topic string, val *structs.RedisQueue) {
	bt_val, _ := json.Marshal(val)
	redis_tool := new(RedisTool)
	topic_redis_key := beego.AppConfig.String(strings.ToUpper(topic) + "_REDIS")
	err := redis_tool.LPush(topic_redis_key, string(bt_val))
	if err != nil {
		Log("push_redis_err", "topic:" + topic + " message:" + err.Error() + " params:" + fmt.Sprint(val))
		fmt.Println(err.Error())
	}
}