import csv
import random
from datetime import datetime, timedelta
import hashlib
import time

# Defina o número de dados que você deseja gerar
num_data_points = 51

# Cria um arquivo CSV para escrever os dados
with open('data.csv', 'w', newline='') as csvfile:
    fieldnames = ["id", "idedge", "idnodo", "hash", "type", "date", "temperatura", "umidade", "rele", "description"]
    writer = csv.DictWriter(csvfile, fieldnames=fieldnames)

    writer.writeheader()

    for id_value in range(0, 51):
        current_time = datetime.now().strftime("%d%H%M%S")  # Formata a hora atual
        numeric_time = ''.join(filter(str.isdigit, current_time))  # Remove caracteres não numéricos

        current_time2 = datetime.now()
        formatted_time = current_time2.strftime("%Y-%m-%d %H:%M")  # Formata a hora atual

        data = {
            "id": f"{numeric_time}",  # Adiciona a hora atual numérica ao id
            "idedge": random.randint(1, 10098),
            "idnodo": random.randint(110, 1006),
            "hash": hashlib.sha512(str(id_value).encode()).hexdigest(),
            "type": random.randint(1, 10),
            "date": formatted_time,  # Formato YYYY-MM-DD HH:MM
            "temperatura": str(random.choice(["20", "21"])),
            "umidade": str(random.choice(["79", "80"])),
            "rele": random.choice(["on", "off"]),
            "description": "Alert data"
        }
        writer.writerow(data)
        
        time.sleep(1)  # Atraso de 1 segundo

print("Dados gerados e salvos em data.csv")
