import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ListProdutosComponent } from './list-produtos.component';

describe('ListProdutosComponent', () => {
  let component: ListProdutosComponent;
  let fixture: ComponentFixture<ListProdutosComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ListProdutosComponent]
    });
    fixture = TestBed.createComponent(ListProdutosComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
