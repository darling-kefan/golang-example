package main

import (
	"fmt"
	"strconv"
	"strings"
	"net/http"
	"encoding/json"
	
	"models"
	pubService "services/ztc/publish"
)

type PubInfo struct {
	Id            string  `json:"id"`
	Title         string  `json:"title"`
	Bid           string  `json:"bid"`
	Ceiling       string  `json:"ceiling"`
	Stime         string  `json:"stime"`
	Etime         string  `json:"etime"`
	Status        string  `json:"status"`
	AdpositionId  string  `json:"adposition_id"`
	Billing       string  `json:"billing"`
	Deliverychans string  `json:"deliverychans"`
	Version       string  `json:"version"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	DeletedAt     string  `json:"deleted_at"`
}

type AdverInfo struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	Tvmid      string  `json:"tvmid"`
	AdverType  string  `json:"advertype"`
	Istop      string  `json:"istop"`
	Shortcom   string  `json:"shortcom"`
}

type CardInfo struct {
	Id     string  `json:"id"`
	Title  string  `json:"title"`
	Logo   string  `json:"logo"`
	Url    string  `json:"url"`
}

type DetectReports struct {
	Status  bool     `json:"status"`
	Reports []string `json:"reports"`
}

type Data struct {
	Publish    *PubInfo          `json:"publish"`
	Advertiser *AdverInfo        `json:"advertiser"`
	Card       *CardInfo         `json:"card"`
	Times      map[string]string `json:"times"`
	Targets    map[string]string `json:"targets"`
	DetectRepo *DetectReports    `json:"detectReports"`
}

type Result struct {
	Code   string `json:"code"`
	Errmsg string `json:"errmsg"`
	Dt     *Data  `json:"data"`
}

func main() {
	http.HandleFunc("/publish/profile", handler)
	http.ListenAndServe(":9001", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	pubid, _ := strconv.Atoi(r.Form["pubid"][0])
	date := r.Form["date"][0]

	// defer models.Conn.Close()
	pubs := models.Publishs().JoinPubsByIds([]int{pubid})
	if len(pubs) == 0 {
		render(w, appErrorf(r.Form["pubid"][0] + " is not exists."))
		return
	}
	pub  := pubService.NewPub(pubs[0])

	pubInfo := new(PubInfo)
	pubInfo.Id = strconv.Itoa(pub.Id)
	pubInfo.Title = pub.Title
	pubInfo.Bid = strconv.Itoa(pub.Bid)
	pubInfo.Ceiling = strconv.Itoa(pub.Ceiling)
	pubInfo.Stime = pub.Stime.Format("2006-01-02 15:04:05")
	pubInfo.Etime = pub.Etime.Format("2006-01-02 15:04:05")
	switch pub.AdpositionId {
	case 1:
		pubInfo.AdpositionId = strconv.Itoa(pub.AdpositionId) + "-卡券"
	case 2:
		pubInfo.AdpositionId = strconv.Itoa(pub.AdpositionId) + "-视频"
	case 3:
		pubInfo.AdpositionId = strconv.Itoa(pub.AdpositionId) + "-开机大图"
	case 4:
		pubInfo.AdpositionId = strconv.Itoa(pub.AdpositionId) + "-Banner广告"
	default:
		pubInfo.AdpositionId = strconv.Itoa(pub.AdpositionId) + "-未知"
	}
	switch pub.Billing {
	case 1:
		pubInfo.Billing = strconv.Itoa(pub.Billing) + "-曝光计费"
	case 2:
		pubInfo.Billing = strconv.Itoa(pub.Billing) + "-点击计费"
	case 3:
		pubInfo.Billing = strconv.Itoa(pub.Billing) + "-到达计费"
	default:
		pubInfo.Billing = strconv.Itoa(pub.Billing) + "-未知"
	}
	switch pub.Status {
	case 1:
		pubInfo.Status = strconv.Itoa(pub.Status) + "-正常"
	case 2:
		pubInfo.Status = strconv.Itoa(pub.Status) + "-暂停"
	case 3:
		pubInfo.Status = strconv.Itoa(pub.Status) + "-停止"
	default:
		pubInfo.Status = strconv.Itoa(pub.Status) + "-未知"
	}
	deliverychans := ""
	for _, delivery := range pub.Deliverychans {
		switch delivery {
		case "1":
			deliverychans += "1-微信,"
		case "2":
			deliverychans += "2-APP,"
		case "3":
			deliverychans += "3-支付宝,"
		}
	}
	deliverychans = strings.TrimRight(deliverychans, ",")
	pubInfo.Deliverychans = deliverychans
	pubInfo.Version = pub.Version
	pubInfo.CreatedAt = pub.CreatedAt.Format("2006-01-02 15:04:05")
	pubInfo.UpdatedAt = pub.UpdatedAt.Format("2006-01-02 15:04:05")
	pubInfo.DeletedAt = pub.DeletedAt.Format("2006-01-02 15:04:05")
	
	adverInfo := new(AdverInfo)
	adverInfo.Id = strconv.Itoa(pub.AdvertiserId)
	adverInfo.Name = pub.Advername
	adverInfo.Tvmid = pub.Tvmid
	switch pub.AdverType {
	case 1:
		adverInfo.AdverType = strconv.Itoa(pub.AdverType) + "-正式"
	case 2:
		adverInfo.AdverType = strconv.Itoa(pub.AdverType) + "-运营"
	case 3:
		switch pub.AdverSubtype {
		case 1:
			adverInfo.AdverType = strconv.Itoa(pub.AdverType) + "-测试账户-技术测试"
		case 2:
			adverInfo.AdverType = strconv.Itoa(pub.AdverType) + "-测试账户-免费账户"
		case 3:
			adverInfo.AdverType = strconv.Itoa(pub.AdverType) + "-测试账户-试投账户"
		}
	}
	if pub.Istop {
		adverInfo.Istop = "置顶"
	} else {
		adverInfo.Istop = "非置顶"
	}
	adverInfo.Shortcom = pub.Shortcom
	
	cardInfo := new(CardInfo)
	cardInfo.Id = strconv.Itoa(pub.CardId)
	cardInfo.Title = pub.Card["title"]
	cardInfo.Logo = pub.Card["logo"]
	cardInfo.Url = pub.Card["url"]

	times := make(map[string]string, 0)
	if len(pub.Times) > 0 {
		for week, v := range pub.Times {
			hours := strings.Join(v, ",")
			times[week] = hours
		}
	}

	targets := make(map[string]string, 0)
	for tartype, tarval := range pub.Targets {
		targets[tartype] = strings.Join(tarval, ",")
	}

	pubtarget := pubService.NewPubTarget(pub, date)
	ok, reports := pubtarget.Detect()

	detectReports := new(DetectReports)
	detectReports.Status  = ok
	if ok {
		detectReports.Reports = []string{}
	} else {
		detectReports.Reports = reports
	}
	
	data := new(Data)
	data.Publish    = pubInfo
	data.Advertiser = adverInfo
	data.Card       = cardInfo
	data.DetectRepo = detectReports
	data.Times      = times
	data.Targets    = targets

	result := new(Result)
	result.Code = "0"
	result.Errmsg = ""
	result.Dt = data
	
	render(w, result)
}

type appError struct {
	Code   string `json:"code"`
	Errmsg string `json:"errmsg"`
}

func appErrorf(errmsg string) *appError {
	return &appError{
		Code:   "500",
		Errmsg: errmsg,
	}
}

func render(w http.ResponseWriter, v interface{}) {
	ret, _ := json.Marshal(v)
	fmt.Fprintf(w, "%s", ret)
}
