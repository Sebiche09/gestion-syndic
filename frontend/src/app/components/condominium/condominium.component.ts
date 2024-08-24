import { Component , Output, EventEmitter} from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { AddressComponent } from '../address/address.component';
import { FtpblueprintComponent } from '../ftpblueprint/ftpblueprint.component';
import { ConciergeComponent } from '../concierge/concierge.component';
import { InformationComponent } from '../information/information.component';
import { ButtonModule } from 'primeng/button';
import { StepsModule } from 'primeng/steps';
import { ConfirmationService, MenuItem, MessageService } from 'primeng/api';
import { ToastModule } from 'primeng/toast';
import { DialogModule } from 'primeng/dialog';
import { ConfirmPopupModule } from 'primeng/confirmpopup';


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
    StepsModule,
    ToastModule,
    DialogModule,
    ConfirmPopupModule
  ],
  providers: [MessageService, ConfirmationService]
})
export class CondominiumComponent {
  @Output() closeDialog = new EventEmitter<void>();
  items: MenuItem[] = [];
  activeIndex: number = 0;
  createCondominiumForm: FormGroup;

  private fromUrlCreateCondominium = environment.apiUrls.condominiumApi;

  constructor(private http: HttpClient, private fb: FormBuilder, public messageService: MessageService, private confirmationService: ConfirmationService) {
    this.createCondominiumForm = this.fb.group({
      informations: this.fb.group({
        name: ['', Validators.required],
        prefix: ['', Validators.required],
        description: ['']
      }),

      address: this.fb.group({
        street: ['', Validators.required],
        number: ['', Validators.required],
        address_complement: [''],
        city: ['', Validators.required],
        postal_code: ['', Validators.required],
        country: ['', Validators.required]
      }),

      ftpBlueprint: this.fb.group({
        blueprint: ['']
      }),

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
        number_concierge: [''],
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
  // Getter pour le formGroup 'concierge'
  get conciergeForm(): FormGroup {
    return this.createCondominiumForm.get('concierge') as FormGroup;
  }
  
  onActiveIndexChange(event: number) {
    this.activeIndex = event;
  }
  resetActiveIndex() {
    this.activeIndex = 0;
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
            label: 'Information',
            command: (event: any) => this.messageService.add({severity:'info', summary:'Information Step', detail: event.item.label})
        },
        {
            label: 'Address',
            command: (event: any) => this.messageService.add({severity:'info', summary:'Address Step', detail: event.item.label})
        },
        {
            label: 'FTP Blueprint',
            command: (event: any) => this.messageService.add({severity:'info', summary:'FTP Blueprint Step', detail: event.item.label})
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

  getErrorSubmit(error: HttpErrorResponse) {
    // Gestion des erreurs
  }
}
