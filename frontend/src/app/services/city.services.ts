import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root',
})
export class CityService {
    private baseCitiesUrl = environment.apiUrls.getCitiesBase;
    private username = environment.apiUrls.username;

  constructor(private http: HttpClient) {}

  getCities(countryCode: string): Observable<any[]> {
    const url = `${this.baseCitiesUrl}?country=${countryCode}&username=${this.username}`;
    return this.http.get<any>(url);
  }
}