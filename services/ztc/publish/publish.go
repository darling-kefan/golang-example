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
	Id            int // 计划id
	Title         string //标题
	AdvertiserId  int // 广告主id
	CardId        int // 创意id
	AdpositionId  int // 广告位，1-卡券、2-视频、3-开机大图、4-Banner广告
	Billing       int //计费模式，1-曝光计费、2-点击计费、3-到达计费
	Bid           int //出价，单位分
	Ceiling       int //日预算，单位分
	AgainNum      float64 //权重

	Status        int //状态，1-正常、2-暂停、3-停止
	VersionNo     int //自投放平台版本号，1-版本v1、2-版本v2
	CardStatus    int //创意状态，1-审核通过、2-等待初审、4-复审不通过、5-等待复审、6-初审不通过
	PromoteRoute  string //推广渠道，1-微信、2-APP、3-支付宝，多个渠道用英文逗号分割

	Advername     string //广告主姓名
	Tvmid         string //tvmid
	Shortcom      string //公司简称
	AdverType     int //广告主类型，1-正式账户、2-运营账户、3-测试账户
	AdverSubtype  int //广告主子类型(仅在advertiser_type=3时有效)，1-技术测试、2-免费帐号、3-试投帐号
	Istop         bool //是否置顶，true-置顶，false-不置顶

	Stime         time.Time //开始投放时间
	Etime         time.Time //结束投放时间
	CreatedAt     time.Time //创建时间
	UpdatedAt     time.Time //更新时间
	DeletedAt     time.Time //删除时间

	Version       string //版本号
	Deliverychans []string //投放渠道
	Targets       map[string][]string //定向信息
	Times         map[string][]string //投放时间
	Card          map[string]string //卡券信息

	Product       int //产品类型，1-基础版、2-专业版
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
			p.Id = v
		} else {
			panic(err.Error())
		}
	}
	if title, ok := pub["title"]; ok {
		p.Title = title
	}
	if advertiserId, ok := pub["advertiser_id"]; ok {
		if v, err := strconv.Atoi(advertiserId); err == nil {
			p.AdvertiserId = v
		}
	}
	if cardId, ok := pub["card_id"]; ok {
		if v, err := strconv.Atoi(cardId); err == nil {
			p.CardId = v
		}
	}
	if adpositionId, ok := pub["adposition_id"]; ok {
		if v, err := strconv.Atoi(adpositionId); err == nil {
			p.AdpositionId = v
		}
	}
	if billing, ok := pub["billing"]; ok {
		if v, err := strconv.Atoi(billing); err == nil {
			p.Billing = v
		}
	}
	if bid, ok := pub["bid"]; ok {
		if v, err := strconv.Atoi(bid); err == nil {
			p.Bid = v
		}
	}
	if ceiling, ok := pub["ceiling"]; ok {
		if v, err := strconv.Atoi(ceiling); err == nil {
			p.Ceiling = v
		}
	}
	if againNum, ok := pub["again_num"]; ok {
		if v, err := strconv.ParseFloat(againNum, 64); err == nil {
			p.AgainNum = v
		}
	}
	if status, ok := pub["status"]; ok {
		if v, err := strconv.Atoi(status); err == nil {
			p.Status = v
		}
	}
	if versionNo, ok := pub["versionNo"]; ok {
		if v, err := strconv.Atoi(versionNo); err == nil {
			p.VersionNo = v
		}
	}
	if cardStatus, ok := pub["cardStatus"]; ok {
		if v, err := strconv.Atoi(cardStatus); err == nil {
			p.CardStatus = v
		}
	}
	if promoteRoute, ok := pub["promoteRoute"]; ok {
		p.PromoteRoute = promoteRoute
	}
	if adverName, ok := pub["advertiser_name"]; ok {
		p.Advername = adverName
	}
	if tvmid, ok := pub["tvmid"]; ok {
		p.Tvmid = tvmid
	}
	if shortcom, ok := pub["shortcom"]; ok {
		p.Shortcom = shortcom
	}
	if adverType, ok := pub["adverType"]; ok {
		if v, err := strconv.Atoi(adverType); err == nil {
			p.AdverType = v
		}
	}
	if adverSubtype, ok := pub["adverSubtype"]; ok {
		if v, err := strconv.Atoi(adverSubtype); err == nil {
			p.AdverSubtype = v
		}
	}
	if istop, ok := pub["istop"]; ok {
		if v, err := strconv.ParseBool(istop); err == nil {
			p.Istop = v
		}
	}
	if stime, ok := pub["stime"]; ok {
		if t, err := time.Parse("2006-01-02 15:04:05", stime); err == nil {
			p.Stime = t
		}
	}
	if etime, ok := pub["etime"]; ok {
		if t, err := time.Parse("2006-01-02 15:04:05", etime); err == nil {
			p.Etime = t
		}
	}
	if createdAt, ok := pub["created_at"]; ok {
		if t, err := time.Parse("2006-01-02 15:04:05", createdAt); err == nil {
			p.CreatedAt = t
		}
	}
	if updatedAt, ok := pub["updated_at"]; ok {
		if t, err := time.Parse("2006-01-02 15:04:05", updatedAt); err == nil {
			p.UpdatedAt = t
		}
	}
	if deletedAt, ok := pub["deleted_at"]; ok {
		if t, err := time.Parse("2006-01-02 15:04:05", deletedAt); err == nil {
			p.DeletedAt = t
		}
	}
}

