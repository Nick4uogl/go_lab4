package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Println("\nPower System Calculator")
		fmt.Println("1. Cable Calculator")
		fmt.Println("2. Short Circuit Calculator")
		fmt.Println("3. Network Calculator")
		fmt.Println("4. Exit")
		fmt.Print("Select an option (1-4): ")
		
		scanner.Scan()
		option := scanner.Text()
		
		switch option {
		case "1":
			runCableCalculator(scanner)
		case "2":
			runShortCircuitCalculator(scanner)
		case "3":
			runNetworkCalculator(scanner)
		case "4":
			fmt.Println("Exiting program.")
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}

// Cable Calculator implementation
func runCableCalculator(scanner *bufio.Scanner) {
	fmt.Println("\n=== Cable Calculator ===")
	
	// Default values
	smDefault := "1300"
	ikDefault := "2500"
	tfDefault := "2.5"
	
	// Get input values
	sm := getInput(scanner, fmt.Sprintf("Enter Sm (MVA) [default: %s]: ", smDefault), smDefault)
	ik := getInput(scanner, fmt.Sprintf("Enter Ik (A) [default: %s]: ", ikDefault), ikDefault)
	tf := getInput(scanner, fmt.Sprintf("Enter tf (s) [default: %s]: ", tfDefault), tfDefault)
	
	// Parse input values
	smVal, _ := strconv.ParseFloat(sm, 64)
	ikVal, _ := strconv.ParseFloat(ik, 64)
	tfVal, _ := strconv.ParseFloat(tf, 64)
	
	// Calculate results
	results := calculateCableParameters(smVal, ikVal, tfVal)
	
	// Display results
	fmt.Println("\nResults:")
	fmt.Printf("Normal mode current: %s A\n", results.normalCurrent)
	fmt.Printf("Post-emergency current: %s A\n", results.postEmergencyCurrent)
	fmt.Printf("Economic cross-section: %s mm²\n", results.economicCrossSection)
	fmt.Printf("Minimum cross-section: %s mm²\n", results.minimumCrossSection)
}

// Short Circuit Calculator implementation
func runShortCircuitCalculator(scanner *bufio.Scanner) {
	fmt.Println("\n=== Short Circuit Calculator ===")
	
	// Default value
	skDefault := "200"
	
	// Get input value
	sk := getInput(scanner, fmt.Sprintf("Enter Short-Circuit Power (Sk) [MVA] [default: %s]: ", skDefault), skDefault)
	
	// Parse input value
	skVal, _ := strconv.ParseFloat(sk, 64)
	
	// Calculate results
	results := calculateShortCircuitParameters(skVal)
	
	// Display results
	fmt.Println("\nResults:")
	fmt.Printf("Xc: %s\n", results.reactorImpedance)
	fmt.Printf("Xt: %s\n", results.transformerImpedance)
	fmt.Printf("Total Resistance: %s\n", results.totalImpedance)
	fmt.Printf("Initial Three-Phase SC Current: %s\n", results.initialShortCircuitCurrent)
}

