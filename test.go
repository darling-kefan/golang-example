package main

import (
	"fmt"
	"models"
	pubservice "services/ztc/publish"
	advservice "services/advertiser"
	_ "strconv"
	_ "reflect"
)

func main() {
	boolRet := advservice.New(90).HasSyncCondition()
	fmt.Println(boolRet)
	return
	//shields := models.Shields().Find(2, []string{"1","2"}, 90,[]string{"*"})
	//fmt.Println(shields)
	// publishs := models.Publishs().JoinPubsByIds([]int{22983,22969})
	publishs := models.Publishs().JoinPubsByIds([]int{22983})
	pub := pubservice.NewPub(publishs[0])
	// fmt.Println(pub)
	//fmt.Println(pub.Product())
	//fmt.Println(pub.Agents())
	pubtarget := pubservice.NewPubTarget(pub, "")
	//fmt.Println(pubtarget)
	ok, reports := pubtarget.Detect()
	fmt.Println(ok, reports)

	//publishs = models.Publishs().PubsByIds([]string{"1000", "22969", "22974"})
	//fmt.Println(publishs)

	//fmt.Println("------------------------------------------------------")
	//publishs := models.Publishs().LastestPubversion(strconv.Itoa(22983), 2)
	//fmt.Println(publishs)

	//fmt.Println("------------------------------------------------------")

	//cards := models.Cards().CardsByIds([]string{"160", "161", "177"})
	//fmt.Println(cards)

	//fmt.Println("------------------------------------------------------")

	//advertisers := models.Advertisers().AdversByIds([]string{"1", "3", "4", "5", "6"}, []string{"id", "tvmid", "name", "mobileno", "email", "type", "subtype"})
	//fmt.Println(advertisers)

	//advQualifis := models.Advertisers().Qualifis([]string{"1", "3", "4", "5", "6"})
	//fmt.Println(advQualifis)
}
