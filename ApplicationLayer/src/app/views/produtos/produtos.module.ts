import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { CadastroProdutosComponent } from './cadastro-produtos/cadastro-produtos.component';
import { ListProdutosComponent } from './list-produtos/list-produtos.component';

import{ProdutosRoutingModule} from './produtos-routing.module'


@NgModule({
  declarations: [
    CadastroProdutosComponent,
    ListProdutosComponent
  ],
  imports: [
    CommonModule,
    ProdutosRoutingModule
  ]
})
export class ProdutosModule { }
