import { Component, Input, Output, EventEmitter, OnInit } from '@angular/core';
import { FormArray, FormGroup, ReactiveFormsModule } from '@angular/forms';
import { TableModule } from 'primeng/table';
import { ToastModule } from 'primeng/toast';
import { TagModule } from 'primeng/tag';
import { ButtonModule } from 'primeng/button';
import { DialogModule } from 'primeng/dialog';
import { CalendarModule } from 'primeng/calendar';
import { DropdownModule } from 'primeng/dropdown';
import { InputMaskModule } from 'primeng/inputmask';
import { InputNumberModule } from 'primeng/inputnumber';
import { CheckboxModule } from 'primeng/checkbox';
import { InputTextModule } from 'primeng/inputtext';

@Component({
  selector: 'app-coproprietaire',
  standalone: true,
  imports: [
    TableModule,
    ToastModule,
    TagModule,
    ButtonModule,
    DialogModule,
    ReactiveFormsModule,
    CalendarModule,
    DropdownModule,
    InputMaskModule,
    InputNumberModule,
    CheckboxModule,
    InputTextModule,
  ],
  templateUrl: './coproprietaire.component.html',
  styleUrls: ['./coproprietaire.component.scss'],
})
export class CoproprietaireComponent implements OnInit{
  @Input() occupants!: FormArray;

  @Output() previous = new EventEmitter<void>();
  @Output() next = new EventEmitter<void>();

  civility: any[] | undefined;
  documentReceivingMethods: any[] | undefined;
  reminderReceivingMethods: any[] | undefined;

  detailsDialogVisible = false;
  selectedOccupantForm: FormGroup | null = null;
  selectedIndex: number | null = null;


  ngOnInit(): void {
    this.civility = [
      { label: 'Monsieur', value: 'Monsieur' },
      { label: 'Madame', value: 'Madame' },
    ];
  
    this.documentReceivingMethods = [
      { label: 'Email', value: 'email' },
      { label: 'Courrier', value: 'courrier' },
      { label: 'Fax', value: 'fax' },
      { label: 'Recommandé', value: 'recommande' },
    ];
  
    this.reminderReceivingMethods = [
      { label: 'Email', value: 'email' },
      { label: 'Courrier', value: 'courrier' },
      { label: 'SMS', value: 'sms' },
      { label: 'Fax', value: 'fax' },
      { label: 'Recommandé', value: 'recommande' },
    ];
  }

  viewDetails(index: number): void {
    const occupant = this.occupants.at(index) as FormGroup;
    if (occupant) {
      this.selectedOccupantForm = occupant;
      this.selectedIndex = index;
      this.detailsDialogVisible = true;
    } else {
      console.error("Occupant non trouvé pour l'index :", index);
    }
  }

  // Fermer le dialogue
  closeDetailsDialog(): void {
    this.detailsDialogVisible = false;
    this.selectedOccupantForm = null!;
    this.selectedIndex = null;
  }

  confirmDetails(): void {
    if (this.selectedIndex !== null) {
      this.selectedOccupantForm?.patchValue({ status: 'validé' });
      console.log('Occupant validé :', this.selectedOccupantForm?.value);
      const occupant = this.occupants.at(this.selectedIndex) as FormGroup;
      occupant.patchValue(this.selectedOccupantForm?.value);
    }
    this.closeDetailsDialog();
  }

  previousStep(): void {
    this.previous.emit();
  }

  nextStep(): void {
    this.next.emit();
  }

  allOccupantsValidated(): boolean {
    return this.occupants.controls.every(
      (occupant) => occupant.get('status')?.value === 'validé'
    );
  }
}
