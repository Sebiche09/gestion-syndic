import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})

export class CondominiumService {
  private apiUrl = 'http://your-backend-api-url';

  private getCivilities = environment.apiUrls.getCivilities;
  private getRemindingMethod = environment.apiUrls.getRemindingMethod;

  constructor(private http: HttpClient) {}

  getCivilityOptions(): Observable<any[]> {
    return this.http.get<any[]>(this.getCivilities);
  }

  getReceivingMethodOptions(): Observable<any[]> {
    return this.http.get<any[]>(this.getRemindingMethod);
  }
}