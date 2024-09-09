import { Component, Output, EventEmitter  } from '@angular/core';
import { FileUploadComponent } from '../upload/upload.component';
import { CardModule } from 'primeng/card';
@Component({
  selector: 'app-cadastre',
  standalone: true,
  imports: [FileUploadComponent, CardModule],
  templateUrl: './cadastre.component.html',
  styleUrl: './cadastre.component.scss'
})
export class CadastreComponent {
  @Output() addressExtracted = new EventEmitter<any>();
  @Output() unitExtracted = new EventEmitter<any>();
  
  // Méthode pour recevoir les données de l'upload depuis FileUploadComponent
  onFileUploaded(response: any) {
    if (response && response.text && response.text.address && response.text.unit) {
      const address = response.text.address;
      const unit = response.text.unit;
      // Émet les données d'adresse vers CondominiumComponent
      this.addressExtracted.emit(address);
      this.unitExtracted.emit(unit);
    }
  }
}
