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
  @Output() lotExtracted = new EventEmitter<any>();
  
  // Méthode pour recevoir les données de l'upload depuis FileUploadComponent
  onFileUploaded(response: any) {
    if (response && response.text && response.text.address) {
      const address = response.text.address;
      const lot = response.text.lot;
      // Émet les données d'adresse vers CondominiumComponent
      this.addressExtracted.emit(address);
      this.lotExtracted.emit(lot);
    }
  }
}
