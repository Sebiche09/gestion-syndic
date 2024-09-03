import { Component } from '@angular/core';
import { FileUploadComponent } from '../upload/upload.component';

@Component({
  selector: 'app-cadastre',
  standalone: true,
  imports: [FileUploadComponent],
  templateUrl: './cadastre.component.html',
  styleUrl: './cadastre.component.scss'
})
export class CadastreComponent {

}
