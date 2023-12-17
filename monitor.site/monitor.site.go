package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	quantidadeDeTestes      = 4
	tempoDeEsperaEmSegundos = 5
	nomeDoArquivoDeSitesTXT = "sites.txt"
	nomeDoArquivoLOG        = "sites.log"
	permissaoNoArquivoLOG   = 0666
	flagDoArquivoLOG        = os.O_CREATE | os.O_RDWR | os.O_APPEND
)

func main() {
	exibirIntroducao()
	for {
		exibirMenuDeComandos()
		comando := lerComandoDoUsuario()
		executarComando(comando)
	}
}

func exibirIntroducao() {
	fmt.Println()
	nome := "Jhoni Conzatti"
	var versao float32 = 1.2
	fmt.Println("Olá sr.", nome)
	fmt.Println("Este programa está na versão", versao)
	fmt.Println()
}

func exibirMenuDeComandos() {
	fmt.Println()
	fmt.Println("Escolha um comando:")
	fmt.Println("1 - Iniciar monitoramento")
	fmt.Println("2 - Exibir logs")
	fmt.Println("0 - Sair do programa")
	fmt.Println()
}

func executarComando(comando int) {
	//if comando == 1 {
	//	fmt.Println("Monitorando...")
	//} else if comando == 2 {
	//	fmt.Println("Exibindo logs...")
	//} else if comando == 0 {
	//	fmt.Println("Saindo...")
	//} else {
	//	fmt.Println("Não conheço este comando")
	//}

	switch comando {
	case 1:
		iniciarMonitoramento()
	case 2:
		exibirLogs()
	case 0:
		sairDoProgramaComSucesso()
	default:
		sairDoProgramaPorComandoNaoReconhecido()
	}
}

func lerComandoDoUsuario() int {
	var comandoLido int
	//fmt.Scanf("%d", &comandoLido)
	_, erro := fmt.Scan(&comandoLido)
	//fmt.Println("O endereço da variável comando é", &comandoLido)
	if erro != nil {
		comandoLido = 99
		fmt.Println("O comando escolhido inválido!", erro.Error())
	} else {
		fmt.Println("O comando escolhido foi", comandoLido)
	}
	fmt.Println()
	return comandoLido
}

func iniciarMonitoramento() {
	fmt.Println("Iniciando monitoramento")

	//Arrays
	//var sites [4]string
	//sites[0] = "https://httpbin.org/status/200"
	//sites[1] = "https://httpbin.org/status/404"
	//sites[2] = "https://www.alura.com.br"
	//sites[3] = "https://bluwaysistemas.ddns.net"

	//sites := [4]string{"https://httpbin.org/status/200", "https://httpbin.org/status/404", "https://www.alura.com.br", "https://bluwaysistemas.ddns.net"}

	//Slices
	//var sites []string
	//sites = append(sites, "https://httpbin.org/status/200")
	//sites = append(sites, "https://httpbin.org/status/404")
	//sites = append(sites, "https://www.alura.com.br")
	//sites = append(sites, "https://bluwaysistemas.ddns.net")

	//sites := []string{"https://httpbin.org/status/200", "https://httpbin.org/status/404", "https://www.alura.com.br", "https://bluwaysistemas.ddns.net", "ds"}

	sites := obterSitesDoArquivo(nomeDoArquivoDeSitesTXT)

	//fmt.Printf("Tipo: %s, tamanho: %d elementos, capacidade: %d\n", reflect.TypeOf(sites), len(sites), cap(sites))
	//fmt.Println("Sites:", sites)

	monitorarSites(sites, quantidadeDeTestes, tempoDeEsperaEmSegundos)
}

