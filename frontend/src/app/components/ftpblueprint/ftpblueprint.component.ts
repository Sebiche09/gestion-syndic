import { Component, Input } from '@angular/core';
import { FormGroup} from '@angular/forms';
import { FloatLabelModule } from 'primeng/floatlabel';
import { InputTextModule } from 'primeng/inputtext';
import { CardModule } from 'primeng/card';
import { FileUploadComponent } from '../upload/upload.component';
@Component({
  selector: 'app-ftpblueprint',
  standalone: true,
  imports: [
    FloatLabelModule,
    InputTextModule,
    CardModule,
    FileUploadComponent],
  templateUrl: './ftpblueprint.component.html',
  styleUrl: './ftpblueprint.component.scss'
})
export class FtpblueprintComponent {
  @Input() ftpBlueprintForm!: FormGroup;
  
} 
