# Abre o arquivo para leitura
with open('/home/specto/requisicoesMB/1EDGE/40/debug_timing_CreateDevice_1234_40.txt', 'r') as file:
    lines = file.readlines()

# Inicializa uma lista para armazenar os tempos de resposta em milissegundos
tempos_resposta_ms = []

# Itera pelas linhas do arquivo
for line in lines:
    # Verifica se a linha contém "Tempo de resposta:"
    if "Tempo de resposta:" in line:
        # Encontra o índice onde o tempo de resposta está
        index_inicio_tempo = line.index("Tempo de resposta:") + len("Tempo de resposta:")
        index_fim_tempo = line.index("ms" if "ms" in line else "s", index_inicio_tempo)
        
        # Extrai o valor do tempo de resposta
        tempo_str = line[index_inicio_tempo:index_fim_tempo]
        
        # Converte o tempo para float e converte para milissegundos se necessário
        tempo_resposta = float(tempo_str)
        if "ms" in line:
            tempo_resposta /= 1000

        print(f"Tempo de resposta: {tempo_resposta:.6f}ms")
        tempos_resposta_ms.append(tempo_resposta)
	
# Calcula a média dos tempos de resposta em milissegundos
media_tempo_resposta_ms = sum(tempos_resposta_ms) / len(tempos_resposta_ms)

# Imprime o resultado
print(f"Tempo médio de resposta: {media_tempo_resposta_ms:.6f}ms")

