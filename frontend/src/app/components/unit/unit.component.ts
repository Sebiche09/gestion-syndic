import { Component, Input, EventEmitter, Output } from '@angular/core';
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
  displayDetailsDialog = false;  
  selectedUnit: FormGroup | null = null;  

  constructor(private fb: FormBuilder) {}

  viewDetails(index: number): void {
    this.selectedUnit = this.units.at(index) as FormGroup;
    this.displayDetailsDialog = true;
  }

  closeDetailsDialog(): void {
    this.displayDetailsDialog = false;
    this.selectedUnit = null;
  }

  // Confirmer les détails et définir le statut à "Validé"
  confirmDetails(): void {
    if (this.selectedUnit) {
      this.selectedUnit.patchValue({
        status: 'Validé'  // Statut défini directement à "Validé"
      });
    }
    this.displayDetailsDialog = false;
    this.selectedUnit = null;
  }

  getOwners(unit: FormGroup): FormArray | null {
    return unit.get('owners') as FormArray || null;
  }

  get unitControls() {
    return this.units.controls;
  }
}
