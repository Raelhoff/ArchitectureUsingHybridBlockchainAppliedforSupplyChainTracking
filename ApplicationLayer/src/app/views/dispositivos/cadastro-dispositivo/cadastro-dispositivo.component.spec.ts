import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CadastroDispositivoComponent } from './cadastro-dispositivo.component';

describe('CadastroDispositivoComponent', () => {
  let component: CadastroDispositivoComponent;
  let fixture: ComponentFixture<CadastroDispositivoComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [CadastroDispositivoComponent]
    });
    fixture = TestBed.createComponent(CadastroDispositivoComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
