package stats

import (
	"github.com/gusanmaz/Synaesthesia_Test/templates"
	"math"
)

type UserStatistics struct {
	SlowerSingleAvg     float64
	SlowerMultiAvg      float64
	SlowerSinglePercent float64
	SlowerMultiPercent  float64
	SubjectSingleAvg    float64
	SubjectMultiAvg     float64
	SubjectRatio        float64
	TopPercent          float64
}

func CalculateUserStatistics(subject []int, population [][]int) UserStatistics {
	//avgArr := make([]int, 10)
	questionNum := len(subject)
	sumArr := make([]float64, questionNum)
	slowerResponses := make([]int, questionNum) // Slower responses of population
	for i, _ := range population {
		for k := 0; k < questionNum; k++ {
			sumArr[k] += float64(population[i][k])
			if subject[k] < population[i][k] {
				slowerResponses[k]++
			}
		}
	}

	subRatio, subSingleAvg, subMultiAvg := getResponseStats(subject)
	topCnt := 0
	for _, v := range population {
		ratio, _, _ := getResponseStats(v)
		// Added = case
		if ratio <= subRatio {
			topCnt++
		}
	}
	topPercent := float64(topCnt) / float64(len(population)) * 100

	slowerSingleResponses := 0
	singleTestCnt := 0
	slowerMultiResponses := 0
	multiTestCnt := 0

	for i, _ := range slowerResponses {
		if i%2 == 1 {
			slowerSingleResponses += slowerResponses[i]
			singleTestCnt++
		} else {
			slowerMultiResponses += slowerResponses[i]
			multiTestCnt++
		}
	}

	var slowerSingleAvg float64 = float64(slowerSingleResponses) / float64(singleTestCnt)
	var SlowerSinglePercent = slowerSingleAvg / float64(len(population)) * 100
	var slowerMultiAvg float64 = float64(slowerMultiResponses) / float64(multiTestCnt)
	var SlowerMultiPercent = slowerMultiAvg / float64(len(population)) * 100
	return UserStatistics{slowerSingleAvg, slowerMultiAvg,
		SlowerSinglePercent, SlowerMultiPercent, subSingleAvg,
		subMultiAvg, subRatio, topPercent}
}

func getResponseStats(resTimes []int) (ratio float64, singleAvg float64, multiAVg float64) {
	singleSum := 0
	multiSum := 0
	singleCnt := 0
	multiCnt := 0
	for i, _ := range resTimes {
		if i%2 == 1 {
			singleSum += resTimes[i]
			singleCnt++
		} else {
			multiSum += resTimes[i]
			multiCnt++
		}
	}
	//var ret2 float64 = float64(singleSum) / float64(len(resTimes))
	//var ret3 float64 = float64(multiSum) / float64(len(resTimes))

	var ret2 float64 = float64(singleSum) / float64(singleCnt)
	var ret3 float64 = float64(multiSum) / float64(multiCnt)

	return (ret2 / ret3), ret2, ret3
}

func InitStats(stats UserStatistics) {
	topPercent := math.Floor(stats.TopPercent*100) / 100
	singlePercent := math.Floor(stats.SlowerSinglePercent*100) / 100
	multiPercent := math.Floor(stats.SlowerMultiPercent*100) / 100
	subjectRatio := math.Floor(stats.SubjectRatio*100) / 100

	Globals := templates.Globals

	Globals["en"]["TopPercent"] = topPercent
	Globals["en"]["SinglePercent"] = singlePercent
	Globals["en"]["MultiPercent"] = multiPercent
	Globals["en"]["SubjectRatio"] = subjectRatio

	Globals["tr"]["TopPercent"] = topPercent
	Globals["tr"]["SinglePercent"] = singlePercent
	Globals["tr"]["MultiPercent"] = multiPercent
	Globals["tr"]["SubjectRatio"] = subjectRatio

	Globals["en"]["PageTitle"] = "Results"
	Globals["tr"]["PageTitle"] = "Sonuçlar"

}
