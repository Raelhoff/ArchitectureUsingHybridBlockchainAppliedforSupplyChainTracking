const fs = require('fs');
const axios = require('axios');
const csv = require('csv-parser');

const inputFile = './data.csv';
const apiUrl = 'http://192.168.0.192:80/alert'; // URL correta
const requestsPerSecond = 1; // 1 requisição por segundo
const totalRequests = 1000; // Total de requisições a serem enviadas
const outputFile = './log.txt'; // Nome do arquivo de log

async function readCsvData(filePath) {
  const jsonDataArray = [];
  return new Promise((resolve, reject) => {
    fs.createReadStream(filePath)
      .pipe(csv())
      .on('data', (row) => {
        jsonDataArray.push(row);
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

    // Criação do arquivo de log
    const logStream = fs.createWriteStream(outputFile, { flags: 'w' });

    let totalTime = 0;
    let minTime = Infinity;
    let maxTime = -Infinity;

    for (let i = 0; i < totalRequests; i++) {
      if (i > 0) {
        await new Promise((resolve) => setTimeout(resolve, 1000 / requestsPerSecond));
      }

      const jsonData = jsonDataArray[i % jsonDataArray.length];
      const currentDatetime = new Date().toISOString();

      // Atualiza a data dentro do objeto de dados
      jsonData.date = currentDatetime;

      const requestStartTime = new Date().getTime();
      const response = await axios.post(apiUrl, jsonData);
      const requestEndTime = new Date().getTime();

      const requestInfo = {
        requestNumber: i + 1,
        requestTime: currentDatetime,
        responseTime: requestEndTime - requestStartTime,
      };

      console.log(`Request ${i + 1}: Status ${response.status}, Response Time ${requestInfo.responseTime} ms`);

      // Escreve os detalhes da requisição no arquivo de log
      logStream.write(JSON.stringify(requestInfo) + '\n');

      totalTime += requestInfo.responseTime;
      if (requestInfo.responseTime < minTime) {
        minTime = requestInfo.responseTime;
      }
      if (requestInfo.responseTime > maxTime) {
        maxTime = requestInfo.responseTime;
      }
    }

    const averageTime = totalTime / totalRequests;

    console.log('All requests completed.');
    console.log('Total Requests:', totalRequests);
    console.log('Average Response Time:', averageTime.toFixed(2), 'ms');
    console.log('Min Response Time:', minTime, 'ms');
    console.log('Max Response Time:', maxTime, 'ms');

    logStream.end();
  } catch (error) {
    console.error('Error:', error.message);
  }
}

sendRequests();
