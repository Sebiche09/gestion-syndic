import { Component, Input } from '@angular/core';
import { FormBuilder, FormGroup, FormArray , Validators} from '@angular/forms';
import { TableModule } from 'primeng/table';

@Component({
  selector: 'app-unit',
  standalone: true,
  imports: [TableModule],
  templateUrl: './unit.component.html',
  styleUrl: './unit.component.scss'
})
export class UnitComponent {
  @Input() units!: FormArray;

  constructor(private fb: FormBuilder) {}
  
  // Méthode pour créer un lot (utilisée pour ajouter un nouveau lot au FormArray)
  createUnit(): FormGroup {
    return this.fb.group({
      unitName: [''],
      cadastralReference: [''],
      unitAddress: this.fb.group({
        street: ['', Validators.required],
        address_complement: [''],
        city: ['', Validators.required],
        postal_code: ['', Validators.required],
        country: ['', Validators.required]
      }),
      unitType: ['', Validators.required],
      foor : ['', Validators.required],
      description : ['', Validators.required],
      quota : ['', Validators.required],
      occupants: this.fb.array([]),
    });
  }
  

  // Ajoute un nouveau lot dans le FormArray
  addUnit(): void {
    this.units.push(this.createUnit());
    console.log(this.units)
  }

  // Supprime un lot du FormArray
  removeUnit(index: number): void {
    this.units.removeAt(index);
  }

  // Retourne les contrôles du FormArray pour itération dans le template
  get unitControls() {
    return this.units.controls;
  }
}