// Network Calculator implementation
func runNetworkCalculator(scanner *bufio.Scanner) {
	fmt.Println("\n=== Power Network Calculator ===")
	
	// Default values
	rsnDefault := "10.65"
	xsnDefault := "24.02"
	rsnMinDefault := "34.88"
	xsnMinDefault := "65.68"
	
	// Get input values
	rsn := getInput(scanner, fmt.Sprintf("Enter Rsn (Ω) [default: %s]: ", rsnDefault), rsnDefault)
	xsn := getInput(scanner, fmt.Sprintf("Enter Xsn (Ω) [default: %s]: ", xsnDefault), xsnDefault)
	rsnMin := getInput(scanner, fmt.Sprintf("Enter Rsn min (Ω) [default: %s]: ", rsnMinDefault), rsnMinDefault)
	xsnMin := getInput(scanner, fmt.Sprintf("Enter Xsn min (Ω) [default: %s]: ", xsnMinDefault), xsnMinDefault)
	
	// Parse input values
	rsnVal, _ := strconv.ParseFloat(rsn, 64)
	xsnVal, _ := strconv.ParseFloat(xsn, 64)
	rsnMinVal, _ := strconv.ParseFloat(rsnMin, 64)
	xsnMinVal, _ := strconv.ParseFloat(xsnMin, 64)
	
	// Calculate results
	results := calculateNetwork(rsnVal, xsnVal, rsnMinVal, xsnMinVal)
	
	// Display results
	fmt.Println("\n110kV bus SC currents (normal/minimum):")
	fmt.Printf("Three-phase: %s/%s A\n", results.iSh3, results.iSh3Min)
	fmt.Printf("Two-phase: %s/%s A\n", results.iSh2, results.iSh2Min)
	
	fmt.Println("\n10kV bus SC currents (normal/minimum):")
	fmt.Printf("Three-phase: %s/%s A\n", results.iShN3, results.iShN3Min)
	fmt.Printf("Two-phase: %s/%s A\n", results.iShN2, results.iShN2Min)
	
	fmt.Println("\nPoint 10 SC currents (normal/minimum):")
	fmt.Printf("Three-phase: %s/%s A\n", results.iLN3, results.iLN3Min)
	fmt.Printf("Two-phase: %s/%s A\n", results.iLN2, results.iLN2Min)
}

// Helper function to get user input with default value
func getInput(scanner *bufio.Scanner, prompt string, defaultValue string) string {
	fmt.Print(prompt)
	scanner.Scan()
	value := scanner.Text()
	if value == "" {
		return defaultValue
	}
	return value
}

// Cable Calculator structures and functions
type CableResults struct {
	normalCurrent       string
	postEmergencyCurrent string
	economicCrossSection string
	minimumCrossSection  string
}

func calculateCableParameters(sm float64, ik float64, tf float64) CableResults {
	im := (sm / 2) / (math.Sqrt(3.0) * 10)
	imPa := 2 * im
	sEk := im / 1.4
	sVsS := (ik * math.Sqrt(tf)) / 92
	
	return CableResults{
		normalCurrent:       fmt.Sprintf("%.1f", im),
		postEmergencyCurrent: fmt.Sprintf("%.0f", imPa),
		economicCrossSection: fmt.Sprintf("%.1f", sEk),
		minimumCrossSection:  fmt.Sprintf("%.0f", sVsS),
	}
}

// Short Circuit Calculator structures and functions
type ShortCircuitResults struct {
	reactorImpedance          string
	transformerImpedance      string
	totalImpedance            string
	initialShortCircuitCurrent string
}

func calculateShortCircuitParameters(sk float64) ShortCircuitResults {
	// Reactor impedance calculation
	xc := math.Pow(10.5, 2) / sk
	// Transformer impedance calculation
	xt := (10.5 / 100) * (math.Pow(10.5, 2) / 6.3)
	// Total impedance
	totalImpedance := xc + xt
	// Initial three-phase short-circuit current
	initialSCCurrent := 10.5 / (math.Sqrt(3.0) * totalImpedance)
	
	return ShortCircuitResults{
		reactorImpedance:          fmt.Sprintf("%.2f", xc),
		transformerImpedance:      fmt.Sprintf("%.2f", xt),
		totalImpedance:            fmt.Sprintf("%.2f", totalImpedance),
		initialShortCircuitCurrent: fmt.Sprintf("%.1f", initialSCCurrent),
	}
}

// Network Calculator structures and functions
type Impedance struct {
	resistance float64
	reactance  float64
	impedance  float64
}

func NewImpedance(resistance float64, reactance float64) Impedance {
	return Impedance{
		resistance: resistance,
		reactance:  reactance,
		impedance:  math.Sqrt(math.Pow(resistance, 2) + math.Pow(reactance, 2)),
	}
}

