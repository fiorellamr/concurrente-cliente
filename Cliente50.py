import socket
import threading
import numpy as np
import time 
class TCPClient50:
    def __init__(self, ip, message_received_callback):
        self.SERVERIP = ip
        self.SERVERPORT = 4444
        self.mMessageListener = message_received_callback
        self.mRun = False
        self.out = None
        self.incoming = None

    def send_message(self, message):
        if self.out is not None and not self.out.closed:
            self.out.write(message + "\n")
            self.out.flush()

    def stop_client(self):
        self.mRun = False

    def run(self):
        self.mRun = True
        try:
            print("TCP Client" + "C: Conectando...")
            with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
                sock.connect((self.SERVERIP, self.SERVERPORT))
                self.out = sock.makefile(mode='w')
                print("TCP Client" + "C: Sent.")
                print("TCP Client" + "C: Done.")
                self.incoming = sock.makefile(mode='r')
                while self.mRun:
                    server_message = self.incoming.readline()
                    if server_message and self.mMessageListener:
                        self.mMessageListener(server_message.rstrip("\n"))
        except Exception as e:
            print("TCP" + "C: Error", e)


class Cliente50:
    sum = [None] *40
    def __init__(self):
        self.mTcpClient = None
        self.sc = None

    def main(self):
        objcli = Cliente50()
        objcli.iniciar()

    def iniciar(self):
        def on_message_received(message):
            self.cliente_recibe(message)

        def tcp_client_thread():
            self.mTcpClient = TCPClient50("127.0.0.1", on_message_received)
            self.mTcpClient.run()

        # Iniciar el hilo del cliente TCP
        thread = threading.Thread(target=tcp_client_thread)
        thread.start()

        salir = "n"
        print("Cliente bandera 01")
        while salir != "s":
            salir = input()
            self.cliente_envia(salir)
        print("Cliente bandera 02")

    def cliente_recibe(self, llego):
        print("CLINTE50 El mensaje::" + llego)
        if "evalua" in llego:
            arrayString = llego.split()
            funcion = arrayString[2]
            n = int(arrayString[5])
            min = int(arrayString[6])
            max = int(arrayString[7])
            dif = (float(arrayString[4]) - float(arrayString[3]))/float(n)
            #print("funcion: ", funcion, "n: ", n, " min: ", min, " max: ",max, " dif: ",dif)
            self.procesar(min,max,funcion,n, dif)

    def cliente_envia(self, envia):
        if self.mTcpClient is not None:
            self.mTcpClient.send_message(envia)
    
    class EvaluadorPolinomios:
        def __init__(self, polinomio):
            self.terminos = polinomio.split("+")

        def evaluarPolinomio(self, x):
            resultado = 0
            for termino in self.terminos:
                partes = termino.strip().split("^")
                coeficiente = float(partes[0].replace("x", "").strip())
                exponente = int(partes[1].strip())
                resultado += coeficiente * x ** exponente
            return resultado

    def procesar(self, a, b, funcion, n,dif):
        inicio = time.time()
        poli = self.EvaluadorPolinomios(funcion)
        N = b - a
        H = 4
        d = N / H

        todos = [None] * 40
        for i in range(H - 1):
            todos[i] = self.tarea0101(i * d + a, i * d + d + a, i, poli, n,dif)
            todos[i].start()

        todos[H - 1] = self.tarea0101(d * (H - 1) + a, b, H - 1, poli, n,dif)
        todos[H - 1].start()
        
        for i in range(H):
            todos[i].join()

        sumatotal = sum(self.sum[:H])
        fin = time.time()
        print("Tiempo: ",fin-inicio,"SUMA TOTAL____:", sumatotal)
        self.cliente_envia("rpta " + str(sumatotal))


    class tarea0101(threading.Thread):
        def __init__(self, min_, max_, id_, poli_, n_,dif_):
            threading.Thread.__init__(self)
            self.max = max_
            self.min = min_
            self.id = id_
            self.poli = poli_
            self.n = n_
            self.dif = dif_

        def run(self):
            integral = 0
            #print(self.min, self.max)
            for i in np.arange(0,(self.max - self.min)/self.dif,1.0):
                XD = self.min + i*self.dif
                integral += self.poli.evaluarPolinomio(XD) * self.dif
                #print("Hilo: ",self.id," Entrada: ",XD,"Salida: ",self.poli.evaluarPolinomio(XD)," n: ",self.n,"max",self.max,"min",self.min,"dif",self.dif)
            Cliente50.sum[self.id] = integral
            #print(" min:", self.min, " max:", self.max, " id:", self.id, " integral:", integral)
            #print("\n")
    class OnMessageReceived:
        def messageReceived(self, message):
            pass
    
    

    def tcp_client_thread(self):
        self.mTcpClient = TCPClient50("127.0.0.1", self.OnMessageReceived())
        self.mTcpClient.run()


# Llamada al método main
if __name__ == "__main__":
    objcli = Cliente50()
    objcli.iniciar()
