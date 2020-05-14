package service

import (
	"encoding/json"
	"fmt"
	"github.com/vincecfl/dex-robot/pkg"
	"github.com/vincecfl/go-common/log"
	"math/rand"
	"time"
)

const (
	dexContractAddr = "TJ86JLUrMEXYQPNXx1tyD1SzxEgPECFpmj"
	btcTokenAddr    = "TEQEni8FCPrmdTQPUKAu1DCpm3ZYESjFg8"
	trxTokenAddr    = "T9yD14Nj9j7xAB4dbGeiX9h8unkKHxuWwb"
	btcOwner        = "TPdBHYrTDiop2fgsmZGDEfNN5SucJADCf4"
	btcOwnerKey     = "514bfc62a1f84b69a46ba6478f991eacb136ef1a2f63a16a66e7f42c14c1de07"
	BUY             = 1
	SELL            = 2
)

func BuyBTCHandle() {
	url := fmt.Sprintf("https://bytego123.cn/dex/api/v1/market/pairOrder4Kline/query?pairID=%v", 1)
	result, err := pkg.Get(url, false, "")
	if err != nil {
		log.Errorf(err, "pkg.Get error")
		return
	}

	resp := &ResultResp{}
	if err := json.Unmarshal([]byte(result), resp); err != nil {
		log.Errorf(err, "json.Unmarshal error")
		return
	}
	if resp.Code != 0 {
		return
	}

	buyList := resp.Data.Buy
	buyLen := len(buyList)
	price := int64(resp.Data.Price * 1e6)

	if price <= 5*1e5 || buyList[buyLen-1].Price*1e6 <= 5*1e5 {
		if err := SetRobotType(1, 1); err != nil {
			log.Errorf(err, "SetRobotType error")
			return
		}
	}

	userAddr := btcOwner
	userKey := btcOwnerKey
	token1 := btcTokenAddr

	if buyLen >= 20 {
		log.Infof("buyLen more than 20")
		return
	}

	buyPrice := int64(0)

	if buyLen == 0 {
		tempPrice := price
		if tempPrice <= 5*1e5 {
			tempPrice = 6 * 1e5
		}
		// 比当前价格少10000
		buyPrice = tempPrice - 10000
	} else if buyLen < 20 {
		tempPrice := int64(buyList[buyLen-1].Price * 1e6)
		if tempPrice <= 5*1e5 {
			tempPrice = 6 * 1e5
		}
		// 比最近一单价格少1000~1500
		buyPrice = tempPrice - RandInt64(1000, 1500)
	}

	if buyPrice > 0 {
		token2 := trxTokenAddr
		amount1 := int64(0)
		if buyPrice <= 1*1e6 {
			amount1 = RandInt64(20, 30) * 1e6
		} else if buyPrice > 1*1e6 && buyPrice <= 2*1e6 {
			amount1 = RandInt64(10, 15) * 1e6
		} else if buyPrice > 2*1e6 {
			amount1 = RandInt64(5, 10) * 1e6
		}
		amount2 := amount1 * buyPrice / 1e6
		err = Buy(true, userAddr, userKey, token1, token2, amount1, amount2, buyPrice, 0)
		if err != nil {
			log.Errorf(err, "Buy error")
			return
		}
		log.Infof("BuyBTCHandle success")
		return
	}

}