func (imp Impedance) Transformed() Impedance {
	kpr := math.Pow(11.0, 2) / math.Pow(115.0, 2)
	return NewImpedance(imp.resistance*kpr, imp.reactance*kpr)
}

type Currents struct {
	threePhaseNormal string
	twoPhaseNormal   string
	threePhaseMin    string
	twoPhaseMin      string
}

type NetworkResults struct {
	iSh3    string
	iSh2    string
	iSh3Min string
	iSh2Min string
	iShN3   string
	iShN2   string
	iShN3Min string
	iShN2Min string
	iLN3    string
	iLN2    string
	iLN3Min string
	iLN2Min string
}

func calculateTransformerReactance() float64 {
	return (11.1 * math.Pow(115.0, 2)) / (100 * 6.3)
}

func calculateImpedances(resistance float64, reactance float64, transformerReactance float64) Impedance {
	return NewImpedance(resistance, reactance+transformerReactance)
}

func formatCurrent(current float64) string {
	return fmt.Sprintf("%.1f", current)
}

func calculateCurrents(voltage float64, normal Impedance, minimum Impedance) Currents {
	threePhaseNormal := formatCurrent(voltage / (math.Sqrt(3.0) * normal.impedance))
	threePhaseNormalVal, _ := strconv.ParseFloat(threePhaseNormal, 64)
	twoPhaseNormal := formatCurrent(threePhaseNormalVal * (math.Sqrt(3.0) / 2))
	
	threePhaseMin := formatCurrent(voltage / (math.Sqrt(3.0) * minimum.impedance))
	threePhaseMinVal, _ := strconv.ParseFloat(threePhaseMin, 64)
	twoPhaseMin := formatCurrent(threePhaseMinVal * (math.Sqrt(3.0) / 2))
	
	return Currents{
		threePhaseNormal: threePhaseNormal,
		twoPhaseNormal:   twoPhaseNormal,
		threePhaseMin:    threePhaseMin,
		twoPhaseMin:      twoPhaseMin,
	}
}

func calculatePoint10Currents(normal Impedance, minimum Impedance) Currents {
	const lineResistance = 12.52 // Total resistance of the line
	const lineReactance = 6.88   // Total reactance of the line

	normalTotal := NewImpedance(normal.resistance+lineResistance, normal.reactance+lineReactance)
	minimumTotal := NewImpedance(minimum.resistance+lineResistance, minimum.reactance+lineReactance)

	return calculateCurrents(11.0, normalTotal, minimumTotal)
}

func calculateNetwork(rsn float64, xsn float64, rsnMin float64, xsnMin float64) NetworkResults {
	xt := calculateTransformerReactance()
	normal := calculateImpedances(rsn, xsn, xt)
	minimum := calculateImpedances(rsnMin, xsnMin, xt)

	currents110kV := calculateCurrents(115.0, normal, minimum)
	currents10kV := calculateCurrents(11.0, normal.Transformed(), minimum.Transformed())
	currentsPoint10 := calculatePoint10Currents(normal, minimum)

	return NetworkResults{
		iSh3:    currents110kV.threePhaseNormal,
		iSh2:    currents110kV.twoPhaseNormal,
		iSh3Min: currents110kV.threePhaseMin,
		iSh2Min: currents110kV.twoPhaseMin,
		iShN3:   currents10kV.threePhaseNormal,
		iShN2:   currents10kV.twoPhaseNormal,
		iShN3Min: currents10kV.threePhaseMin,
		iShN2Min: currents10kV.twoPhaseMin,
		iLN3:    currentsPoint10.threePhaseNormal,
		iLN2:    currentsPoint10.twoPhaseNormal,
		iLN3Min: currentsPoint10.threePhaseMin,
		iLN2Min: currentsPoint10.twoPhaseMin,
	}
}