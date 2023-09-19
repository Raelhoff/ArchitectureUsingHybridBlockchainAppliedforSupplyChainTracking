import re

# Abrir o arquivo para leitura
with open("/home/specto/requisicoesMB/1EDGE/40/debug_timing_CreateDevice_1234_40_trm.txt", "r") as arquivo:
    linhas = arquivo.readlines()

# Extrair os tempos de resposta médios dos dados lidos
tempos_de_resposta = []
pattern = r"Tempo de Resposta Médio: (\d+\.\d+)(ms|s)"
for linha in linhas:
    match = re.search(pattern, linha)
    if match:
        valor = float(match.group(1))
        unidade = match.group(2)
        
        if unidade == "s":
            valor *= 1000  # Convertendo para milissegundos
        
        tempos_de_resposta.append(valor)

# Calcular a média dos tempos de resposta, apenas se houver tempos
if tempos_de_resposta:
    media = sum(tempos_de_resposta) / len(tempos_de_resposta)
    print("Média do Tempo de Resposta Médio:", media, "ms")
else:
    print("Nenhum tempo de resposta encontrado no arquivo.")

