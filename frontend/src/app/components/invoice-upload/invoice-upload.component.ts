import { Component } from '@angular/core';
import { MessageService } from 'primeng/api';
import { FileUploadEvent, FileUploadModule } from 'primeng/fileupload';
import { ToastModule } from 'primeng/toast';
import { CommonModule } from '@angular/common';
import { UploadService } from '../../services/upload.service';

@Component({
  selector: 'app-invoice-upload',
  standalone: true,
  templateUrl: './invoice-upload.component.html',
  styleUrls: ['./invoice-upload.component.scss'],
  imports: [FileUploadModule, ToastModule, CommonModule],
  providers: [MessageService]
})

export class FileUploadComponent {
  uploadedFiles: any[] = [];
  uploadUrl = 'http://localhost:8080/upload';

  constructor(private uploadService: UploadService, private messageService: MessageService) {}

  onUpload(event: FileUploadEvent) {
    this.uploadService.uploadFiles(event.files).subscribe({
      next: (response) => {
        this.uploadedFiles.push(...event.files);
        this.messageService.add({ severity: 'info', summary: 'Files Uploaded', detail: '' });
      },
      error: (error) => {
        this.messageService.add({ severity: 'error', summary: 'Upload Failed', detail: error.message });
      }
    });
  }
}
