import { Component, Input, signal } from '@angular/core';
import { FormBuilder, FormGroup, FormArray } from '@angular/forms';
import { TableModule } from 'primeng/table';
import { TagModule } from 'primeng/tag';
import { DialogModule } from 'primeng/dialog';
import { ButtonModule } from 'primeng/button';

@Component({
  selector: 'app-unit',
  standalone: true,
  imports: [TableModule, TagModule, DialogModule, ButtonModule],
  templateUrl: './unit.component.html',
  styleUrls: ['./unit.component.scss']
})
export class UnitComponent {
  @Input() units!: FormArray;

  displayDetailsDialog = signal(false);  
  selectedUnit = signal<FormGroup | null>(null);  

  constructor(private fb: FormBuilder) {}

  viewDetails(index: number): void {
    this.selectedUnit.set(this.units.at(index) as FormGroup);
    this.displayDetailsDialog.set(true);
  }

  closeDetailsDialog(): void {
    this.displayDetailsDialog.set(false);
    this.selectedUnit.set(null);
  }

  confirmDetails(): void {
    if (this.selectedUnit()) {
      this.selectedUnit()?.patchValue({ status: 'Validé' });
    }
    this.closeDetailsDialog();
  }

  getOwners(unit: FormGroup): FormArray | null {
    return unit.get('owners') as FormArray || null;
  }
  
  get unitControls() {
    return this.units.controls;
  }
}
