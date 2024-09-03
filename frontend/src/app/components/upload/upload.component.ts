import { Component, Input, signal, computed, inject } from '@angular/core';
import { MessageService, PrimeNGConfig } from 'primeng/api';
import { UploadService } from '../../services/upload.service';
import { FileUploadModule } from 'primeng/fileupload';
import { ButtonModule } from 'primeng/button';
import { CommonModule } from '@angular/common';
import { BadgeModule } from 'primeng/badge';
import { ProgressBarModule } from 'primeng/progressbar';
import { ToastModule } from 'primeng/toast';

@Component({
  selector: 'app-upload',
  standalone: true,
  templateUrl: './upload.component.html',
  styleUrls: ['./upload.component.scss'],
  imports: [
    FileUploadModule, 
    ButtonModule, 
    BadgeModule, 
    ProgressBarModule, 
    ToastModule, 
    CommonModule
  ],
})
export class FileUploadComponent {
  @Input() fileType: string = ''; // Input pour spécifier le type de fichier (ex: 'cadastre', 'facture')
  reponse = "";
  files = signal<File[]>([]);
  totalSize = computed(() => this.files().reduce((sum, file) => sum + file.size, 0));
  totalSizePercent = computed(() => this.totalSize() / 10);

  private config = inject(PrimeNGConfig);
  private messageService = inject(MessageService);
  private uploadService = inject(UploadService); // Injection du service d'upload

  choose(event: Event, callback: () => void) {
    callback();
  }

  onRemoveTemplatingFile(event: Event, file: File, removeFileCallback: (event: Event, index: number) => void, index: number) {
    removeFileCallback(event, index);
    this.files.update(files => files.filter((_, i) => i !== index));
  }

  onClearTemplatingUpload(clear: () => void) {
    clear();
    this.files.set([]);
  }

  onTemplatedUpload() {
    if (this.files().length > 0) {
      console.log(this.files());
      this.files().forEach(file => {
        this.uploadService.uploadFile(file, this.fileType).subscribe({
          next: (response) => {
            console.log(response)
            this.messageService.add({ severity: 'success', summary: 'Success', detail: 'File Uploaded', life: 3000 });
          },
          error: (error) => {
            this.messageService.add({ severity: 'error', summary: 'Error', detail: error.message, life: 3000 });
          }
        });
      });
    }
  }

  onSelectedFiles(event: any) {
    this.files.set(event.currentFiles);
  }

  uploadEvent(callback: () => void) {
    callback();
  }

  formatSize(bytes: number): string {
    const k = 1024;
    const dm = 3;

    const sizes = this.config.translation?.fileSizeTypes || ['Bytes', 'KB', 'MB', 'GB', 'TB'];

    if (bytes === 0) return `0 ${sizes[0]}`;

    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return `${(bytes / Math.pow(k, i)).toFixed(dm)} ${sizes[i]}`;
  }
}
