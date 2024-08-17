import { Component, OnInit } from '@angular/core';
import { CondominiumService } from '../../services/condominium.service';
import { FormControl, FormGroup, TouchedChangeEvent, Validators, FormBuilder } from '@angular/forms';
import { ReactiveFormsModule } from '@angular/forms';
import { FloatLabelModule } from 'primeng/floatlabel';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';

@Component({
    selector: 'condominium',
    templateUrl: './condominium.component.html',
    styleUrl: './condominium.component.scss',
    standalone: true,
    imports: [ReactiveFormsModule, FloatLabelModule],
    providers: [HttpClient]
})

export class CondominiumComponent implements OnInit {

  private fromUrlCreateCondominium = environment.apiUrls.condominiumApi;

  public civilityTypes: any[] = [];
  public receivingMethods: any[] = [];

  formCategoriesName = {
    informations: [
      {id: "name", name: "Name", type: "text"},
      {id: "prefix", name: "Prefix", type: "text"},
      {id: "description", name: "Description", type: "text"},
    ],

    address: [
      {id: "street", name: "Street", type: "text"},
      {id: "number", name: "Number", type: "text"},
      {id: "addressComplement", name: "Address Complement", type: "text"},
      {id: "city", name: "City", type: "text"},
      {id: "postalCode", name: "Postal code", type: "text"},
      {id: "country", name: "Country", type: "text"},
    ],

    ftpBlueprint: [
      {id: "blueprint", name: "Blueprint FTP", type: "text"},
    ],

    concierge: [
      {id: "name", name: "Name", type: "text"},
      {id: "surname", name: "Surname", type: "text"},
      {id: "email", name: "Email", type: "text"},
      {id: "phone", name: "Phone", type: "text"},
      {id: "corporation", name: "Corporation", type: "text"},
      {id: "iban", name: "IBAN", type: "text"},
      {id: "birthdate", name: "Birthdate", type: "text"},
      {id: "civility", name: "Civility", type: "selector"},
      {id: "documentReceivingMethod", name: "Document Receiving Method", type: "selector"},
      {id: "reminderDelay", name: "Reminder Delay", type: "text"},
      {id: "reminderReceivingMethod", name: "Reminder Receiving Method", type: "selector"},
    ],
  };

  createCondominiumForm = new FormGroup({
    informations : new FormGroup({
      name : new FormControl(''),
      prefix : new FormControl(''),
      description : new FormControl(''),
    }),

    address : new FormGroup({
      street: new FormControl(''),
      number : new FormControl(''),
      addressComplement : new FormControl(''),
      city: new FormControl(''),
      postalCode: new FormControl(''),
      country : new FormControl(''),
    }),

    ftpBlueprint : new FormGroup({
      blueprint: new FormControl(''),
    }),

    concierge : new FormGroup({
      name: new FormControl(''),
      surname: new FormControl(''),
      email: new FormControl(''),
      corporation: new FormControl(''), //false par défaut
      phone: new FormControl(''),
      iban: new FormControl(''),
      birthdate: new FormControl(''),
      civility: new FormControl(''), //requête pour le selecteur (la table est pré-remplie)
      documentReceivingMethod: new FormControl(''), // ""
      reminderDelay: new FormControl(''),
      reminderReceivingMethod: new FormControl(''),
    }),
  });

  constructor(private http: HttpClient, private condominiumService: CondominiumService) {}

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

    this.condominiumService.getReceivingMethodOptions().subscribe({
      next: (data) => {
        this.receivingMethods = data;
      },
      error: (error) => {
        console.error('Failed to load receiving method options', error);
      }
    });
  }

  getCategories() {
    return Object.entries(this.formCategoriesName).map(([key, value]) => ({ key, value }));
  }

  ngOnInit(): void {
    this.loadOptions();
  }

  onSubmit(): void {
    console.log('Form Data:', this.createCondominiumForm.value);

    // Récupérer les données du formulaire
    const formData = this.createCondominiumForm.value;

    // Envoyer une requête HTTP POST avec les données du formulaire
    this.http.post(this.fromUrlCreateCondominium, formData).subscribe({
      next: (response) => {
        console.log('Réponse de l\'API:', response);
      },
      error: (error) => {
        console.error('Erreur lors de l\'envoi:', error);
      }
    });

    this.createCondominiumForm.reset();
  }
}