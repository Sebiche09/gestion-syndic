import { Component , Output, EventEmitter, signal, computed } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule, FormArray } from '@angular/forms';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { AddressComponent } from '../address/address.component';
import { FtpblueprintComponent } from '../ftpblueprint/ftpblueprint.component';
import { ConciergeComponent } from '../concierge/concierge.component';
import { InformationComponent } from '../information/information.component';
import { CadastreComponent } from '../cadastre/cadastre.component';
import { UnitComponent } from '../unit/unit.component';
import { ButtonModule } from 'primeng/button';
import { StepsModule } from 'primeng/steps';
import { ConfirmationService, MenuItem, MessageService } from 'primeng/api';
import { ToastModule } from 'primeng/toast';
import { DialogModule } from 'primeng/dialog';
import { ConfirmPopupModule } from 'primeng/confirmpopup';
import { CoproprietaireComponent } from "../coproprietaire/coproprietaire.component";
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
    ConciergeComponent,
    CadastreComponent,
    UnitComponent,
    StepsModule,
    ToastModule,
    DialogModule,
    ConfirmPopupModule,
    CoproprietaireComponent,
    
],
  providers: [MessageService, ConfirmationService]
})
export class CondominiumComponent {
  @Output() closeDialog = new EventEmitter<void>();
  items: MenuItem[] = [];
  activeIndex = signal(0);
  createCondominiumForm: FormGroup;

  private fromUrlCreateCondominium = environment.apiUrls.condominiumApi;

  constructor(private http: HttpClient,
     private fb: FormBuilder,
     public messageService: MessageService,
     private confirmationService: ConfirmationService,
    private uniqueCheckService: UniqueCheckService) {
    this.createCondominiumForm = this.fb.group({
      informations: this.fb.group({
        name: ['', [Validators.required, Validators.minLength(3)],[UniqueValidator.checkNameUniqueness(this.uniqueCheckService)]],
        prefix: ['', [Validators.required], [UniqueValidator.checkPrefixUniqueness(this.uniqueCheckService)]],
        description: ['', [Validators.maxLength(500)]]
      }),

      address: this.fb.group({
        street: ['', Validators.required],
        address_complement: [''],
        city: ['', Validators.required],
        postal_code: ['', Validators.required],
        country: ['', Validators.required]
      }),

      ftpBlueprint: this.fb.group({
        blueprint: ['']
      }),

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
        country_concierge: ['']
      })
    });
  }

  // Getter pour le formGroup 'informations'
  get informationsForm(): FormGroup {
    return this.createCondominiumForm.get('informations') as FormGroup;
  }
  // Getter pour le formGroup 'address'
  get addressForm(): FormGroup {
    return this.createCondominiumForm.get('address') as FormGroup;
  }
  // Getter pour le formGroup 'ftpBlueprint'
  get ftpBlueprintForm(): FormGroup {
    return this.createCondominiumForm.get('ftpBlueprint') as FormGroup;
  }
  // Getter pour le formGroup 'lot'
  get units(): FormArray {
    return this.createCondominiumForm.get('units') as FormArray;
  }
  // Getter pour le formGroup 'concierge'
  get conciergeForm(): FormGroup {
    return this.createCondominiumForm.get('concierge') as FormGroup;
  }
  
  isStepValid(): boolean {
    // Vérifiez les validations pour chaque étape en fonction de l'index actif
    if (this.activeIndex() === 0) {
      const informationsForm = this.createCondominiumForm.get('informations');
      return informationsForm?.valid ?? false;
    } else if (this.activeIndex() === 2) {
      const addressForm = this.createCondominiumForm.get('address');
      return addressForm?.valid ?? false;
    }
    return true;
  }

  onActiveIndexChange(index: number) {
    this.activeIndex.set(index);
  }
  resetStep(){
    this.activeIndex.set(0);
  }
  nextStep() {
    this.activeIndex.update(i => i + 1);
  }

  previousStep() {
    this.activeIndex.update(i => i - 1);
  }
  confirmation(event: Event) {
    this.confirmationService.confirm({
        target: event.target as EventTarget,
        message: "Etes-vous sur de confirmer l'envoi du formulaire?",
        icon: 'pi pi-exclamation-triangle',
        accept: () => {
            this.messageService.add({ severity: 'info', summary: 'Confirmed', detail: 'Formulaire confirmé', life: 3000 });
            this.onSubmit();
        },
        reject: () => {
            this.messageService.add({ severity: 'error', summary: 'Rejected', detail: 'Formulaire rejeté', life: 3000 });
        }
    });
}

  ngOnInit() {
    this.items = [
        {
            label: 'Informations',
            command: (event: any) => this.messageService.add({severity:'info', summary:'Information Step', detail: event.item.label})
        },
        {
          label: 'Cadastre',
          command: (event: any) => this.messageService.add({severity:'info', summary:'Cadastre Step', detail: event.item.label})
        },
        {
            label: 'Adresse',
            command: (event: any) => this.messageService.add({severity:'info', summary:'Address Step', detail: event.item.label})
        },
        {
            label: 'Plans',
            command: (event: any) => this.messageService.add({severity:'info', summary:'FTP Blueprint Step', detail: event.item.label})
        },
        {
          label: 'Lots',
          command: (event: any) => this.messageService.add({severity:'info', summary:'Lot Step', detail: event.item.label})
        },
        {
          label: 'Co-propriétaires',
          command: (event: any) => this.messageService.add({severity:'info', summary:'Co-propriétaire Step', detail: event.item.label})
        },
        {
            label: 'Concierge',
            command: (event: any) => this.messageService.add({severity:'info', summary:'Concierge Step', detail: event.item.label})
        }
    ];
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
        console.log('Réponse de l\'API:', response);
        this.createCondominiumForm.reset();
      },
      error: (error: HttpErrorResponse) => {
        console.error('Erreur lors de l\'envoi:', error);
        this.getErrorSubmit(error);
      }
    });
  }
  onAddressExtracted(address: any) {
    if (address) {
      // Met à jour le formulaire d'adresse avec les valeurs récupérées
      this.addressForm.patchValue({
        street: address.street,
        postal_code: address.postal_code,
        city: address.city,
        country: address.country
      });
    }
  }
  onUnitExtracted(unit: any) {
    // Vider les unités existantes avant d'ajouter de nouvelles données
    this.units.clear();  
    if (unit && typeof unit === 'object') {
      // Parcourir les clés de l'objet 'unit'
      Object.keys(unit).forEach((unitKey: string) => {
        const unitDataArray = unit[unitKey];
    
        // Vérifier si la valeur associée à la clé est un tableau (et non 'null')
        if (Array.isArray(unitDataArray) && unitDataArray.length > 0) {
          unitDataArray.forEach((unitData: any) => {
            console.log(unitData.title)
            // Créer un nouveau FormGroup pour chaque unité
            const newUnit = this.fb.group({
              cadastralReference: [unitKey || '', Validators.required],
              unitType: [unitData.title || '', Validators.required],
              unitAddress: this.fb.group({
                street: [unitData.address?.street || '', Validators.required],
                postal_code: [unitData.address?.postal_code || '', Validators.required],
                city: [unitData.address?.city || '', Validators.required],
                country: [unitData.address?.country || '', Validators.required],
              }),
            });
    
            // Ajouter le nouveau lot au FormArray
            this.units.push(newUnit);
          });
        }
      });
    }
    
  }
  
  getErrorSubmit(error: HttpErrorResponse) {
    // Gestion des erreurs
  }
}
