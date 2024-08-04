import { Routes } from '@angular/router';
import { InvoiceUploadComponent } from './components/invoice-upload/invoice-upload.component';

export const routes: Routes = [
  { path: 'upload', component: InvoiceUploadComponent },
  { path: '', redirectTo: '/upload', pathMatch: 'full' }
];
