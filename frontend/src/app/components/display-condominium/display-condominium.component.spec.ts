import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DisplayCondominiumComponent } from './display-condominium.component';

describe('DisplayCondominiumComponent', () => {
  let component: DisplayCondominiumComponent;
  let fixture: ComponentFixture<DisplayCondominiumComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DisplayCondominiumComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(DisplayCondominiumComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
