import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';
@Injectable({
  providedIn: 'root'
})
export class UploadService {
  private uploadUrl = environment.apiUrls.uploadApi;

  constructor(private http: HttpClient) {}

  uploadFiles(files: File[]): Observable<any> {
    const formData = new FormData();
    for (let file of files) {
      formData.append('files', file, file.name);
    }
    return this.http.post(this.uploadUrl, formData);
  }
}