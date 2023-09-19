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

	// Inicializar uma lista para armazenar os tempos
	tempos := make([]float64, 0)

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
			if strings.HasSuffix(tempoStr, "ms") || strings.HasSuffix(tempoStr, "s") {
				tempo, err := parseTempo(tempoStr)
				if err == nil {
					tempos = append(tempos, tempo)
				}
			}
		}
	}

	// Calcular a média
	media := calcularMedia(tempos)

	fmt.Printf("Tempo de Resposta Médio Total: %.6f ms\n", media)
}

func parseTempo(tempoStr string) (float64, error) {
	tempoStr = strings.TrimSpace(tempoStr)
	if strings.HasSuffix(tempoStr, "ms") {
		return strconv.ParseFloat(strings.TrimSuffix(tempoStr, "ms"), 64)
	} else if strings.HasSuffix(tempoStr, "s") {
		tempo, err := strconv.ParseFloat(strings.TrimSuffix(tempoStr, "s"), 64)
		if err == nil {
			return tempo * 1000, nil // Converter segundos para milissegundos
		}
	}
	return 0, fmt.Errorf("formato de tempo inválido")
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
