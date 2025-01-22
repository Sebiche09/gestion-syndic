import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class UnitService {
  constructor(private http: HttpClient) { }

  getUnits(): Observable<any> {
    const url = environment.apiUrls.unit;
    return this.http.get(url);
  }
}
