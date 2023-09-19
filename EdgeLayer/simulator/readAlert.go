package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type Confirmation struct {
	TempoEnvio     time.Time
	TempoRecebido  time.Time
	TempoResposta  time.Duration
}

func main() {
	filePath := "/home/workspace/bitburket/client/logs/debug_timing_sendAlertDeviceNotResponde_1234.txt" // Update with the correct file path

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close()

	var confirmations []Confirmation

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		confirmation, err := parseConfirmation(line)
		if err != nil {
			fmt.Println("Erro ao analisar a linha:", err)
			continue
		}

		confirmations = append(confirmations, confirmation)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Erro durante a leitura do arquivo:", err)
		return
	}

	if len(confirmations) == 0 {
		fmt.Println("Nenhum dado de confirmação encontrado.")
		return
	}

	max := maxDuration(confirmations)
	min := minDuration(confirmations)
	avg := avgDuration(confirmations)

	fmt.Printf("Tempo de Resposta Máximo: %s\n", max)
	fmt.Printf("Tempo de Resposta Mínimo: %s\n", min)
	fmt.Printf("Tempo de Resposta Médio: %s\n", avg)
}

func parseConfirmation(line string) (Confirmation, error) {
	parts := strings.Split(line, ", ")

	if len(parts) != 3 {
		return Confirmation{}, fmt.Errorf("linha mal formatada")
	}

	enviadoStr := strings.TrimPrefix(parts[0], "Enviado em: ")
	enviado, err := time.Parse("02/01/2006 15:04:05", enviadoStr)
	if err != nil {
		return Confirmation{}, fmt.Errorf("erro ao analisar tempo de envio: %v", err)
	}

	recebidoStr := strings.TrimPrefix(parts[1], "Recebido em: ")
	recebido, err := time.Parse("02/01/2006 15:04:05", recebidoStr)
	if err != nil {
		return Confirmation{}, fmt.Errorf("erro ao analisar tempo de recebimento: %v", err)
	}

	tempoStr := strings.TrimPrefix(parts[2], "Tempo de resposta: ")
	tempo, err := time.ParseDuration(tempoStr)
	if err != nil {
		return Confirmation{}, fmt.Errorf("erro ao analisar tempo de resposta: %v", err)
	}

	return Confirmation{
		TempoEnvio:     enviado,
		TempoRecebido:  recebido,
		TempoResposta:  tempo,
	}, nil
}

func maxDuration(confirmations []Confirmation) time.Duration {
	max := confirmations[0].TempoResposta
	for _, confirmation := range confirmations {
		if confirmation.TempoResposta > max {
			max = confirmation.TempoResposta
		}
	}
	return max
}

func minDuration(confirmations []Confirmation) time.Duration {
	min := confirmations[0].TempoResposta
	for _, confirmation := range confirmations {
		if confirmation.TempoResposta < min {
			min = confirmation.TempoResposta
		}
	}
	return min
}

func avgDuration(confirmations []Confirmation) time.Duration {
	total := int64(0)
	for _, confirmation := range confirmations {
		total += confirmation.TempoResposta.Nanoseconds()
	}
	average := total / int64(len(confirmations))
	return time.Duration(average)
}

