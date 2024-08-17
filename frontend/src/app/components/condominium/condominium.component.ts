import { Component, OnInit } from '@angular/core';
import { CondominiumService } from '../../services/condominium.service';
import { FormControl, FormGroup, TouchedChangeEvent, Validators } from '@angular/forms';
import { ReactiveFormsModule } from '@angular/forms';
import { FloatLabelModule } from 'primeng/floatlabel';

@Component({
    selector: 'condominium',
    templateUrl: './condominium.component.html',
    styleUrl: './condominium.component.scss',
    standalone: true,
    imports: [ReactiveFormsModule, FloatLabelModule],
    //providers: [FormGroup, FormControl]
})

export class CondominiumComponent implements OnInit {

  formCategoriesName = {
    formInformationsElements: [
      {id: "name", name: "Name", type: "text"},
      {id: "prefix", name: "Prefix", type: "text"},
      {id: "description", name: "Description", type: "text"},
    ],

    formAddressElements: [
      {id: "street", name: "Street", type: "text"},
      {id: "number", name: "Number", type: "text"},
      {id: "addressComplement", name: "Address Complement", type: "text"},
      {id: "city", name: "City", type: "text"},
      {id: "postalCode", name: "Postal code", type: "text"},
      {id: "country", name: "Country", type: "text"},
    ],

    formBlueprintElements: [
      {id: "blueprint", name: "Blueprint FTP", type: "text"},
    ],

    formConciergeElements: [
      {id: "name", name: "Name", type: "text"},
      {id: "surname", name: "Surname", type: "text"},
      {id: "email", name: "Email", type: "text"},
      {id: "phone", name: "Phone", type: "text"},
      {id: "corporation", name: "Corporation", type: "text"},
      {id: "iban", name: "IBAN", type: "text"},
      {id: "birthdate", name: "Birthdate", type: "text"},
      {id: "civility", name: "Civility", type: "text"},
      {id: "documentReceivingMethod", name: "Document Receiving Method", type: "text"},
      {id: "reminderDelay", name: "Reminder Delay", type: "text"},
      {id: "reminderReceivingMethod", name: "Reminder Receiving Method", type: "text"},
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

  //constructor(public CondominiumService: CondominiumService) {}

  getCategories() {
    return Object.entries(this.formCategoriesName).map(([key, value]) => ({ key, value }));
  }

  ngOnInit(): void {
    this.createCondominiumForm.events.subscribe( (event) => {
      if (event instanceof TouchedChangeEvent) {
        console.log('Touched: ' + event.touched)
      }
    });
  }

  onSubmit(): void {
    console.log('Form Data:', this.createCondominiumForm.value);

    this.createCondominiumForm.reset();
  }
}