import { Component } from '@angular/core';
import { NgForm } from '@angular/forms';
import { AlertService } from 'src/app/services/alert.service';
import { Alert } from 'src/app/models/alert';

@Component({
  selector: 'app-listalerta',
  templateUrl: './listalerta.component.html',
  styleUrls: ['./listalerta.component.scss']
})
export class ListalertaComponent {
  alerta: Alert = new Alert();
  alertas: Alert[] = []; // Initializing the property with an empty array

  constructor(private alertService: AlertService) {}

  ngOnInit() {
    this.getAlerts();
  }

        // Chama o serviço para obtém todos os alertas
        getAlerts() {
          this.alertService.getAlerts().subscribe((alertas: Alert[]) => {
            this.alertas = alertas;
          });
        }

}
