import { Routes } from '@angular/router';
import { FileUploadComponent } from './components/invoice-upload/invoice-upload.component';

export const routes: Routes = [
  { path: 'upload', component: FileUploadComponent },
  { path: '', redirectTo: '/upload', pathMatch: 'full' }
];
