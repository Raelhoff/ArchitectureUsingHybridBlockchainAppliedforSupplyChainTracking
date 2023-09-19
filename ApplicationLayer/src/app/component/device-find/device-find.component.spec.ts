import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DeviceFindComponent } from './device-find.component';

describe('DeviceFindComponent', () => {
  let component: DeviceFindComponent;
  let fixture: ComponentFixture<DeviceFindComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [DeviceFindComponent]
    });
    fixture = TestBed.createComponent(DeviceFindComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
