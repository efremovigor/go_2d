package main

import (
    "github.com/faiface/pixel"
    "math/rand"
    "time"
    "image/color"
    "golang.org/x/image/colornames"
)

//Индекс начала генерации
var indexOfBeginGenerate = "1|1"
//Ширина кубов влабиринта
var factor = 10
//Ширина лабиринта
var lengthX = 101
//Высота либиринта
var lengthY = 77
//Лист перекрестков(для генерации)
var crossways []string
var src = rand.NewSource(time.Now().UnixNano())
var r = rand.New(src)

//Индексы лабиринта
var collection = map[string]Square{}
//Индексы начала и конца лабиринта
var startIndex, finishIndex string

//Цвета либинта
var colorMap = map[int]color.RGBA{
    0: colornames.Black,
    1: colornames.White,
    2: colornames.Gold,
    3: colornames.Green,
}
//Куб
type Square struct {
    state int
    color color.RGBA
    rect  pixel.Rect
    x     int
    y     int
    pass  bool
}
//назначение кода цвета
func (s *Square) getColor() color.RGBA {
    return colorMap[s.state]
}

//Генерация матрицы лабиринта
func generate(lengthX int, lengthY int) {
    //Create blank
    collection = make(map[string]Square)
    for y := 0; y < lengthY; y++ {
        for x := 0; x < lengthX; x++ {
            rect := pixel.R(float64(x*factor), float64(y*factor), float64(factor+x*factor), float64(factor+y*factor))
            square := Square{x: x, y: y, rect: rect}
            if x+1 == lengthX || x == 0 || y+1 == lengthY || y == 0 || x%2 != 1 || y%2 != 1 {
                square.state = 0
            } else {
                square.state = 1
            }
            collection[getIndex(x, y)] = square
        }
    }
    //Create paths
    updateWay(indexOfBeginGenerate)

    for len(crossways) > 0 {
        crossway := crossways[0]
        crossways = crossways[1:]
        updateWay(crossway)
    }

    for _ ,item := range collection {
        item.pass = false
        item.save()
    }
    setStart()
    setFinish()
    return
}

//выбор точки старта
func setStart() {
    var x, y int
    if r.Intn(1) == 1 {
        x = oddRandom(lengthX)
        y = 0
    } else {
        x = 0
        y = oddRandom(lengthY)
    }
    startIndex = getIndex(x, y)
    element := collection[startIndex]
    element.state = 2
    element.pass = true
    element.save()
}

//Выбор точки финиша
func setFinish() {
    element := collection[startIndex]
    x := lengthX - (element.x + 1)
    y := lengthY - (element.y + 1)
    finishIndex = getIndex(x, y)
    element = collection[finishIndex]
    element.state = 3
    element.pass = true
    element.save()
}


func (s Square) save() {
    collection[getIndex(s.x,s.y)] = s
}

//Делает одну итерацию генерации(генерирует ветку маршрута)
func updateWay(indexWay string) {
    siblings := getSiblings(indexWay)
    if len(siblings) > 1 {
        crossways = append(crossways, indexWay)
    }
    for len(siblings) > 0 {
        if len(siblings) > 1 {
            crossways = append(crossways, indexWay)
        }
        if len(siblings) > 0 {
            nextIndex := siblings[r.Intn(len(siblings))]
            for _, index := range getMediator(collection[indexWay], collection[nextIndex]) {
                element := collection[index]
                element.state = 1
                element.pass = true
                element.save()
                indexWay = index
            }
        }
        siblings = getSiblings(indexWay)
    }
}

//получить сосе
func getSiblings(index string) (list []string) {
    rect := collection[index]
    if current, ok := collection[getIndex(rect.x, rect.y-2)]; ok && current.pass == false {
        list = append(list, getIndex(rect.x, rect.y-2))
    }
    if current, ok := collection[getIndex(rect.x, rect.y+2)]; ok && current.pass == false {
        list = append(list, getIndex(rect.x, rect.y+2))
    }
    if current, ok := collection[getIndex(rect.x-2, rect.y)]; ok && current.pass == false {
        list = append(list, getIndex(rect.x-2, rect.y))
    }
    if current, ok := collection[getIndex(rect.x+2, rect.y)]; ok && current.pass == false {
        list = append(list, getIndex(rect.x+2, rect.y))
    }
    return
}

func getMediator(square1 Square, square2 Square) (list []string) {
    for square1.x > square2.x {
        square1.x--
        list = append(list, getIndex(square1.x, square1.y))
    }

    for square1.x < square2.x {
        square1.x++
        list = append(list, getIndex(square1.x, square1.y))
    }

    for square1.y > square2.y {
        square1.y--
        list = append(list, getIndex(square1.x, square1.y))
    }

    for square1.y < square2.y {
        square1.y++
        list = append(list, getIndex(square1.x, square1.y))
    }
    return
}