import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { HomeComponent } from './component/home/home.component';
import { DeviceFindComponent } from './component/device-find/device-find.component';
import { DeviceCreateComponent } from './component/device-create/device-create.component';
import { DeviceListComponent } from './component/device-list/device-list.component';

const APP_ROUTES: Routes = [
  {path: '', component: HomeComponent},
  {path: 'createdevice', component: DeviceCreateComponent},
  {path: 'listadevice', component: DeviceListComponent},
  {path: 'deviceupdate/:id', component: DeviceFindComponent},
];

@NgModule({
    imports: [RouterModule.forRoot(APP_ROUTES)],
    exports: [RouterModule]
  })
  export class AppRoutingModule { }
  
  export const routing = AppRoutingModule;