func (p *Pub) setPubversion() {
	pubversion := models.Publishs().LastestPubversion(p.Id, p.VersionNo)
	if version, ok := pubversion["version"].(string); ok {
		p.Version = version
	}
	p.Targets = make(map[string][]string)
	if targets, ok := pubversion["target"].([]interface{}); ok {
		for _, v := range targets {
			if target, ok := v.(map[string]interface{}); ok {
				tarkey := strconv.FormatFloat(target["id"].(float64), 'f', 0, 64)
				if tarval, ok := target["value"].([]interface{}); ok {
					for _, vv := range tarval {
						switch vvv := vv.(type) {
						case float64:
							p.Targets[tarkey] = append(p.Targets[tarkey], strconv.FormatFloat(vvv, 'f', 0, 64))
						case string:
							p.Targets[tarkey] = append(p.Targets[tarkey], vvv)
						}
					}
				}
			}
		}
	}
	p.Times = make(map[string][]string)
	if times, ok := pubversion["time"].(map[string]interface{}); ok {
		for tk, tv := range times {
			if tvv, ok := tv.([]interface{}); ok {
				for _, tvvv := range tvv {
					tvvvv := tvvv.(string)
					p.Times[tk] = append(p.Times[tk], tvvvv[0:2])
				}
			}
		}
	}
	p.Deliverychans = make([]string, 3)
	if delichans, ok := pubversion["deliverychannels"].([]interface{}); ok {
		for dk, dv := range delichans {
			dvv := dv.(float64)
			p.Deliverychans[dk] = strconv.FormatFloat(dvv, 'f', 0, 64)
		}
	}
	p.Card = make(map[string]string, 4)
	if cards, ok := pubversion["card"].([]interface{}); ok {
		if card, ok := cards[0].(map[string]interface{}); ok {
			for ck, cv := range card {
				if ck == "title" || ck == "logo" || ck == "url" {
					p.Card[ck] = cv.(string)
				} else if (ck == "type") {
					switch cvv := cv.(type) {
					case float64:
						p.Card[ck] = strconv.FormatFloat(cvv, 'f', 0, 64)
					case string:
						p.Card[ck] = cvv
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
		if v == p.AdvertiserId {
			isBasic = true
		}
	}
	if isBasic == true {
		p.Product = PRODUCT_BASIC
	} else {
		p.Product = PRODUCT_PROFESSIONAL
	}
}

func (p *Pub) Agents() (agents []int) {
	if len(p.PromoteRoute) <= 0 {
		agents = []int{1}
	} else {
		pr := strings.Split(p.PromoteRoute, ",")
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
	return p.AdpositionId
}
