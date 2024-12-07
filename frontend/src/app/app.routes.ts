import { Routes } from '@angular/router';
import { FileUploadComponent } from './components/upload/upload.component';
import { CondominiumComponent } from './components/condominium/condominium.component';
import { DisplayCondominiumComponent } from './components/display-condominium/display-condominium.component';
import { AccueilComponent } from './components/accueil/accueil.component';

export const routes: Routes = [
  { path: 'upload', component: FileUploadComponent },
  { path: 'createcondominium', component: CondominiumComponent },
  { path: 'displaycondominium', component: DisplayCondominiumComponent },
  {path: 'accueil', component: AccueilComponent},
  { path: '', redirectTo: '/accueil', pathMatch: 'full' }
];
