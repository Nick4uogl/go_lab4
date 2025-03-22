package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
)

// Структури та типи даних
type CableResults struct {
	NormalCurrent        string
	PostEmergencyCurrent string
	EconomicCrossSection string
	MinimumCrossSection  string
}

type ShortCircuitResults struct {
	ReactorImpedance           string
	TransformerImpedance       string
	TotalImpedance             string
	InitialShortCircuitCurrent string
}

type NetworkResults struct {
	ISh3     string
	ISh2     string
	ISh3Min  string
	ISh2Min  string
	IShN3    string
	IShN2    string
	IShN3Min string
	IShN2Min string
	ILN3     string
	ILN2     string
	ILN3Min  string
	ILN2Min  string
}

type Impedance struct {
	Resistance float64
	Reactance  float64
	Impedance  float64
}

// Функції для Impedance
func NewImpedance(resistance, reactance float64) Impedance {
	return Impedance{
		Resistance: resistance,
		Reactance:  reactance,
		Impedance:  math.Sqrt(math.Pow(resistance, 2) + math.Pow(reactance, 2)),
	}
}

func (imp Impedance) Transformed() Impedance {
	kpr := math.Pow(11.0, 2) / math.Pow(115.0, 2)
	return NewImpedance(imp.Resistance*kpr, imp.Reactance*kpr)
}

// Функція розрахунку параметрів кабелю
func calculateCableParameters(sm, ik, tf float64) CableResults {
	im := (sm / 2) / (math.Sqrt(3.0) * 10)
	imPa := 2 * im
	sEk := im / 1.4
	sVsS := (ik * math.Sqrt(tf)) / 92

	return CableResults{
		NormalCurrent:        fmt.Sprintf("%.1f", im),
		PostEmergencyCurrent: fmt.Sprintf("%.0f", imPa),
		EconomicCrossSection: fmt.Sprintf("%.1f", sEk),
		MinimumCrossSection:  fmt.Sprintf("%.0f", sVsS),
	}
}

// Функція розрахунку параметрів короткого замикання
func calculateShortCircuitParameters(sk float64) ShortCircuitResults {
	// Розрахунок імпедансу реактора
	xc := math.Pow(10.5, 2) / sk
	// Розрахунок імпедансу трансформатора
	xt := (10.5 / 100) * (math.Pow(10.5, 2) / 6.3)
	// Загальний імпеданс
	totalImpedance := xc + xt
	// Початковий струм трифазного короткого замикання
	initialSCCurrent := 10.5 / (math.Sqrt(3.0) * totalImpedance)

	return ShortCircuitResults{
		ReactorImpedance:           fmt.Sprintf("%.2f", xc),
		TransformerImpedance:       fmt.Sprintf("%.2f", xt),
		TotalImpedance:             fmt.Sprintf("%.2f", totalImpedance),
		InitialShortCircuitCurrent: fmt.Sprintf("%.1f", initialSCCurrent),
	}
}

// Допоміжні функції для розрахунку мережі
func calculateTransformerReactance() float64 {
	return (11.1 * math.Pow(115.0, 2)) / (100 * 6.3)
}

func calculateImpedances(resistance, reactance, transformerReactance float64) Impedance {
	return NewImpedance(resistance, reactance+transformerReactance)
}

func formatCurrent(current float64) string {
	return fmt.Sprintf("%.1f", current)
}

func calculateCurrents(voltage float64, normal, minimum Impedance) (string, string, string, string) {
	threePhaseNormal := formatCurrent(voltage / (math.Sqrt(3.0) * normal.Impedance))
	threePhaseNormalVal, _ := strconv.ParseFloat(threePhaseNormal, 64)
	twoPhaseNormal := formatCurrent(threePhaseNormalVal * (math.Sqrt(3.0) / 2))

	threePhaseMin := formatCurrent(voltage / (math.Sqrt(3.0) * minimum.Impedance))
	threePhaseMinVal, _ := strconv.ParseFloat(threePhaseMin, 64)
	twoPhaseMin := formatCurrent(threePhaseMinVal * (math.Sqrt(3.0) / 2))

	return threePhaseNormal, twoPhaseNormal, threePhaseMin, twoPhaseMin
}

