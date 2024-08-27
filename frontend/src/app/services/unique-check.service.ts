import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root',
})
export class UniqueCheckService {
  private checkUniqueness = environment.apiUrls.checkUniqueness; // Remplacez par votre API

  constructor(private http: HttpClient) {}

  checkNameUniqueness(name: string): Observable<boolean> {
    return this.http.get<boolean>(`${this.checkUniqueness}?name=${name}`);
  }

  checkPrefixUniqueness(prefix: string): Observable<boolean> {
    return this.http.get<boolean>(`${this.checkUniqueness}?prefix=${prefix}`);
  }
}
