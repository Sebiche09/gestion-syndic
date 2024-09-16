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
  @Output() textExtracted = new EventEmitter<any>();
  
  // Méthode pour recevoir les données de l'upload depuis FileUploadComponent
  onFileUploaded(response: any) {
    if (response && response.text) {
      const text = response.text;
      this.textExtracted.emit(text);
    }
  }
}
