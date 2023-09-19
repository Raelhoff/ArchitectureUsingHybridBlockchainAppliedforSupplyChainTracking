import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { ListDispositivosComponent } from './list-dispositivos/list-dispositivos.component';
import { CadastroDispositivoComponent } from './cadastro-dispositivo/cadastro-dispositivo.component';


const routes: Routes = [
  {
    path: '',
    data: {
      title: 'Dispositivos'
    },
    children: [
      {
        path: '',
        pathMatch: 'full',
        redirectTo: 'dispositivos'
      },
      {
        path: 'cadastro',
        component: CadastroDispositivoComponent,
        data: {
          title: 'Lista de Dispositivos'
        }
      },
      {
        path: 'lista',
        component: ListDispositivosComponent,
        data: {
          title: 'Lista de Dispositivos'
        }
      },
    ]
  }
];


@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class DispositivosRoutingModule {
}
