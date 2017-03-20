package services.ztc.publish

import (
	"fmt"
	"strings"
	_ "strconv"

	"models"
)

type pub struct {
	id int // 计划id
	advertiser_id int // 广告主id
	card_id int // 创意id
	adposition_id int // 广告位，1-卡券、2-视频、3-开机大图、4-Banner广告
	status int //状态，1-正常、2-暂停、3-停止
	check int //审核状态，1-等待审核、2-审核不通过、3-审核通过
	billing int //计费模式，1-曝光计费、2-点击计费、3-到达计费
	bid float64 //出价，单位分
	ceiling float64 //日预算，单位分
	
}


