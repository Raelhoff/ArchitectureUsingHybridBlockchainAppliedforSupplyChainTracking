const fs = require('fs');
const axios = require('axios');
const crypto = require('crypto');
const csv = require('csv-parser');

const inputFile = './data.csv';
const apiUrl = 'http://192.168.0.192:80/create-assets'; // Updated URL
const requestsPerSecond = 20; // 1 requisição por segundo
const totalRequests = 300; // Total de requisições a serem enviadas
const outputFile = './log.txt'; // Nome do arquivo de log

async function readCsvData(filePath) {
  const jsonDataArray = [];
  return new Promise((resolve, reject) => {
    fs.createReadStream(filePath)
      .pipe(csv())
      .on('data', (row) => {
        const currentDateTime = new Date().toISOString();
        const randomData = Math.random().toString();
        const dataToHash = currentDateTime + randomData;

        const hash = crypto.createHash('sha512').update(dataToHash).digest('hex');

        const jsonData = {
          id: hash,
          idedge: '7390',
          idnodo: '827',
          hash: hash,
          type: '10',
          date: currentDateTime,
          temperatura: '',
          umidade: '',
          rele: 'off',
          description: 'Alert data',
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

    // Criação do arquivo de log
    const logStream = fs.createWriteStream(outputFile, { flags: 'w' });

    let totalResponseTime = 0;
    let minResponseTime = Infinity;
    let maxResponseTime = -Infinity;

    for (let i = 0; i < totalRequests; i++) {
      if (i > 0) {
        await new Promise((resolve) => setTimeout(resolve, 1000 / requestsPerSecond));
      }

      const requestStartTime = new Date().getTime();
      const response = await axios.post(apiUrl, jsonDataArray[i % jsonDataArray.length]);
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

      // Escreve os detalhes da requisição no arquivo de log
      logStream.write(JSON.stringify(requestInfo) + '\n');
    }

    const averageResponseTime = totalResponseTime / totalRequests;

    // Write the minimum, maximum, and average response times to the log file
    logStream.write(`Minimum Response Time: ${minResponseTime} ms\n`);
    logStream.write(`Maximum Response Time: ${maxResponseTime} ms\n`);
    logStream.write(`Average Response Time: ${averageResponseTime.toFixed(1)} ms\n`);

    console.log(`Minimum Response Time: ${minResponseTime} ms`);
    console.log(`Maximum Response Time: ${maxResponseTime} ms`);
    console.log(`Average Response Time: ${averageResponseTime.toFixed(1)} ms`);

    logStream.end();
  } catch (error) {
    console.error('Error:', error.message);
  }
}

sendRequests();