func calculatePoint10Currents(normal, minimum Impedance) (string, string, string, string) {
	const lineResistance = 12.52 // Загальний опір лінії
	const lineReactance = 6.88   // Загальна реактивність лінії

	normalTotal := NewImpedance(normal.Resistance+lineResistance, normal.Reactance+lineReactance)
	minimumTotal := NewImpedance(minimum.Resistance+lineResistance, minimum.Reactance+lineReactance)

	return calculateCurrents(11.0, normalTotal, minimumTotal)
}

// Функція розрахунку мережі
func calculateNetwork(rsn, xsn, rsnMin, xsnMin float64) NetworkResults {
	xt := calculateTransformerReactance()
	normal := calculateImpedances(rsn, xsn, xt)
	minimum := calculateImpedances(rsnMin, xsnMin, xt)

	iSh3, iSh2, iSh3Min, iSh2Min := calculateCurrents(115.0, normal, minimum)
	iShN3, iShN2, iShN3Min, iShN2Min := calculateCurrents(11.0, normal.Transformed(), minimum.Transformed())
	iLN3, iLN2, iLN3Min, iLN2Min := calculatePoint10Currents(normal, minimum)

	return NetworkResults{
		ISh3:     iSh3,
		ISh2:     iSh2,
		ISh3Min:  iSh3Min,
		ISh2Min:  iSh2Min,
		IShN3:    iShN3,
		IShN2:    iShN2,
		IShN3Min: iShN3Min,
		IShN2Min: iShN2Min,
		ILN3:     iLN3,
		ILN2:     iLN2,
		ILN3Min:  iLN3Min,
		ILN2Min:  iLN2Min,
	}
}

// Обробники HTTP запитів
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func cableCalculatorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sm, _ := strconv.ParseFloat(r.FormValue("sm"), 64)
	ik, _ := strconv.ParseFloat(r.FormValue("ik"), 64)
	tf, _ := strconv.ParseFloat(r.FormValue("tf"), 64)

	results := calculateCableParameters(sm, ik, tf)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{
		"normalCurrent": "%s",
		"postEmergencyCurrent": "%s",
		"economicCrossSection": "%s",
		"minimumCrossSection": "%s"
	}`, results.NormalCurrent, results.PostEmergencyCurrent,
		results.EconomicCrossSection, results.MinimumCrossSection)
}

func shortCircuitCalculatorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sk, _ := strconv.ParseFloat(r.FormValue("sk"), 64)

	results := calculateShortCircuitParameters(sk)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{
		"reactorImpedance": "%s",
		"transformerImpedance": "%s",
		"totalImpedance": "%s",
		"initialShortCircuitCurrent": "%s"
	}`, results.ReactorImpedance, results.TransformerImpedance,
		results.TotalImpedance, results.InitialShortCircuitCurrent)
}

func networkCalculatorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rsn, _ := strconv.ParseFloat(r.FormValue("rsn"), 64)
	xsn, _ := strconv.ParseFloat(r.FormValue("xsn"), 64)
	rsnMin, _ := strconv.ParseFloat(r.FormValue("rsnMin"), 64)
	xsnMin, _ := strconv.ParseFloat(r.FormValue("xsnMin"), 64)

	results := calculateNetwork(rsn, xsn, rsnMin, xsnMin)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{
		"iSh3": "%s",
		"iSh2": "%s",
		"iSh3Min": "%s",
		"iSh2Min": "%s",
		"iShN3": "%s",
		"iShN2": "%s",
		"iShN3Min": "%s",
		"iShN2Min": "%s",
		"iLN3": "%s",
		"iLN2": "%s",
		"iLN3Min": "%s",
		"iLN2Min": "%s"
	}`, results.ISh3, results.ISh2, results.ISh3Min, results.ISh2Min,
		results.IShN3, results.IShN2, results.IShN3Min, results.IShN2Min,
		results.ILN3, results.ILN2, results.ILN3Min, results.ILN2Min)
}

func main() {
	// Статичні файли
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Маршрутизація
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/calculate/cable", cableCalculatorHandler)
	http.HandleFunc("/calculate/shortcircuit", shortCircuitCalculatorHandler)
	http.HandleFunc("/calculate/network", networkCalculatorHandler)

	// Запуск сервера
	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}