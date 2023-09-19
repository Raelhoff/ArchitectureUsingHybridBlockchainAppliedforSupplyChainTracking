const fs = require('fs');
const axios = require('axios');
const csv = require('csv-parser');
const crypto = require('crypto');

const inputFile = './data.csv';
const apiUrl = 'http://192.168.0.192:80/alertHashBack'; // URL atualizada
const requestsPerSecond = 100; // 1 requisição por segundo
const totalRequests = 1500; // Total de requisições a serem enviadas
const outputFile = './log.txt'; // Nome do arquivo de log

let currentRequestId = 1005+1500; // Variável para manter o ID único

// Função para gerar um hash único baseado em data, minuto, segundo e milissegundo
function generateUniqueHash() {
  const currentDate = new Date();
  const dataToHash = `${currentDate.getFullYear()}${currentDate.getMonth()}${currentDate.getDate()}${currentDate.getHours()}${currentDate.getMinutes()}${currentDate.getSeconds()}${currentDate.getMilliseconds()}`;
  
  const hash = crypto.createHash('sha512'); // Crie um objeto de hash SHA-512
  hash.update(dataToHash); // Atualize o hash com os dados formatados
  return hash.digest('hex'); // Retorne o hash em formato hexadecimal
}

// Função para gerar um ID único crescente
function generateUniqueId() {
  return ++currentRequestId;
}

async function readCsvData(filePath) {
  const jsonDataArray = [];

  return new Promise((resolve, reject) => {
    fs.createReadStream(filePath)
      .pipe(csv())
      .on('data', (row) => {
        const currentDateTime = new Date().toISOString();

        const jsonData = {
          id: generateUniqueHash(),
          idedge: parseInt(row.idedge, 10),
          idnodo: parseInt(row.idnodo, 10),
          hash: row.hash.toString(),
          type: parseInt(row.type, 10),
          date: currentDateTime.toString(),
          temperatura: row.temperatura.toString(),
          umidade: row.umidade.toString(),
          rele: row.rele.toString(),
          description: row.description.toString(),
        };

        jsonDataArray.push(jsonData);
      })
      .on('end', () => {
        resolve(jsonDataArray);
      })
      .on('error', (error) => {
        reject(error);
      });
  });
}

async function sendRequests() {
  try {
    const jsonDataArray = await readCsvData(inputFile);

    // Criar o arquivo de log
    const logStream = fs.createWriteStream(outputFile, { flags: 'w' });

    let totalResponseTime = 0;
    let minResponseTime = Infinity;
    let maxResponseTime = -Infinity;

    for (let i = 0; i < totalRequests; i++) {
      if (i > 0) {
        await new Promise((resolve) => setTimeout(resolve, 1000 / requestsPerSecond));
      }

      const requestStartTime = new Date().getTime();
      const jsonData = jsonDataArray[i % jsonDataArray.length];

      const response = await axios.post(apiUrl, jsonData);
      const requestEndTime = new Date().getTime();

      const responseTime = requestEndTime - requestStartTime;
      totalResponseTime += responseTime;

      if (responseTime < minResponseTime) {
        minResponseTime = responseTime;
      }

      if (responseTime > maxResponseTime) {
        maxResponseTime = responseTime;
      }

      const requestInfo = {
        requestNumber: i + 1,
        requestTime: new Date(requestStartTime).toISOString(),
        responseTime: responseTime,
      };

      console.log(`Request ${i + 1}: Status ${response.status}, Response Time ${requestInfo.responseTime} ms`);

      // Escrever os detalhes da solicitação no arquivo de log
      logStream.write(JSON.stringify(requestInfo) + '\n');
    }

    const averageResponseTime = totalResponseTime / totalRequests;

    // Escrever os tempos de resposta mínimo, máximo e médio no arquivo de log
    logStream.write(`Minimum Response Time: ${minResponseTime} ms\n`);
    logStream.write(`Maximum Response Time: ${maxResponseTime} ms\n`);
    logStream.write(`Average Response Time: ${averageResponseTime.toFixed(1)} ms\n`);

    console.log(`Minimum Response Time: ${minResponseTime} ms`);
    console.log(`Maximum Response Time: ${maxResponseTime} ms`);
    console.log(`Average Response Time: ${averageResponseTime.toFixed(1)} ms`);
    console.log(`Process completed at: ${new Date().toLocaleTimeString()}`);

    logStream.end();
  } catch (error) {
    console.error('Error:', error.message);
  }
}


// Inicia a função sendRequests() para executar a cada 2 horas e 30 minutos (10 vezes)
//const interval = setInterval(sendRequests, 2 * 60 * 60 * 1000 + 30 * 60 * 1000);

// Inicia a função sendRequests() para executar a cada 2 horas (10 vezes)
sendRequests()

//const interval = setInterval(sendRequests, 2 * 60 * 60 * 1000);

// Inicia a função sendRequests() para executar a cada 1 hora (10 vezes)
//const interval = setInterval(sendRequests, 60 * 60 * 1000);

// Inicia a função sendRequests() para executar a cada 45 minutos (10 vezes)//
//const interval = setInterval(sendRequests, 45 * 60 * 1000);

// Inicia a função sendRequests() para executar a cada 1 hora e 15 minutos (10 vezes)
const interval = setInterval(sendRequests, 1 * 60 * 60 * 1000 + 15 * 60 * 1000);
