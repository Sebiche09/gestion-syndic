import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})

export class CondominiumService {
  private uploadUrl = environment.apiUrls.condominiumApi;

  ticketInformation = {
    personalInformation: {
        firstname: '',
        lastname: '',
        age: null
    },
    seatInformation: {
        class: null,
        wagon: null,
        seat: null
    },
    paymentInformation: {
        cardholderName: '',
        cardholderNumber: '',
        date: '',
        cvv: '',
        remember: false
    }
  };

  constructor(private http: HttpClient) { }

  getTicketInformation() {
    return this.ticketInformation;
  }

  setTicketInformation(ticketInformation: { personalInformation: { firstname: string; lastname: string; age: null; }; seatInformation: { class: null; wagon: null; seat: null; }; paymentInformation: { cardholderName: string; cardholderNumber: string; date: string; cvv: string; remember: boolean; }; }) {
    this.ticketInformation = ticketInformation;
  }

  submitForm(formData: any): Observable<any> {
    const headers = new HttpHeaders({ 'Content-Type': 'application/json' });
    
    return this.http.post<any>(this.uploadUrl, formData, { headers });
  }
}
