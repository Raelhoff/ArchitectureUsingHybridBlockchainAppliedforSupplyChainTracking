package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type Confirmation struct {
	OK            bool
	TempoResposta time.Duration
}

func main() {
	filePath := "/home/workspace/bitburket/simulator/logs/debug_timing_CreateDevice_9876.txt"

	// Read and process the log file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close()

	var durations []time.Duration

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		confirmation, err := parseConfirmation(line)
		if err != nil {
			fmt.Println("Erro ao analisar a linha:", err)
			continue
		}

		durations = append(durations, confirmation.TempoResposta)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Erro durante a leitura do arquivo:", err)
		return
	}

	if len(durations) == 0 {
		fmt.Println("Nenhum dado de tempo de resposta encontrado.")
		return
	}

	max := maxDuration(durations)
	min := minDuration(durations)
	avg := avgDuration(durations)

	fmt.Printf("Tempo de Resposta Máximo: %s\n", max)
	fmt.Printf("Tempo de Resposta Mínimo: %s\n", min)
	fmt.Printf("Tempo de Resposta Médio: %s\n", avg)

	// Create a backup file
	backupFilePath := "/home/workspace/bitburket/simulator/logs/debug_timing_CreateDevice_9876_bkp.txt"
	backupFile, err := os.OpenFile(backupFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo de backup:", err)
		return
	}
	defer backupFile.Close()

	// Write the content of the original file to the backup file
	originalFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer originalFile.Close()

	scannerOriginal := bufio.NewScanner(originalFile)
	for scannerOriginal.Scan() {
		line := scannerOriginal.Text()
		_, writeErr := backupFile.WriteString(line + "\n")
		if writeErr != nil {
			fmt.Println("Erro ao escrever no arquivo de backup:", writeErr)
			return
		}
	}

	if err := scannerOriginal.Err(); err != nil {
		fmt.Println("Erro durante a leitura do arquivo original:", err)
		return
	}

	// Write max, min, and avg response times to the backup file
	_, writeErr := backupFile.WriteString(fmt.Sprintf("Tempo de Resposta Máximo: %s\n", max))
	if writeErr != nil {
		fmt.Println("Erro ao escrever no arquivo de backup:", writeErr)
		return
	}
	_, writeErr = backupFile.WriteString(fmt.Sprintf("Tempo de Resposta Mínimo: %s\n", min))
	if writeErr != nil {
		fmt.Println("Erro ao escrever no arquivo de backup:", writeErr)
		return
	}
	_, writeErr = backupFile.WriteString(fmt.Sprintf("Tempo de Resposta Médio: %s\n", avg))
	if writeErr != nil {
		fmt.Println("Erro ao escrever no arquivo de backup:", writeErr)
		return
	}

	_, writeErr = backupFile.WriteString(fmt.Sprintf("------------------------------------------------------------------------------\n"))
	if writeErr != nil {
		fmt.Println("Erro ao escrever no arquivo de backup:", writeErr)
		return
	}

	// Close the backup file before removing the original file
	backupFile.Close()

	// Delete the original log file
	err = os.Remove(filePath)
	if err != nil {
		fmt.Println("Erro ao deletar o arquivo original:", err)
		return
	}

	fmt.Printf("Conteúdo do arquivo %s copiado para %s.\n", filePath, backupFilePath)
}

func parseConfirmation(line string) (Confirmation, error) {
	parts := strings.Split(line, ", ")

	if len(parts) != 6 {
		return Confirmation{}, fmt.Errorf("linha mal formatada")
	}

	okPart := strings.TrimSpace(parts[0])
	if okPart != "Confirmation: true" {
		return Confirmation{}, fmt.Errorf("confirmação não OK")
	}

	var tempoResposta time.Duration
	for _, part := range parts {
		if strings.HasPrefix(part, "Tempo de resposta:") {
			tempoStr := strings.TrimSpace(strings.TrimPrefix(part, "Tempo de resposta:"))
			duration, err := time.ParseDuration(tempoStr)
			if err != nil {
				return Confirmation{}, fmt.Errorf("erro ao analisar tempo de resposta: %v", err)
			}
			tempoResposta = duration
			break
		}
	}

	if tempoResposta == 0 {
		return Confirmation{}, fmt.Errorf("nenhum dado de tempo de resposta encontrado")
	}

	return Confirmation{OK: true, TempoResposta: tempoResposta}, nil
}

func maxDuration(durations []time.Duration) time.Duration {
	max := durations[0]
	for _, duration := range durations {
		if duration > max {
			max = duration
		}
	}
	return max
}

func minDuration(durations []time.Duration) time.Duration {
	min := durations[0]
	for _, duration := range durations {
		if duration < min {
			min = duration
		}
	}
	return min
}

func avgDuration(durations []time.Duration) time.Duration {
	total := int64(0)
	for _, duration := range durations {
		total += duration.Nanoseconds()
	}
	average := total / int64(len(durations))
	return time.Duration(average)
}
