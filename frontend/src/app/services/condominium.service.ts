import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';
import { FormBuilder, FormGroup, Validators, FormArray } from '@angular/forms';
import { UniqueValidator } from '../validators/unique-validator';
import { UniqueCheckService } from '../services/unique-check.service';
@Injectable({
  providedIn: 'root'
})

export class CondominiumService {
  private fb = inject(FormBuilder);
  private uniqueCheckService = inject(UniqueCheckService);

  constructor(private http: HttpClient) {}

  submitCondominium(formData: any): Observable<any> {
    const url = environment.apiUrls.condominiumApi;
    return this.http.post(url, formData);
  }

  createCondominiumForm(): FormGroup {
    return this.fb.group({
      informations: this.fb.group({
        name: [
          'test', 
          [Validators.required, Validators.minLength(3), Validators.maxLength(50)],
          [UniqueValidator.checkNameUniqueness(this.uniqueCheckService)],
        ],
        prefix: [
          'test',
          [Validators.required, Validators.maxLength(50)],
          [UniqueValidator.checkPrefixUniqueness(this.uniqueCheckService)],
        ],
        description: ['test', [Validators.maxLength(500)]],
      }),
      address: this.fb.group({
        street: ['', Validators.required],
        address_complement: [''],
        city: ['', Validators.required],
        postal_code: ['', Validators.required],
        country: ['Belgique', Validators.required],
      }),
      units: this.fb.array([]),
      occupants: this.fb.array([]),
    });
  }

  addOccupant(occupants: FormArray, occupantData: any): void {
    const newOccupant = this.fb.group({
      name: [occupantData.name || '', Validators.required],
      surname: [occupantData.surname || '', Validators.required],
      birthDate: [occupantData.birthDate || '01.01.2024', Validators.required],
      email: [occupantData.email || 'test', [Validators.email]],
      corporation: [occupantData.corporation || false],
      phone: [occupantData.phone || '32470542056'],
      iban: [occupantData.iban || 'BE68539007547034'],
      civility: [occupantData.civility || 'Monsieur', Validators.required],
      address: this.fb.group({
        street: [occupantData.address?.street || '', Validators.required],
        postal_code: [occupantData.address?.postal_code || '', Validators.required],
        city: [occupantData.address?.city || '', Validators.required],
        country: [occupantData.address?.country || '', Validators.required],
      }),
      document_receiving_method: ['fax'],
      reminder_delay: [0],
      reminder_receiving_method: ['fax'],
      isConcierge: [occupantData.concierge || false],
      status: ['brouillon'],
    });
    
    occupants.push(newOccupant);
  }

  addUnit(units: FormArray, unitData: any): void {
    const newUnit = this.fb.group({
      cadastralReference: [unitData.unitKey || '', Validators.required],
      unitType: [unitData.unitType || '', Validators.required],
      floor: [unitData.floor || 0, Validators.required],
      description: [unitData.description || ''],
      quota: [unitData.quota || 0],
      status: ['brouillon'],
      unitAddress: this.fb.group({
        street: [unitData.street || '', Validators.required],
        complement: [unitData.complement || ''],
        postal_code: [unitData.postal_code || '', Validators.required],
        city: [unitData.city || '', Validators.required],
        country: [unitData.country || '', Validators.required],
      }),
      owners: this.fb.array([]), 
    });
    const ownersArray = newUnit.get('owners') as FormArray;
    unitData.owners.forEach((owner: any) => {
      console.log('Adding owner:', owner.name, owner.surname);
      const ownerGroup = this.fb.group({
        name: [owner.name || '', Validators.required],
        surname: [owner.surname || '', Validators.required],
        title: [owner.title || '', Validators.required], 
        quota: [owner.quota || 0], 
        administrator: [owner.administrator || false],
      });
      ownersArray.push(ownerGroup);
    });
  
    units.push(newUnit);
  }
  
}