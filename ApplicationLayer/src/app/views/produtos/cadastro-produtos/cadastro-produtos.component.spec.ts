import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CadastroProdutosComponent } from './cadastro-produtos.component';

describe('CadastroProdutosComponent', () => {
  let component: CadastroProdutosComponent;
  let fixture: ComponentFixture<CadastroProdutosComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [CadastroProdutosComponent]
    });
    fixture = TestBed.createComponent(CadastroProdutosComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
