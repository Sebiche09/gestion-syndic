import { Component, OnInit , Input } from '@angular/core';
import { FormGroup , ReactiveFormsModule, FormsModule} from '@angular/forms';
import { CommonModule } from '@angular/common';
import { FloatLabelModule } from 'primeng/floatlabel';
import { InputTextModule } from 'primeng/inputtext';
import { InputMaskModule } from 'primeng/inputmask';
import { InputSwitchModule } from 'primeng/inputswitch';
import { CalendarModule } from 'primeng/calendar';
import { DropdownModule } from 'primeng/dropdown';
import { KeyFilterModule } from 'primeng/keyfilter';

import { CardModule } from 'primeng/card';


import { CondominiumService } from '../../services/condominium.service';

@Component({
  selector: 'app-concierge',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    FormsModule,
    FloatLabelModule, 
    InputTextModule, 
    InputMaskModule, 
    InputSwitchModule, 
    CalendarModule, 
    DropdownModule, 
    KeyFilterModule,
    CardModule],
  templateUrl: './concierge.component.html',
  styleUrl: './concierge.component.scss'
})
export class ConciergeComponent implements OnInit{
  @Input() conciergeForm!: FormGroup;
  isThereConcierge = false;


  // Function to switch the concierge form
  conciergeSwitch(event: any) {
    const isThereConciergeControl = this.conciergeForm.get('isThereConcierge');
    if (isThereConciergeControl) {
      isThereConciergeControl.setValue(event.checked);
      
      if (event.checked) {
        // Activer tous les contrôles du formulaire concierge
        this.conciergeForm.enable();
      } else {
        // Désactiver tous les contrôles du formulaire concierge
        this.conciergeForm.disable();
        // Garder le contrôle `isThereConcierge` activé pour que l'utilisateur puisse réactiver le formulaire
        isThereConciergeControl.enable();
      }
    }
  }

  constructor(private condominiumService: CondominiumService) {}
  public civilityTypes: any[] = [];
  public reminderReceivingMethods: any[] = [];
  public documentReceivingMethods: any[] = [];

  ngOnInit(): void{
    this.loadOptions();
  }

  //Get civilities types and receiving methods from DB
  loadOptions(): void {
    this.condominiumService.getCivilityOptions().subscribe({
      next: (data) => {
        this.civilityTypes = data;
        console.log('Civility options loaded', data);
      },
      error: (error) => {
        console.error('Failed to load civility options', error);
      }
    });

    this.condominiumService.getReminderReceivingMethodOptions().subscribe({
      next: (data) => {
        this.reminderReceivingMethods = data;
      },
      error: (error) => {
        console.error('Failed to load reminder receiving method options', error);
      }
    });

    this.condominiumService.getDocumentReceivingMethodOptions().subscribe({
      next: (data) => {
        this.documentReceivingMethods = data;
      },
      error: (error) => {
        console.error('Failed to load receiving method options', error);
      }
    });
  }

  getLabel(field: string): string {
    const labels: Record<string, string> = {
      name: 'Prénom',
      surname: 'Nom',
      email: 'Adresse mail',
      corporation: 'Entreprise',
      phone: 'Téléphone',
      iban: 'Iban',
      birthdate: 'Date de naissance',
      civility: 'Civilité',
      document_receiving_method: 'Méthode d\'envoi des documents',
      reminder_delay: 'Délai rappels',
      reminder_receiving_method: 'Méthode d\'envoi rappels'
    };
    return labels[field] || field;
  }
}
