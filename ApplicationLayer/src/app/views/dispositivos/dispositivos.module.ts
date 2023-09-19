import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ListDispositivosComponent } from './list-dispositivos/list-dispositivos.component';
import { ReactiveFormsModule } from '@angular/forms';

import { DocsComponentsModule } from '@docs-components/docs-components.module';


import { DispositivosRoutingModule } from './dispositivos-routing.module';


import {
  ButtonGroupModule,
  ButtonModule,
  CardModule,
  CollapseModule,
  DropdownModule,
  FormModule,
  GridModule,
  NavbarModule,
  NavModule,
  SharedModule,
  UtilitiesModule
} from '@coreui/angular';


import { IconModule } from '@coreui/icons-angular';
import { CadastroDispositivoComponent } from './cadastro-dispositivo/cadastro-dispositivo.component';

@NgModule({
  declarations: [
    ListDispositivosComponent,
    CadastroDispositivoComponent
  ],
  imports: [
    CommonModule,
    DispositivosRoutingModule,
    ButtonModule,
    ButtonGroupModule,
    GridModule,
    IconModule,
    CardModule,
    UtilitiesModule,
    DropdownModule,
    SharedModule,
    FormModule,
    ReactiveFormsModule,
    DocsComponentsModule,
    NavbarModule,
    CollapseModule,
    NavModule,
    NavbarModule
  ]
})


export class DispositivosModule { }
