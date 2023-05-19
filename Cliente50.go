package main

import (
	"bufio"
	"fmt"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type TCPClient50 struct {
	serverMsg        string
	SERVERIP         string
	SERVERPORT       int
	mMessageListener func(string)
	mRun             bool
	out              *bufio.Writer
	in               *bufio.Reader
}

func NewTCPClient50(ip string, listener func(string)) *TCPClient50 {
	return &TCPClient50{
		SERVERIP:         ip,
		SERVERPORT:       4444,
		mMessageListener: listener,
	}
}

func (c *TCPClient50) sendMessage(message string) {
	if c.out != nil {
		c.out.WriteString(message + "\n")
		c.out.Flush()
	}
}

func (c *TCPClient50) stopClient() {
	c.mRun = false
}

func (c *TCPClient50) run() {
	c.mRun = true
	serverAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", c.SERVERIP, c.SERVERPORT))
	if err != nil {
		fmt.Println("TCP: Error", err)
		return
	}

	fmt.Println("TCP Client: Conectando...")
	socket, err := net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		fmt.Println("TCP: Error", err)
		return
	}

	c.out = bufio.NewWriter(socket)
	c.in = bufio.NewReader(socket)

	defer socket.Close()

	fmt.Println("TCP Client: Sent.")
	fmt.Println("TCP Client: Done.")

	for c.mRun {
		serverMsg, err := c.in.ReadString('\n')
		if err != nil {
			fmt.Println("TCP S: Error", err)
			return
		}

		serverMsg = strings.TrimSpace(serverMsg)
		if serverMsg != "" && c.mMessageListener != nil {
			c.mMessageListener(serverMsg)
		}
	}
}

func main() {
	client := Cliente50{}
	client.iniciar()
}

type Cliente50 struct {
	sum        [40]float64
	mTcpClient *TCPClient50
	sc         *bufio.Scanner
}

func (c *Cliente50) iniciar() {
	go func() {
		c.mTcpClient = NewTCPClient50("127.0.0.1", func(message string) {
			c.ClienteRecibe(message)
		})
		c.mTcpClient.run()
	}()

	salir := "n"
	c.sc = bufio.NewScanner(os.Stdin)
	for salir != "s" {
		c.sc.Scan()
		salir = c.sc.Text()
		c.ClienteEnvia(salir)
	}
}

func (c *Cliente50) ClienteRecibe(llego string) {
	fmt.Println("CLINTE50 El mensaje:", llego)
	if strings.Contains(llego, "evalua") {
		arrayString := strings.Fields(llego)
		funcion := arrayString[2]
		n := parseInt(arrayString[8])
		min := parseInt(arrayString[6])
		max := parseInt(arrayString[7])
		dif := (parseFloat(arrayString[4]) - parseFloat(arrayString[3])) / float64(n)
		c.procesar(min, max, funcion, n, dif)
	}
}

func (c *Cliente50) ClienteEnvia(envia string) {
	if c.mTcpClient != nil {
		c.mTcpClient.sendMessage(envia)
	}
}

func (c *Cliente50) procesar(a int, b int, funcion string, n int, dif float64) {
	start := time.Now()
	var wg sync.WaitGroup
	poli := NewEvaluadorPolinomios(funcion)
	N := b - a
	H := 2
	d := N / H
	todos := make([]*tarea0101, H)

	for i := 0; i < H-1; i++ {
		wg.Add(1)
		todos[i] = &tarea0101{
			min:  float64(i*d + a),
			max:  float64(i*d + d + a),
			id:   i,
			poli: poli,
			n:    n,
			dif:  dif,
		}
		go todos[i].run(&wg)
	}
	wg.Add(1)
	todos[H-1] = &tarea0101{
		min:  float64((d * (H - 1)) + a),
		max:  float64(b),
		id:   H - 1,
		poli: poli,
		n:    n,
		dif:  dif,
	}
	go todos[H-1].run(&wg)

	for i := 0; i <= H-1; i++ {
		todos[i].wait()
	}
	wg.Wait()

	sumatotal := 0.0
	for i := 0; i < H; i++ {
		sumatotal += todos[i].sum[i]
	}
	timeElapsed := time.Since(start)
	fmt.Println("Tiempo: ", timeElapsed, "SUMA TOTAL____:", sumatotal)
	c.ClienteEnvia(fmt.Sprintf("rpta %f", sumatotal))
}

type tarea0101 struct {
	max  float64
	min  float64
	n    int
	dif  float64
	id   int
	poli EvaluadorPolinomios
	sum  [40]float64
}

func (t *tarea0101) run(wg *sync.WaitGroup) {
	defer wg.Done()
	integral := 0.0

	for i := 0; i < int((t.max-t.min)/t.dif); i++ {
		XD := t.min + float64(i)*t.dif
		integral += t.poli.evaluarPolinomio(XD) * t.dif
		//fmt.Printf("Hilo: %d Entrada: %f Salida: %f n: %d max: %f min: %f dif: %f\n", t.id, XD, t.poli.evaluarPolinomio(XD), t.n, t.max, t.min, t.dif)
	}
	//fmt.Printf("%f", integral)
	t.sum[t.id] = integral
}

func (t *tarea0101) wait() {
	// no-op
}

type EvaluadorPolinomios struct {
	terminos []string
}

func NewEvaluadorPolinomios(polinomio string) EvaluadorPolinomios {
	return EvaluadorPolinomios{terminos: strings.Split(polinomio, "+")}
}

func (e *EvaluadorPolinomios) evaluarPolinomio(x float64) float64 {
	resultado := 0.0

	for _, termino := range e.terminos {
		partes := strings.Split(strings.TrimSpace(termino), "^")
		coeficiente := parseFloat(strings.Replace(partes[0], "x", "", -1))
		exponente := parseInt(partes[1])
		resultado += coeficiente * math.Pow(x, float64(exponente))
	}

	return resultado
}

func parseInt(s string) int {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		fmt.Println("Error converting to int:", err)
		return 0
	}
	return int(val)
}

func parseFloat(s string) float64 {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		fmt.Println("Error converting to float64:", err)
		return 0.0
	}
	return val
}
