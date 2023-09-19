import { Component } from '@angular/core';
import { DeviceService } from 'src/app/services/device.service';
import { Device } from 'src/app/models/device';
import { NgForm } from '@angular/forms';

@Component({
  selector: 'app-list-dispositivos',
  templateUrl: './list-dispositivos.component.html',
  styleUrls: ['./list-dispositivos.component.scss']
})
export class ListDispositivosComponent {
  device: Device = new Device();
  devices: Device[] = []; // Initializing the property with an empty array

  constructor(private deviceService: DeviceService) {}

  ngOnInit() {
    this.getDevices();
  }

      // defini se um carro será criado ou atualizado
      saveDevice(form: NgForm) {
        if (this.device.ID !== undefined) {
          this.deviceService.updateCar(this.device).subscribe(() => {
          //  this.cleanForm(form);
          });
        } else {
          this.deviceService.saveCar(this.device).subscribe(() => {
           // this.cleanForm(form);
          });
        }
      }

      // Chama o serviço para obtém todos os devices
      getDevices() {
      this.deviceService.getDevices().subscribe((devices: Device[]) => {
        this.devices = devices;
      });
    }

    // deleta um device
    deleteDevice(device: Device) {
      this.deviceService.deleteCar(device).subscribe(() => {
        this.getDevices();
      });
    }

    // copia o carro para ser editado.
    editDevice(device: Device) {
      this.device = { ...device };
    }

  }

