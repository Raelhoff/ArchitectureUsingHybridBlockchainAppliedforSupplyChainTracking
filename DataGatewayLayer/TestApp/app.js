const fs = require('fs');
const axios = require('axios');
const csv = require('csv-parser');

const inputFile = './data.csv';
const apiUrl = 'http://192.168.0.192:80/create-assets'; // Substitua pela URL correta
const requestsPerSecond = 1; // Número de requisições por segundo
const maxExecutionTime = 60; // Tempo máximo de execução em segundos
const numThreads = 4; // Número de threads para enviar as requisições em paralelo

async function readCsvData(filePath) {
  const jsonDataArray = [];
  return new Promise((resolve, reject) => {
    fs.createReadStream(filePath)
      .pipe(csv())
      .on('data', (row) => {
        const jsonData = {
          id: row.id,
          idedge: row.idedge,
          idnodo: row.idnodo,
          hash: row.hash,
          type: row.type,
          date: row.date,
          temperatura: row.temperatura,
          umidade: row.umidade,
          rele: row.rele,
          description: row.description,
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

async function sendPostRequestsInParallel(dataArray) {
  const startTime = new Date().getTime();
  const requestResults = [];

  for (let i = 0; i < numThreads; i++) {
    sendRequestsInThread(dataArray, startTime, requestResults);
  }

  await Promise.all(requestResults);

  return requestResults;
}

async function sendRequestsInThread(dataArray, startTime, requestResults) {
  const threadRequests = dataArray.length / numThreads;

  for (let i = 0; i < threadRequests; i++) {
    if ((new Date().getTime() - startTime) / 1000 >= maxExecutionTime) {
      console.log('Max execution time reached.');
      return; // Sair do loop se o tempo máximo for atingido
    }

    const jsonData = dataArray[i];
    const requestInfo = {
      requestData: jsonData,
      success: false,
      responseTime: null,
    };

    try {
      const requestStartTime = new Date().getTime();
      const response = await axios.post(apiUrl, jsonData);
      const requestEndTime = new Date().getTime();

      if (response.status === 200) {
        requestInfo.success = true;
        requestInfo.responseTime = requestEndTime - requestStartTime;
        console.log('POST response:', response.data);
      } else {
        console.error('Error posting data:', response.status, response.statusText);
      }
    } catch (error) {
      if (error.response && error.response.status === 500) {
        console.error('Error posting data:', error.message);
        requestInfo.errorMessage = error.message;
      } else {
        console.error('Unexpected error:', error.message);
      }
    }

    requestResults.push(Promise.resolve(requestInfo));
  }
}

async function main() {
  try {
    const jsonDataArray = await readCsvData(inputFile);
    const requestResults = await sendPostRequestsInParallel(jsonDataArray);

    const successfulRequests = requestResults.filter(result => result.success);
    const failedRequests = requestResults.filter(result => !result.success);

    const totalRequests = requestResults.length;
    const successfulRequestCount = successfulRequests.length;
    const failedRequestCount = failedRequests.length;

    const successfulResponseTimeSum = successfulRequests.reduce((sum, result) => sum + result.responseTime, 0);
    const averageSuccessfulResponseTime = successfulResponseTimeSum / successfulRequestCount;
    const minSuccessfulResponseTime = Math.min(...successfulRequests.map(result => result.responseTime));
    const maxSuccessfulResponseTime = Math.max(...successfulRequests.map(result => result.responseTime));

    console.log('All requests completed.');
    console.log('Total Requests:', totalRequests);
    console.log('Successful Requests:', successfulRequestCount);
    console.log('Failed Requests:', failedRequestCount);
    console.log('Average Successful Response Time:', averageSuccessfulResponseTime.toFixed(2), 'ms');
    console.log('Min Successful Response Time:', minSuccessfulResponseTime, 'ms');
    console.log('Max Successful Response Time:', maxSuccessfulResponseTime, 'ms');

    failedRequests.forEach((failedRequest, index) => {
      console.log(`Failed Request ${index + 1}:`);
      console.log('  Data:', failedRequest.requestData);
      console.log('  Error:', failedRequest.errorMessage);
    });
  } catch (error) {
    console.error('Error:', error.message);
  }
}

main();
