import { AbstractControl, AsyncValidatorFn, ValidationErrors } from '@angular/forms';
import { Observable, of } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
import { UniqueCheckService } from '../services/unique-check.service';

export class UniqueValidator {
  static checkNameUniqueness(uniqueCheckService: UniqueCheckService): AsyncValidatorFn {
    return (control: AbstractControl): Observable<ValidationErrors | null> => {
      if (!control.value) {
        return of(null);
      }
      return uniqueCheckService.checkNameUniqueness(control.value).pipe(
        map((response: any) => {return response.isTaken ? { nameTaken: true } : null;}),
        catchError((error) => {
          console.error('Error during uniqueness check (name):', error); // Log errors
          return of(null);
        })
      );
    };
  }

  static checkPrefixUniqueness(uniqueCheckService: UniqueCheckService): AsyncValidatorFn {
    return (control: AbstractControl): Observable<ValidationErrors | null> => {
      if (!control.value) {
        return of(null);
      }
      return uniqueCheckService.checkPrefixUniqueness(control.value).pipe(
        map((response: any) => {return response.isTaken ? { prefixTaken: true } : null;}),
        catchError((error) => {
          console.error('Error during uniqueness check (prefix):', error); // Log errors
          return of(null);
        })
      );
    };
  }
}
