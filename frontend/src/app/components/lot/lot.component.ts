import { Component, Input } from '@angular/core';
import { FormBuilder, FormGroup, FormArray , Validators} from '@angular/forms';
import { TableModule } from 'primeng/table';

@Component({
  selector: 'app-lot',
  standalone: true,
  imports: [TableModule],
  templateUrl: './lot.component.html',
  styleUrl: './lot.component.scss'
})
export class LotComponent {
  @Input() lots!: FormArray;

  constructor(private fb: FormBuilder) {}
  
  // Méthode pour créer un lot (utilisée pour ajouter un nouveau lot au FormArray)
  createLot(): FormGroup {
    return this.fb.group({
      lotName: [''],
      cadastralReference: [''],
      lotAddress: this.fb.group({
        street: ['', Validators.required],
        address_complement: [''],
        city: ['', Validators.required],
        postal_code: ['', Validators.required],
        country: ['', Validators.required]
      })
    });
  }

  // Ajoute un nouveau lot dans le FormArray
  addLot(): void {
    this.lots.push(this.createLot());
  }

  // Supprime un lot du FormArray
  removeLot(index: number): void {
    this.lots.removeAt(index);
  }

  // Retourne les contrôles du FormArray pour itération dans le template
  get lotControls() {
    return this.lots.controls;
  }
}
