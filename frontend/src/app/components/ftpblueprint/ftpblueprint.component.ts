import { Component, Input,Output, EventEmitter } from '@angular/core';
import { FormGroup} from '@angular/forms';
import { FloatLabelModule } from 'primeng/floatlabel';
import { InputTextModule } from 'primeng/inputtext';
import { CardModule } from 'primeng/card';
import { FileUploadComponent } from '../upload/upload.component';
import { ButtonModule } from 'primeng/button';
@Component({
  selector: 'app-ftpblueprint',
  standalone: true,
  imports: [
    FloatLabelModule,
    InputTextModule,
    CardModule,
    FileUploadComponent,
    ButtonModule],
  templateUrl: './ftpblueprint.component.html',
  styleUrl: './ftpblueprint.component.scss'
})
export class FtpblueprintComponent {
  @Output() next = new EventEmitter<void>();
  @Output() previous = new EventEmitter<void>();
  
  previousStep() {
    this.previous.emit();
  }
  nextStep() {
      this.next.emit();
  }
} 
