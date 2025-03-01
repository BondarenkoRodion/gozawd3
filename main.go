package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

// integral обчислює визначений інтеграл методом Сімпсона
func integral(f func(float64) float64, a, b float64, n int) float64 {
	h := (b - a) / float64(n)
	sum := f(a) + f(b)

	for i := 1; i < n; i++ {
		x := a + float64(i)*h
		if i%2 == 0 {
			sum += 2 * f(x)
		} else {
			sum += 4 * f(x)
		}
	}

	return sum * (h / 3)
}

// HTML-шаблон
var tmpl = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html lang="uk">
<head>
    <meta charset="UTF-8">
    <title>Калькулятор прибутку</title>
</head>
<body>
    <h2>Розрахунок прибутку</h2>
    <form method="POST">
        <label>Pc: <input type="number" step="any" name="Pc" required></label><br>
        <label>До покращення (s1): <input type="number" step="any" name="s1" required></label><br>
        <label>Після покращення (s2): <input type="number" step="any" name="s2" required></label><br>
        <label>B: <input type="number" step="any" name="B" required></label><br>
        <button type="submit">Обчислити</button>
    </form>
    {{if .}}
        <h3>Результати:</h3>
        <p>Прибуток до покращення: {{.Before}}</p>
        <p>Прибуток після покращення: {{.After}}</p>
    {{end}}
</body>
</html>
`))

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		Pc, _ := strconv.ParseFloat(r.FormValue("Pc"), 64)
		s1, _ := strconv.ParseFloat(r.FormValue("s1"), 64)
		s2, _ := strconv.ParseFloat(r.FormValue("s2"), 64)
		B, _ := strconv.ParseFloat(r.FormValue("B"), 64)

		dW1 := integral(func(p float64) float64 {
			return (1 / (s1 * math.Sqrt(2*math.Pi))) * math.Exp(-math.Pow(p-Pc, 2)/(2*math.Pow(s1, 2)))
		}, 4.75, 5.25, 1000000)
		W1 := Pc * 24 * dW1
		Pryb1 := W1 * B
		W21 := Pc * 24 * (1 - dW1)
		Sch1 := W21 * B
		Riznytsia1 := Pryb1 - Sch1

		dW2 := integral(func(p float64) float64 {
			return (1 / (s2 * math.Sqrt(2*math.Pi))) * math.Exp(-math.Pow(p-Pc, 2)/(2*math.Pow(s2, 2)))
		}, 4.75, 5.25, 1000000)
		W12 := Pc * 24 * dW2
		Pryb2 := W12 * B
		W22 := Pc * 24 * (1 - dW2)
		Sch2 := W22 * B
		Riznytsia2 := Pryb2 - Sch2

		tmpl.Execute(w, map[string]float64{"Before": Riznytsia1, "After": Riznytsia2})
		return
	}
	tmpl.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Сервер запущено на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
