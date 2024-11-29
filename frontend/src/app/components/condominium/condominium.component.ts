import { Component, Output, EventEmitter, signal, inject } from '@angular/core';
import {FormGroup, FormArray, ReactiveFormsModule} from '@angular/forms';
//import primeng modules
import { ConfirmationService, MenuItem, MessageService } from 'primeng/api';
import { ButtonModule } from 'primeng/button';
import { StepsModule } from 'primeng/steps';
import { DialogModule } from 'primeng/dialog';
import { ConfirmPopupModule } from 'primeng/confirmpopup';
import { AddressComponent } from '../address/address.component';
import { FtpblueprintComponent } from '../ftpblueprint/ftpblueprint.component';
import { InformationComponent } from '../information/information.component';
import { CadastreComponent } from '../cadastre/cadastre.component';
import { UnitComponent } from '../unit/unit.component';
import { CondominiumService } from '../../services/condominium.service';
import { CoproprietaireComponent } from '../coproprietaire/coproprietaire.component';
import { title } from 'node:process';
@Component({
  selector: 'app-condominium',
  templateUrl: './condominium.component.html',
  styleUrls: ['./condominium.component.scss'],
  standalone: true,
  imports: [
    ReactiveFormsModule,
    ButtonModule,
    AddressComponent,
    InformationComponent,
    FtpblueprintComponent,
    CadastreComponent,
    UnitComponent,
    StepsModule,
    DialogModule,
    ConfirmPopupModule,
    CoproprietaireComponent
  ],
  providers: [MessageService, ConfirmationService],
})
export class CondominiumComponent {
  @Output() closeRequest = new EventEmitter<void>(); // Événement pour fermer le dialog parent (app.component.html)

  private messageService = inject(MessageService);
  private confirmationService = inject(ConfirmationService);

  items: MenuItem[] = [
    'Informations',
    'Cadastre',
    'Adresse',
    'Plans',
    'Co-propriétaires',
    'Lots',
  ].map((label) => ({ label })); // Création des titres des étapes

  activeIndex = signal(0); // Index de l'étape active

  createCondominiumForm: FormGroup; // Formulaire de création de copropriété

  constructor(private condominiumService: CondominiumService) {
    this.createCondominiumForm = this.condominiumService.createCondominiumForm();
  }

  //--------------------- GETTER --------------------------
  get informationsForm(): FormGroup {
    return this.createCondominiumForm.get('informations') as FormGroup;
  }
  get addressForm(): FormGroup {
    return this.createCondominiumForm.get('address') as FormGroup;
  }
  get ftpBlueprintForm(): FormGroup {
    return this.createCondominiumForm.get('ftpBlueprint') as FormGroup;
  }
  get units(): FormArray {
    return this.createCondominiumForm.get('units') as FormArray;
  }
  get occupants(): FormArray {
    return this.createCondominiumForm.get('occupants') as FormArray;
  }
  //---------------------------------------------------------
  
  resetStep() {
    this.activeIndex.set(0);
  }
  nextStep() {
    this.activeIndex.update((i) => i + 1);
  }
  previousStep() {
    this.activeIndex.update((i) => i - 1);
  }
  clearFormAndIndex() {
    this.createCondominiumForm.reset();
    this.activeIndex.set(0);
  }

  confirmation(event: Event) {
    this.confirmationService.confirm({
      target: event.target as EventTarget,
      message: "Etes-vous sûr de confirmer l'envoi du formulaire?",
      icon: 'pi pi-exclamation-triangle',
      accept: () => {
        this.messageService.add({
          severity: 'info',
          summary: 'Confirmed',
          detail: 'Formulaire confirmé',
          life: 3000,
        });
        this.onSubmit();
      },
      reject: () => {
        this.messageService.add({
          severity: 'error',
          summary: 'Rejected',
          detail: 'Formulaire rejeté',
          life: 3000,
        });
      },
    });
  }

  onSubmit(): void {
    this.closeRequest.emit();
    this.condominiumService.submitCondominium(this.createCondominiumForm.value)
      .subscribe({
        next: (response) => {
          console.log("Réponse de l'API:", response);
          this.createCondominiumForm.reset();
        },
        error: (error) => {
          console.error("Erreur lors de l'envoi:", error);
          this.messageService.add({
            severity: 'error',
            summary: "Erreur lors de l'envoi",
            detail: error.message,
          });
        },
      });
  }
  
  onTextExtracted(text: any) {
    if (text) {
      this.addressForm.patchValue({
        street: text.address.street,
        postal_code: text.address.postal_code,
        city: text.address.city,
        country: text.address.country,
      });
  
      const occupantsArray = this.createCondominiumForm.get('occupants') as FormArray;
  
      const unitData = text.unit;
      if (unitData && typeof unitData === 'object') {
        Object.keys(unitData).forEach((unitKey: string) => {
          const unitInfo = unitData[unitKey];
          if (Array.isArray(unitInfo.owners) && unitInfo.owners.length > 0) {
            const transformedOwners = unitInfo.owners.map((owner: any) => {
              let occupantIndex = occupantsArray.controls.findIndex((control) =>
                control.get('name')?.value === owner.last_name &&
                control.get('surname')?.value === owner.first_name &&
                control.get('title')?.value === owner.title

              );
  
              if (occupantIndex === -1) {
                this.condominiumService.addOccupant(occupantsArray, {
                  name: owner.last_name,
                  surname: owner.first_name,
                  title: owner.title,
                });
                occupantIndex = occupantsArray.length - 1;
              }
  
              // Inclure nom et prénom directement
              return {
                name: owner.last_name,
                surname: owner.first_name,
                title: owner.title || '',
                quota: owner.quota || 0,
                administrator: owner.administrator || false,
              };
            });
  
            unitInfo.owners = transformedOwners;
  
            this.condominiumService.addUnit(this.units, {
              unitKey,
              ...unitInfo,
            });
          } else {
            console.warn(`Unité ${unitKey} n'a pas de propriétaires définis.`);
          }
        });
      }
    }
  }
  
  
  
  
    
}
