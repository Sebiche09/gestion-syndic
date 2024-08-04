import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse, HttpEventType } from '@angular/common/http';
import { catchError, map } from 'rxjs/operators';
import { Observable, throwError } from 'rxjs';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class UploadService {
  private baseUrl = environment.apiUrls.uploadApi;

  constructor(private http: HttpClient) { }

  upload(file: File, metadata: any): Observable<any> {
    // Create a FormData object for the file and metadata
    const formData: FormData = new FormData();
    formData.append('file', file);

    for (const key in metadata) {
      // Check if the key exists in the metadata object
      if (metadata.hasOwnProperty(key)) {
        formData.append(key, metadata[key]);
      }
    }
    // Create a POST request to the API endpoint
    return this.http.post<any>(this.baseUrl, formData, {
      // report information about the upload progress
      reportProgress: true,
      // request return relative event information
      observe: 'events'
    }).pipe(
      map(event => {
        switch (event.type) {
          case HttpEventType.UploadProgress:
            // Calculate the progress percentage
            const progress = Math.round(100 * event.loaded / event.total!);
            return { status: 'progress', message: progress };
          // Return the response body
          case HttpEventType.Response:
            return event.body;
          default:
            return `Unhandled event: ${event.type}`;
        }
      }),
      catchError(this.handleError)
    );
  }

  private handleError(error: HttpErrorResponse) {
    let errorMessage = '';
    if (error.error instanceof ErrorEvent) {
      // A client-side or network error occurred.
      errorMessage = `An error occurred: ${error.error.message}`;
    } else {
      // The backend returned an unsuccessful response code.
      errorMessage = `Server returned code: ${error.status}, error message is: ${error.message}`;
    }
    return throwError(errorMessage);
  }
}
