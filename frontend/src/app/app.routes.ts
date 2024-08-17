import { Routes } from '@angular/router';
import { FileUploadComponent } from './components/invoice-upload/invoice-upload.component';
import { CondominiumComponent } from './components/condominium/condominium.component';

export const routes: Routes = [
  { path: 'upload', component: FileUploadComponent },
  { path: 'condominium', component: CondominiumComponent },
  { path: '', redirectTo: '/condominium', pathMatch: 'full' }
];
