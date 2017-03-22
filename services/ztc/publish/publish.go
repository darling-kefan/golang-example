package publish

import (
	_ "fmt"
	"time"
	"strconv"
	"strings"
	
	"models"
)

const (
	PRODUCT_BASIC, PRODUCT_PROFESSIONAL = 1, 2
	DISPLAY_CARD, DISPLAY_VIDEO, DISPLAY_KJDT, DISPLAY_BANNER = 1, 2, 3, 4
)

type Pub struct {
	id            int // 计划id
	title         string //标题
	advertiserId  int // 广告主id
	cardId        int // 创意id
	adpositionId  int // 广告位，1-卡券、2-视频、3-开机大图、4-Banner广告
	billing       int //计费模式，1-曝光计费、2-点击计费、3-到达计费
	bid           float64 //出价，单位分
	ceiling       float64 //日预算，单位分
	againNum      float64 //权重

	status        int //状态，1-正常、2-暂停、3-停止
	versionNo     int //自投放平台版本号，1-版本v1、2-版本v2
	cardStatus    int //创意状态，1-审核通过、2-等待初审、4-复审不通过、5-等待复审、6-初审不通过
	promoteRoute  string //推广渠道，1-微信、2-APP、3-支付宝，多个渠道用英文逗号分割
	
	adverType     int //广告主类型，1-正式账户、2-运营账户、3-测试账户
	adverSubtype  int //广告主子类型(仅在advertiser_type=3时有效)，1-技术测试、2-免费帐号、3-试投帐号
	istop         bool //是否置顶，true-置顶，false-不置顶

	stime         time.Time //开始投放时间
	etime         time.Time //结束投放时间
	createdAt     time.Time //创建时间
	updatedAt     time.Time //更新时间
	deletedAt     time.Time //删除时间

	version       string //版本号
	deliverychans []string //投放渠道
	targets       map[string][]string //定向信息
	times         map[string][]string //投放时间
	card          map[string]interface{} //卡券信息

	product       int //产品类型，1-基础版、2-专业版
}

func NewPub(pub map[string]string) (p *Pub) {
	p = new(Pub)
	p.setPub(pub)
	p.setPubversion()
	p.setProduct()
	return
}

func (p *Pub) setPub(pub map[string]string) {
	if id, ok := pub["id"]; ok {
		if v, err := strconv.Atoi(id); err == nil {
			p.id = v
		} else {
			panic(err.Error())
		}
	}
	if title, ok := pub["title"]; ok {
		p.title = title
	}
	if advertiserId, ok := pub["advertiser_id"]; ok {
		if v, err := strconv.Atoi(advertiserId); err == nil {
			p.advertiserId = v
		}
	}
	if cardId, ok := pub["card_id"]; ok {
		if v, err := strconv.Atoi(cardId); err == nil {
			p.cardId = v
		}
	}
	if adpositionId, ok := pub["adposition_id"]; ok {
		if v, err := strconv.Atoi(adpositionId); err == nil {
			p.adpositionId = v
		}
	}
	if billing, ok := pub["billing"]; ok {
		if v, err := strconv.Atoi(billing); err == nil {
			p.billing = v
		}
	}
	if bid, ok := pub["bid"]; ok {
		if v, err := strconv.ParseFloat(bid, 64); err == nil {
			p.bid = v
		}
	}
	if ceiling, ok := pub["ceiling"]; ok {
		if v, err := strconv.ParseFloat(ceiling, 64); err == nil {
			p.ceiling = v
		}
	}
	if againNum, ok := pub["again_num"]; ok {
		if v, err := strconv.ParseFloat(againNum, 64); err == nil {
			p.againNum = v
		}
	}
	if status, ok := pub["status"]; ok {
		if v, err := strconv.Atoi(status); err == nil {
			p.status = v
		}
	}
	if versionNo, ok := pub["versionNo"]; ok {
		if v, err := strconv.Atoi(versionNo); err == nil {
			p.versionNo = v
		}
	}
	if cardStatus, ok := pub["cardStatus"]; ok {
		if v, err := strconv.Atoi(cardStatus); err == nil {
			p.cardStatus = v
		}
	}
	if promoteRoute, ok := pub["promoteRoute"]; ok {
		p.promoteRoute = promoteRoute
	}
	if adverType, ok := pub["adverType"]; ok {
		if v, err := strconv.Atoi(adverType); err == nil {
			p.adverType = v
		}
	}
	if adverSubtype, ok := pub["adverSubtype"]; ok {
		if v, err := strconv.Atoi(adverSubtype); err == nil {
			p.adverSubtype = v
		}
	}
	if istop, ok := pub["istop"]; ok {
		if v, err := strconv.ParseBool(istop); err == nil {
			p.istop = v
		}
	}
	if stime, ok := pub["stime"]; ok {
		if t, err := time.Parse("2006-01-02 15:04:05", stime); err == nil {
			p.stime = t
		}
	}
	if etime, ok := pub["etime"]; ok {
		if t, err := time.Parse("2006-01-02 15:04:05", etime); err == nil {
			p.etime = t
		}
	}
	if createdAt, ok := pub["created_at"]; ok {
		if t, err := time.Parse("2006-01-02 15:04:05", createdAt); err == nil {
			p.createdAt = t
		}
	}
	if updatedAt, ok := pub["updated_at"]; ok {
		if t, err := time.Parse("2006-01-02 15:04:05", updatedAt); err == nil {
			p.updatedAt = t
		}
	}
	if deletedAt, ok := pub["deleted_at"]; ok {
		if t, err := time.Parse("2006-01-02 15:04:05", deletedAt); err == nil {
			p.deletedAt = t
		}
	}
}

