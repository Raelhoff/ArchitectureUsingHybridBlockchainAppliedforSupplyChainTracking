import csv
import random
from datetime import datetime
import hashlib

# Defina o número de dados que você deseja gerar
num_data_points = 51

# Cria um arquivo CSV para escrever os dados
with open('data.csv', 'w', newline='') as csvfile:
    fieldnames = ["id", "idedge", "idnodo", "hash", "type", "date", "temperatura", "umidade", "rele", "description"]
    writer = csv.DictWriter(csvfile, fieldnames=fieldnames)

    writer.writeheader()

    for id_value in range(0, 3000):
        data = {
            "id": id_value,
            "idedge": random.randint(1, 10098),
            "idnodo": random.randint(110, 1006),
            "hash": hashlib.sha512(str(id_value).encode()).hexdigest(),  # Calcula hash SHA-512
            "type": random.randint(1, 10),
            "date": datetime.now().isoformat(),
            "temperatura": str(random.choice(["20", "21"])),
            "umidade": str(random.choice(["79", "80"])),
            "rele": random.choice(["on", "off"]),
            "description": "Alert data"
        }
        writer.writerow(data)

print("Dados gerados e salvos em data.csv")
