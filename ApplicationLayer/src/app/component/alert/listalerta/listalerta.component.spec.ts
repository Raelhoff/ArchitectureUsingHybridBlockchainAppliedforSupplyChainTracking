import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ListalertaComponent } from './listalerta.component';

describe('ListalertaComponent', () => {
  let component: ListalertaComponent;
  let fixture: ComponentFixture<ListalertaComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ListalertaComponent]
    });
    fixture = TestBed.createComponent(ListalertaComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
