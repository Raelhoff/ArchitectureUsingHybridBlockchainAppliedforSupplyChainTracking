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
        const jsonData = {
          id: 1,
          idedge: 2,
          idnodo: 3,
          hash: 'abc123',
          type: 4,
          date: new Date().toISOString(), // Use current date and time
          temperatura: '25',
          umidade: '60',
          rele: 'on',
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
  
      // Create the log file
      const logStream = fs.createWriteStream(outputFile, { flags: 'w' });
  
      let totalResponseTime = 0;
      let minResponseTime = Infinity;
      let maxResponseTime = -Infinity;
  
      for (let i = 0; i < totalRequests; i++) {
        if (i > 0) {
          await new Promise((resolve) => setTimeout(resolve, 1000 / requestsPerSecond));
        }
  
        const jsonData = jsonDataArray[i % jsonDataArray.length];
  
        const requestStartTime = new Date().getTime();
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
  
        // Write the details of the request to the log file
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
