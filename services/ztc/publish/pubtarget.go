package publish

import (
	"fmt"
	"log"
	"time"
	"strconv"
	"strings"
	"github.com/garyburd/redigo/redis"

	"config"
	"models"
	advService "services/advertiser"
)

type PubTarget struct {
	pub  *Pub
	date string
}

func NewPubTarget(pub interface{}, date string) (pt *PubTarget) {
	var p *Pub
	switch v := pub.(type) {
	case map[string]string:
		p = new(Pub)
		p.setPub(v)
		p.setPubversion()
		p.setProduct()
	case *Pub:
		p = v
	}
	if len(date) == 0 {
		date = time.Now().Format("2006-01-02")
	}
	pt = &PubTarget{pub: p, date: date}
	return
}

func (pt *PubTarget) Detect() (ret bool, reports []string) {
	ret = true
	// 检查日期
	date, _  := time.Parse("2006-01-02", pt.date)
	if date.After(pt.pub.etime) || date.Before(pt.pub.stime) {
		ret = false
		msg := fmt.Sprintf("Time interval: %s ~ %s, args: %s", pt.pub.stime.Format("2006-01-02"), pt.pub.etime.Format("2006-01-02"), date);
		reports = append(reports, msg)
	}

	// 检查小时
	var week int
	if weekday := date.Weekday(); weekday == time.Sunday {
		week = 7
	} else {
		week = int(weekday)
	}
	isCurHourValid := false
	if len(pt.pub.times) == 0 {
		isCurHourValid = true
	}
	hours, ok := pt.pub.times[strconv.Itoa(week)]
	if ok {
		hour := time.Now().Hour()
		for _, v := range hours {
			if strconv.Itoa(hour) == strings.TrimLeft(v, "0") {
				isCurHourValid = true
			}
		}
	}
	if isCurHourValid == false {
		ret = false
		reports = append(reports, fmt.Sprintf("Hour points: %v", pt.pub.times))
	}

	// 计划状态
	if pt.pub.status != 1 {
		ret = false
		switch pt.pub.status {
		case 2:
			reports = append(reports, "Publish status: 2(pause)")
		case 3:
			reports = append(reports, "Publish status: 3(suspend)")
		default:
			reports = append(reports, fmt.Sprintf("Publish status: %d(unknonw)", pt.pub.status))
		}
	}
	
	// 计划是否被删除
	if !pt.pub.deletedAt.IsZero() {
		ret = false
		reports = append(reports, "Publish is deleted")
	}
	
	// 验证创意审核状态
	if pt.pub.cardStatus != 1 {
		ret = false
		switch pt.pub.cardStatus {
		case 2, 5:
			reports = append(reports, "Card waiting for check")
		case 4, 6:
			reports = append(reports, "Card check not passed")
		default:
			reports = append(reports, fmt.Sprintf("Card status is unknown(%d)", pt.pub.cardStatus))
		}
	}
	
	// 验证创意url、logo合法性
	displayWay := pt.pub.DisplayWay()
	// 视频创意不验证logo
	if displayWay != DISPLAY_VIDEO {
		logo, ok := pt.pub.card["logo"]
		if !ok {
			ret = false
			reports = append(reports, "Card logo is not existed.")
		} else {
			if logoStr, ok := logo.(string); ok {
				logoStr = strings.ToLower(logoStr)
				if !strings.HasPrefix(logoStr, "http://") && !strings.HasPrefix(logoStr, "https://") {
					ret = false
					reports = append(reports, "Card logo is valid.")
				}
			}
		}
	}
	// 视频和开机大图不验证url
	if displayWay != DISPLAY_VIDEO || displayWay != DISPLAY_KJDT {
		url, ok := pt.pub.card["url"]
		if !ok {
			ret = false
			reports = append(reports, "Card url is not existed.")
		} else {
			if urlStr, ok := url.(string); ok {
				urlStr = strings.ToLower(urlStr)
				if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
					ret = false
					reports = append(reports, "Card url is valid.")
				}
			}
		}
	}

	// 判定广告主是否被屏蔽
	shields := models.Shields().Find(3, []string{"1","2"}, pt.pub.advertiserId, []string{"id"})
	if len(shields) > 0 {
		ret = false
		reports = append(reports, "Advertiser: " + strconv.Itoa(pt.pub.advertiserId) + " is in shields table")
	}

	// 判断广告主资质
	if !advService.New(pt.pub.advertiserId).HasSyncCondition() {
		ret = false
		reports = append(reports, "Advertiser: " + strconv.Itoa(pt.pub.advertiserId) + " qualification is not enough.")
	}

	c, err := redis.Dial("tcp", config.Get("REDIS", "127.0.0.1:6379"))
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	
	// 判断广告主是否还有余额(测试广告主不检测广告主余额)
	if pt.pub.adverType != 3 {
		advkey := "budget:adv:" + strconv.Itoa(pt.pub.advertiserId)
		ret, err := c.Do("HGETALL", advkey)
		if err != nil {
			log.Fatal(err)
		}
		if r, ok := ret.([]interface{}); ok {
			budget := make(map[string]string, len(r) / 2)
			for i := 0; i < len(r); i = i + 2 {
				vk, _ := r[i].([]byte)
				vv, _ := r[i + 1].([]byte)
				budget[string(vk)] = string(vv)
			}
			amount := 0.0
			if v, isPresent := budget["amount"]; isPresent {
				amount, _ = strconv.ParseFloat(v, 64)
			} else {
				ret = false
				reports = append(reports, advkey + " is not exist, or field amount not exist.")
			}
			consume := 0.0
			if v, isPresent := budget["consume"]; isPresent {
				consume, _ = strconv.ParseFloat(v, 64)
			}
			if consume >= amount {
				ret = false
				reports = append(reports, "Advertiser: " + strconv.Itoa(pt.pub.advertiserId) + " amount not enough.")
			}
		}
	}

	// 判断预算是否花超
	if time.Now().Format("2006-01-02") == pt.date {
		pubkey := "budget:pub:" + strconv.Itoa(pt.pub.id) + ":" + time.Now().Format("20060102")
		ret, err := c.Do("HGETALL", pubkey)
		if err != nil {
			log.Fatal(err)
		}
		if r, ok := ret.([]interface{}); ok {
			budget := make(map[string]string, len(r) / 2)
			for i := 0; i < len(r); i = i + 2 {
				vk, _ := r[i].([]byte)
				vv, _ := r[i + 1].([]byte)
				budget[string(vk)] = string(vv)
			}
			ceiling := 0.0
			if v, isPresent := budget["ceiling"]; isPresent {
				ceiling, _ = strconv.ParseFloat(v, 64)
			}
			consume := 0.0
			if v, isPresent := budget["consume"]; isPresent {
				consume, _ = strconv.ParseFloat(v, 64)
			}
			if consume >= ceiling {
				ret = false
				reports = append(reports, "The consume more than the ceiling.")
			}
		}
	}
	
	return
}
