import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CoproprietaireComponent } from './coproprietaire.component';

describe('CoproprietaireComponent', () => {
  let component: CoproprietaireComponent;
  let fixture: ComponentFixture<CoproprietaireComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CoproprietaireComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CoproprietaireComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
