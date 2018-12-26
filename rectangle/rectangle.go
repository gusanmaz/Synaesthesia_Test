package rectangle

import (
	"encoding/json"
	"math"
	"math/rand"
	"time"
)

const (
	xInc      float64 = 0.001
	yInc      float64 = 0.001
	xLen      float64 = 0.048 //0.032
	yLen      float64 = 0.09  //0.06
	rectCnt   int     = 30    // 50
	choiceNum int     = 4
)

var (
	class1Ratios []float64
	rectangles   []Rectangle
	class1Cnt    int
	class2Cnt    int
	randomer     *rand.Rand
)

type Rectangle struct {
	X1, Y1, X2, Y2 float64
	Class          int8
}

type TestData struct {
	Rectangles  []Rectangle
	Class1Cnt   int
	Class2Cnt   int
	AnswerClass int
	Answer      int
	Choices     []int
	Class1Color string
	Class2Color string
}

func Init() {
	class1Ratios = []float64{0.06, 0.12, 0.18} // 0.03 0.06 0.09
	rectangles = make([]Rectangle, rectCnt)
	randomer = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func firstNDigits(num float64, n int) float64 {
	wholeNum := int(num * math.Pow(10, float64(n)))
	wholeFNum := float64(wholeNum)
	return wholeFNum / (math.Pow(10, float64(n)))
}

func CreateNewArrangement() {
	//randomer := rand.New(rand.NewSource(time.Now().UnixNano()))

	var colCnt int = int(1.000 / xInc)
	var rowCnt int = int(1.000 / yInc)

	secondBuf := make([][]byte, rowCnt)
	for i, _ := range secondBuf {
		secondBuf[i] = make([]byte, colCnt)
	}
	//rectangles := make([]Rectangle, rectCnt)
	ratio1Select := randomer.Int() % len(class1Ratios)
	ratio1 := class1Ratios[ratio1Select]

	class1Cnt = 0
	class2Cnt = 0

	for i := 0; i < rectCnt; i++ {
		num := randomer.Float64()
		if num > ratio1 {
			class2Cnt++
			continue
		}
		class1Cnt++
	}

	// New Code
	if class1Cnt == 0 {
		class1Cnt = randomer.Int() % 5
		class1Cnt++
	}

	curRectCnt := 0
outerLoop:
	for curRectCnt < rectCnt {
		x1 := firstNDigits(randomer.Float64(), 3)
		y1 := firstNDigits(randomer.Float64(), 3)
		x2 := x1 + xLen
		y2 := y1 + yLen
		if x2 >= 1.000 || y2 >= 1.000 {
			continue
		}

		x1Cell := int(x1 / xInc)
		y1Cell := int(y1 / yInc)
		x2Cell := int(x2 / xInc)
		y2Cell := int(y2 / yInc)

		for row := x1Cell; row <= x2Cell; row++ {
			for col := y1Cell; col < y2Cell; col++ {
				if secondBuf[row][col] == 1 {
					continue outerLoop
				}
			}
		}

		for row := x1Cell; row <= x2Cell; row++ {
			for col := y1Cell; col < y2Cell; col++ {
				secondBuf[row][col] = 1
			}
		}
		if curRectCnt < class1Cnt {
			rectangles[curRectCnt] = Rectangle{x1, y1, x2, y2, 1}
		} else {
			rectangles[curRectCnt] = Rectangle{x1, y1, x2, y2, 2}
		}
		curRectCnt++
	}
}

func GetClassCnt(c int) int {
	if c == 1 {
		return class1Cnt
	}
	if c == 2 {
		return class2Cnt
	}
	return -1
}

func GetMultipleChoices(class, choiceNum int) []int {
	answer := GetClassCnt(class)
	answerInd := randomer.Int() % choiceNum
	choices := make([]int, choiceNum)
	choices[0] = answer - answerInd

	for i, _ := range choices {
		choices[i] = choices[0] + i
	}

	inc := 0
	/*if choices[0] < 0{
		inc = -choices[0]
		for i,_ := range choices{
			choices[i] = choices[i] + inc
		}
	}*/
	if choices[0] <= 0 {
		inc = -choices[0] + 1
		for i, _ := range choices {
			choices[i] = choices[i] + inc
		}
	}

	return choices
}

func GetJSObject() string {
	jsObj, _ := json.Marshal(rectangles)
	return string(jsObj)
}

func GetTestJSON(color1, color2 string) string {
	Init()
	CreateNewArrangement()
	choices := GetMultipleChoices(1, choiceNum)
	testData := TestData{}
	testData.Rectangles = rectangles
	testData.Class1Cnt = class1Cnt
	testData.Class2Cnt = class2Cnt
	testData.Answer = class1Cnt
	testData.AnswerClass = 1
	testData.Choices = choices
	testData.Class1Color = color1
	testData.Class2Color = color2
	jsObj, _ := json.Marshal(&testData)
	return string(jsObj)
}
