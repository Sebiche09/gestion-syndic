import { Component, OnInit } from '@angular/core';
import { CondominiumService } from '../../services/condominium.service';
import { FormControl, FormGroup, FormsModule, Validators } from '@angular/forms';
import { ReactiveFormsModule } from '@angular/forms';
import { FloatLabelModule } from 'primeng/floatlabel';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';

@Component({
    selector: 'condominium',
    templateUrl: './condominium.component.html',
    styleUrl: './condominium.component.scss',
    standalone: true,
    imports: [ReactiveFormsModule, FloatLabelModule, FormsModule],
    providers: [HttpClient]
})

export class CondominiumComponent implements OnInit {

  private fromUrlCreateCondominium = environment.apiUrls.condominiumApi;

  isCheckboxChecked: boolean = false;

  public civilityTypes: any[] = [];
  public reminderReceivingMethods: any[] = [];
  public documentReceivingMethods: any[] = [];

  formCategoriesName = {
    informations: [
      {id: "name", name: "Name", type: "text"},
      {id: "prefix", name: "Prefix", type: "text"},
      {id: "description", name: "Description", type: "text"},
    ],

    address: [
      {id: "street", name: "Street", type: "text"},
      {id: "number", name: "Number", type: "text"},
      {id: "address_complement", name: "Address Complement", type: "text"},
      {id: "city", name: "City", type: "text"},
      {id: "postal_code", name: "Postal code", type: "text"},
      {id: "country", name: "Country", type: "text"},
    ],

    ftpBlueprint: [
      {id: "blueprint", name: "Blueprint FTP", type: "text"},
    ],

    concierge: [
      {id: "name", name: "Name", type: "text"},
      {id: "surname", name: "Surname", type: "text"},
      {id: "email", name: "Email", type: "email"},
      {id: "phone", name: "Phone", type: "text"},
      {id: "corporation", name: "Corporation", type: "checkbox"},
      {id: "iban", name: "IBAN", type: "text"},
      {id: "birthdate", name: "Birthdate", type: "date"},
      {id: "civility", name: "Civility", type: "selector"},
      {id: "document_receiving_method", name: "Document Receiving Method", type: "selector"},
      {id: "reminder_delay", name: "Reminder Delay", type: "number"},
      {id: "reminder_receiving_method", name: "Reminder Receiving Method", type: "selector"},

      // Concierge address
      {id: "street_concierge", name: "Street", type: "text"},
      {id: "number_concierge", name: "Number", type: "text"},
      {id: "address_complement_concierge", name: "Address Complement", type: "text"},
      {id: "city_concierge", name: "City", type: "text"},
      {id: "postal_code_concierge", name: "Postal code", type: "text"},
      {id: "country_concierge", name: "Country", type: "text"},
    ],
  };

  createCondominiumForm = new FormGroup({
    informations: new FormGroup({
      name: new FormControl('aaaaaaaaaaa'),
      prefix: new FormControl('aaaaaaaaaaa'),
      description : new FormControl('aaaaaaaaaaa'),
    }),

    address: new FormGroup({
      street: new FormControl('aaaaaaaaaaa'),
      number: new FormControl('aaaaaaaaaaa'),
      address_complement: new FormControl('aaaaaaaaaaa'),
      city: new FormControl('aaaaaaaaaaa'),
      postal_code: new FormControl('aaaaaaaaaaa'),
      country: new FormControl('aaaaaaaaaaa'),
    }),

    ftpBlueprint: new FormGroup({
      blueprint: new FormControl('aaaaaaaaaaa'),
    }),

    concierge: new FormGroup({
      name: new FormControl('aaaaaaaaaaa'),
      surname: new FormControl('aaaaaaaaaaa'),
      email: new FormControl('aaaaaaaaaaa'),
      corporation: new FormControl(false), //false par défaut
      phone: new FormControl('aaaaaaaaaaa'),
      iban: new FormControl('aaaaaaaaaaa'),
      birthdate: new FormControl('aaaaaaaaaaa'),
      civility: new FormControl(''), //requête pour le selecteur (la table est pré-remplie)
      document_receiving_method: new FormControl(''), // ""
      reminder_delay: new FormControl('40'),
      reminder_receiving_method: new FormControl(''),

      //concierge address
      street_concierge: new FormControl('aaaaaaaaaaa'),
      number_concierge: new FormControl('aaaaaaaaaaa'),
      address_complement_concierge: new FormControl('aaaaaaaaaaa'),
      city_concierge: new FormControl('aaaaaaaaaaa'),
      postal_code_concierge: new FormControl('aaaaaaaaaaa'),
      country_concierge: new FormControl('aaaaaaaaaaa'),
    }),
  });

  constructor(private http: HttpClient, private condominiumService: CondominiumService) {}

  //fonction init
  ngOnInit(): void {
    this.loadOptions();
  }

  //Get civilities types and receiving methods from DB
  loadOptions(): void {
    this.condominiumService.getCivilityOptions().subscribe({
      next: (data) => {
        this.civilityTypes = data;
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

  onCheckboxChange(event: any) {
    this.isCheckboxChecked = event.target.checked;
  }

  getCategories() {
    return Object.entries(this.formCategoriesName).map(([key, value]) => ({ key, value }));
  }

  onSubmit(): void {
    console.log('Form Data:', this.createCondominiumForm.value);

    // Récupérer les données du formulaire
    let formData = this.createCondominiumForm.value;

    // convert in boolean
    if(formData.concierge) {
      formData.concierge.corporation = !!formData.concierge.corporation;
    }

    // Envoyer une requête HTTP POST avec les données du formulaire
    this.http.post(this.fromUrlCreateCondominium, formData).subscribe({
      next: (response) => {
        console.log('Réponse de l\'API:', response);
      },
      error: (error) => {
        console.error('Erreur lors de l\'envoi:', error);
      }
    });

    //this.createCondominiumForm.reset();
  }
}