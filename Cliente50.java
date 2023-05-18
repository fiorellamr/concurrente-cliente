package sumaNred;

import java.util.Scanner;
import java.util.*;
import sumaNred.TCPClient50;

class Cliente50 {
    public double sum[] = new double[40];
    TCPClient50 mTcpClient;
    Scanner sc;

    public static void main(String[] args) {
        Cliente50 objcli = new Cliente50();
        objcli.iniciar();
    }

    void iniciar() {
        new Thread(
                new Runnable() {
                    @Override
                    public void run() {
                        mTcpClient = new TCPClient50("127.0.0.1",
                                new TCPClient50.OnMessageReceived() {
                                    @Override
                                    public void messageReceived(String message) {
                                        ClienteRecibe(message);
                                    }
                                });
                        mTcpClient.run();
                    }
                }).start();

        String salir = "n";
        sc = new Scanner(System.in);
        // System.out.println("Cliente bandera 01");
        while (!salir.equals("s")) {
            salir = sc.nextLine();
            ClienteEnvia(salir);
        }
        // System.out.println("Cliente bandera 02");

    }

    void ClienteRecibe(String llego) {
        System.out.println("CLINTE50 El mensaje::" + llego);
        if (llego.trim().contains("evalua")) {
            String arrayString[] = llego.split("\\s+");
            String funcion = arrayString[2];
            int n = Integer.parseInt(arrayString[8]);
            int min = Integer.parseInt(arrayString[6]);
            int max = Integer.parseInt(arrayString[7]);
            float dif = (Float.parseFloat(arrayString[4]) - Float.parseFloat(arrayString[3])) / n;
            procesar(min, max, funcion, n, dif);
        }
    }

    void ClienteEnvia(String envia) {
        if (mTcpClient != null) {
            mTcpClient.sendMessage(envia);
        }
    }

    void procesar(float a, float b, String funcion, int n, float dif) {
        EvaluadorPolinomios poli = new EvaluadorPolinomios(funcion);
        float N = (b - a);
        int H = 2;// luego aumentar
        float d = (float) ((N) / H);
        // n = n/H;
        // INSTANCIA LAS TAREAS EN HILOS E INICIALIZA LOS HILOS
        Thread todos[] = new Thread[40];
        for (int i = 0; i < (H - 1); i++) {
            // System.out.println("a:" + (i * d + a) + "b" + (i * d + d + a) + " i" + i);
            todos[i] = new tarea0101((i * d + a), (i * d + d + a), i, poli, n, dif);
            todos[i].start();
        }
        // System.out.println("a" + ((d * (H - 1)) + a) + "b" + (b + 1) + " i" + (H -
        // 1));

        todos[H - 1] = new tarea0101(((d * (H - 1)) + a), b, (H - 1), poli, n, dif);
        todos[H - 1].start();

        for (int i = 0; i <= (H - 1); i++) {
            try {
                todos[i].join();
            } catch (InterruptedException ex) {
                System.out.println("error" + ex);
            }
        }

        double sumatotal = 0;
        for (int i = 0; i < H; i++) {
            sumatotal = sumatotal + sum[i];
        }
        System.out.println("SUMA TOTAL____:" + sumatotal);
        ClienteEnvia("rpta " + sumatotal);
    }

    public class tarea0101 extends Thread {
        public float max, min, n, dif;
        public int id;
        public EvaluadorPolinomios poli;

        tarea0101(float min_, float max_, int id_, EvaluadorPolinomios poli_, float n_, float dif_) {
            max = max_;
            min = min_;
            id = id_;
            poli = poli_;
            n = n_;
            dif = dif_;
        }

        public void run() {
            double integral = 0;

            for (int i = 0; i < (max - min) / dif; i++) {
                float XD = min + (i) * dif;
                integral = integral + poli.evaluarPolinomio(XD) * dif;
                System.out.println("Hilo: " + id + "Entrada: " + XD + "Salida: " + poli.evaluarPolinomio(XD) + "n: " + n
                        + "max: " + max + "min: " + min + "dif: " + dif);
            }
            sum[id] = integral;

        }
    }

}

class EvaluadorPolinomios {
    public List<String> terminos;

    EvaluadorPolinomios(String polinomio) {
        terminos = Arrays.asList(polinomio.split("\\+"));
    }

    public double evaluarPolinomio(double x) {
        double resultado = 0;
        for (String termino : terminos) {
            String[] partes = termino.trim().split("\\^");
            double coeficiente = Double.parseDouble(partes[0].replace("x", "").trim());
            int exponente = Integer.parseInt(partes[1].trim());
            resultado += coeficiente * Math.pow(x, exponente);
        }
        return resultado;
    }
}
