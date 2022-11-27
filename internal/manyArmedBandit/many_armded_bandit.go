package manyarmedbandit

import (
	"crypto/rand"
	"math"
	"math/big"

	ms "github.com/PalPalych7/OtusProjectWork/internal/mainstructs"
)

type BanditStruct struct {
	BanditConfig ms.BanditConfig
	Cend         float32 // цена деления (на какой процент уменьшаем "случайную величину" за 1 показ)
}

func New(bc ms.BanditConfig) *BanditStruct {
	cend := float32(100-bc.FinalRandomPecent) / float32(bc.PartialLearningCount)
	return &BanditStruct{bc, cend}
}

func randInt(maxV int) int {
	v, _ := rand.Int(rand.Reader, big.NewInt(int64(maxV+1)))
	randVal := int(v.Int64())
	return randVal
}

func kvadrProc(arrBS []ms.BannerStruct) int {
	arrSumKvProc := make([]float64, 0)
	var curKvProc float64
	var sumKvProc float64
	for _, v := range arrBS {
		if v.ShowCount > 0 {
			curKvProc = math.Pow(float64(v.ClickCount)/float64(v.ShowCount)*100, 2)
		} else {
			curKvProc = 0
		}
		sumKvProc += curKvProc
		arrSumKvProc = append(arrSumKvProc, sumKvProc)
	}
	sumKvProcInt := int(sumKvProc)
	randVal := randInt(sumKvProcInt)
	res := 0
	for i, v := range arrSumKvProc {
		if float64(randVal) <= v {
			res = i
			return res
		}
	}
	return res
}

func (b *BanditStruct) GetBannerNum(arrStruct []ms.BannerStruct) int {
	showSum := 0
	var res int
	for _, v := range arrStruct {
		showSum += v.ShowCount
	}
	if showSum <= b.BanditConfig.FullLearnigCount {
		// режим обучения
		return randInt(len(arrStruct) - 1)
	}
	var randomPecent float32
	if showSum <= b.BanditConfig.FullLearnigCount+b.BanditConfig.PartialLearningCount {
		// вычисляем "вероятностный процент" линейно от 100 до минимального
		randomPecent = float32(100) - float32(showSum-b.BanditConfig.FullLearnigCount)*b.Cend
	} else {
		// "вероятностный процент" - минимальный из конфига
		randomPecent = float32(b.BanditConfig.FinalRandomPecent)
	}
	if float32(randInt(101)) < randomPecent {
		res = randInt(len(arrStruct) - 1)
	} else {
		// подбор согласно квадратичным вероятностям просмотров
		res = kvadrProc(arrStruct)
	}
	return res
}
