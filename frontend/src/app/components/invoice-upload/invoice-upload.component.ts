import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { UploadService } from '../../services/upload.service';

@Component({
  selector: 'app-invoice-upload',
  standalone: true,
  templateUrl: './invoice-upload.component.html',
  styleUrls: ['./invoice-upload.component.scss'],
  imports: [CommonModule, FormsModule]
})
export class InvoiceUploadComponent {
  file: File | null = null;
  progress: number | null = null;

  constructor(private uploadService: UploadService) { }

  onFileChange(event: any) {
    if (event.target.files.length > 0) {
      this.file = event.target.files[0];
    }
  }

  onSubmit() {
    if (this.file) {
      const metadata = {
        // Add other metadata fields here
      };

      this.uploadService.upload(this.file, metadata).subscribe(
        event => {
          if (typeof event === 'object' && 'status' in event && event.status === 'progress') {
            this.progress = event.message;
          } else {
            console.log('Upload successful:', event);
          }
        },
        error => {
          console.error('Upload error:', error);
        }
      );
    }
  }
}
