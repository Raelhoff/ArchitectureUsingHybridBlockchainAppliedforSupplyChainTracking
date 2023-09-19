package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"
)

func main() {
	// Caminho do arquivo
	caminhoArquivo := "/home/specto/requisicoesMB/1EDGE/1/debug_timing_CreateDevice_1234_1_requisicao_trm.txt"

	// Inicializar listas para armazenar tempos
	temposMS := make([]float64, 0)  // Tempos em milissegundos
	temposS := make([]float64, 0)   // Tempos em segundos

	// Ler o arquivo
	bytesArquivo, err := ioutil.ReadFile(caminhoArquivo)
	if err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
		return
	}

	linhas := strings.Split(string(bytesArquivo), "\n")
	for _, linha := range linhas {
		if strings.Contains(linha, "Tempo de Resposta Médio:") {
			tempoStr := strings.TrimSpace(strings.Split(linha, ":")[1])
			if strings.HasSuffix(tempoStr, "ms") {
				tempo, err := strconv.ParseFloat(strings.TrimSuffix(tempoStr, "ms"), 64)
				if err == nil {
					temposMS = append(temposMS, tempo)
				}
			} else if strings.HasSuffix(tempoStr, "s") {
				tempo, err := strconv.ParseFloat(strings.TrimSuffix(tempoStr, "s"), 64)
				if err == nil {
					temposS = append(temposS, tempo)
				}
			}
		}
	}

	// Calcular a média
	mediaMS := calcularMedia(temposMS)
	mediaS := calcularMedia(temposS)

	fmt.Printf("Tempo de Resposta Médio Total em Milissegundos: %.6f ms\n", mediaMS)
	fmt.Printf("Tempo de Resposta Médio Total em Segundos: %.6f s\n", mediaS)
}

func calcularMedia(tempos []float64) float64 {
	if len(tempos) == 0 {
		return 0.0
	}

	soma := 0.0
	for _, tempo := range tempos {
		soma += tempo
	}

	return soma / float64(len(tempos))
}

