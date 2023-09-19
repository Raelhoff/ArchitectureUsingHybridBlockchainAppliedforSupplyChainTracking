import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { map, retry, catchError } from 'rxjs/operators';
import { Alert } from '../models/alert';


interface ApiAlertResponse {
  status: number;
  message: string;
  data: {
    data: Alert[];
  };
}

@Injectable({
  providedIn: 'root'
})

export class AlertService {
  private ipAddress = '192.168.0.192:80';
  private apiUrl = `http://${this.ipAddress}/alert`; // Replace with your API URL

  url = '192.168.0.192:80/alert'; // api rest fake
  //  url = 'http://localhost:8080/REST/resources/cars';

  // injetando o HttpClient
  constructor(private httpClient: HttpClient) { }

  // Headers
  httpOptions = {
    headers: new HttpHeaders({ 'Content-Type': 'application/json' })
  }

  
  // Obtem todos os Alertas
  getAlerts(): Observable<Alert[]> {
    return this.httpClient.get<ApiAlertResponse>(this.apiUrl)
      .pipe(
        retry(2),
        catchError(this.handleError),
        map(response => response.data.data) // Extract the 'data' array from the response
      );
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
