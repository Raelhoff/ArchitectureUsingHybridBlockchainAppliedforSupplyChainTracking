import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { map, retry, catchError } from 'rxjs/operators';
import { Device } from '../models/device';


interface ApiResponse {
  status: number;
  message: string;
  data: {
    data: Device[];
  };
}

@Injectable({
  providedIn: 'root'
})

export class DeviceService {

  url = 'http://192.168.0.192/device'; // api rest fake
  //  url = 'http://localhost:8080/REST/resources/cars';

  // injetando o HttpClient
  constructor(private httpClient: HttpClient) { }

  // Headers
  httpOptions = {
    headers: new HttpHeaders({ 'Content-Type': 'application/json' })
  }

  // Obtem todos os Devices
  getDevices(): Observable<Device[]> {
    return this.httpClient.get<ApiResponse>(this.url)
      .pipe(
        retry(2),
        catchError(this.handleError),
        map(response => response.data.data) // Extract the 'data' array from the response
      );
  }

  // Obtem um Device pelo id
  getDeviceById(id: number): Observable<Device> {
    return this.httpClient.get<Device>(this.url + '/' + id)
      .pipe(
        retry(2),
        catchError(this.handleError)
      )
  }

  // salva um Device
  saveCar(device: Device): Observable<Device> {
    return this.httpClient.post<Device>(this.url, JSON.stringify(device), this.httpOptions)
      .pipe(
        catchError(this.handleError)
      )
  }

  // utualiza um Device
  updateCar(device: Device): Observable<Device> {
    return this.httpClient.put<Device>(this.url + '/' + device.ID, JSON.stringify(device), this.httpOptions)
      .pipe(
        catchError(this.handleError)
      )
  }

  // deleta um device
  deleteCar(device: Device) {
    return this.httpClient.delete<Device>(this.url + '/' + device.ID, this.httpOptions)
      .pipe(
        catchError(this.handleError)
      )
  }

  // Manipulação de erros
  handleError(error: HttpErrorResponse) {
    let errorMessage = '';
    if (error.error instanceof ErrorEvent) {
      // Erro ocorreu no lado do client
      errorMessage = error.error.message;
    } else {
      // Erro ocorreu no lado do servidor
      errorMessage = `Código do erro: ${error.status}, ` + `menssagem: ${error.message}`;
    }
    console.log(errorMessage);
    return throwError(errorMessage);
  };

}
