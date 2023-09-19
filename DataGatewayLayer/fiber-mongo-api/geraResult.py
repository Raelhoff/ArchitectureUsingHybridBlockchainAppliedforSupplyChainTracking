import os
from datetime import datetime

# Lê o arquivo transaction_log.txt
with open("transaction_log.txt", "r") as arquivo_entrada:
    linhas = arquivo_entrada.readlines()

# Pega as datas da primeira e última linha
primeira_linha = linhas[0]
ultima_linha = linhas[-1]

# Extrai as datas das linhas usando expressões regulares
import re
padrao_data = r"Data: (\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})"
data_inicio_str = re.search(padrao_data, primeira_linha).group(1)
data_fim_str = re.search(padrao_data, ultima_linha).group(1)

# Converte as datas de string para objetos datetime
data_inicio = datetime.strptime(data_inicio_str, "%Y-%m-%d %H:%M:%S")
data_fim = datetime.strptime(data_fim_str, "%Y-%m-%d %H:%M:%S")

# Calcula a diferença de tempo em segundos
diferenca_tempo = (data_fim - data_inicio).total_seconds()

# Abre o arquivo de saída em modo de adição
with open("resultado.txt", "a") as arquivo_saida:
    # Escreve o resultado no arquivo
    arquivo_saida.write(str(diferenca_tempo) + "s\n")

# Remove o arquivo transaction_log.txt
os.remove("transaction_log.txt")
