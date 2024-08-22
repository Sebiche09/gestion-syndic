import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FtpblueprintComponent } from './ftpblueprint.component';

describe('FtpblueprintComponent', () => {
  let component: FtpblueprintComponent;
  let fixture: ComponentFixture<FtpblueprintComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [FtpblueprintComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(FtpblueprintComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