func obterSitesDoArquivo(nomeDoArquivo string) []string {
	//conteudoDoArquivo, erro := os.ReadFile(arquivoTXT)
	//if conteudoDoArquivo != nil {
	//	fmt.Println(string(conteudoDoArquivo))
	//}
	var sites []string
	arquivoDeSitesTXT, erro := os.Open(nomeDoArquivo)
	if erro != nil {
		fmt.Printf("Erro ao abrir arquivo %s! %s\n", nomeDoArquivo, erro.Error())
	}
	if arquivoDeSitesTXT != nil {
		leitorDoArquivo := bufio.NewReader(arquivoDeSitesTXT)
		for {
			linhaDoArquivo, erro := leitorDoArquivo.ReadString('\n')
			if erro != nil {
				if erro != io.EOF {
					fmt.Printf("Erro ao ler linha do arquivo %s! %s\n", nomeDoArquivo, erro.Error())
				}
				break
			}
			linhaDoArquivo = strings.TrimSpace(linhaDoArquivo)
			if linhaDoArquivo != "" {
				sites = append(sites, linhaDoArquivo)
			}
		}
		arquivoDeSitesTXT.Close()
	}
	return sites
}

func monitorarSites(sites []string, quantidadeDeTestes uint, tempoDeEsperaEmSegundos time.Duration) {
	if len(sites) > 0 {
		fmt.Printf("Monitorando %d sites, %d vezes a cada %d segundos\n", len(sites), quantidadeDeTestes, tempoDeEsperaEmSegundos)
		for i := 0; i < int(quantidadeDeTestes); i++ {
			fmt.Printf("Monitoramento %d de %d\n", i+1, quantidadeDeTestes)
			fmt.Println("-----------------------------------------------")
			for j, site := range sites {
				fmt.Printf("Testando site %d: %s\n", j+1, site)
				testarSite(site)
			}
			if i < int(quantidadeDeTestes)-1 {
				time.Sleep(tempoDeEsperaEmSegundos * time.Second)
			}
			fmt.Println()
		}
		fmt.Println("Fim do monitoramento de sites.")
	} else {
		fmt.Println("Nenhum site encontrado para monitorar!")
	}
}

func testarSite(site string) {
	resposta, erro := http.Get(site)
	var mensagem string
	online := false
	if erro != nil {
		mensagem = fmt.Sprintf("Site \"%s\" está com problemas. Erro: %s", site, erro.Error())
	}
	if resposta != nil {
		online = resposta.StatusCode == 200
		if online {
			mensagem = fmt.Sprintf("Site \"%s\" foi carregado com sucesso!", site)
		} else {
			mensagem = fmt.Sprintf("Site \"%s\" está com problemas. Status Code: %d", site, resposta.StatusCode)
		}
	}
	registrarLog(site, online, mensagem)
	fmt.Println(mensagem)
	fmt.Println("-----------------------------------------------")
}

func registrarLog(site string, online bool, mensagem string) {
	arquivoLog, erro := os.OpenFile(nomeDoArquivoLOG, flagDoArquivoLOG, permissaoNoArquivoLOG)
	if erro != nil {
		fmt.Printf("Erro ao abrir arquivo %s! %s\n", arquivoLog.Name(), erro.Error())
	}
	if arquivoLog != nil {
		arquivoLog.WriteString(fmt.Sprintf("%s - %s - online: %t. %s\n", time.Now().Format("02/01/2006 15:04:05"), site, online, mensagem))
		arquivoLog.Close()
	}
}

func exibirLogs() {
	arquivoLog, erro := os.ReadFile(nomeDoArquivoLOG)
	if erro != nil {
		fmt.Printf("Erro ao ler arquivo %s! %s\n", nomeDoArquivoLOG, erro.Error())
	}
	if arquivoLog != nil {
		fmt.Print(string(arquivoLog))
	}
}

func sairDoProgramaComSucesso() {
	fmt.Println("Saindo...")
	fmt.Println()
	os.Exit(0)
}

func sairDoProgramaPorComandoNaoReconhecido() {
	fmt.Println("Comando não reconhecido!")
	fmt.Println()
	os.Exit(1)
}