func SellBTCHandle() {
	url := fmt.Sprintf("https://bytego123.cn/dex/api/v1/market/pairOrder4Kline/query?pairID=%v", 1)
	result, err := pkg.Get(url, false, "")
	if err != nil {
		log.Errorf(err, "pkg.Get error")
		return
	}
	resp := &ResultResp{}
	if err := json.Unmarshal([]byte(result), resp); err != nil {
		log.Errorf(err, "json.Unmarshal error")
		return
	}
	if resp.Code != 0 {
		return
	}

	sellList := resp.Data.Sell
	sellLen := len(sellList)
	price := int64(resp.Data.Price * 1e6)

	if price >= 30*1e5 || sellList[sellLen-1].Price*1e6 >= 30*1e5 {
		if err := SetRobotType(1, 2); err != nil {
			log.Errorf(err, "SetRobotType error")
			return
		}
	}

	userAddr := btcOwner
	userKey := btcOwnerKey
	token1 := btcTokenAddr

	if sellLen >= 20 {
		log.Infof("sellLen more than 20")
		return
	}

	sellPrice := int64(0)

	if sellLen == 0 {
		tempPrice := price
		if tempPrice >= 30*1e5 {
			tempPrice = 29 * 1e5
		}
		// 比当前价格多10000
		sellPrice = tempPrice + 10000
	} else if sellLen < 20 {
		tempPrice := int64(sellList[sellLen-1].Price * 1e6)
		if tempPrice >= 30*1e5 {
			tempPrice = 29 * 1e5
		}
		// 最近一单价格多1000~1500
		sellPrice = tempPrice + RandInt64(1000, 1500)
	}

	if sellPrice > 0 {
		token2 := trxTokenAddr
		amount1 := int64(0)
		if sellPrice <= 1*1e6 {
			amount1 = RandInt64(20, 30) * 1e6
		} else if sellPrice > 1*1e6 && sellPrice <= 2*1e6 {
			amount1 = RandInt64(10, 15) * 1e6
		} else if sellPrice > 2*1e6 {
			amount1 = RandInt64(5, 10) * 1e6
		}
		amount2 := amount1 * sellPrice / 1e6
		err := Approve(btcTokenAddr, userAddr, userKey, dexContractAddr, amount1)
		if err != nil {
			log.Errorf(err, "Approve error")
			return
		}
		err = Sell(false, userAddr, userKey, token1, token2, amount1, amount2, sellPrice, 0)
		if err != nil {
			log.Errorf(err, "sell error")
			return
		}
		log.Infof("SellBTCHandle success")
	}

}

func TradeBTCHandle() {
	url := fmt.Sprintf("https://bytego123.cn/dex/api/v1/market/pairOrder4Kline/query?pairID=%v", 1)
	result, err := pkg.Get(url, false, "")
	if err != nil {
		log.Errorf(err, "pkg.Get error")
		return
	}

	resp := &ResultResp{}
	if err := json.Unmarshal([]byte(result), resp); err != nil {
		log.Errorf(err, "json.Unmarshal error")
		return
	}
	if resp.Code != 0 {
		return
	}

	buyList := resp.Data.Buy
	buyLen := len(buyList)
	sellList := resp.Data.Sell
	sellLen := len(sellList)

	log.Infof("TradeBTCHandle buyLen:%v, sellLen:%v", buyLen, sellLen)

	if buyLen <= 6 || sellLen <= 6 {
		log.Infof("buyLen or sellLen less than 6")
		return
	}

	robotType := GetRobotType(1)
	if robotType == 0 {
		log.Errorf(nil, "robotType is 0")
		return
	}

	log.Infof("robotType:%v", robotType)

	currentTime := time.Now().Unix()

	orderType := BUY
	rand := RandInt64(1, 101)
	// robotType为2 卖单为主
	if robotType == 2 {
		time60 := currentTime % (60 * 60)
		time15 := currentTime % (15 * 60)
		if time15 <= 300 && time60 < 2700 {
			sell4Five(buyList)
			return
		} else if time15 <= 300 && time60 >= 2700 {
			buy4Five(sellList)
			return
		}
		if rand <= 30 {
			orderType = BUY
		} else {
			orderType = SELL
		}
	} else if robotType == 1 {
		time60 := currentTime % (60 * 60)
		time15 := currentTime % (15 * 60)
		if time15 <= 300 && time60 < 2700 {
			buy4Five(sellList)
			return
		} else if time15 <= 300 && time60 >= 2700 {
			sell4Five(buyList)
			return
		}
		if rand <= 70 {
			orderType = BUY
		} else {
			orderType = SELL
		}
	}

	if orderType == BUY {
		buy(sellList)
	} else {
		sell(buyList)
	}
	return
}

