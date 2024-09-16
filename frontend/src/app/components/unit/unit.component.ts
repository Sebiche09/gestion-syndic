import { Component, Input } from '@angular/core';
import { FormBuilder, FormGroup, FormArray } from '@angular/forms';
import { TableModule } from 'primeng/table';
import { TagModule } from 'primeng/tag';
import { DialogModule } from 'primeng/dialog';

@Component({
  selector: 'app-unit',
  standalone: true,
  imports: [TableModule, TagModule, DialogModule],
  templateUrl: './unit.component.html',
  styleUrls: ['./unit.component.scss']
})
export class UnitComponent {
  @Input() units!: FormArray;

  displayDetailsDialog = false;  // Gère l'état de la boîte de dialogue des détails
  selectedUnit: FormGroup | null = null;  // Référence à l'unité sélectionnée

  constructor(private fb: FormBuilder) {}

  // Afficher les détails du lot
  viewDetails(index: number): void {
    this.selectedUnit = this.units.at(index) as FormGroup;
    this.displayDetailsDialog = true;  // Ouvrir uniquement la boîte de dialogue des détails
  }

  // Fermer uniquement la boîte de dialogue des détails sans confirmation
  closeDetailsDialog(): void {
    this.displayDetailsDialog = false;  // Ferme seulement le dialogue de détails
    this.selectedUnit = null;  // Réinitialiser l'unité sélectionnée
  }

  // Confirmer les informations et mettre à jour le statut
  confirmDetails(): void {
    if (this.selectedUnit) {
      this.selectedUnit.patchValue({
        status: 'validé'  // Met à jour le statut de l'unité
      });
    }
    this.displayDetailsDialog = false;  // Fermer le dialogue de détails
    this.selectedUnit = null;
  }

  // Retourner la liste des propriétaires pour un lot
  getOwners(unit: FormGroup): FormArray | null {
    return unit.get('owners') as FormArray || null;
  }

  // Supprimer un lot du FormArray
  removeUnit(index: number): void {
    this.units.removeAt(index);
  }

  // Retourne les contrôles du FormArray pour itération dans le template
  get unitControls() {
    return this.units.controls;
  }
}
