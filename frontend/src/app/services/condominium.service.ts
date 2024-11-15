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

  private getCivilities = environment.apiUrls.getCivilities;
  private getDocumentRemindingMethod = environment.apiUrls.getDocumentRemindingMethod;
  private getReminderRemindingMethod = environment.apiUrls.getReminderRemindingMethod;

  constructor(private http: HttpClient) {}

  getCivilityOptions(): Observable<any[]> {
    return this.http.get<any[]>(this.getCivilities);
  }

  getDocumentReceivingMethodOptions(): Observable<any[]> {
    return this.http.get<any[]>(this.getDocumentRemindingMethod);
  }

  getReminderReceivingMethodOptions(): Observable<any[]> {
    return this.http.get<any[]>(this.getReminderRemindingMethod);
  }

  submitCondominium(formData: any): Observable<any> {
    const url = environment.apiUrls.condominiumApi;
    return this.http.post(url, formData);
  }

  createCondominiumForm(): FormGroup {
    return this.fb.group({
      informations: this.fb.group({
        name: [
          '', 
          [Validators.required, Validators.minLength(3), Validators.maxLength(50)],
          [UniqueValidator.checkNameUniqueness(this.uniqueCheckService)],
        ],
        prefix: [
          '',
          [Validators.required, Validators.maxLength(50)],
          [UniqueValidator.checkPrefixUniqueness(this.uniqueCheckService)],
        ],
        description: ['', [Validators.maxLength(500)]],
      }),
      address: this.fb.group({
        street: ['', Validators.required],
        address_complement: [''],
        city: ['', Validators.required],
        postal_code: ['', Validators.required],
        country: ['Belgique', Validators.required],
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

  addUnit(units: FormArray, unitData: any): void {
    const newUnit = this.fb.group({
      cadastralReference: [unitData.unitKey || '', Validators.required],
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
      const ownerGroup = this.fb.group({
        lastName: [owner.last_name || '', Validators.required],
        firstName: [owner.first_name || '', Validators.required],
        title: [owner.title || '', Validators.required],
        address: this.fb.group({
          street: [owner.address?.street || '', Validators.required],
          postal_code: [owner.address?.postal_code || '', Validators.required],
          city: [owner.address?.city || '', Validators.required],
          country: [owner.address?.country || '', Validators.required],
        }),
      });
      ownersArray.push(ownerGroup);
    });

    units.push(newUnit);
  }
}