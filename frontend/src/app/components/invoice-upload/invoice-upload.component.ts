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

  // Inject the UploadService
  constructor(private uploadService: UploadService) { }

  // Handle file input changes
  onFileChange(event: any) {
    // Get the file from the event
    if (event.target.files.length > 0) {
      // Set the first file in the list (file)
      this.file = event.target.files[0];
    }
  }

  onSubmit() {
    // Check if a file is selected
    if (this.file) {
      this.uploadService.upload(this.file, {}).subscribe(
        {
          next: (event) => {
            // Handle progress events
            if (typeof event === 'object' && 'status' in event && event.status === 'progress') {
              // Set the progress percentage
              this.progress = event.message;
            } else {
              console.log('Upload successful:', event);
            }
          },
          error: (error) => {
            console.error('Upload error:', error);
          }
        }
      );
    }
  }
}