func (p *Pub) setPubversion() {
	pubversion := models.Publishs().LastestPubversion(p.id, p.versionNo)
	if version, ok := pubversion["version"].(string); ok {
		p.version = version
	}
	p.targets = make(map[string][]string)
	if targets, ok := pubversion["target"].([]interface{}); ok {
		for _, v := range targets {
			if target, ok := v.(map[string]interface{}); ok {
				tarkey := strconv.FormatFloat(target["id"].(float64), 'f', 0, 64)
				if tarval, ok := target["value"].([]interface{}); ok {
					for _, vv := range tarval {
						switch vvv := vv.(type) {
						case float64:
							p.targets[tarkey] = append(p.targets[tarkey], strconv.FormatFloat(vvv, 'f', 0, 64))
						case string:
							p.targets[tarkey] = append(p.targets[tarkey], vvv)
						}
					}
				}
			}
		}
	}
	p.times = make(map[string][]string)
	if times, ok := pubversion["time"].(map[string]interface{}); ok {
		for tk, tv := range times {
			if tvv, ok := tv.([]interface{}); ok {
				for _, tvvv := range tvv {
					tvvvv := tvvv.(string)
					p.times[tk] = append(p.times[tk], tvvvv[0:2])
				}
			}
		}
	}
	p.deliverychans = make([]string, 3)
	if delichans, ok := pubversion["deliverychannels"].([]interface{}); ok {
		for dk, dv := range delichans {
			dvv := dv.(float64)
			p.deliverychans[dk] = strconv.FormatFloat(dvv, 'f', 0, 64)
		}
	}
	p.card = make(map[string]interface{}, 4)
	if cards, ok := pubversion["card"].([]interface{}); ok {
		if card, ok := cards[0].(map[string]interface{}); ok {
			for ck, cv := range card {
				if ck == "title" || ck == "logo" || ck == "url" {
					p.card[ck] = cv.(string)
				} else if (ck == "type") {
					switch cvv := cv.(type) {
					case float64:
						p.card[ck] = strconv.FormatFloat(cvv, 'f', 0, 64)
					case string:
						p.card[ck] = cvv
					}
				}
			}
		}
	}
}

func (p *Pub) setProduct() {
	isBasic := false
	basicAdvids := models.Advversions().BasicAdverIds()
	for _, v := range basicAdvids {
		if v == p.advertiserId {
			isBasic = true
		}
	}
	if isBasic == true {
		p.product = PRODUCT_BASIC
	} else {
		p.product = PRODUCT_PROFESSIONAL
	}
}

func (p *Pub) Product() int {
	return p.product
}

func (p *Pub) Agents() (agents []int) {
	if len(p.promoteRoute) <= 0 {
		agents = []int{1}
	} else {
		pr := strings.Split(p.promoteRoute, ",")
		for _, v := range pr {
			if vv, err := strconv.Atoi(v); err == nil {
				agents = append(agents, vv)
			}
		}
	}
	return
}

// DisplayWay get the ad display way
func (p *Pub) DisplayWay() int {
	return p.adpositionId
}
