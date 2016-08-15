package main

import (
	"bufio"
	"fmt"
	"strings"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

const MAX_LOCAL_ASTEROID_NUMBER = 450134
const RESONANCE_AXIS = 5.204
const SKIP_LINE_COUNT = 6

func getTrojanAsteroidsFromWeb(url string) []string {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal("Error during download data")
		os.Exit(-1)
	}
	trojanContent := doc.Find("#main pre").Text()
	var re = regexp.MustCompile("\\([0-9]+\\)")
	var matches = re.FindAllString(trojanContent, -1)
	var asteroidNumbers []string
	for _, val := range matches {
		asteroidNumbers = append(asteroidNumbers, val[1:len(val)-1])
	}
	return asteroidNumbers
}

func checkFile(filepath string) {
	if _, err := os.Stat(filepath); err != nil {
		fmt.Println(filepath, "doesn't exist")
		os.Exit(-1)
	}
}

type FileDesc struct {
	file    *os.File
	scanner *bufio.Scanner
}

func getFileDesc(filepath string) *FileDesc {
	var file, err = os.Open(filepath)
	if err != nil {
		log.Fatal("File cannot be read")
		os.Exit(-1)
	}
	return &FileDesc{file: file, scanner: bufio.NewScanner(file)}
}

func getAsteroidsFromFile(filepath string) []string {
	checkFile(filepath)
	var asteroidNumbers []string
	var fileDesc = getFileDesc(filepath)
	var scanner = fileDesc.scanner
	for scanner.Scan() {
		asteroidNumbers = append(asteroidNumbers, scanner.Text())
	}
	fileDesc.file.Close()
	return asteroidNumbers
}

func getAsteroidsDifference(webAsteroidNumbers, localAsteroidNumbers []string) []string {
	var diffAsteroidNumbers []string

	for _, webAst := range webAsteroidNumbers {
		if v, _ := strconv.Atoi(webAst); v > MAX_LOCAL_ASTEROID_NUMBER {
			continue
		}
		for _, localAst := range localAsteroidNumbers {
			if webAst == localAst {
				break
			}
		}
		fmt.Println("Asteroid " + webAst + " doesn't exist")
		diffAsteroidNumbers = append(diffAsteroidNumbers, webAst)
	}
	return diffAsteroidNumbers
}


type AxisInfo struct {
	axis float32
	axisDiff float32
}

func getAxisesFromCatalog(filepath string) map[string]float64 {
	checkFile(filepath)
	var fileDesc = getFileDesc(filepath)
	var scanner = fileDesc.scanner
	var count = 0
	var res map[string]float64
	for scanner.Scan() {
		if count < SKIP_LINE_COUNT {
			count++
			continue
		}
		var asteroidData = scanner.Text()
		var data = strings.Split(asteroidData, " ")
		var asteroidNumber = data[0][1:len(data[0])-1]
		var asteroidAxis = data[2]
		axis, err := strconv.ParseFloat(asteroidAxis, 32)
		if err != nil {
			log.Fatal("error during casting " + asteroidAxis)
		}
		res[asteroidNumber] = axis
	}
	fileDesc.file.Close()
	return res
}

func main() {
	var filepath = os.Args[1]
	var catalogFilepath = os.Args[2]
	url := "http://www.minorplanetcenter.org/iau/lists/JupiterTrojans.html"
	var webAsteroidNumbers = getTrojanAsteroidsFromWeb(url)
	var localAsteroidNumbers = getAsteroidsFromFile(filepath)
	var diffAsteroidNumbers = getAsteroidsDifference(webAsteroidNumbers, localAsteroidNumbers)
	var axises = getAxisesFromCatalog(catalogFilepath)

	for _, val := range diffAsteroidNumbers {
		var axis = axises[val]
		fmt.Println(val + " " + axis + " " + axis - RESONANCE_AXIS)
	}
}
