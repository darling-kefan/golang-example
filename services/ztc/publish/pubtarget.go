package publish

import (
	_ "fmt"
	"time"
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
	if weekday := date.Weekday(); weekday == time.Sunday {
		week := 7
	} else {
		week := int(weekday)
	}
	isCurHourValid = false
	if hours, ok := pt.pub.times[strconv.Itoa(week)]; ok {
		hour := time.Now().Hour()
		for _, v := range hours {
			if strconv.Itoa(hour) == strings.TrimLeft(v, "0") {
				isCurHourValid = true
			}
		}
	} else {
		isCurHourValid = true
	}
	if isCurHourValid == false {
		ret = false
		reports = append(reports, fmt.Sprintf("Hour points: %v", hours))
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
		if logo, ok := p.pub.card["logo"]; !ok {
			ret = false
			reports = append(reports, "Card logo is not existed.")
		}
		if logoStr, err := logo.(string); err == nil {
			logoStr = strings.ToLower(logoStr)
			if !strings.HasPrefix(logoStr, "http://") || !strings.HasPrefix(logoStr, "https://") {
				ret = false
				reports = append(reports, "Card logo is valid.")
			}
		}
	}
	// 视频和开机大图不验证url
	if displayWay != DISPLAY_VIDEO || displayWay != DISPLAY_KJDT {
		if url, ok := p.pub.card["url"]; !ok {
			ret = false
			reports = append(reports, "Card url is not existed.")
		}
		if urlStr, err := url.(string); err == nil {
			urlStr = strings.ToLower(urlStr)
			if !strings.HasPrefix(urlStr, "http://") || !strings.HasPrefix(urlStr, "https://") {
				ret = false
				reports = append(reports, "Card url is valid.")
			}
		}
	}

	// 判定广告主是否被屏蔽
	

	// 判断广告主是否还有余额(测试广告主不检测广告主余额)

	// 判断预算是否花超

	// 判断广告主资质
	
	return
}