func buy(sellList []*PairOrderModel) error {
	userAddr := btcOwner
	userKey := btcOwnerKey
	token1 := btcTokenAddr
	buyPrice := int64(sellList[0].Price * 1e6)
	token2 := trxTokenAddr
	amount1 := RandInt64(20, 30) * 1e6
	amount2 := amount1 * buyPrice / 1e6
	err := Buy(true, userAddr, userKey, token1, token2, amount1, amount2, buyPrice, 0)
	if err != nil {
		log.Errorf(err, "Buy error")
		return err
	}
	log.Infof("TradeBTCHandle buy success")
	return nil
}

func buy4Five(sellList []*PairOrderModel) error {
	userAddr := btcOwner
	userKey := btcOwnerKey
	token1 := btcTokenAddr
	buyPrice := int64(sellList[4].Price * 1e6)
	token2 := trxTokenAddr
	amount1 := int64(0)
	for i := 0; i <= 4; i++ {
		amount1 += int64(sellList[i].TotalQuoteAmount * 1e6)
	}
	amount1 += 20 * 1e6
	amount2 := amount1 * buyPrice / 1e6
	err := Buy(true, userAddr, userKey, token1, token2, amount1, amount2, buyPrice, 0)
	if err != nil {
		log.Errorf(err, "Buy error")
		return err
	}
	log.Infof("TradeBTCHandle buy4Five success")
	return nil
}

func sell(buyList []*PairOrderModel) error {
	userAddr := btcOwner
	userKey := btcOwnerKey
	token1 := btcTokenAddr
	sellPrice := int64(buyList[0].Price * 1e6)
	token2 := trxTokenAddr
	amount1 := RandInt64(20, 30) * 1e6
	amount2 := amount1 * sellPrice / 1e6
	err := Approve(btcTokenAddr, userAddr, userKey, dexContractAddr, amount1)
	if err != nil {
		log.Errorf(err, "Approve error")
		return err
	}
	err = Sell(false, userAddr, userKey, token1, token2, amount1, amount2, sellPrice, 0)
	if err != nil {
		log.Errorf(err, "sell error")
		return err
	}
	log.Infof("TradeBTCHandle sell success")
	return nil
}

func sell4Five(buyList []*PairOrderModel) error {
	userAddr := btcOwner
	userKey := btcOwnerKey
	token1 := btcTokenAddr
	sellPrice := int64(buyList[4].Price * 1e6)
	token2 := trxTokenAddr
	amount1 := RandInt64(20, 30) * 1e6
	for i := 0; i <= 4; i++ {
		amount1 += int64(buyList[i].TotalQuoteAmount * 1e6)
	}
	amount1 += 20 * 1e6
	amount2 := amount1 * sellPrice / 1e6
	err := Approve(btcTokenAddr, userAddr, userKey, dexContractAddr, amount1)
	if err != nil {
		log.Errorf(err, "Approve error")
		return err
	}
	err = Sell(false, userAddr, userKey, token1, token2, amount1, amount2, sellPrice, 0)
	if err != nil {
		log.Errorf(err, "sell error")
		return err
	}
	log.Infof("TradeBTCHandle sell4Five success")
	return nil
}

func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}

type ResultResp struct {
	Code int           `json:"code"`
	Data PairOrderResp `json:"data"`
}

type PairOrderResp struct {
	Buy   []*PairOrderModel `json:"buy"`
	Sell  []*PairOrderModel `json:"sell"`
	Price float64           `json:"price"`
}

type PairOrderModel struct {
	Price            float64 `json:"price"`
	TotalQuoteAmount float64 `json:"totalQuoteAmount"`
	TotalBaseAmount  float64 `json:"totalBaseAmount"`
	TotalOrder       int64   `json:"totalOrder"`
}
