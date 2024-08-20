import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})

export class CondominiumService {

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
}