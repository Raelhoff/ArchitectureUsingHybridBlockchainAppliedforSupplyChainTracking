import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { CadastroProdutosComponent } from './cadastro-produtos/cadastro-produtos.component';
import { ListProdutosComponent } from './list-produtos/list-produtos.component';

const routes: Routes = [
  {
    path: '',
    data: {
      title: 'Produtos'
    },
    children: [
      {
        path: '',
        pathMatch: 'full',
        redirectTo: 'produtos'
      },
      {
        path: 'cadastro',
        component: CadastroProdutosComponent,
        data: {
          title: 'Produtos Cadastro'
        }
      },
      {
        path: 'lista',
        component: ListProdutosComponent,
        data: {
          title: 'Produtos Lista'
        }
      },
    ]
  }
];


@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class ProdutosRoutingModule {
}
