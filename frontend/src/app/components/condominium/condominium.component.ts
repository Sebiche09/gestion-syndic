import { Component, Output, EventEmitter, signal, inject } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule, FormArray} from '@angular/forms';
//import primeng modules
import { ConfirmationService, MenuItem, MessageService } from 'primeng/api';
import { ButtonModule } from 'primeng/button';
import { StepsModule } from 'primeng/steps';
import { DialogModule } from 'primeng/dialog';
import { ConfirmPopupModule } from 'primeng/confirmpopup';
//import standalone components
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { AddressComponent } from '../address/address.component';
import { FtpblueprintComponent } from '../ftpblueprint/ftpblueprint.component';
import { InformationComponent } from '../information/information.component';
import { CadastreComponent } from '../cadastre/cadastre.component';
import { UnitComponent } from '../unit/unit.component';
// Import unique validator
import { UniqueCheckService } from '../../services/unique-check.service';
import { UniqueValidator } from '../../validators/unique-validator';

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
  ],
  providers: [MessageService, ConfirmationService],
})
export class CondominiumComponent {
  @Output() closeDialog = new EventEmitter<void>(); // Événement pour fermer le dialog parent (app.component.html)

  // Injection de dépendances via `inject` pour plus de clarté
  private fb = inject(FormBuilder);
  private http = inject(HttpClient);
  private messageService = inject(MessageService);
  private confirmationService = inject(ConfirmationService);
  private uniqueCheckService = inject(UniqueCheckService);

  items: MenuItem[] = [
    'Informations',
    'Cadastre',
    'Adresse',
    'Plans',
    'Lots',
    'Co-propriétaires',
    'Concierge',
  ].map((label) => ({ label })); // Création des titres des étapes

  activeIndex = signal(0); // Index de l'étape active

  createCondominiumForm: FormGroup; // Formulaire de création de copropriété

  private fromUrlCreateCondominium = environment.apiUrls.condominiumApi;

  constructor() {
    this.createCondominiumForm = this.initForm();
  }

  //----------------- Crée un formulaire vide ------------------
  private initForm(): FormGroup {
    return this.fb.group({
      informations: this.fb.group({
        name: [
          '',
          [Validators.required, Validators.minLength(3)],
          [UniqueValidator.checkNameUniqueness(this.uniqueCheckService)],
        ],
        prefix: [
          '',
          [Validators.required],
          [UniqueValidator.checkPrefixUniqueness(this.uniqueCheckService)],
        ],
        exercice: ['', [Validators.required]],  
        description: ['', [Validators.maxLength(500)]],
      }),
      address: this.fb.group({
        street: ['', Validators.required],
        address_complement: [''],
        city: ['', Validators.required],
        postal_code: ['', Validators.required],
        country: ['', Validators.required],
      }),
      ftpBlueprint: this.fb.group({ blueprint: [''] }),
      units: this.fb.array([]),
      concierge: this.fb.group({
        isThereConcierge: [false],
        name: [''],
        surname: [''],
        email: [''],
        corporation: [false],
        phone: [''],
        iban: [''],
        birthdate: [''],
        civility: [''],
        document_receiving_method: [''],
        reminder_delay: [40],
        reminder_receiving_method: [''],
        street_concierge: [''],
        address_complement_concierge: [''],
        city_concierge: [''],
        postal_code_concierge: [''],
        country_concierge: [''],
      }),
    });
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
  get conciergeForm(): FormGroup {
    return this.createCondominiumForm.get('concierge') as FormGroup;
  }
  //---------------------------------------------------------

  isStepValid(): boolean {
    return this.activeIndex() === 0 ? this.informationsForm.valid : true;
  }
  resetStep() {
    this.activeIndex.set(0);
  }
  nextStep() {
    this.activeIndex.update((i) => i + 1);
  }
  previousStep() {
    this.activeIndex.update((i) => i - 1);
  }

  confirmation(event: Event) {
    this.confirmationService.confirm({
      target: event.target as EventTarget,
      message: "Etes-vous sur de confirmer l'envoi du formulaire?",
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
    this.closeDialog.emit();
    console.log('Formulaire soumis:', this.createCondominiumForm.value);
    let formData = this.createCondominiumForm.value;

    if (formData.concierge) {
      formData.concierge.corporation = !!formData.concierge.corporation;
    }

    this.http.post(this.fromUrlCreateCondominium, formData).subscribe({
      next: (response) => {
        console.log("Réponse de l'API:", response);
        this.createCondominiumForm.reset();
      },
      error: (error: HttpErrorResponse) => {
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
      console.log(text);

      // Met à jour le formulaire d'adresse avec les valeurs récupérées (adresse principale du bâtiment)
      this.addressForm.patchValue({
        street: text.address.street,
        postal_code: text.address.postal_code,
        city: text.address.city,
        country: text.address.country,
      });

      // Vider les unités existantes
      this.units.clear();

      // Parcourir les unités (lots) dans 'unit'
      const unitData = text.unit; // Récupère les unités depuis le JSON
      if (unitData && typeof unitData === 'object') {
        Object.keys(unitData).forEach((unitKey: string) => {
          const unitInfo = unitData[unitKey];

          // Récupère le complément d'adresse
          const complement = unitInfo?.complement || '';

          // Si 'unitInfo.owners' existe et est un tableau (donc non 'null')
          if (Array.isArray(unitInfo.owners) && unitInfo.owners.length > 0) {
            // Crée un formulaire pour chaque unité
            const newUnit = this.fb.group({
              cadastralReference: [unitKey || '', Validators.required],
              status: ['brouillon'],
              unitAddress: this.fb.group({
                // Combiner l'adresse principale avec le complément d'adresse
                street: [text.address.street || '', Validators.required],
                complement: [complement],
                postal_code: [
                  text.address.postal_code || '',
                  Validators.required,
                ],
                city: [text.address.city || '', Validators.required],
                country: [text.address.country || '', Validators.required],
              }),
              // Crée un FormArray pour les propriétaires
              owners: this.fb.array([]), // Initialisation vide, on le remplira après
            });

            // Ajoute les propriétaires au FormArray
            const ownersArray = newUnit.get('owners') as FormArray;
            unitInfo.owners.forEach((owner: any) => {
              const ownerGroup = this.fb.group({
                lastName: [owner.last_name || '', Validators.required],
                firstName: [owner.first_name || '', Validators.required],
                title: [owner.title || '', Validators.required],
                address: this.fb.group({
                  street: [owner.address?.street || '', Validators.required],
                  postal_code: [
                    owner.address?.postal_code || '',
                    Validators.required,
                  ],
                  city: [owner.address?.city || '', Validators.required],
                  country: [owner.address?.country || '', Validators.required],
                }),
              });

              // Ajoute chaque propriétaire au FormArray
              ownersArray.push(ownerGroup);
            });

            // Ajouter le nouveau lot avec les propriétaires au FormArray des unités
            this.units.push(newUnit);
          }
        });
      }
    }
  }
}
