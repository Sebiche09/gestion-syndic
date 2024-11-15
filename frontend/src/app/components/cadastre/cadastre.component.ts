import { Component, Output, EventEmitter  } from '@angular/core';
import { FileUploadComponent } from '../upload/upload.component';
import { CardModule } from 'primeng/card';
import { ButtonModule } from 'primeng/button';
@Component({
  selector: 'app-cadastre',
  standalone: true,
  imports: [FileUploadComponent, CardModule, ButtonModule],
  templateUrl: './cadastre.component.html',
  styleUrl: './cadastre.component.scss'
})
export class CadastreComponent {
  @Output() textExtracted = new EventEmitter<any>();
  @Output() next = new EventEmitter<void>();
  @Output() previous = new EventEmitter<void>();
  
  // Méthode pour recevoir les données de l'upload depuis FileUploadComponent
  onFileUploaded(response: any) {
    if (response && response.text) {
      const text = response.text;
      this.textExtracted.emit(text);
    }
  }
  previousStep() {
      this.previous.emit();
  }
  nextStep() {
      this.next.emit();
  }
}
